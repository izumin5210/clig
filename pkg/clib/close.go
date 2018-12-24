package clib

var closers []func()

// Close closes cli utilities.
func Close() {
	for _, f := range closers {
		f()
	}
}

func addCloseFunc(f func()) {
	closers = append(closers, f)
}
