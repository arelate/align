package render

import (
	"github.com/arelate/align/data"
	"github.com/arelate/align/render/view_models"
	"io"
	"sort"
)

func Wikis(w io.Writer) error {

	var err error
	rdx, err = rdx.RefreshReader()
	if err != nil {
		return err
	}

	wikiPrimaryImages := make(map[string]string)

	slugs := rdx.Keys(data.WikiPrimaryImageProperty)
	sort.Strings(slugs)

	for _, slug := range slugs {
		if primaryImage, ok := rdx.GetFirstVal(data.WikiPrimaryImageProperty, slug); ok {
			wikiPrimaryImages[slug] = primaryImage
		}
	}

	wvm, err := view_models.NewWikisViewModel(slugs, wikiPrimaryImages)
	if err != nil {
		return err
	}

	if err := tmpl.ExecuteTemplate(w, "wikis-slug", wvm); err != nil {
		return err
	}

	return nil
}
