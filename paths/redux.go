package paths

import (
	"github.com/arelate/align/data"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/pathways"
)

func NewReduxReader() (kvas.ReadableRedux, error) {
	amd, err := pathways.GetAbsDir(Metadata)
	if err != nil {
		return nil, err
	}

	return kvas.NewReduxReader(amd, data.AllReduxProperties()...)
}

func NewReduxWriter() (kvas.WriteableRedux, error) {
	amd, err := pathways.GetAbsDir(Metadata)
	if err != nil {
		return nil, err
	}

	return kvas.NewReduxWriter(amd, data.AllReduxProperties()...)
}
