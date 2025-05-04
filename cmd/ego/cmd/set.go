package cmd

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/juanpablocruz/minervaid/internal/config"
	"github.com/juanpablocruz/minervaid/internal/credentials"
	"github.com/juanpablocruz/minervaid/internal/vault"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Issue a credential for all stored attributes",
	Long:  "Store/update a metadata attribute and issue a new Verifiable Credential with all attributes for the active identity.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key, val := args[0], args[1]

		// Load CLI config
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		rootDir := cfg.RootDir
		if rootDir == "" {
			rootDir = "store"
		}
		vaultDir := filepath.Join(rootDir, cfg.Active)

		// Load existing attributes
		attrFile := filepath.Join(vaultDir, "attributes.json")
		attrs := make(map[string]interface{})
		if data, err := os.ReadFile(attrFile); err == nil {
			err = json.Unmarshal(data, &attrs)
			if err != nil {
				return fmt.Errorf("invalid attributes.json: %w", err)
			}
		}

		// Update attribute map
		attrs[key] = val

		// Save updated attributes.json
		if err := os.MkdirAll(vaultDir, fs.FileMode(0700)); err != nil {
			return fmt.Errorf("create vault dir: %w", err)
		}
		outData, err := json.MarshalIndent(attrs, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal attributes: %w", err)
		}
		if err := os.WriteFile(attrFile, outData, fs.FileMode(0600)); err != nil {
			return fmt.Errorf("write attributes.json: %w", err)
		}

		// Load identity keys
		v := vault.NewVault(vaultDir)
		didDoc, priv, err := v.Load()
		if err != nil {
			return fmt.Errorf("load vault: %w", err)
		}

		// Extract DID
		var doc map[string]interface{}
		if err := json.Unmarshal(didDoc, &doc); err != nil {
			return fmt.Errorf("parse did.json: %w", err)
		}
		did, _ := doc["id"].(string)

		// Issue credential with full attributes
		credID := time.Now().UTC().Format("20060102T150405Z")
		cred := credentials.NewCredential(credID, did, attrs)
		if err := cred.SignCredential(priv, did+"#keys-1"); err != nil {
			return fmt.Errorf("sign credential: %w", err)
		}

		// Save credential
		credDir := filepath.Join(vaultDir, "credentials")
		if err := os.MkdirAll(credDir, fs.FileMode(0700)); err != nil {
			return fmt.Errorf("create credentials dir: %w", err)
		}
		store := &credentials.FileStore{Dir: credDir}
		if err := store.Save(cred); err != nil {
			return fmt.Errorf("save credential: %w", err)
		}

		cmd.Printf("Credential '%s' issued with attributes: %v", credID, attrs)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
