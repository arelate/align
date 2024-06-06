package rest

import (
	"github.com/arelate/align/render"
	"net/http"
)

func GetWikis(w http.ResponseWriter, r *http.Request) {

	// GET /wikis/{slug}

	if err := render.Wikis(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
