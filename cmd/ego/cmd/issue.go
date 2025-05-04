package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/juanpablocruz/minervaid/internal/config"
	"github.com/juanpablocruz/minervaid/internal/credentials"
	"github.com/juanpablocruz/minervaid/internal/vault"
	"github.com/spf13/cobra"
)

var (
	credID      string
	altVaultDir string
)

// issueCmd issues a new Verifiable Credential based on stored attributes.json
var issueCmd = &cobra.Command{
	Use:   "issue [--id <id>] [--out <directory>]",
	Short: "Issue a new verifiable credential for current attributes",
	Long: `Issue a new Verifiable Credential whose subject is the full set of attributes
stored in attributes.json for the active identity vault.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load CLI config
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		rootDir := cfg.RootDir
		if rootDir == "" {
			rootDir = "store"
		}

		// Determine vault directory
		vaultDir := altVaultDir
		if vaultDir == "" {
			vaultDir = filepath.Join(rootDir, cfg.Active)
		}
		// Read attributes.json
		attrPath := filepath.Join(vaultDir, "attributes.json")
		data, err := os.ReadFile(attrPath)
		if err != nil {
			return fmt.Errorf("read attributes.json: %w", err)
		}
		var attrs map[string]interface{}
		if err := json.Unmarshal(data, &attrs); err != nil {
			return fmt.Errorf("invalid attributes.json: %w", err)
		}

		// Determine credential ID
		id := credID
		if id == "" {
			id = time.Now().UTC().Format("20060102T150405Z")
		}

		// Load DID and private key
		v := vault.NewVault(vaultDir)
		didDoc, priv, err := v.Load()
		if err != nil {
			return fmt.Errorf("load vault: %w", err)
		}

		// Parse DID
		var doc map[string]interface{}
		if err := json.Unmarshal(didDoc, &doc); err != nil {
			return fmt.Errorf("parse did.json: %w", err)
		}
		did, ok := doc["id"].(string)
		if !ok {
			return fmt.Errorf("did.json missing 'id'")
		}

		// Create and sign credential
		cred := credentials.NewCredential(id, did, attrs)
		if err := cred.SignCredential(priv, did+"#keys-1"); err != nil {
			return fmt.Errorf("sign credential: %w", err)
		}

		// Save credential
		credDir := filepath.Join(vaultDir, "credentials")
		if err := os.MkdirAll(credDir, 0700); err != nil {
			return fmt.Errorf("create credentials dir: %w", err)
		}
		store := &credentials.FileStore{Dir: credDir}
		if err := store.Save(cred); err != nil {
			return fmt.Errorf("save credential: %w", err)
		}

		cmd.Printf("Credential '%s' issued\n", id)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(issueCmd)
	issueCmd.Flags().StringVar(&credID, "id", "", "Credential ID (optional)")
	issueCmd.Flags().StringVar(&altVaultDir, "out", "", "Directory of the vault (optional, uses active if unset)")
}
