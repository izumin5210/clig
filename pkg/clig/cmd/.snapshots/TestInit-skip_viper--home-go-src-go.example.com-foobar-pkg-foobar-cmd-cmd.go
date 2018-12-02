package cmd

import (
	"github.com/izumin5210/clig/pkg/cli"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"k8s.io/utils/exec"

	"go.example.com/foobar/pkg/foobar"
)

func NewDefaultFoobarCommand(wd cli.Path, build cli.Build) *cobra.Command {
	return NewFoobarCommand(&foobar.Ctx{
		WorkingDir: wd,
		IO:         cli.Stdio(),
		FS:         afero.NewOsFs(),
		Exec:       exec.Interface(),
		Build:      build,
	})
}

func NewFoobarCommand(ctx *foobar.Ctx) *cobra.Command {
	cmd := &cobra.Command{
		Use: ctx.Build.AppName,
		PersistentPreRunE: func(c *cobra.Command, args []string) error {
			return errors.WithStack(ctx.Init())
		},
	}

	cli.AddLoggingFlags(cmd)

	cmd.AddCommand(
		cli.NewVersionCommand(ctx.IO, ctx.Build),
	)

	return cmd
}
