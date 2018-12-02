package main

import (
	"fmt"
	"os"

	"github.com/izumin5210/clig/pkg/cli"

	"go.example.com/foobar/pkg/foobar/cmd"
)

const (
	appName = "foobar"
	version = "v0.0.1"
)

var (
	revision, buildDate string
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	cmd := cmd.NewDefaultFoobarCommand(cli.Build{
		AppName:   appName,
		Version:   version,
		Revision:  revision,
		BuildDate: buildDate,
	})

	return cmd.Execute()
}

