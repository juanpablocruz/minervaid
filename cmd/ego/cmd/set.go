package cmd

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/juanpablocruz/minervaid/internal/config"
	"github.com/spf13/cobra"
)

// setCmd updates the attributes.json for the active identity vault
// It does not issue a credential; use `ego issue` to generate one from attributes.
var setCmd = &cobra.Command{
	Use:   "set <key> <value> [--out <directory>]",
	Short: "Set a metadata attribute for the active identity",
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

		// Determine vault directory
		vaultDir := altVaultDir
		if vaultDir == "" {
			vaultDir = filepath.Join(rootDir, cfg.Active)
		}

		// Ensure vault exists
		if _, err := os.Stat(filepath.Join(vaultDir, "did.json")); err != nil {
			return fmt.Errorf("vault '%s' not found: %w", cfg.Active, err)
		}

		// Load existing attributes
		attrFile := filepath.Join(vaultDir, "attributes.json")
		attrs := make(map[string]interface{})
		if data, err := os.ReadFile(attrFile); err == nil {
			if err := json.Unmarshal(data, &attrs); err != nil {
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

		cmd.Printf("Attribute '%s' set to '%s'\n", key, val)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
	setCmd.Flags().StringVar(&altVaultDir, "out", "", "Directory of the vault (optional, overrides active)")
}
