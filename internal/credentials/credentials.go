package credentials

import (
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Credential represents a W3C Verifiable Credential.
type Credential struct {
	Context           []string               `json:"@context"`
	ID                string                 `json:"id"`
	Type              []string               `json:"type"`
	Issuer            string                 `json:"issuer"`
	IssuanceDate      time.Time              `json:"issuanceDate"`
	CredentialSubject map[string]interface{} `json:"credentialSubject"`
	Proofs            []json.RawMessage      `json:"proof"`
}

// SignatureProof is the Ed25519 signature proof.
type SignatureProof struct {
	Type               string `json:"type"`
	Created            string `json:"created"`
	ProofPurpose       string `json:"proofPurpose"`
	VerificationMethod string `json:"verificationMethod"`
	JWS                string `json:"jws"`
}

// Challenge describes a generic proof challenge.
type Challenge struct {
	Type   string                 `json:"type"`
	Params map[string]interface{} `json:"params"`
}

// NewCredential builds a new Verifiable Credential.
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

// AttachProof applies a zero-knowledge proof based on a generic Challenge JSON.
// Currently supports proof type "range" to generate a bulletproof range proof.
func (c *Credential) AttachProof(challengeJSON []byte) error {
	// Parse generic challenge
	var ch Challenge
	if err := json.Unmarshal(challengeJSON, &ch); err != nil {
		return fmt.Errorf("invalid challenge JSON: %w", err)
	}

	subj := c.CredentialSubject
	switch ch.Type {
	case "range":
		// expect Params: field (string), min (number or string)
		fldVal, ok := ch.Params["field"]
		fld, okStr := fldVal.(string)
		if !ok || !okStr || fld == "" {
			return fmt.Errorf("range challenge missing or invalid 'field'")
		}

		minRaw, ok := ch.Params["min"]
		if !ok {
			return fmt.Errorf("range challenge missing 'min'")
		}
		// min may be float64 or string
		var minVal uint64
		switch mv := minRaw.(type) {
		case float64:
			minVal = uint64(mv)
		case string:
			parsed, err := strconv.ParseUint(mv, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid numeric string for 'min': %w", err)
			}
			minVal = parsed
		default:
			return fmt.Errorf("unsupported type %T for 'min'", minRaw)
		}

		// extract value
		raw, ok := subj[fld]
		if !ok {
			return fmt.Errorf("credentialSubject missing '%s'", fld)
		}

		// value may be numeric or string
		var val uint64
		switch v := raw.(type) {
		case float64:
			val = uint64(v)
		case int:
			val = uint64(v)
		case int64:
			val = uint64(v)
		case uint64:
			val = v
		case string:
			parsed, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid numeric string for field '%s': %w", fld, err)
			}
			val = parsed
		default:
			return fmt.Errorf("unsupported type %T for field '%s'", raw, fld)
		}

		// generate proof
		rp, err := GenerateRangeProof(val, minVal)
		if err != nil {
			return fmt.Errorf("generating range proof: %w", err)
		}

		// marshal proof and append
		b, err := json.Marshal(rp)
		if err != nil {
			return fmt.Errorf("marshaling proof JSON: %w", err)
		}
		c.Proofs = append(c.Proofs, b)

		// remove raw field
		delete(c.CredentialSubject, fld)

	default:
		return fmt.Errorf("unsupported proof type '%s'", ch.Type)
	}
	return nil
}

// SignCredential signs the credential with Ed25519 and appends a signature proof.
func (c *Credential) SignCredential(priv ed25519.PrivateKey, verificationMethod string) error {
	// clone without proofs
	tmp := *c
	tmp.Proofs = nil

	// marshal for signing
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
