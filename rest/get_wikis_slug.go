package rest

import (
	"fmt"
	"github.com/arelate/align/paths"
	"github.com/arelate/align/render"
	"github.com/arelate/align/render/view_models"
	"io"
	"net/http"
)

func GetWikisSlug(w http.ResponseWriter, r *http.Request) {

	// GET /wikis/{slug}

	slug := r.PathValue("slug")

	var err error

	nkv, err := paths.NavigationKeyValues()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !nkv.Has(slug) {
		mpUrl := fmt.Sprintf("/%s/%s", slug, view_models.MainPage)
		http.Redirect(w, r, mpUrl, http.StatusTemporaryRedirect)
		return
	}

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
