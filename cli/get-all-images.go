package cli

import (
	"github.com/arelate/align/paths"
	"github.com/boggydigital/nod"
	"net/url"
)

func GetAllImagesHandler(u *url.URL) error {
	q := u.Query()

	slug := q.Get("slug")
	force := q.Has("force")

	return GetAllImages(slug, force)
}

func GetAllImages(slug string, force bool) error {

	gaia := nod.NewProgress("getting all images for %s...", slug)
	defer gaia.End()

	kv, err := paths.DataKeyValues(slug)
	if err != nil {
		return gaia.EndWithError(err)
	}

	pages := kv.Keys()

	gaia.TotalInt(len(pages))

	for _, page := range pages {

		page, err = url.PathUnescape(page)
		if err != nil {
			return gaia.EndWithError(err)
		}

		if err := GetImages(kv, slug, page, force); err != nil {
			return gaia.EndWithError(err)
		}

		gaia.Increment()
	}

	gaia.EndWithResult("done")
	return nil
}
