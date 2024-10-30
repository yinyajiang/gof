package ofapi

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	"github.com/yinyajiang/gof"
)

func HttpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}

func ApiURL(format string, a ...any) string {
	urlpath := fmt.Sprintf(format, a...)
	return gof.OFApiDomain + ApiURLPath(urlpath)
}

func ApiURLPath(format string, a ...any) string {
	urlpath := fmt.Sprintf(format, a...)

	if !strings.HasPrefix(urlpath, "/") {
		panic("urlpath must start with / : " + urlpath)
	}
	if !strings.HasPrefix(gof.OFApiPathBase, "/") {
		panic("OFApiPathBase must start with / : " + gof.OFApiPathBase)
	}
	if strings.HasSuffix(gof.OFApiPathBase, "/") {
		panic("OFApiPathBase must not end with / : " + gof.OFApiPathBase)
	}

	if strings.HasPrefix(urlpath, gof.OFApiPathBase) {
		return urlpath
	}
	return gof.OFApiPathBase + urlpath
}
