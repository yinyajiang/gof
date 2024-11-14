package common

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/yinyajiang/gof"
	"golang.org/x/exp/rand"
)

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

func HttpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy: func(*http.Request) (*url.URL, error) {
				return gof.Proxy(), nil
			},
		},
	}
}

var lastRequestTime time.Time

func HttpDo(req *http.Request, readAll ...bool) (*http.Response, []byte, error) {
	if !gof.IsDisableTimeInterval() {
		now := time.Now()
		if lastRequestTime.IsZero() {
			lastRequestTime = now
		} else {
			since := now.Sub(lastRequestTime)
			if since < gof.MaxTimeInterval {
				time.Sleep(time.Duration(rand.Int63n(int64(gof.MaxTimeInterval - since))))
			}
		}
		lastRequestTime = now
	}

	resp, err := HttpClient().Do(req)
	if err != nil {
		return resp, nil, err
	}
	if IsSuccessfulStatusCode(resp.StatusCode) {
		if len(readAll) > 0 && readAll[0] {
			body, e := io.ReadAll(resp.Body)
			resp.Body.Close()
			return resp, body, e
		}
		return resp, nil, nil
	}

	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	err = fmt.Errorf("[%s] [%s] failed, err: %v, status: %d, body: %s", req.Method, req.URL.String(), err, resp.StatusCode, string(body))
	return nil, nil, err
}

func HttpGet(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	_, data, err := HttpDo(req, true)
	return data, err
}

func HttpGetUnmarshal(url string, pointer any) error {
	data, err := HttpGet(url)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, pointer)
	if err != nil {
		err = fmt.Errorf("unmarshal %s failed, err: %v, data: %s", url, err, string(data))
	}
	return err
}

func HttpComposeParams(urlpath string, params any) string {
	switch params := any(params).(type) {
	case string:
		params = strings.TrimLeft(params, "?")
		if params != "" {
			if strings.Contains(urlpath, "?") {
				urlpath = urlpath + "&" + params
			} else {
				urlpath = urlpath + "?" + params
			}
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
	return urlpath
}

func ParseHttpFileInfo(resp *http.Response) HttpFileInfo {
	var headerInfo HttpFileInfo
	lastModified, err := time.Parse(time.RFC1123, resp.Header.Get("Last-Modified"))
	if err != nil {
		headerInfo.LastModified = time.Now()
		fmt.Printf("parse last modified: %v\n", err)
	} else {
		headerInfo.LastModified = lastModified.Local()
	}
	headerInfo.ContentLength = resp.ContentLength
	return headerInfo
}

func ConvertCookieToNetscape(cookieStr string, domain string) string {
	var result strings.Builder
	// Write header comment
	result.WriteString("# Netscape HTTP Cookie File\n")
	result.WriteString("# https://curl.haxx.se/rfc/cookie_spec.html\n")
	result.WriteString("# This is a generated file!  Do not edit.\n\n")

	// Parse cookies
	cookies := strings.Split(cookieStr, ";")
	for _, cookie := range cookies {
		cookie = strings.TrimSpace(cookie)
		if cookie == "" {
			continue
		}

		parts := strings.SplitN(cookie, "=", 2)
		if len(parts) != 2 {
			continue
		}

		name := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Format: domain HTTP_ONLY path SECURE expiry name value
		// Using default values: HTTP_ONLY=FALSE, path=/, SECURE=FALSE, expiry=0 (session cookie)
		line := fmt.Sprintf("%s\tFALSE\t/\tFALSE\t0\t%s\t%s\n",
			domain, name, value)
		result.WriteString(line)
	}
	return result.String()
}
