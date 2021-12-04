package cli

import (
	"fmt"
	"path/filepath"

	"github.com/13k/dm/internal/config"
	"github.com/13k/dm/internal/markdown"
	"github.com/13k/dm/internal/util"
)

type options struct {
	logPath      string
	basePath     string
	outputPath   string
	latest       bool
	latestMode   string
	slackChannel string
}

func parseOptions(opts *options) (*config.Config, error) {
	var err error

	cfg := &config.Config{
		SlackChannel: opts.slackChannel,
	}

	if opts.outputPath == "" {
		return nil, fmt.Errorf("output path cannot be empty")
	}

	if opts.basePath != "" && opts.latest {
		return nil, fmt.Errorf("base path and latest are mutually exclusive")
	}

	cfg.OutputPath, err = filepath.Abs(opts.outputPath)
	if err != nil {
		return nil, fmt.Errorf("could not determine absolute path to %q: %w", opts.outputPath, err)
	}

	latestBy, err := util.LatestFileModeFromString(opts.latestMode)
	if err != nil {
		return nil, fmt.Errorf("error parsing latest file mode: %w", err)
	}

	cfg.InputPath = opts.basePath

	if cfg.InputPath != "" {
		cfg.InputPath, err = filepath.Abs(cfg.InputPath)
		if err != nil {
			return nil, fmt.Errorf("could not determine absolute path to %q: %w", cfg.InputPath, err)
		}
	} else if opts.latest {
		searchDir := filepath.Dir(cfg.OutputPath)
		latestInfo, err := util.FindLatestFile(searchDir, latestBy, markdown.Extensions)

		if err != nil {
			return nil, fmt.Errorf("could not find latest file: %w", err)
		}

		if latestInfo == nil {
			return nil, fmt.Errorf("could not find any latest file by %s in %q", latestBy, searchDir)
		}

		cfg.InputPath = filepath.Join(searchDir, latestInfo.Name())
	}

	return cfg, nil
}
