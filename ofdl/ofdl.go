package ofdl

import (
	"context"
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/duke-git/lancet/v2/maputil"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	mrecover "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/yinyajiang/gof"
	"github.com/yinyajiang/gof/common"
	"github.com/yinyajiang/gof/ofapi"
	"github.com/yinyajiang/gof/ofapi/model"
	"github.com/yinyajiang/gof/ofdrm"
)

type DependentTools struct {
	YtDlpPath  string
	FFMpegPath string
	Mp4Decrypt string
}

type Config struct {
	AuthInfo ofapi.AuthInfo
	CacheDir string

	OptionalRulesURI []string
	CachePriority    bool

	WVDURI                    string
	RawWVDIDURI               string
	RawPrivateKeyURI          string
	OptionalCDRMProjectServer []string

	DependentTools DependentTools
}

type OFDl struct {
	api      *ofapi.OFAPI
	drmapi   *ofdrm.OFDRM
	cacheDir string

	dependentTools DependentTools
}

func NewOFDL(config Config) (*OFDl, error) {
	if config.DependentTools.YtDlpPath == "" {
		config.DependentTools.YtDlpPath = "yt-dlp"
	}
	if config.DependentTools.FFMpegPath == "" {
		config.DependentTools.FFMpegPath = "ffmpeg"
	}
	if config.DependentTools.Mp4Decrypt == "" {
		config.DependentTools.Mp4Decrypt = "mp4decrypt"
	}

	api, err := ofapi.NewOFAPI(ofapi.Config{
		AuthInfo:         config.AuthInfo,
		OptionalRulesURI: config.OptionalRulesURI,
		RulesCacheDir:    path.Join(config.CacheDir, "of_apis"),
		CachePriority:    config.CachePriority,
	})
	if err != nil {
		return nil, err
	}
	drmapi, err := ofdrm.NewOFDRM(api.Req(), ofdrm.OFDRMConfig{
		WVDOption: ofdrm.DRMWVDOption{
			ClientIDURI:         config.RawWVDIDURI,
			ClientPrivateKeyURI: config.RawPrivateKeyURI,
			ClientCacheDir:      path.Join(config.CacheDir, "of_drms"),
			CachePriority:       config.CachePriority,
		},
		OptionalCDRMProjectServer: config.OptionalCDRMProjectServer,
	})
	if err != nil {
		return nil, err
	}
	dl := &OFDl{
		api:            api,
		drmapi:         drmapi,
		dependentTools: config.DependentTools,
		cacheDir:       path.Join(config.CacheDir, "of_dls"),
	}
	if err := dl.checkDependentTools(); err != nil {
		return nil, err
	}
	return dl, nil
}

func (dl *OFDl) OFAPI() *ofapi.OFAPI {
	return dl.api
}

func (dl *OFDl) OFDRM() *ofdrm.OFDRM {
	return dl.drmapi
}

func (dl *OFDl) Serve(ctx context.Context, addr string, debug bool) {
	app := fiber.New(
		fiber.Config{
			DisableStartupMessage: true,
			Immutable:             true,
		},
	)
	go func() {
		<-ctx.Done()
		app.Shutdown()
	}()
	if debug {
		app.Use(logger.New())
	} else {
		app.Use(mrecover.New())
	}
	app.Hooks().OnShutdown(func() error {
		log.Println("shutdown, close")
		return nil
	})
	dl.AddFiberRoutes(app)
	app.Listen(addr)
}

func (dl *OFDl) AddFiberRoutes(router fiber.Router) {
	addOFDlFiberRoutes(dl, router)
}

