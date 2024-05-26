package render

import (
	"encoding/json"
	"errors"
	"github.com/arelate/align/paths"
	"github.com/arelate/align/render/view_models"
	"github.com/arelate/southern_light/ign_integration"
	"io"
)

func WikisSlugPage(slug, page string, w io.Writer) error {

	var err error
	if _, ok := keyValues[slug]; !ok {
		keyValues[slug], err = paths.DataKeyValues(slug)
		if err != nil {
			return err
		}
	}

	kv := keyValues[slug]

	wp, err := kv.Get(page)
	if err != nil {
		return err
	}
	if wp == nil {
		return errors.New("page not found: " + page)
	}
	defer wp.Close()

	var wikiProps ign_integration.WikiProps
	if err := json.NewDecoder(wp).Decode(&wikiProps); err != nil {
		return err
	}

	wpvm := view_models.NewWikiPageViewModel(slug, &wikiProps)

	if err := tmpl.ExecuteTemplate(w, "wikis-slug-page", wpvm); err != nil {
		return err
	}

	return nil
}
