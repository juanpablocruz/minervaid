package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/juanpablocruz/minervaid/internal/config"
	"github.com/juanpablocruz/minervaid/internal/credentials"
	"github.com/spf13/cobra"
)

// revokeCmd marks a credential as revoked in the vault
var revokeCmd = &cobra.Command{
	Use:   "revoke <credID> [--out <vaultDir>]",
	Short: "Revoke a credential",
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
		// Initialize revocation list
		rl, err := credentials.NewRevocationList(filepath.Join(vaultDir, "revocations.json"))
		if err != nil {
			return fmt.Errorf("initialize revocation list: %w", err)
		}
		if err := rl.Revoke(credID); err != nil {
			return fmt.Errorf("revoke credential: %w", err)
		}
		cmd.Printf("Credential '%s' revoked\n", credID)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(revokeCmd)
	revokeCmd.Flags().StringVar(&altVaultDir, "out", "", "Vault directory (optional, uses active)")
}
