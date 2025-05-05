package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var callbackPort int

var authCallbackCmd = &cobra.Command{
	Use:   "auth-callback --port <port>",
	Short: "Start a local HTTP server to capture OIDC4VP callback",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		port := callbackPort
		http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
			idToken := r.URL.Query().Get("id_token")
			if idToken == "" {
				// maybe fragment
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, "id_token missing")
				return
			}
			fmt.Fprintln(w, "Authentication successful. You can close this window.")
			// Print to stdout for CLI capture
			fmt.Printf("%s", idToken)
			// Shutdown after a moment
			go func() { time.Sleep(100 * time.Millisecond); os.Exit(0) }()
		})
		addr := fmt.Sprintf(":%d", port)
		log.Printf("Listening on %s for callback...", addr)
		http.ListenAndServe(addr, nil)
		return nil
	},
}

func init() {
	authCallbackCmd.Flags().IntVar(&callbackPort, "port", 8080, "Port to listen for callback")
	rootCmd.AddCommand(authCallbackCmd)
}
