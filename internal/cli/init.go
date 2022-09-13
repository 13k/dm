package cli

import (
	"github.com/spf13/cobra"

	"github.com/13k/dm/internal/config"
	"github.com/13k/dm/internal/util"
)

func init() {
	cobra.OnInitialize(onInit)
}

func onInit() {
	if err := configureLogger(); err != nil {
		util.Fatal("failed to configure logger: %v", err)
	}

	if err := config.LoadFile(); err != nil {
		util.Fatal("failed to load configuration file: %v", err)
	}
}
