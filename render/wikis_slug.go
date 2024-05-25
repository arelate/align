package render

import (
	"github.com/arelate/align/render/view_models"
	"io"
)

func WikisPage(slug string, w io.Writer) error {

	wn, err := WikiNavigation(slug)
	if err != nil {
		return err
	}

	wsvm := view_models.NewWikiSlugViewModel(slug, wn)

	if err := tmpl.ExecuteTemplate(w, "wikis-slug", wsvm); err != nil {
		return err
	}

	return nil
}
