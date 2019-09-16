package cmd

import (
	"github.com/izumin5210/clig/pkg/clib"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"k8s.io/utils/exec"

	"go.example.com/foobar/pkg/foobar"
)

func NewDefaultFoobarCommand(wd clib.Path, build clib.Build) *cobra.Command {
	return NewFoobarCommand(&foobar.Ctx{
		WorkingDir: wd,
		IO:         clib.Stdio(),
		FS:         afero.NewOsFs(),
		Exec:       exec.New(),
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

	clib.SetIO(cmd, ctx.IO)
	clib.AddLoggingFlags(cmd)

	cmd.AddCommand(
		clib.NewVersionCommand(ctx.Build),
	)

	return cmd
}

