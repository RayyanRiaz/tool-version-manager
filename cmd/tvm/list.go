package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List commands for tools and versions",
}

var listLocalCmd = &cobra.Command{
	Use:   "local <tool-id>",
	Short: "List all locally installed versions of a tool",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		toolID := args[0]

		tool, tvm, err := getToolWithTVM(toolID)
		if err != nil {
			return err
		}

		versions, err := tvm.GetAllLocalVersions(tool)
		if err != nil {
			return fmt.Errorf("failed to get local versions for %s: %w", toolID, err)
		}

		if len(versions) == 0 {
			fmt.Printf("No local versions found for %s\n", toolID)
			return nil
		}

		fmt.Printf("Local versions for %s:\n", toolID)
		for _, version := range versions {
			fmt.Printf("  %s\n", version)
		}

		return nil
	},
}

var listRemoteCmd = &cobra.Command{
	Use:   "remote <tool-id>",
	Short: "List all remote versions of a tool",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		toolID := args[0]

		tool, tvm, err := getToolWithTVM(toolID)
		if err != nil {
			return err
		}

		versions, err := tvm.GetAllRemoteVersions(tool)
		if err != nil {
			return fmt.Errorf("failed to get remote versions for %s: %w", toolID, err)
		}

		if len(versions) == 0 {
			fmt.Printf("No remote versions found for %s\n", toolID)
			return nil
		}

		fmt.Printf("Remote versions for %s:\n", toolID)
		for _, version := range versions {
			fmt.Printf("  %s\n", version)
		}

		return nil
	},
}

func init() {
	listCmd.AddCommand(listLocalCmd)
	listCmd.AddCommand(listRemoteCmd)

	RootCmd.AddCommand(listCmd)
}
