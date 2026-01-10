package config

import (
	"fmt"
	"os"

	"rayyanriaz/tool-version-manager/pkg/models"
	"rayyanriaz/tool-version-manager/pkg/utils"
)

type LocalFileConfig struct {
	configFilePath              string                    `json:"-"`
	Tools                       models.UniqueToolWrappers `json:"tools"`
	DownloadsDir                string                    `json:"downloads_dir,omitempty"`
	SymlinksDir                 string                    `json:"symlinks_dir,omitempty"`
	GitHubToken                 string                    `json:"github_token,omitempty"`
	RemoteVersionsCacheFilePath string                    `json:"remote_versions_cache_file_path,omitempty"`
}

func NewLocalFileConfig(configPath string) *LocalFileConfig {
	var config LocalFileConfig
	config.configFilePath = configPath

	return &config
}

func (c *LocalFileConfig) GetTools() models.UniqueToolWrappers {
	return c.Tools
}

func (c *LocalFileConfig) Load() error {
	if err := utils.LoadFile(c.configFilePath, c); err != nil {
		return fmt.Errorf("failed to load config file %s: %w", c.configFilePath, err)
	}

	if c.DownloadsDir == "" {
		c.DownloadsDir = "./tvm_cache"
	}
	if c.SymlinksDir == "" {
		c.SymlinksDir = "./bin"
	}
	if c.RemoteVersionsCacheFilePath == "" {
		c.RemoteVersionsCacheFilePath = "./.tools.state.yaml"
	}

	// Environment variable takes precedence over config file
	if envToken := os.Getenv("GITHUB_TOKEN"); envToken != "" {
		c.GitHubToken = envToken
	}

	if err := c.ensureDirectories(); err != nil {
		return fmt.Errorf("failed to ensure directories: %w", err)
	}
	return nil
	// return utils.LoadFile(c.configFilePath, c)
}

func (c *LocalFileConfig) Save() error {
	return utils.SaveFile(c.configFilePath, c)
}

func (l *LocalFileConfig) ensureDirectories() error {
	directories := []string{l.DownloadsDir, l.SymlinksDir}

	for _, dir := range directories {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return nil
}

var _ models.Config = (*LocalFileConfig)(nil)
