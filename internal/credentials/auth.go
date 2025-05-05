package credentials

import (
	"crypto/ed25519"
	crand "crypto/rand"
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/juanpablocruz/minervaid/internal/store"
	"github.com/mr-tron/base58"
)

// AuthenticationChallenge represents a nonce challenge for DID auth
type AuthenticationChallenge struct {
	Type      string    `json:"type"`
	Challenge string    `json:"challenge"`
	Domain    string    `json:"domain,omitempty"`
	IssuedAt  time.Time `json:"issuedAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// AuthenticationResponse contains the signed challenge
type AuthenticationResponse struct {
	Challenge *AuthenticationChallenge `json:"challenge"`
	Proof     *SignatureProof          `json:"proof"`
}

// NewAuthChallenge generates a fresh challenge for a domain, valid for a duration
// Uses crypto/rand for secure randomness
func NewAuthChallenge(domain string, validity time.Duration) *AuthenticationChallenge {
	nonce := make([]byte, 16)
	if _, err := crand.Read(nonce); err != nil {
		panic(fmt.Sprintf("failed to generate nonce: %v", err))
	}
	ch := &AuthenticationChallenge{
		Type:      "AuthenticationChallenge",
		Challenge: base58.Encode(nonce),
		Domain:    domain,
		IssuedAt:  time.Now().UTC(),
		ExpiresAt: time.Now().Add(validity).UTC(),
	}
	return ch
}

// RespondAuthChallenge loads the holder's private key from storeDir/keystore.json,
// signs the serialized challenge, and returns a response containing the proof.
func RespondAuthChallenge(did string, ch *AuthenticationChallenge, storeDir string) (*AuthenticationResponse, error) {
	ksPath := filepath.Join(storeDir, "keystore.json")
	ks := store.LoadKeyStore(ksPath)
	privEnc, ok := ks[did]
	if !ok {
		return nil, fmt.Errorf("unknown DID: %s", did)
	}
	privBytes, err := base58.Decode(privEnc)
	if err != nil {
		return nil, fmt.Errorf("decoding private key: %w", err)
	}
	// Serialize challenge
	data, err := json.Marshal(ch)
	if err != nil {
		return nil, fmt.Errorf("marshaling challenge: %w", err)
	}
	// Sign with Ed25519
	sig := ed25519.Sign(privBytes, data)
	proof := &SignatureProof{
		Type:               "Ed25519Signature2018",
		Created:            time.Now().UTC().Format(time.RFC3339),
		ProofPurpose:       "authentication",
		VerificationMethod: did + "#keys-1",
		JWS:                fmt.Sprintf("%x", sig),
	}
	return &AuthenticationResponse{Challenge: ch, Proof: proof}, nil
}
