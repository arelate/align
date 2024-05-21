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

		"GET /wikis/{slug}":        Log(http.HandlerFunc(GetWikis)),
		"GET /wikis/{slug}/{page}": Log(http.HandlerFunc(GetWikis)),

		"GET /image/{slug}/{a}/{bc}/{filename}": Log(http.HandlerFunc(GetImage)),

		"GET /": Log(http.NotFoundHandler()),
	}

	for p, h := range patternHandlers {
		http.HandleFunc(p, h.ServeHTTP)
	}
}
