package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/juanpablocruz/minervaid/internal/credentials"
	"github.com/mr-tron/base58"
	"github.com/spf13/cobra"
)

var (
	presDid    string
	presCreds  string
	presReveal string
)

var newPresentationCmd = &cobra.Command{
	Use:   "new-presentation",
	Short: "Generate a Verifiable Presentation",
	RunE: func(cmd *cobra.Command, args []string) error {
		if presDid == "" || presCreds == "" {
			return fmt.Errorf("--did and --creds are required")
		}
		if err := os.MkdirAll(storeDir, 0755); err != nil {
			return err
		}
		ks := loadKeyStore(filepath.Join(storeDir, "keystore.json"))
		privEnc, ok := ks[presDid]
		if !ok {
			return fmt.Errorf("unknown DID: %s", presDid)
		}
		privBytes, err := base58.Decode(privEnc)
		if err != nil {
			return err
		}
		ids := strings.Split(presCreds, ",")
		var creds []credentials.Credential
		for _, id := range ids {
			c, err := (&credentials.FileStore{Dir: filepath.Join(storeDir, "credentials")}).Get(id)
			if err != nil {
				return err
			}
			creds = append(creds, *c)
		}
		if presReveal != "" {
			fields := strings.Split(presReveal, ",")
			for i, c := range creds {
				filtered := make(map[string]interface{})
				for _, f := range fields {
					if v, ok := c.CredentialSubject[f]; ok {
						filtered[f] = v
					}
				}
				c.CredentialSubject = filtered
				creds[i] = c
			}
		}
		pres := credentials.NewPresentation(creds, presDid)
		if err := pres.SignPresentation(privBytes, presDid+"#keys-1"); err != nil {
			return err
		}
		outDir := filepath.Join(storeDir, "presentations")
		if err := os.MkdirAll(outDir, 0755); err != nil {
			return err
		}
		id := time.Now().UTC().Format("20060102T150405Z")
		data, _ := pres.ToJSON()
		if err := os.WriteFile(filepath.Join(outDir, id+".json"), data, 0644); err != nil {
			return err
		}
		fmt.Println(id)
		return nil
	},
}

func init() {
	newPresentationCmd.Flags().StringVar(&presDid, "did", "", "Holder DID")
	newPresentationCmd.Flags().StringVar(&presCreds, "creds", "", "Comma-separated credential IDs")
	newPresentationCmd.Flags().StringVar(&presReveal, "reveal", "", "Fields to reveal")
	rootCmd.AddCommand(newPresentationCmd)
}
