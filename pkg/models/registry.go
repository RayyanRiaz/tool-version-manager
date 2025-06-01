package models

import "fmt"

type tvmRegistry map[string]ToolVersionManager
type configRegistry map[string]Config

type ToolRegistry struct {
	*tvmRegistry
	*configRegistry
}

func (r *ToolRegistry) RegisterTVM(toolType string, manager ToolVersionManager) {
	if _, exists := (*r.tvmRegistry)[toolType]; exists {
		panic("Tool already registered: " + toolType)
	}
	(*r.tvmRegistry)[toolType] = manager
}

func (r *ToolRegistry) RegisterConfig(toolType string, config Config) {
	if _, exists := (*r.configRegistry)[toolType]; exists {
		panic("Config service already registered for tool type: " + toolType)
	}
	(*r.configRegistry)[toolType] = config
}

func (r *ToolRegistry) GetTVM(toolType string) (ToolVersionManager, error) {
	if manager, exists := (*r.tvmRegistry)[toolType]; exists {
		return manager, nil
	}
	return nil, fmt.Errorf("Tool manager not registered for %s. Allowed types are: %v", toolType, ToolRegistrar.GetRegisteredToolTypes())
}

func (r *ToolRegistry) GetRegisteredToolTypes() []string {
	types := []string{}
	for k := range *r.tvmRegistry {
		types = append(types, k)
	}
	return types
}
func (r *ToolRegistry) GetConfig(toolType string) (Config, error) {
	if config, exists := (*r.configRegistry)[toolType]; exists {
		return config, nil
	}
	return nil, fmt.Errorf("Config service not registered for tool type: %s", toolType)
}

func (r *ToolRegistry) GetRegisteredConfigTypes() []string {
	types := []string{}
	for k := range *r.configRegistry {
		types = append(types, k)
	}
	return types
}

var ToolRegistrar = ToolRegistry{
	tvmRegistry:    &tvmRegistry{},
	configRegistry: &configRegistry{},
}
