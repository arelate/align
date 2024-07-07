package view_models

import (
	"github.com/arelate/align/data"
	"github.com/boggydigital/kevlar"
	"html/template"
	"path"
)

type WikisSlugPageViewModel struct {
	Slug          string
	WikiName      string
	PageTitle     string
	PublishDate   string
	UpdatedAt     string
	Entities      []template.HTML
	PrevPageTitle string
	PrevPageUrl   string
	NextPageTitle string
	NextPageUrl   string
}

func NewWikiPageViewModel(slug, page string, rdx kevlar.ReadableRedux) *WikisSlugPageViewModel {

	wpvm := &WikisSlugPageViewModel{
		Slug: slug,
	}

	if wikiName, ok := rdx.GetLastVal(data.WikiNameProperty, slug); ok {
		wpvm.WikiName = wikiName
	}

	sp := path.Join(slug, page)

	if pageTitle, ok := rdx.GetLastVal(data.PageTitleProperty, sp); ok {
		wpvm.PageTitle = pageTitle
	}
	if nextPageUrl, ok := rdx.GetLastVal(data.PageNextPageUrlProperty, sp); ok {
		wpvm.NextPageUrl = nextPageUrl
		if npt, ok := rdx.GetLastVal(data.PageTitleProperty, path.Join(slug, nextPageUrl)); ok && npt != "" {
			wpvm.NextPageTitle = npt
		} else {
			wpvm.NextPageTitle = nextPageUrl
		}
	}
	if prevPageUrl, ok := rdx.GetLastVal(data.PagePrevPageUrlProperty, sp); ok {
		wpvm.PrevPageUrl = prevPageUrl
		if ppt, ok := rdx.GetLastVal(data.PageTitleProperty, path.Join(slug, prevPageUrl)); ok && ppt != "" {
			wpvm.PrevPageTitle = ppt
		} else {
			wpvm.PrevPageTitle = prevPageUrl
		}
	}
	if publishDate, ok := rdx.GetLastVal(data.PagePublishDateProperty, sp); ok {
		wpvm.PublishDate = publishDate
	}
	if updatedAt, ok := rdx.GetLastVal(data.PageUpdatedAtProperty, sp); ok {
		wpvm.UpdatedAt = updatedAt
	}

	if htmlEntries, ok := rdx.GetAllValues(data.PageHTMLEntriesProperty, sp); ok {
		for _, entry := range htmlEntries {
			wpvm.Entities = append(wpvm.Entities, template.HTML(entry))
		}
	}

	return wpvm
}
