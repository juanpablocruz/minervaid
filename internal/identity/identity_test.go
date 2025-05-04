package identity

import (
	"crypto/ed25519"
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {
	pub, priv, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("Error in GenerateKeyPair: %v", err)
	}

	if len(pub) != ed25519.PublicKeySize {
		t.Errorf("Size of public key invalid: %d", len(pub))
	}

	if len(priv) != ed25519.PrivateKeySize {
		t.Errorf("Size of private key invalid: %d", len(pub))
	}
}

func TestGenerateDID(t *testing.T) {
	pub, _, _ := GenerateKeyPair()
	did := GenerateDID(pub)
	if len(did) == 0 || did[:9] != "did:key:z" {
		t.Errorf("DID invalid: %s (%s)", did, did[:9])
	}
}

func TestEncodePrivateKey(t *testing.T) {
	_, priv, _ := GenerateKeyPair()
	enc := EncodePrivateKey(priv)
	if len(enc) == 0 {
		t.Error("Encoded private key is empty")
	}
}
