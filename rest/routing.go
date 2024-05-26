package rest

import (
	"github.com/boggydigital/nod"
	"net/http"
)

var (
	Log = nod.RequestLog
)

func HandleFuncs() {

	patternHandlers := map[string]http.Handler{

		"GET /wikis":               Log(http.HandlerFunc(GetWikis)),
		"GET /wikis/{slug}":        Log(http.HandlerFunc(GetWikisSlug)),
		"GET /wikis/{slug}/":       Log(http.HandlerFunc(GetWikisSlug)),
		"GET /wikis/{slug}/{page}": Log(http.HandlerFunc(GetWikisSlugPage)),

		"GET /image/{slug}/{a}/{bc}/{filename}": Log(http.HandlerFunc(GetImage)),

		"GET /": Log(http.RedirectHandler("/wikis", http.StatusPermanentRedirect)),
	}

	for p, h := range patternHandlers {
		http.HandleFunc(p, h.ServeHTTP)
	}
}
