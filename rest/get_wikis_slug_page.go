package rest

import (
	"github.com/arelate/align/render"
	"net/http"
	"net/url"
)

func GetWikisSlugPage(w http.ResponseWriter, r *http.Request) {

	// GET /wikis/{slug}/{page}

	slug := r.PathValue("slug")
	page := r.PathValue("page")

	var err error

	page, err = url.PathUnescape(page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := render.WikisSlugPage(slug, page, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
