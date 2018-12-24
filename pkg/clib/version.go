package clib

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

// Build is a container for the application build information.
type Build struct {
	AppName   string
	Version   string
	Revision  string
	BuildDate string
}

// NewVersionCommand create a new cobra.Command to print the version information.
func NewVersionCommand(io IO, cfg Build) *cobra.Command {
	return &cobra.Command{
		Use:           "version",
		Short:         "Print the version information",
		Long:          "Print the version information",
		SilenceErrors: true,
		SilenceUsage:  true,
		Run: func(cmd *cobra.Command, _ []string) {
			buf := bytes.NewBufferString(cfg.AppName + " " + cfg.Version)
			buf.WriteString(" (")
			var meta []string
			for _, c := range []string{cfg.BuildDate, cfg.Revision, runtime.Version(), fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)} {
				if c != "" {
					meta = append(meta, c)
				}
			}
			buf.WriteString(strings.Join(meta, " "))
			buf.WriteString(")")
			fmt.Fprintln(io.Out(), buf.String())
		},
	}
}
