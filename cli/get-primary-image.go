package cli

import (
	"encoding/json"
	"github.com/arelate/align/paths"
	"github.com/arelate/align/render/view_models"
	"github.com/arelate/southern_light/ign_integration"
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathways"
	"net/url"
)

func GetPrimaryImageHandler(u *url.URL) error {
	q := u.Query()

	slug := q.Get("slug")
	force := q.Has("force")

	return GetPrimaryImage(nil, slug, force)
}

func GetPrimaryImage(kv kvas.KeyValues, slug string, force bool) error {

	gpia := nod.NewProgress("getting primary image for %s...", slug)
	defer gpia.End()

	var err error
	if kv == nil {
		kv, err = paths.DataKeyValues(slug)
		if err != nil {
			return gpia.EndWithError(err)
		}
	}

	wp, err := kv.Get(view_models.MainPage)
	if err != nil {
		return gpia.EndWithError(err)
	}
	defer wp.Close()

	var wikiProps ign_integration.WikiProps

	if err := json.NewDecoder(wp).Decode(&wikiProps); err != nil {
		return gpia.EndWithError(err)
	}

	piu, err := url.Parse(wikiProps.PrimaryImageUrl())
	if err != nil {
		return gpia.EndWithError(err)
	}

	spid, err := pathways.GetAbsRelDir(paths.PrimaryImages)
	if err != nil {
		return gpia.EndWithError(err)
	}

	dc := dolo.DefaultClient
	rp := piu.Path

	if err := dc.Download(piu, force, nil, spid, rp); err != nil {
		return gpia.EndWithError(err)
	}

	gpia.EndWithResult("done")

	return nil
}
