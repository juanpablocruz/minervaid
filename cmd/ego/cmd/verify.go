package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/juanpablocruz/minervaid/internal/credentials"
	"github.com/spf13/cobra"
)

var credentialFile string

var verifyCmd = &cobra.Command{
	Use:   "verify --file <credential.json>",
	Short: "Verify a verifiable credential",
	Long:  "Verify the signature of a credential JSON file.",
	RunE: func(cmd *cobra.Command, args []string) error {
		file := credentialFile
		if file == "" {
			return fmt.Errorf("--file must be provided")
		}
		data, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read credential file: %w", err)
		}
		var cred credentials.Credential
		if err := json.Unmarshal(data, &cred); err != nil {
			return fmt.Errorf("invalid credential JSON: %w", err)
		}
		if err := credentials.VerifyCredential(&cred); err != nil {
			cmd.Printf("Credential verification failed: %v\n", err)
			os.Exit(1)
		}
		cmd.Println("Credential is valid ✅")
		return nil
	},
}

func init() {
	verifyCmd.Flags().StringVar(&credentialFile, "file", "", "Path to credential JSON file (required)")
	verifyCmd.MarkFlagRequired("file")
}
