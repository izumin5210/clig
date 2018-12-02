package cmd

import (
	"github.com/spf13/cobra"

	"github.com/izumin5210/clig/pkg/cli"
	"github.com/izumin5210/clig/pkg/clig"
)

func NewDefaultCligCommand(build cli.Build) *cobra.Command {
	return NewCligCommand(&clig.Ctx{
		IO:    cli.Stdio(),
		Build: build,
	})
}

func NewCligCommand(ctx *clig.Ctx) *cobra.Command {
	cmd := &cobra.Command{
		Use: ctx.Build.AppName,
	}

	cli.AddLoggingFlags(cmd)

	cmd.AddCommand(
		cli.NewVersionCommand(ctx.IO, ctx.Build),
	)

	return cmd
}
