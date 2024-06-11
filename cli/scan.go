package cli

import (
	"github.com/boggydigital/nod"
	"net/url"
	"strings"
)

func ScanHandler(u *url.URL) error {
	slugs := strings.Split(u.Query().Get("slug"), ",")
	all := u.Query().Has("all")
	return Scan(all, slugs...)
}

func Scan(all bool, slugs ...string) error {
	sa := nod.NewProgress("scanning manuals...")
	defer sa.End()

	sa.EndWithResult("done")

	return nil
}
