package config

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/spf13/viper"

	"github.com/13k/dm/internal/util"
	"github.com/13k/dm/meta"
)

var ErrConfig = errors.New("configuration error")

func init() {
	defaultConfigDir := filepath.Join(xdg.ConfigHome, meta.AppName)
	defaultOutputPath := util.Cwd().String()
	defaultLatestMode := util.LatestFileByName.String()

	viper.AddConfigPath(defaultConfigDir)
	viper.SetConfigName("config")

	viper.SetDefault("output_path", defaultOutputPath)
	viper.SetDefault("latest_mode", defaultLatestMode)
}

func LoadFile() error {
	err := viper.ReadInConfig()

	if err != nil {
		var notExistErr viper.ConfigFileNotFoundError

		if !errors.As(err, &notExistErr) {
			return err
		}
	} else {
		log.Printf("config.LoadFile() -- loaded configuration from file %q", viper.ConfigFileUsed())
	}

	return nil
}

type Config struct {
	InputPath    util.Path
	OutputPath   util.Path
	Latest       bool
	LatestMode   string
	SlackChannel string

	LatestBy util.LatestFileMode
}

func New() (*Config, error) { //nolint:cyclop,funlen
	cfg := &Config{
		InputPath:    util.NewPath(viper.GetString("input_path")),
		OutputPath:   util.NewPath(viper.GetString("output_path")),
		Latest:       viper.GetBool("latest"),
		LatestMode:   viper.GetString("latest_mode"),
		SlackChannel: viper.GetString("slack_channel"),
	}

	var (
		err   error
		isDir bool
	)

	if cfg.OutputPath == "" {
		return nil, fmt.Errorf("%w: output path cannot be empty", ErrConfig)
	}

	cfg.LatestBy, err = util.LatestFileModeFromString(cfg.LatestMode)
	if err != nil {
		return nil, err
	}

	if cfg.InputPath == "" && cfg.Latest {
		cfg.InputPath = util.Cwd()
	}

	if cfg.InputPath != "" {
		cfg.InputPath, err = cfg.InputPath.Abs()
		if err != nil {
			return nil, err
		}
	}

	if cfg.InputPath != "" {
		isDir, err = cfg.InputPath.IsDir()
		if err != nil {
			return nil, err
		}

		if isDir {
			cfg.InputPath, err = util.SearchLatestDocumentFile(cfg.InputPath, cfg.LatestBy)
			if err != nil {
				return nil, err
			}
		}
	}

	isDir, err = cfg.OutputPath.IsDir()
	if err != nil {
		return nil, err
	}

	if isDir {
		defaultOutputFilename := util.TodayString() + ".md"
		cfg.OutputPath = cfg.OutputPath.Join(defaultOutputFilename)
	} else {
		cfg.OutputPath, err = cfg.OutputPath.Abs()
		if err != nil {
			return nil, err
		}
	}

	return cfg, nil
}
