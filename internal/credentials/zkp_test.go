package credentials

import (
	"encoding/json"
	"testing"
)

func TestGenerateRangeProofSuccess(t *testing.T) {
	value := uint64(30)
	min := uint64(18)
	rp, err := GenerateRangeProof(value, min)
	if err != nil {
		t.Fatalf("GenerateRangeProof failed: %v", err)
	}
	// proof metadata check
	if rp.Min != min {
		t.Errorf("expected Min=%d, got %d", min, rp.Min)
	}
	if rp.Range == 0 {
		t.Error("expected non-zero Range")
	}
	if rp.Type != "BulletproofRangeProof" {
		t.Errorf("unexpected Type: %s", rp.Type)
	}
}

func TestGenerateRangeProofBelowMin(t *testing.T) {
	_, err := GenerateRangeProof(10, 18)
	if err == nil {
		t.Error("expected error when value < min, got nil")
	}
}

func TestVerifyRangeProof(t *testing.T) {
	value := uint64(25)
	min := uint64(18)
	rp, err := GenerateRangeProof(value, min)
	if err != nil {
		t.Fatalf("GenerateRangeProof failed: %v", err)
	}
	// marshal/unmarshal to simulate JSON transit
	data, err := json.Marshal(rp)
	if err != nil {
		t.Fatalf("failed to marshal RangeProof: %v", err)
	}
	var rp2 RangeProof
	if err := json.Unmarshal(data, &rp2); err != nil {
		t.Fatalf("failed to unmarshal RangeProof: %v", err)
	}
	// now verify
	if err := VerifyRangeProof(&rp2); err != nil {
		t.Errorf("VerifyRangeProof failed: %v", err)
	}
}

func TestVerifyRangeProofInvalid(t *testing.T) {
	value, min := uint64(20), uint64(18)
	rp, err := GenerateRangeProof(value, min)
	if err != nil {
		t.Fatalf("GenerateRangeProof failed: %v", err)
	}
	// Actually corrupt the proof
	rp.Proof.V = nil
	if err := VerifyRangeProof(rp); err == nil {
		t.Error("expected VerifyRangeProof to fail on corrupted proof")
	}
}
