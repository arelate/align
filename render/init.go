package render

import (
	"embed"
	"github.com/arelate/align/paths"
	"github.com/boggydigital/kvas"
	"html/template"
)

var (
	tmpl *template.Template
	//go:embed "templates/*.gohtml"
	templates embed.FS

	keyValues map[string]kvas.KeyValues

	rdx kvas.ReadableRedux
)

func Init() error {

	keyValues = make(map[string]kvas.KeyValues)

	var err error
	rdx, err = paths.NewReduxReader()
	if err != nil {
		return err
	}

	tmpl, err =
		template.
			New("").
			ParseFS(templates, "templates/*.gohtml")
	if err != nil {
		return err
	}

	return nil
}
