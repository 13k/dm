package util

import (
	"io/fs"
	"sort"
)

var _ sort.Interface = FileInfosByNameAsc(nil)

type FileInfosByNameAsc []fs.FileInfo

func (s FileInfosByNameAsc) Len() int           { return len(s) }
func (s FileInfosByNameAsc) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s FileInfosByNameAsc) Less(i, j int) bool { return s[i].Name() < s[j].Name() }

var _ sort.Interface = FileInfosByNameDesc(nil)

type FileInfosByNameDesc []fs.FileInfo

func (s FileInfosByNameDesc) Len() int           { return len(s) }
func (s FileInfosByNameDesc) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s FileInfosByNameDesc) Less(i, j int) bool { return s[i].Name() > s[j].Name() }

var _ sort.Interface = FileInfosByModTimeAsc(nil)

type FileInfosByModTimeAsc []fs.FileInfo

func (s FileInfosByModTimeAsc) Len() int      { return len(s) }
func (s FileInfosByModTimeAsc) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s FileInfosByModTimeAsc) Less(i, j int) bool {
	return s[i].ModTime().UnixNano() < s[j].ModTime().UnixNano()
}

var _ sort.Interface = FileInfosByModTimeDesc(nil)

type FileInfosByModTimeDesc []fs.FileInfo

func (s FileInfosByModTimeDesc) Len() int      { return len(s) }
func (s FileInfosByModTimeDesc) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s FileInfosByModTimeDesc) Less(i, j int) bool {
	return s[i].ModTime().UnixNano() > s[j].ModTime().UnixNano()
}
