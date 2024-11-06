package ofdl

import (
	"fmt"
	"path"
	"runtime"
	"strings"

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
	AuthInfo gof.AuthInfo
	CacheDir string

	OptionalRulesURL []string
	CachePriority    bool

	ClientIDURL               string
	ClientPrivateKeyURL       string
	OptionalCDRMProjectServer []string

	DependentTools DependentTools
}

type OFDl struct {
	api    *ofapi.OFAPI
	drmapi *ofdrm.OFDRM

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
		OptionalRulesURL: config.OptionalRulesURL,
		RulesCacheDir:    path.Join(config.CacheDir, "of_rules"),
		CachePriority:    config.CachePriority,
	})
	if err != nil {
		return nil, err
	}
	drmapi, err := ofdrm.NewOFDRM(api.Req(), ofdrm.OFDRMConfig{
		ClientIDURL:               config.ClientIDURL,
		ClientPrivateKeyURL:       config.ClientPrivateKeyURL,
		OptionalCDRMProjectServer: config.OptionalCDRMProjectServer,
		ClientCacheDir:            path.Join(config.CacheDir, "of_client"),
		CachePriority:             config.CachePriority,
	})
	if err != nil {
		return nil, err
	}
	dl := &OFDl{
		api:            api,
		drmapi:         drmapi,
		dependentTools: config.DependentTools,
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

	isScrapeAll := isOFHomeURL(url) || ofurlMatchs(url, reSubscriptions)
	if isScrapeAll {
		results, err = dl.scrapeAll()
		return results, false, err
	}

	//chart
	chartID, ok := ofurlFinds(url, "ID", reSingleChat)
	if ok {
		chats := []scrapeIdentifier{
			{
				id:       chartID,
				hintName: "",
			},
		}
		results, err = dl.scrapeChats(chats)
		return results, false, err
	}

	//charts
	if ok = ofurlMatchs(url, reChats); ok {
		subs, err := dl.api.GetSubscriptions(ofapi.SubscritionTypeActive)
		if err != nil {
			return nil, false, err
		}
		chats := []scrapeIdentifier{}
		for _, sub := range subs {
			chats = append(chats, scrapeIdentifier{
				id:       sub.ID,
				hintName: sub.Username,
			})
		}
		results, err = dl.scrapeChats(chats)
		return results, false, err
	}

	//collections list
	listID, ok := ofurlFinds(url, "ID", reUserList)
	if ok {
		userList, err := dl.api.GetCollectionsListUsers(listID)
		if err != nil {
			return nil, false, err
		}
		users := []scrapeIdentifier{}
		for _, user := range userList {
			users = append(users, scrapeIdentifier{
				id:       user.ID,
				hintName: user.Username,
			})
		}
		results, err = dl.scrapeUsers(users, ofapi.UserMediasAll)
		return results, false, err
	}

	//post
	postID, userName, ok := ofurlFinds2(url, "PostID", "UserName", reSinglePost)
	if ok {
		post, err := dl.api.GetPost(postID)
		if err != nil {
			return nil, false, err
		}
		results, err = dl.collectMedias(userName, post)
		return results, true, err
	}

	//user, user media
	userName, ok = ofurlFinds(url, "UserName", reUser, reUserMedia)
	if ok {
		results, err = dl.scrapeUserByName(userName, ofapi.UserMediasAll)
		return results, false, err
	}

	//user photos
	userName, ok = ofurlFinds(url, "UserName", reUserPhotos)
	if ok {
		results, err = dl.scrapeUserByName(userName, ofapi.UserMediasPhoto)
		return results, false, err
	}

	//user videos
	userName, ok = ofurlFinds(url, "UserName", reUserVideos)
	if ok {
		results, err = dl.scrapeUserByName(userName, ofapi.UserMediasVideo)
		return results, false, err
	}

	//all bookmarks
	if ok = ofurlMatchs(url, reAllBookmarks); ok {
		bookmarks, err := dl.api.GetAllBookmarkes()
		if err != nil {
			return nil, false, err
		}
		results, err = collecMutilMedias(dl, "bookmarks", bookmarks)
		return results, false, err
	}

	//single bookmark
	bookmarkID, ok := ofurlFinds(url, "ID", reSingleBookmark)
	if ok {
		bookmarks, err := dl.api.GetBookmark(bookmarkID)
		if err != nil {
			return nil, false, err
		}
		results, err = collecMutilMedias(dl, "bookmark", bookmarks)
		return results, false, err
	}

	results, err = dl.scrapeAll()
	return results, false, err
}

