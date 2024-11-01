package ofdrm

import "github.com/yinyajiang/gof/ofapi/model"

type DRMInfo struct {
	model.Drm
	MediaID int64
	PostID  int64
}
