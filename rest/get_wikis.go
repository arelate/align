package rest

import (
	"encoding/json"
	"fmt"
	"github.com/arelate/align/paths"
	"github.com/arelate/southern_light/ign_integration"
	"github.com/boggydigital/kvas"
	"html/template"
	"net/http"
	"strings"
	"time"
)

type WikiPageViewModel struct {
	Title         string
	PublishDate   string
	UpdatedAt     string
	Entities      []template.HTML
	PrevPageLabel string
	PrevPageUrl   string
	NextPageLabel string
	NextPageUrl   string
}

func GetWikis(w http.ResponseWriter, r *http.Request) {

	// GET /wikis/{slug}/{page}

	slug := r.PathValue("slug")
	page := r.PathValue("page")

	if page == "" {
		mainPageUrl := fmt.Sprintf("/wikis/%s/Main_Page", slug)
		http.Redirect(w, r, mainPageUrl, http.StatusPermanentRedirect)
		return
	}

	if _, ok := keyValues[slug]; !ok {
		sdd, err := paths.AbsDataSlugDir(slug)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		keyValues[slug], err = kvas.ConnectLocal(sdd, kvas.JsonExt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	kv := keyValues[slug]

	wp, err := kv.Get(page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer wp.Close()

	var wikiProps ign_integration.WikiProps
	if err := json.NewDecoder(wp).Decode(&wikiProps); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	wpvm := NewWikiPageViewModel(&wikiProps)

	if err := tmpl.ExecuteTemplate(w, "wiki-page", wpvm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func rewriteImageLinks(html string) string {
	return strings.Replace(html, "https://oyster.ignimgs.com/mediawiki/apis.ign.com", "/image", -1)
}

func NewWikiPageViewModel(wp *ign_integration.WikiProps) *WikiPageViewModel {
	page := wp.Props.PageProps.Page

	wpvm := &WikiPageViewModel{
		Title:       page.Title,
		PublishDate: page.PublishDate.Format(time.RFC1123),
		UpdatedAt:   page.UpdatedAt.Format(time.RFC1123),
	}

	for _, he := range wp.HTMLEntities() {
		wpvm.Entities = append(wpvm.Entities, template.HTML(rewriteImageLinks(he.Values.Html)))

		for _, iv := range he.ImageValues {
			wpvm.Entities = append(wpvm.Entities, template.HTML("<img src='"+rewriteImageLinks(iv.Original)+"' />"))
		}
	}

	wpvm.NextPageLabel = wp.Props.PageProps.Page.Page.NextPage.Label
	wpvm.NextPageUrl = wp.NextPageUrl()

	wpvm.PrevPageLabel = wp.Props.PageProps.Page.Page.PrevPage.Label
	wpvm.PrevPageUrl = wp.PreviousPageUrl()

	return wpvm
}
