package cmd

import (
	"bytes"
	"context"
	"errors"
	"go/build"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

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
				{Path: root.Join(".gitignore").String(), Template: templateGitignore},
				{Path: root.Join("Makefile").String(), Template: templateMakefile},
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

			ctx := context.TODO()

			run := func(ctx context.Context, name string, args ...string) error {
				cmd := c.Exec.CommandContext(ctx, name, args...)
				cmd.SetStdin(c.IO.In())
				cmd.SetStdout(c.IO.Out())
				cmd.SetStderr(c.IO.Err())
				cmd.SetDir(root.String())
				zap.L().Debug("exec command", zap.String("cmd", name), zap.Strings("args", args), zap.Stringer("dir", root))
				return cmd.Run()
			}

			err = run(ctx, "dep", "init")
			if err != nil {
				return err
			}

			if _, err := c.Exec.LookPath("gex"); err != nil {
				err = run(ctx, "go", "get", "github.com/izumin5210/gex/cmd/gex")
				if err != nil {
					return err
				}
			}

			err = run(ctx, "gex", "--add", "github.com/mitchellh/gox")
			if err != nil {
				return err
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
		zap.L().Debug("create a directory", zap.String("dir", dir))
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

	zap.L().Debug("create a new flie", zap.String("path", e.Path))
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
	templateMakefile = mustCreateTemplate("makefile", `PATH := ${PWD}/bin:${PATH}
export PATH

.DEFAULT_GOAL := all

REVISION ?= $(shell git describe --always)
BUILD_DATE ?= $(shell date +'%Y-%m-%dT%H:%M:%SZ')

GO_BUILD_FLAGS := -v
GO_TEST_FLAGS := -v -timeout 2m
GO_COVER_FLAGS := -coverprofile coverage.txt -covermode atomic
SRC_FILES := $(shell go list -f {{"'{{range .GoFiles}}{{printf \"%s/%s\\n\" $$.Dir .}}{{end}}'"}} ./...)

XC_ARCH := 386 amd64
XC_OS := darwin linux windows


#  App
#----------------------------------------------------------------
BIN_DIR := ./bin
OUT_DIR := ./dist
GENERATED_BINS :=
PACKAGES :=

define cmd-tmpl

$(eval NAME := $(notdir $(1)))
$(eval OUT := $(addprefix $(BIN_DIR)/,$(NAME)))
$(eval LDFLAGS := -ldflags "-X main.revision=$(REVISION) -X main.buildDate=$(BUILD_DATE)")

$(OUT): $(SRC_FILES)
	go build $(GO_BUILD_FLAGS) $(LDFLAGS) -o $(OUT) $(1)

.PHONY: $(NAME)
$(NAME): $(OUT)

.PHONY: $(NAME)-package
$(NAME)-package: $(NAME)
	gox \
		$(LDFLAGS) \
		-os="$(XC_OS)" \
		-arch="$(XC_ARCH)" \
		-output="$(OUT_DIR)/$(NAME)_{{"{{.OS}}_{{.Arch}}"}}" \
		$(1)

$(eval GENERATED_BINS += $(OUT))
$(eval PACKAGES += $(NAME)-package)

endef

$(foreach src,$(wildcard ./cmd/*),$(eval $(call cmd-tmpl,$(src))))


#  Commands
#----------------------------------------------------------------
.PHONY: all
all: $(GENERATED_BINS)

.PHONY: packages
packages: $(PACKAGES)

.PHONY: setup
setup:
ifdef CI
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
endif
	dep ensure -v -vendor-only
	@go get github.com/izumin5210/gex/cmd/gex
	gex --build --verbose

.PHONY: clean
clean:
	rm -rf $(BIN_DIR)/*

.PHONY: gen
gen:
	go generate ./...

.PHONY: test
test:
	go test $(GO_TEST_FLAGS) ./...

.PHONY: cover
cover:
	go test $(GO_TEST_FLAGS) $(GO_COVER_FLAGS) ./...
`)
	templateGitignore = mustCreateTemplate("gitignore", `/bin
/dist
/vendor
`)
)
