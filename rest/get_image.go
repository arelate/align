package rest

import (
	"github.com/arelate/align/paths"
	"net/http"
	"os"
	"path/filepath"
)

func GetImage(w http.ResponseWriter, r *http.Request) {

	// GET /image/{slug}/{a}/{bc}/{filename}

	slug := r.PathValue("slug")
	a := r.PathValue("a")
	bc := r.PathValue("bc")
	filename := r.PathValue("filename")

	sid, err := paths.AbsImagesSlugDir(slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	absImageFilename := filepath.Join(sid, a, bc, filename)
	if _, err := os.Stat(absImageFilename); err == nil {
		http.ServeFile(w, r, absImageFilename)
	} else if os.IsNotExist(err) {
		http.Error(w, absImageFilename, http.StatusNotFound)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
