package common

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/yinyajiang/gof"
)

func MustMarshalJSON(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

func ParseVideoMPDInfo(dashVideoURL string) (gof.MPDInfo, error) {
	split := strings.Split(dashVideoURL, ",")
	if len(split) != 6 {
		return gof.MPDInfo{}, fmt.Errorf("invalid video URL format: %s", dashVideoURL)
	}
	mpdurl := split[0]
	policy := split[1]
	signature := split[2]
	keyPairID := split[3]
	mediaid := split[4]
	postid := split[5]
	return gof.MPDInfo{
		MPDURL:    mpdurl,
		Policy:    policy,
		Signature: signature,
		KeyPairID: keyPairID,
		MediaID:   mediaid,
		PostID:    postid,
	}, nil
}
