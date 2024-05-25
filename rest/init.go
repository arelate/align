package rest

import "github.com/boggydigital/kvas"

var staticsKeyValues map[string]kvas.KeyValues

func Init() error {

	staticsKeyValues = make(map[string]kvas.KeyValues)

	return nil
}
