package cli

import (
	"github.com/arelate/align/paths"
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathways"
	"net/url"
)

func MigrateHandler(_ *url.URL) error {
	return Migrate()
}

func Migrate() error {

	ma := nod.NewProgress("migrating data...")
	defer ma.End()

	targets := []pathways.AbsDir{paths.Data, paths.Metadata, paths.Navigation, paths.Pages}

	ma.TotalInt(len(targets))

	for _, target := range targets {

		dir, err := pathways.GetAbsDir(target)
		if err != nil {
			return err
		}

		if err := kevlar.MigrateAll(dir); err != nil {
			return ma.EndWithError(err)
		}

		ma.Increment()
	}

	ma.EndWithResult("done")

	return nil
}
