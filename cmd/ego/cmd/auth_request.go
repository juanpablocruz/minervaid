package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/spf13/cobra"
)

var (
	clientID            string
	redirectURI         string
	authorizeEndpoint   string
	nonce               string
	presentationDefFile string
)

var authRequestCmd = &cobra.Command{
	Use:   "auth-request",
	Short: "Generate OIDC4VP authorization URL",
	Long:  `Build the URL for an OIDC4VP /authorize request using a presentation definition and nonce.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load presentation_definition JSON
		var pd interface{}
		data, err := os.ReadFile(presentationDefFile)
		if err != nil {
			return fmt.Errorf("read presentation definition: %w", err)
		}
		if err := json.Unmarshal(data, &pd); err != nil {
			return fmt.Errorf("invalid JSON in presentation definition: %w", err)
		}

		// Build URL
		u, err := url.Parse(authorizeEndpoint)
		if err != nil {
			return fmt.Errorf("invalid authorize endpoint: %w", err)
		}
		q := u.Query()
		q.Set("response_type", "id_token")
		q.Set("client_id", clientID)
		q.Set("redirect_uri", redirectURI)
		q.Set("scope", "openid")
		q.Set("nonce", nonce)
		// embed presentation_definition as encoded JSON
		pdBytes, _ := json.Marshal(pd)
		q.Set("presentation_definition", string(pdBytes))
		u.RawQuery = q.Encode()

		cmd.Println(u.String())
		return nil
	},
}

func init() {
	authRequestCmd.Flags().StringVar(&authorizeEndpoint, "endpoint", "", "OIDC authorize endpoint (required)")
	authRequestCmd.Flags().StringVar(&clientID, "client-id", "", "Client ID (required)")
	authRequestCmd.Flags().StringVar(&redirectURI, "redirect-uri", "", "Redirect URI (required)")
	authRequestCmd.Flags().StringVar(&presentationDefFile, "presentation-definition", "", "Path to presentation definition JSON (required)")
	authRequestCmd.Flags().StringVar(&nonce, "nonce", "", "Nonce value (required)")
	authRequestCmd.MarkFlagRequired("endpoint")
	authRequestCmd.MarkFlagRequired("client-id")
	authRequestCmd.MarkFlagRequired("redirect-uri")
	authRequestCmd.MarkFlagRequired("presentation-definition")
	authRequestCmd.MarkFlagRequired("nonce")
	rootCmd.AddCommand(authRequestCmd)
}
