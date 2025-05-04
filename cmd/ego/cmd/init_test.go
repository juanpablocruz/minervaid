package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitCommand(t *testing.T) {
	// 1. Prepare a temp directory
	tmp := t.TempDir()

	// 2. Capture stdout/stderr
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	// 3. Run: ego init --name testid --out <tmp>
	rootCmd.SetArgs([]string{"init", "--name", "testid", "--out", tmp})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	out := buf.String()
	if !strings.HasPrefix(out, "Identity '") {
		t.Errorf("unexpected output:\n%s", out)
	}

	// 4. Check did.json
	didPath := filepath.Join(tmp, "did.json")
	data, err := os.ReadFile(didPath)
	if err != nil {
		t.Fatalf("did.json not created: %v", err)
	}
	var doc map[string]interface{}
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("invalid JSON in did.json: %v", err)
	}
	idVal, ok := doc["id"].(string)
	if !ok {
		t.Fatalf("did.json missing \"id\" field")
	}
	if !strings.HasPrefix(idVal, "did:key:") && !strings.HasPrefix(idVal, "did:web:") {
		t.Errorf("unexpected DID in did.json: %s", idVal)
	}

	// 5. Check keystore.json
	ksPath := filepath.Join(tmp, "keystore.json")
	data2, err := os.ReadFile(ksPath)
	if err != nil {
		t.Fatalf("keystore.json not created: %v", err)
	}
	var ks map[string]string
	if err := json.Unmarshal(data2, &ks); err != nil {
		t.Fatalf("invalid JSON in keystore.json: %v", err)
	}
	pk, ok := ks["privateKey"]
	if !ok {
		t.Fatalf("keystore.json missing \"privateKey\" field")
	}
	if len(pk) == 0 {
		t.Error("privateKey is empty")
	}
}
