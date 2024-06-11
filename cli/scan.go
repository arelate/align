package cli

import (
	"github.com/arelate/align/data"
	"github.com/arelate/align/paths"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathways"
	"net/url"
	"os"
	"path/filepath"
)

func ScanHandler(u *url.URL) error {
	slug := u.Query().Get("slug")
	title := u.Query().Get("title")
	return Scan(slug, title)
}

func Scan(slug, title string) error {
	sa := nod.Begin("scanning manuals...")
	defer sa.End()

	slugManuals := make(map[string][]string)

	if manuals, err := scanManuals(slug); err != nil {
		return sa.EndWithError(err)
	} else {
		slugManuals[slug] = manuals
	}

	rdx, err := paths.NewReduxWriter()
	if err != nil {
		return sa.EndWithError(err)
	}

	if title != "" {
		if err := rdx.AddValues(data.NavigationTitleProperty, slug, title); err != nil {
			return sa.EndWithError(err)
		}
	}

	if err := rdx.BatchReplaceValues(data.ManualsProperty, slugManuals); err != nil {
		return sa.EndWithError(err)
	}

	sa.EndWithResult("done")

	return nil
}

func scanManuals(slug string) ([]string, error) {

	md, err := pathways.GetAbsDir(paths.Manuals)
	if err != nil {
		return nil, err
	}

	smd := filepath.Join(md, slug)

	slugManualsDir, err := os.Open(smd)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	defer slugManualsDir.Close()

	return slugManualsDir.Readdirnames(-1)
}
