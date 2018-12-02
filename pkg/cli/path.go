package cli

import "path/filepath"

// Path represents a filepath in filesystem.
type Path string

func (p Path) String() string { return string(p) }

// Join joins path elements to the path.
func (p Path) Join(elem ...string) Path {
	return Path(filepath.Join(append([]string{p.String()}, elem...)...))
}
