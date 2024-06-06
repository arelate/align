package rest

import (
	"github.com/arelate/align/render"
	"net/http"
)

func GetWikisSlug(w http.ResponseWriter, r *http.Request) {

	// GET /wikis/{slug}

	slug := r.PathValue("slug")

	if err := render.WikisPage(slug, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
