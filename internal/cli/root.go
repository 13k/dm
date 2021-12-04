package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/13k/dm/internal/config"
	"github.com/13k/dm/internal/meta"
	"github.com/13k/dm/internal/ui/app"
	"github.com/13k/dm/internal/util"
)

var (
	cwd     string
	rawOpts rawOptions
)

type rawOptions struct {
	logPath    string
	basePath   string
	outputPath string
	latest     bool
	latestMode string
}

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
	logPath := rawOpts.logPath

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
		&rawOpts.logPath,
		"log", "L", "",
		`log file`,
	)
	cmd.Flags().StringVarP(
		&rawOpts.basePath,
		"base", "b", "",
		`base file from which to load initial notes ("-" reads from stdin)`,
	)
	cmd.Flags().BoolVarP(
		&rawOpts.latest,
		"latest", "l", false,
		`use latest (lexically by filename) notes file (in the same directory as '--output') as '--base' file`,
	)
	cmd.Flags().StringVarP(
		&rawOpts.latestMode,
		"latest-mode", "m", util.LatestFileByName.String(),
		`mode to search for '--latest' (available: "name", "modified")`,
	)
	cmd.Flags().StringVarP(
		&rawOpts.outputPath,
		"output", "o", defaultOutputPath(),
		`output file`,
	)

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	log.Printf("rootCmd.run() -- rawOpts: %#+v", rawOpts)

	opts, err := newOptions(&rawOpts)
	if err != nil {
		return fmt.Errorf("options error: %w", err)
	}

	log.Printf("rootCmd.run() -- opts: %#+v", opts)

	cfg := &config.Config{
		InputPath:  opts.inputPath,
		OutputPath: opts.outputPath,
	}

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
