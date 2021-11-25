package util_test

import (
	"errors"
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

func createFile(name string, data []byte, mode fs.FileMode, mtime time.Time) *fstest.MapFile {
	if name == "" {
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

	testFS[name] = f

	return f
}

func createEmptyFile(name string, mode fs.FileMode, mtime time.Time) *fstest.MapFile {
	return createFile(name, nil, mode, mtime)
}

func touchFile(name string) *fstest.MapFile { //nolint:unparam
	return createEmptyFile(name, 0, time.Time{})
}

func TestFindLatestFile(t *testing.T) {
	t.Cleanup(teardownFS)

	setupFS()

	touchFile("dm.md")
	touchFile("dm.mkd")
	touchFile("latest-by-name.md")
	touchFile("latest-by-name.txt")
	touchFile("latest-by-mtime.md")
	touchFile("latest-by-mtime.txt")

	mdExts := []string{"md", "mkd"}

	testCases := []struct {
		dirname  string
		mode     util.LatestFileMode
		exts     []string
		expected string
		err      string
	}{
		{"/", util.LatestFileMode(-1), nil, "", "invalid latest file mode"},
		{"/", util.LatestFileMode(666), nil, "", "invalid latest file mode"},
		{"invalid", util.LatestFileByName, nil, "", "is not absolute"},
		{"/invalid", util.LatestFileByName, nil, "", "file does not exist"},
		{"/", util.LatestFileByName, []string{"[1-]"}, "", "pattern"},
		{"/", util.LatestFileByModTime, []string{"invalid"}, "", ""},
		{"/", util.LatestFileByModTime, nil, "latest-by-mtime.txt", ""},
		{"/", util.LatestFileByModTime, mdExts, "latest-by-mtime.md", ""},
		{"/", util.LatestFileByModTime, []string{"mkd"}, "dm.mkd", ""},
		{"/", util.LatestFileByName, nil, "latest-by-name.txt", ""},
		{"/", util.LatestFileByName, mdExts, "latest-by-name.md", ""},
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
