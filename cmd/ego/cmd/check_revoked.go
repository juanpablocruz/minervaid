package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/juanpablocruz/minervaid/internal/config"
	"github.com/juanpablocruz/minervaid/internal/credentials"
	"github.com/spf13/cobra"
)

// checkRevokedCmd checks if a credential is revoked
var checkRevokedCmd = &cobra.Command{
	Use:   "check-revoked <credID> [--out <vaultDir>]",
	Short: "Check if a credential is revoked",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		credID := args[0]
		// Load config
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		rootDir := cfg.RootDir
		if rootDir == "" {
			rootDir = "store"
		}
		vaultDir := altVaultDir
		if vaultDir == "" {
			vaultDir = filepath.Join(rootDir, cfg.Active)
		}
		// Load revocation list
		rl, err := credentials.NewRevocationList(filepath.Join(vaultDir, "revocations.json"))
		if err != nil {
			return fmt.Errorf("initialize revocation list: %w", err)
		}
		if rl.IsRevoked(credID) {
			cmd.Printf("Credential '%s' is revoked", credID)
			return nil
		} else {
			cmd.Printf("Credential '%s' is not revoked", credID)
			return nil
		}
	},
}

func init() {
	rootCmd.AddCommand(checkRevokedCmd)
	checkRevokedCmd.Flags().StringVar(&altVaultDir, "out", "", "Vault directory (optional, uses active)")
}
