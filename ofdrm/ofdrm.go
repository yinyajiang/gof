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
	"time"

	"github.com/antchfx/xmlquery"
	"github.com/duke-git/lancet/v2/slice"
	widevine "github.com/iyear/gowidevine"
	"github.com/iyear/gowidevine/widevinepb"
	"github.com/yinyajiang/gof/common"
	"github.com/yinyajiang/gof/ofapi"
)

type OFDRM struct {
	req *ofapi.Req
	cfg OFDRMConfig
}

type OFDRMConfig struct {
	ClientID                  []byte
	ClientPrivateKey          []byte
	ClientIDURL               string
	ClientPrivateKeyURL       string
	ClientCacheDir            string
	CachePriority             bool
	OptionalCDRMProjectServer []string
}

func NewOFDRM(req *ofapi.Req, config OFDRMConfig) (*OFDRM, error) {
	if len(config.ClientID) == 0 || len(config.ClientPrivateKey) == 0 {
		clientID, clientPrivateKey, err := LoadClient(config.ClientCacheDir, config.ClientIDURL, config.ClientPrivateKeyURL, config.CachePriority)
		if err != nil {
			return nil, err
		}
		config.ClientID = clientID
		config.ClientPrivateKey = clientPrivateKey
	}
	return &OFDRM{
		req: req,
		cfg: config,
	}, nil
}

func (c *OFDRM) GetVideoDecryptedKeyAuto(drm DRMInfo) (string, error) {
	var clientErr error
	var key string
	if len(c.cfg.ClientID) != 0 {
		key, clientErr = c.GetVideoDecryptedKeyByClient(drm)
		if clientErr == nil {
			return key, nil
		}
		fmt.Println("failed to get decrypted key by client: ", clientErr)
	}

	key, serverErr := c.GetVideoDecryptedKeyByServer(drm)
	if serverErr == nil {
		return key, nil
	}
	fmt.Println("failed to get decrypted key by server: ", serverErr)

	if clientErr != nil {
		return "", clientErr
	}
	return "", serverErr
}

func (c *OFDRM) GetVideoDecryptedKeyByClient(drm DRMInfo) (string, error) {
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

func (c *OFDRM) GetVideoLastModified(drm DRMInfo) (time.Time, error) {
	lastModified, err := c.unsignedGetHeader(drm, "Last-Modified")
	if err != nil {
		fmt.Println("failed to get last modified: ", err)
		return time.Now(), err
	}
	lastModifiedTime, err := time.Parse(time.RFC1123, lastModified)
	if err != nil {
		return time.Now(), fmt.Errorf("parse last modified: %w", err)
	}
	return lastModifiedTime.Local(), nil
}

func (c *OFDRM) GetVideoDecryptedKeyByServer(drm DRMInfo) (string, error) {
	const fixedServerURL = "https://cdrm-project.com/"

	serverURLs := c.cfg.OptionalCDRMProjectServer
	if !slice.Contain(serverURLs, fixedServerURL) {
		serverURLs = append(serverURLs, fixedServerURL)
	}

	pssh, err := c.getDRMPSSH(drm)
	if err != nil {
		return "", err
	}
	maxAttempts := 30
	for i := 0; i < maxAttempts; i++ {
		serverURL := serverURLs[i%len(serverURLs)]
		decryptedKey, err := c.getVideoDecryptedKeyByServer(serverURL, pssh, drm)
		if err == nil {
			return decryptedKey, nil
		}
		fmt.Printf("try %d/%d, failed to get decrypted key from %s: %s\n", i+1, maxAttempts, serverURL, err)
	}
	return "", fmt.Errorf("all servers failed")
}

func (c *OFDRM) getVideoDecryptedKeyByServer(serverURL, pssh string, drm DRMInfo) (string, error) {
	data := common.MustMarshalJSON(map[string]string{
		"PSSH":        pssh,
		"License URL": ofapi.ApiURL(c.drmURLPath(drm)),
		"Headers":     common.MustUnmarshalJSONStr(c.req.SignedHeaders(c.drmURLPath(drm))),
		"JSON":        "",
		"Cookies":     "",
		"Data":        "",
		"Proxy":       "",
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
	content := string(body)
	if strings.Contains(strings.ToLower(content), "error") {
		return "", fmt.Errorf("failed to get decrypted key: %s", content)
	}
	var result map[string]any
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}
	msg, ok := result["Message"]
	if !ok {
		return "", fmt.Errorf("no message")
	}
	return strings.TrimSpace(msg.(string)), nil
}

func (c *OFDRM) drmURLPath(drm DRMInfo) string {
	return ofapi.ApiURLPath("/users/media/%d/drm/post/%d?type=widevine", drm.MediaID, drm.PostID)
}

func (c *OFDRM) getDRMPSSH(drm DRMInfo) (string, error) {
	data, err := c.unsignedGet(drm)
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
		widevine.FromRaw(c.cfg.ClientID, c.cfg.ClientPrivateKey),
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

func (c *OFDRM) unsignedGet(drm DRMInfo) (body []byte, err error) {
	_, body, err = c.unsignedGetResp(drm, true)
	return
}

func (c *OFDRM) unsignedGetHeader(drm DRMInfo, key string) (h string, err error) {
	resp, _, err := c.unsignedGetResp(drm, false)
	if err != nil {
		return "", err
	}
	resp.Body.Close()
	return resp.Header.Get(key), nil
}

func (c *OFDRM) unsignedGetResp(drm DRMInfo, readAll ...bool) (resp *http.Response, body []byte, err error) {
	req, err := http.NewRequest("GET", drm.Drm.Manifest.Dash, nil)
	if err != nil {
		return nil, nil, err
	}
	header := c.req.UnsignedHeaders(map[string]string{
		"Cookie": fmt.Sprintf("CloudFront-Policy=%s; CloudFront-Signature=%s; CloudFront-Key-Pair-Id=%s",
			drm.Drm.Signature.Dash.CloudFrontPolicy,
			drm.Drm.Signature.Dash.CloudFrontSignature,
			drm.Drm.Signature.Dash.CloudFrontKeyPairID),
	})
	common.AddHeaders(req, nil, header)
	return common.HttpDo(req, readAll...)
}
