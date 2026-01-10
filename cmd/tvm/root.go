package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "tvm",
	Short: "Tool Version Manager - Manage versions of development tools",
	Long: `TVM (Tool Version Manager) helps you manage multiple versions of development tools.
You can install, link, unlink, and switch between different versions of tools like ripgrep, fd, and more.`,
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "Path to configuration file (default: $TVM_CONFIG or tools.yaml)")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	// Set up logging and bootstrap after flags are parsed
	RootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if verbose {
			// Set up a logger that outputs to stdout as well as a debug file
			logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}))
			// Optionally, you can also log to a file
			file, err := os.OpenFile("tvm.debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err == nil {
				fileHandler := slog.NewTextHandler(file, &slog.HandlerOptions{
					Level: slog.LevelDebug,
				})
				logger = slog.New(fileHandler)
			} else {
				slog.Warn("Failed to open debug log file, logging to stdout only", "error", err)
			}
			slog.SetDefault(logger)
		}

		// Bootstrap after flags are parsed so configPath is available
		return bootstrap()
	}
}
