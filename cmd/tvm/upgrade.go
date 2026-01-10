package cmd

import (
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"rayyanriaz/tool-version-manager/pkg/models"

	"github.com/spf13/cobra"
)

var force bool
var all bool

var upgradeCmd = &cobra.Command{
	Use:   "upgrade <tool-id>",
	Short: "Upgrade a tool to the latest version",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !all && len(args) == 0 {
			return fmt.Errorf("you must provide either a tool ID or use the --all flag")
		}

		if all && len(args) > 0 {
			return fmt.Errorf("cannot use --all with specific tool IDs")
		}

		allTools, err := getAllTools()
		if err != nil {
			return fmt.Errorf("failed to get all tools: %w", err)
		}
		if all {
			// Get all tools
			toolIDs := make([]string, len(allTools))
			for i, tool := range allTools {
				toolIDs[i] = tool.Wrapped.GetId()
			}
			slog.Debug("Upgrading", "tools", toolIDs)
			return upgradeTools(toolIDs)
		}

		toolIDs := strings.Split(args[0], ",")
		for i := range toolIDs {
			toolIDs[i] = strings.TrimSpace(toolIDs[i])
			if toolIDs[i] == "" {
				return fmt.Errorf("invalid empty tool ID in input")
			}
			// toolId should exist in the tools list
			if _, err := getToolById(toolIDs[i]); err != nil {
				return fmt.Errorf("tool %s does not exist: %w", toolIDs[i], err)
			}
		}

		return upgradeTools(toolIDs)
	},
}

func upgradeTools(toolIDs []string) error {

	var wg sync.WaitGroup
	errs := make([]error, len(toolIDs))
	var mu sync.Mutex
	for i, toolID := range toolIDs {
		wg.Add(1)
		go func(i int, toolID string) {
			defer wg.Done()
			err := upgradeTool(toolID)
			mu.Lock()
			errs[i] = err
			if err == nil {
				fmt.Printf("Successfully upgraded %s\n", toolID)
			} else {
				fmt.Printf("Failed to upgrade %s: %v\n", toolID, err)
			}
			mu.Unlock()
		}(i, toolID)
	}

	wg.Wait()

	var finalErr error
	for _, err := range errs {
		if err != nil {
			if finalErr == nil {
				finalErr = fmt.Errorf("errors occurred during upgrade: %v", err)
			} else {
				finalErr = fmt.Errorf("%w; %v", finalErr, err)
			}
		}
	}
	if finalErr != nil {
		slog.Error("Upgrade errors", "errors", finalErr)
		return finalErr
	} else {
		slog.Info("All tools upgraded successfully")
		fmt.Println("All specified tools upgraded successfully.")
		return nil
	}

}

func upgradeTool(toolID string) error {
	tool, tvm, err := getToolWithTVM(toolID)
	if err != nil {
		return fmt.Errorf("failed to get tool %s: %w", toolID, err)
	}

	// Get latest version
	latestVersion, err := tvm.GetLatestRemoteVersion(tool)
	if err != nil {
		return fmt.Errorf("failed to get latest version for %s: %w", toolID, err)
	}

	// Update cache with latest version
	_ = updateCachedLatestVersion(toolID, latestVersion)

	// Get current linked version
	linkInfo, err := tvm.GetLinkInfo(tool)
	var currentVersion models.ToolVersion
	if err != nil {
		// If linkInfo is nil or Version is empty, treat as no current version
		if linkInfo == nil || linkInfo.Version == "" {
			currentVersion = ""
			// No current version, proceed with installation
		} else {
			return fmt.Errorf("failed to get current version for %s: %w", toolID, err)
		}
		// No current version, proceed with installation
	} else {
		currentVersion = linkInfo.Version
	}

	// Compare versions if there's a current version
	if currentVersion != "" {
		result, err := tvm.CompareVersions(tool, currentVersion, latestVersion)
		if err != nil {
			return fmt.Errorf("failed to compare versions for %s: %w", toolID, err)
		}

		if result >= 0 && !force {
			fmt.Printf("%s is already at the latest version (%s)\n", toolID, currentVersion)
			return nil
		}
	}

	fmt.Printf("Upgrading %s to version %s...\n", toolID, latestVersion)

	// check if the tool is already installed
	localVersions, err := tvm.GetAllLocalVersions(tool)
	if err != nil {
		slog.Debug("Warning: failed to get local versions for tool %s: %v", toolID, err)
	}
	isInstalled := false
	for _, v := range localVersions {
		if v == latestVersion {
			isInstalled = true
			break
		}
	}
	if isInstalled {
		slog.Debug("Tool %s version %s is already installed", toolID, latestVersion)
	} else {
		err = tvm.InstallToolForVersion(tool, latestVersion)
		if err != nil {
			return fmt.Errorf("failed to install %s version %s: %w", toolID, latestVersion, err)
		}
	}

	// Link latest version
	err = tvm.LinkTool(tool, latestVersion)
	if err != nil {
		return fmt.Errorf("failed to link %s version %s: %w", toolID, latestVersion, err)
	}

	fmt.Printf("Successfully upgraded %s to version %s\n", toolID, latestVersion)
	return nil

}

func init() {
	upgradeCmd.Flags().BoolVarP(&force, "force", "f", false, "Force link even if another version is already linked")
	upgradeCmd.Flags().BoolVarP(&all, "all", "a", false, "Upgrade all tools to their latest versions")

	RootCmd.AddCommand(upgradeCmd)
}
