package identity

import (
	"crypto/ed25519"
	"fmt"

	"github.com/mr-tron/base58"
)

func GenerateKeyPair() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	return ed25519.GenerateKey(nil)
}

func GenerateDID(pub ed25519.PublicKey) string {
	prefix := []byte{0xED, 0x01}
	data := append(prefix, pub...)
	enc := base58.Encode(data)
	return fmt.Sprintf("did:key:z%s", enc)
}

func EncodePrivateKey(priv ed25519.PrivateKey) string {
	return base58.Encode(priv)
}
