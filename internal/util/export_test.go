package util

import (
	"io/fs"
)

func GetFS() fs.FS {
	return testFs
}

func SetFS(fsys fs.FS) {
	testFs = fsys
}
