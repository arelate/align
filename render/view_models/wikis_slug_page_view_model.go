package view_models

import (
	"github.com/arelate/southern_light/ign_integration"
	"html/template"
	"strings"
)

type WikisSlugPageViewModel struct {
	Slug          string
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

func rewriteImageLinks(html string) string {
	return strings.Replace(html, "https://oyster.ignimgs.com/mediawiki/apis.ign.com", "/image", -1)
}

func rewriteOriginLinks(html string) string {
	return strings.Replace(html, "https://www.ign.com", "", -1)
}

func disableStyles(html string) string {
	return strings.Replace(html, "style=", "data-style=", -1)
}

func NewWikiPageViewModel(slug string, wp *ign_integration.WikiProps) *WikisSlugPageViewModel {
	page := wp.Props.PageProps.Page

	wpvm := &WikisSlugPageViewModel{
		Slug:         slug,
		WikiPageName: page.Name,
		PageTitle:    page.Page.Title,
		PublishDate:  page.PublishDate.Format("Jan 2, 2006"),
		UpdatedAt:    page.UpdatedAt.Format("Jan 2, 2006"),
	}

	for _, he := range wp.HTMLEntities() {
		content := he.Values.Html

		content = rewriteOriginLinks(content)
		content = rewriteImageLinks(content)
		content = disableStyles(content)

		wpvm.Entities = append(wpvm.Entities, template.HTML(content))

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
