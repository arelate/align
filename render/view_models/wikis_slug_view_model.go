package view_models

import (
	"fmt"
	"github.com/arelate/southern_light/ign_integration"
	"html/template"
	"net/url"
)

const MainPage = "Main_Page"

type WikisSlugViewModel struct {
	GuideTitle string
	Slug       string
	Items      []template.HTML
}

func NewWikiSlugViewModel(slug string, wikiNavigation []ign_integration.WikiNavigation) *WikisSlugViewModel {
	wsvm := &WikisSlugViewModel{
		GuideTitle: GuideTitle(wikiNavigation),
		Slug:       slug,
		Items:      make([]template.HTML, 0),
	}

	for _, wn := range wikiNavigation {
		wsvm.Items = append(wsvm.Items, template.HTML(WikiNavigationHTML(slug, wn)))
	}

	return wsvm
}

func GuideTitle(wn []ign_integration.WikiNavigation) string {
	for _, w := range wn {
		return w.Label
	}
	return ""
}

func WikiNavigationHTML(slug string, wn ign_integration.WikiNavigation) string {

	u := url.PathEscape(wn.Url)
	if u == "" {
		u = MainPage
	}

	link := fmt.Sprintf("<a href='/wikis/%s/%s'>%s</a>", slug, u, wn.Label)
	if len(wn.SubNav) > 0 {
		link += "<ul>"
	}
	for _, sn := range wn.SubNav {
		link += "<li>" + WikiNavigationHTML(slug, sn) + "</li>"
	}
	if len(wn.SubNav) > 0 {
		link += "</ul>"
	}
	return link
}
