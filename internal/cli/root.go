package cli

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/13k/dm/internal/config"
	"github.com/13k/dm/internal/ui/app"
	"github.com/13k/dm/internal/util"
	"github.com/13k/dm/meta"
)

func rootCmd() *cobra.Command { //nolint:funlen
	cmd := &cobra.Command{
		Use:   "dm [flags]",
		Short: "Create daily meeting notes",
		Long: `Create daily meeting notes.

Options:

If input and latest are not given, starts with an empty document.

If latest is given and input is not given, input defaults to current directory.

If input is a directory, latest is automatically enabled. The input file is then resolved using
latest-method.

If output is not given, it defaults to the current directory.

If output is a directory, the output file is "<output>/<current_date>.md".
		`,
		Version:       meta.Version,
		RunE:          run,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.Flags().StringP(
		"log", "L",
		"",
		`log file`,
	)
	cmd.Flags().StringP(
		"input", "i",
		viper.GetString("input_path"),
		`input file from which to load initial notes ("-" reads from stdin) or directory`,
	)
	cmd.Flags().StringP(
		"output", "o",
		viper.GetString("output_path"),
		`output file or directory`,
	)
	cmd.Flags().BoolP(
		"latest", "l",
		viper.GetBool("latest"),
		`use latest file as input file (in the same directory as input directory)`,
	)
	cmd.Flags().StringP(
		"latest-mode", "m",
		viper.GetString("latest_mode"),
		`mode to search for latest file (available: "name", "modified")`,
	)
	cmd.Flags().StringP(
		"slack", "s",
		viper.GetString("slack_channel"),
		`slack channel to publish notes`,
	)

	util.Must(viper.BindPFlag("log_path", cmd.Flags().Lookup("log")))
	util.Must(viper.BindPFlag("input_path", cmd.Flags().Lookup("input")))
	util.Must(viper.BindPFlag("output_path", cmd.Flags().Lookup("output")))
	util.Must(viper.BindPFlag("latest", cmd.Flags().Lookup("latest")))
	util.Must(viper.BindPFlag("latest_mode", cmd.Flags().Lookup("latest-mode")))
	util.Must(viper.BindPFlag("slack_channel", cmd.Flags().Lookup("slack")))

	return cmd
}

func run(_ *cobra.Command, _ []string) error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("failed to create configuration: %w", err)
	}

	log.Printf("rootCmd.run() -- config: %#+v", cfg)

	model := app.New(cfg)

	if _, err := tea.NewProgram(model).Run(); err != nil {
		return fmt.Errorf("failed to initialize ui: %w", err)
	}

	return nil
}

func Execute() {
	util.Must(rootCmd().Execute())
}
