package cmd

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/izumin5210/clig/pkg/cli"
	clitesting "github.com/izumin5210/clig/pkg/cli/testing"
	"github.com/spf13/afero"
	"k8s.io/utils/exec"
	exectesting "k8s.io/utils/exec/testing"

	"github.com/izumin5210/clig/pkg/clig"
)

func TestInit(t *testing.T) {
	defer func(p string) { BuildContext.GOPATH = p }(BuildContext.GOPATH)
	BuildContext.GOPATH = "/home/go"
	wd := cli.Path("/home/go/src/go.example.com")

	createFakeCmd := func(name string, args ...string) *exectesting.FakeCmd {
		cmd := &exectesting.FakeCmd{
			RunScript: []exectesting.FakeRunAction{func() ([]byte, []byte, error) { return nil, nil, nil }},
		}
		exectesting.InitFakeCmd(cmd, name, args...)
		return cmd
	}

	cases := []struct {
		test         string
		args         []string
		files        []string
		excmds       []*exectesting.FakeCmd
		lookPathFunc func(string) (string, error)
	}{
		{
			test: "simple",
			args: []string{"foobar"},
			files: []string{
				"foobar/.gitignore",
				"foobar/.reviewdog.yml",
				"foobar/.travis.yml",
				"foobar/Makefile",
				"foobar/cmd/foobar/main.go",
				"foobar/pkg/foobar/config.go",
				"foobar/pkg/foobar/context.go",
				"foobar/pkg/foobar/cmd/cmd.go",
			},
			excmds: []*exectesting.FakeCmd{
				createFakeCmd("dep", "init"),
				createFakeCmd("gex",
					"--add", "github.com/mitchellh/gox",
					"--add", "github.com/haya14busa/reviewdog/cmd/reviewdog",
					"--add", "github.com/kisielk/errcheck",
					"--add", "github.com/srvc/wraperr/cmd/wraperr",
					"--add", "golang.org/x/lint/golint",
					"--add", "honnef.co/go/tools/cmd/megacheck",
					"--add", "mvdan.cc/unparam",
				),
			},
		},
		{
			test: "gex has already been installed",
			args: []string{"foobar"},
			files: []string{
				"foobar/.gitignore",
				"foobar/.reviewdog.yml",
				"foobar/.travis.yml",
				"foobar/Makefile",
				"foobar/cmd/foobar/main.go",
				"foobar/pkg/foobar/config.go",
				"foobar/pkg/foobar/context.go",
				"foobar/pkg/foobar/cmd/cmd.go",
			},
			excmds: []*exectesting.FakeCmd{
				createFakeCmd("dep", "init"),
				createFakeCmd("go", "get", "github.com/izumin5210/gex/cmd/gex"),
				createFakeCmd("gex",
					"--add", "github.com/mitchellh/gox",
					"--add", "github.com/haya14busa/reviewdog/cmd/reviewdog",
					"--add", "github.com/kisielk/errcheck",
					"--add", "github.com/srvc/wraperr/cmd/wraperr",
					"--add", "golang.org/x/lint/golint",
					"--add", "honnef.co/go/tools/cmd/megacheck",
					"--add", "mvdan.cc/unparam",
				),
			},
			lookPathFunc: func(cmd string) (string, error) {
				if cmd == "gex" {
					return "", errors.New("error")
				}
				return cmd, nil
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.test, func(t *testing.T) {
			fs := afero.NewMemMapFs()

			fexec := &exectesting.FakeExec{LookPathFunc: tc.lookPathFunc}
			if fexec.LookPathFunc == nil {
				fexec.LookPathFunc = func(cmd string) (string, error) { return cmd, nil }
			}
			for _, c := range tc.excmds {
				c := c
				fexec.CommandScript = append(fexec.CommandScript, func(cmd string, args ...string) exec.Cmd {
					if got, want := append([]string{cmd}, args...), c.Argv; !reflect.DeepEqual(got, want) {
						t.Errorf("called command is %v, want %v", got, want)
					}
					return c
				})
			}

			ctx := &clig.Ctx{
				WorkingDir: wd,
				IO:         clitesting.NewFakeIO(),
				FS:         fs,
				Exec:       fexec,
			}

			cmd := newInitCommand(ctx)
			cmd.SetArgs(tc.args)
			err := cmd.Execute()

			if err != nil {
				t.Fatalf("failed to execute command: %v", err)
			}

			files := make(map[string]struct{}, len(tc.files))
			for _, f := range tc.files {
				files[wd.Join(f).String()] = struct{}{}
			}

			afero.Walk(fs, wd.String(), func(path string, info os.FileInfo, err error) error {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if info.IsDir() {
					return nil
				}
				if _, ok := files[path]; ok {
					delete(files, path)
					t.Run(path, func(t *testing.T) {
						data, err := afero.ReadFile(fs, path)
						if err != nil {
							t.Errorf("failed to read %q: %v", path, err)
						}
						cupaloy.SnapshotT(t, string(data))
					})
				} else {
					t.Errorf("unexpected file is created: %q", path)
				}
				return nil
			})

			if got, want := fexec.CommandCalls, len(tc.excmds); got != want {
				t.Errorf("called %d external commands, want %d", got, want)
			}
		})
	}
}
