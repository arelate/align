package cli

import (
	"github.com/arelate/align/paths"
	"github.com/arelate/southern_light/ign_integration"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathways"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
)

const (
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Safari/605.1.15"
	mainPage  = "Main_Page"
)

func GetSourcePageHandler(u *url.URL) error {

	slug := u.Query().Get("slug")
	page := u.Query().Get("page")
	force := u.Query().Has("force")

	return GetSourcePage(slug, page, force)
}

func GetSourcePage(slug, page string, force bool) error {

	if page == "" {
		page = mainPage
	}

	gsca := nod.Begin("getting source page %s...", path.Join(slug, page))
	defer gsca.End()

	spd, err := pathways.GetAbsDir(paths.SourcePages)
	if err != nil {
		return gsca.EndWithError(err)
	}

	spd = filepath.Join(spd, slug)

	kv, err := kvas.ConnectLocal(spd, kvas.HtmlExt)
	if err != nil {
		return gsca.EndWithError(err)
	}

	if kv.Has(page) && !force {
		gsca.EndWithResult("already exist")
		return nil
	}

	u := ign_integration.WikiUrl(slug, page)

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return gsca.EndWithError(err)
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return gsca.EndWithError(err)
	}
	defer resp.Body.Close()

	if err := kv.Set(page, resp.Body); err != nil {
		return gsca.EndWithError(err)
	}

	gsca.EndWithResult("done")

	return nil
}
