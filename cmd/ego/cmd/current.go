package cmd

import (
	"github.com/juanpablocruz/minervaid/internal/config"
	"github.com/spf13/cobra"
)

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show the name of the active identity",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if cfg.Active == "" {
			cmd.Println("No active identity. Use 'ego use <name>' or run 'ego init'")
		} else {
			cmd.Printf("Active identity: %s\n", cfg.Active)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(currentCmd)
}
