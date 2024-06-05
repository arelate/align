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

	//if _, ok := staticsKeyValues[slug]; !ok {
	//	staticsKeyValues[slug], err = paths.StaticsKeyValues(slug)
	//	if err != nil {
	//		http.Error(w, err.Error(), http.StatusInternalServerError)
	//		return
	//	}
	//}
	//
	//skv := staticsKeyValues[slug]
	//
	//if skv.Has(page) {
	//
	//	rc, err := skv.Get(page)
	//	if err != nil {
	//		http.Error(w, err.Error(), http.StatusInternalServerError)
	//		return
	//	}
	//
	//	if _, err := io.Copy(w, rc); err != nil {
	//		rc.Close()
	//		http.Error(w, err.Error(), http.StatusInternalServerError)
	//		return
	//	}
	//
	//	rc.Close()
	//	return
	//}

	if err := render.WikisSlugPage(slug, page, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
