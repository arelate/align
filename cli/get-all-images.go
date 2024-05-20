package cli

import (
	"github.com/arelate/align/paths"
	"github.com/boggydigital/kvas"
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

	sdd, err := paths.AbsDataSlugDir(slug)
	if err != nil {
		return gaia.EndWithError(err)
	}

	kv, err := kvas.ConnectLocal(sdd, kvas.JsonExt)
	if err != nil {
		return gaia.EndWithError(err)
	}

	pages := kv.Keys()

	gaia.TotalInt(len(pages))

	for _, page := range pages {

		if err := GetImages(kv, slug, page, force); err != nil {
			return gaia.EndWithError(err)
		}

		gaia.Increment()
	}

	gaia.EndWithResult("done")
	return nil
}
