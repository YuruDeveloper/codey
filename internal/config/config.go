package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Providers       map[string]json.RawMessage `json:"providers"`
	CurrentProvider string                     `json:"current_provider"`
	CurrentModel    string                     `json:"current_model"`
}

func New() *Config {
	return &Config{
		Providers:       make(map[string]json.RawMessage, 1),
		CurrentProvider: "",
		CurrentModel:    "",
	}
}

func (instance *Config) GetPath() string {
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, ".codey", "config.json")
	return path
}

func (instance *Config) Load() error {

	data, err := os.ReadFile(instance.GetPath())
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, instance)
	return err
}

func (instance *Config) Save() error {
	os.MkdirAll(filepath.Dir(instance.GetPath()), 0755)

	data, _ := json.MarshalIndent(instance, "", " ")
	return os.WriteFile(instance.GetPath(), data, 0600)
}

func (instance *Config) GetProviderAuth(name string) json.RawMessage {
	return instance.Providers[name]
}

func (instance *Config) SetProviderAuth(name string, data json.RawMessage) {
	instance.Providers[name] = data
}
