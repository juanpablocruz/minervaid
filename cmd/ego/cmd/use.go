package cmd

import (
	"github.com/juanpablocruz/minervaid/internal/config"
	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:   "use <vaultName> [<rootDir>]",
	Short: "Select an identity vault to use",
	Long:  "Set the active identity vault by name (and optional root directory)",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		cfg.Active = args[0]
		if len(args) == 2 {
			cfg.RootDir = args[1]
		}
		if err := cfg.Save(); err != nil {
			return err
		}
		cmd.Printf("Active identity set to '%s'\n", cfg.Active)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(useCmd)
}
