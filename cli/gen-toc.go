package cli

import (
	"github.com/arelate/align/paths"
	"github.com/arelate/align/render"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"net/url"
	"strings"
)

func GenTOCHandler(u *url.URL) error {
	slug := u.Query().Get("slug")
	force := u.Query().Has("force")

	return GenTOC(slug, force)
}

func GenTOC(slug string, force bool) error {
	gta := nod.Begin("generating static TOC for %s...", slug)
	defer gta.End()

	skv, err := paths.StaticsKeyValues(slug)
	if err != nil {
		return gta.EndWithError(err)
	}

	if skv.Has(slug) && !force {
		gta.EndWithResult("already exists")
		return nil
	}

	if err := setStaticTOC(slug, skv); err != nil {
		return gta.EndWithError(err)
	}

	gta.EndWithResult("done")

	return nil
}

func setStaticTOC(slug string, kv kvas.KeyValues) error {
	sb := new(strings.Builder)

	if err := render.WikisPage(slug, sb); err != nil {
		return err
	}

	return kv.Set(slug, strings.NewReader(sb.String()))
}
