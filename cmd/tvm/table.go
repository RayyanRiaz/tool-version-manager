package cmd

import (
	"fmt"
	"sort"
	"strings"

	"rayyanriaz/tool-version-manager/pkg/models"

	"github.com/spf13/cobra"
)

type ToolTableRow struct {
	Name            string
	Type            string
	LinkedVersion   string
	LinkedAt        string
	LocalVersions   []string
	LocalCount      int
	LatestRemote    string
	UpdateAvailable bool
}

var (
	showRemote   bool
	sortBy       string
	outputFormat string
)

const NA = "NA"

var tableCmd = &cobra.Command{
	Use:   "table",
	Short: "Show a beautiful tabular view of all tools",
	Long: `Display all configured tools in a beautiful table format with information about
linked versions, installation dates, local versions, and latest remote versions (from cache).
Use --remote to fetch fresh latest versions from remote sources and update the cache.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		tools, err := getAllTools()
		if err != nil {
			return fmt.Errorf("failed to get tools: %w", err)
		}

		// Collect data for all tools
		var rows []ToolTableRow

		for _, toolWrapper := range tools {
			tool := toolWrapper.Wrapped
			tvm, err := models.ToolRegistrar.GetTVM(tool.GetType())
			if err != nil {
				return fmt.Errorf("failed to get TVM for tool %s: %w", tool.GetId(), err)
			}
			row := ToolTableRow{
				Name: tool.GetId(),
				Type: tool.GetType(),
			}

			// Get linked version and date
			if linkInfo, err := tvm.GetLinkInfo(tool); err == nil && linkInfo != nil && linkInfo.Version != "" {
				row.LinkedVersion = string(linkInfo.Version)
				row.LinkedAt = formatLinkedAt(string(linkInfo.LinkedAt))
			} else {
				row.LinkedVersion = NA
				row.LinkedAt = NA
			}

			// Get local versions
			if localVersions, err := tvm.GetAllLocalVersions(tool); err == nil {
				row.LocalCount = len(localVersions)
				var versionStrs []string
				for _, v := range localVersions {
					versionStrs = append(versionStrs, string(v))
				}
				row.LocalVersions = versionStrs
			} else {
				row.LocalCount = 0
			}

			// Get latest remote version - from cache or fresh fetch
			if showRemote {
				// Fresh fetch from remote and update cache
				if latestVersion, err := tvm.GetLatestRemoteVersion(tool); err == nil {
					row.LatestRemote = string(latestVersion)
					// Update cache with fresh value
					_ = updateCachedLatestVersion(tool.GetId(), latestVersion)

					// Check if update is available
					if row.LinkedVersion != NA {
						if result, err := tvm.CompareVersions(tool, models.ToolVersion(row.LinkedVersion), latestVersion); err == nil {
							row.UpdateAvailable = result < 0
						}
					}
				} else {
					row.LatestRemote = NA
				}
			} else {
				// Use cached version if available
				if cachedVersion, found := getCachedLatestVersion(tool.GetId()); found {
					row.LatestRemote = string(cachedVersion)

					// Check if update is available
					if row.LinkedVersion != NA {
						if result, err := tvm.CompareVersions(tool, models.ToolVersion(row.LinkedVersion), cachedVersion); err == nil {
							row.UpdateAvailable = result < 0
						}
					}
				} else {
					row.LatestRemote = NA
				}
			}

			rows = append(rows, row)
		}

		// Sort rows
		switch sortBy {
		case "name":
			sort.Slice(rows, func(i, j int) bool {
				return rows[i].Name < rows[j].Name
			})
		case "type":
			sort.Slice(rows, func(i, j int) bool {
				if rows[i].Type == rows[j].Type {
					return rows[i].Name < rows[j].Name
				}
				return rows[i].Type < rows[j].Type
			})
		case "linked":
			sort.Slice(rows, func(i, j int) bool {
				if rows[i].LinkedVersion == rows[j].LinkedVersion {
					return rows[i].Name < rows[j].Name
				}
				return rows[i].LinkedVersion < rows[j].LinkedVersion
			})
		case "count":
			sort.Slice(rows, func(i, j int) bool {
				if rows[i].LocalCount == rows[j].LocalCount {
					return rows[i].Name < rows[j].Name
				}
				return rows[i].LocalCount > rows[j].LocalCount
			})
		}

		return displayTable(rows)
	},
}

// ANSI color codes
const (
	colorReset = "\033[0m"
	colorGreen = "\033[32m"
)

// formatLinkedAt truncates the timestamp to show only up to seconds (YYYY-MM-DD HH:MM:SS)
func formatLinkedAt(linkedAt string) string {
	if linkedAt == "" || linkedAt == NA {
		return linkedAt
	}
	if len(linkedAt) >= 19 {
		return linkedAt[:19]
	}
	return linkedAt
}

func displayTable(rows []ToolTableRow) error {
	// Build headers and widths based on verbose mode
	// Normal mode: Tool, Linked Version, Latest Remote, Update, Linked At, Local Count, Local Versions
	// Verbose mode: Tool, Type, Linked Version, Latest Remote, Update, Linked At, Local Count, Local Versions
	var headers []string
	var widths []int

	if verbose {
		headers = []string{"Tool", "Type", "Linked Version", "Latest Remote", "Update", "Linked At", "Local Count", "Local Versions"}
		widths = []int{8, 6, 15, 14, 6, 12, 6, 20}
	} else {
		headers = []string{"Tool", "Linked Version", "Latest Remote", "Update", "Linked At", "Local Count", "Local Versions"}
		widths = []int{8, 15, 14, 6, 12, 6, 20}
	}

	// Update widths based on content
	for i, header := range headers {
		if len(header) > widths[i] {
			widths[i] = len(header)
		}
	}

	for _, row := range rows {
		updateStatus := "No"
		if row.UpdateAvailable {
			updateStatus = "Yes"
		}

		// Format local versions
		localVersionsStr := "-"
		if len(row.LocalVersions) > 0 {
			if len(row.LocalVersions) <= 3 {
				localVersionsStr = strings.Join(row.LocalVersions, ", ")
			} else {
				localVersionsStr = strings.Join(row.LocalVersions[:3], ", ") + fmt.Sprintf(" (+%d more)", len(row.LocalVersions)-3)
			}
		}

		var data []string
		if verbose {
			data = []string{
				row.Name,
				row.Type,
				row.LinkedVersion,
				row.LatestRemote,
				updateStatus,
				row.LinkedAt,
				fmt.Sprintf("%d", row.LocalCount),
				localVersionsStr,
			}
		} else {
			data = []string{
				row.Name,
				row.LinkedVersion,
				row.LatestRemote,
				updateStatus,
				row.LinkedAt,
				fmt.Sprintf("%d", row.LocalCount),
				localVersionsStr,
			}
		}

		// Update widths
		for i, col := range data {
			if i < len(widths) && len(col) > widths[i] {
				widths[i] = len(col)
			}
		}
	}

	// Print header
	printTableRow(headers, widths, true, nil)

	// Print separator
	var sep []string
	for _, width := range widths {
		sep = append(sep, strings.Repeat("-", width))
	}
	printTableRow(sep, widths, false, nil)

	// Print data rows
	for _, row := range rows {
		updateStatus := "No"
		if row.UpdateAvailable {
			updateStatus = "Yes"
		}

		// Format local versions
		localVersionsStr := "-"
		if len(row.LocalVersions) > 0 {
			if len(row.LocalVersions) <= 3 {
				localVersionsStr = strings.Join(row.LocalVersions, ", ")
			} else {
				localVersionsStr = strings.Join(row.LocalVersions[:3], ", ") + fmt.Sprintf(" (+%d more)", len(row.LocalVersions)-3)
			}
		}

		var data []string
		var updateColIndex int
		if verbose {
			data = []string{
				row.Name,
				row.Type,
				row.LinkedVersion,
				row.LatestRemote,
				updateStatus,
				row.LinkedAt,
				fmt.Sprintf("%d", row.LocalCount),
				localVersionsStr,
			}
			updateColIndex = 4
		} else {
			data = []string{
				row.Name,
				row.LinkedVersion,
				row.LatestRemote,
				updateStatus,
				row.LinkedAt,
				fmt.Sprintf("%d", row.LocalCount),
				localVersionsStr,
			}
			updateColIndex = 3
		}

		// Set color for update column if update is available
		var colColors map[int]string
		if row.UpdateAvailable {
			colColors = map[int]string{updateColIndex: colorGreen}
		}

		printTableRow(data, widths, false, colColors)
	}

	return nil
}

func printTableRow(cols []string, widths []int, isHeader bool, colColors map[int]string) {
	for i, col := range cols {
		if i > 0 {
			fmt.Print("  ")
		}

		width := widths[i]

		// Check if this column should be colored
		color := ""
		reset := ""
		if colColors != nil {
			if c, ok := colColors[i]; ok {
				color = c
				reset = colorReset
			}
		}

		if isHeader {
			fmt.Printf("%-*s", width, strings.ToUpper(col))
		} else {
			fmt.Printf("%s%-*s%s", color, width, col, reset)
		}
	}
	fmt.Println()
}

func init() {
	tableCmd.Flags().BoolVarP(&showRemote, "remote", "r", false, "Fetch fresh latest versions from remote (updates cache)")
	tableCmd.Flags().StringVarP(&sortBy, "sort", "s", "name", "Sort by: name, type, linked, count")
	tableCmd.Flags().StringVarP(&outputFormat, "format", "f", "table", "Output format: table, json")

	RootCmd.AddCommand(tableCmd)
}
