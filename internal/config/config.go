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

var (
	defaultConfigDir      string
	defaultOutputFilename string
	defaultOutputPath     string
	defaultLatestMode     string = util.LatestFileByName.String()
)

func init() {
	defaultConfigDir = filepath.Join(xdg.ConfigHome, meta.AppName)
	defaultOutputFilename = fmt.Sprintf("%s.md", util.TodayString())
	defaultOutputPath = filepath.Join(util.Cwd, defaultOutputFilename)

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
	InputPath    string
	OutputPath   string
	Latest       bool
	LatestMode   string
	SlackChannel string

	LatestBy util.LatestFileMode
}

func New() (*Config, error) { //nolint: funlen
	cfg := &Config{
		InputPath:    viper.GetString("input_path"),
		OutputPath:   viper.GetString("output_path"),
		Latest:       viper.GetBool("latest"),
		LatestMode:   viper.GetString("latest_mode"),
		SlackChannel: viper.GetString("slack_channel"),
	}

	var (
		err   error
		isDir bool
	)

	if cfg.OutputPath == "" {
		return nil, fmt.Errorf("output path cannot be empty")
	}

	cfg.LatestBy, err = util.LatestFileModeFromString(cfg.LatestMode)

	if err != nil {
		return nil, err
	}

	if cfg.InputPath == "" && cfg.Latest {
		cfg.InputPath = util.Cwd
	}

	if cfg.InputPath != "" { //nolint: nestif
		isDir, err = util.IsDir(cfg.InputPath)

		if err != nil {
			return nil, err
		}

		if isDir {
			cfg.InputPath, err = util.SearchLatestDocumentFile(cfg.InputPath, cfg.LatestBy)

			if err != nil {
				return nil, err
			}
		} else {
			cfg.InputPath, err = util.AbsFilepath(cfg.InputPath)

			if err != nil {
				return nil, err
			}
		}
	}

	isDir, err = util.IsDir(cfg.OutputPath)

	if err != nil {
		return nil, err
	}

	if isDir {
		cfg.OutputPath = filepath.Join(cfg.OutputPath, defaultOutputFilename)
	} else {
		cfg.OutputPath, err = util.AbsFilepath(cfg.OutputPath)

		if err != nil {
			return nil, err
		}
	}

	return cfg, nil
}
