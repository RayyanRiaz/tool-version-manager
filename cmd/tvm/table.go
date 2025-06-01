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
linked versions, installation dates, local versions, and optionally latest remote versions.`,
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
				row.LinkedAt = string(linkInfo.LinkedAt)
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

			// In case the user wants remote info
			if showRemote {
				if latestVersion, err := tvm.GetLatestRemoteVersion(tool); err == nil {
					row.LatestRemote = string(latestVersion)

					// Check if update is available
					if row.LinkedVersion != NA {
						if result, err := tvm.CompareVersions(tool, models.ToolVersion(row.LinkedVersion), latestVersion); err == nil {
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

func displayTable(rows []ToolTableRow) error {
	// Calculate column widths
	widths := []int{8, 6, 15, 12, 6} // min widths for: Tool, Type, Linked Version, Linked At, Local Count
	if showRemote {
		widths = append(widths, 14, 8) // Latest Remote, Update Available
	}
	widths = append(widths, 20) // Local Versions

	// Headers
	headers := []string{"Tool", "Type", "Linked Version", "Linked At", "Local Count"}
	if showRemote {
		headers = append(headers, "Latest Remote", "Update Available")
	}
	headers = append(headers, "Local Versions")

	// Update widths based on content
	for i, header := range headers {
		if len(header) > widths[i] {
			widths[i] = len(header)
		}
	}

	for _, row := range rows {
		data := []string{
			row.Name,
			row.Type,
			row.LinkedVersion,
			row.LinkedAt,
			fmt.Sprintf("%d", row.LocalCount),
		}

		if showRemote {
			updateStatus := "No"
			if row.UpdateAvailable {
				updateStatus = "Yes"
			}
			data = append(data, row.LatestRemote, updateStatus)
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
		data = append(data, localVersionsStr)

		// Update widths
		for i, col := range data {
			if i < len(widths) && len(col) > widths[i] {
				widths[i] = len(col)
			}
		}
	}

	// Print header
	printRow(headers, widths, true)

	// Print separator
	var sep []string
	for _, width := range widths {
		sep = append(sep, strings.Repeat("-", width))
	}
	printRow(sep, widths, false)

	// Print data rows
	for _, row := range rows {
		data := []string{
			row.Name,
			row.Type,
			row.LinkedVersion,
			row.LinkedAt,
			fmt.Sprintf("%d", row.LocalCount),
		}

		if showRemote {
			updateStatus := "No"
			if row.UpdateAvailable {
				updateStatus = "Yes"
			}
			data = append(data, row.LatestRemote, updateStatus)
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
		data = append(data, localVersionsStr)

		printRow(data, widths, false)
	}

	return nil
}

func printRow(cols []string, widths []int, isHeader bool) {
	for i, col := range cols {
		if i > 0 {
			fmt.Print("  ")
		}

		width := widths[i]
		if isHeader {
			fmt.Printf("%-*s", width, strings.ToUpper(col))
		} else {
			fmt.Printf("%-*s", width, col)
		}
	}
	fmt.Println()
}

func init() {
	tableCmd.Flags().BoolVarP(&showRemote, "remote", "r", false, "Fetch and show latest remote versions")
	tableCmd.Flags().StringVarP(&sortBy, "sort", "s", "name", "Sort by: name, type, linked, count")
	tableCmd.Flags().StringVarP(&outputFormat, "format", "f", "table", "Output format: table, json")

	RootCmd.AddCommand(tableCmd)
}
