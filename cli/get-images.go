package cli

import (
	"encoding/json"
	"github.com/arelate/align/paths"
	"github.com/arelate/southern_light/ign_integration"
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"net/url"
	"path/filepath"
	"strings"
)

func GetImagesHandler(u *url.URL) error {
	q := u.Query()

	slug := q.Get("slug")
	page := q.Get("page")
	force := q.Has("force")

	return GetImages(nil, slug, page, force)
}

func GetImages(kv kvas.KeyValues, slug, page string, force bool) error {

	gia := nod.NewProgress("getting images for %s...", filepath.Join(slug, page))
	defer gia.End()

	var err error
	if kv == nil {
		kv, err = paths.DataKeyValues(slug)
		if err != nil {
			return gia.EndWithError(err)
		}
	}

	wp, err := kv.Get(page)
	if err != nil {
		return gia.EndWithError(err)
	}
	defer wp.Close()

	var wikiProps ign_integration.WikiProps

	if err := json.NewDecoder(wp).Decode(&wikiProps); err != nil {
		return gia.EndWithError(err)
	}

	imageUrls := make(map[string]any)

	for _, he := range wikiProps.HTMLEntities() {
		ius, err := he.ImageUrls()
		if err != nil {
			return gia.EndWithError(err)
		}

		for _, iu := range ius {
			imageUrls[iu] = nil
		}
	}

	gia.TotalInt(len(imageUrls))

	dc := dolo.DefaultClient

	sid, err := paths.AbsImagesSlugDir(slug)
	if err != nil {
		return gia.EndWithError(err)
	}

	for imageUrl := range imageUrls {

		u, err := url.Parse(imageUrl)
		if err != nil {
			return gia.EndWithError(err)
		}

		// download without query params
		u.RawQuery = ""

		if rp := relPath(u.Path, slug); rp != "" {
			if err := dc.Download(u, force, nil, sid, rp); err != nil {
				gia.Error(err)
				continue
			}
		}

		gia.Increment()
	}

	gia.EndWithResult("done")
	return nil
}

func relPath(path, slug string) string {
	if _, rp, ok := strings.Cut(path, slug); ok {
		return rp
	} else {
		return ""
	}
}
