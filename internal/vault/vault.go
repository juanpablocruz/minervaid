package vault

import (
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/juanpablocruz/minervaid/internal/identity"
	"github.com/mr-tron/base58"
)

const (
	didFilename      = "did.json"
	keystoreFilename = "keystore.json"
)

type Vault struct {
	BaseDir string
}

// NewVault prepares a Vault rooted at baseDir (does not create files yet)
func NewVault(baseDir string) *Vault {
	return &Vault{BaseDir: baseDir}
}

// Init writes a new DID Document and encodes the private key in base58, saving both to disk.
func (v *Vault) Init(didDoc []byte, privKey ed25519.PrivateKey) error {
	if err := os.MkdirAll(v.BaseDir, fs.ModePerm); err != nil {
		return fmt.Errorf("create vault dir: %w", err)
	}

	// Write DID document
	if err := os.WriteFile(
		filepath.Join(v.BaseDir, didFilename),
		didDoc,
		0600,
	); err != nil {
		return fmt.Errorf("write did.json: %w", err)
	}

	// Encode and write private key
	enc := identity.EncodePrivateKey(privKey)
	k := map[string]string{"privateKey": enc}
	ksBytes, err := json.MarshalIndent(k, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal keystore: %w", err)
	}
	if err := os.WriteFile(
		filepath.Join(v.BaseDir, keystoreFilename),
		ksBytes,
		0600,
	); err != nil {
		return fmt.Errorf("write keystore.json: %w", err)
	}

	return nil
}

// Load reads did.json and keystore.json, decodes the private key from base58, and returns the contents.
func (v *Vault) Load() ([]byte, ed25519.PrivateKey, error) {
	// Read DID document
	didDoc, err := os.ReadFile(filepath.Join(v.BaseDir, didFilename))
	if err != nil {
		return nil, nil, fmt.Errorf("read did.json: %w", err)
	}

	// Read keystore entry
	ksBytes, err := os.ReadFile(filepath.Join(v.BaseDir, keystoreFilename))
	if err != nil {
		return nil, nil, fmt.Errorf("read keystore.json: %w", err)
	}
	var k map[string]string
	if err := json.Unmarshal(ksBytes, &k); err != nil {
		return nil, nil, fmt.Errorf("unmarshal keystore: %w", err)
	}
	enc, ok := k["privateKey"]
	if !ok {
		return nil, nil, fmt.Errorf("keystore missing 'privateKey'")
	}

	// Decode private key
	rawPriv, err := base58.Decode(enc)
	if err != nil {
		return nil, nil, fmt.Errorf("decode private key: %w", err)
	}

	return didDoc, ed25519.PrivateKey(rawPriv), nil
}
