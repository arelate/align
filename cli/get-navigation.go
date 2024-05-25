package cli

import (
	"errors"
	"github.com/arelate/align/paths"
	"github.com/arelate/align/render/view_models"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathways"
	"io"
	"net/url"
	"strings"
)

var (
	ErrNavigationNotFound = errors.New("navigation not found")
)

func GetNavigationHandler(u *url.URL) error {
	slug := u.Query().Get("slug")
	force := u.Query().Has("force")

	return GetNavigation(slug, force)
}

func GetNavigation(slug string, force bool) error {

	gna := nod.Begin("getting navigation for %s...", slug)
	defer gna.End()

	snd, err := pathways.GetAbsDir(paths.Navigation)
	if err != nil {
		return gna.EndWithError(err)
	}

	nkv, err := kvas.NewKeyValues(snd, kvas.JsonExt)
	if err != nil {
		return gna.EndWithError(err)
	}

	if nkv.Has(slug) && !force {
		gna.EndWithResult("already exist")
		return nil
	}

	dkv, err := paths.DataKeyValues(slug)
	if err != nil {
		return gna.EndWithError(err)
	}

	mprc, err := dkv.Get(view_models.MainPage)
	if err != nil {
		return gna.EndWithError(err)
	}
	defer mprc.Close()

	sb := new(strings.Builder)
	if _, err := io.Copy(sb, mprc); err != nil {
		return gna.EndWithError(err)
	}

	if err := getSetNavigation(slug, sb.String(), nkv); err != nil {
		return gna.EndWithError(err)
	}

	gna.EndWithResult("done")

	return nil
}

func getSetNavigation(slug string, data string, kv kvas.KeyValues) error {
	if _, rem, ok := strings.Cut(data, "\"navigation\":"); ok {
		if nav, _, ok := strings.Cut(rem, ",\"videos:"); ok {
			return kv.Set(slug, strings.NewReader(nav))
		}
	}
	return ErrNavigationNotFound
}
