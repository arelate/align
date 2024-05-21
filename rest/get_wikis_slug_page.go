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
)

type WikiPageViewModel struct {
	WikiPageName  string
	PageTitle     string
	PublishDate   string
	UpdatedAt     string
	Entities      []template.HTML
	PrevPageLabel string
	PrevPageUrl   string
	NextPageLabel string
	NextPageUrl   string
}

func GetWikisSlugPage(w http.ResponseWriter, r *http.Request) {

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
		WikiPageName: page.Name,
		PageTitle:    page.Page.Title,
		PublishDate:  page.PublishDate.Format("Jan 2, 2006"),
		UpdatedAt:    page.UpdatedAt.Format("Jan 2, 2006"),
	}

	for _, he := range wp.HTMLEntities() {
		wpvm.Entities = append(wpvm.Entities, template.HTML(rewriteImageLinks(he.Values.Html)))

		imagesContent := ""
		for _, iv := range he.ImageValues {
			imagesContent += "<img src='" + rewriteImageLinks(iv.Original) + "' />"
		}

		if imagesContent != "" {
			wpvm.Entities = append(wpvm.Entities, template.HTML(imagesContent))
		}
	}

	wpvm.NextPageLabel = wp.Props.PageProps.Page.Page.NextPage.Label
	wpvm.NextPageUrl = wp.NextPageUrl()

	wpvm.PrevPageLabel = wp.Props.PageProps.Page.Page.PrevPage.Label
	wpvm.PrevPageUrl = wp.PreviousPageUrl()

	return wpvm
}

func chapterTitle(wp *ign_integration.WikiProps) (string, string) {
	title := wp.Props.PageProps.Page.Page.Title
	if parts := strings.Split(title, " - "); len(parts) == 2 {
		return "", parts[0]
	} else if len(parts) > 2 {
		return parts[0], parts[1]
	}
	return "", ""
}
