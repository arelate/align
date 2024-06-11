package cli

import (
	"github.com/arelate/align/render/view_models"
	"net/url"
)

func SyncHandler(u *url.URL) error {
	slug := u.Query().Get("slug")
	title := u.Query().Get("title")
	force := u.Query().Has("force")

	return Sync(slug, title, force)
}

func Sync(slug, title string, force bool) error {

	if err := GetPage(slug, view_models.MainPage, force); err != nil {
		return err
	}

	if err := GetData(slug, view_models.MainPage, force); err != nil {
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

	if err := Reduce(false, slug); err != nil {
		return err
	}

	if err := Scan(slug, title); err != nil {
		return err
	}

	if err := Backup(); err != nil {
		return err
	}

	return nil
}
