package ofdl

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/yinyajiang/gof"
	"github.com/yinyajiang/gof/ofapi/model"
	"github.com/yinyajiang/gof/ofdrm"
)

func correctURL(url string) string {
	if i := strings.Index(url, "?"); i != -1 {
		url = url[:i]
	}
	if !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}
	url = strings.Replace(strings.TrimSpace(url), "www.", "", 1)
	return url
}

func isHomeURL(url string) bool {
	url = correctURL(url)
	return strings.TrimRight(url, "/") == strings.TrimRight(gof.OFPostDomain, "/")
}

func reUrlPathParse(ori string, rePath string, minSplit int) (pu *url.URL, splitPaths []string, err error) {
	u := correctURL(ori)
	if rePath != "" && !strings.HasPrefix(rePath, "/") {
		rePath = "/" + rePath
	}
	re := `(?i)` + regexp.QuoteMeta(gof.OFPostDomain) + rePath
	regex := regexp.MustCompile(re)
	if !regex.MatchString(u) {
		return nil, nil, fmt.Errorf("invalid url: %s, for regex: %s", ori, re)
	}
	pu, err = url.Parse(u)
	if err != nil {
		return nil, nil, err
	}
	splitPaths = strings.Split(strings.TrimLeft(pu.Path, "/"), "/")
	if minSplit > 0 {
		if len(splitPaths) < minSplit {
			return nil, nil, fmt.Errorf("invalid url path length: %d, %s", len(splitPaths), ori)
		}
	}
	return pu, splitPaths, nil
}

func composeDRMURL(mediaID int64, postID int64, drm model.DRM) string {
	return strings.Join([]string{
		drm.Manifest.Dash,
		drm.Signature.Dash.CloudFrontPolicy,
		drm.Signature.Dash.CloudFrontSignature,
		drm.Signature.Dash.CloudFrontKeyPairID,
		fmt.Sprint(mediaID),
		fmt.Sprint(postID),
	}, ",")
}

func splitDRMURL(drmUrl string) ofdrm.DRMInfo {
	split := strings.Split(drmUrl, ",")
	if len(split) != 6 {
		return ofdrm.DRMInfo{}
	}
	manifest := split[0]
	policy := split[1]
	signature := split[2]
	keyPairID := split[3]
	mediaid, err := strconv.ParseInt(split[4], 10, 64)
	if err != nil {
		return ofdrm.DRMInfo{}
	}
	postid, err := strconv.ParseInt(split[5], 10, 64)
	if err != nil {
		return ofdrm.DRMInfo{}
	}
	return ofdrm.DRMInfo{
		PostID:  postid,
		MediaID: mediaid,
		DRM: model.DRM{
			Manifest: model.Manifest{
				Dash: manifest,
			},
			Signature: model.Signature{
				Dash: model.CloudFront{
					CloudFrontPolicy:    policy,
					CloudFrontSignature: signature,
					CloudFrontKeyPairID: keyPairID,
				},
			},
		},
	}
}

func times(times ...time.Time) time.Time {
	if len(times) == 0 {
		return time.Time{}
	}
	for _, t := range times {
		if !t.IsZero() {
			return t
		}
	}
	return time.Time{}
}

type collecFunc func() (string, []model.Post, error)

func parallelCollecPostsMedias(dl *OFDl, funs []collecFunc) ([]DownloadableMedia, error) {
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
			userName, posts, err := fun()
			lock.Lock()
			defer lock.Unlock()

			var medias []DownloadableMedia
			if err == nil {
				medias, err = collecPostsMedias(dl, userName, posts)
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

func collecPostsMedias(dl *OFDl, hintUser string, posts []model.Post) ([]DownloadableMedia, error) {
	if len(posts) == 0 {
		return nil, fmt.Errorf("posts is empty")
	}
	results := []DownloadableMedia{}
	for _, post := range posts {
		medias, e := dl.collectPostMedia(hintUser, post)
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
