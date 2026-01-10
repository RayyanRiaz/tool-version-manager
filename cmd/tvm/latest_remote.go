package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var latestCmd = &cobra.Command{
	Use:   "latest <tool-id>",
	Short: "Show the latest available version of a tool",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		toolID := args[0]

		tool, tvm, err := getToolWithTVM(toolID)
		if err != nil {
			return err
		}

		version, err := tvm.GetLatestRemoteVersion(tool)
		if err != nil {
			return fmt.Errorf("failed to get latest version for %s: %w", toolID, err)
		}

		// Update cache with latest version
		_ = updateCachedLatestVersion(toolID, version)

		fmt.Println(version)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(latestCmd)
}
