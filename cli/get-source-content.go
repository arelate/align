package cli

import (
	"fmt"
	"github.com/boggydigital/nod"
	"net/url"
	"strings"
)

func GetSourceContentHandler(u *url.URL) error {

	slugs := strings.Split(u.Query().Get("slug"), ",")

	return GetSourceContent(slugs...)
}

func GetSourceContent(slugs ...string) error {

	gsca := nod.NewProgress("getting source content...")
	defer gsca.End()

	gsca.TotalInt(len(slugs))

	for _, slug := range slugs {
		fmt.Println(slug)

		gsca.Increment()
	}

	gsca.EndWithResult("done")

	return nil
}
