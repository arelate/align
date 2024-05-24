package cli

import "net/url"

func GenAllStaticsHandler(u *url.URL) error {
	slug := u.Query().Get("slug")
	force := u.Query().Has("force")

	return GenAllStatics(slug, force)
}

func GenAllStatics(slug string, force bool) error {
	return nil
}
