package cmd

import (
	"fmt"

	"rayyanriaz/tool-version-manager/pkg/models"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install <tool-id> <version>",
	Short: "Install a specific version of a tool",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		toolID := args[0]
		version := models.ToolVersion(args[1])

		tool, tvm, err := getToolWithTVM(toolID)
		if err != nil {
			return err
		}

		fmt.Printf("Installing %s version %s...\n", toolID, version)

		err = tvm.InstallToolForVersion(tool, version)
		if err != nil {
			return fmt.Errorf("failed to install %s version %s: %w", toolID, version, err)
		}

		fmt.Printf("Successfully installed %s version %s\n", toolID, version)
		return nil
	},
}

var linkCmd = &cobra.Command{
	Use:   "link <tool-id> <version>",
	Short: "Link a tool version (make it active)",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		toolID := args[0]
		version := models.ToolVersion(args[1])

		tool, tvm, err := getToolWithTVM(toolID)
		if err != nil {
			return err
		}

		fmt.Printf("Linking %s version %s...\n", toolID, version)

		err = tvm.LinkTool(tool, version)
		if err != nil {
			return fmt.Errorf("failed to link %s version %s: %w", toolID, version, err)
		}

		fmt.Printf("Successfully linked %s version %s\n", toolID, version)
		return nil
	},
}

var unlinkCmd = &cobra.Command{
	Use:   "unlink <tool-id>",
	Short: "Unlink a tool (remove active version)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		toolID := args[0]

		tool, tvm, err := getToolWithTVM(toolID)
		if err != nil {
			return err
		}

		fmt.Printf("Unlinking %s...\n", toolID)

		err = tvm.UnlinkTool(tool)
		if err != nil {
			return fmt.Errorf("failed to unlink %s: %w", toolID, err)
		}

		fmt.Printf("Successfully unlinked %s\n", toolID)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(installCmd)
	RootCmd.AddCommand(linkCmd)
	RootCmd.AddCommand(unlinkCmd)
}
