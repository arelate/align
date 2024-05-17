package paths

import (
	"github.com/boggydigital/pathways"
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
