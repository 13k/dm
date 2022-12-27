package util

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type Path string

func NewPath(p string) Path {
	return Path(filepath.ToSlash(p))
}

func (p Path) String() string {
	return string(p)
}

func (p Path) Join(parts ...string) Path {
	args := []string{p.String()}
	args = append(args, parts...)

	return NewPath(filepath.Join(args...))
}

func (p Path) Root() (Path, error) {
	abs, err := filepath.Abs(p.String())

	if err != nil {
		return "", err
	}

	root := "/"
	vol := filepath.VolumeName(abs)

	if vol != "" {
		root = vol + "/"
	}

	return NewPath(root), nil
}

func (p Path) IsAbs() bool {
	return filepath.IsAbs(p.String())
}

func (p Path) Abs() (Path, error) {
	abs, err := filepath.Abs(p.String())

	if err != nil {
		return "", err
	}

	return NewPath(abs), nil
}

func (p Path) Rel(parent Path) (Path, error) {
	rel, err := filepath.Rel(parent.String(), p.String())

	if err != nil {
		return "", err
	}

	return NewPath(rel), nil
}

func (p Path) FS() fs.FS {
	if testFs != nil {
		return testFs
	}

	return os.DirFS(p.String())
}

func (p Path) IsDir() (bool, error) {
	fi, err := os.Stat(p.String())

	if err != nil {
		return false, fmt.Errorf("failed to get filesystem information from %q: %w", p, err)
	}

	return fi.IsDir(), nil
}

func (p Path) ReadDir() ([]fs.DirEntry, error) {
	root, err := p.Root()

	if err != nil {
		return nil, err
	}

	fsys := root.FS()
	rel, err := p.Rel(root)

	if err != nil {
		return nil, err
	}

	entries, err := fs.ReadDir(fsys, rel.String())

	if err != nil {
		return nil, err
	}

	return entries, nil
}
