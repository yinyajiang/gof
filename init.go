package gof

import "strings"

func init() {
	if !strings.HasPrefix(OFApiPathBase, "/") {
		panic("OFApiPathBase must start with / : " + OFApiPathBase)
	}
	if strings.HasSuffix(OFApiPathBase, "/") {
		panic("OFApiPathBase must not end with / : " + OFApiPathBase)
	}
	if !strings.HasPrefix(OFApiDomain, "https://") {
		panic("OFApiDomain must start with https:// : " + OFApiDomain)
	}
	if strings.HasSuffix(OFApiDomain, "/") {
		panic("OFApiDomain must not end with / : " + OFApiDomain)
	}
}
