package cmd

import (
	"fmt"
	"os"

	"rayyanriaz/tool-version-manager/pkg/impl/config"
	scriptdriventvm "rayyanriaz/tool-version-manager/pkg/impl/scriptdriven_tvm"
	"rayyanriaz/tool-version-manager/pkg/models"
)

var (
	configPath    string
	verbose       bool
	configService config.LocalFileConfig
)

func bootstrap() error {
	// defaults - priority: CLI flag > ENV > default
	if configPath == "" {
		configPath = os.Getenv("TVM_CONFIG")
	}
	if configPath == "" {
		configPath = "tools.yaml"
	}

	/*
		Bootstrap different TVMs with root info.
		Until a better design is made, I am just assigning the missing globals in this function
	*/

	// for scriptdriven.ScriptsDrivenTVM
	cfg := config.NewLocalFileConfig(configPath)
	models.ToolRegistrar.RegisterConfig("scripts_driven", cfg)
	models.ToolRegistrar.RegisterTVM("scripts_driven", scriptdriventvm.NewScriptsDrivenTVM())
	if err := cfg.Load(); err != nil {
		return fmt.Errorf("failed to create config service for script_driven: %w", err)
	}

	return nil
}

func getToolById(toolID string) (models.Tool, error) {

	for _, tool_type := range models.ToolRegistrar.GetRegisteredConfigTypes() {
		cfg, err := models.ToolRegistrar.GetConfig(tool_type)
		if err != nil {
			return nil, fmt.Errorf("failed to get config service for tool type '%s': %w", tool_type, err)
		}
		for _, toolWrapper := range cfg.GetTools() {
			if toolWrapper.Wrapped.GetId() == toolID {
				return toolWrapper.Wrapped, nil
			}
		}
	}
	// if we reach here, we didn't find the tool in any registered config service
	return nil, fmt.Errorf("tool '%s' not found in any configuration", toolID)
}

func getToolWithTVM(toolId string) (models.Tool, models.ToolVersionManager, error) {
	tool, err := getToolById(toolId)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get tool '%s': %w", toolId, err)
	}

	tvm, err := models.ToolRegistrar.GetTVM(tool.GetType())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get TVM for tool type '%s': %w", tool.GetType(), err)
	}

	return tool, tvm, nil
}

func getAllTools() (models.UniqueToolWrappers, error) {
	var allTools models.UniqueToolWrappers

	for _, tool_type := range models.ToolRegistrar.GetRegisteredConfigTypes() {
		cfg, err := models.ToolRegistrar.GetConfig(tool_type)
		if err != nil {
			return nil, fmt.Errorf("failed to get config for tool type '%s': %w", tool_type, err)
		}
		allTools = append(allTools, cfg.GetTools()...)
	}

	return allTools, nil
}
