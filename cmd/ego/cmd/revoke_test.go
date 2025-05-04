package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestRevokeCommand(t *testing.T) {
	tmp := t.TempDir()
	// init vault
	rootCmd.SetArgs([]string{"init", "--name", "v", "--out", tmp})
	if err := Execute(); err != nil {
		t.Fatal(err)
	}
	// create a dummy VC file
	credDir := filepath.Join(tmp, "credentials")
	os.MkdirAll(credDir, 0700)
	dummy := []byte(`{"id":"c1"}`)
	os.WriteFile(filepath.Join(credDir, "c1.json"), dummy, 0600)
	// revoke it
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"revoke", "c1", "--out", tmp})
	if err := Execute(); err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("Credential 'c1' revoked")) {
		t.Errorf("unexpected output %s", buf.String())
	}
}
