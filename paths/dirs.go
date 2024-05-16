package paths

import (
	"github.com/boggydigital/pathways"
)

const DefaultAlignRootDir = "/usr/share/align"

const (
	Backups        pathways.AbsDir = "backups"
	Images         pathways.AbsDir = "images"
	Metadata       pathways.AbsDir = "metadata"
	ReducedContent pathways.AbsDir = "reduced-content"
	SourceContent  pathways.AbsDir = "source-content"
)

var AllAbsDirs = []pathways.AbsDir{
	Backups,
	Images,
	Metadata,
	ReducedContent,
	SourceContent,
}
