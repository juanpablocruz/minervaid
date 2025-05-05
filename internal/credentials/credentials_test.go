package credentials

import (
	"crypto/ed25519"
	"encoding/json"
	"testing"
)

func TestCredentialAttachRangeProof(t *testing.T) {
	subj := map[string]interface{}{"id": "did:example:holder", "age": 25}
	cred := NewCredential("cred1", "did:example:issuer", subj)
	// Build generic range challenge JSON
	challenge := Challenge{
		Type:   "range",
		Params: map[string]interface{}{"field": "age", "min": 18},
	}
	challengeJSON, _ := json.Marshal(challenge)
	if err := cred.AttachProof(challengeJSON); err != nil {
		t.Fatalf("AttachProof failed: %v", err)
	}
	// age should be removed
	if _, ok := cred.CredentialSubject["age"]; ok {
		t.Error("age field should be removed after proof")
	}
	// one proof attached
	if len(cred.Proofs) != 1 {
		t.Errorf("expected 1 proof, got %d", len(cred.Proofs))
	}
	// Unmarshal into RangeProof type from zkp.go
	var rp RangeProof
	if err := json.Unmarshal(cred.Proofs[0], &rp); err != nil {
		t.Fatalf("Unmarshal RangeProof: %v", err)
	}
	if rp.Type != "BulletproofRangeProof" {
		t.Errorf("unexpected proof type: %s", rp.Type)
	}
}

func TestSignCredential(t *testing.T) {
	subj := map[string]interface{}{"id": "did:example:holder"}
	cred := NewCredential("cred2", "did:example:issuer", subj)
	_, priv, _ := ed25519.GenerateKey(nil)
	// append dummy proof
	cred.Proofs = append(cred.Proofs, json.RawMessage(`{"type":"Dummy"}`))

	if err := cred.SignCredential(priv, "did:example:issuer#key-1"); err != nil {
		t.Fatalf("SignCredential failed: %v", err)
	}
	if len(cred.Proofs) != 2 {
		t.Errorf("expected 2 proofs, got %d", len(cred.Proofs))
	}
	// Verify signature proof is last
	var sp SignatureProof
	if err := json.Unmarshal(cred.Proofs[1], &sp); err != nil {
		t.Fatalf("Unmarshal SignatureProof: %v", err)
	}
	if sp.Type != "Ed25519Signature2018" {
		t.Errorf("unexpected signature proof type: %s", sp.Type)
	}
}
