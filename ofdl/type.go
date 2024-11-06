package ofdl

import (
	"time"
)

type DownloadableMedia struct {
	PostID      int64
	MediaID     int64
	DownloadURL string
	Type        string
	Time        time.Time
	Title       string
	IsDrm       bool
}

type scrapeIdentifier struct {
	id       any
	hintName string
}
