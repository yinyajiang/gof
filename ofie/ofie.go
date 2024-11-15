package ofie

import (
	"context"
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"path"
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

type Config struct {
	OFAuthConfig                           ofapi.OFAuthConfig
	OFDRMConfig                            ofdrm.OFDRMConfig
	CacheDir                               string
	Debug                                  bool
	CacheSeconds                           int
	PreferMediaTypeWhenExtractAllMediasURL string //video,photo,all
}

type OFIE struct {
	api                                    *ofapi.OFAPI
	drmapi                                 *ofdrm.OFDRM
	cacheDir                               string
	cacheSeconds                           int
	preferMediaTypeWhenExtractAllMediasURL string
}

func NewOFIE(config Config) (*OFIE, error) {
	if config.OFAuthConfig.RulesCacheDir == "" {
		config.OFAuthConfig.RulesCacheDir = path.Join(config.CacheDir, "of_apis")
	}
	if config.OFDRMConfig.WVDOption.ClientCacheDir == "" {
		config.OFDRMConfig.WVDOption.ClientCacheDir = path.Join(config.CacheDir, "of_drms")
	}
	if config.Debug {
		gof.SetDebug(true)
	}
	api := ofapi.NewOFAPI(config.OFAuthConfig)
	drmapi, err := ofdrm.NewOFDRM(api.Req(), config.OFDRMConfig)
	if err != nil {
		return nil, err
	}

	if strings.Contains(config.PreferMediaTypeWhenExtractAllMediasURL, "video") {
		config.PreferMediaTypeWhenExtractAllMediasURL = string(ofapi.UserVideos)
	} else if strings.Contains(config.PreferMediaTypeWhenExtractAllMediasURL, "photo") {
		config.PreferMediaTypeWhenExtractAllMediasURL = string(ofapi.UserPhotos)
	} else {
		config.PreferMediaTypeWhenExtractAllMediasURL = string(ofapi.UserAll)
	}

	ie := &OFIE{
		api:                                    api,
		drmapi:                                 drmapi,
		cacheDir:                               path.Join(config.CacheDir, "of_ies"),
		cacheSeconds:                           config.CacheSeconds,
		preferMediaTypeWhenExtractAllMediasURL: config.PreferMediaTypeWhenExtractAllMediasURL,
	}
	return ie, nil
}

func (ie *OFIE) OFAPI() *ofapi.OFAPI {
	return ie.api
}

func (ie *OFIE) OFDRM() *ofdrm.OFDRM {
	return ie.drmapi
}

func (ie *OFIE) Serve(ctx context.Context, addr string) {
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
	if gof.IsDebug() {
		app.Use(logger.New())
	} else {
		app.Use(mrecover.New())
	}
	app.Hooks().OnShutdown(func() error {
		log.Println("shutdown, close")
		return nil
	})
	app.Hooks().OnListen(func(c fiber.ListenData) error {
		log.Println("listen on => ", addr)
		return nil
	})
	ie.AddFiberRoutes(app)
	err := app.Listen(addr)
	if err != nil {
		log.Println("listen error => ", err)
	}
}

func (ie *OFIE) AddFiberRoutes(router fiber.Router) {
	addOFIEFiberRoutes(ie, router)
}

func (ie *OFIE) ExtractMedias(url string, option ExtractOption) (ret ExtractResult, err error) {
	if url == "" {
		url = gof.OFPostDomain
	}
	if !isOFURL(url) {
		return ExtractResult{}, fmt.Errorf("not a valid of url: %s", url)
	}
	defer collectTitle(&ret)

	type cachedMediaInfo struct {
		Medias      []MediaInfo
		IsSingleURL bool
		Time        time.Time
		Title       string
	}
	cached := cachedMediaInfo{}
	if !option.DisableCache && ie.cacheUnmarshal("medias", url, &cached) && (ie.cacheSeconds < 0 || cached.Time.After(time.Now().Add(-time.Duration(ie.cacheSeconds)*time.Second))) {
		return ExtractResult{
			Medias:      cached.Medias,
			IsSingleURL: cached.IsSingleURL,
			IsFromCache: true,
			Title:       cached.Title,
		}, nil
	}

	defer func() {
		if err == nil {
			cached.Medias = ret.Medias
			cached.IsSingleURL = ret.IsSingleURL
			cached.Time = time.Now()
			cached.Title = strings.Split(ret.Medias[0].Title, titleSeparator)[0]
			ie.cacheMarshal("medias", url, cached)
		} else {
			ie.cacheDelete("medias", url)
			err = ie.convertApiError(err)
		}
	}()

	if ofurlMatchs(url, reSubscriptions, reHome) {
		ret.Medias, err = ie.extractUser("", "")
		return
	}

	//chart
	founds, ok := ofurlFinds(nil, []string{"ID"}, url, reChat)
	if ok {
		ret.Medias, err = ie.extractChat(founds["ID"])
		return
	}

	//collections list
	founds, ok = ofurlFinds(nil, []string{"ID"}, url, reCollectionsList)
	if ok {
		ret.Medias, err = ie.extractCollectionsList(founds["ID"])
		return
	}

	//post
	founds, ok = ofurlFinds([]string{"ID", "UserName"}, nil, url, reSinglePost)
	if ok {
		post, e := ie.api.GetPost(founds["ID"])
		if e != nil {
			return ExtractResult{}, e
		}
		ret.Medias, err = ie.collectMedias(founds["UserName"], post)
		ret.IsSingleURL = true
		return
	}

	//user
	founds, ok = ofurlFinds([]string{"UserName"}, []string{"MediaType"}, url, reUserWithMediaType)
	if ok {
		ret.Medias, err = ie.extractUser(founds["UserName"], founds["MediaType"])
		return
	}

	//bookmarks
	founds, ok = ofurlFinds(nil, []string{"ID", "MediaType"}, url, reBookmarksWithMediaType)
	if ok {
		ret.Medias, err = ie.extractBookmarks(founds["ID"], founds["MediaType"])
		return
	}

	ret.Medias, err = ie.extractUser("", "")
	return
}

func (ie *OFIE) FetchFileInfo(mediaURI string) (info common.HttpFileInfo, err error) {
	defer func() {
		err = ie.convertApiError(err)
	}()

	if !isDrmURL(mediaURI) {
		return ie.api.Req().GetFileInfo(mediaURI)
	}
	return ie.drmapi.GetFileInfo(parseDrmUri(mediaURI))
}

func (ie *OFIE) GetNonDRMSecrets() (NonDRMSecrets, error) {
	return NonDRMSecrets{
		Headers: map[string]string{
			"User-Agent": ie.api.UserAgent(),
		},
	}, nil
}

func (ie *OFIE) FetchDRMSecrets(mediaURI string, disableCache_ ...bool) (ret DRMSecrets, err error) {
	defer func() {
		if err != nil {
			err = ie.convertApiError(err)
		}
	}()

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
	if !disableCache && ie.cacheUnmarshal("secrets", drminfo.DRM.Manifest.Dash, &secrets) {
		return drmSecretsFromCacheFun(secrets), nil
	}

	decript, err := ie.drmapi.GetDecryptedKeyAuto(drminfo)
	if err != nil {
		return DRMSecrets{}, err
	}
	headers := ie.drmapi.HTTPHeaders(drminfo)
	secrets = cachedSecrets{
		DecryptKey: decript,
		Headers:    headers,
		Time:       time.Now(),
	}
	ie.cacheMarshal("secrets", drminfo.DRM.Manifest.Dash, secrets)
	return drmSecretsFromCacheFun(secrets), nil
}

func (ie *OFIE) extractUser(allEmptryOrUserName string, allEmptryOrMediaType string) ([]MediaInfo, error) {
	users := []extractIdentifier{}
	if allEmptryOrUserName == "" {
		subs, err := ie.api.GetSubscriptions(ofapi.SubscritionTypeActive)
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
		usr, err := ie.api.GetUserByUsername(allEmptryOrUserName)
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
	return ie.extractUsersByIdentifier(users, allEmptryOrMediaType)
}

func (ie *OFIE) extractBookmarks(allEmptryOrID string, allEmptryOrMediaType string) ([]MediaInfo, error) {
	if allEmptryOrMediaType == "" && ie.preferMediaTypeWhenExtractAllMediasURL != "" {
		result, err := ie._extractBookmarks(allEmptryOrID, ie.preferMediaTypeWhenExtractAllMediasURL)
		if err == nil {
			return result, nil
		}
	}
	return ie._extractBookmarks(allEmptryOrID, allEmptryOrMediaType)
}

func (ie *OFIE) _extractBookmarks(allEmptryOrID string, allEmptryOrMediaType string) ([]MediaInfo, error) {
	if allEmptryOrID == "" {
		bookmarks, err := ie.api.GetAllBookmarkes(ofapi.BookmarkMedia(allEmptryOrMediaType))
		if err != nil {
			return nil, err
		}
		return ie.collecMutilMedias("bookmarks"+titleSeparator+allEmptryOrMediaType, bookmarks)
	}
	bookmarks, err := ie.api.GetBookmark(allEmptryOrID, ofapi.BookmarkMedia(allEmptryOrMediaType))
	if err != nil {
		return nil, err
	}
	return ie.collecMutilMedias("bookmark"+titleSeparator+allEmptryOrMediaType, bookmarks)
}

func (ie *OFIE) extractCollectionsList(allEmptryOrID string) ([]MediaInfo, error) {
	if allEmptryOrID == "" {
		return ie.extractUser("", "")
	} else {
		userList, err := ie.api.GetCollectionsListUsers(allEmptryOrID)
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
		return ie.extractUsersByIdentifier(users, "")
	}
}

type extractIdentifier struct {
	id        any
	hintTitle string
}

func (ie *OFIE) extractUsersByIdentifier(users []extractIdentifier, allEmptryOrMediaType string) ([]MediaInfo, error) {
	if allEmptryOrMediaType == "" && ie.preferMediaTypeWhenExtractAllMediasURL != "" {
		result, err := ie._extractUsersByIdentifier(users, ie.preferMediaTypeWhenExtractAllMediasURL)
		if err == nil {
			return result, nil
		}
	}
	return ie._extractUsersByIdentifier(users, allEmptryOrMediaType)
}

func (ie *OFIE) _extractUsersByIdentifier(users []extractIdentifier, allEmptryOrMediaType string) ([]MediaInfo, error) {
	funs := []extractFunc{}
	for _, user := range users {
		funs = append(funs, func() (string, []model.Post, error) {
			userID, e := toInt64(user.id)
			if e != nil {
				return "", nil, e
			}
			posts, e := ie.api.GetUserMedias(userID, ofapi.UserMedias(allEmptryOrMediaType))
			return user.hintTitle, posts, e
		})
	}
	return ie.parallelExtractFunc(funs)
}

func (ie *OFIE) extractChat(allEmptryOrID string) ([]MediaInfo, error) {
	if allEmptryOrID == "" {
		return ie.extractUser("", "")
	} else {
		chatID, err := toInt64(allEmptryOrID)
		if err != nil {
			return nil, err
		}
		posts, err := ie.api.GetChatMessages(chatID)
		if err != nil {
			return nil, err
		}
		return ie.collecMutilMedias("", posts)
	}
}

type extractFunc func() (hintTitle string, posts []model.Post, error error)

func (ie *OFIE) parallelExtractFunc(funs []extractFunc) ([]MediaInfo, error) {
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
				medias, err = ie.collecMutilMedias(hintTitle, posts)
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

func (ie *OFIE) collecMutilMedias(hintTitle string, posts []model.Post) ([]MediaInfo, error) {
	if len(posts) == 0 {
		return nil, fmt.Errorf("posts is empty")
	}

	results := []MediaInfo{}
	for _, post := range posts {
		medias, e := ie.collectMedias(hintTitle, post)
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

func (ie *OFIE) collectMedias(hintTitle string, post model.Post) ([]MediaInfo, error) {
	if len(post.Media) == 0 {
		return nil, fmt.Errorf("no media found")
	}
	hintTitle = strings.Trim(hintTitle, titleSeparator)

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
			Title:   strings.TrimLeft(fmt.Sprintf("%s%s%d%s%d", hintTitle, titleSeparator, post.ID, titleSeparator, i), titleSeparator),
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

func (ie *OFIE) cacheMarshal(storage, key string, v any) error {
	keymd5 := fmt.Sprintf("%x", md5.Sum([]byte(key)))
	cachePath := path.Join(ie.cacheDir, storage, keymd5)
	return common.FileMarshal(cachePath, v)
}

func (ie *OFIE) cacheDelete(storage, key string) {
	keymd5 := fmt.Sprintf("%x", md5.Sum([]byte(key)))
	cachePath := path.Join(ie.cacheDir, storage, keymd5)
	os.Remove(cachePath)
}

func (ie *OFIE) cacheUnmarshal(storage, key string, pv any) bool {
	keymd5 := fmt.Sprintf("%x", md5.Sum([]byte(key)))
	cachePath := path.Join(ie.cacheDir, storage, keymd5)
	if !fileutil.IsExist(cachePath) {
		return false
	}
	return common.FileUnmarshal(cachePath, pv) == nil
}

func (ie *OFIE) convertApiError(err error) error {
	if err == nil {
		return nil
	}
	e := ie.api.CheckAuth()
	if e != nil {
		return fmt.Errorf("try to Sign in again : %w", e)
	}
	return err
}

const titleSeparator = "_"