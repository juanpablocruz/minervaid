package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/juanpablocruz/minervaid/internal/credentials"
	"github.com/spf13/cobra"
)

var (
	revokeID string
)

var revokeCredCmd = &cobra.Command{
	Use:   "revoke-cred",
	Short: "Revoke a credential by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		if revokeID == "" {
			return fmt.Errorf("--id is required")
		}
		if err := os.MkdirAll(storeDir, 0755); err != nil {
			return err
		}
		path := filepath.Join(storeDir, "revocations.json")
		rl, err := credentials.NewRevocationList(path)
		if err != nil {
			return err
		}
		if err := rl.Revoke(revokeID); err != nil {
			return err
		}
		fmt.Printf("Credential %s revoked", revokeID)
		return nil
	},
}

var listRevokedCmd = &cobra.Command{
	Use:   "list-revoked",
	Short: "List all revoked credential IDs",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := filepath.Join(storeDir, "revocations.json")
		rl, err := credentials.NewRevocationList(path)
		if err != nil {
			return err
		}
		for _, id := range rl.List() {
			fmt.Println(id)
		}
		return nil
	},
}

var checkRevokedCmd = &cobra.Command{
	Use:   "check-revoked",
	Short: "Check if a credential is revoked",
	RunE: func(cmd *cobra.Command, args []string) error {
		if revokeID == "" {
			return fmt.Errorf("--id is required")
		}
		path := filepath.Join(storeDir, "revocations.json")
		rl, err := credentials.NewRevocationList(path)
		if err != nil {
			return err
		}
		if rl.IsRevoked(revokeID) {
			fmt.Printf("Credential %s is revoked", revokeID)
		} else {
			fmt.Printf("Credential %s is not revoked", revokeID)
		}
		return nil
	},
}

func init() {
	revokeCredCmd.Flags().StringVar(&revokeID, "id", "", "Credential ID to revoke or check")
	listRevokedCmd.Flags().StringVar(&storeDir, "store", storeDir, "Directory for storing data")
	checkRevokedCmd.Flags().StringVar(&revokeID, "id", "", "Credential ID to check")
	rootCmd.AddCommand(revokeCredCmd, listRevokedCmd, checkRevokedCmd)
}
