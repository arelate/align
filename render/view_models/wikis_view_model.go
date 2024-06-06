package view_models

import (
	"fmt"
	"html/template"
	"net/url"
	"path"
)

func NewWikisViewModel(order []string, wikiPrimaryImages map[string]string) (*WikisSlugViewModel, error) {
	wsvm := &WikisSlugViewModel{
		Title:    "All Guides",
		Wrapping: true,
	}

	for _, slug := range order {

		piu, err := url.Parse(wikiPrimaryImages[slug])
		if err != nil {
			return nil, err
		}

		u := fmt.Sprintf("<a href='/wikis/%s'><img src='/%s' title='%s' /></a>",
			slug,
			path.Join("primary_image", piu.Path),
			slug)

		wsvm.Items = append(wsvm.Items,
			template.HTML(u))
	}

	return wsvm, nil
}
