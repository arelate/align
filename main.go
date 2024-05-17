package main

import (
	"bytes"
	_ "embed"
	"github.com/arelate/align/cli"
	"github.com/arelate/align/paths"
	"github.com/boggydigital/clo"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathways"
	"os"
)

var (
	//go:embed "cli-commands.txt"
	cliCommands []byte
	//go:embed "cli-help.txt"
	cliHelp []byte
)

const (
	dirOverridesFilename = "directories.txt"
)

func main() {
	nod.EnableStdOutPresenter()

	ea := nod.Begin("align is serving gaming guides...")
	defer ea.End()

	if err := pathways.Setup(dirOverridesFilename,
		paths.DefaultAlignRootDir,
		nil,
		paths.AllAbsDirs...); err != nil {
		_ = ea.EndWithError(err)
		os.Exit(1)
	}

	defs, err := clo.Load(
		bytes.NewBuffer(cliCommands),
		bytes.NewBuffer(cliHelp),
		nil)

	if err != nil {
		_ = ea.EndWithError(err)
		os.Exit(1)
	}

	clo.HandleFuncs(map[string]clo.Handler{
		"backup":          cli.BackupHandler,
		"get-source-page": cli.GetSourcePageHandler,
		"reduce-page":     cli.ReducePageHandler,
		"serve":           cli.ServeHandler,
		"version":         cli.VersionHandler,
	})

	if err := defs.AssertCommandsHaveHandlers(); err != nil {
		_ = ea.EndWithError(err)
		os.Exit(1)
	}

	if err := defs.Serve(os.Args[1:]); err != nil {
		_ = ea.EndWithError(err)
		os.Exit(1)
	}

}
