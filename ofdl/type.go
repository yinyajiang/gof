package ofdl

import (
	"time"

	"github.com/yinyajiang/gof/common"
	"github.com/yinyajiang/gof/ofdrm"
)

type DownloadableMedia struct {
	PostID      int64
	MediaID     int64
	DownloadURL string
	Type        string
	Time        time.Time
	Title       string

	_isDrm    bool
	_drmapi   *ofdrm.OFDRM
	_fileinfo *common.HttpFileInfo
}

func (dm *DownloadableMedia) IsDrm() bool {
	return dm._isDrm
}

func (dm *DownloadableMedia) DrmInfo() ofdrm.DRMInfo {
	if !dm.IsDrm() {
		return ofdrm.DRMInfo{}
	}
	return splitDRMURL(dm.DownloadURL)
}

func (dm *DownloadableMedia) FetchFileInfo() (info common.HttpFileInfo, err error) {
	if dm._fileinfo != nil {
		return *dm._fileinfo, nil
	}

	defer func() {
		if err == nil {
			dm._fileinfo = &info
		}
	}()

	if !dm.IsDrm() {
		return dm._drmapi.Req().GetFileInfo(dm.DownloadURL)
	}
	return dm._drmapi.GetFileInfo(dm.DrmInfo())
}
