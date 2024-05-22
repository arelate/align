package cli

import (
	"encoding/json"
	"errors"
	"github.com/arelate/align/paths"
	"github.com/arelate/southern_light/ign_integration"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"net/url"
	"strconv"
	"time"
)

const (
	defaultThrottleMs = int64(1000)
)

func GetAllPagesHandler(u *url.URL) error {
	q := u.Query()

	slug := q.Get("slug")
	force := q.Has("force")
	var throttle = defaultThrottleMs
	if tstr := q.Get("throttle"); tstr != "" {
		if ti, err := strconv.ParseInt(tstr, 10, 64); err == nil {
			throttle = ti
		}
	}

	return GetAllPages(slug, throttle, force)
}

func GetAllPages(slug string, throttle int64, force bool) error {

	gapa := nod.Begin("getting all pages for %s...", slug)
	defer gapa.End()

	sdd, err := paths.AbsDataSlugDir(slug)
	if err != nil {
		return gapa.EndWithError(err)
	}

	kv, err := kvas.ConnectLocal(sdd, kvas.JsonExt)
	if err != nil {
		return gapa.EndWithError(err)
	}

	pages := map[string]bool{mainPage: false}

	for morePages(pages) {
		page := nextPage(pages)

		page, err = url.PathUnescape(page)
		if err != nil {
			return gapa.EndWithError(err)
		}

		urls, err := getUrls(kv, slug, page, throttle, force)
		if err != nil {
			return gapa.EndWithError(err)
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

func getUrls(kv kvas.KeyValues, slug, page string, throttle int64, force bool) ([]string, error) {
	if err := GetPage(slug, page, force); err != nil {
		return nil, err
	}

	if err := GetData(slug, page, force); err != nil {
		return nil, err
	}

	// throttle requests
	time.Sleep(time.Duration(throttle) * time.Millisecond)

	if err := kv.IndexRefresh(); err != nil {
		return nil, err
	}

	wikiPage, err := kv.Get(page)
	if err != nil {
		return nil, err
	}

	if wikiPage == nil {
		return nil, errors.New("page not found: " + page)
	}

	defer wikiPage.Close()

	var wikiProps ign_integration.WikiProps
	if err := json.NewDecoder(wikiPage).Decode(&wikiProps); err != nil {
		return nil, err
	}

	urls := make([]string, 0)
	for _, he := range wikiProps.HTMLEntities() {
		purls, err := he.PageUrls(slug)
		if err != nil {
			return nil, err
		}
		if npu := wikiProps.NextPageUrl(); npu != "" {
			purls = append(purls, npu)
		}
		if ppu := wikiProps.PreviousPageUrl(); ppu != "" {
			purls = append(purls, ppu)
		}
		urls = append(urls, purls...)
	}

	return urls, nil
}
