package ofdl

import (
	"crypto/md5"
	"fmt"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/duke-git/lancet/v2/maputil"
	"github.com/duke-git/lancet/v2/slice"
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

	ClientIDURI               string
	ClientPrivateKeyURI       string
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
		ClientIDURI:               config.ClientIDURI,
		ClientPrivateKeyURI:       config.ClientPrivateKeyURI,
		OptionalCDRMProjectServer: config.OptionalCDRMProjectServer,
		ClientCacheDir:            path.Join(config.CacheDir, "of_drms"),
		CachePriority:             config.CachePriority,
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

func (dl *OFDl) ScrapeMedias(url string) (results []DownloadableMedia, isSingleURL bool, err error) {
	if url == "" {
		url = gof.OFPostDomain
	}
	if !isOFURL(url) {
		return nil, false, fmt.Errorf("not a valid of url: %s", url)
	}

	if ofurlMatchs(url, reSubscriptions, reHome) {
		results, err = dl.scrapeUser("", "")
		return results, false, err
	}

	//chart
	founds, ok := ofurlFinds(nil, []string{"ID"}, url, reChat)
	if ok {
		results, err = dl.scrapeChat(founds["ID"])
		return results, false, err
	}

	//collections list
	founds, ok = ofurlFinds(nil, []string{"ID"}, url, reCollectionsList)
	if ok {
		results, err = dl.scrapeCollectionsList(founds["ID"])
		return results, false, err
	}

	//post
	founds, ok = ofurlFinds([]string{"ID", "UserName"}, nil, url, reSinglePost)
	if ok {
		post, err := dl.api.GetPost(founds["ID"])
		if err != nil {
			return nil, false, err
		}
		results, err = dl.collectMedias(founds["UserName"], post)
		return results, true, err
	}

	//user
	founds, ok = ofurlFinds([]string{"UserName"}, []string{"MediaType"}, url, reUserWithMediaType)
	if ok {
		results, err = dl.scrapeUser(founds["UserName"], founds["MediaType"])
		return results, false, err
	}

	//bookmarks
	founds, ok = ofurlFinds(nil, []string{"ID", "MediaType"}, url, reBookmarksWithMediaType)
	if ok {
		results, err = dl.scrapeBookmarks(founds["ID"], founds["MediaType"])
		return results, false, err
	}

	results, err = dl.scrapeUser("", "")
	return results, false, err
}

func (dl *OFDl) FetchFileInfo(downloadURL string) (info common.HttpFileInfo, err error) {
	if !isDrmURL(downloadURL) {
		return dl.api.Req().GetFileInfo(downloadURL)
	}
	return dl.drmapi.GetFileInfo(parseDRMURL(downloadURL))
}

func (dl *OFDl) FetchDRMSecrets(downloadURL string, disableCache_ ...bool) (DRMSecrets, error) {
	type cachedSecrets struct {
		DecryptKey string
		Headers    map[string]string
		Time       time.Time
	}

	drmSecretsFromCacheFun := func(secrets cachedSecrets, cookiesCache string) DRMSecrets {
		cookies := secrets.Headers["Cookie"]
		delete(secrets.Headers, "Cookie")
		return DRMSecrets{
			DecryptKey:         secrets.DecryptKey,
			Headers:            secrets.Headers,
			Cookies:            cookies,
			NetscapeCookieFile: cookiesCache,
			TimeStamp:          secrets.Time,
		}
	}

	disableCache := len(disableCache_) > 0 && disableCache_[0]

	var secrets cachedSecrets
	drminfo := parseDRMURL(downloadURL)
	urlmd5 := fmt.Sprintf("%x", md5.Sum([]byte(drminfo.DRM.Manifest.Dash)))
	secretsCachePath := path.Join(dl.cacheDir, "secrets", urlmd5)
	netscapeCookiePath := secretsCachePath + ".cookies"
	if !disableCache && fileutil.IsExist(secretsCachePath) && fileutil.IsExist(netscapeCookiePath) {
		err := common.FileUnmarshal(secretsCachePath, &secrets)
		if err == nil {
			return drmSecretsFromCacheFun(secrets, netscapeCookiePath), nil
		}
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
	common.FileMarshal(secretsCachePath, secrets)
	netscape := common.ConvertCookieToNetscape(secrets.Headers["Cookie"], gof.OFApiDomain)
	err = common.WriteFile(netscapeCookiePath, []byte(netscape))
	if err != nil {
		return DRMSecrets{}, err
	}
	return drmSecretsFromCacheFun(secrets, netscapeCookiePath), nil
}

func (dl *OFDl) Download(dir string, media DownloadableMedia) error {
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
	drminfo := parseDRMURL(media.DownloadURL)
	for k, v := range dl.drmapi.HTTPHeaders(drminfo) {
		if k == "Accept" {
			continue
		}
		args = append(args, "--add-header")
		args = append(args, fmt.Sprintf("%s:%s", k, v))
	}
	args = append(args, drminfo.DRM.Manifest.Dash)
	fmt.Println(args)
	dl.FetchDRMSecrets(media.DownloadURL)
	return nil
}

func (dl *OFDl) scrapeUser(allEmptryOrUserName string, allEmptryOrMediaType string) ([]DownloadableMedia, error) {
	users := []scrapeIdentifier{}
	if allEmptryOrUserName == "" {
		subs, err := dl.api.GetSubscriptions(ofapi.SubscritionTypeActive)
		if err != nil {
			return nil, err
		}
		for _, sub := range subs {
			users = append(users, scrapeIdentifier{
				id:        sub.ID,
				hintTitle: sub.Username,
			})
		}

	} else {
		usr, err := dl.api.GetUserByUsername(allEmptryOrUserName)
		if err != nil {
			return nil, err
		}
		users = []scrapeIdentifier{
			{
				id:        usr.ID,
				hintTitle: allEmptryOrUserName,
			},
		}
	}
	return dl.scrapeUsersByIdentifier(users, allEmptryOrMediaType)
}

func (dl *OFDl) scrapeBookmarks(allEmptryOrID string, allEmptryOrMediaType string) ([]DownloadableMedia, error) {
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

func (dl *OFDl) scrapeCollectionsList(allEmptryOrID string) ([]DownloadableMedia, error) {
	if allEmptryOrID == "" {
		return dl.scrapeUser("", "")
	} else {
		userList, err := dl.api.GetCollectionsListUsers(allEmptryOrID)
		if err != nil {
			return nil, err
		}
		users := []scrapeIdentifier{}
		for _, user := range userList {
			users = append(users, scrapeIdentifier{
				id:        user.ID,
				hintTitle: user.Username,
			})
		}
		return dl.scrapeUsersByIdentifier(users, "")
	}
}

type scrapeIdentifier struct {
	id        any
	hintTitle string
}

func (dl *OFDl) scrapeUsersByIdentifier(users []scrapeIdentifier, allEmptryOrMediaType string) ([]DownloadableMedia, error) {
	funs := []collecFunc{}
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
	return dl.parallelCollecFunc(funs)
}

func (dl *OFDl) scrapeChat(allEmptryOrID string) ([]DownloadableMedia, error) {
	if allEmptryOrID == "" {
		return dl.scrapeUser("", "")
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

type collecFunc func() (hintTitle string, posts []model.Post, error error)

func (dl *OFDl) parallelCollecFunc(funs []collecFunc) ([]DownloadableMedia, error) {
	ch := make(chan struct{}, 5)
	results := []DownloadableMedia{}
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

			var medias []DownloadableMedia
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

func (dl *OFDl) collecMutilMedias(hintTitle string, posts []model.Post) ([]DownloadableMedia, error) {
	if len(posts) == 0 {
		return nil, fmt.Errorf("posts is empty")
	}

	results := []DownloadableMedia{}
	for _, post := range posts {
		medias, e := dl.collectMedias(hintTitle, post)
		if e == nil {
			results = append(results, medias...)
		}
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no media found")
	}
	slice.SortBy(results, func(i, j DownloadableMedia) bool {
		return i.Time.After(j.Time)
	})
	return results, nil
}

func (dl *OFDl) collectMedias(hintTitle string, post model.Post) ([]DownloadableMedia, error) {
	if len(post.Media) == 0 {
		return nil, fmt.Errorf("no media found")
	}
	hintTitle = strings.Trim(hintTitle, ".")

	mediaSet := make(map[int64]DownloadableMedia)
	for i, media := range post.Media {
		if !media.CanView || media.Files == nil {
			continue
		}
		if hintTitle == "" {
			hintTitle = post.FromUser.Username
		}
		dm := DownloadableMedia{
			PostID:  post.ID,
			MediaID: media.ID,
			Type:    media.Type,
			Time:    times(media.CreatedAt, post.CreatedAt, post.PostedAt),
			Title:   strings.TrimLeft(fmt.Sprintf("%s.%d.%d", hintTitle, post.ID, i), "."),
		}

		if media.Files.Drm == nil {
			if media.Files.Full != nil {
				dm.DownloadURL = media.Files.Full.URL
			} else if media.Files.Preview != nil {
				dm.DownloadURL = media.Files.Preview.URL
			}
			dm.IsDrm = false
			if strings.Contains(dm.DownloadURL, "upload") {
				continue
			}
		} else {
			dm.DownloadURL = composeDRMURL(ofdrm.DRMInfo{
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
	slice.SortBy(results, func(i, j DownloadableMedia) bool {
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
