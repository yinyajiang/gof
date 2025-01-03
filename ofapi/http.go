package ofapi

import (
	"fmt"
	"strings"

	"github.com/yinyajiang/gof"
)

func ApiURL(format string, a ...any) string {
	return gof.OFApiDomain + ApiURLPath(format, a...)
}

func ApiURLPath(format string, a ...any) string {
	urlpath := fmt.Sprintf(format, a...)

	if !strings.HasPrefix(gof.OFApiDomain, "https://www.") {
		urlpath = strings.Replace(urlpath, "www.", "", 1)
	}
	if strings.HasPrefix(urlpath, gof.OFApiDomain) {
		urlpath, _ = strings.CutPrefix(urlpath, gof.OFApiDomain)
	}
	if !strings.HasPrefix(urlpath, "/") {
		urlpath = "/" + urlpath
	}

	if strings.HasPrefix(urlpath, gof.OFApiPathBase) {
		return urlpath
	}
	return gof.OFApiPathBase + urlpath
}
