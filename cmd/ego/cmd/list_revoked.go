package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/juanpablocruz/minervaid/internal/config"
	"github.com/juanpablocruz/minervaid/internal/credentials"
	"github.com/spf13/cobra"
)

// listRevokedCmd lists all revoked credential IDs
var listRevokedCmd = &cobra.Command{
	Use:   "list-revoked [--out <vaultDir>]",
	Short: "List all revoked credentials",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
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
		for _, id := range rl.List() {
			cmd.Println(id)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listRevokedCmd)
	listRevokedCmd.Flags().StringVar(&altVaultDir, "out", "", "Vault directory (optional, uses active)")
}
