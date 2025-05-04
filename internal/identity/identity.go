package identity

import (
	"crypto/ed25519"
	"encoding/json"
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

func BuildDIDDocument(did string, pub ed25519.PublicKey) ([]byte, error) {
	vmID := did + "#keys-1"
	pubBase58 := base58.Encode(pub)
	doc := map[string]interface{}{
		"@context": []interface{}{"https://www.w3.org/ns/did/v1"},
		"id":       did,
		"verificationMethod": []map[string]interface{}{{
			"id":              vmID,
			"type":            "Ed25519VerificationKey2018",
			"controller":      did,
			"publicKeyBase58": pubBase58,
		}},
		"authentication": []string{vmID},
	}
	return json.MarshalIndent(doc, "", "  ")
}
