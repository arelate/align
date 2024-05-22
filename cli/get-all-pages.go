package cli

import (
	"bytes"
	"encoding/json"
	"github.com/arelate/align/paths"
	"github.com/arelate/southern_light/ign_integration"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
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

	spd, err := paths.AbsPagesSlugDir(slug)
	if err != nil {
		return gapa.EndWithError(err)
	}

	skv, err := kvas.ConnectLocal(spd, kvas.HtmlExt)
	if err != nil {
		return gapa.EndWithError(err)
	}

	rpd, err := paths.AbsDataSlugDir(slug)
	if err != nil {
		return gapa.EndWithError(err)
	}

	rkv, err := kvas.ConnectLocal(rpd, kvas.JsonExt)
	if err != nil {
		return gapa.EndWithError(err)
	}

	//pages := map[string]bool{mainPage: false}
	pages := map[string]bool{mainPage: false}

	for morePages(pages) {
		page := nextPage(pages)

		page, err = url.PathUnescape(page)
		if err != nil {
			return gapa.EndWithError(err)
		}

		urls, err := getUrls(skv, rkv, slug, page, throttle, force)
		if err != nil {
			return gapa.EndWithError(err)
		}

		for _, u := range urls {
			u, err = url.PathUnescape(u)
			if err != nil {
				return gapa.EndWithError(err)
			}
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

func getUrls(skv, rkv kvas.KeyValues, slug, page string, throttle int64, force bool) ([]string, error) {

	gua := nod.Begin("getting page and data for %s...", filepath.Join(slug, page))
	defer gua.End()

	var err error
	var bts []byte

	gotPage, gotData := false, false

	buf := bytes.NewBuffer(bts)

	if !skv.Has(page) || force {
		err = getSetPageContent(skv, slug, page, buf)
		// throttle requests
		time.Sleep(time.Duration(throttle) * time.Millisecond)
	} else {
		gotPage = true
		src, err := skv.Get(page)
		if err != nil {
			return nil, err
		}
		_, err = buf.ReadFrom(src)
		src.Close()
	}
	if err != nil {
		return nil, err
	}

	data := ""

	if !rkv.Has(page) || force {
		data, err = getSetReducedContent(page, buf, rkv)
	} else {
		gotData = true
		rrc, err := rkv.Get(page)
		if err != nil {
			return nil, err
		}

		buf.Reset()
		_, err = buf.ReadFrom(rrc)
		rrc.Close()
		data = buf.String()
	}
	if err != nil {
		return nil, err
	}

	var wikiProps ign_integration.WikiProps
	if err := json.NewDecoder(strings.NewReader(data)).Decode(&wikiProps); err != nil {
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

	result := "done"
	if gotPage && gotData {
		result = "already exist"
	}

	gua.EndWithResult(result)

	return urls, nil
}
