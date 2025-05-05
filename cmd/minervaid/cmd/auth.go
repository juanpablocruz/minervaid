package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/juanpablocruz/minervaid/internal/credentials"
	"github.com/spf13/cobra"
)

var (
	authDomain string
	authExpiry time.Duration
	authDid    string
	authFile   string
)

var authChallengeCmd = &cobra.Command{
	Use:   "auth-challenge",
	Short: "Generate an authentication challenge",
	RunE: func(cmd *cobra.Command, args []string) error {
		ch := credentials.NewAuthChallenge(authDomain, authExpiry)
		data, _ := json.MarshalIndent(ch, "", "  ")
		fmt.Println(string(data))
		return nil
	},
}

var authRespondCmd = &cobra.Command{
	Use:   "auth-respond",
	Short: "Respond to an authentication challenge",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := os.ReadFile(authFile)
		if err != nil {
			return err
		}
		var ch credentials.AuthenticationChallenge
		if err := json.Unmarshal(data, &ch); err != nil {
			return err
		}
		resp, err := credentials.RespondAuthChallenge(authDid, &ch, storeDir)
		if err != nil {
			return err
		}
		out, _ := json.MarshalIndent(resp, "", "  ")
		fmt.Println(string(out))
		return nil
	},
}

func init() {
	authChallengeCmd.Flags().StringVar(&authDomain, "domain", "", "Challenge domain")
	authChallengeCmd.Flags().DurationVar(&authExpiry, "expiry", 5*time.Minute, "Challenge expiry duration")
	authRespondCmd.Flags().StringVar(&authDid, "did", "", "Holder DID")
	authRespondCmd.Flags().StringVar(&authFile, "file", "", "Path to challenge JSON")
	rootCmd.AddCommand(authChallengeCmd, authRespondCmd)
}