func (dl *OFDl) ExtractMedias(url string, disableCache_ ...bool) (ret ExtractResult, err error) {
	if url == "" {
		url = gof.OFPostDomain
	}
	if !isOFURL(url) {
		return ExtractResult{}, fmt.Errorf("not a valid of url: %s", url)
	}
	disableCache := len(disableCache_) > 0 && disableCache_[0]

	type cachedMediaInfo struct {
		Medias      []MediaInfo
		IsSingleURL bool
		Time        time.Time
	}
	cached := cachedMediaInfo{}
	if !disableCache && dl.cacheUnmarshal("medias", url, &cached) && cached.Time.After(time.Now().Add(-time.Hour*24)) {
		return ExtractResult{
			Medias:      cached.Medias,
			IsSingleURL: cached.IsSingleURL,
			IsFromCache: true,
		}, nil
	}

	defer func() {
		if err == nil {
			cached.Medias = ret.Medias
			cached.IsSingleURL = ret.IsSingleURL
			cached.Time = time.Now()
			dl.cacheMarshal("medias", url, cached)
		} else {
			dl.cacheDelete("medias", url)
		}
	}()

	if ofurlMatchs(url, reSubscriptions, reHome) {
		ret.Medias, err = dl.extractUser("", "")
		return
	}

	//chart
	founds, ok := ofurlFinds(nil, []string{"ID"}, url, reChat)
	if ok {
		ret.Medias, err = dl.extractChat(founds["ID"])
		return
	}

	//collections list
	founds, ok = ofurlFinds(nil, []string{"ID"}, url, reCollectionsList)
	if ok {
		ret.Medias, err = dl.extractCollectionsList(founds["ID"])
		return
	}

	//post
	founds, ok = ofurlFinds([]string{"ID", "UserName"}, nil, url, reSinglePost)
	if ok {
		post, e := dl.api.GetPost(founds["ID"])
		if e != nil {
			return ExtractResult{}, e
		}
		ret.Medias, err = dl.collectMedias(founds["UserName"], post)
		ret.IsSingleURL = true
		return
	}

	//user
	founds, ok = ofurlFinds([]string{"UserName"}, []string{"MediaType"}, url, reUserWithMediaType)
	if ok {
		ret.Medias, err = dl.extractUser(founds["UserName"], founds["MediaType"])
		return
	}

	//bookmarks
	founds, ok = ofurlFinds(nil, []string{"ID", "MediaType"}, url, reBookmarksWithMediaType)
	if ok {
		ret.Medias, err = dl.extractBookmarks(founds["ID"], founds["MediaType"])
		return
	}

	ret.Medias, err = dl.extractUser("", "")
	return
}

func (dl *OFDl) FetchFileInfo(mediaURI string) (info common.HttpFileInfo, err error) {
	if !isDrmURL(mediaURI) {
		return dl.api.Req().GetFileInfo(mediaURI)
	}
	return dl.drmapi.GetFileInfo(parseDrmUri(mediaURI))
}

func (dl *OFDl) GetNonDRMSecrets() (NonDRMSecrets, error) {
	return NonDRMSecrets{
		Headers: map[string]string{
			"User-Agent": dl.api.UserAgent(),
		},
	}, nil
}

