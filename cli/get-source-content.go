package cli

import (
	"errors"
	"github.com/arelate/align/paths"
	"github.com/arelate/southern_light/ign_integration"
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/kvas_dolo"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathways"
	"golang.org/x/exp/maps"
	"net/url"
	"strings"
)

func GetSourceContentHandler(u *url.URL) error {

	slugs := strings.Split(u.Query().Get("slug"), ",")
	force := u.Query().Has("force")

	return GetSourceContent(force, slugs...)
}

func GetSourceContent(force bool, slugs ...string) error {

	gsca := nod.NewProgress("getting source content...")
	defer gsca.End()

	gsca.TotalInt(len(slugs))

	scd, err := pathways.GetAbsDir(paths.SourceContent)
	if err != nil {
		return gsca.EndWithError(err)
	}

	kv, err := kvas.ConnectLocal(scd, kvas.HtmlExt)
	if err != nil {
		return gsca.EndWithError(err)
	}

	urls := make([]*url.URL, 0, len(slugs))

	for _, slug := range slugs {
		urls = append(urls, ign_integration.WikiUrl(slug))
		gsca.Increment()
	}

	dc := dolo.DefaultClient

	indexSetter := kvas_dolo.NewIndexSetter(kv, slugs...)

	if errs := dc.GetSet(urls, indexSetter, gsca, force); len(errs) > 0 {
		errorUrls := make([]string, 0, len(errs))
		for _, i := range maps.Keys(errs) {
			errorUrls = append(errorUrls, urls[i].String())
		}
		return gsca.EndWithError(errors.New("errors GetSet for URLs: " + strings.Join(errorUrls, ",")))
	}

	gsca.EndWithResult("done")

	return nil
}
