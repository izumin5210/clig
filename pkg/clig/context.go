package clig

import (
	"github.com/spf13/afero"

	"github.com/izumin5210/clig/pkg/cli"
)

type Ctx struct {
	IO cli.IO
	FS afero.Fs

	Build cli.Build
}
