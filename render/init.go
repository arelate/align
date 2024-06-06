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

	rdx kvas.ReadableRedux
)

func Init() error {

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
