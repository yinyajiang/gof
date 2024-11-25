package ofie

import (
	"time"
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

type ExtractOption struct {
	DisableCache bool
}

type DRMSecrets struct {
	MPDURL        string
	DecryptKey    string
	Headers       map[string]string
	Cookies       map[string]string
	CookiesString string
}

type FetchDRMSecretsOption struct {
	DisableCache bool
	MustClient   bool
}

type NonDRMSecrets struct {
	Headers map[string]string
}
