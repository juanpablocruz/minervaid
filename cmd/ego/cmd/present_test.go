package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestPresentCommand_AllCredentials(t *testing.T) {
	// Setup temp vault
	tmpDir := t.TempDir()
	// Initialize vault
	rootCmd.SetArgs([]string{"init", "--name", "testvault", "--out", tmpDir})
	if err := Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}
	// Set two attributes
	rootCmd.SetArgs([]string{"set", "age", "30", "--out", tmpDir})
	if err := Execute(); err != nil {
		t.Fatalf("set age failed: %v", err)
	}
	rootCmd.SetArgs([]string{"set", "name", "Bob", "--out", tmpDir})
	if err := Execute(); err != nil {
		t.Fatalf("set name failed: %v", err)
	}
	// Issue a credential snapshot
	rootCmd.SetArgs([]string{"issue", "--out", tmpDir, "--id", "vc1"})
	if err := Execute(); err != nil {
		t.Fatalf("issue vc1 failed: %v", err)
	}
	// Modify attribute and issue another
	rootCmd.SetArgs([]string{"set", "email", "bob@example.com", "--out", tmpDir})
	if err := Execute(); err != nil {
		t.Fatalf("set email failed: %v", err)
	}
	rootCmd.SetArgs([]string{"issue", "--out", tmpDir, "--id", "vc2"})
	if err := Execute(); err != nil {
		t.Fatalf("issue vc2 failed: %v", err)
	}

	// Run present with all credentials
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"present", "--out", tmpDir})
	if err := Execute(); err != nil {
		t.Fatalf("present failed: %v", err)
	}
	// Check output message
	out := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte("Presentation '")) {
		t.Errorf("unexpected output: %s", out)
	}

	// Find the generated presentation file
	presDir := filepath.Join(tmpDir, "presentations")
	files, err := os.ReadDir(presDir)
	if err != nil || len(files) != 1 {
		t.Fatalf("expected one presentation file, got %v, err=%v", files, err)
	}
	presFile := filepath.Join(presDir, files[0].Name())

	// Load and parse presentation JSON
	data, err := os.ReadFile(presFile)
	if err != nil {
		t.Fatalf("read presentation file failed: %v", err)
	}
	var pres map[string]interface{}
	if err := json.Unmarshal(data, &pres); err != nil {
		t.Fatalf("invalid presentation JSON: %v", err)
	}

	// Verify verifiableCredential field contains two entries
	vcs, ok := pres["verifiableCredential"].([]interface{})
	if !ok {
		t.Fatalf("verifiableCredential missing or wrong type: %T", pres["verifiableCredential"])
	}
	if len(vcs) != 2 {
		t.Errorf("expected 2 credentials in presentation, got %d", len(vcs))
	}
}

func TestPresentCommand_SelectiveDisclosure(t *testing.T) {
	// Setup temp vault
	tmpDir := t.TempDir()
	// Initialize vault
	rootCmd.SetArgs([]string{"init", "--name", "demo", "--out", tmpDir})
	if err := Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}
	// Set attributes
	rootCmd.SetArgs([]string{"set", "age", "45", "--out", tmpDir})
	if err := Execute(); err != nil {
		t.Fatalf("set failed: %v", err)
	}
	rootCmd.SetArgs([]string{"set", "email", "test@x.com", "--out", tmpDir})
	if err := Execute(); err != nil {
		t.Fatalf("set failed: %v", err)
	}
	// Issue initial VC
	rootCmd.SetArgs([]string{"issue", "--out", tmpDir, "--id", "vcA"})
	if err := Execute(); err != nil {
		t.Fatalf("issue failed: %v", err)
	}

	// Run present with selective reveal (only email)
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"present", "--out", tmpDir, "--reveal", "email"})
	if err := Execute(); err != nil {
		t.Fatalf("present selective failed: %v", err)
	}

	// Load the presentation file
	files, _ := os.ReadDir(filepath.Join(tmpDir, "presentations"))
	data, _ := os.ReadFile(filepath.Join(tmpDir, "presentations", files[0].Name()))
	var pres map[string]interface{}
	json.Unmarshal(data, &pres)
	vcs := pres["verifiableCredential"].([]interface{})
	// Only one VC subject, with only email field
	subj := vcs[0].(map[string]interface{})["credentialSubject"].(map[string]interface{})
	if len(subj) != 1 || subj["email"] != "test@x.com" {
		t.Errorf("selective disclosure wrong subject: %v", subj)
	}
}
