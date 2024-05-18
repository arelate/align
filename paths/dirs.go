package paths

import (
	"github.com/boggydigital/pathways"
	"path/filepath"
)

const DefaultAlignRootDir = "/usr/share/align"

const (
	Backups      pathways.AbsDir = "backups"
	Images       pathways.AbsDir = "images"
	Metadata     pathways.AbsDir = "metadata"
	SourcePages  pathways.AbsDir = "source-pages"
	ReducedPages pathways.AbsDir = "reduced-pages"
)

var AllAbsDirs = []pathways.AbsDir{
	Backups,
	Images,
	Metadata,
	SourcePages,
	ReducedPages,
}

func AbsSourcePagesDir(slug string) (string, error) {
	if spd, err := pathways.GetAbsDir(SourcePages); err == nil {
		return filepath.Join(spd, slug), nil
	} else {
		return "", err
	}
}

func AbsReducedPagesDir(slug string) (string, error) {
	if spd, err := pathways.GetAbsDir(ReducedPages); err == nil {
		return filepath.Join(spd, slug), nil
	} else {
		return "", err
	}
}
