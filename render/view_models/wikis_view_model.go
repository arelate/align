package view_models

import (
	"fmt"
	"html/template"
)

func NewWikisViewModel(order []string, titles map[string]string) (*WikisSlugViewModel, error) {
	wsvm := &WikisSlugViewModel{
		Title: "All Guides",
	}

	for _, slug := range order {
		u := fmt.Sprintf("<a href='/wikis/%s'>%s</a>", slug, titles[slug])
		wsvm.Items = append(wsvm.Items, template.HTML(u))
	}

	return wsvm, nil
}
