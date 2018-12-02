package clig

import (
	"github.com/spf13/afero"

	"github.com/izumin5210/clig/pkg/cli"
)

type Ctx struct {
	WorkingDir cli.Path
	IO         cli.IO
	FS         afero.Fs

	Build cli.Build
}
