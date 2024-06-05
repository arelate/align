package rest

import (
	"github.com/arelate/align/paths"
	"github.com/arelate/align/render"
	"net/http"
	"sort"
)

func GetWikis(w http.ResponseWriter, r *http.Request) {

	// GET /wikis/{slug}

	nkv, err := paths.NavigationKeyValues()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	wikis := nkv.Keys()

	sort.Strings(wikis)

	if err := render.Wikis(wikis, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
