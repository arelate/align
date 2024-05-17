package cli

import "net/url"

func ParseLinkedPagesHandler(u *url.URL) error {
	slug := u.Query().Get("slug")
	page := u.Query().Get("page")
	return ParseLinkedPages(slug, page)
}

func ParseLinkedPages(slug, page string) error {

	// get next page
	// get previous page
	// process HTML content and extract pages

	return nil
}
