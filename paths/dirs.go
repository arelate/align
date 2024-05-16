package paths

import (
	"github.com/boggydigital/pathways"
)

const DefaultBoilerplateRootDir = "/usr/share/boilerplate"

const (
	Backups  pathways.AbsDir = "backups"
	Metadata pathways.AbsDir = "metadata"
)

var AllAbsDirs = []pathways.AbsDir{
	Backups,
	Metadata,
}
