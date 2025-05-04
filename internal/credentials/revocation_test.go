package credentials

import (
	"os"
	"testing"
)

func TestRevocationList(t *testing.T) {
	tmp, err := os.CreateTemp("", "revoke*.json")
	if err != nil {
		t.Fatal(err)
	}
	path := tmp.Name()
	tmp.Close()
	defer os.Remove(path)

	rl, err := NewRevocationList(path)
	if err != nil {
		t.Fatalf("Init revocation list: %v", err)
	}
	if rl.IsRevoked("foo") {
		t.Error("foo should not be revoked yet")
	}
	if err := rl.Revoke("foo"); err != nil {
		t.Fatalf("Failed to revoke foo: %v", err)
	}
	if !rl.IsRevoked("foo") {
		t.Error("foo should be revoked")
	}
	ids := rl.List()
	if len(ids) != 1 || ids[0] != "foo" {
		t.Errorf("unexpected list: %v", ids)
	}
	// revoking again should error
	if err := rl.Revoke("foo"); err == nil {
		t.Error("expected error on double revoke")
	}
}
