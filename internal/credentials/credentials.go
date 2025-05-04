package credentials

import (
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"time"
)

type Credential struct {
	Context           []string               `json:"@context"`
	ID                string                 `json:"id"`
	Type              []string               `json:"type"`
	Issuer            string                 `json:"issuer"`
	IssuanceDate      time.Time              `json:"issuanceDate"`
	CredentialSubject map[string]interface{} `json:"credentialSubject"`
	Proofs            []json.RawMessage      `json:"proof"`
}

// SignatureProof is the Ed25519 signature proof
type SignatureProof struct {
	Type               string `json:"type"`
	Created            string `json:"created"`
	ProofPurpose       string `json:"proofPurpose"`
	VerificationMethod string `json:"verificationMethod"`
	JWS                string `json:"jws"`
}

// RangeProof wraps a zero-knowledge proof of age >= MinAge
type RangeProof struct {
	Type   string  `json:"type"`
	MinAge uint64  `json:"minAge"`
	Proof  ZKProof `json:"proof"`
}

func NewCredential(id, issuer string, subject map[string]interface{}) *Credential {
	return &Credential{
		Context:           []string{"https://www.w3.org/2018/credentials/v1"},
		ID:                id,
		Type:              []string{"VerifiableCredential"},
		Issuer:            issuer,
		IssuanceDate:      time.Now().UTC(),
		CredentialSubject: subject,
		Proofs:            []json.RawMessage{},
	}
}

func (c *Credential) AttachAgeProof(minAge uint64) error {
	raw, ok := c.CredentialSubject["age"]
	if !ok {
		return fmt.Errorf("subject missing 'age' field")
	}
	var ageVal uint64
	switch v := raw.(type) {
	case float64:
		ageVal = uint64(v)
	case int:
		ageVal = uint64(v)
	case int64:
		ageVal = uint64(v)
	case uint64:
		ageVal = v
	default:
		return fmt.Errorf("invalid age type %T", raw)
	}

	proof, err := GenerateAgeProof(ageVal, minAge)
	if err != nil {
		return fmt.Errorf("generating age proof: %w", err)
	}
	rp := RangeProof{Type: "BulletproofRangeProof", MinAge: minAge, Proof: *proof}
	data, err := json.Marshal(rp)
	if err != nil {
		return fmt.Errorf("marshaling range proof: %w", err)
	}
	c.Proofs = append(c.Proofs, data)
	delete(c.CredentialSubject, "age")
	return nil
}

func (c *Credential) SignCredential(priv ed25519.PrivateKey, verificationMethod string) error {
	// Serialize credential without the signature proof
	tmp := *c
	tmp.Proofs = nil
	data, err := json.Marshal(&tmp)
	if err != nil {
		return err
	}
	sig := ed25519.Sign(priv, data)
	sp := SignatureProof{
		Type:               "Ed25519Signature2018",
		Created:            time.Now().UTC().Format(time.RFC3339),
		ProofPurpose:       "assertionMethod",
		VerificationMethod: verificationMethod,
		JWS:                fmt.Sprintf("%x", sig),
	}
	b, err := json.Marshal(sp)
	if err != nil {
		return fmt.Errorf("marshaling signature proof: %w", err)
	}
	c.Proofs = append(c.Proofs, b)
	return nil
}
