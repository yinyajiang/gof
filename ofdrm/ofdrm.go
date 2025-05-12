package ofdrm

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/antchfx/xmlquery"
	"github.com/duke-git/lancet/v2/slice"
	widevine "github.com/iyear/gowidevine"
	"github.com/iyear/gowidevine/widevinepb"
	"github.com/yinyajiang/gof/common"
	"github.com/yinyajiang/gof/ofapi"
)

type OFDRM struct {
	req               *ofapi.Req
	cdrmProjectServer []string
	wvd               *wvdSt
}

type OFDRMConfig struct {
	WVDOption                 DRMWVDOption
	OptionalCDRMProjectServer []string
}

func NewOFDRM(req *ofapi.Req, config OFDRMConfig) (*OFDRM, error) {
	wvd, err := loadWVD(config.WVDOption)
	if err != nil && len(config.OptionalCDRMProjectServer) == 0 {
		return nil, err
	}
	return &OFDRM{
		req:               req,
		cdrmProjectServer: config.OptionalCDRMProjectServer,
		wvd:               wvd,
	}, nil
}

func (c *OFDRM) Req() *ofapi.Req {
	return c.req
}

func (c *OFDRM) WVD() *wvdSt {
	return c.wvd
}

func (c *OFDRM) GetDecryptedKeyAuto(drm DRMInfo) (string, error) {
	var clientErr error
	var key string
	if c.wvd != nil {
		key, clientErr = c.GetDecryptedKeyByClient(drm)
		if clientErr == nil {
			return key, nil
		}
		fmt.Println("failed to get decrypted key by client: ", clientErr)
	}

	key, serverErr := c.GetDecryptedKeyCDMProject(drm)
	if serverErr == nil {
		return key, nil
	} else {
		fmt.Println("failed to get decrypted key by server: ", serverErr)
	}

	key, oferr := c.GetDecryptedKeyByOFDL(drm)
	if oferr == nil {
		return key, nil
	} else {
		fmt.Println("failed to get decrypted key by ofdl: ", oferr)
	}

	if clientErr != nil {
		return "", clientErr
	}
	return "", serverErr
}

func (c *OFDRM) GetDecryptedKeyByClient(drm DRMInfo) (string, error) {
	pssh, err := c.getDRMPSSH(drm)
	if err != nil {
		return "", err
	}

	keys, err := c.getWidevineKeys(c.drmURLPath(drm), pssh)
	if err != nil {
		return "", err
	}

	key := keys[0]
	decryptedKey := strings.ToLower(hex.EncodeToString(key.ID)) + ":" + strings.ToLower(hex.EncodeToString(key.Key))
	return decryptedKey, nil
}

func (c *OFDRM) GetFileInfo(drm DRMInfo) (common.HttpFileInfo, error) {
	resp, _, err := c.drmHttpGetResp(drm, false)
	if err != nil {
		return common.HttpFileInfo{}, err
	}
	resp.Body.Close()
	return common.ParseHttpFileInfo(resp), nil
}

const publicCDRMProjectServer = "https://cdrm-project.com/api/decrypt"

func (c *OFDRM) GetDecryptedKeyCDMProject(drm DRMInfo) (string, error) {
	const fixedServerURL = publicCDRMProjectServer

	serverURLs := c.cdrmProjectServer
	if !slice.Contain(serverURLs, fixedServerURL) {
		serverURLs = append(serverURLs, fixedServerURL)
	}

	pssh, err := c.getDRMPSSH(drm)
	if err != nil {
		return "", err
	}
	maxAttempts := 30 * len(serverURLs)
	for i := 0; i < maxAttempts; i++ {
		serverURL := serverURLs[i%len(serverURLs)]
		decryptedKey, err := c.getVideoDecryptedKeyByCDMProject(serverURL, pssh, drm)
		if err == nil {
			return decryptedKey, nil
		}
		fmt.Printf("try %d/%d, failed to get decrypted key from %s: %s\n", i+1, maxAttempts, serverURL, err)
	}
	return "", fmt.Errorf("all servers failed")
}

func (c *OFDRM) _getVideoDecryptedKeyByCDMProject(serverURL, pssh, licurl string, headers map[string]string) (string, error) {
	data := common.MustMarshalJSON(map[string]string{
		"pssh":    pssh,
		"licurl":  licurl,
		"headers": common.MustUnmarshalJSONStr(headers),
		"cookies": "",
		"data":    "",
	})
	req, err := http.NewRequest("POST", serverURL, io.NopCloser(bytes.NewReader(data)))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	_, body, err := common.HttpDo(req, true)
	if err != nil {
		return "", err
	}
	var result map[string]any
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", fmt.Errorf("parse response: %w, contnet:%s", err, string(body))
	}
	if !strings.EqualFold(fmt.Sprint(result["status"]), "success") {
		return "", fmt.Errorf("status: %s, contnet:%s", result["status"], string(body))
	}

	msg, ok := result["message"]
	if !ok {
		return "", fmt.Errorf("no message")
	}
	return strings.TrimSpace(msg.(string)), nil
}

