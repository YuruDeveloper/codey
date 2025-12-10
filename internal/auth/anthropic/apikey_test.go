package anthropicAuth

import (
	"encoding/json"
	"testing"

	"github.com/YuruDeveloper/codey/internal/config"
)

func TestNewApiKeyAuth(t *testing.T) {
	// 1. Setup
	authData := AuthData{
		Type: ApiKey,
		Key:  "test-api-key",
	}

	// 2. Execute
	apiKeyAuth := NewApiKeyAuth(authData)

	// 3. Assert
	if apiKeyAuth == nil {
		t.Fatal("NewApiKeyAuth returned nil")
	}

	expectedKey := "test-api-key"
	actualKey := apiKeyAuth.Key()
	if actualKey != expectedKey {
		t.Errorf("Key() mismatch: got '%s', want '%s'", actualKey, expectedKey)
	}
}

func TestApiKeyAuth_Save(t *testing.T) {
	// 1. Setup
	apiKeyAuth := &ApiKeyAuth{}
	apiKeyAuth.SetApiKey("my-saved-key")
	cfg := config.New()

	// 2. Execute
	apiKeyAuth.Save(cfg)

	// 3. Assert
	rawData := cfg.GetProviderAuth(name) // 'name' is the constant "anthropic"
	if rawData == nil {
		t.Fatal("Saved auth data is nil")
	}

	var loadedAuthData AuthData
	err := json.Unmarshal(rawData, &loadedAuthData)
	if err != nil {
		t.Fatalf("Failed to unmarshal saved auth data: %v", err)
	}

	if loadedAuthData.Type != ApiKey {
		t.Errorf("Expected auth type to be '%s', got '%s'", ApiKey, loadedAuthData.Type)
	}
	if loadedAuthData.Key != "my-saved-key" {
		t.Errorf("Expected auth key to be 'my-saved-key', got '%s'", loadedAuthData.Key)
	}
}
