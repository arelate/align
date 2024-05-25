package render

import (
	"encoding/json"
	"fmt"
	"github.com/arelate/align/paths"
	"github.com/arelate/align/render/view_models"
	"github.com/arelate/southern_light/ign_integration"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/pathways"
)

func WikiNavigation(slug string) ([]ign_integration.WikiNavigation, error) {

	snd, err := pathways.GetAbsDir(paths.Navigation)
	if err != nil {
		return nil, err
	}

	nkv, err := kvas.NewKeyValues(snd, kvas.JsonExt)
	if err != nil {
		return nil, err
	}

	toc, err := nkv.Get(slug)
	if err != nil {
		return nil, err
	}

	if toc == nil {
		return nil, fmt.Errorf("no toc for %s", slug)
	}

	defer toc.Close()

	var wikiNavigation []ign_integration.WikiNavigation

	if err := json.NewDecoder(toc).Decode(&wikiNavigation); err != nil {
		return nil, err
	}

	return wikiNavigation, nil
}

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
