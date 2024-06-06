package nav

import (
	"encoding/json"
	"fmt"
	"github.com/arelate/align/paths"
	"github.com/arelate/southern_light/ign_integration"
)

func WikiNavigation(slug string) ([]ign_integration.WikiNavigation, error) {

	nkv, err := paths.NavigationKeyValues()
	if err != nil {
		return nil, err
	}

	toc, err := nkv.Get(slug)
	if err != nil {
		return nil, err
	}

	if toc == nil {
		return nil, fmt.Errorf("no toc for %s", slug)
	}

	defer toc.Close()

	var wikiNavigation []ign_integration.WikiNavigation

	if err := json.NewDecoder(toc).Decode(&wikiNavigation); err != nil {
		return nil, err
	}

	return wikiNavigation, nil
}
