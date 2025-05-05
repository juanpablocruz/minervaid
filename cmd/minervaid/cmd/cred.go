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

var newCredCmd = &cobra.Command{
	Use:   "new-cred",
	Short: "Issue a new Verifiable Credential",
	RunE: func(cmd *cobra.Command, args []string) error {
		if credDid == "" || credSubject == "" {
			return fmt.Errorf("--did and --subject are required")
		}
		if err := os.MkdirAll(storeDir, 0755); err != nil {
			return err
		}
		ks := loadKeyStore(filepath.Join(storeDir, "keystore.json"))
		privEnc, ok := ks[credDid]
		if !ok {
			return fmt.Errorf("unknown DID: %s", credDid)
		}
		privBytes, err := base58.Decode(privEnc)
		if err != nil {
			return err
		}
		// parse subject
		subjData := []byte(credSubject)
		if subjData[0] == '@' {
			subjData, err = os.ReadFile(string(subjData[1:]))
			if err != nil {
				return err
			}
		}
		var subj map[string]interface{}
		if err := json.Unmarshal(subjData, &subj); err != nil {
			return err
		}
		id := credID
		if id == "" {
			id = time.Now().UTC().Format("20060102T150405Z")
		}
		cred := credentials.NewCredential(id, credDid, subj)
		if zkpMinAge > 0 {
			if err := cred.AttachAgeProof(zkpMinAge); err != nil {
				return err
			}
		}
		if err := cred.SignCredential(privBytes, credDid+"#keys-1"); err != nil {
			return err
		}
		credDir := filepath.Join(storeDir, "credentials")
		if err := os.MkdirAll(credDir, 0755); err != nil {
			return err
		}
		store := &credentials.FileStore{Dir: credDir}
		if err := store.Save(cred); err != nil {
			return err
		}
		fmt.Println(id)
		return nil
	},
}

func init() {
	newCredCmd.Flags().StringVar(&credDid, "did", "", "Issuer DID")
	newCredCmd.Flags().StringVar(&credSubject, "subject", "", "Subject JSON or @file")
	newCredCmd.Flags().StringVar(&credID, "id", "", "Credential ID (optional)")
	newCredCmd.Flags().Uint64Var(&zkpMinAge, "zkp-min-age", 0, "Generate ZKP age proof for age ≥ this value")
	rootCmd.AddCommand(newCredCmd)
}
