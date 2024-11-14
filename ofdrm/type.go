package ofdrm

import (
	"github.com/yinyajiang/gof/ofapi/model"
)

type DRMInfo struct {
	model.DRM
	MediaID int64
	PostID  int64
}

type DRMWVDOption struct {
	WVD              []byte // wvd,
	RawWVDID         []byte // wvd client id
	RawWVDPrivateKey []byte // wvd private key

	WVDURI              string // wvd uri
	ClientIDURI         string // wvd client id uri
	ClientPrivateKeyURI string // wvd private key uri

	ClientCacheDir      string
	ClientCachePriority bool
}
