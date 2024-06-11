package cli

import (
	"fmt"
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

	manuals, err := scanManuals(slug)
	if err != nil {
		return sa.EndWithError(err)
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

	if err := rdx.ReplaceValues(data.ManualsProperty, slug, manuals...); err != nil {
		return sa.EndWithError(err)
	}

	result := "done"
	if len(manuals) > 0 {
		result = fmt.Sprintf("found %d manual(s)", len(manuals))
	}

	sa.EndWithResult(result)

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