func (c *OFDRM) getVideoDecryptedKeyByCDMProject(serverURL, pssh string, drm DRMInfo) (string, error) {
	return c._getVideoDecryptedKeyByCDMProject(serverURL, pssh, ofapi.ApiURL(c.drmURLPath(drm)), c.req.SignedHeaders(c.drmURLPath(drm)))
}

func (c *OFDRM) drmURLPath(drm DRMInfo) string {
	return ofapi.ApiURLPath("/users/media/%d/drm/post/%d?type=widevine", drm.MediaID, drm.PostID)
}

func (c *OFDRM) getDRMPSSH(drm DRMInfo) (string, error) {
	data, err := c.drmHttpGet(drm)
	if err != nil {
		return "", err
	}
	doc, err := xmlquery.Parse(bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	psshEles := doc.SelectElements("//cenc:pssh")
	if len(psshEles) < 1 {
		return "", fmt.Errorf("//cenc:pssh count < 1")
	}
	pssh := strings.TrimSpace(psshEles[1].InnerText())
	return pssh, nil
}

func (c *OFDRM) getWidevineKeys(urlpath, pssh string) ([]*widevine.Key, error) {
	device, err := widevine.NewDevice(
		widevine.FromWVD(bytes.NewReader(c.wvd.WVD())),
	)
	if err != nil {
		return nil, fmt.Errorf("create device: %w", err)
	}
	cdm := widevine.NewCDM(device)

	psshData, err := base64.StdEncoding.DecodeString(pssh)
	if err != nil {
		return nil, fmt.Errorf("decode base64: %w", err)
	}

	wpssh, err := widevine.NewPSSH(psshData)
	if err != nil {
		return nil, fmt.Errorf("parse pssh: %w", err)
	}

	cert, err := c.loadWidevineServiceCert(urlpath)
	if err != nil {
		return nil, fmt.Errorf("get service cert: %w", err)
	}
	challenge, parseLicenseFunc, err := cdm.GetLicenseChallenge(wpssh, widevinepb.LicenseType_AUTOMATIC, true, cert)
	if err != nil {
		return nil, fmt.Errorf("get license challenge: %w", err)
	}
	license, err := c.req.Post(urlpath, nil, challenge)
	if err != nil {
		return nil, err
	}
	keys, err := parseLicenseFunc(license)
	if err != nil {
		return nil, fmt.Errorf("parse license: %w", err)
	}
	keys = slice.Filter(keys, func(_ int, key *widevine.Key) bool {
		return key.Type == widevinepb.License_KeyContainer_CONTENT
	})
	if len(keys) == 0 {
		keys = slice.Filter(keys, func(_ int, key *widevine.Key) bool {
			return key.Type != widevinepb.License_KeyContainer_SIGNING
		})
	}
	if len(keys) == 0 {
		return nil, fmt.Errorf("no keys")
	}
	return keys, nil
}

func (c *OFDRM) loadWidevineServiceCert(urlpath string) (*widevinepb.DrmCertificate, error) {
	serviceCert, err := c.req.Post(urlpath, nil, widevine.ServiceCertificateRequest)
	if err != nil {
		return nil, err
	}
	cert, err := widevine.ParseServiceCert(serviceCert)
	if err != nil {
		return nil, fmt.Errorf("parse service cert: %w", err)
	}
	return cert, nil
}

func (c *OFDRM) DRMHTTPHeaders(drm DRMInfo) map[string]string {
	return c.req.UnsignedHeaders(map[string]string{
		"Cookie": fmt.Sprintf("CloudFront-Policy=%s; CloudFront-Signature=%s; CloudFront-Key-Pair-Id=%s",
			drm.DRM.Signature.Dash.CloudFrontPolicy,
			drm.DRM.Signature.Dash.CloudFrontSignature,
			drm.DRM.Signature.Dash.CloudFrontKeyPairID),
	})
}

func (c *OFDRM) drmHttpGetResp(drm DRMInfo, readAll ...bool) (resp *http.Response, body []byte, err error) {
	req, err := http.NewRequest("GET", drm.DRM.Manifest.Dash, nil)
	if err != nil {
		return nil, nil, err
	}
	common.AddHeaders(req, nil, c.DRMHTTPHeaders(drm))
	return common.HttpDo(req, readAll...)
}

func (c *OFDRM) drmHttpGet(drm DRMInfo) (body []byte, err error) {
	_, body, err = c.drmHttpGetResp(drm, true)
	return
}

func (c *OFDRM) GetDecryptedKeyByOFDL(drm DRMInfo) (string, error) {
	pssh, err := c.getDRMPSSH(drm)
	if err != nil {
		return "", err
	}
	data := common.MustMarshalJSON(map[string]string{
		"pssh":       pssh,
		"licenceURL": ofapi.ApiURL(c.drmURLPath(drm)),
		"headers":    common.MustUnmarshalJSONStr(c.req.SignedHeaders(c.drmURLPath(drm))),
	})
	req, err := http.NewRequest("POST", "https://ofdl.tools/WV", io.NopCloser(bytes.NewReader(data)))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	_, body, err := common.HttpDo(req, true)
	if err != nil {
		return "", err
	}
	dekey := strings.TrimSpace(string(body))
	if strings.Contains(dekey, ":") {
		return dekey, nil
	}
	return "", fmt.Errorf("invalid key: %s", dekey)
}
