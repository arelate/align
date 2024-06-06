package paths

import (
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
)

const (
	PrimaryImages pathways.RelDir = "_primary"
)

var AllAbsDirs = []pathways.AbsDir{
	Backups,
	Data,
	Images,
	Metadata,
	Navigation,
	Pages,
}

var RelToAbsDirs = map[pathways.RelDir]pathways.AbsDir{
	PrimaryImages: Images,
}

func AbsImagesSlugDir(slug string) (string, error) {
	if isd, err := pathways.GetAbsDir(Images); err == nil {
		return filepath.Join(isd, slug), nil
	} else {
		return "", err
	}
}
