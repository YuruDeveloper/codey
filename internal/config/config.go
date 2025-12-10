package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	appError "github.com/YuruDeveloper/codey/internal/error"
	"github.com/YuruDeveloper/codey/internal/ports"
)

var _ ports.AppConfig = (*Config)(nil)

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

func (instance *Config) getPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "" , appError.NewError(appError.FailFindHomeDir,err)
	}
	path := filepath.Join(home, ".codey", "config.json")
	return path , nil
}

func (instance *Config) Load() error {
	path  , err:= instance.getPath()
	if err != nil {
		return err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return appError.NewError(appError.FailReadFile,err)
	}
	err = json.Unmarshal(data, instance)
	if err != nil {
		return appError.NewError(appError.JsonUnMarshalError,err)
	}
	return nil
}

func (instance *Config) Save() error {
	path , err := instance.getPath()
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return appError.NewError(appError.FailMakeFolder,err)
	}
	data, err := json.MarshalIndent(instance, "", " ")
	if err != nil {
		return appError.NewError(appError.JsonMarshalError,err)
	}
	err = os.WriteFile(path,data,0600)
	if err != nil {
		return appError.NewError(appError.FailMakeFile,err)
	}
	return nil
}

func (instance *Config) GetProviderAuth(name string) json.RawMessage {
	return instance.Providers[name]
}

func (instance *Config) SetProviderAuth(name string, data json.RawMessage) {
	instance.Providers[name] = data
}
