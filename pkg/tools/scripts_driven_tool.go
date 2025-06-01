package tools

// import (
// 	"fmt"
// 	"log/slog"
// 	"strings"

// 	"rayyanriaz/tool-version-manager/pkg/config"
// 	"rayyanriaz/tool-version-manager/pkg/models"
// 	"rayyanriaz/tool-version-manager/pkg/state"
// 	"rayyanriaz/tool-version-manager/pkg/utils"
// )

// type ScriptsDrivenTool struct {
// 	models.ToolBase `yaml:",inline"`
// 	Source          struct {
// 		Scripts map[string][]utils.ScriptStep `json:"scripts"`
// 	} `json:"source"`
// 	Extra map[string]interface{} `json:"extra,omitempty"`
// }

// type ScriptsDrivenTVM struct {
// 	StateService  state.ToolStateService
// 	ConfigService config.ConfigService
// }

// func (t *ScriptsDrivenTVM) GetToolFactory() (models.ToolFactory, error) {
// 	slog.Debug("Creating ScriptsDrivenTool factory")
// 	return func() models.Tool {
// 		return &ScriptsDrivenTool{}
// 	}, nil
// }

// func (t *ScriptsDrivenTVM) ExecuteScriptSteps(tool models.Tool, name string, argToFirstStep string) (string, error) {
// 	slog.Debug("Executing steps for tool", "tool", tool.GetId(), "name", name, "argToFirstStep", argToFirstStep)
// 	sdTool, ok := tool.(*ScriptsDrivenTool)
// 	if !ok {
// 		return "", fmt.Errorf("invalid tool type: expected *ScriptDrivenTool, got %T", tool)
// 	}

// 	steps, ok := sdTool.Source.Scripts[name]
// 	if !ok {
// 		return "", fmt.Errorf("script %q not defined for tool %s", name, sdTool.GetId())
// 	}

// 	vars := map[string]any{
// 		"Config": map[string]any{
// 			"DownloadsDir": t.ConfigService.GetConfigFile().DownloadsDir,
// 			"SymlinksDir":  t.ConfigService.GetConfigFile().SymlinksDir,
// 			"StateFile":    t.ConfigService.GetConfigFile().StateFilePath,
// 		},
// 		"Tool":        sdTool,
// 		"Arg":         argToFirstStep,
// 		"StepOutputs": make(map[string]string),
// 	}

// 	out, err := utils.ExecuteBashScriptSteps(steps, vars)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to execute script steps for tool %s: %w", tool.GetId(), err)
// 	}
// 	slog.Debug("Executed script steps", "tool", tool.GetId(), "name", name, "output", out)

// 	return out, nil

// }

// func (t *ScriptsDrivenTVM) GetLinkedVersion(tool models.Tool) (models.ToolVersion, error) {
// 	// todo: proper error types. because sometimes we want to use "" tool version to indicate "not linked"
// 	toolState, err := t.StateService.GetToolState(tool.GetId())
// 	if err != nil {
// 		return "", fmt.Errorf("failed to get tool state: %w", err)
// 	}
// 	if toolState == nil {
// 		return "", fmt.Errorf("tool %s not found in state file", tool.GetId())
// 	}
// 	if toolState.ActiveLocalVersion == "" {
// 		return "", fmt.Errorf("tool %s is not linked", tool.GetId())
// 	}
// 	return toolState.ActiveLocalVersion, nil
// }

// func (t *ScriptsDrivenTVM) GetAllLocalVersions(tool models.Tool) ([]models.ToolVersion, error) {
// 	out, err := t.ExecuteScriptSteps(tool, "getAllLocalVersions", "")
// 	if err != nil {
// 		return nil, err
// 	}
// 	var vs []models.ToolVersion
// 	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
// 		vs = append(vs, models.ToolVersion(line))
// 	}
// 	return vs, nil
// }

// func (t *ScriptsDrivenTVM) GetAllRemoteVersions(tool models.Tool) ([]models.ToolVersion, error) {
// 	out, err := t.ExecuteScriptSteps(tool, "getAllRemoteVersions", "")
// 	if err != nil {
// 		return nil, err
// 	}
// 	var vs []models.ToolVersion
// 	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
// 		vs = append(vs, models.ToolVersion(line))
// 	}
// 	return vs, nil
// }

