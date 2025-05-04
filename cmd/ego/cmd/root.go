package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "ego",
	Short: "Ego - CLI for SSI",
	Long:  "Ego is a CLI tool to manage SSI.",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(issueCmd)
	rootCmd.AddCommand(verifyCmd)
}
