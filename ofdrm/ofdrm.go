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
	"github.com/yinyajiang/gof"
	"github.com/yinyajiang/gof/common"
	"github.com/yinyajiang/gof/ofapi"
)

type OFDRM struct {
	req *ofapi.Req
	cfg OFDRMConfig
}

type OFDRMConfig struct {
	ClientID          []byte
	ClientPrivateKey  []byte
	CDRMProjectServer []string
}

func NewOFDRM(req *ofapi.Req, config OFDRMConfig) *OFDRM {
	return &OFDRM{
		req: req,
		cfg: config,
	}
}

func (c *OFDRM) GetVideoDecryptedKeyAuto(dashVideoURL string) (string, error) {
	useClient := len(c.cfg.ClientID) != 0 && len(c.cfg.ClientPrivateKey) != 0
	useServer := len(c.cfg.CDRMProjectServer) != 0

	if !useClient && !useServer {
		return "", fmt.Errorf("not config client id or private key, and CDRMProjectServer")
	}
	if useClient {
		key, err := c.GetVideoDecryptedKeyByClient(dashVideoURL)
		if err == nil {
			return key, nil
		}
		fmt.Println("failed to get decrypted key by client: ", err)
	}
	if useServer {
		return c.GetVideoDecryptedKeyByServer(dashVideoURL)
	}
	return "", nil
}

func (c *OFDRM) GetVideoDecryptedKeyByClient(dashVideoURL string) (string, error) {
	mpdInfo, err := common.ParseVideoMPDInfo(dashVideoURL)
	if err != nil {
		return "", err
	}
	pssh, err := c.getDRMPSSH(mpdInfo)
	if err != nil {
		return "", err
	}

	keys, err := c.getWidevineKeys(c.drmURLPath(mpdInfo), pssh)
	if err != nil {
		return "", err
	}

	key := keys[0]
	decryptedKey := strings.ToLower(hex.EncodeToString(key.ID)) + ":" + strings.ToLower(hex.EncodeToString(key.Key))
	return decryptedKey, nil
}

func (c *OFDRM) GetVideoLastModified(dashVideoURL string) (time.Time, error) {
	mpdInfo, err := common.ParseVideoMPDInfo(dashVideoURL)
	if err != nil {
		return time.Time{}, err
	}
	header, err := c.req.MPDGetHeader(mpdInfo)
	if err != nil {
		return time.Time{}, err
	}
	lastModified := header.Get("Last-Modified")
	if lastModified == "" {
		return time.Now(), nil
	}
	lastModifiedTime, err := time.Parse(time.RFC1123, lastModified)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse last modified: %w", err)
	}
	return lastModifiedTime.Local(), nil
}

func (c *OFDRM) GetVideoDecryptedKeyByServer(dashVideoURL string) (string, error) {
	mpdInfo, err := common.ParseVideoMPDInfo(dashVideoURL)
	if err != nil {
		return "", err
	}
	pssh, err := c.getDRMPSSH(mpdInfo)
	if err != nil {
		return "", err
	}
	maxAttempts := 30
	for i := 0; i < maxAttempts; i++ {
		serverURL := c.cfg.CDRMProjectServer[i%len(c.cfg.CDRMProjectServer)]
		decryptedKey, err := c.getVideoDecryptedKeyByServer(serverURL, pssh, mpdInfo)
		if err == nil {
			return decryptedKey, nil
		}
		fmt.Printf("try %d/%d, failed to get decrypted key from %s: %s\n", i+1, maxAttempts, serverURL, err)
	}
	return "", fmt.Errorf("all servers failed")
}

func (c *OFDRM) getVideoDecryptedKeyByServer(serverURL, pssh string, mpdInfo gof.MPDURLInfo) (string, error) {
	data := common.MustMarshalJSON(map[string]string{
		"PSSH":        pssh,
		"License URL": ofapi.ApiURL(c.drmURLPath(mpdInfo)),
		"Headers":     string(common.MustMarshalJSON(c.req.AuthHeaders(c.drmURLPath(mpdInfo)))),
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

func (c *OFDRM) drmURLPath(mpdInfo gof.MPDURLInfo) string {
	return ofapi.ApiURLPath("/users/media/%s/drm/post/%s?type=widevine", mpdInfo.MediaID, mpdInfo.PostID)
}

func (c *OFDRM) getDRMPSSH(mpdInfo gof.MPDURLInfo) (string, error) {
	data, err := c.req.MPDGet(mpdInfo)
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
