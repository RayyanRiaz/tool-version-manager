package scriptdriventvm

import (
	"fmt"
	"path/filepath"
	"strings"

	"rayyanriaz/tool-version-manager/pkg/models"
	"rayyanriaz/tool-version-manager/pkg/utils"
)

type ScriptsDrivenTool struct {
	models.ToolBase `yaml:",inline"`
	Source          struct {
		Scripts struct {
			GetAllLocalVersions    []utils.ScriptStep `json:"getAllLocalVersions"`
			GetAllRemoteVersions   []utils.ScriptStep `json:"getAllRemoteVersions"`
			GetLatestRemoteVersion []utils.ScriptStep `json:"getLatestRemoteVersion"`
			FetchToolForVersion    []utils.ScriptStep `json:"fetchToolForVersion"`
			GetLinkInfo            []utils.ScriptStep `json:"getLinkInfo"`
			LinkTool               []utils.ScriptStep `json:"linkTool"`
			UnlinkTool             []utils.ScriptStep `json:"unlinkTool"`
		} `json:"scripts"`
	} `json:"source"`
	Extra map[string]interface{} `json:"extra,omitempty"`
}

func (t ScriptsDrivenTool) ShellFriendlySymlinks() string {
	b := strings.Builder{}
	for _, symlink := range t.Symlinks {
		sanitizedFrom := strings.TrimSpace(symlink.From)
		sanitizedTo := strings.TrimSpace(symlink.To)
		if sanitizedTo == "" {
			sanitizedTo = filepath.Base(sanitizedFrom)
		}
		fmt.Fprintf(&b, "%s:%s\n", sanitizedFrom, sanitizedTo)
	}

	// Remove the trailing newline if it exists
	s := b.String()
	if len(s) > 0 && s[len(s)-1] == '\n' {
		s = s[:len(s)-1]
	}
	return s
}
