package clib

import (
	"io"
	"os"
	"runtime"

	colorable "github.com/mattn/go-colorable"
	"github.com/spf13/cobra"
)

// IO contains an input reader, an output writer and an error writer.
type IO struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}

// Stdio returns a standard IO object.
func Stdio() *IO {
	var (
		inR  io.Reader = os.Stdin
		outW io.Writer = os.Stdout
		errW io.Writer = os.Stderr
	)
	if runtime.GOOS == "windows" {
		outW = colorable.NewColorableStdout()
		errW = colorable.NewColorableStderr()
	}
	return NewIO(inR, outW, errW)
}

func NewIO(inR io.Reader, outW io.Writer, errW io.Writer) *IO {
	return &IO{
		In:  inR,
		Out: outW,
		Err: errW,
	}
}

// SetIO set an IO to *cobra.Command.
func SetIO(c *cobra.Command, io *IO) {
	c.SetIn(io.In)
	c.SetOut(io.Out)
	c.SetErr(io.Err)
}

// GetIO extract an IO object from *cobra.Command.
func GetIO(c *cobra.Command) *IO {
	return NewIO(
		c.InOrStdin(),
		c.OutOrStdout(),
		c.ErrOrStderr(),
	)
}
