package cli

import (
	"fmt"
	"github.com/arelate/align/rest"
	"github.com/boggydigital/nod"
	"net/http"
	"net/url"
	"strconv"
)

func ServeHandler(u *url.URL) error {
	portStr := u.Query().Get("port")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return err
	}

	stderr := u.Query().Has("stderr")

	return Serve(port, stderr)
}

func Serve(port int, stderr bool) error {

	if stderr {
		nod.EnableStdErrLogger()
		nod.DisableOutput(nod.StdOut)
	}

	rest.HandleFuncs()

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
