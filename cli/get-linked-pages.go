package cli

import "net/url"

func GetLinkedPagesHandler(u *url.URL) error {
	slug := u.Query().Get("slug")
	page := u.Query().Get("page")

	return GetLinkedPages(slug, page)
}

func GetLinkedPages(slug, page string) error {

	// get next page
	// get previous page
	// process HTML content and extract pages

	return nil
}
