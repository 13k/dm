package util

import (
	"io/fs"
)

func GetFS() fs.FS {
	return fsys
}

func SetFS(testFsys fs.FS) {
	fsys = testFsys
}
