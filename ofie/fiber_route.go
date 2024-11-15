package ofie

import (
	"encoding/json"
	"strings"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/gofiber/fiber/v2"
	"github.com/yinyajiang/gof"
)

type ofFiberRoute struct {
	ie     *OFIE
	router fiber.Router
}

func addOFIEFiberRoutes(ie *OFIE, router fiber.Router) {
	r := &ofFiberRoute{
		ie:     ie,
		router: router,
	}
	r.registerRoutes()
}

func (r *ofFiberRoute) registerRoutes() {
	r.router.Get("/of/extract", r.extract)
	r.router.Post("/of/extract", r.extract)

	r.router.Get("/of/fileinfo", r.fileinfo)
	r.router.Post("/of/fileinfo", r.fileinfo)

	r.router.Get("/of/drmsecrets", r.drmSecrets)
	r.router.Post("/of/drmsecrets", r.drmSecrets)

	r.router.Get("/of/nondrmsecrets", r.nonDrmSecrets)
	r.router.Post("/of/nondrmsecrets", r.nonDrmSecrets)
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

	if len(req.MediaFilter) > 0 {
		filter := map[string]struct{}{}
		for _, t := range req.MediaFilter {
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
		filtered := slice.Filter(result.Medias, func(_ int, m MediaInfo) bool {
			ty := ""
			if m.IsDrm {
				ty = "drm-" + strings.ToLower(m.Type)
			} else {
				ty = "non-drm-" + strings.ToLower(m.Type)
			}
			_, ok := filter[ty]
			return ok
		})
		if len(filtered) != 0 {
			result.Medias = filtered
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
	secrets, err := r.ie.FetchDRMSecrets(req.MediaURI, req.DisableCache)
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
