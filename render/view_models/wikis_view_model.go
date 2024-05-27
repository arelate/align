package view_models

import (
	"fmt"
	"html/template"
)

func NewWikisViewModel(wikis []string) *WikisSlugViewModel {
	wsvm := &WikisSlugViewModel{
		GuideTitle: "All Guides",
	}

	for _, w := range wikis {
		u := fmt.Sprintf("<a href='/wikis/%s'>%s</a>", w, w)
		wsvm.Items = append(wsvm.Items,
			template.HTML(u))
	}

	return wsvm
}
