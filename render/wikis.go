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

	titles := make(map[string]string)

	slugs := rdx.Keys(data.NavigationTitleProperty)
	sort.Strings(slugs)

	for _, slug := range slugs {
		if title, ok := rdx.GetFirstVal(data.NavigationTitleProperty, slug); ok {
			titles[slug] = title
		}
	}

	wvm, err := view_models.NewWikisViewModel(slugs, titles)
	if err != nil {
		return err
	}

	if err := tmpl.ExecuteTemplate(w, "wikis-slug", wvm); err != nil {
		return err
	}

	return nil
}