func (dl *OFDl) FetchFileInfo(downloadURL string) (info common.HttpFileInfo, err error) {
	if !isDrmURL(downloadURL) {
		return dl.api.Req().GetFileInfo(downloadURL)
	}
	return dl.drmapi.GetFileInfo(parseDRMURL(downloadURL))
}

func (dl *OFDl) FetchDRMDecrypt(downloadURL string) (string, error) {
	return dl.drmapi.GetDecryptedKeyAuto(parseDRMURL(downloadURL))
}

func (dl *OFDl) scrapeAll() ([]DownloadableMedia, error) {
	subs, err := dl.api.GetSubscriptions(ofapi.SubscritionTypeActive)
	if err != nil {
		return nil, err
	}
	users := []scrapeIdentifier{}
	for _, sub := range subs {
		users = append(users, scrapeIdentifier{
			id:       sub.ID,
			hintName: sub.Username,
		})
	}
	return dl.scrapeUsers(users, ofapi.UserMediasAll)
}

func (dl *OFDl) scrapeUserByName(userName string, userMedias ofapi.UserMedias) ([]DownloadableMedia, error) {
	usr, err := dl.api.GetUserByUsername(userName)
	if err != nil {
		return nil, err
	}
	return dl.scrapeUsers([]scrapeIdentifier{
		{
			id:       usr.ID,
			hintName: userName,
		},
	}, userMedias)
}

func (dl *OFDl) scrapeUsers(users []scrapeIdentifier, userMedias ofapi.UserMedias) ([]DownloadableMedia, error) {
	funs := []collecFunc{}
	for _, user := range users {
		funs = append(funs, func() (string, []model.Post, error) {
			userID, e := toInt64(user.id)
			if e != nil {
				return "", nil, e
			}
			posts, e := dl.api.GetUserMedias(userID, userMedias)
			return user.hintName, posts, e
		})
	}
	return parallelCollecPostsMedias(dl, funs)
}

func (dl *OFDl) scrapeChats(chats []scrapeIdentifier) ([]DownloadableMedia, error) {
	funs := []collecFunc{}
	for _, chat := range chats {
		funs = append(funs, func() (string, []model.Post, error) {
			chatID, e := toInt64(chat.id)
			if e != nil {
				return "", nil, e
			}
			posts, e := dl.api.GetChatMessages(chatID)
			return "", posts, e
		})
	}
	return parallelCollecPostsMedias(dl, funs)
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
	return nil
}

func (dl *OFDl) collectMedias(hintName string, post model.Post) ([]DownloadableMedia, error) {
	if len(post.Media) == 0 {
		return nil, fmt.Errorf("no media found")
	}

	mediaSet := make(map[int64]DownloadableMedia)
	for i, media := range post.Media {
		if !media.CanView || media.Files == nil {
			continue
		}
		if hintName == "" {
			hintName = post.FromUser.Username
		}
		dm := DownloadableMedia{
			PostID:  post.ID,
			MediaID: media.ID,
			Type:    media.Type,
			Time:    times(media.CreatedAt, post.CreatedAt, post.PostedAt),
			Title:   strings.TrimLeft(fmt.Sprintf("%s.%d.%d", hintName, post.ID, i), "."),
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
