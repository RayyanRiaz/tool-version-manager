package cmd

import (
	"fmt"
	"strings"
	"sync"

	"rayyanriaz/tool-version-manager/pkg/models"

	"github.com/spf13/cobra"
)

var fetchAll bool

var fetchCmd = &cobra.Command{
	Use:   "fetch [tool-id]",
	Short: "Fetch and cache the latest remote versions for tools",
	Long: `Fetch the latest available versions from remote sources and cache them locally.
This allows other commands like 'table' to show remote version info without making network requests.

Examples:
  tvm fetch --all       # Fetch latest versions for all tools
  tvm fetch ripgrep     # Fetch latest version for a specific tool
  tvm fetch rg,fzf,fd   # Fetch latest versions for multiple tools`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !fetchAll && len(args) == 0 {
			return fmt.Errorf("you must provide either a tool ID or use the --all flag")
		}

		if fetchAll && len(args) > 0 {
			return fmt.Errorf("cannot use --all with specific tool IDs")
		}

		allTools, err := getAllTools()
		if err != nil {
			return fmt.Errorf("failed to get all tools: %w", err)
		}

		var toolIDs []string
		if fetchAll {
			toolIDs = make([]string, len(allTools))
			for i, tool := range allTools {
				toolIDs[i] = tool.Wrapped.GetId()
			}
		} else {
			toolIDs = strings.Split(args[0], ",")
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
		}

		return fetchLatestVersions(toolIDs)
	},
}

func fetchLatestVersions(toolIDs []string) error {
	var wg sync.WaitGroup
	results := make([]struct {
		toolID  string
		version models.ToolVersion
		err     error
	}, len(toolIDs))

	for i, toolID := range toolIDs {
		wg.Add(1)
		go func(i int, toolID string) {
			defer wg.Done()
			tool, tvm, err := getToolWithTVM(toolID)
			if err != nil {
				results[i] = struct {
					toolID  string
					version models.ToolVersion
					err     error
				}{toolID, "", err}
				return
			}

			version, err := tvm.GetLatestRemoteVersion(tool)
			results[i] = struct {
				toolID  string
				version models.ToolVersion
				err     error
			}{toolID, version, err}
		}(i, toolID)
	}

	wg.Wait()

	// Process results and update cache
	var hasErrors bool
	for _, result := range results {
		if result.err != nil {
			fmt.Printf("Failed to fetch %s: %v\n", result.toolID, result.err)
			hasErrors = true
			continue
		}

		if err := updateCachedLatestVersion(result.toolID, result.version); err != nil {
			fmt.Printf("Failed to cache %s: %v\n", result.toolID, err)
			hasErrors = true
			continue
		}

		fmt.Printf("%s: %s\n", result.toolID, result.version)
	}

	if hasErrors {
		return fmt.Errorf("some tools failed to fetch")
	}

	fmt.Printf("\nSuccessfully fetched and cached latest versions for %d tools\n", len(toolIDs))
	return nil
}

func init() {
	fetchCmd.Flags().BoolVarP(&fetchAll, "all", "a", false, "Fetch latest versions for all tools")
	RootCmd.AddCommand(fetchCmd)
}
