package ofapi

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/yinyajiang/gof/common"
)

func initPublishTimeParam(addParam map[string]string, timePoint time.Time, timeDirection TimeDirection) map[string]string {
	if timePoint.IsZero() {
		timePoint = time.Now()
		timeDirection = TimeDirectionBefore
	}
	initTime := strconv.FormatInt(timePoint.Unix(), 10) + ".000000"

	if addParam == nil {
		addParam = make(map[string]string)
	}
	param := addParam
	if timeDirection == TimeDirectionBefore {
		param["beforePublishTime"] = initTime
	} else {
		param["afterPublishTime"] = initTime
	}
	return param
}

func updatePublishTimeParam(param map[string]string, timeDirection TimeDirection, moreMarker moreMarker) {
	if timeDirection == TimeDirectionBefore {
		param["beforePublishTime"] = moreMarker.TailMarker
	} else {
		if moreMarker.HeadMarker > moreMarker.TailMarker {
			param["afterPublishTime"] = moreMarker.HeadMarker
		} else {
			param["afterPublishTime"] = moreMarker.TailMarker
		}
	}
}

type moreList[T any] struct {
	HasMore  bool `json:"hasMore"`
	List     []T  `json:"list"`
	Counters any  `json:"counters"`
	moreMarker
}

type moreMarker struct {
	HeadMarker string `json:"headMarker"`
	TailMarker string `json:"tailMarker"`
}

/*
"user_id": "4045962599",
"user_agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36",
"x_bc": "1sgsadgsgsdgsdg14gsadgasdf6",
"cookie": "sess=abcdefg; auth_id=4045962599;"
*/

/*
user_id:={} || user_agent:={} || x_bc:={} || cookie:={ sess={};auth_id={} }
*/
func parseOFAuthInfo(authInfo string) OFAuthInfo {
	authinfo := OFAuthInfo{}
	for _, s := range strings.Split(authInfo, "||") {
		kv := strings.SplitN(s, ":=", 2)
		if len(kv) != 2 {
			continue
		}
		switch strings.ToLower(strings.TrimSpace(kv[0])) {
		case "user_id":
			authinfo.UserID = strings.TrimSpace(kv[1])
		case "user_agent":
			authinfo.UserAgent = strings.TrimSpace(kv[1])
		case "x_bc":
			authinfo.X_BC = strings.TrimSpace(kv[1])
		case "cookie":
			authinfo.Cookie += ";" + strings.TrimSpace(kv[1])
		}
	}
	return correctAuthInfo(authinfo)
}

func authInfoToString(authInfo OFAuthInfo) string {
	return fmt.Sprintf("user_id:=%s || user_agent:=%s || x_bc:=%s || cookie:=%s", authInfo.UserID, authInfo.UserAgent, authInfo.X_BC, authInfo.Cookie)
}

func correctAuthInfo(authInfo OFAuthInfo) OFAuthInfo {
	if authInfo.UserID == "" {
		authInfo.UserID = common.FindCookie(authInfo.Cookie, "auth_id")
	} else {
		authInfo.Cookie += ";" + "user_id=" + authInfo.UserID
	}
	authInfo.Cookie = correctCookie(authInfo.Cookie)
	return authInfo
}

func correctCookie(cookie string) string {
	sess := ""
	auth_id := ""
	common.ForeachCookie(cookie, func(k, v string) bool {
		switch strings.ToLower(k) {
		case "sess":
			sess = v
		case "auth_id":
			auth_id = v
		case "user_id":
			if auth_id == "" {
				auth_id = v
			}
		}
		return true
	})
	if sess != "" && auth_id != "" {
		return sess + ";" + auth_id
	}
	return ""
}
