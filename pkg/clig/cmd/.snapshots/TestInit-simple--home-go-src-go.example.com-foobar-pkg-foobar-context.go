package foobar

import (
	"github.com/izumin5210/clig/pkg/cli"
)

type Ctx struct {
	WorkingDir cli.Path
	IO         cli.IO

	Build cli.Build
}

