package foobar

import (
	"github.com/izumin5210/clig/pkg/clib"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"go.uber.org/zap"
	"k8s.io/utils/exec"
)

type Ctx struct {
	WorkingDir clib.Path
	IO         *clib.IO
	FS         afero.Fs
	Exec       exec.Interface

	Build  clib.Build
	Config *Config
}

func (c *Ctx) Init() error {

	return nil
}

