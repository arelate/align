package render

import (
	"github.com/arelate/align/render/view_models"
	"io"
)

func Wikis(wikis []string, w io.Writer) error {

	wvm := view_models.NewWikisViewModel(wikis)

	if err := tmpl.ExecuteTemplate(w, "wikis-slug", wvm); err != nil {
		return err
	}

	return nil
}
