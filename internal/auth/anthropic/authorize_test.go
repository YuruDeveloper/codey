package anthropicAuth

import (
	"net/url"
	"testing"
)

func TestAuthorizeURL(t *testing.T) {
	testCases := []struct {
		name         string
		mode         AuthMode
		expectedHost string
	}{
		{"Console Mode", console, "console.anthropic.com"},
		{"Max Mode", Max, "claude.ai"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			authURL, pkce := AuthorizeURL(tc.mode)

			if pkce == nil {
				t.Fatal("AuthorizeURL returned a nil PKCE struct")
			}
			if authURL == "" {
				t.Fatal("AuthorizeURL returned an empty URL string")
			}

			parsedURL, err := url.Parse(authURL)
			if err != nil {
				t.Fatalf("Failed to parse returned URL: %v", err)
			}

			if parsedURL.Host != tc.expectedHost {
				t.Errorf("Expected host %s, got %s", tc.expectedHost, parsedURL.Host)
			}

			query := parsedURL.Query()
			expectedParams := []string{
				"client_id",
				"response_type",
				"redirect_uri",
				"scope",
				"code_challenge",
				"code_challenge_method",
				"state",
				"code",
			}

			for _, param := range expectedParams {
				if !query.Has(param) {
					t.Errorf("URL is missing required query parameter: %s", param)
				}
			}

			if query.Get("code_challenge") != pkce.Challenge {
				t.Error("URL code_challenge does not match PKCE challenge")
			}
			if query.Get("state") != pkce.Verifier {
				t.Error("URL state does not match PKCE verifier")
			}
			if query.Get("code_challenge_method") != "S256" {
				t.Error("URL code_challenge_method is not S256")
			}
		})
	}
}
