package models

type UniqueToolWrappers []ToolWrapper

type Config interface {
	Load() error
	Save() error
	GetTools() UniqueToolWrappers
}
