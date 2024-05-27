package cli

import (
	"github.com/arelate/align/paths"
	"github.com/arelate/align/render/view_models"
	"github.com/arelate/southern_light/ign_integration"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"io"
	"net/http"
	"net/url"
	"path"
)

const (
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36"
)

func GetPageHandler(u *url.URL) error {
	slug := u.Query().Get("slug")
	page := u.Query().Get("page")
	force := u.Query().Has("force")

	return GetPage(slug, page, force)
}

func GetPage(slug, page string, force bool) error {

	if page == "" {
		page = view_models.MainPage
	}

	gsca := nod.Begin("getting source page %s...", path.Join(slug, page))
	defer gsca.End()

	page, err := url.PathUnescape(page)
	if err != nil {
		return gsca.EndWithError(err)
	}

	kv, err := paths.PagesKeyValues(slug)
	if err != nil {
		return gsca.EndWithError(err)
	}

	if kv.Has(page) && !force {
		gsca.EndWithResult("already exist")
		return nil
	}

	if err := getSetPageContent(kv, slug, page, nil); err != nil {
		return gsca.EndWithError(err)
	}

	gsca.EndWithResult("done")

	return nil
}

func getSetPageContent(kv kvas.KeyValues, slug, page string, dst io.Writer) error {

	u := ign_integration.WikiPageUrl(slug, page)

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// only save successful results
	if resp.StatusCode != http.StatusOK {
		return nil
	}

	var rdr io.Reader = resp.Body
	if dst != nil {
		rdr = io.TeeReader(resp.Body, dst)
	}

	return kv.Set(page, rdr)
}
