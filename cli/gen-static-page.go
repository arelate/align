package cli

import (
	"github.com/arelate/align/paths"
	"github.com/arelate/align/render"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"net/url"
	"path/filepath"
	"strings"
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

	skv, err := paths.StaticsKeyValues(slug)
	if err != nil {
		return gspa.EndWithError(err)
	}

	if skv.Has(page) && !force {
		gspa.EndWithResult("already exists")
		return nil
	}

	if err := setStaticPage(slug, page, skv); err != nil {
		return gspa.EndWithError(err)
	}

	gspa.EndWithResult("done")

	return nil
}

func setStaticPage(slug, page string, kv kvas.KeyValues) error {
	sb := new(strings.Builder)

	if err := render.WikisSlugPage(slug, page, sb); err != nil {
		return err
	}

	return kv.Set(page, strings.NewReader(sb.String()))
}
