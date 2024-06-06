package data

const (
	WikiNameProperty         = "wiki-name"
	WikiPrimaryImageProperty = "wiki-primary-image"
	WikiPages                = "wiki-pages"

	// slug relative dir

	PageTitleProperty       = "page-title"
	PagePublishDateProperty = "page-publish-date"
	PageUpdatedAtProperty   = "page-updated-at"
	PageHTMLEntriesProperty = "page-html-entries"
	PageNextPageUrlProperty = "page-next-page-url"
	PagePrevPageUrlProperty = "page-prev-page-url"
	PageMissingProperty     = "page-missing"

	NavigationProperty      = "navigation"
	NavigationTitleProperty = "navigation-title"
	SubNavProperty          = "sub-nav"
)

func AllReduxProperties() []string {
	return []string{
		WikiNameProperty,
		WikiPrimaryImageProperty,
		WikiPages,
		PageTitleProperty,
		PagePublishDateProperty,
		PageUpdatedAtProperty,
		PageHTMLEntriesProperty,
		PageNextPageUrlProperty,
		PagePrevPageUrlProperty,
		PageMissingProperty,
		NavigationProperty,
		NavigationTitleProperty,
		SubNavProperty,
	}
}
