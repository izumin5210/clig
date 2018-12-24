package clig

import (
	"github.com/spf13/afero"
	"k8s.io/utils/exec"

	"github.com/izumin5210/clig/pkg/clib"
)

type Ctx struct {
	WorkingDir clib.Path
	IO         clib.IO
	FS         afero.Fs
	Exec       exec.Interface

	Build clib.Build
}
