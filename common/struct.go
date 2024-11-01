package common

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/yinyajiang/gof"
)

type MPDURLInfo struct {
	MPDURL    string
	Policy    string
	Signature string
	KeyPairID string
	MediaID   string
	PostID    string
}

func ParseVideoMPDInfo(dashVideoURL string) (MPDURLInfo, error) {
	split := strings.Split(dashVideoURL, ",")
	if len(split) != 6 {
		return MPDURLInfo{}, fmt.Errorf("invalid video URL format: %s", dashVideoURL)
	}
	mpdurl := split[0]
	policy := split[1]
	signature := split[2]
	keyPairID := split[3]
	mediaid := split[4]
	postid := split[5]
	return MPDURLInfo{
		MPDURL:    mpdurl,
		Policy:    policy,
		Signature: signature,
		KeyPairID: keyPairID,
		MediaID:   mediaid,
		PostID:    postid,
	}, nil
}

func ParseSinglePostURL(postURL string) (gof.PostURLInfo, error) {
	postURL = strings.Replace(strings.TrimSpace(postURL), "www.", "", 1)
	re := regexp.MustCompile(`(?i)https://onlyfans\.com/[0-9]+/[A-Za-z0-9]+`)
	if !re.MatchString(postURL) {
		return gof.PostURLInfo{}, fmt.Errorf("invalid post URL format: %s", postURL)
	}
	u, err := url.Parse(postURL)
	if err != nil {
		return gof.PostURLInfo{}, err
	}
	split := strings.Split(strings.TrimLeft(u.Path, "/"), "/")
	if len(split) < 2 {
		return gof.PostURLInfo{}, fmt.Errorf("invalid post URL format, length: %d: %s", len(split), postURL)
	}
	return gof.PostURLInfo{
		PostID:   split[0],
		UserName: split[1],
	}, nil
}

func MustMarshalJSON(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}
