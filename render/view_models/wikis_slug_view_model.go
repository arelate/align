package view_models

import (
	"fmt"
	"github.com/arelate/align/data"
	"github.com/boggydigital/kvas"
	"html/template"
	"path"
)

const MainPage = "Main_Page"

type WikisSlugViewModel struct {
	Title    string
	Slug     string
	Wrapping bool
	Items    []template.HTML
}

func NewWikiSlugViewModel(slug string, rdx kvas.ReadableRedux) *WikisSlugViewModel {

	wsvm := &WikisSlugViewModel{
		Slug:  slug,
		Items: make([]template.HTML, 0),
	}

	if navTitle, ok := rdx.GetFirstVal(data.NavigationTitleProperty, slug); ok {
		wsvm.Title = navTitle
	}

	if nav, ok := rdx.GetAllValues(data.NavigationProperty, slug); ok {
		for _, pageUrl := range nav {
			wsvm.Items = append(wsvm.Items, template.HTML(WikiNavigationHTML(slug, pageUrl, rdx)))
		}
	}

	return wsvm
}

func WikiNavigationHTML(slug, pageUrl string, rdx kvas.ReadableRedux) string {

	if pageUrl == "" {
		pageUrl = MainPage
	}

	pageTitle := ""
	if pt, ok := rdx.GetFirstVal(data.PageTitleProperty, path.Join(slug, pageUrl)); ok {
		pageTitle = pt
	}

	link := fmt.Sprintf("<a href='/wikis/%s/%s'>%s</a>", slug, pageUrl, pageTitle)

	if subNav, ok := rdx.GetAllValues(data.SubNavProperty, path.Join(slug, pageUrl)); ok {
		if len(subNav) > 0 {
			link += "<ul>"
		}
		for _, sn := range subNav {
			link += "<li>" + WikiNavigationHTML(slug, sn, rdx) + "</li>"
		}
		if len(subNav) > 0 {
			link += "</ul>"
		}
	}

	return link
}
