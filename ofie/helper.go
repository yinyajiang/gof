package ofie

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/yinyajiang/gof/ofapi/model"
	"github.com/yinyajiang/gof/ofdrm"
)

func composeDrmUri(drm ofdrm.DRMInfo) string {
	return strings.Join([]string{
		drm.Manifest.Dash,
		drm.Signature.Dash.CloudFrontPolicy,
		drm.Signature.Dash.CloudFrontSignature,
		drm.Signature.Dash.CloudFrontKeyPairID,
		fmt.Sprint(drm.MediaID),
		fmt.Sprint(drm.PostID),
	}, ",")
}

func parseDrmUri(drmUri string) ofdrm.DRMInfo {
	split := strings.Split(drmUri, ",")
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
	return parseDrmUri(url).Manifest.Dash != ""
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

func toInt64(id any) (int64, error) {
	str := fmt.Sprint(id)
	return strconv.ParseInt(str, 10, 64)
}

func collectTitle(result *ExtractResult) {
	if result == nil || len(result.Medias) == 0 || result.Title != "" {
		return
	}
	result.Title = strings.Split(result.Medias[0].Title, titleSeparator)[0]
}

func is404Error(err error) bool {
	return err != nil && strings.Contains(err.Error(), "404")
}
