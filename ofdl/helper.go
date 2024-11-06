package ofdl

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/yinyajiang/gof/ofapi/model"
	"github.com/yinyajiang/gof/ofdrm"
)

func composeDRMURL(drm ofdrm.DRMInfo) string {
	return strings.Join([]string{
		drm.Manifest.Dash,
		drm.Signature.Dash.CloudFrontPolicy,
		drm.Signature.Dash.CloudFrontSignature,
		drm.Signature.Dash.CloudFrontKeyPairID,
		fmt.Sprint(drm.MediaID),
		fmt.Sprint(drm.PostID),
	}, ",")
}

func parseDRMURL(drmUrl string) ofdrm.DRMInfo {
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

func isDrmURL(url string) bool {
	return parseDRMURL(url).Manifest.Dash != ""
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
			hintName, posts, err := fun()
			lock.Lock()
			defer lock.Unlock()

			var medias []DownloadableMedia
			if err == nil {
				medias, err = collecMutilMedias(dl, hintName, posts)
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

func collecMutilMedias(dl *OFDl, hintName string, posts []model.Post) ([]DownloadableMedia, error) {
	if len(posts) == 0 {
		return nil, fmt.Errorf("posts is empty")
	}
	results := []DownloadableMedia{}
	for _, post := range posts {
		medias, e := dl.collectMedias(hintName, post)
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

func toInt64(id any) (int64, error) {
	str := fmt.Sprint(id)
	return strconv.ParseInt(str, 10, 64)
}
