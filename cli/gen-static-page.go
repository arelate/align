package cli

import (
	"github.com/boggydigital/nod"
	"net/url"
	"path/filepath"
)

func GenStaticPageHandler(u *url.URL) error {
	slug := u.Query().Get("slug")
	page := u.Query().Get("page")
	force := u.Query().Has("force")

	return GenStaticPage(slug, page, force)
}

func GenStaticPage(slug, page string, force bool) error {

	gspa := nod.Begin("generating static page for %s...", filepath.Join(slug, page))
	defer gspa.End()

	//skv, err := paths.StaticsKeyValues(slug)
	//if err != nil {
	//	return gspa.EndWithError(err)
	//}

	gspa.EndWithResult("done")

	return nil
}
