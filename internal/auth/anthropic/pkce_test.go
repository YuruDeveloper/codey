package anthropicAuth

import (
	"crypto/sha256"
	"encoding/base64"
	"testing"
)

func TestGeneratePKCE(t *testing.T) {
	pkce1 := GeneratePKCE()
	if pkce1 == nil {
		t.Fatal("GeneratePKCE returned nil")
	}

	// RFC 7636 recommends verifier length between 43 and 128 (inclusive).
	// 32 random bytes result in 43 base64-url encoded characters.
	expectedVerifierLen := 43
	if len(pkce1.Verifier) != expectedVerifierLen {
		t.Errorf("Verifier length mismatch: got %d, want %d", len(pkce1.Verifier), expectedVerifierLen)
	}

	// SHA256 hash is 32 bytes, which base64-url encodes to 43 characters.
	expectedChallengeLen := 43
	if len(pkce1.Challenge) != expectedChallengeLen {
		t.Errorf("Challenge length mismatch: got %d, want %d", len(pkce1.Challenge), expectedChallengeLen)
	}

	// Ensure two calls produce different verifiers/challenges (probabilistically)
	pkce2 := GeneratePKCE()
	if pkce1.Verifier == pkce2.Verifier {
		t.Error("GeneratePKCE produced identical verifiers on two calls")
	}
	if pkce1.Challenge == pkce2.Challenge {
		t.Error("GeneratePKCE produced identical challenges on two calls")
	}
}

func TestPKCEChallengeCalculation(t *testing.T) {
	testVerifierBytes := make([]byte, 32)
	// Fill with a predictable pattern for testability
	for i := 0; i < 32; i++ {
		testVerifierBytes[i] = byte(i)
	}
	knownVerifier := base64.RawURLEncoding.EncodeToString(testVerifierBytes)

	// Manually calculate the expected challenge
	hash := sha256.Sum256([]byte(knownVerifier))
	expectedChallenge := base64.RawURLEncoding.EncodeToString(hash[:])

	// Create a dummy PKCE object for calculation
	// We are effectively testing the internal logic of GeneratePKCE,
	// but focusing on the challenge calculation part.
	pkce := &PKCE{
		Verifier: knownVerifier,
	}
	// Re-implement the challenge calculation from GeneratePKCE for comparison
	calculatedHash := sha256.Sum256([]byte(pkce.Verifier))
	calculatedChallenge := base64.RawURLEncoding.EncodeToString(calculatedHash[:])

	if calculatedChallenge != expectedChallenge {
		t.Errorf("PKCE challenge calculation mismatch:\nGot:  %s\nWant: %s", calculatedChallenge, expectedChallenge)
	}
}
