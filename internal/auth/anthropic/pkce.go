package anthropicAuth

import (
	"crypto/rand"
	"crypto/sha256"
    "encoding/base64"
)

type PKCE struct {
	Verifier string
	Challenge string
}

func GeneratePKCE() *PKCE {
	verifierBytes := make([]byte,32)
	if _ , err := rand.Read(verifierBytes); err != nil {
		return nil
	}

	verifierString := base64.RawURLEncoding.EncodeToString(verifierBytes)

	hash := sha256.Sum256([]byte(verifierString))

	challenge := base64.RawURLEncoding.EncodeToString(hash[:])

	return &PKCE{
		Verifier: verifierString,
		Challenge: challenge,
	}
}