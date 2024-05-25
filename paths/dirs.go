package paths

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/pathways"
	"path/filepath"
)

const DefaultAlignRootDir = "/usr/share/align"

const (
	Backups    pathways.AbsDir = "backups"
	Data       pathways.AbsDir = "data"
	Images     pathways.AbsDir = "images"
	Metadata   pathways.AbsDir = "metadata"
	Navigation pathways.AbsDir = "navigation"
	Pages      pathways.AbsDir = "pages"
	Static     pathways.AbsDir = "statics"
)

var AllAbsDirs = []pathways.AbsDir{
	Backups,
	Data,
	Images,
	Metadata,
	Navigation,
	Pages,
	Static,
}

func slugKeyValues(slug string, absDir pathways.AbsDir, ext string) (kvas.KeyValues, error) {
	if spd, err := pathways.GetAbsDir(absDir); err == nil {
		aspd := filepath.Join(spd, slug)
		return kvas.NewKeyValues(aspd, ext)
	} else {
		return nil, err
	}
}

func PagesKeyValues(slug string) (kvas.KeyValues, error) {
	return slugKeyValues(slug, Pages, kvas.HtmlExt)
}

func DataKeyValues(slug string) (kvas.KeyValues, error) {
	return slugKeyValues(slug, Data, kvas.JsonExt)
}

func StaticsKeyValues(slug string) (kvas.KeyValues, error) {
	return slugKeyValues(slug, Static, kvas.HtmlExt)
}

func AbsImagesSlugDir(slug string) (string, error) {
	if isd, err := pathways.GetAbsDir(Images); err == nil {
		return filepath.Join(isd, slug), nil
	} else {
		return "", err
	}
}
