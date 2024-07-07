package view_models

import (
	"fmt"
	"github.com/arelate/align/data"
	"github.com/boggydigital/kevlar"
	"html/template"
	"net/url"
	"path"
)

const MainPage = "Main_Page"

type WikisSlugViewModel struct {
	Title   string
	Slug    string
	Items   []template.HTML
	Manuals []string
}

func NewWikiSlugViewModel(slug string, rdx kevlar.ReadableRedux) (*WikisSlugViewModel, error) {

	wsvm := &WikisSlugViewModel{
		Slug:  slug,
		Items: make([]template.HTML, 0),
	}

	if navTitle, ok := rdx.GetLastVal(data.NavigationTitleProperty, slug); ok {
		wsvm.Title = navTitle
	}

	if nav, ok := rdx.GetAllValues(data.NavigationProperty, slug); ok {
		for _, pageUrl := range nav {
			wnh, err := WikiNavigationHTML(slug, pageUrl, rdx)
			if err != nil {
				return nil, err
			}
			wsvm.Items = append(wsvm.Items, template.HTML(wnh))
		}
	}

	if manuals, ok := rdx.GetAllValues(data.ManualsProperty, slug); ok {
		wsvm.Manuals = manuals
	}

	return wsvm, nil
}

func WikiNavigationHTML(slug, pageUrl string, rdx kevlar.ReadableRedux) (string, error) {

	if pageUrl == "" {
		pageUrl = MainPage
	}

	upu, err := url.PathUnescape(pageUrl)
	if err != nil {
		return pageUrl, err
	}

	pageTitle := ""
	if pt, ok := rdx.GetLastVal(data.PageTitleProperty, path.Join(slug, pageUrl)); ok {
		pageTitle = pt
	}

	missingLink := rdx.HasValue(data.PageMissingProperty, slug, upu)

	link := ""
	if pageTitle != "" {
		attr := ""
		if !missingLink {
			attr = fmt.Sprintf("href='/wikis/%s/%s'", slug, pageUrl)
		} else {
			attr = "class='subtle'"
		}

		link = fmt.Sprintf("<a %s>%s</a>", attr, pageTitle)
	} else {
		pageUrlTitle := pageUrl
		if pu, err := url.PathUnescape(pageUrl); err == nil {
			pageUrlTitle = pu
		}
		link = fmt.Sprintf("<span class='subtle'>%s</span>", pageUrlTitle)
	}

	if subNav, ok := rdx.GetAllValues(data.SubNavProperty, path.Join(slug, pageUrl)); ok {
		if len(subNav) > 0 {
			link += "<ul>"
		}
		for _, sn := range subNav {
			subLink, err := WikiNavigationHTML(slug, sn, rdx)
			if err != nil {
				return link, err
			}
			link += "<li>" + subLink + "</li>"
		}
		if len(subNav) > 0 {
			link += "</ul>"
		}
	}

	return link, nil
}
