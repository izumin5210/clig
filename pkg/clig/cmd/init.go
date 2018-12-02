package cmd

import (
	"bytes"
	"errors"
	"go/build"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/izumin5210/clig/pkg/clig"
)

var (
	BuildContext = build.Default
)

func newInitCommand(c *clig.Ctx) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "init",
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			name := args[0]
			root := c.WorkingDir.Join(name)
			pkg, err := getImportPath(root.String())
			if err != nil {
				return err
			}

			params := struct {
				Name    string
				Package string
			}{Name: name, Package: pkg}

			entries := []*entry{
				{Path: root.Join("cmd", params.Name, "main.go").String(), Template: templateMain},
				{Path: root.Join("pkg", params.Name, "context.go").String(), Template: templateCtx},
				{Path: root.Join("pkg", params.Name, "cmd", "cmd.go").String(), Template: templateCmd},
			}

			for _, e := range entries {
				err = e.Create(c.FS, params)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	return cmd
}

func getImportPath(rootPath string) (importPath string, err error) {
	for _, gopath := range filepath.SplitList(BuildContext.GOPATH) {
		prefix := filepath.Join(gopath, "src") + string(filepath.Separator)
		// FIXME: should not use strings.HasPrefix
		if strings.HasPrefix(rootPath, prefix) {
			importPath = filepath.ToSlash(strings.Replace(rootPath, prefix, "", 1))
			break
		}
	}
	if importPath == "" {
		err = errors.New("failed to get the import path")
	}
	return
}

type entry struct {
	Template *template.Template
	Path     string
}

func (e *entry) Create(fs afero.Fs, params interface{}) error {
	dir := filepath.Dir(e.Path)
	if ok, err := afero.DirExists(fs, dir); err != nil {
		return err
	} else if !ok {
		err = fs.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	buf := new(bytes.Buffer)
	err := e.Template.Execute(buf, params)
	if err != nil {
		return err
	}

	err = afero.WriteFile(fs, e.Path, buf.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}

func mustCreateTemplate(name, tmpl string) *template.Template {
	return template.Must(template.New(name).Funcs(funcMap).Parse(tmpl))
}

var (
	funcMap      = template.FuncMap{"ToCamel": strcase.ToCamel}
	templateMain = mustCreateTemplate("main", `package main

import (
	"fmt"
	"os"

	"github.com/izumin5210/clig/pkg/cli"

	"{{.Package}}/pkg/{{.Name}}/cmd"
)

const (
	appName = "{{.Name}}"
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
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	cmd := cmd.NewDefault{{ToCamel .Name}}Command(cli.Path(wd), cli.Build{
		AppName:   appName,
		Version:   version,
		Revision:  revision,
		BuildDate: buildDate,
	})

	return cmd.Execute()
}
`)
	templateCtx = mustCreateTemplate("ctx", `package {{.Name}}

import (
	"github.com/izumin5210/clig/pkg/cli"
)

type Ctx struct {
	WorkingDir cli.Path
	IO         cli.IO

	Build cli.Build
}
`)
	templateCmd = mustCreateTemplate("cmd", `package cmd

import (
	"github.com/spf13/cobra"
	"github.com/izumin5210/clig/pkg/cli"

	"{{.Package}}/pkg/{{.Name}}"
)

func NewDefault{{ToCamel .Name}}Command(wd cli.Path, build cli.Build) *cobra.Command {
	return New{{ToCamel .Name}}Command(&{{.Name}}.Ctx{
		WorkingDir: wd,
		IO:         cli.Stdio(),
		Build:      build,
	})
}

func New{{ToCamel .Name}}Command(ctx *{{.Name}}.Ctx) *cobra.Command {
	cmd := &cobra.Command{
		Use: ctx.Build.AppName,
	}

	cli.AddLoggingFlags(cmd)

	cmd.AddCommand(
		cli.NewVersionCommand(ctx.IO, ctx.Build),
	)

	return cmd
}
`)
)
