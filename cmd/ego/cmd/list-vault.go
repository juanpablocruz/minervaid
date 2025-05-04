package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var listVaultCmd = &cobra.Command{
	Use:   "list-vaults [<rootDir>]",
	Short: "List all vaults",
	Long:  "List all vault directories in the given root directory (default 'store').",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		rootDir := "store"
		if len(args) == 1 {
			rootDir = args[0]
		}
		entries, err := os.ReadDir(rootDir)
		if err != nil {
			return fmt.Errorf("read root directory: %w", err)
		}
		for _, e := range entries {
			if e.IsDir() {
				// check did.json inside
				if _, err := os.Stat(filepath.Join(rootDir, e.Name(), "did.json")); err == nil {
					cmd.Println(e.Name())
				}
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listVaultCmd)
}
