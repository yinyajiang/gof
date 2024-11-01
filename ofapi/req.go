package ofapi

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/yinyajiang/gof"
	"github.com/yinyajiang/gof/common"
)

type Req struct {
	authInfo gof.AuthInfo
	rules    rules
}

func (r *Req) AuthHeaders(urlpath string) map[string]string {
	urlpath = ApiURLPath(urlpath)
	timestamp := time.Now().UTC().UnixMilli()
	hashBytes := sha1.Sum([]byte(strings.Join([]string{r.rules.StaticParam, fmt.Sprintf("%d", timestamp), urlpath, r.authInfo.UserID}, "\n")))
	hashString := strings.ToLower(hex.EncodeToString(hashBytes[:]))
	checksum := slice.Reduce(r.rules.ChecksumIndexes, func(_ int, number int, accumulator int) int {
		return accumulator + int(hashString[number])
	}, 0) + r.rules.ChecksumConstant
	sign := r.rules.Prefix + ":" + hashString + ":" + strings.ToLower(fmt.Sprintf("%X", checksum)) + ":" + r.rules.Suffix
	header := map[string]string{
		"accept":     "application/json, text/plain",
		"app-token":  r.rules.AppToken,
		"cookie":     r.authInfo.Cookie,
		"sign":       sign,
		"time":       fmt.Sprintf("%d", timestamp),
		"user-id":    r.authInfo.UserID,
		"user-agent": r.authInfo.UserAgent,
		"x-bc":       r.authInfo.X_BC,
	}
	return header
}

func (r *Req) Post(urlpath string, params any, body []byte) (data []byte, err error) {
	req, err := r.buildAuthRequest("POST", urlpath, params, body)
	if err != nil {
		return nil, err
	}
	_, data, err = common.HttpDo(req, true)
	return
}

func (r *Req) Get(urlpath string, params any) (data []byte, err error) {
	req, err := r.buildAuthRequest("GET", urlpath, params)
	if err != nil {
		return nil, err
	}
	_, data, err = common.HttpDo(req, true)
	return
}

func (r *Req) GetUnmashel(urlpath string, params any, pointer any) (err error) {
	data, err := r.Get(urlpath, params)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, pointer)
	if err != nil {
		fmt.Printf("urlpath: %s, data unmarshal error: %v\n", urlpath, err)
		fmt.Println(string(data))
	}
	return err
}

func (r *Req) buildAuthRequest(method, urlpath string, params any, body_ ...[]byte) (*http.Request, error) {
	switch params := any(params).(type) {
	case string:
		params = strings.TrimLeft(params, "?")
		if params != "" {
			urlpath = urlpath + "?" + params
		}
	case map[string]string:
		if len(params) > 0 {
			query := url.Values{}
			for k, v := range params {
				query.Add(k, v)
			}
			if strings.Contains(urlpath, "?") {
				urlpath = urlpath + "&" + query.Encode()
			} else {
				urlpath = urlpath + "?" + query.Encode()
			}
		}
	}

	var body io.Reader
	if len(body_) > 0 {
		body = io.NopCloser(bytes.NewReader(body_[0]))
	}
	req, err := http.NewRequest(method, ApiURL(urlpath), body)
	if err != nil {
		return nil, err
	}
	r.AddAuthHeaders(req, urlpath)
	return req, nil
}

func (r *Req) AddAuthHeaders(req *http.Request, urlpath string) {
	common.AddHeaders(req, r.AuthHeaders(urlpath), nil)
}

func (r *Req) NoSignHeaders() map[string]string {
	return map[string]string{
		"User-Agent": r.authInfo.UserAgent,
		"Accept":     "*/*",
		"X-BC":       r.authInfo.X_BC,
		"Cookie":     strings.TrimPrefix(r.authInfo.Cookie, ";"),
	}
}
