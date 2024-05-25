package render

import (
	"errors"
	"github.com/arelate/align/paths"
	"github.com/arelate/align/render/view_models"
	"io"
)

func WikisPage(slug string, w io.Writer) error {

	var err error
	if _, ok := keyValues[slug]; !ok {
		keyValues[slug], err = paths.DataKeyValues(slug)
		if err != nil {
			return err
		}
	}

	kv := keyValues[slug]

	wp, err := kv.Get(slug)
	if err != nil {
		return err
	}
	if wp == nil {
		return errors.New("toc not found: " + slug)
	}

	defer wp.Close()

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
