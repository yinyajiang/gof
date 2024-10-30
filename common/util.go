package common

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/yinyajiang/gof"
)

func HttpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}

func AddHeaders(req *http.Request, addHeaders, setHeaders map[string]string) {
	for k, v := range addHeaders {
		req.Header.Add(k, v)
	}
	for k, v := range setHeaders {
		req.Header.Set(k, v)
	}
}

func IsSuccessfulStatusCode(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}

func MustMarshalJSON(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

func ParseVideoMPDInfo(dashVideoURL string) (gof.VideoMPDInfo, error) {
	split := strings.Split(dashVideoURL, ",")
	if len(split) != 6 {
		return gof.VideoMPDInfo{}, fmt.Errorf("invalid video URL format: %s", dashVideoURL)
	}
	mpdurl := split[0]
	policy := split[1]
	signature := split[2]
	keyPairID := split[3]
	mediaid := split[4]
	postid := split[5]
	return gof.VideoMPDInfo{
		MPDURL:    mpdurl,
		Policy:    policy,
		Signature: signature,
		KeyPairID: keyPairID,
		MediaID:   mediaid,
		PostID:    postid,
	}, nil
}
