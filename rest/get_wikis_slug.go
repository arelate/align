package rest

import (
	"github.com/arelate/align/paths"
	"github.com/arelate/align/render"
	"io"
	"net/http"
)

func GetWikisSlug(w http.ResponseWriter, r *http.Request) {

	// GET /wikis/{slug}

	slug := r.PathValue("slug")

	var err error

	if _, ok := staticsKeyValues[slug]; !ok {
		staticsKeyValues[slug], err = paths.StaticsKeyValues(slug)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	skv := staticsKeyValues[slug]

	if skv.Has(slug) {

		rc, err := skv.Get(slug)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, err := io.Copy(w, rc); err != nil {
			rc.Close()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rc.Close()
		return
	}

	if err := render.WikisPage(slug, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
