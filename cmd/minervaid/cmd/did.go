package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/juanpablocruz/minervaid/internal/identity"
	"github.com/spf13/cobra"
)

var newDidCmd = &cobra.Command{
	Use:   "new-did",
	Short: "Generate a new DID",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := os.MkdirAll(storeDir, 0755); err != nil {
			return err
		}
		pub, priv, err := identity.GenerateKeyPair()
		if err != nil {
			return err
		}
		did := identity.GenerateDID(pub)
		privEnc := identity.EncodePrivateKey(priv)
		ksPath := filepath.Join(storeDir, "keystore.json")
		ks := loadKeyStore(ksPath)
		ks[did] = privEnc
		saveKeyStore(ksPath, ks)
		fmt.Println(did)
		return nil
	},
}

var listDidsCmd = &cobra.Command{
	Use:   "list-dids",
	Short: "List all DIDs",
	RunE: func(cmd *cobra.Command, args []string) error {
		ksPath := filepath.Join(storeDir, "keystore.json")
		ks := loadKeyStore(ksPath)
		for did := range ks {
			fmt.Println(did)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(newDidCmd, listDidsCmd)
}
