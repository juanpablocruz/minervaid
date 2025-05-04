package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/juanpablocruz/minervaid/internal/config"
	"github.com/juanpablocruz/minervaid/internal/credentials"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show DID and list of credentials for active identity",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		rootDir := cfg.RootDir
		if rootDir == "" {
			rootDir = "store"
		}
		vaultDir := filepath.Join(rootDir, cfg.Active)

		// Print DID document
		didData, err := os.ReadFile(filepath.Join(vaultDir, "did.json"))
		if err != nil {
			return fmt.Errorf("read did.json: %w", err)
		}
		cmd.Println(string(didData))

		// List credentials and their subjects
		credDir := filepath.Join(vaultDir, "credentials")
		files, _ := os.ReadDir(credDir)
		for _, f := range files {
			if filepath.Ext(f.Name()) == ".json" {
				data, _ := os.ReadFile(filepath.Join(credDir, f.Name()))
				var c credentials.Credential
				json.Unmarshal(data, &c)
				cmd.Printf("- %s: %v\n", c.ID, c.CredentialSubject)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
