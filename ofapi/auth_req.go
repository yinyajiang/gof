package ofapi

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
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

func AuthHeaders(urlpath string, auth gof.AuthInfo, rules gof.Rules) map[string]string {
	urlpath = ApiURLPath(urlpath)

	timestamp := time.Now().UTC().UnixMilli()
	hashBytes := sha1.Sum([]byte(strings.Join([]string{rules.StaticParam, fmt.Sprintf("%d", timestamp), urlpath, auth.UserID}, "\n")))
	hashString := strings.ToLower(hex.EncodeToString(hashBytes[:]))
	checksum := slice.Reduce(rules.ChecksumIndexes, func(_ int, number int, accumulator int) int {
		return accumulator + int(hashString[number])
	}, 0) + rules.ChecksumConstant
	sign := rules.Prefix + ":" + hashString + ":" + strings.ToLower(fmt.Sprintf("%X", checksum)) + ":" + rules.Suffix
	header := map[string]string{
		"accept":     "application/json, text/plain",
		"app-token":  rules.AppToken,
		"cookie":     auth.Cookie,
		"sign":       sign,
		"time":       fmt.Sprintf("%d", timestamp),
		"user-id":    auth.UserID,
		"user-agent": auth.UserAgent,
		"x-bc":       auth.X_BC,
	}
	return header
}

func OFApiAuthPost[P Params](urlpath string, params P, auth gof.AuthInfo, rules gof.Rules, body ...[]byte) (data []byte, err error) {
	req, err := buildAuthRequest("POST", urlpath, params, auth, rules, body...)
	if err != nil {
		return nil, err
	}
	resp, err := HttpClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if !common.IsSuccessfulStatusCode(resp.StatusCode) {
		return nil, fmt.Errorf("failed to post data: %s, urlpath: %s", resp.Status, urlpath)
	}
	return io.ReadAll(resp.Body)
}

func OFApiAuthGet[P Params](urlpath string, params P, auth gof.AuthInfo, rules gof.Rules) (data []byte, err error) {
	req, err := buildAuthRequest("GET", urlpath, params, auth, rules)
	if err != nil {
		return nil, err
	}
	resp, err := HttpClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if !common.IsSuccessfulStatusCode(resp.StatusCode) {
		return nil, fmt.Errorf("failed to get data: %s, urlpath: %s", resp.Status, urlpath)
	}
	return io.ReadAll(resp.Body)
}

type Params interface {
	string | map[string]string
}

func buildAuthRequest[P Params](method, urlpath string, params P, auth gof.AuthInfo, rules gof.Rules, body_ ...[]byte) (*http.Request, error) {
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
			urlpath = urlpath + "?" + query.Encode()
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
	AddAuthHeaders(req, urlpath, auth, rules)
	return req, nil
}

func AddAuthHeaders(req *http.Request, urlpath string, auth gof.AuthInfo, rules gof.Rules) {
	common.AddHeaders(req, AuthHeaders(urlpath, auth, rules), nil)
}
