package render

import (
	"github.com/arelate/align/render/view_models"
	"io"
)

func WikisPage(slug string, w io.Writer) error {

	var err error
	rdx, err = rdx.RefreshReader()
	if err != nil {
		return err
	}

	wsvm := view_models.NewWikiSlugViewModel(slug, rdx)

	if err := tmpl.ExecuteTemplate(w, "wikis-slug", wsvm); err != nil {
		return err
	}

	return nil
}
