package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var currentCmd = &cobra.Command{
	Use:   "current <tool-id>",
	Short: "Show the currently linked version of a tool",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		toolID := args[0]

		tool, tvm, err := getToolWithTVM(toolID)
		if err != nil {
			return err
		}

		linkInfo, err := tvm.GetLinkInfo(tool)
		if err != nil {
			return fmt.Errorf("failed to get linked version for %s: %w", toolID, err)
		}

		if linkInfo.Version == "" {
			fmt.Printf("No version linked for %s\n", toolID)
		} else {
			fmt.Printf("Current version of %s: %s\n. Linked at: %s", toolID, linkInfo.Version, linkInfo.LinkedAt)
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(currentCmd)
}
