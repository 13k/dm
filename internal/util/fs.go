package util

import (
	"errors"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"github.com/gobwas/glob"

	"github.com/13k/dm/internal/markdown"
)

var testFs fs.FS //nolint:gochecknoglobals // Used in tests

var (
	ErrAbsPathRequired       = errors.New("absolute path is required")
	ErrInvalidLatestFileMode = errors.New("invalid latest file mode")
	ErrFileNotFound          = errors.New("file not found")
)

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

	return latestFileModeLowerBound, fmt.Errorf("%w: %q", ErrInvalidLatestFileMode, s)
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

func FindLatestFile(dir Path, mode LatestFileMode, exts []string) (fs.FileInfo, error) {
	if mode <= latestFileModeLowerBound || mode >= latestFileModeUpperBound {
		return nil, fmt.Errorf("%w: %v", ErrInvalidLatestFileMode, mode)
	}

	infos, err := readDirInfos(dir)
	if err != nil {
		return nil, err
	}

	infos, err = filterFileInfosByExts(infos, exts)
	if err != nil {
		return nil, err
	}

	if len(infos) == 0 {
		return nil, nil //nolint:nilnil // Valid optional return value
	}

	sortFileInfosByLatest(infos, mode)

	return infos[0], nil
}

func readDirInfos(dir Path) ([]fs.FileInfo, error) {
	if !dir.IsAbs() {
		return nil, fmt.Errorf("%w: %q", ErrAbsPathRequired, dir)
	}

	entries, err := dir.ReadDir()
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

	var (
		filtered []fs.FileInfo

		pattern = fmt.Sprintf("*.{%s}", strings.Join(exts, ","))
		g, err  = glob.Compile(pattern)
	)

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

func SearchLatestDocumentFile(dir Path, mode LatestFileMode) (Path, error) {
	latestInfo, err := FindLatestFile(dir, mode, markdown.Extensions)
	if err != nil {
		return "", fmt.Errorf("failed to find latest document file in directory %q: %w", dir, err)
	}

	if latestInfo == nil {
		return "", fmt.Errorf("failed to find latest document file by %s in directory %q: %w", mode, dir, ErrFileNotFound)
	}

	latestPath := dir.Join(latestInfo.Name())

	return latestPath, nil
}
