package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"text/template"

	"github.com/goccy/go-yaml"
)

func LoadJSONFile[T any](filePath string, out *T) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("failed to unmarshal JSON from file %s: %w", filePath, err)
	}
	return nil
}

func LoadYAMLFile[T any](filePath string, out *T) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	if len(bytes.TrimSpace(data)) == 0 {
		return fmt.Errorf("YAML file %s is empty", filePath)
	}
	if err := yaml.UnmarshalWithOptions(data, out, yaml.AllowDuplicateMapKey()); err != nil {
		return fmt.Errorf("failed to unmarshal YAML from file %s: %w", filePath, err)
	}
	return nil
}

func SaveJSONFile[T any](filePath string, data *T) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %w", err)
	}
	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write JSON to file %s: %w", filePath, err)
	}
	return nil
}

func RenderTemplate(tmplStr string, data map[string]any) (string, error) {
	tmpl, err := template.New("cmd").Parse(tmplStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func SaveYAMLFile[T any](filePath string, data *T) error {

	yamlData, err := yaml.MarshalWithOptions(data, yaml.UseLiteralStyleIfMultiline(true))

	// yamlData, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data to YAML: %w", err)
	}
	if err := os.WriteFile(filePath, yamlData, 0644); err != nil {
		return fmt.Errorf("failed to write YAML to file %s: %w", filePath, err)
	}
	return nil
}

func SaveFile[T any](filePath string, data *T) error {
	return SaveYAMLFile(filePath, data)
}
func LoadFile[T any](filePath string, out *T) error {
	return LoadYAMLFile(filePath, out)
}
