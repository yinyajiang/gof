package ofapi

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/yinyajiang/gof"
	"github.com/yinyajiang/gof/common"
)

func MPDHeaders(authInfo gof.AuthInfo, mpdInfo gof.MPDURLInfo) map[string]string {
	return map[string]string{
		"User-Agent": authInfo.UserAgent,
		"Accept":     "*/*",
		"X-BC":       authInfo.X_BC,
		"Cookie": fmt.Sprintf("CloudFront-Policy=%s; CloudFront-Signature=%s; CloudFront-Key-Pair-Id=%s; %s;",
			mpdInfo.Policy,
			mpdInfo.Signature,
			mpdInfo.KeyPairID,
			strings.TrimPrefix(authInfo.Cookie, ";"),
		),
	}
}

func AddMPDHeaders(req *http.Request, authInfo gof.AuthInfo, mpdInfo gof.MPDURLInfo) {
	common.AddHeaders(req, nil, MPDHeaders(authInfo, mpdInfo))
}

func OFApiMPDGet(authInfo gof.AuthInfo, mpdInfo gof.MPDURLInfo) (body []byte, err error) {
	_, body, err = OFApiMPDGetResp(authInfo, mpdInfo, true)
	return
}

func OFApiMPDGetHeader(authInfo gof.AuthInfo, mpdInfo gof.MPDURLInfo) (header http.Header, err error) {
	resp, _, err := OFApiMPDGetResp(authInfo, mpdInfo, false)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	return resp.Header, nil
}

func OFApiMPDGetResp(authInfo gof.AuthInfo, mpdInfo gof.MPDURLInfo, readAll ...bool) (resp *http.Response, body []byte, err error) {
	req, err := http.NewRequest("GET", mpdInfo.MPDURL, nil)
	if err != nil {
		return nil, nil, err
	}
	AddMPDHeaders(req, authInfo, mpdInfo)
	return common.HttpDo(req, readAll...)
}
