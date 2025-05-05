package cmd

import (
	"fmt"
	"os"

	"github.com/juanpablocruz/minervaid/internal/store"
	"github.com/spf13/cobra"
)

func loadKeyStore(path string) map[string]string {
	return store.LoadKeyStore(path)
}

func saveKeyStore(path string, ks map[string]string) {
	store.SaveKeyStore(path, ks)
}

var (
	storeDir string
	rootCmd  = &cobra.Command{
		Use:   "minervaid",
		Short: "MinervaID CLI",
		Long:  "Modular SSI CLI for DIDs, VCs, Presentations, Revocation and Authentication",
	}
)

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&storeDir, "store", "./store", "Directory for storing data")
}
