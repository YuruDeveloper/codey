package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestConfig_LoadAndGetProviderAuth(t *testing.T) {
	// 1. Setup: Create a temporary home directory and a dummy config.json
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	configDir := filepath.Join(tempDir, ".codey")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create temp config dir: %v", err)
	}

	dummyAuthConfig := `{"api_key": "dummy-key"}`
	dummyConfigJSON := []byte(`{
		"current_provider": "test_provider",
		"current_model": "test_model_1",
		"providers": {
			"test_provider": ` + dummyAuthConfig + `
		}
	}`)

	configPath := filepath.Join(configDir, "config.json")
	if err := os.WriteFile(configPath, dummyConfigJSON, 0644); err != nil {
		t.Fatalf("Failed to write dummy config file: %v", err)
	}

	// 2. Execute: Load the config
	cfg := New()
	err := cfg.Load()
	if err != nil {
		t.Fatalf("cfg.Load() failed: %v", err)
	}

	// 3. Assert: Check the loaded values
	if cfg.CurrentProvider != "test_provider" {
		t.Errorf("Expected CurrentProvider to be 'test_provider', got '%s'", cfg.CurrentProvider)
	}
	if cfg.CurrentModel != "test_model_1" {
		t.Errorf("Expected CurrentModel to be 'test_model_1', got '%s'", cfg.CurrentModel)
	}

	// Assert GetProviderAuth
	authData := cfg.GetProviderAuth("test_provider")
	if authData == nil {
		t.Fatalf("GetProviderAuth('test_provider') returned nil")
	}

	// Unmarshal both to compare content, as json.RawMessage is a byte slice
	var expectedAuth, actualAuth map[string]string
	if err := json.Unmarshal([]byte(dummyAuthConfig), &expectedAuth); err != nil {
		t.Fatalf("Failed to unmarshal dummy auth config: %v", err)
	}
	if err := json.Unmarshal(authData, &actualAuth); err != nil {
		t.Fatalf("Failed to unmarshal actual auth data: %v", err)
	}

	if !reflect.DeepEqual(actualAuth, expectedAuth) {
		t.Errorf("GetProviderAuth mismatch:\nGot:  %v\nWant: %v", actualAuth, expectedAuth)
	}
}

func TestConfig_SaveAndSetProviderAuth(t *testing.T) {
	// 1. Setup: Create a temporary home directory
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	// 2. Execute: Create a config, set auth data, and save it
	cfg := New()
	cfg.CurrentProvider = "saved_provider"
	cfg.CurrentModel = "saved_model"

	authData := json.RawMessage(`{"api_key": "saved-key"}`)
	cfg.SetProviderAuth("saved_provider", authData)

	err := cfg.Save()
	if err != nil {
		t.Fatalf("cfg.Save() failed: %v", err)
	}

	// 3. Assert: Read the saved file and verify its contents
	configPath := filepath.Join(tempDir, ".codey", "config.json")
	savedData, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read saved config file: %v", err)
	}

	var loadedCfg Config
	if err := json.Unmarshal(savedData, &loadedCfg); err != nil {
		t.Fatalf("Failed to unmarshal saved config data: %v", err)
	}

	// Compare fields individually because of json.RawMessage formatting differences
	if loadedCfg.CurrentProvider != cfg.CurrentProvider {
		t.Errorf("CurrentProvider mismatch: got %s, want %s", loadedCfg.CurrentProvider, cfg.CurrentProvider)
	}
	if loadedCfg.CurrentModel != cfg.CurrentModel {
		t.Errorf("CurrentModel mismatch: got %s, want %s", loadedCfg.CurrentModel, cfg.CurrentModel)
	}

	// Compare provider auth data by unmarshalling it
	loadedAuthData := loadedCfg.GetProviderAuth("saved_provider")
	originalAuthData := cfg.GetProviderAuth("saved_provider")

	var loadedAuth, originalAuth map[string]any
	if err := json.Unmarshal(loadedAuthData, &loadedAuth); err != nil {
		t.Fatalf("Failed to unmarshal loaded auth data: %v", err)
	}
	if err := json.Unmarshal(originalAuthData, &originalAuth); err != nil {
		t.Fatalf("Failed to unmarshal original auth data: %v", err)
	}

	if !reflect.DeepEqual(loadedAuth, originalAuth) {
		t.Errorf("Provider auth data mismatch:\nGot:  %+v\nWant: %+v", loadedAuth, originalAuth)
	}
}
