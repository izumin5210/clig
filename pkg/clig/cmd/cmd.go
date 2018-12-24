package cmd

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"k8s.io/utils/exec"

	"github.com/izumin5210/clig/pkg/clib"
	"github.com/izumin5210/clig/pkg/clig"
)

func NewDefaultCligCommand(wd clib.Path, build clib.Build) *cobra.Command {
	return NewCligCommand(&clig.Ctx{
		WorkingDir: wd,
		IO:         clib.Stdio(),
		FS:         afero.NewOsFs(),
		Exec:       exec.New(),
		Build:      build,
	})
}

func NewCligCommand(ctx *clig.Ctx) *cobra.Command {
	cmd := &cobra.Command{
		Use: ctx.Build.AppName,
	}

	clib.AddLoggingFlags(cmd)

	cmd.AddCommand(
		newInitCommand(ctx),
		clib.NewVersionCommand(ctx.IO, ctx.Build),
	)

	return cmd
}
