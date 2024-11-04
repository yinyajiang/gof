package ofdl

import (
	"strings"

	"github.com/duke-git/lancet/v2/maputil"
	"github.com/yinyajiang/gof/common"
	"github.com/yinyajiang/gof/ofapi/model"
	"github.com/yinyajiang/gof/ofdrm"
)

type DownloadableMedia struct {
	PostID        int64
	MediaID       int64
	DownloadedURL string
	DownloadedDrm *model.Drm
	Type          string
}

func (dm *DownloadableMedia) IsDrm() bool {
	return dm.DownloadedDrm != nil
}

func (dm *DownloadableMedia) DrmInfo() ofdrm.DRMInfo {
	if dm.DownloadedDrm == nil {
		return ofdrm.DRMInfo{
			PostID:  dm.PostID,
			MediaID: dm.MediaID,
		}
	}
	return ofdrm.DRMInfo{
		PostID:  dm.PostID,
		MediaID: dm.MediaID,
		Drm:     *dm.DownloadedDrm,
	}
}

func (dm *DownloadableMedia) FileInfo(ofdrm *ofdrm.OFDRM) (common.HttpFileInfo, error) {
	if !dm.IsDrm() {
		return ofdrm.Req().GetFileInfo(dm.DownloadedURL)
	}
	return ofdrm.GetFileInfo(dm.DrmInfo())
}

func CollatePostMedia(post model.Post) []DownloadableMedia {
	if len(post.Media) == 0 {
		return nil
	}

	mediaSet := make(map[int64]DownloadableMedia)
	for _, media := range post.Media {
		if !media.CanView || media.Files == nil {
			continue
		}
		if media.Files.Drm == nil {
			dm := DownloadableMedia{
				PostID:  post.ID,
				MediaID: media.ID,
				Type:    media.Type,
			}
			if media.Files.Full != nil {
				dm.DownloadedURL = media.Files.Full.URL
			} else if media.Files.Preview != nil {
				dm.DownloadedURL = media.Files.Preview.URL
			}
			if strings.Contains(dm.DownloadedURL, "upload") {
				continue
			}
			mediaSet[media.ID] = dm
		} else {
			dm := DownloadableMedia{
				PostID:        post.ID,
				MediaID:       media.ID,
				Type:          media.Type,
				DownloadedDrm: media.Files.Drm,
			}
			mediaSet[media.ID] = dm
		}
	}
	return maputil.Values(mediaSet)
}
