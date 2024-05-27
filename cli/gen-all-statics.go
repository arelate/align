package cli

import (
	"github.com/arelate/align/paths"
	"github.com/arelate/align/render"
	"github.com/boggydigital/nod"
	"net/url"
)

func GenAllStaticsHandler(u *url.URL) error {
	slug := u.Query().Get("slug")

	return GenAllStatics(slug)
}

func GenAllStatics(slug string) error {

	gasa := nod.NewProgress("generating all statics for %s...", slug)
	defer gasa.End()

	wn, err := render.WikiNavigation(slug)
	if err != nil {
		return gasa.EndWithError(err)
	}

	pages := render.AllLinks(wn)

	skv, err := paths.StaticsKeyValues(slug)
	if err != nil {
		return gasa.EndWithError(err)
	}

	if err := setStaticTOC(slug, skv); err != nil {
		return gasa.EndWithError(err)
	}

	gasa.TotalInt(len(pages))

	for _, page := range pages {

		page, err = url.PathUnescape(page)
		if err != nil {
			return gasa.EndWithError(err)
		}

		if !skv.Has(page) {
			gasa.Increment()
			continue
		}

		if err := setStaticPage(slug, page, skv); err != nil {
			return gasa.EndWithError(err)
		}

		gasa.Increment()
	}

	gasa.EndWithResult("done")

	return nil
}
