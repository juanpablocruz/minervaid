package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestListRevokedCommand(t *testing.T) {
	tmp := t.TempDir()
	// init vault
	rootCmd.SetArgs([]string{"init", "--name", "v3", "--out", tmp})
	if err := Execute(); err != nil {
		t.Fatal(err)
	}
	// prepare revocation list
	rl := `[
	  "a1",
	  "b2"
	]`
	os.WriteFile(filepath.Join(tmp, "revocations.json"), []byte(rl), 0600)
	// list-revoked
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"list-revoked", "--out", tmp})
	if err := Execute(); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte("a1")) || !bytes.Contains(buf.Bytes(), []byte("b2")) {
		t.Errorf("unexpected output %s", out)
	}
}
