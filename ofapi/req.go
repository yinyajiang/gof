package ofapi

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/duke-git/lancet/v2/maputil"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/yinyajiang/gof"
	"github.com/yinyajiang/gof/common"
)

type Req struct {
	authInfo gof.AuthInfo
	rules    rules
}

func (r *Req) UserAgent() string {
	return r.authInfo.UserAgent
}

func (r *Req) Post(urlpath string, params any, body []byte) (data []byte, err error) {
	req, err := r.buildSignedRequest("POST", urlpath, params, body)
	if err != nil {
		return nil, err
	}
	_, data, err = common.HttpDo(req, true)
	return
}

func (r *Req) Get(urlpath string, params any) (data []byte, err error) {
	req, err := r.buildSignedRequest("GET", urlpath, params)
	if err != nil {
		return nil, err
	}
	_, data, err = common.HttpDo(req, true)
	return
}

func (r *Req) GetUnmarshal(urlpath string, params any, pointer any) (err error) {
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

func (r *Req) GetFileInfo(u string) (common.HttpFileInfo, error) {
	if common.MaybeDrmURL(u) {
		err := fmt.Errorf("[warning] url(%s) maybe drm url, use ofdrm.GetFileInfo instead", u)
		return common.HttpFileInfo{}, err
	}
	req, err := r.buildUARequest("GET", u, nil)
	if err != nil {
		return common.HttpFileInfo{}, err
	}
	resp, _, err := common.HttpDo(req, false)
	if err != nil {
		return common.HttpFileInfo{}, err
	}
	resp.Body.Close()
	return common.ParseHttpFileInfo(resp), nil
}

func (r *Req) buildSignedRequest(method, urlpath string, params any, body_ ...[]byte) (*http.Request, error) {
	req, err := r.buildRequest(method, &urlpath, params, body_...)
	if err != nil {
		return nil, err
	}
	common.AddHeaders(req, r.SignedHeaders(urlpath), nil)
	return req, nil
}

func (r *Req) buildUARequest(method, urlpath string, params any, body_ ...[]byte) (*http.Request, error) {
	req, err := r.buildRequest(method, &urlpath, params, body_...)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", r.authInfo.UserAgent)
	return req, nil
}

func (r *Req) buildRequest(method string, urlpath *string, params any, body_ ...[]byte) (*http.Request, error) {
	if urlpath == nil {
		return nil, errors.New("urlpath is nil")
	}
	*urlpath = common.HttpComposeParams(*urlpath, params)
	var body io.Reader
	if len(body_) > 0 {
		body = io.NopCloser(bytes.NewReader(body_[0]))
	}
	return http.NewRequest(method, ApiURL(*urlpath), body)
}

func (r *Req) UnsignedHeaders(mergedHeaders map[string]string) map[string]string {
	cookie := strings.TrimPrefix(r.authInfo.Cookie, ";")
	if mergedHeaders != nil && mergedHeaders["Cookie"] != "" {
		cookie = strings.TrimSuffix(cookie, ";") + ";" + strings.TrimPrefix(mergedHeaders["Cookie"], ";")
		delete(mergedHeaders, "Cookie")
	}

	return maputil.Merge(map[string]string{
		"User-Agent": r.authInfo.UserAgent,
		"Accept":     "*/*",
		"X-BC":       r.authInfo.X_BC,
		"Cookie":     cookie,
	}, mergedHeaders)
}

func (r *Req) SignedHeaders(urlpath string) map[string]string {
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
