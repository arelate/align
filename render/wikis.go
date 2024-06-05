package render

import (
	"github.com/arelate/align/paths"
	"github.com/arelate/align/render/view_models"
	"io"
)

func Wikis(wikis []string, w io.Writer) error {

	for _, slug := range wikis {
		var err error
		if _, ok := keyValues[slug]; !ok {
			keyValues[slug], err = paths.DataKeyValues(slug)
			if err != nil {
				return err
			}
		}
	}

	wvm, err := view_models.NewWikisViewModel(wikis, keyValues)
	if err != nil {
		return err
	}

	if err := tmpl.ExecuteTemplate(w, "wikis-slug", wvm); err != nil {
		return err
	}

	return nil
}
