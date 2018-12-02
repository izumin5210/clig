package main

import (
	"fmt"
	"os"

	"github.com/izumin5210/clig/pkg/cli"
	"github.com/izumin5210/clig/pkg/clig/cmd"
)

const (
	appName = "clig"
	version = "v0.1.0"
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
	cmd := cmd.NewDefaultCligCommand(cli.Build{
		AppName:   appName,
		Version:   version,
		Revision:  revision,
		BuildDate: buildDate,
	})

	return cmd.Execute()
}
