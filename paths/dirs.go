package paths

import (
	"github.com/boggydigital/pathways"
	"path/filepath"
)

const DefaultAlignRootDir = "/usr/share/align"

const (
	Backups  pathways.AbsDir = "backups"
	Images   pathways.AbsDir = "images"
	Metadata pathways.AbsDir = "metadata"
	Pages    pathways.AbsDir = "pages"
	Data     pathways.AbsDir = "data"
)

var AllAbsDirs = []pathways.AbsDir{
	Backups,
	Images,
	Metadata,
	Pages,
	Data,
}

func AbsPagesSlugDir(slug string) (string, error) {
	return absSlugDir(slug, Pages)
}

func AbsDataSlugDir(slug string) (string, error) {
	return absSlugDir(slug, Data)
}

func AbsImagesSlugDir(slug string) (string, error) {
	return absSlugDir(slug, Images)
}

func absSlugDir(slug string, absDir pathways.AbsDir) (string, error) {
	if spd, err := pathways.GetAbsDir(absDir); err == nil {
		return filepath.Join(spd, slug), nil
	} else {
		return "", err
	}
}
