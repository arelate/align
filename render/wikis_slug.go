package render

import (
	"encoding/json"
	"fmt"
	"github.com/arelate/align/paths"
	"github.com/arelate/align/render/view_models"
	"github.com/arelate/southern_light/ign_integration"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/pathways"
	"io"
)

func WikisPage(slug string, w io.Writer) error {
	snd, err := pathways.GetAbsDir(paths.Navigation)
	if err != nil {
		return err
	}

	nkv, err := kvas.ConnectLocal(snd, kvas.JsonExt)
	if err != nil {
		return err
	}

	toc, err := nkv.Get(slug)
	if err != nil {
		return err
	}

	if toc == nil {
		return fmt.Errorf("no toc for %s", slug)
	}

	defer toc.Close()

	var wikiNavigation []ign_integration.WikiNavigation

	if err := json.NewDecoder(toc).Decode(&wikiNavigation); err != nil {
		return err
	}

	wsvm := view_models.NewWikiSlugViewModel(slug, wikiNavigation)

	if err := tmpl.ExecuteTemplate(w, "wikis-slug", wsvm); err != nil {
		return err
	}

	return nil
}
