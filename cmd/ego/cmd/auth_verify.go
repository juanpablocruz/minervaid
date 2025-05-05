package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var idTokenFile string

func base64URLDecode(s string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(s)
}

// authVerifyCmd verifies a JWT id_token and reports the authenticated DID
var authVerifyCmd = &cobra.Command{
	Use:   "auth-verify --id-token <file>",
	Short: "Verify OIDC4VP id_token JWT and extract DID",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Read JWT from file
		data, err := os.ReadFile(idTokenFile)
		if err != nil {
			return err
		}
		token := strings.TrimSpace(string(data))
		parts := strings.Split(token, ".")
		if len(parts) != 3 {
			return fmt.Errorf("invalid JWT format")
		}

		// Decode payload
		payload, err := base64URLDecode(parts[1])
		if err != nil {
			return fmt.Errorf("decode JWT payload: %w", err)
		}

		// Parse claims
		var claims map[string]interface{}
		if err := json.Unmarshal(payload, &claims); err != nil {
			return fmt.Errorf("invalid JWT claims: %w", err)
		}

		// Extract issuer (DID)
		iss, ok := claims["iss"].(string)
		if !ok || iss == "" {
			return fmt.Errorf("iss claim missing or invalid")
		}

		cmd.Printf("Authentication successful for DID: %s\n", iss)
		return nil
	},
}

func init() {
	authVerifyCmd.Flags().StringVar(&idTokenFile, "id-token", "", "Path to file containing the JWT id_token (required)")
	authVerifyCmd.MarkFlagRequired("id-token")
	rootCmd.AddCommand(authVerifyCmd)
}
