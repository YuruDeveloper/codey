package anthropicAuth

type authType string

const (
	OAuth  = authType("oauth")
	ApiKey = authType("apikey")
)

const (
	name = "anthropic"
)

type AuthData struct {
	Type authType `json:"type"`

	// OAuth 필드들
	Refresh string `json:"refresh,omitempty"`
	Access  string `json:"access,omitempty"`
	Expires int64  `json:"expires,omitempty"`

	// API Key 필드
	Key string `json:"key,omitempty"`
}
