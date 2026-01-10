package scriptdriventvm

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"rayyanriaz/tool-version-manager/pkg/impl/config"
	"rayyanriaz/tool-version-manager/pkg/models"
	"rayyanriaz/tool-version-manager/pkg/utils"
)

type ScriptsDrivenTVM struct {
	configService *config.LocalFileConfig
}

func NewScriptsDrivenTVM() *ScriptsDrivenTVM {
	slog.Debug("Creating new ScriptsDrivenTVM")
	configService, err := models.ToolRegistrar.GetConfig("scripts_driven")
	if err != nil {
		slog.Error("Failed to get config service for ScriptsDrivenTVM", "error", err)
		return nil
	}
	return &ScriptsDrivenTVM{
		configService: configService.(*config.LocalFileConfig),
	}
}

func (t *ScriptsDrivenTVM) CreateNewTool() models.Tool {
	slog.Debug("Creating new ScriptsDrivenTool")
	return &ScriptsDrivenTool{}
}

func (t *ScriptsDrivenTVM) buildTemplateVars(tool models.Tool, argToFirstStep string) map[string]any {
	vars := map[string]any{
		"Config": map[string]any{
			"DownloadsDir": t.configService.DownloadsDir,
			"SymlinksDir":  t.configService.SymlinksDir,
			"GitHubToken":  t.configService.GitHubToken,
		},
		"Tool": tool,
		"Arg":  argToFirstStep,
	}
	return vars
}

func (t *ScriptsDrivenTVM) GetLinkInfo(tool models.Tool) (*models.ToolLinkInfo, error) {

	script := tool.(*ScriptsDrivenTool).Source.Scripts.GetLinkInfo
	vars := t.buildTemplateVars(tool, "")
	out, err := utils.ExecuteBashScriptSteps(script, vars)
	if err != nil {
		return nil, fmt.Errorf("failed to get link info for tool %s: %w", tool.GetId(), err)
	}

	var linkInfo models.ToolLinkInfo
	if err := json.Unmarshal([]byte(out), &linkInfo); err != nil {
		return nil, fmt.Errorf("failed to parse link info for tool %s: %w", tool.GetId(), err)
	}
	return &linkInfo, nil
}

func (t *ScriptsDrivenTVM) GetAllLocalVersions(tool models.Tool) ([]models.ToolVersion, error) {
	script := tool.(*ScriptsDrivenTool).Source.Scripts.GetAllLocalVersions
	vars := t.buildTemplateVars(tool, "")

	out, err := utils.ExecuteBashScriptSteps(script, vars)
	if err != nil {
		return nil, fmt.Errorf("failed to get all local versions for tool %s: %w", tool.GetId(), err)
	}
	var vs []models.ToolVersion
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		vs = append(vs, models.ToolVersion(line))
	}
	return vs, nil
}

func (t *ScriptsDrivenTVM) GetAllRemoteVersions(tool models.Tool) ([]models.ToolVersion, error) {
	script := tool.(*ScriptsDrivenTool).Source.Scripts.GetAllRemoteVersions
	vars := t.buildTemplateVars(tool, "")
	out, err := utils.ExecuteBashScriptSteps(script, vars)
	if err != nil {
		return nil, fmt.Errorf("failed to get all remote versions for tool %s: %w", tool.GetId(), err)
	}
	var vs []models.ToolVersion
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		vs = append(vs, models.ToolVersion(line))
	}
	return vs, nil
}

func (t *ScriptsDrivenTVM) GetLatestRemoteVersion(tool models.Tool) (models.ToolVersion, error) {
	script := tool.(*ScriptsDrivenTool).Source.Scripts.GetLatestRemoteVersion
	vars := t.buildTemplateVars(tool, "")
	out, err := utils.ExecuteBashScriptSteps(script, vars)
	return models.ToolVersion(strings.TrimSpace(out)), err
}

func (t *ScriptsDrivenTVM) CompareVersions(tool models.Tool, v1, v2 models.ToolVersion) (int, error) {
	comparator := &models.ToolComparerWithVersionParsing{}
	res, err := comparator.CompareVersions(tool, v1, v2)
	return res, err
}

func (t *ScriptsDrivenTVM) InstallToolForVersion(tool models.Tool, version models.ToolVersion) error {
	script := tool.(*ScriptsDrivenTool).Source.Scripts.FetchToolForVersion
	vars := t.buildTemplateVars(tool, string(version))
	out, err := utils.ExecuteBashScriptSteps(script, vars)
	if err != nil {
		return fmt.Errorf("failed to install tool %s for version %s: %w", tool.GetId(), version, err)
	}
	if out != "" {
		return fmt.Errorf("script output: %s", out)
	}
	return nil
}

func (t *ScriptsDrivenTVM) LinkTool(tool models.Tool, version models.ToolVersion) error {
	if version == "" {
		return fmt.Errorf("version cannot be empty")
	}

	toolInfo, err := t.GetLinkInfo(tool)
	if err != nil { // log the error but continue
		return fmt.Errorf("failed to get linked version for tool %s: %w", tool.GetId(), err)
	}
	currentVersion := toolInfo.Version

	linkScript := tool.(*ScriptsDrivenTool).Source.Scripts.LinkTool
	vars := t.buildTemplateVars(tool, string(version))

	var successfullyLinked = false
	defer func() {
		if !successfullyLinked {
			if currentVersion != "" {
				// revert to the previous version
				_, revertErr := utils.ExecuteBashScriptSteps(linkScript, vars)
				if revertErr != nil {
					fmt.Printf("failed to revert to previous version %s: %v\n", currentVersion, revertErr)
				}
			} else {
				// if there was no previous version, just unlink the tool
				if unlinkErr := t.UnlinkTool(tool); unlinkErr != nil {
					fmt.Printf("failed to unlink tool %s: %v\n", tool.GetId(), unlinkErr)
				}
			}
		}
	}()
	_, err = utils.ExecuteBashScriptSteps(linkScript, vars)
	if err != nil {
		return fmt.Errorf("failed to link tool %s to version %s: %w", tool.GetId(), version, err)
	}
	successfullyLinked = true
	return nil
}

func (t *ScriptsDrivenTVM) UnlinkTool(tool models.Tool) error {
	var successfullyUnlinked = false

	toolInfo, err := t.GetLinkInfo(tool)
	if err != nil { // log the error but continue
		fmt.Printf("failed to get linked version for tool %s: %v\n", tool.GetId(), err)
	}
	if err != nil {
		return fmt.Errorf("failed to get linked version: %w", err)
	}
	currentVersion := toolInfo.Version
	if currentVersion == "" {
		return fmt.Errorf("tool %s is not linked to any version", tool.GetId())
	}

	defer func() {
		if !successfullyUnlinked {
			// if unlinking fails, try to revert to the previous version
			if revertErr := t.LinkTool(tool, currentVersion); revertErr != nil {
				fmt.Printf("failed to revert to previous version %s: %v\n", currentVersion, revertErr)
			}
		}
	}()

	unlinkScript := tool.(*ScriptsDrivenTool).Source.Scripts.UnlinkTool
	vars := t.buildTemplateVars(tool, string(currentVersion))

	_, err = utils.ExecuteBashScriptSteps(unlinkScript, vars)
	if err != nil {
		return fmt.Errorf("failed to unlink tool %s: %w", tool.GetId(), err)
	}
	successfullyUnlinked = true

	return nil
}

var _ models.ToolVersionManager = (*ScriptsDrivenTVM)(nil)
