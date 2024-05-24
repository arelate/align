package rest

import (
	"github.com/arelate/align/render"
	"net/http"
)

func GetWikisSlugPage(w http.ResponseWriter, r *http.Request) {

	// GET /wikis/{slug}/{page}

	slug := r.PathValue("slug")
	page := r.PathValue("page")

	if err := render.WikisSlugPage(slug, page, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
