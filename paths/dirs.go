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

func AbsPagesDir(slug string) (string, error) {
	if spd, err := pathways.GetAbsDir(Pages); err == nil {
		return filepath.Join(spd, slug), nil
	} else {
		return "", err
	}
}

func AbsDataDir(slug string) (string, error) {
	if spd, err := pathways.GetAbsDir(Data); err == nil {
		return filepath.Join(spd, slug), nil
	} else {
		return "", err
	}
}
