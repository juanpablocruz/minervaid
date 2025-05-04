package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestIssueCommand_NewAttributesCredential(t *testing.T) {
	// 1. Setup temp vault
	tmpDir := t.TempDir()
	// Initialize vault
	rootCmd.SetArgs([]string{"init", "--name", "testvault", "--out", tmpDir})
	if err := Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// 2. Write attributes.json
	attrs := map[string]interface{}{"age": "32", "name": "Alice"}
	attrBytes, _ := json.Marshal(attrs)
	if err := os.WriteFile(filepath.Join(tmpDir, "attributes.json"), attrBytes, 0600); err != nil {
		t.Fatalf("write attributes.json failed: %v", err)
	}

	// 3. Run issue command with altVaultDir and custom ID
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"issue", "--out", tmpDir, "--id", "cred123"})
	if err := Execute(); err != nil {
		t.Fatalf("issue failed: %v", err)
	}

	// 4. Assert output
	out := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte("Credential 'cred123' issued")) {
		t.Errorf("unexpected output: %s", out)
	}

	// 5. Check credential file
	credFile := filepath.Join(tmpDir, "credentials", "cred123.json")
	data, err := os.ReadFile(credFile)
	if err != nil {
		t.Fatalf("credential file not created: %v", err)
	}

	// 6. Validate credential content
	var cred map[string]interface{}
	if err := json.Unmarshal(data, &cred); err != nil {
		t.Fatalf("invalid JSON in credential file: %v", err)
	}
	// Check subject equals attrs
	subj, ok := cred["credentialSubject"].(map[string]interface{})
	if !ok {
		t.Fatalf("credentialSubject missing or wrong type: %v", cred["credentialSubject"])
	}
	for k, v := range attrs {
		if subj[k] != v {
			t.Errorf("expected subject[%s] = %v, got %v", k, v, subj[k])
		}
	}
}
