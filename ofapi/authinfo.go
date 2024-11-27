package ofapi

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/yinyajiang/gof/common"
)

/*
"user_id": "4045962599",
"user_agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36",
"x_bc": "1sgsadgsgsdgsdg14gsadgasdf6",
"cookie": "sess=abcdefg; auth_id=4045962599;"
*/

/*
user_id:={} || user_agent:={} || x_bc:={} || cookie:={ sess={};auth_id={} }
*/
func String2AuthInfo(authInfo string) OFAuthInfo {
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
			if !strings.HasSuffix(authinfo.Cookie, ";") {
				authinfo.Cookie += ";"
			}
			authinfo.Cookie += strings.TrimSpace(kv[1])
		}
	}
	return correctAuthInfo(authinfo)
}

func Raw2AuthInfo(ua, cookiefile string) (OFAuthInfo, error) {
	cookies, err := common.ParseCookieFile(cookiefile)
	if err != nil {
		return OFAuthInfo{}, err
	}
	cookiestr := ""
	for k, v := range cookies {
		cookiestr += fmt.Sprintf("%s=%s;", k, v)
	}
	return correctAuthInfo(OFAuthInfo{
		Cookie:    cookiestr,
		UserAgent: ua,
	}), nil
}

func loadCacheAuthInfo(cacheDir string) (OFAuthInfo, error) {
	data, err := os.ReadFile(filepath.Join(cacheDir, "auth"))
	if err != nil {
		return OFAuthInfo{}, err
	}
	var auth OFAuthInfo
	err = json.Unmarshal(data, &auth)
	return auth, err
}

func cacheAuthInfo(cacheDir string, auth OFAuthInfo) {
	data, err := json.Marshal(auth)
	if err != nil {
		fmt.Printf("marshal auth failed, err: %v\n", err)
		return
	}
	fileutil.CreateDir(cacheDir)
	os.WriteFile(filepath.Join(cacheDir, "auth"), data, 0644)
}

func authInfo2String(authInfo OFAuthInfo) string {
	return fmt.Sprintf("user_id:=%s || user_agent:=%s || x_bc:=%s || cookie:=%s", authInfo.UserID, authInfo.UserAgent, authInfo.X_BC, authInfo.Cookie)
}

func correctAuthInfo(authInfo OFAuthInfo) OFAuthInfo {
	if authInfo.UserID == "" {
		authInfo.UserID = common.FindCookie(authInfo.Cookie, "auth_id")
	} else {
		if !strings.HasSuffix(authInfo.Cookie, ";") {
			authInfo.Cookie += ";"
		}
		authInfo.Cookie += "user_id=" + authInfo.UserID
	}
	if authInfo.X_BC == "" {
		authInfo.X_BC = common.FindCookie(authInfo.Cookie, "fp")
	}
	authInfo.Cookie = _trimCookie(authInfo.Cookie)
	return authInfo
}

func _trimCookie(cookie string) string {
	sess := ""
	auth_id := ""
	common.ForeachCookie(cookie, func(k, v string) bool {
		switch strings.ToLower(k) {
		case "sess":
			if v != "" {
				sess = v
			}
		case "auth_id":
			if v != "" {
				auth_id = v
			}
		case "user_id":
			if auth_id == "" && v != "" {
				auth_id = v
			}
		}
		return true
	})
	if sess != "" && auth_id != "" {
		return fmt.Sprintf("sess=%s;auth_id=%s", sess, auth_id)
	}
	return ""
}
