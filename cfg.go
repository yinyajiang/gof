package gof

import "net/url"

func SetDebug(d bool) {
	debug = d
}

func IsDebug() bool {
	return debug
}

func Proxy() *url.URL {
	return proxy
}

func ProxyString() string {
	if proxy == nil {
		return ""
	}
	return proxy.String()
}

func SetProxy(proxyURL string) {
	if proxyURL == "" {
		return
	}
	u, err := url.Parse(proxyURL)
	if err != nil {
		return
	}
	proxy = u
}

var debug bool
var proxy *url.URL
