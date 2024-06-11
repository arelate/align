package rest

import (
	"github.com/arelate/align/paths"
	"github.com/boggydigital/pathways"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

func GetPrimaryImage(w http.ResponseWriter, r *http.Request) {

	// GET /primary_image/{yyyy}/{mm}/{dd}/{image}

	yyyy := r.PathValue("yyyy")
	mm := r.PathValue("mm")
	dd := r.PathValue("dd")
	image := r.PathValue("image")

	// make sure we're working with the filename and not a path
	image = path.Base(image)

	spid, err := pathways.GetAbsRelDir(paths.PrimaryImages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	relPath := path.Join(yyyy, mm, dd, image)

	absImageFilename := filepath.Join(spid, relPath)
	if _, err := os.Stat(absImageFilename); err == nil {
		http.ServeFile(w, r, absImageFilename)
	} else if os.IsNotExist(err) {
		http.Error(w, absImageFilename, http.StatusNotFound)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
