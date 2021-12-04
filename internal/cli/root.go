package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/13k/dm/internal/meta"
	"github.com/13k/dm/internal/ui/app"
	"github.com/13k/dm/internal/util"
)

var (
	cwd  string
	opts options
)

func init() {
	var err error

	cwd, err = os.Getwd()

	if err != nil {
		util.Fatal("could not determine current working directory: %w", err)
	}

	cobra.OnInitialize(onInit)
}

func onInit() {
	configureLogger()
}

func configureLogger() {
	logPath := opts.logPath

	if logPath == "" {
		logPath = os.DevNull
	}

	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o666)
	if err != nil {
		util.Fatal("could not open log file %q: %w", logPath, err)
	}

	log.SetOutput(f)
}

func defaultOutputPath() string {
	basename := fmt.Sprintf("%s.md", util.Today())

	return filepath.Join(cwd, basename)
}

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "dm [flags]",
		Short:         "Create daily meeting notes",
		Version:       meta.Version,
		RunE:          run,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.Flags().StringVarP(
		&opts.logPath,
		"log", "L", "",
		`log file`,
	)
	cmd.Flags().StringVarP(
		&opts.basePath,
		"base", "b", "",
		`base file from which to load initial notes ("-" reads from stdin)`,
	)
	cmd.Flags().BoolVarP(
		&opts.latest,
		"latest", "l", false,
		`use latest (lexically by filename) notes file (in the same directory as '--output') as '--base' file`,
	)
	cmd.Flags().StringVarP(
		&opts.latestMode,
		"latest-mode", "m", util.LatestFileByName.String(),
		`mode to search for '--latest' (available: "name", "modified")`,
	)
	cmd.Flags().StringVarP(
		&opts.outputPath,
		"output", "o", defaultOutputPath(),
		`output file`,
	)
	cmd.Flags().StringVarP(
		&opts.slackChannel,
		"slack", "s", "",
		`slack channel to publish notes`,
	)

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	log.Printf("rootCmd.run() -- opts: %#+v", opts)

	cfg, err := parseOptions(&opts)
	if err != nil {
		return fmt.Errorf("options error: %w", err)
	}

	log.Printf("rootCmd.run() -- opts: %#+v", opts)

	model := app.NewModel(cfg)

	if err := tea.NewProgram(model).Start(); err != nil {
		return fmt.Errorf("could not initialize ui: %w", err)
	}

	return nil
}

func Execute() {
	if err := rootCmd().Execute(); err != nil {
		util.Fatal("Error: %v", err)
	}
}
