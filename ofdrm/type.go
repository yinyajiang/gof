package ofdrm

import (
	"github.com/yinyajiang/gof/ofapi/model"
)

type DRMInfo struct {
	model.DRM
	MediaID int64
	PostID  int64
}

// uri: url、filepath、[]byte
type DRMWVDOption struct {
	WVDURI         any
	WVDMd5URIIfZip any

	ClientIDURI         any // wvd client id uri
	ClientPrivateKeyURI any // wvd private key uri

	WVDCacheDir string
}
