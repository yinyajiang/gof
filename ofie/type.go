package ofie

import (
	"fmt"
	"strings"
	"time"

	"github.com/yinyajiang/gof"
)

type MediaInfo struct {
	PostID   int64
	MediaID  int64
	MediaURI string
	Type     string
	Time     time.Time
	Title    string
	IsDrm    bool
}

type ExtractResult struct {
	Medias      []MediaInfo
	IsSingleURL bool
	IsFromCache bool
	Title       string
}

func (m MediaInfo) PostURL() string {
	return fmt.Sprintf("%s/%d/%s", gof.OFPostDomain, m.PostID, strings.Split(m.Title, ".")[0])
}

type DRMSecrets struct {
	MPDURL        string
	DecryptKey    string
	Headers       map[string]string
	Cookies       map[string]string
	CookiesString string
}

type NonDRMSecrets struct {
	Headers map[string]string
}
