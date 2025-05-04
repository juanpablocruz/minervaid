package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestCheckRevokedCommand(t *testing.T) {
	tmp := t.TempDir()
	// init vault
	rootCmd.SetArgs([]string{"init", "--name", "v2", "--out", tmp})
	if err := Execute(); err != nil {
		t.Fatal(err)
	}
	// prepare revocation list
	rlPath := filepath.Join(tmp, "revocations.json")
	os.WriteFile(rlPath, []byte(`["x1"]`), 0600)
	// check revoked
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"check-revoked", "x1", "--out", tmp})
	_ = Execute() // exit 0
	if !bytes.Contains(buf.Bytes(), []byte("is revoked")) {
		t.Errorf("expected revoked message, got %s", buf.String())
	}
	// check not revoked
	buf.Reset()
	rootCmd.SetArgs([]string{"check-revoked", "x2", "--out", tmp})
	_ = Execute()
	if !bytes.Contains(buf.Bytes(), []byte("is not revoked")) {
		t.Errorf("expected not revoked message, got %s", buf.String())
	}
}
