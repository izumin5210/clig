package cmd

import (
	"github.com/spf13/cobra"
	"github.com/izumin5210/clig/pkg/cli"

	"go.example.com/foobar/pkg/foobar"
)

func NewDefaultFoobarCommand(wd cli.Path, build cli.Build) *cobra.Command {
	return NewFoobarCommand(&foobar.Ctx{
		WorkingDir: wd,
		IO:         cli.Stdio(),
		Build:      build,
	})
}

func NewFoobarCommand(ctx *foobar.Ctx) *cobra.Command {
	cmd := &cobra.Command{
		Use: ctx.Build.AppName,
	}

	cli.AddLoggingFlags(cmd)

	cmd.AddCommand(
		cli.NewVersionCommand(ctx.IO, ctx.Build),
	)

	return cmd
}

