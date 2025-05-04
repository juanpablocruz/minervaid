package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/juanpablocruz/minervaid/internal/config"
	"github.com/spf13/cobra"
)

// Usage: ego list [<vaultName> [<rootDir>]]

var listCmd = &cobra.Command{
	Use:   "list [<vaultName> [<rootDir>]]",
	Short: "List all issued credentials in a vault",
	Long: `List IDs of all credentials stored in a vault.

With no args, uses the active vault from ~/.ego/config.json.
Otherwise: list <vaultName> [<rootDir>] looks under <rootDir>/<vaultName>/credentials (default rootDir="store").`,
	Args: cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Determine vaultName and rootDir
		var vaultName, rootDir string
		switch len(args) {
		case 0:
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			vaultName = cfg.Active
			rootDir = cfg.RootDir
		case 1:
			vaultName = args[0]
		case 2:
			vaultName = args[0]
			rootDir = args[1]
		}
		if vaultName == "" {
			return fmt.Errorf("no vault specified and no active vault; use 'ego use' or pass vaultName")
		}
		if rootDir == "" {
			rootDir = "store"
		}
		target := filepath.Join(rootDir, vaultName)

		// Read credential files
		credDir := filepath.Join(target, "credentials")
		files, err := os.ReadDir(credDir)
		if err != nil {
			if os.IsNotExist(err) {
				return nil // no credentials yet
			}
			return fmt.Errorf("reading credentials: %w", err)
		}
		for _, f := range files {
			if !f.IsDir() && filepath.Ext(f.Name()) == ".json" {
				cmd.Println(f.Name()[:len(f.Name())-5])
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
