package paths

import (
	"github.com/arelate/align/data"
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/pathways"
)

func NewReduxReader() (kevlar.ReadableRedux, error) {
	amd, err := pathways.GetAbsDir(Metadata)
	if err != nil {
		return nil, err
	}

	return kevlar.NewReduxReader(amd, data.AllReduxProperties()...)
}

func NewReduxWriter() (kevlar.WriteableRedux, error) {
	amd, err := pathways.GetAbsDir(Metadata)
	if err != nil {
		return nil, err
	}

	return kevlar.NewReduxWriter(amd, data.AllReduxProperties()...)
}
