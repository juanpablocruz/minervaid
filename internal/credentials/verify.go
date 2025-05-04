package credentials

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mr-tron/base58"
)

func ResolveDidKeyPub(did string) (ed25519.PublicKey, error) {
	const prefix = "did:key:z"
	if !strings.HasPrefix(did, prefix) {
		return nil, fmt.Errorf("unsupported DID method")
	}
	data := did[len(prefix):]
	raw, err := base58.Decode(data)
	if err != nil {
		return nil, err
	}
	// strip multicodec prefix 0xed01
	if len(raw) < 2 || raw[0] != 0xed || raw[1] != 0x01 {
		return nil, fmt.Errorf("invalid multicodec prefix")
	}
	return ed25519.PublicKey(raw[2:]), nil
}

func VerifyCredential(cred *Credential) error {
	if len(cred.Proofs) == 0 {
		return fmt.Errorf("no proof present in credential")
	}
	// signature proof is the last proof in the array
	last := cred.Proofs[len(cred.Proofs)-1]
	var sp SignatureProof
	if err := json.Unmarshal(last, &sp); err != nil {
		return fmt.Errorf("unmarshal signature proof: %w", err)
	}
	sigBytes, err := hex.DecodeString(sp.JWS)
	if err != nil {
		return err
	}
	// serialize credential without proofs
	tmp := *cred
	tmp.Proofs = nil
	data, err := json.Marshal(tmp)
	if err != nil {
		return err
	}
	pub, err := ResolveDidKeyPub(cred.Issuer)
	if err != nil {
		return err
	}
	if !ed25519.Verify(pub, data, sigBytes) {
		return fmt.Errorf("invalid credential signature")
	}
	return nil
}

func VerifyPresentation(pres *Presentation) error {
	if len(pres.Proofs) == 0 {
		return fmt.Errorf("no proof present in presentation")
	}
	// signature proof is the last proof in the array
	last := pres.Proofs[len(pres.Proofs)-1]
	var sp SignatureProof
	if err := json.Unmarshal(last, &sp); err != nil {
		return fmt.Errorf("unmarshal signature proof: %w", err)
	}
	sigBytes, err := hex.DecodeString(sp.JWS)
	if err != nil {
		return err
	}
	// serialize presentation without proofs
	tmp := *pres
	tmp.Proofs = nil
	data, err := json.Marshal(tmp)
	if err != nil {
		return err
	}
	pub, err := ResolveDidKeyPub(pres.Holder)
	if err != nil {
		return err
	}
	if !ed25519.Verify(pub, data, sigBytes) {
		return fmt.Errorf("invalid presentation signature")
	}
	// verify all embedded credentials
	for _, vc := range pres.VerifiableCredential {
		if err := VerifyCredential(&vc); err != nil {
			return fmt.Errorf("embedded credential %s failed: %w", vc.ID, err)
		}
	}
	return nil
}
