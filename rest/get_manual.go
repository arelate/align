package rest

import (
	"github.com/arelate/align/data"
	"github.com/boggydigital/pathways"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

func GetManual(w http.ResponseWriter, r *http.Request) {

	// GET /manual/{slug}/{filename}

	slug := r.PathValue("slug")
	filename := r.PathValue("filename")

	// make sure we're working with the filename and not a path
	filename = path.Base(filename)

	md, err := pathways.GetAbsDir(data.ManualsProperty)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	absManualFilename := filepath.Join(md, slug, filename)
	if _, err := os.Stat(absManualFilename); err == nil {
		http.ServeFile(w, r, absManualFilename)
	} else if os.IsNotExist(err) {
		http.Error(w, absManualFilename, http.StatusNotFound)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
