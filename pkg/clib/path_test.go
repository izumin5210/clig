package clib_test

import (
	"testing"

	"github.com/izumin5210/clig/pkg/clib"
)

func TestPath_String(t *testing.T) {
	pathStr := "/go/src/awesomeapp"
	path := clib.Path(pathStr)

	if got, want := path.String(), pathStr; got != want {
		t.Errorf("String() returned %q, want %q", got, want)
	}
}

func TestPath_Join(t *testing.T) {
	path := clib.Path("/go/src/awesomeapp")

	if got, want := path.Join("cmd", "server"), clib.Path("/go/src/awesomeapp/cmd/server"); got != want {
		t.Errorf("Join() returned %q, want %q", got, want)
	}
}
