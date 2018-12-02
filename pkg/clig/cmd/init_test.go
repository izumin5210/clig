package cmd

import (
	"os"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/izumin5210/clig/pkg/cli"
	clitesting "github.com/izumin5210/clig/pkg/cli/testing"
	"github.com/spf13/afero"

	"github.com/izumin5210/clig/pkg/clig"
)

func TestInit(t *testing.T) {
	cases := []struct {
		test  string
		args  []string
		files []string
	}{
		{
			test: "simple",
			args: []string{"foobar"},
			files: []string{
				"foobar/cmd/foobar/main.go",
				"foobar/pkg/foobar/context.go",
				"foobar/pkg/foobar/cmd/cmd.go",
			},
		},
	}

	defer func(p string) { BuildContext.GOPATH = p }(BuildContext.GOPATH)
	BuildContext.GOPATH = "/home/go"

	for _, tc := range cases {
		wd := cli.Path("/home/go/src/go.example.com")

		t.Run(tc.test, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			ctx := &clig.Ctx{
				WorkingDir: wd,
				IO:         clitesting.NewFakeIO(),
				FS:         fs,
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
		})
	}
}
