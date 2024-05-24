package render

import (
	"embed"
	"github.com/boggydigital/kvas"
	"html/template"
)

var (
	tmpl *template.Template
	//go:embed "templates/*.gohtml"
	templates embed.FS

	keyValues map[string]kvas.KeyValues
)

func Init() error {

	keyValues = make(map[string]kvas.KeyValues)

	tmpl = template.Must(
		template.
			New("").
			ParseFS(templates, "templates/*.gohtml"))

	return nil
}
