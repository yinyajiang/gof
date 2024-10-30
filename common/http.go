package common

import "net/http"

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
