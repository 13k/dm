package util_test

import (
	"errors"
	"fmt"
	"io/fs"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/13k/dm/internal/util"
)

var (
	actualFS = util.GetFS()
	testFS   = fstest.MapFS{}

	rootTestFS util.Path
)

func setupFS() {
	if testFS == nil {
		testFS = fstest.MapFS{}
	}

	util.SetFS(testFS)
}

func teardownFS() {
	util.SetFS(actualFS)
}

func createFile(path util.Path, data []byte, mode fs.FileMode, mtime time.Time) *fstest.MapFile {
	if path == "" {
		panic(errors.New("cannot create file with empty name"))
	}

	if mode == 0 {
		mode = 0o666
	}

	if mtime.IsZero() {
		mtime = time.Now()
	}

	f := &fstest.MapFile{
		Data:    data,
		Mode:    mode,
		ModTime: mtime,
	}

	testFS[path.String()] = f

	return f
}

func createEmptyFile(path util.Path, mode fs.FileMode, mtime time.Time) *fstest.MapFile {
	return createFile(path, nil, mode, mtime)
}

func touchFile(path util.Path, mtime time.Time) *fstest.MapFile { //nolint:unparam
	return createEmptyFile(path, 0, mtime)
}

func TestFindLatestFile(t *testing.T) {
	t.Cleanup(teardownFS)

	setupFS()

	touchFile("dm.md", time.Now())
	touchFile("dm.mkd", time.Now().Add(1*time.Second))
	touchFile("latest-by-name.md", time.Now().Add(2*time.Second))
	touchFile("latest-by-name.txt", time.Now().Add(3*time.Second))
	touchFile("latest-by-mtime.md", time.Now().Add(4*time.Second))
	touchFile("latest-by-mtime.txt", time.Now().Add(5*time.Second))

	entries, err := rootTestFS.ReadDir()

	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		fi, err := entry.Info()

		if err != nil {
			panic(err)
		}

		eType := "f"

		if fi.IsDir() {
			eType = "d"
		}

		fmt.Printf("%30s  %s  %s\n", fi.Name(), eType, fi.ModTime())
	}

	mdExts := []string{"md", "mkd"}

	testCases := []struct {
		dirname  util.Path
		mode     util.LatestFileMode
		exts     []string
		expected string
		err      string
	}{
		{rootTestFS, util.LatestFileMode(-1), nil, "", "invalid latest file mode"},
		{rootTestFS, util.LatestFileMode(666), nil, "", "invalid latest file mode"},
		{"invalid", util.LatestFileByName, nil, "", "is not absolute"},
		{rootTestFS.Join("invalid"), util.LatestFileByName, nil, "", "file does not exist"},
		{rootTestFS, util.LatestFileByName, []string{"[1-]"}, "", "pattern"},
		{rootTestFS, util.LatestFileByModTime, []string{"invalid"}, "", ""},
		{rootTestFS, util.LatestFileByModTime, nil, "latest-by-mtime.txt", ""},
		{rootTestFS, util.LatestFileByModTime, mdExts, "latest-by-mtime.md", ""},
		{rootTestFS, util.LatestFileByModTime, []string{"mkd"}, "dm.mkd", ""},
		{rootTestFS, util.LatestFileByName, nil, "latest-by-name.txt", ""},
		{rootTestFS, util.LatestFileByName, mdExts, "latest-by-name.md", ""},
	}

	for tcidx, tc := range testCases {
		actual, err := util.FindLatestFile(tc.dirname, tc.mode, tc.exts)

		if tc.err == "" {
			if err != nil {
				t.Errorf("case #%d: unexpected error: %v", tcidx, err)
				continue
			}
		} else {
			if err == nil {
				t.Errorf("case #%d: expected error matching %q, actual: nil", tcidx, tc.err)
				continue
			}

			if !strings.Contains(err.Error(), tc.err) {
				t.Errorf("case #%d: expected error matching %q, actual: %q", tcidx, tc.err, err.Error())
				continue
			}
		}

		if tc.expected == "" {
			if actual != nil {
				t.Errorf("case #%d: expected nil, actual: %q", tcidx, actual.Name())
				continue
			}
		} else {
			if actual == nil {
				t.Errorf("case #%d: expected: %q, actual: nil", tcidx, tc.expected)
				continue
			}

			if actual.Name() != tc.expected {
				t.Errorf("case #%d: expected %q, actual: %q", tcidx, tc.expected, actual.Name())
				continue
			}
		}
	}
}
