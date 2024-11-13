package ofdl

import (
	"encoding/json"
	"strings"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/gofiber/fiber/v2"
	"github.com/yinyajiang/gof"
)

type ofFiberRoute struct {
	dl     *OFDl
	router fiber.Router
}

func addOFDlFiberRoutes(dl *OFDl, router fiber.Router) {
	r := &ofFiberRoute{
		dl:     dl,
		router: router,
	}
	r.registerRoutes()
}

func (r *ofFiberRoute) registerRoutes() {
	r.router.Get("/extract", r.extract)
	r.router.Post("/extract", r.extract)

	r.router.Get("/fileinfo", r.fileinfo)
	r.router.Post("/fileinfo", r.fileinfo)

	r.router.Get("/drmsecrets", r.drmSecrets)
	r.router.Post("/drmsecrets", r.drmSecrets)

	r.router.Get("/nondrmsecrets", r.nonDrmSecrets)
	r.router.Post("/nondrmsecrets", r.nonDrmSecrets)
}

func (r *ofFiberRoute) extract(c *fiber.Ctx) error {
	var req struct {
		URL          string `json:"url"`
		DisableCache bool   `json:"disable_cache"`
	}
	err := r.bodyUnmarshal(c, &req)
	if err != nil {
		return r.statusError(c, err)
	}
	result, err := r.dl.ExtractMedias(req.URL, req.DisableCache)
	if err != nil {
		return r.statusError(c, err)
	}

	//drm video, drm photo, normal video
	drmVideo := MediaInfo{}
	drmPhoto := MediaInfo{}
	normalVideo := MediaInfo{}
	for _, m := range result.Medias {
		if m.IsDrm {
			if strings.EqualFold(m.Type, "video") && drmVideo.Type == "" {
				drmVideo = m
			} else if (strings.EqualFold(m.Type, "photo") || strings.EqualFold(m.Type, "gif")) && drmPhoto.Type == "" {
				drmPhoto = m
			}
		} else if strings.EqualFold(m.Type, "video") && normalVideo.Type == "" {
			normalVideo = m
		}
	}
	result.Medias = slice.Filter([]MediaInfo{drmVideo, drmPhoto, normalVideo}, func(_ int, m MediaInfo) bool {
		return m.Type != ""
	})

	return r.statusSuccess(c, fiber.Map{
		"ExtractResult": result,
		"Proxy":         gof.ProxyString(),
	})
}

func (r *ofFiberRoute) fileinfo(c *fiber.Ctx) error {
	var req struct {
		MediaURI string `json:"media_uri"`
	}
	err := r.bodyUnmarshal(c, &req)
	if err != nil {
		return r.statusError(c, err)
	}
	info, err := r.dl.FetchFileInfo(req.MediaURI)
	if err != nil {
		return r.statusError(c, err)
	}
	return r.statusSuccess(c, info)
}

func (r *ofFiberRoute) drmSecrets(c *fiber.Ctx) error {
	var req struct {
		MediaURI     string `json:"media_uri"`
		DisableCache bool   `json:"disable_cache"`
	}
	err := r.bodyUnmarshal(c, &req)
	if err != nil {
		return r.statusError(c, err)
	}
	secrets, err := r.dl.FetchDRMSecrets(req.MediaURI, req.DisableCache)
	if err != nil {
		return r.statusError(c, err)
	}
	return r.statusSuccess(c, secrets)
}

func (r *ofFiberRoute) nonDrmSecrets(c *fiber.Ctx) error {
	secrets, err := r.dl.GetNonDRMSecrets()
	if err != nil {
		return r.statusError(c, err)
	}
	return r.statusSuccess(c, secrets)
}

func (r *ofFiberRoute) bodyUnmarshal(c *fiber.Ctx, p any) error {
	body := c.Body()
	return json.Unmarshal(body, p)
}

func (r *ofFiberRoute) statusError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": err.Error()})
}

func (r *ofFiberRoute) statusSuccess(c *fiber.Ctx, p any) error {
	return c.Status(fiber.StatusOK).JSON(p)
}
