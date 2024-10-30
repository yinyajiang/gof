package ofapi

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/yinyajiang/gof"
	"github.com/yinyajiang/gof/common"
)

func AddMPDHeaders(req *http.Request, authInfo gof.AuthInfo, mpdInfo gof.VideoMPDInfo) {
	common.AddHeaders(req, nil, map[string]string{
		"User-Agent": authInfo.UserAgent,
		"Accept":     "*/*",
		"X-BC":       authInfo.X_BC,
		"Cookie": fmt.Sprintf("CloudFront-Policy=%s; CloudFront-Signature=%s; CloudFront-Key-Pair-Id=%s; %s;",
			mpdInfo.Policy,
			mpdInfo.Signature,
			mpdInfo.KeyPairID,
			strings.TrimPrefix(authInfo.Cookie, ";"),
		),
	})
}

func OFApiMPDGet(authInfo gof.AuthInfo, mpdInfo gof.VideoMPDInfo) (body []byte, err error) {
	resp, err := OFApiMPDGetResp(authInfo, mpdInfo)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func OFApiMPDGetHeader(authInfo gof.AuthInfo, mpdInfo gof.VideoMPDInfo) (header http.Header, err error) {
	resp, err := OFApiMPDGetResp(authInfo, mpdInfo)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if !common.IsSuccessfulStatusCode(resp.StatusCode) {
		return nil, fmt.Errorf("failed to get data: %s", resp.Status)
	}
	return resp.Header, nil
}

func OFApiMPDGetResp(authInfo gof.AuthInfo, mpdInfo gof.VideoMPDInfo) (resp *http.Response, err error) {
	client := HttpClient()
	req, err := http.NewRequest("GET", mpdInfo.MPDURL, nil)
	if err != nil {
		return nil, err
	}
	AddMPDHeaders(req, authInfo, mpdInfo)
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	if !common.IsSuccessfulStatusCode(resp.StatusCode) {
		resp.Body.Close()
		return nil, fmt.Errorf("failed to get data: %s", resp.Status)
	}
	return resp, nil
}
