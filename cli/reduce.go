package cli

import "net/url"

func ReduceHandler(u *url.URL) error {
	slug := u.Query().Get("slug")
	force := u.Query().Has("force")

	return Reduce(slug, force)
}

func Reduce(slug string, force bool) error {
	return nil
}
