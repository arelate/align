package cli

import (
	"github.com/arelate/align/paths"
	"net/url"
)

func RegenAllStaticsHandler(u *url.URL) error {
	return RegenAllStatics()
}

func RegenAllStatics() error {

	nkv, err := paths.NavigationKeyValues()
	if err != nil {
		return err
	}

	for _, slug := range nkv.Keys() {
		if err := GenAllStatics(slug); err != nil {
			return err
		}
	}

	return nil
}
