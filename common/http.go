package common

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
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
		},
	}
}

func HttpDo(req *http.Request, readAll ...bool) (*http.Response, []byte, error) {
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
