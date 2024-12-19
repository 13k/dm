package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

func configureLogger() error {
	logPath := viper.GetString("log_path")

	if logPath == "" {
		logPath = os.DevNull
	}

	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o666)
	if err != nil {
		return fmt.Errorf("could not open log file %q: %w", logPath, err)
	}

	log.SetOutput(f)

	return nil
}
