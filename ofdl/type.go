package ofdl

import (
	"fmt"
	"strings"
	"time"

	"github.com/yinyajiang/gof"
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

func (m DownloadableMedia) PostURL() string {
	return fmt.Sprintf("%s/%d/%s", gof.OFPostDomain, m.PostID, strings.Split(m.Title, ".")[0])
}

type DRMSecrets struct {
	DecryptKey         string
	Headers            map[string]string
	Cookies            string
	NetscapeCookieFile string
	TimeStamp          time.Time
}
