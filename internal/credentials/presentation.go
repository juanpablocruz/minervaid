package credentials

import (
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"time"
)

type Presentation struct {
	Context              []string          `json:"@context"`
	Type                 []string          `json:"type"`
	VerifiableCredential []Credential      `json:"verifiableCredential"`
	Holder               string            `json:"holder,omitempty"`
	Proofs               []json.RawMessage `json:"proof"`
}

func NewPresentation(creds []Credential, holder string) *Presentation {
	return &Presentation{
		Context:              []string{"https://www.w3.org/2018/credentials/v1"},
		Type:                 []string{"VerifiablePresentation"},
		VerifiableCredential: creds,
		Holder:               holder,
		Proofs:               []json.RawMessage{},
	}
}

func (p *Presentation) SignPresentation(priv ed25519.PrivateKey, verificationMethod string) error {
	// Serialize presentation without existing signature proofs
	tmp := *p
	tmp.Proofs = nil
	data, err := json.Marshal(tmp)
	if err != nil {
		return err
	}
	sig := ed25519.Sign(priv, data)
	sp := SignatureProof{
		Type:               "Ed25519Signature2018",
		Created:            time.Now().UTC().Format(time.RFC3339),
		ProofPurpose:       "authentication",
		VerificationMethod: verificationMethod,
		JWS:                fmt.Sprintf("%x", sig),
	}
	proofBytes, err := json.Marshal(sp)
	if err != nil {
		return fmt.Errorf("marshaling signature proof: %w", err)
	}
	p.Proofs = append(p.Proofs, proofBytes)
	return nil
}

// ToJSON returns the formatted JSON of the presentation
func (p *Presentation) ToJSON() ([]byte, error) {
	return json.MarshalIndent(p, "", "  ")
}
