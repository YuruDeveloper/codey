package anthropicAuth

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/YuruDeveloper/codey/internal/config"
	"golang.org/x/oauth2"
)

func TestNewOAuthAuth(t *testing.T) {
	// 1. Setup
	cfg := config.New()
	authData := AuthData{
		Type:    OAuth,
		Access:  "test-access-token",
		Refresh: "test-refresh-token",
		Expires: time.Now().Add(1 * time.Hour).Unix(),
	}
	rawData, err := json.Marshal(authData)
	if err != nil {
		t.Fatalf("Failed to marshal auth data: %v", err)
	}
	cfg.SetProviderAuth(name, rawData) // 'name' is "anthropic"

	// 2. Execute
	oauthAuth := NewOAuthAuth(cfg)

	// 3. Assert
	if oauthAuth == nil {
		t.Fatal("NewOAuthAuth returned nil")
	}

	if oauthAuth.Key() != "test-access-token" {
		t.Errorf("Key() mismatch: got '%s', want '%s'", oauthAuth.Key(), "test-access-token")
	}

	// Check internal token state indirectly via Save
	savedCfg := config.New()
	oauthAuth.Save(savedCfg)
	savedRawData := savedCfg.GetProviderAuth(name)

	var savedAuthData AuthData
	if err := json.Unmarshal(savedRawData, &savedAuthData); err != nil {
		t.Fatalf("Failed to unmarshal saved data: %v", err)
	}

	// Correctly compare Expires, allowing for small differences in time representation
	if savedAuthData.Expires != authData.Expires {
		t.Errorf("Saved Expires does not match initial Expires.\nGot: %d\nWant: %d", savedAuthData.Expires, authData.Expires)
	}
	savedAuthData.Expires = authData.Expires // Nullify expiry diff for DeepEqual
	if !reflect.DeepEqual(savedAuthData, authData) {
		t.Errorf("Saved auth data does not match initial data.\nGot: %+v\nWant:%+v", savedAuthData, authData)
	}
}

func TestOAuthAuth_Save(t *testing.T) {
	// 1. Setup
	token := &oauth2.Token{
		AccessToken:  "my-access-token",
		RefreshToken: "my-refresh-token",
		Expiry:       time.Now().Add(30 * time.Minute),
	}
	oauthAuth := &OAuthAuth{token: token}
	cfg := config.New()

	// 2. Execute
	oauthAuth.Save(cfg)

	// 3. Assert
	rawData := cfg.GetProviderAuth(name)
	if rawData == nil {
		t.Fatal("Saved auth data is nil")
	}

	var loadedAuthData AuthData
	if err := json.Unmarshal(rawData, &loadedAuthData); err != nil {
		t.Fatalf("Failed to unmarshal saved auth data: %v", err)
	}

	expectedAuthData := AuthData{
		Type:    OAuth,
		Access:  "my-access-token",
		Refresh: "my-refresh-token",
		Expires: token.Expiry.Unix(),
	}

	if loadedAuthData.Expires != expectedAuthData.Expires {
		t.Errorf("Saved Expires does not match initial Expires.\nGot: %d\nWant: %d", loadedAuthData.Expires, expectedAuthData.Expires)
	}
	loadedAuthData.Expires = expectedAuthData.Expires // Nullify expiry diff for DeepEqual
	if !reflect.DeepEqual(loadedAuthData, expectedAuthData) {
		t.Errorf("Saved auth data mismatch.\nGot:  %+v\nWant: %+v", loadedAuthData, expectedAuthData)
	}
}

func TestOAuthAuth_Update(t *testing.T) {
	t.Run("Token is valid", func(t *testing.T) {
		// 1. Setup
		originalToken := &oauth2.Token{
			AccessToken:  "valid-access-token",
			RefreshToken: "valid-refresh-token",
			Expiry:       time.Now().Add(1 * time.Hour),
		}
		oauthAuth := &OAuthAuth{token: originalToken}

		// 2. Execute
		oauthAuth.Update(context.Background()) // Assumes Update now takes a context

		// 3. Assert
		if oauthAuth.Key() != "valid-access-token" {
			t.Errorf("Token was refreshed when it shouldn't have been. Got key: %s", oauthAuth.Key())
		}
	})

	t.Run("Token is expired and refresh succeeds", func(t *testing.T) {
		// 1. Setup: Mock server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/token" {
				t.Fatalf("Expected to request '/token', got: %s", r.URL.Path)
			}
			body, _ := io.ReadAll(r.Body)
			if !strings.Contains(string(body), "grant_type=refresh_token") {
				t.Error("Request body does not contain grant_type=refresh_token")
			}
			if !strings.Contains(string(body), "refresh_token=expired-refresh-token") {
				t.Error("Request body does not contain correct refresh_token")
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"access_token":  "new-access-token",
				"refresh_token": "new-refresh-token",
				"expires_in":    3600,
			})
		}))
		defer server.Close()

		// 2. Setup: Expired token and auth object
		expiredToken := &oauth2.Token{
			AccessToken:  "expired-access-token",
			RefreshToken: "expired-refresh-token",
			Expiry:       time.Now().Add(-1 * time.Hour), // Expired
		}
		mockOauth2Config := &oauth2.Config{
			ClientID: ClientID,
			Endpoint: oauth2.Endpoint{TokenURL: server.URL + "/token"},
		}
		oauthAuth := &OAuthAuth{
			token:        expiredToken,
			oauth2Config: mockOauth2Config,
		}

		// 3. Setup: Context with mock http client
		mockClient := server.Client()
		ctx := context.WithValue(context.Background(), oauth2.HTTPClient, mockClient)

		// 4. Execute
		oauthAuth.Update(ctx) // Assumes Update now takes a context

		// 5. Assert
		if oauthAuth.Key() != "new-access-token" {
			t.Errorf("Expected key to be 'new-access-token', got '%s'", oauthAuth.Key())
		}
	})
}
