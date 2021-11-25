package cli

import (
	"fmt"
	"path/filepath"

	"github.com/13k/dm/internal/markdown"
	"github.com/13k/dm/internal/util"
)

type options struct {
	inputPath  string
	outputPath string
}

func newOptions(raw *rawOptions) (*options, error) {
	var err error

	opts := &options{}

	if raw.outputPath == "" {
		return nil, fmt.Errorf("output path cannot be empty")
	}

	if raw.basePath != "" && raw.latest {
		return nil, fmt.Errorf("base path and latest are mutually exclusive")
	}

	opts.outputPath, err = filepath.Abs(raw.outputPath)
	if err != nil {
		return nil, fmt.Errorf("could not determine absolute path to %q: %w", raw.outputPath, err)
	}

	latestBy, err := util.LatestFileModeFromString(raw.latestMode)
	if err != nil {
		return nil, fmt.Errorf("error parsing latest file mode: %w", err)
	}

	opts.inputPath = raw.basePath

	if opts.inputPath != "" {
		opts.inputPath, err = filepath.Abs(opts.inputPath)
		if err != nil {
			return nil, fmt.Errorf("could not determine absolute path to %q: %w", opts.inputPath, err)
		}
	} else if raw.latest {
		searchDir := filepath.Dir(opts.outputPath)
		latestInfo, err := util.FindLatestFile(searchDir, latestBy, markdown.Extensions)

		if err != nil {
			return nil, fmt.Errorf("could not find latest file: %w", err)
		}

		if latestInfo == nil {
			return nil, fmt.Errorf("could not find any latest file by %s in %q", latestBy, searchDir)
		}

		opts.inputPath = filepath.Join(searchDir, latestInfo.Name())
	}

	return opts, nil
}
