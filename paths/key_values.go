package paths

import (
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/pathways"
	"path/filepath"
)

func slugKeyValues(slug string, absDir pathways.AbsDir, ext string) (kevlar.KeyValues, error) {
	if spd, err := pathways.GetAbsDir(absDir); err == nil {
		aspd := filepath.Join(spd, slug)
		return kevlar.NewKeyValues(aspd, ext)
	} else {
		return nil, err
	}
}

func PagesKeyValues(slug string) (kevlar.KeyValues, error) {
	return slugKeyValues(slug, Pages, kevlar.HtmlExt)
}

func DataKeyValues(slug string) (kevlar.KeyValues, error) {
	return slugKeyValues(slug, Data, kevlar.JsonExt)
}

func NavigationKeyValues() (kevlar.KeyValues, error) {
	nd, err := pathways.GetAbsDir(Navigation)
	if err != nil {
		return nil, err
	}
	return kevlar.NewKeyValues(nd, kevlar.JsonExt)
}
