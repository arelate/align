package cli

import (
	"encoding/json"
	"fmt"
	"github.com/arelate/align/paths"
	"github.com/arelate/southern_light/ign_integration"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"net/url"
	"time"
)

func GetAllPagesHandler(u *url.URL) error {
	slug := u.Query().Get("slug")
	force := u.Query().Has("force")

	return GetAllPages(slug, force)
}

func GetAllPages(slug string, force bool) error {

	gapa := nod.Begin("getting all pages for %s...", slug)
	defer gapa.End()

	sdd, err := paths.AbsDataDir(slug)
	if err != nil {
		return gapa.EndWithError(err)
	}

	kv, err := kvas.ConnectLocal(sdd, kvas.JsonExt)
	if err != nil {
		return gapa.EndWithError(err)
	}

	pages := map[string]bool{mainPage: false}

	var wikiProps ign_integration.WikiProps

	for morePages(pages) {

		page := nextPage(pages)

		if err := GetPage(slug, page, force); err != nil {
			return gapa.EndWithError(err)
		}

		if err := GetData(slug, page, force); err != nil {
			return gapa.EndWithError(err)
		}

		time.Sleep(1000 * time.Millisecond)

		if err := kv.IndexRefresh(); err != nil {
			return gapa.EndWithError(err)
		}

		wikiPage, err := kv.Get(page)
		if err != nil {
			return gapa.EndWithError(err)
		}

		if wikiPage == nil {
			return gapa.EndWithError(fmt.Errorf("page %s not found", page))
		}

		if err := json.NewDecoder(wikiPage).Decode(&wikiProps); err != nil {
			return gapa.EndWithError(err)
		}

		wikiPage.Close()

		urls := make([]string, 0)
		for _, he := range wikiProps.HTMLEntities() {
			purls, err := he.PageUrls(slug)
			if err != nil {
				return gapa.EndWithError(err)
			}
			if npu := wikiProps.NextPageUrl(); npu != "" {
				purls = append(purls, npu)
			}
			if ppu := wikiProps.PreviousPageUrl(); ppu != "" {
				purls = append(purls, ppu)
			}
			urls = append(urls, purls...)
		}

		for _, u := range urls {
			if got := pages[u]; !got {
				pages[u] = false
			}
		}

		pages[page] = true
	}

	gapa.EndWithResult("done")
	return nil
}

func nextPage(pages map[string]bool) string {
	for p, g := range pages {
		if !g {
			return p
		}
	}
	return ""
}

func morePages(pages map[string]bool) bool {
	for _, g := range pages {
		if !g {
			return true
		}
	}
	return false
}
