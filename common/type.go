package common

import "time"

type HttpFileInfo struct {
	LastModified  time.Time
	ContentLength int64
}
