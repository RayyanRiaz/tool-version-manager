package models

import (
	"encoding/json"
	"fmt"
	"strings"
)

/*
First we define (un)marshal methods for ToolWrapper to handle serialization and deserialization.
This allows us to wrap any Tool implementation and serialize it as JSON or YAML.
*/

func (t ToolWrapper) MarshalJSON() ([]byte, error) {
	if t.Wrapped == nil {
		return nil, fmt.Errorf("wrapped tool is nil")
	}

	return json.Marshal(t.Wrapped)
}

func (t ToolWrapper) MarshalYAML() (any, error) {
	if t.Wrapped == nil {
		return nil, fmt.Errorf("wrapped tool is nil")
	}
	return t.Wrapped, nil
}

type toolWrapperPeek struct {
	Type string `json:"type"`
}

func (t *ToolWrapper) unmarshal(unmarshal func(any) error) error {
	var peek toolWrapperPeek
	if err := unmarshal(&peek); err != nil {
		return err
	}

	tvm, err := ToolRegistrar.GetTVM(peek.Type)
	if err != nil {
		return err
	}

	tool := tvm.CreateNewTool()
	if err := unmarshal(tool); err != nil {
		return err
	}

	t.Wrapped = tool
	return nil
}

func (t *ToolWrapper) UnmarshalJSON(data []byte) error {
	return t.unmarshal(func(v any) error {
		return json.Unmarshal(data, v)
	})
}

func (t *ToolWrapper) UnmarshalYAML(unmarshal func(any) error) error {
	return t.unmarshal(unmarshal)
}

/*
Next, we define the unmarshal methods for UniqueToolWrappers.
This type is a slice of ToolWrapper that ensures all tools have unique IDs.
*/

func (tools *UniqueToolWrappers) unmarshal(unmarshal func(any) error) error {
	// Create an alias to avoid recursion, unmarshal into it, and copy the values to tools
	type Alias UniqueToolWrappers
	aux := &Alias{}
	if err := unmarshal(aux); err != nil {
		return err
	}
	*tools = UniqueToolWrappers(*aux)

	// validate the uniqueness of IDs
	ids := make(map[string]bool)
	for _, tool := range *tools {
		id := tool.Wrapped.GetId()
		if ids[id] {
			return fmt.Errorf("duplicate tool id: %s", id)
		}
		if strings.TrimSpace(id) == "" {
			return fmt.Errorf("tool id cannot be empty")
		}

		if strings.TrimSpace(id) != id {
			return fmt.Errorf("remove spaces around id: %s", id)
		}

		if id == "all" {
			return fmt.Errorf("`all` is a reserved tool-type")
		}
		ids[id] = true
	}
	return nil

}

func (tools *UniqueToolWrappers) UnmarshalJSON(data []byte) error {
	return tools.unmarshal(func(v any) error {
		return json.Unmarshal(data, v)
	})
}

func (tools *UniqueToolWrappers) UnmarshalYAML(unmarshal func(any) error) error {
	return tools.unmarshal(unmarshal)
}
