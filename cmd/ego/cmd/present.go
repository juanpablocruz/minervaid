package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/juanpablocruz/minervaid/internal/config"
	"github.com/juanpablocruz/minervaid/internal/credentials"
	"github.com/juanpablocruz/minervaid/internal/vault"
	"github.com/spf13/cobra"
)

var (
	credsFlag     string
	revealFlag    string
	outVault      string
	zkpChallenges []string
)

// presentCmd creates a Verifiable Presentation from existing credentials
var presentCmd = &cobra.Command{
	Use:   "present [--creds <id,id,...>] [--reveal <field,field,...>] [--out <directory>]",
	Short: "Create a Verifiable Presentation",
	Long: `Load one or more VCs from vault credentials, optionally apply selective disclosure,
and sign a Verifiable Presentation.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load CLI config
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		rootDir := cfg.RootDir
		if rootDir == "" {
			rootDir = "store"
		}

		// Determine vaultDir
		vaultDir := outVault
		if vaultDir == "" {
			vaultDir = filepath.Join(rootDir, cfg.Active)
		}

		// Load DID and key
		v := vault.NewVault(vaultDir)
		didDoc, priv, err := v.Load()
		if err != nil {
			return fmt.Errorf("load vault: %w", err)
		}
		// Parse DID
		var didDocMap map[string]interface{}
		if err := json.Unmarshal(didDoc, &didDocMap); err != nil {
			return fmt.Errorf("parse did.json: %w", err)
		}
		did, _ := didDocMap["id"].(string)

		// Determine credential IDs
		var ids []string
		if credsFlag != "" {
			ids = strings.Split(credsFlag, ",")
		} else {
			// load all
			credDir := filepath.Join(vaultDir, "credentials")
			files, err := os.ReadDir(credDir)
			if err != nil {
				return fmt.Errorf("read credentials dir: %w", err)
			}
			for _, f := range files {
				if !f.IsDir() && filepath.Ext(f.Name()) == ".json" {
					ids = append(ids, strings.TrimSuffix(f.Name(), ".json"))
				}
			}
		}

		// Load credentials
		var credsList []credentials.Credential
		for _, id := range ids {
			c, err := (&credentials.FileStore{Dir: filepath.Join(vaultDir, "credentials")}).Get(id)
			if err != nil {
				return fmt.Errorf("load credential %s: %w", id, err)
			}
			credsList = append(credsList, *c)
		}

		// Apply selective disclosure
		if revealFlag != "" {
			fields := strings.Split(revealFlag, ",")
			for i, cred := range credsList {
				filtered := make(map[string]interface{})
				for _, f := range fields {
					if v, ok := cred.CredentialSubject[f]; ok {
						filtered[f] = v
					}
				}
				cred.CredentialSubject = filtered
				credsList[i] = cred
			}
		}
		// Apply ZKP challenges
		for _, entry := range zkpChallenges {
			parts := strings.Split(entry, ":")
			if len(parts) != 3 {
				return fmt.Errorf("invalid zkp spec %q; expected type:field:param", entry)
			}
			proofType, field, rawParam := parts[0], parts[1], parts[2]
			switch proofType {
			case "range":
				minVal, err := strconv.ParseUint(rawParam, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid min in zkp %q: %w", entry, err)
				}
				ch := credentials.Challenge{
					Type:   "range",
					Params: map[string]interface{}{"field": field, "min": minVal},
				}
				chBytes, _ := json.Marshal(ch)
				for i := range credsList {
					if err := credsList[i].AttachProof(chBytes); err != nil {
						return fmt.Errorf("attaching %q to %s: %w", entry, credsList[i].ID, err)
					}
				}
			default:
				return fmt.Errorf("unsupported proof type %q", proofType)
			}
		}

		// If we’ve attached ZKPs but no explicit reveal, redact all other fields
		if len(zkpChallenges) > 0 && revealFlag == "" {
			for i := range credsList {
				// zero out credentialSubject so nobody sees the raw values
				credsList[i].CredentialSubject = map[string]interface{}{}
			}
		}

		// Build presentation
		pres := credentials.NewPresentation(credsList, did)
		if err := pres.SignPresentation(priv, did+"#keys-1"); err != nil {
			return fmt.Errorf("sign presentation: %w", err)
		}

		// Save presentation
		presDir := filepath.Join(vaultDir, "presentations")
		if err := os.MkdirAll(presDir, 0700); err != nil {
			return fmt.Errorf("create presentations dir: %w", err)
		}
		presID := time.Now().UTC().Format("20060102T150405Z")
		outFile := filepath.Join(presDir, presID+".json")
		data, err := pres.ToJSON()
		if err != nil {
			return fmt.Errorf("marshal presentation: %w", err)
		}
		if err := os.WriteFile(outFile, data, 0600); err != nil {
			return fmt.Errorf("write presentation: %w", err)
		}

		cmd.Printf("Presentation '%s' created\n", presID)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(presentCmd)
	presentCmd.Flags().StringVar(&credsFlag, "creds", "", "Comma-separated credential IDs (default: all)")
	presentCmd.Flags().StringVar(&revealFlag, "reveal", "", "Comma-separated fields to reveal (optional)")
	presentCmd.Flags().StringVar(&outVault, "out", "", "Vault directory (optional, uses active)")
	presentCmd.Flags().
		StringSliceVar(&zkpChallenges, "zkp", nil,
			"Add a zero-knowledge proof from `<type>:<field>:<param>`; can be repeated")
}
