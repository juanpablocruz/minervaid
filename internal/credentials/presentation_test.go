package credentials

import (
	"crypto/ed25519"
	"encoding/json"
	"testing"
)

func TestPresentationSign(t *testing.T) {
	holder := "did:example:holder"
	// Prepare a dummy credential
	subj := map[string]interface{}{"id": holder}
	cred := *NewCredential("cred1", "did:example:issuer", subj)
	// Sign the credential so VC contains a signature proof
	_, privIssuer, _ := ed25519.GenerateKey(nil)
	cred.SignCredential(privIssuer, "did:example:issuer#key-1")

	pres := NewPresentation([]Credential{cred}, holder)
	_, privHolder, _ := ed25519.GenerateKey(nil)
	if err := pres.SignPresentation(privHolder, holder+"#key-1"); err != nil {
		t.Fatalf("SignPresentation failed: %v", err)
	}
	if len(pres.Proofs) != 1 {
		t.Fatalf("expected 1 proof, got %d", len(pres.Proofs))
	}
	var sp SignatureProof
	if err := json.Unmarshal(pres.Proofs[0], &sp); err != nil {
		t.Fatalf("unmarshal signature proof: %v", err)
	}
	if sp.Type != "Ed25519Signature2018" {
		t.Errorf("unexpected proof type: %s", sp.Type)
	}
}
