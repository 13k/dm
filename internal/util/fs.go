package util

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gobwas/glob"
)

var fsys fs.FS = os.DirFS("/") // for testing

const (
	latestFileModeLowerBound LatestFileMode = iota - 1
	LatestFileByName
	LatestFileByModTime
	latestFileModeUpperBound
)

type LatestFileMode int

func LatestFileModeFromString(s string) (LatestFileMode, error) {
	switch s {
	case "name":
		return LatestFileByName, nil
	case "modified":
		return LatestFileByModTime, nil
	}

	return latestFileModeLowerBound, fmt.Errorf("invalid latest file mode %q", s)
}

func (m LatestFileMode) String() string {
	switch m {
	case latestFileModeLowerBound, latestFileModeUpperBound:
	case LatestFileByName:
		return "name"
	case LatestFileByModTime:
		return "modified"
	}

	return "<invalid>"
}

func FindLatestFile(dirname string, mode LatestFileMode, exts []string) (fs.FileInfo, error) {
	if mode <= latestFileModeLowerBound || mode >= latestFileModeUpperBound {
		return nil, fmt.Errorf("invalid latest file mode (%v)", mode)
	}

	infos, err := readDirInfos(dirname)
	if err != nil {
		return nil, err
	}

	infos, err = filterFileInfosByExts(infos, exts)
	if err != nil {
		return nil, err
	}

	if len(infos) == 0 {
		return nil, nil
	}

	sortFileInfosByLatest(infos, mode)

	return infos[0], nil
}

func readDirInfos(dirname string) ([]fs.FileInfo, error) {
	if !filepath.IsAbs(dirname) {
		return nil, fmt.Errorf("path %q is not absolute", dirname)
	}

	// trim starting "/"
	dirname = dirname[1:]

	if dirname == "" {
		dirname = "."
	}

	entries, err := fs.ReadDir(fsys, dirname)
	if err != nil {
		return nil, err
	}

	result := make([]fs.FileInfo, len(entries))

	for i := range entries {
		result[i], err = entries[i].Info()
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func filterFileInfosByExts(infos []fs.FileInfo, exts []string) ([]fs.FileInfo, error) {
	if len(exts) == 0 {
		return infos, nil
	}

	var filtered []fs.FileInfo

	pattern := fmt.Sprintf("*.{%s}", strings.Join(exts, ","))
	g, err := glob.Compile(pattern)

	if err != nil {
		return nil, fmt.Errorf("error compiling glob pattern %q: %w", pattern, err)
	}

	for _, fi := range infos {
		if g.Match(fi.Name()) {
			filtered = append(filtered, fi)
		}
	}

	return filtered, nil
}

func sortFileInfosByLatest(infos []fs.FileInfo, mode LatestFileMode) {
	var sortBy sort.Interface

	switch mode {
	case latestFileModeLowerBound, latestFileModeUpperBound:
		panic("unreachable")
	case LatestFileByName:
		sortBy = FileInfosByNameDesc(infos)
	case LatestFileByModTime:
		sortBy = FileInfosByModTimeDesc(infos)
	}

	sort.Sort(sortBy)
}
