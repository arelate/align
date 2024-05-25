package cli

import "net/url"

func SyncHandler(u *url.URL) error {
	slug := u.Query().Get("slug")
	force := u.Query().Has("force")

	return Sync(slug, force)
}

func Sync(slug string, force bool) error {

	if err := GetPage(slug, "", force); err != nil {
		return err
	}

	if err := GetNavigation(slug, force); err != nil {
		return err
	}

	if err := GetAllPages(slug, defaultThrottleMs, force); err != nil {
		return err
	}

	if err := GetAllImages(slug, force); err != nil {
		return err
	}

	if err := GenTOC(slug, force); err != nil {
		return err
	}

	if err := GenAllStatics(slug, force); err != nil {
		return err
	}

	return nil
}
