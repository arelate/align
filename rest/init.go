package rest

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

	//metadataDir, err := pasu.GetAbsDir(paths.Metadata)
	//if err != nil {
	//	return err
	//}

	//if rdx, err = kvas.NewReduxReader(metadataDir, data.AllProperties()...); err != nil {
	//	return err
	//}

	//if progressRdx, err = kvas.NewReduxWriter(metadataDir,
	//	data.VideoProgressProperty,
	//	data.VideoEndedProperty,
	//	data.PlaylistNewVideosProperty,
	//); err != nil {
	//	return err
	//}

	keyValues = make(map[string]kvas.KeyValues)

	tmpl = template.Must(
		template.
			New("").
			ParseFS(templates, "templates/*.gohtml"))

	return nil
}
