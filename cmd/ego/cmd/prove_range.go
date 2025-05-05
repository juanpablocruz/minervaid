package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/juanpablocruz/minervaid/internal/credentials"
	"github.com/spf13/cobra"
)

var (
	proofField string
	proofMin   uint64
	credPath   string
	storeDir   string
)

// proveRangeCmd generates a Bulletproof range proof for a numeric field in a credential JSON.
var proveRangeCmd = &cobra.Command{
	Use:   "prove-range --field <field> --min <minValue> --cred <path> [--store <dir>]",
	Short: "Generate a Bulletproof range proof for a credential field",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Read credential file
		data, err := os.ReadFile(credPath)
		if err != nil {
			return fmt.Errorf("reading credential: %w", err)
		}

		// Parse JSON into generic map
		var cred map[string]interface{}
		if err := json.Unmarshal(data, &cred); err != nil {
			return fmt.Errorf("invalid credential JSON: %w", err)
		}

		// Extract credentialSubject
		subj, ok := cred["credentialSubject"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("missing credentialSubject in credential")
		}

		// Extract numeric value
		raw, ok := subj[proofField]
		if !ok {
			return fmt.Errorf("field '%s' not found in credentialSubject", proofField)
		}
		// Only numeric types supported
		var value uint64
		switch v := raw.(type) {
		case float64:
			value = uint64(v)
		case int:
			value = uint64(v)
		case int64:
			value = uint64(v)
		case uint64:
			value = v
		default:
			return fmt.Errorf("unsupported field type %T for range proof", raw)
		}

		// Generate RangeProof
		rp, err := credentials.GenerateRangeProof(value, proofMin)
		if err != nil {
			return fmt.Errorf("generating range proof: %w", err)
		}

		// Marshal to JSON
		out, err := json.MarshalIndent(rp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshaling proof JSON: %w", err)
		}

		// Print JSON proof
		fmt.Println(string(out))
		return nil
	},
}

func init() {
	proveRangeCmd.Flags().StringVar(&proofField, "field", "", "CredentialSubject field to prove (required)")
	proveRangeCmd.Flags().Uint64Var(&proofMin, "min", 0, "Minimum value for range proof (required)")
	proveRangeCmd.Flags().StringVar(&credPath, "cred", "", "Path to credential JSON file (required)")
	proveRangeCmd.Flags().StringVar(&storeDir, "store", "./store", "Store directory for context")
	_ = proveRangeCmd.MarkFlagRequired("field")
	_ = proveRangeCmd.MarkFlagRequired("min")
	_ = proveRangeCmd.MarkFlagRequired("cred")
	rootCmd.AddCommand(proveRangeCmd)
}
