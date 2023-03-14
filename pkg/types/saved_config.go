package types

// SavedConfig ...
type SavedConfig struct {
	Current string            `json:"current" yaml:"current"`
	Items   []SavedConfigItem `json:"items" yaml:"items"`
}

// SavedConfigItem ...
type SavedConfigItem struct {
	Name    string `json:"name" yaml:"name"`
	Content string `json:"content" yaml:"content"`
}
