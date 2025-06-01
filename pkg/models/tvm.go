package models

type ToolFactory func() Tool

type ToolLinkInfo struct {
	Version  ToolVersion `json:"version"`
	LinkedAt string      `json:"linked_at"`
}

type ToolDiscovery interface {
	GetAllLocalVersions(tool Tool) ([]ToolVersion, error)
	GetAllRemoteVersions(tool Tool) ([]ToolVersion, error)
	GetLatestRemoteVersion(tool Tool) (ToolVersion, error)
}

type ToolLinker interface {
	LinkTool(tool Tool, version ToolVersion) error
	UnlinkTool(tool Tool) error
	GetLinkInfo(tool Tool) (*ToolLinkInfo, error)
}

type ToolInstaller interface {
	InstallToolForVersion(tool Tool, version ToolVersion) error
}

type ToolComparer interface {
	CompareVersions(tool Tool, v1 ToolVersion, v2 ToolVersion) (int, error)
}

type ToolVersionManager interface {
	ToolDiscovery
	ToolLinker
	ToolInstaller
	ToolComparer
	CreateNewTool() Tool
}
