package cli

import "net/url"

func ParseLinkedPagesHandler(u *url.URL) error {
	slug := u.Query().Get("slug")
	page := u.Query().Get("page")
	return ParseLinkedPages(slug, page)
}

func ParseLinkedPages(slug, page string) error {
	return nil
}
