package cli

import (
	"encoding/json"
	"github.com/arelate/align/data"
	"github.com/arelate/align/nav"
	"github.com/arelate/align/paths"
	"github.com/arelate/align/render/view_models"
	"github.com/arelate/southern_light/ign_integration"
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/nod"
	"net/url"
	"path"
	"strings"
)

func ReduceHandler(u *url.URL) error {
	slugs := strings.Split(u.Query().Get("slug"), ",")
	all := u.Query().Has("all")
	return Reduce(all, slugs...)
}

func Reduce(all bool, slugs ...string) error {

	ra := nod.NewProgress("reducing data...")
	defer ra.End()

	if all {
		nkv, err := paths.NavigationKeyValues()
		if err != nil {
			return ra.EndWithError(err)
		}
		slugs, err = nkv.Keys()
		if err != nil {
			return ra.EndWithError(err)
		}
	}

	ra.TotalInt(len(slugs))

	for _, slug := range slugs {
		if err := reduceSlug(slug); err != nil {
			return ra.EndWithError(err)
		}
		ra.Increment()
	}

	ra.EndWithResult("done")

	return nil
}

func reduceSlug(slug string) error {

	rsa := nod.NewProgress(" %s", slug)
	defer rsa.End()

	reductions := make(map[string]map[string][]string)
	for _, p := range data.AllReduxProperties() {
		reductions[p] = make(map[string][]string)
	}

	rdx, err := paths.NewReduxWriter()
	if err != nil {
		return rsa.EndWithError(err)
	}

	dkv, err := paths.DataKeyValues(slug)
	if err != nil {
		return rsa.EndWithError(err)
	}

	// wiki

	reductions[data.WikiPages][slug], err = dkv.Keys()
	if err != nil {
		return rsa.EndWithError(err)
	}

	mainPage, err := getWikiPage(view_models.MainPage, dkv)
	if err != nil {
		return rsa.EndWithError(err)
	}

	reductions[data.WikiNameProperty][slug] = []string{mainPage.Props.PageProps.Page.Name}
	reductions[data.WikiPrimaryImageProperty][slug] = []string{mainPage.PrimaryImageUrl()}

	// pages

	pages, err := dkv.Keys()
	if err != nil {
		return rsa.EndWithError(err)
	}

	rsa.TotalInt(len(pages))

	for _, page := range pages {

		wp, err := getWikiPage(page, dkv)
		if err != nil {
			return rsa.EndWithError(err)
		}

		sp := path.Join(slug, page)

		reductions[data.PageTitleProperty][sp] = []string{wp.PageTitle()}
		reductions[data.PageNextPageUrlProperty][sp] = []string{wp.NextPageUrl()}
		reductions[data.PagePrevPageUrlProperty][sp] = []string{wp.PreviousPageUrl()}
		reductions[data.PagePublishDateProperty][sp] = []string{wp.PublishDate().Format("Jan 2, 2006")}
		reductions[data.PageUpdatedAtProperty][sp] = []string{wp.UpdatedAt().Format("Jan 2, 2006")}

		htmlEntities := make([]string, 0)
		for _, he := range wp.HTMLEntities() {
			content := he.Values.Html

			content = rewriteOriginLinks(content)
			content = rewriteImageLinks(content)
			content = disableStyles(content)

			if content != "" {
				htmlEntities = append(htmlEntities, content)
			}

			imagesContent := make([]string, 0, len(he.ImageValues))
			for _, iv := range he.ImageValues {
				imagesContent = append(imagesContent, "<img src='"+rewriteImageLinks(iv.Original)+"' />")
			}

			htmlEntities = append(htmlEntities, imagesContent...)
		}

		reductions[data.PageHTMLEntriesProperty][sp] = htmlEntities

		rsa.Increment()
	}

	// navigation

	wikiNavigation, err := nav.WikiNavigation(slug)
	if err != nil {
		return rsa.EndWithError(err)
	}

	wikiNav := make([]string, 0, len(wikiNavigation))
	navTitle := ""

	for _, wn := range wikiNavigation {
		if navTitle == "" {
			navTitle = wn.Label
		}
		setNavigationSubNav(slug, &wn, reductions)
		wikiNav = append(wikiNav, wn.Url)
	}

	reductions[data.NavigationTitleProperty][slug] = []string{navTitle}
	reductions[data.NavigationProperty][slug] = wikiNav

	for _, link := range nav.AllLinks(wikiNavigation) {
		link, err = url.PathUnescape(link)
		if err != nil {
			return rsa.EndWithError(err)
		}
		has, err := dkv.Has(link)
		if err != nil {
			return rsa.EndWithError(err)
		}
		if !has {
			reductions[data.PageMissingProperty][slug] = append(reductions[data.PageMissingProperty][slug], link)
		}
	}

	for property := range reductions {
		if err := rdx.BatchReplaceValues(property, reductions[property]); err != nil {
			return rsa.EndWithError(err)
		}
	}

	rsa.EndWithResult("done")

	return nil
}

func getWikiPage(page string, kv kevlar.KeyValues) (*ign_integration.WikiProps, error) {
	wikiPage, err := kv.Get(page)
	if err != nil {
		return nil, err
	}
	defer wikiPage.Close()

	var wikiProps ign_integration.WikiProps

	err = json.NewDecoder(wikiPage).Decode(&wikiProps)
	return &wikiProps, err
}

func rewriteImageLinks(html string) string {
	return strings.Replace(html, "https://oyster.ignimgs.com/mediawiki/apis.ign.com", "/image", -1)
}

func rewriteOriginLinks(html string) string {
	return strings.Replace(html, "https://www.ign.com", "", -1)
}

func disableStyles(html string) string {
	return strings.Replace(html, "style=", "data-style=", -1)
}

func setNavigationSubNav(slug string, wn *ign_integration.WikiNavigation, reductions map[string]map[string][]string) {

	su := path.Join(slug, wn.Url)

	if len(reductions[data.PageTitleProperty][su]) == 0 {
		reductions[data.PageTitleProperty][su] = []string{wn.Label}
	}

	if len(wn.SubNav) == 0 {
		return
	}

	if wn.Url == "" {
		wn.Url = view_models.MainPage
	}

	subnav := make([]string, 0, len(wn.SubNav))
	for _, sn := range wn.SubNav {
		subnav = append(subnav, sn.Url)
		setNavigationSubNav(slug, &sn, reductions)
	}

	reductions[data.SubNavProperty][su] = subnav

}
