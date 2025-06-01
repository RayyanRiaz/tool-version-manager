package models

// symlinks
type ToolSymlink struct {
	From string `json:"from"`
	To   string `json:"to,omitempty"`
}

// version. simple string for now
type ToolVersion string

// Tool interface
type Tool interface {
	GetId() string
	GetType() string
	GetSymlinks() []ToolSymlink
}

type ToolBase struct {
	Id       string        `json:"id"`
	Type     string        `json:"type"`
	Symlinks []ToolSymlink `json:"symlinks,omitempty"`
}

func (t ToolBase) GetId() string {
	return t.Id
}

func (t ToolBase) GetType() string {
	return t.Type
}

func (t ToolBase) GetSymlinks() []ToolSymlink {
	return t.Symlinks
}

type ToolWrapper struct {
	Wrapped Tool
}

var _ Tool = (*ToolBase)(nil)
