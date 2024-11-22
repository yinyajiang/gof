package ofie

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/gofiber/fiber/v2"
	"github.com/yinyajiang/gof"
	"github.com/yinyajiang/gof/common"
	"github.com/yinyajiang/gof/ofapi"
)

const AUTH_PATH = "/of/auth"
const EXTRACT_PATH = "/of/extract"
const FILEINFO_PATH = "/of/fileinfo"
const DRM_SECRETS_PATH = "/of/drmsecrets"
const NON_DRM_SECRETS_PATH = "/of/nondrmsecrets"

type ofFiberRoute struct {
	ie           *OFIE
	router       fiber.Router
	preferFilter []string
}

// video, audio, photo, drm-video, drm-photo, drm-audio, non-drm-video, non-drm-photo, non-drm-audio
func addOFIEFiberRoutes(ie *OFIE, router fiber.Router, preferFilter ...string) {
	r := &ofFiberRoute{
		ie:           ie,
		router:       router,
		preferFilter: common.CleanEmptryString(preferFilter),
	}
	r.registerRoutes()
}

func (r *ofFiberRoute) registerRoutes() {
	r.router.Get(EXTRACT_PATH, r.extract)
	r.router.Post(EXTRACT_PATH, r.extract)

	r.router.Get(FILEINFO_PATH, r.fileinfo)
	r.router.Post(FILEINFO_PATH, r.fileinfo)

	r.router.Get(DRM_SECRETS_PATH, r.drmSecrets)
	r.router.Post(DRM_SECRETS_PATH, r.drmSecrets)

	r.router.Get(NON_DRM_SECRETS_PATH, r.nonDrmSecrets)
	r.router.Post(NON_DRM_SECRETS_PATH, r.nonDrmSecrets)

	r.router.Get(AUTH_PATH, r.auth)
	r.router.Post(AUTH_PATH, r.auth)
}

func (r *ofFiberRoute) extract(c *fiber.Ctx) error {
	var req struct {
		URL          string
		DisableCache bool
		MediaFilter  []string //video, audio, photo, drm-video, drm-photo, drm-audio, non-drm-video, non-drm-photo, non-drm-audio
	}
	err := r.bodyUnmarshal(c, &req)
	if err != nil {
		return r.statusError(c, err)
	}
	result, err := r.ie.ExtractMedias(req.URL, ExtractOption{
		DisableCache: req.DisableCache,
	})
	if err != nil {
		return r.statusError(c, err)
	}

	filterArr := common.CleanEmptryString(req.MediaFilter)
	if len(filterArr) == 0 {
		filterArr = r.preferFilter
	}

	if len(filterArr) > 0 {
		filter := map[string]struct{}{}
		for _, t := range filterArr {
			if strings.EqualFold(t, "video") {
				filter["drm-video"] = struct{}{}
				filter["non-drm-video"] = struct{}{}
			} else if strings.EqualFold(t, "photo") {
				filter["drm-photo"] = struct{}{}
				filter["non-drm-photo"] = struct{}{}
				filter["drm-gif"] = struct{}{}
				filter["non-drm-gif"] = struct{}{}
			} else if strings.EqualFold(t, "audio") {
				filter["drm-audio"] = struct{}{}
				filter["non-drm-audio"] = struct{}{}
			} else {
				filter[strings.ToLower(t)] = struct{}{}
			}
		}
		filteredMedias := slice.Filter(result.Medias, func(_ int, m MediaInfo) bool {
			ty := ""
			if m.IsDrm {
				ty = "drm-" + strings.ToLower(m.Type)
			} else {
				ty = "non-drm-" + strings.ToLower(m.Type)
			}
			_, ok := filter[ty]
			return ok
		})
		if len(filteredMedias) != 0 {
			result.Medias = filteredMedias
		}
	}

	return r.statusSuccess(c, fiber.Map{
		"ExtractResult": result,
		"Proxy":         gof.ProxyString(),
	})
}

func (r *ofFiberRoute) fileinfo(c *fiber.Ctx) error {
	var req struct {
		MediaURI string
	}
	err := r.bodyUnmarshal(c, &req)
	if err != nil {
		return r.statusError(c, err)
	}
	info, err := r.ie.FetchFileInfo(req.MediaURI)
	if err != nil {
		return r.statusError(c, err)
	}
	return r.statusSuccess(c, info)
}

func (r *ofFiberRoute) drmSecrets(c *fiber.Ctx) error {
	var req struct {
		MediaURI     string
		DisableCache bool
	}
	err := r.bodyUnmarshal(c, &req)
	if err != nil {
		return r.statusError(c, err)
	}
	secrets, err := r.ie.FetchDRMSecrets(req.MediaURI, FetchDRMSecretsOption{
		DisableCache: req.DisableCache,
	})
	if err != nil {
		return r.statusError(c, err)
	}
	return r.statusSuccess(c, secrets)
}

func (r *ofFiberRoute) nonDrmSecrets(c *fiber.Ctx) error {
	secrets, err := r.ie.GetNonDRMSecrets()
	if err != nil {
		return r.statusError(c, err)
	}
	return r.statusSuccess(c, secrets)
}

func (r *ofFiberRoute) auth(c *fiber.Ctx) error {
	body := c.Body()
	var authInfo ofapi.OFAuthInfo
	err := json.Unmarshal(body, &authInfo)
	if err == nil {
		err = r.ie.api.Auth(authInfo)
	} else {
		err = r.ie.api.AuthByString(string(body))
	}
	if err != nil {
		return r.statusError(c, err)
	}
	return r.statusSuccess(c, nil)
}

func (r *ofFiberRoute) bodyUnmarshal(c *fiber.Ctx, p any) error {
	body := c.Body()
	return json.Unmarshal(body, p)
}

func (r *ofFiberRoute) statusError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": err.Error()})
}

func (r *ofFiberRoute) statusSuccess(c *fiber.Ctx, p any) error {
	if p == nil {
		return c.Status(fiber.StatusOK).SendString("success")
	}
	return c.Status(fiber.StatusOK).JSON(p)
}

//utils

type OFClientHelper struct {
	ServerAddr string
}

func (h *OFClientHelper) Auth(authInfo any) error {
	by, err := h.post(AUTH_PATH, authInfo)
	if err != nil {
		return err
	}
	if strings.Contains(string(by), "error") {
		return errors.New(string(by))
	}
	return nil
}

func (h *OFClientHelper) post(path string, p any) ([]byte, error) {
	var body []byte
	var err error
	if p != nil {
		switch p := p.(type) {
		case string:
			body = []byte(p)
		case []byte:
			body = p
		default:
			body, err = json.Marshal(p)
			if err != nil {
				return nil, err
			}
		}
	}

	if h.ServerAddr == "" {
		return nil, errors.New("serverAddr is empty")
	}
	serverAddr := strings.TrimSuffix(h.ServerAddr, "/")
	resp, err := http.Post(serverAddr+path, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