// func (t *ScriptsDrivenTVM) GetLatestRemoteVersion(tool models.Tool) (models.ToolVersion, error) {
// 	out, err := t.ExecuteScriptSteps(tool, "getLatestRemoteVersion", "")
// 	return models.ToolVersion(strings.TrimSpace(out)), err
// }

// func (t *ScriptsDrivenTVM) CompareVersions(tool models.Tool, v1, v2 models.ToolVersion) (int, error) {
// 	comparator := &models.ToolComparerWithVersionParsing{}
// 	res, err := comparator.CompareVersions(tool, v1, v2)
// 	return res, err
// }

// func (t *ScriptsDrivenTVM) InstallToolForVersion(tool models.Tool, version models.ToolVersion) error {
// 	_, err := t.ExecuteScriptSteps(tool, "fetchToolForVersion", string(version))
// 	return err
// }

// func (t *ScriptsDrivenTVM) LinkTool(tool models.Tool, version models.ToolVersion, force bool) error {
// 	if version == "" {
// 		return fmt.Errorf("version cannot be empty")
// 	}

// 	if tool == nil || tool.GetId() == "" {
// 		return fmt.Errorf("tool cannot be nil or have an empty ID")
// 	}

// 	currentVersion, err := t.GetLinkedVersion(tool)
// 	if err != nil { // log the error but continue
// 		fmt.Printf("failed to get linked version for tool %s: %v\n", tool.GetId(), err)
// 	}

// 	if currentVersion == version && !force {
// 		return fmt.Errorf("tool %s is already linked to version %s", tool.GetId(), version)
// 	}

// 	var successfullyLinked = false
// 	defer func() {
// 		if !successfullyLinked {
// 			if currentVersion != "" {
// 				// revert to the previous version
// 				_, revertErr := t.ExecuteScriptSteps(tool, "linkTool", string(currentVersion))
// 				if revertErr != nil {
// 					fmt.Printf("failed to revert to previous version %s: %v\n", currentVersion, revertErr)
// 				}
// 			} else {
// 				// if there was no previous version, just unlink the tool
// 				if unlinkErr := t.UnlinkTool(tool); unlinkErr != nil {
// 					fmt.Printf("failed to unlink tool %s: %v\n", tool.GetId(), unlinkErr)
// 				}
// 			}
// 		}
// 	}()
// 	_, err = t.ExecuteScriptSteps(tool, "linkTool", string(version))
// 	if err != nil {
// 		return fmt.Errorf("failed to link tool %s to version %s: %w", tool.GetId(), version, err)
// 	}
// 	err = t.StateService.SetToolState(tool.GetId(), version, force)
// 	if err != nil {
// 		return fmt.Errorf("failed to update tool state for %s: %w", tool.GetId(), err)
// 	}
// 	successfullyLinked = true
// 	return nil
// }

// func (t *ScriptsDrivenTVM) UnlinkTool(tool models.Tool) error {
// 	var successfullyUnlinked = false
// 	currentVersion, err := t.GetLinkedVersion(tool)
// 	if err != nil {
// 		return fmt.Errorf("failed to get linked version: %w", err)
// 	}
// 	if currentVersion == "" {
// 		return fmt.Errorf("tool %s is not linked to any version", tool.GetId())
// 	}
// 	defer func() {
// 		if !successfullyUnlinked {
// 			// if unlinking fails, try to revert to the previous version
// 			if revertErr := t.LinkTool(tool, currentVersion, true); revertErr != nil {
// 				fmt.Printf("failed to revert to previous version %s: %v\n", currentVersion, revertErr)
// 			}
// 		}
// 	}()
// 	_, err = t.ExecuteScriptSteps(tool, "unlinkTool", "")
// 	if err != nil {
// 		return fmt.Errorf("failed to unlink tool %s: %w", tool.GetId(), err)
// 	}
// 	err = t.StateService.SetToolState(tool.GetId(), "", false)
// 	if err != nil {
// 		return fmt.Errorf("failed to update tool state for %s: %w", tool.GetId(), err)
// 	}
// 	successfullyUnlinked = true

// 	return nil
// }

// func init() {
// 	models.ToolRegistrar.RegisterToolVersionManager("scripts_driven", &ScriptsDrivenTVM{})
// }
// //
