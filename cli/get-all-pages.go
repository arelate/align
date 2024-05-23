package cli

import (
	"bytes"
	"encoding/json"
	"github.com/arelate/align/paths"
	"github.com/arelate/southern_light/ign_integration"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathways"
	"io"
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

// getUrls implements getting page, data, navigation in one function
// it uses the same principal code as relevant cmds, and supports
// getting existing data instead of fetching it every time.
// Here's the sequence:
// 1) Get page - from the storage if it exists or origin
// 2) Get data - from the storage if it exists or extracted
// from the buffered data result of getting page
// 3) Get navigation - only for the main page and only from results
// 4) Decode data JSON and get all the links - previous, next pages and <a href>
func getUrls(skv, rkv kvas.KeyValues, slug, page string, throttle int64, force bool) ([]string, error) {

	gua := nod.Begin("getting page and data for %s...", filepath.Join(slug, page))
	defer gua.End()

	var err error
	var sr io.Reader

	gotPage, gotData := false, false

	if !skv.Has(page) || force {
		buf := bytes.NewBuffer(make([]byte, 0, 512))
		err = getSetPageContent(skv, slug, page, buf)
		sr = buf
		// throttle requests
		time.Sleep(time.Duration(throttle) * time.Millisecond)
	} else {
		gotPage = true
		sr, err = skv.Get(page)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}
	if src, ok := sr.(io.ReadCloser); ok {
		defer src.Close()
	}

	data := ""

	if !rkv.Has(page) || force {
		data, err = getSetReducedContent(page, sr, rkv)
	} else {
		gotData = true
		rrc, err := rkv.Get(page)
		if err != nil {
			return nil, err
		}
		defer rrc.Close()
		sb := new(strings.Builder)
		_, err = io.Copy(sb, rrc)
		data = sb.String()

	}
	if err != nil {
		return nil, err
	}

	if page == mainPage && data != "" {

		snd, err := pathways.GetAbsDir(paths.Navigation)
		if err != nil {
			return nil, err
		}

		nkv, err := kvas.ConnectLocal(snd, kvas.JsonExt)
		if err != nil {
			return nil, err
		}

		if err := getSetNavigation(slug, data, nkv); err != nil {
			return nil, err
		}
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
