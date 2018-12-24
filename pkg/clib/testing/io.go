package testing

import (
	"bytes"
	"io"
)

// FakeIO is a fake implementation of the IO interface using `bytes.Buffer`s.
type FakeIO struct {
	InBuf  *bytes.Buffer
	OutBuf *bytes.Buffer
	ErrBuf *bytes.Buffer
}

func (i *FakeIO) In() io.Reader  { return i.InBuf }
func (i *FakeIO) Out() io.Writer { return i.OutBuf }
func (i *FakeIO) Err() io.Writer { return i.ErrBuf }

// NewFakeIO returns a new FakeIO object.
func NewFakeIO() *FakeIO {
	return &FakeIO{
		InBuf:  new(bytes.Buffer),
		OutBuf: new(bytes.Buffer),
		ErrBuf: new(bytes.Buffer),
	}
}
