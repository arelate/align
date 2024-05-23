package rest

import (
	"encoding/json"
	"fmt"
	"github.com/arelate/align/paths"
	"github.com/arelate/southern_light/ign_integration"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/pathways"
	"html/template"
	"net/http"
	"net/url"
)

type WikiSlugViewModel struct {
	GuideTitle string
	Items      []template.HTML
}

func GetWikisSlug(w http.ResponseWriter, r *http.Request) {

	// GET /wikis/{slug}

	slug := r.PathValue("slug")

	snd, err := pathways.GetAbsDir(paths.Navigation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	nkv, err := kvas.ConnectLocal(snd, kvas.JsonExt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	toc, err := nkv.Get(slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if toc == nil {
		http.Redirect(w, r, fmt.Sprintf("/wikis/%s/Main_Page", slug), http.StatusTemporaryRedirect)
		return
	}

	defer toc.Close()

	var wikiNavigation []ign_integration.WikiNavigation

	if err := json.NewDecoder(toc).Decode(&wikiNavigation); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	wsvm := NewWikiSlugViewModel(slug, wikiNavigation)

	if err := tmpl.ExecuteTemplate(w, "wikis-slug", wsvm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func NewWikiSlugViewModel(slug string, wikiNavigation []ign_integration.WikiNavigation) *WikiSlugViewModel {
	wsvm := &WikiSlugViewModel{
		GuideTitle: GuideTitle(wikiNavigation),
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
		u = "Main_Page"
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
