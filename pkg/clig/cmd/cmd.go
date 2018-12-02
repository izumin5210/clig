package cmd

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/izumin5210/clig/pkg/cli"
	"github.com/izumin5210/clig/pkg/clig"
)

func NewDefaultCligCommand(wd cli.Path, build cli.Build) *cobra.Command {
	return NewCligCommand(&clig.Ctx{
		WorkingDir: wd,
		IO:         cli.Stdio(),
		FS:         afero.NewOsFs(),
		Build:      build,
	})
}

func NewCligCommand(ctx *clig.Ctx) *cobra.Command {
	cmd := &cobra.Command{
		Use: ctx.Build.AppName,
	}

	cli.AddLoggingFlags(cmd)

	cmd.AddCommand(
		newInitCommand(ctx),
		cli.NewVersionCommand(ctx.IO, ctx.Build),
	)

	return cmd
}