func (dl *OFDl) FetchDRMSecrets(mediaURI string, disableCache_ ...bool) (DRMSecrets, error) {
	type cachedSecrets struct {
		DecryptKey string
		Headers    map[string]string
		Time       time.Time
	}
	drminfo := parseDrmUri(mediaURI)

	drmSecretsFromCacheFun := func(secrets cachedSecrets) DRMSecrets {
		cookieString := secrets.Headers["Cookie"]
		delete(secrets.Headers, "Cookie")
		cookies := map[string]string{}
		for _, pairs := range strings.Split(cookieString, ";") {
			kv := strings.SplitN(pairs, "=", 2)
			if len(kv) == 2 {
				cookies[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
			}
		}
		return DRMSecrets{
			MPDURL:        drminfo.DRM.Manifest.Dash,
			DecryptKey:    secrets.DecryptKey,
			Headers:       secrets.Headers,
			Cookies:       cookies,
			CookiesString: cookieString,
		}
	}

	disableCache := len(disableCache_) > 0 && disableCache_[0]

	var secrets cachedSecrets
	if !disableCache && dl.cacheUnmarshal("secrets", drminfo.DRM.Manifest.Dash, &secrets) {
		return drmSecretsFromCacheFun(secrets), nil
	}

	decript, err := dl.drmapi.GetDecryptedKeyAuto(drminfo)
	if err != nil {
		return DRMSecrets{}, err
	}
	headers := dl.drmapi.HTTPHeaders(drminfo)
	secrets = cachedSecrets{
		DecryptKey: decript,
		Headers:    headers,
		Time:       time.Now(),
	}
	dl.cacheMarshal("secrets", drminfo.DRM.Manifest.Dash, secrets)
	return drmSecretsFromCacheFun(secrets), nil
}

func (dl *OFDl) Download(dir string, media MediaInfo) error {
	if !media.IsDrm {
		return nil
	}
	args := []string{
		"--no-part",
		"--restrict-filenames",
		"-o",
		fmt.Sprintf(`%s/%%(title)s.%%(ext)s`, dir),
		"--format",
		"bestvideo[ext=mp4]+bestaudio[ext=m4a]/best[ext=mp4]/best[ext=m4a]",
	}
	drminfo := parseDrmUri(media.MediaURI)
	for k, v := range dl.drmapi.HTTPHeaders(drminfo) {
		if k == "Accept" {
			continue
		}
		args = append(args, "--add-header")
		args = append(args, fmt.Sprintf("%s:%s", k, v))
	}
	args = append(args, drminfo.DRM.Manifest.Dash)
	fmt.Println(args)
	dl.FetchDRMSecrets(media.MediaURI)
	return nil
}

func (dl *OFDl) extractUser(allEmptryOrUserName string, allEmptryOrMediaType string) ([]MediaInfo, error) {
	users := []extractIdentifier{}
	if allEmptryOrUserName == "" {
		subs, err := dl.api.GetSubscriptions(ofapi.SubscritionTypeActive)
		if err != nil {
			return nil, err
		}
		for _, sub := range subs {
			users = append(users, extractIdentifier{
				id:        sub.ID,
				hintTitle: sub.Username,
			})
		}

	} else {
		usr, err := dl.api.GetUserByUsername(allEmptryOrUserName)
		if err != nil {
			return nil, err
		}
		users = []extractIdentifier{
			{
				id:        usr.ID,
				hintTitle: allEmptryOrUserName,
			},
		}
	}
	return dl.extractUsersByIdentifier(users, allEmptryOrMediaType)
}

func (dl *OFDl) extractBookmarks(allEmptryOrID string, allEmptryOrMediaType string) ([]MediaInfo, error) {
	if allEmptryOrID == "" {
		bookmarks, err := dl.api.GetAllBookmarkes(ofapi.BookmarkMedia(allEmptryOrMediaType))
		if err != nil {
			return nil, err
		}
		return dl.collecMutilMedias("bookmarks."+allEmptryOrMediaType, bookmarks)
	}
	bookmarks, err := dl.api.GetBookmark(allEmptryOrID, ofapi.BookmarkMedia(allEmptryOrMediaType))
	if err != nil {
		return nil, err
	}
	return dl.collecMutilMedias("bookmark."+allEmptryOrMediaType, bookmarks)
}

func (dl *OFDl) extractCollectionsList(allEmptryOrID string) ([]MediaInfo, error) {
	if allEmptryOrID == "" {
		return dl.extractUser("", "")
	} else {
		userList, err := dl.api.GetCollectionsListUsers(allEmptryOrID)
		if err != nil {
			return nil, err
		}
		users := []extractIdentifier{}
		for _, user := range userList {
			users = append(users, extractIdentifier{
				id:        user.ID,
				hintTitle: user.Username,
			})
		}
		return dl.extractUsersByIdentifier(users, "")
	}
}

type extractIdentifier struct {
	id        any
	hintTitle string
}

func (dl *OFDl) extractUsersByIdentifier(users []extractIdentifier, allEmptryOrMediaType string) ([]MediaInfo, error) {
	funs := []extractFunc{}
	for _, user := range users {
		funs = append(funs, func() (string, []model.Post, error) {
			userID, e := toInt64(user.id)
			if e != nil {
				return "", nil, e
			}
			posts, e := dl.api.GetUserMedias(userID, ofapi.UserMedias(allEmptryOrMediaType))
			return user.hintTitle, posts, e
		})
	}
	return dl.parallelExtractFunc(funs)
}

func (dl *OFDl) extractChat(allEmptryOrID string) ([]MediaInfo, error) {
	if allEmptryOrID == "" {
		return dl.extractUser("", "")
	} else {
		chatID, err := toInt64(allEmptryOrID)
		if err != nil {
			return nil, err
		}
		posts, err := dl.api.GetChatMessages(chatID)
		if err != nil {
			return nil, err
		}
		return dl.collecMutilMedias("", posts)
	}
}

type extractFunc func() (hintTitle string, posts []model.Post, error error)

func (dl *OFDl) parallelExtractFunc(funs []extractFunc) ([]MediaInfo, error) {
	ch := make(chan struct{}, 5)
	results := []MediaInfo{}
	var firstErr error
	var lock sync.Mutex
	var wg sync.WaitGroup
	for _, fun := range funs {
		ch <- struct{}{}
		wg.Add(1)
		go func() {
			defer func() {
				<-ch
				wg.Done()
			}()
			hintTitle, posts, err := fun()
			lock.Lock()
			defer lock.Unlock()

			var medias []MediaInfo
			if err == nil {
				medias, err = dl.collecMutilMedias(hintTitle, posts)
			}
			if err != nil {
				if firstErr == nil {
					firstErr = err
				}
			} else {
				results = append(results, medias...)
			}
		}()
	}
	wg.Wait()
	if len(results) != 0 {
		return results, nil
	}
	return results, firstErr
}

func (dl *OFDl) collecMutilMedias(hintTitle string, posts []model.Post) ([]MediaInfo, error) {
	if len(posts) == 0 {
		return nil, fmt.Errorf("posts is empty")
	}

	results := []MediaInfo{}
	for _, post := range posts {
		medias, e := dl.collectMedias(hintTitle, post)
		if e == nil {
			results = append(results, medias...)
		}
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no media found")
	}
	slice.SortBy(results, func(i, j MediaInfo) bool {
		return i.Time.After(j.Time)
	})
	return results, nil
}

func (dl *OFDl) collectMedias(hintTitle string, post model.Post) ([]MediaInfo, error) {
	if len(post.Media) == 0 {
		return nil, fmt.Errorf("no media found")
	}
	hintTitle = strings.Trim(hintTitle, ".")

	mediaSet := make(map[int64]MediaInfo)
	for i, media := range post.Media {
		if !media.CanView || media.Files == nil {
			continue
		}
		if hintTitle == "" {
			hintTitle = post.FromUser.Username
		}
		dm := MediaInfo{
			PostID:  post.ID,
			MediaID: media.ID,
			Type:    media.Type,
			Time:    times(media.CreatedAt, post.CreatedAt, post.PostedAt),
			Title:   strings.TrimLeft(fmt.Sprintf("%s.%d.%d", hintTitle, post.ID, i), "."),
		}

		if media.Files.Drm == nil {
			if media.Files.Full != nil {
				dm.MediaURI = media.Files.Full.URL
			} else if media.Files.Preview != nil {
				dm.MediaURI = media.Files.Preview.URL
			}
			dm.IsDrm = false
			if strings.Contains(dm.MediaURI, "upload") {
				continue
			}
		} else {
			dm.MediaURI = composeDrmUri(ofdrm.DRMInfo{
				DRM:     *media.Files.Drm,
				MediaID: media.ID,
				PostID:  post.ID,
			})
			dm.IsDrm = true
		}
		mediaSet[media.ID] = dm
	}
	if len(mediaSet) == 0 {
		return nil, fmt.Errorf("no can view media found")
	}
	results := maputil.Values(mediaSet)
	slice.SortBy(results, func(i, j MediaInfo) bool {
		return i.Time.After(j.Time)
	})
	return results, nil
}

func (dl *OFDl) checkDependentTools() error {
	addExe := func(path *string) {
		//is windows
		if strings.EqualFold(runtime.GOOS, "windows") {
			if !strings.HasSuffix(*path, ".exe") {
				*path = *path + ".exe"
			}
		} else {
			if strings.HasSuffix(*path, ".exe") {
				*path = strings.TrimSuffix(*path, ".exe")
			}
		}
	}
	addExe(&dl.dependentTools.YtDlpPath)
	addExe(&dl.dependentTools.FFMpegPath)
	addExe(&dl.dependentTools.Mp4Decrypt)

	if !fileutil.IsExist(dl.dependentTools.YtDlpPath) {
		return fmt.Errorf("ytdlp not found")
	}
	if !fileutil.IsExist(dl.dependentTools.FFMpegPath) {
		return fmt.Errorf("ffmpeg not found")
	}
	if !fileutil.IsExist(dl.dependentTools.Mp4Decrypt) {
		return fmt.Errorf("mp4decrypt not found")
	}
	return nil
}

func (dl *OFDl) cacheMarshal(storage, key string, v any) error {
	keymd5 := fmt.Sprintf("%x", md5.Sum([]byte(key)))
	cachePath := path.Join(dl.cacheDir, storage, keymd5)
	return common.FileMarshal(cachePath, v)
}

func (dl *OFDl) cacheDelete(storage, key string) {
	keymd5 := fmt.Sprintf("%x", md5.Sum([]byte(key)))
	cachePath := path.Join(dl.cacheDir, storage, keymd5)
	os.Remove(cachePath)
}

func (dl *OFDl) cacheUnmarshal(storage, key string, pv any) bool {
	keymd5 := fmt.Sprintf("%x", md5.Sum([]byte(key)))
	cachePath := path.Join(dl.cacheDir, storage, keymd5)
	if !fileutil.IsExist(cachePath) {
		return false
	}
	return common.FileUnmarshal(cachePath, pv) == nil
}
