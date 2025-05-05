package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/juanpablocruz/minervaid/internal/credentials"
	"github.com/mr-tron/base58"
	"github.com/spf13/cobra"
)

var (
	credDid     string
	credSubject string
	credID      string
	zkpMinAge   uint64
)

// newCredCmd issues a new Verifiable Credential, optionally attaching a range proof.
var newCredCmd = &cobra.Command{
	Use:   "new-cred --did <did> --subject <json|@file> [--id <id>] [--zkp-min-age <n>]",
	Short: "Issue a new Verifiable Credential",
	RunE: func(cmd *cobra.Command, args []string) error {
		if credDid == "" || credSubject == "" {
			return fmt.Errorf("--did and --subject are required")
		}
		// Ensure store directory exists
		if err := os.MkdirAll(storeDir, 0755); err != nil {
			return err
		}

		// Load issuer private key
		ks := loadKeyStore(filepath.Join(storeDir, "keystore.json"))
		privEnc, ok := ks[credDid]
		if !ok {
			return fmt.Errorf("unknown DID: %s", credDid)
		}
		privBytes, err := base58.Decode(privEnc)
		if err != nil {
			return fmt.Errorf("decoding private key: %w", err)
		}

		// Parse subject JSON
		subjData := []byte(credSubject)
		if len(subjData) > 0 && subjData[0] == '@' {
			subjData, err = os.ReadFile(string(subjData[1:]))
			if err != nil {
				return fmt.Errorf("reading subject file: %w", err)
			}
		}
		var subj map[string]interface{}
		if err := json.Unmarshal(subjData, &subj); err != nil {
			return fmt.Errorf("invalid subject JSON: %w", err)
		}

		// Determine credential ID
		id := credID
		if id == "" {
			id = time.Now().UTC().Format("20060102T150405Z")
		}

		// Create credential
		cred := credentials.NewCredential(id, credDid, subj)

		// Optionally attach a range proof via generic challenge
		if zkpMinAge > 0 {
			ch := credentials.Challenge{
				Type:   "range",
				Params: map[string]interface{}{"field": "age", "min": zkpMinAge},
			}
			chJSON, err := json.Marshal(ch)
			if err != nil {
				return fmt.Errorf("building ZKP challenge JSON: %w", err)
			}
			if err := cred.AttachProof(chJSON); err != nil {
				return fmt.Errorf("attaching proof: %w", err)
			}
		}

		// Sign credential
		if err := cred.SignCredential(privBytes, credDid+"#keys-1"); err != nil {
			return fmt.Errorf("signing credential: %w", err)
		}

		// Save credential to file
		credDir := filepath.Join(storeDir, "credentials")
		if err := os.MkdirAll(credDir, 0755); err != nil {
			return fmt.Errorf("creating credentials dir: %w", err)
		}
		store := &credentials.FileStore{Dir: credDir}
		if err := store.Save(cred); err != nil {
			return fmt.Errorf("saving credential: %w", err)
		}

		fmt.Println(id)
		return nil
	},
}

func init() {
	newCredCmd.Flags().StringVar(&credDid, "did", "", "Issuer DID (required)")
	newCredCmd.Flags().StringVar(&credSubject, "subject", "", "Subject JSON or @file (required)")
	newCredCmd.Flags().StringVar(&credID, "id", "", "Credential ID (optional)")
	newCredCmd.Flags().Uint64Var(&zkpMinAge, "zkp-min-age", 0, "Generate ZKP proof for age >= this value")
	newCredCmd.Flags().StringVar(&storeDir, "store", "./store", "Directory for storing data")
	_ = newCredCmd.MarkFlagRequired("did")
	_ = newCredCmd.MarkFlagRequired("subject")
	rootCmd.AddCommand(newCredCmd)
}
