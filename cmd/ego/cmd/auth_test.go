package cmd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAuthRequestCommand(t *testing.T) {
	// Prepare temp presentation_definition file
	tmp := t.TempDir()
	pd := map[string]interface{}{"id": "test_pd"}
	pdBytes, _ := json.Marshal(pd)
	pdFile := filepath.Join(tmp, "pd.json")
	if err := os.WriteFile(pdFile, pdBytes, 0600); err != nil {
		t.Fatalf("write presentation_definition: %v", err)
	}

	// Execute ego auth-request
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	execArgs := []string{
		"auth-request",
		"--endpoint", "https://as.test/auth",
		"--client-id", "cid123",
		"--redirect-uri", "http://localhost:1111/callback",
		"--presentation-definition", pdFile,
		"--nonce", "n1",
	}
	rootCmd.SetArgs(execArgs)
	if err := Execute(); err != nil {
		t.Fatalf("auth-request failed: %v", err)
	}
	out := strings.TrimSpace(buf.String())

	u, err := url.Parse(out)
	if err != nil {
		t.Fatalf("invalid URL output: %v", err)
	}
	q := u.Query()
	if got := q.Get("client_id"); got != "cid123" {
		t.Errorf("expected client_id=cid123, got %s", got)
	}
	if got := q.Get("nonce"); got != "n1" {
		t.Errorf("expected nonce=n1, got %s", got)
	}
	// Check presentation_definition round-trip
	gotPD := q.Get("presentation_definition")
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(gotPD), &parsed); err != nil {
		t.Errorf("invalid pd JSON in URL: %v", err)
	}
	if parsed["id"] != "test_pd" {
		t.Errorf("presentation_definition content mismatch: %v", parsed)
	}
}

func TestAuthVerifyCommand(t *testing.T) {
	// 1. Initialize vault
	tmp := t.TempDir()
	rootCmd.SetArgs([]string{"init", "--name", "vault", "--out", tmp})
	if err := Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}
	// 2. Create a credential and presentation
	rootCmd.SetArgs([]string{"set", "age", "42", "--out", tmp})
	Execute()
	rootCmd.SetArgs([]string{"issue", "--out", tmp, "--id", "vc1"})
	Execute()
	rootCmd.SetArgs([]string{"present", "--out", tmp})
	Execute()
	// 3. Load the generated presentation
	presDir := filepath.Join(tmp, "presentations")
	files, err := os.ReadDir(presDir)
	if err != nil || len(files) == 0 {
		t.Fatalf("no presentation generated: %v", err)
	}
	presBytes, err := os.ReadFile(filepath.Join(presDir, files[0].Name()))
	if err != nil {
		t.Fatalf("read presentation failed: %v", err)
	}
	var pres interface{}
	if err := json.Unmarshal(presBytes, &pres); err != nil {
		t.Fatalf("invalid presentation JSON: %v", err)
	}

	// 4. Build a dummy JWT with alg=none, embed VP and iss
	payload := map[string]interface{}{
		"iss": "did:key:zDUMMY",
		"vp":  pres,
	}
	plBytes, _ := json.Marshal(payload)
	b64pl := base64.RawURLEncoding.EncodeToString(plBytes)
	head := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none"}`))
	jwt := fmt.Sprintf("%s.%s.", head, b64pl)

	jwtFile := filepath.Join(tmp, "token.jwt")
	if err := os.WriteFile(jwtFile, []byte(jwt), 0600); err != nil {
		t.Fatalf("write id_token file: %v", err)
	}

	// 5. Run auth-verify
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"auth-verify", "--id-token", jwtFile})
	if err := Execute(); err != nil {
		t.Fatalf("auth-verify failed: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Authentication successful for DID:") {
		t.Errorf("unexpected output: %s", out)
	}
}
