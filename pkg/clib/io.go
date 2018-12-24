package clib

import (
	"io"
	"os"
	"runtime"

	colorable "github.com/mattn/go-colorable"
)

// IO contains an input reader, an output writer and an error writer.
type IO interface {
	In() io.Reader
	Out() io.Writer
	Err() io.Writer
}

// IOContainer is a basic implementation of the IO interface.
type IOContainer struct {
	InR  io.Reader
	OutW io.Writer
	ErrW io.Writer
}

func (i *IOContainer) In() io.Reader  { return i.InR }
func (i *IOContainer) Out() io.Writer { return i.OutW }
func (i *IOContainer) Err() io.Writer { return i.ErrW }

// Stdio returns a standard IO object.
func Stdio() IO {
	io := &IOContainer{
		InR:  os.Stdin,
		OutW: os.Stdout,
		ErrW: os.Stderr,
	}
	if runtime.GOOS == "windows" {
		io.OutW = colorable.NewColorableStdout()
		io.ErrW = colorable.NewColorableStderr()
	}
	return io
}
