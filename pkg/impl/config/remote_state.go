package config

import (
	"os"
	"time"

	"rayyanriaz/tool-version-manager/pkg/models"
	"rayyanriaz/tool-version-manager/pkg/utils"
)

// ToolVersionCache holds cached version info for a single tool
type ToolVersionCache struct {
	LatestVersion models.ToolVersion `json:"latest_version"`
	LastChecked   time.Time          `json:"last_checked"`
}

// RemoteVersionsCache holds cached latest versions for all tools
type RemoteVersionsCache struct {
	filePath string                      `json:"-"`
	Tools    map[string]ToolVersionCache `json:"tools"`
}

// NewRemoteVersionsCache creates a new cache instance
func NewRemoteVersionsCache(filePath string) *RemoteVersionsCache {
	return &RemoteVersionsCache{
		filePath: filePath,
		Tools:    make(map[string]ToolVersionCache),
	}
}

// Load reads the cache from disk. Returns nil error if file doesn't exist.
func (c *RemoteVersionsCache) Load() error {
	if _, err := os.Stat(c.filePath); os.IsNotExist(err) {
		// File doesn't exist yet, that's fine
		c.Tools = make(map[string]ToolVersionCache)
		return nil
	}

	if err := utils.LoadFile(c.filePath, c); err != nil {
		// If file exists but is corrupted/empty, start fresh
		c.Tools = make(map[string]ToolVersionCache)
		return nil
	}

	if c.Tools == nil {
		c.Tools = make(map[string]ToolVersionCache)
	}

	return nil
}

// Save writes the cache to disk
func (c *RemoteVersionsCache) Save() error {
	return utils.SaveFile(c.filePath, c)
}

// GetCachedVersion returns the cached latest version for a tool, or empty if not cached
func (c *RemoteVersionsCache) GetCachedVersion(toolID string) (models.ToolVersion, time.Time, bool) {
	if cache, ok := c.Tools[toolID]; ok {
		return cache.LatestVersion, cache.LastChecked, true
	}
	return "", time.Time{}, false
}

// SetCachedVersion updates the cached latest version for a tool
func (c *RemoteVersionsCache) SetCachedVersion(toolID string, version models.ToolVersion) {
	c.Tools[toolID] = ToolVersionCache{
		LatestVersion: version,
		LastChecked:   time.Now(),
	}
}
