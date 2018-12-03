package foobar

import (
	"github.com/izumin5210/clig/pkg/cli"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"go.uber.org/zap"
	"k8s.io/utils/exec"
)

type Ctx struct {
	WorkingDir cli.Path
	IO         cli.IO
	FS         afero.Fs
	Exec       exec.Interface

	Build  cli.Build
	Config *Config
}

func (c *Ctx) Init() error {

	return nil
}

