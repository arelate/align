package nav

import (
	"github.com/arelate/align/render/view_models"
	"github.com/arelate/southern_light/ign_integration"
)

func AllLinks(wikiNavigation []ign_integration.WikiNavigation) []string {
	if len(wikiNavigation) == 0 {
		return nil
	}
	links := make([]string, 0, len(wikiNavigation))
	for _, wiki := range wikiNavigation {
		link := wiki.Url
		if link == "" {
			link = view_models.MainPage
		}
		links = append(links, link)
		links = append(links, AllLinks(wiki.SubNav)...)
	}
	return links
}
