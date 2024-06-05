package paths

import (
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/pathways"
	"path/filepath"
)

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

func NavigationKeyValues() (kvas.KeyValues, error) {
	nd, err := pathways.GetAbsDir(Navigation)
	if err != nil {
		return nil, err
	}
	return kvas.NewKeyValues(nd, kvas.JsonExt)
}

func StaticsKeyValues(slug string) (kvas.KeyValues, error) {
	return slugKeyValues(slug, Static, kvas.HtmlExt)
}
