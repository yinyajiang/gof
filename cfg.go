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

func SetEnableTimeInterval(t bool) {
	isDisableTimeInterval = !t
}

func IsDisableTimeInterval() bool {
	return isDisableTimeInterval
}

var debug bool
var proxy *url.URL
var isDisableTimeInterval bool
