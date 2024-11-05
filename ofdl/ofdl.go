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

func (dl *OFDl) ScrapeHome() ([]DownloadableMedia, error) {
	subs, err := dl.api.GetSubscriptions(ofapi.SubscritionTypeActive)
	if err != nil {
		return nil, err
	}
	funs := []collecFunc{}
	for _, sub := range subs {
		funs = append(funs, func() (string, []model.Post, error) {
			posts, e := dl.api.GetUserPosts(sub.ID)
			return sub.Username, posts, e
		})
	}
	return parallelCollecPostsMedias(dl, funs)
}

func (dl *OFDl) ScrapeUserMedia(url string) ([]DownloadableMedia, error) {
	_, split, err := reUrlPathParse(url, `/[A-Za-z0-9\.]+`, 1)
	if err != nil {
		return nil, err
	}
	userName := split[0]
	usr, err := dl.api.GetUserByUsername(userName)
	if err != nil {
		return nil, err
	}
	posts, err := dl.api.GetUserPosts(usr.ID)
	if err != nil {
		return nil, err
	}
	return collecPostsMedias(dl, userName, posts)
}

func (dl *OFDl) ScrapePostMedia(url string) ([]DownloadableMedia, error) {
	_, split, err := reUrlPathParse(url, `/[0-9]+/[A-Za-z0-9\.]+`, 2)
	if err != nil {
		return nil, err
	}
	postID := split[0]
	userName := split[1]
	post, err := dl.api.GetPost(postID)
	if err != nil {
		return nil, err
	}
	return dl.collectPostMedia(userName, post)
}

func (dl *OFDl) Download(dir string, media DownloadableMedia) error {
	if !media.IsDrm() {
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
	drminfo := media.DrmInfo()
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

func (dl *OFDl) collectPostMedia(hintUser string, post model.Post) ([]DownloadableMedia, error) {
	if len(post.Media) == 0 {
		return nil, fmt.Errorf("no media found")
	}

	mediaSet := make(map[int64]DownloadableMedia)
	for i, item := range post.Media {
		if !item.CanView || item.Files == nil {
			continue
		}
		dm := DownloadableMedia{
			PostID:  post.ID,
			MediaID: item.ID,
			Type:    item.Type,
			Time:    times(item.CreatedAt, post.CreatedAt, post.PostedAt),
			Title:   strings.TrimLeft(fmt.Sprintf("%s.%d.%d", hintUser, post.ID, i), "."),
			_drmapi: dl.drmapi,
		}

		if item.Files.Drm == nil {
			if item.Files.Full != nil {
				dm.DownloadURL = item.Files.Full.URL
			} else if item.Files.Preview != nil {
				dm.DownloadURL = item.Files.Preview.URL
			}
			dm._isDrm = false
			if strings.Contains(dm.DownloadURL, "upload") {
				continue
			}
		} else {
			dm.DownloadURL = composeDRMURL(item.ID, post.ID, *item.Files.Drm)
			dm._isDrm = true
		}
		mediaSet[item.ID] = dm
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
