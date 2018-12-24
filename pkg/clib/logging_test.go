package clib_test

import (
	"strings"
	"testing"

	"github.com/izumin5210/clig/pkg/clib"
	"github.com/spf13/cobra"
)

func TestLogging(t *testing.T) {
	cases := []struct {
		args      []string
		mode      clib.LoggingMode
		isDebug   bool
		isVerbose bool
	}{
		{
			mode: clib.LoggingNop,
		},
		{
			args:      []string{"-v"},
			mode:      clib.LoggingVerbose,
			isVerbose: true,
		},
		{
			args:      []string{"--verbose"},
			mode:      clib.LoggingVerbose,
			isVerbose: true,
		},
		{
			args:    []string{"--debug"},
			mode:    clib.LoggingDebug,
			isDebug: true,
		},
	}

	for _, tc := range cases {
		t.Run(strings.Join(tc.args, " "), func(t *testing.T) {
			defer clib.Close()

			var (
				mode               clib.LoggingMode
				isDebug, isVerbose bool
			)

			cmd := &cobra.Command{
				Run: func(*cobra.Command, []string) {
					mode = clib.Logging()
					isDebug = clib.IsDebug()
					isVerbose = clib.IsVerbose()
				},
			}

			clib.AddLoggingFlags(cmd)
			cmd.SetArgs(tc.args)
			err := cmd.Execute()

			if err != nil {
				t.Errorf("Execute() returned an error: %v", err)
			}

			if got, want := mode, tc.mode; got != want {
				t.Errorf("LoggingMode() returned %v, want %v", got, want)
			}

			if got, want := isVerbose, tc.isVerbose; got != want {
				t.Errorf("IsVerbose() returned %t, want %t", got, want)
			}

			if got, want := isDebug, tc.isDebug; got != want {
				t.Errorf("IsDebug() returned %t, want %t", got, want)
			}
		})
	}
}
