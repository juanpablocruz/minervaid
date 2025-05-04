package credentials

import "testing"

func TestAgeProofFlow(t *testing.T) {
	proof, err := GenerateAgeProof(25, 18)
	if err != nil {
		t.Fatalf("GenerateAgeProof failed: %v", err)
	}
	if err := VerifyAgeProof(proof); err != nil {
		t.Errorf("VerifyAgeProof failed: %v", err)
	}
	if _, err := GenerateAgeProof(16, 18); err == nil {
		t.Error("expected error for age<minAge")
	}
}
