package clig

import (
	"github.com/spf13/afero"
	"k8s.io/utils/exec"

	"github.com/izumin5210/clig/pkg/cli"
)

type Ctx struct {
	WorkingDir cli.Path
	IO         cli.IO
	FS         afero.Fs
	Exec       exec.Interface

	Build cli.Build
}
