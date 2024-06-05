package render

import (
	"github.com/arelate/align/render/view_models"
	"io"
)

func WikisSlugPage(slug, page string, w io.Writer) error {

	var err error
	rdx, err = rdx.RefreshReader()
	if err != nil {
		return err
	}

	wpvm := view_models.NewWikiPageViewModel(slug, page, rdx)

	if err := tmpl.ExecuteTemplate(w, "wikis-slug-page", wpvm); err != nil {
		return err
	}

	return nil
}
