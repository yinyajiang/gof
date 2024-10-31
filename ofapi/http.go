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
	if !strings.HasPrefix(gof.OFApiPathBase, "/") {
		panic("OFApiPathBase must start with / : " + gof.OFApiPathBase)
	}
	if strings.HasSuffix(gof.OFApiPathBase, "/") {
		panic("OFApiPathBase must not end with / : " + gof.OFApiPathBase)
	}
	if !strings.HasPrefix(gof.OFApiDomain, "https://") {
		panic("OFApiDomain must start with https:// : " + gof.OFApiDomain)
	}
	if strings.HasSuffix(gof.OFApiDomain, "/") {
		panic("OFApiDomain must not end with / : " + gof.OFApiDomain)
	}

	urlpath := fmt.Sprintf(format, a...)

	if !strings.HasPrefix(gof.OFApiDomain, "https://www.") {
		urlpath = strings.Replace(urlpath, "www.", "", 1)
	}
	if strings.HasPrefix(urlpath, gof.OFApiDomain) {
		urlpath, _ = strings.CutPrefix(urlpath, gof.OFApiDomain)
	}

	if !strings.HasPrefix(urlpath, "/") {
		panic("urlpath must start with / : " + urlpath)
	}

	if strings.HasPrefix(urlpath, gof.OFApiPathBase) {
		return urlpath
	}
	return gof.OFApiPathBase + urlpath
}
