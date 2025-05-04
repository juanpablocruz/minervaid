package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/juanpablocruz/minervaid/internal/identity"
	"github.com/juanpablocruz/minervaid/internal/vault"
	"github.com/spf13/cobra"
)

var (
	name      string
	webDomain string
	vaultDir  string
)

var initCmd = &cobra.Command{
	Use:   "init --name <name> [--web domain.com] [--out <dir>]",
	Short: "Create a new identity",
	RunE: func(cmd *cobra.Command, args []string) error {
		pub, priv, err := identity.GenerateKeyPair()
		if err != nil {
			return fmt.Errorf("generate key pair: %w", err)
		}

		var did string
		if webDomain != "" {
			did = fmt.Sprintf("did:web:%s", webDomain)
		} else {
			did = identity.GenerateDID(pub)
		}

		didDoc, err := identity.BuildDIDDocument(did, pub)
		if err != nil {
			return fmt.Errorf("build did document: %w", err)
		}

		target := vaultDir
		if target == "" {
			target = filepath.Join("store", name)
		}

		v := vault.NewVault(target)
		if err := v.Init(didDoc, priv); err != nil {
			return fmt.Errorf("initialize vault: %w", err)
		}

		// 6. Success message
		cmd.Printf("Identity '%s' created at %s\n", did, target)
		return nil
	},
}

func init() {
	initCmd.Flags().StringVar(&name, "name", "", "Name of the identity (required)")
	initCmd.Flags().StringVar(&webDomain, "web", "", "Domain for did:web method (optional)")
	initCmd.Flags().StringVar(&vaultDir, "out", "", "Directory in which to create the identity (optional)")
	initCmd.MarkFlagRequired("name")
}
