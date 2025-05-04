package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/juanpablocruz/minervaid/internal/credentials"
	"github.com/juanpablocruz/minervaid/internal/identity"
	"github.com/mr-tron/base58"
)

const (
	keystoreFile   = "keystore.json"
	credentialsDir = "credentials"
	revocationFile = "revocations.json"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: minervaid <command> [options]")
		fmt.Println("Commands: new-did, list-dids, new-cred, list-creds, get-cred")
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "new-did":
		newDID(os.Args[2:])
	case "list-dids":
		listDIDs(os.Args[2:])
	case "new-cred":
		newCred(os.Args[2:])
	case "list-creds":
		listCreds(os.Args[2:])
	case "get-cred":
		getCred(os.Args[2:])
	case "new-presentation":
		newPresentation(os.Args[2:])
	case "list-presents":
		listPresents(os.Args[2:])
	case "get-presentation":
		getPresentation(os.Args[2:])
	case "verify-cred":
		verifyCred(os.Args[2:])
	case "verify-presentation":
		verifyPresentation(os.Args[2:])
	case "revoke-cred":
		revokeCred(os.Args[2:])
	case "list-revoked":
		listRevoked(os.Args[2:])
	case "check-revoked":
		checkRevoked(os.Args[2:])
	default:
		log.Fatalf("Unknown command: %s", cmd)
	}
}

func newDID(args []string) {
	fs := flag.NewFlagSet("new-did", flag.ExitOnError)
	store := fs.String("store", "./store", "Directory for storing data")
	fs.Parse(args)

	if err := os.MkdirAll(*store, 0755); err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}
	pub, priv, err := identity.GenerateKeyPair()
	if err != nil {
		log.Fatalf("Error generating keypair: %v", err)
	}
	did := identity.GenerateDID(pub)
	privEnc := identity.EncodePrivateKey(priv)

	ksPath := filepath.Join(*store, keystoreFile)
	ks := loadKeyStore(ksPath)
	ks[did] = privEnc
	saveKeyStore(ksPath, ks)

	fmt.Println(did)
}

func listDIDs(args []string) {
	fs := flag.NewFlagSet("list-dids", flag.ExitOnError)
	store := fs.String("store", "./store", "Directory for storing data")
	fs.Parse(args)

	ksPath := filepath.Join(*store, keystoreFile)
	ks := loadKeyStore(ksPath)
	for did := range ks {
		fmt.Println(did)
	}
}

func newCred(args []string) {
	fs := flag.NewFlagSet("new-cred", flag.ExitOnError)
	store := fs.String("store", "./store", "Directory for storing data")
	didFlag := fs.String("did", "", "Issuer DID")
	subjFlag := fs.String("subject", "", "Subject JSON or @file")
	idFlag := fs.String("id", "", "Credential ID (optional)")
	fs.Parse(args)

	if *didFlag == "" || *subjFlag == "" {
		fs.Usage()
		os.Exit(1)
	}
	if err := os.MkdirAll(*store, 0755); err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}

	// Load issuer private key
	ks := loadKeyStore(filepath.Join(*store, keystoreFile))
	privEnc, ok := ks[*didFlag]
	if !ok {
		log.Fatalf("Unknown DID: %s", *didFlag)
	}
	privBytes, err := base58.Decode(privEnc)
	if err != nil {
		log.Fatalf("Decoding private key: %v", err)
	}

	// Parse subject JSON
	subjData := []byte(*subjFlag)
	if subjData[0] == '@' {
		subjData, err = os.ReadFile(string(subjData[1:]))
		if err != nil {
			log.Fatalf("Reading subject file: %v", err)
		}
	}
	var subj map[string]interface{}
	if err := json.Unmarshal(subjData, &subj); err != nil {
		log.Fatalf("Invalid subject JSON: %v", err)
	}

	// Assign an ID if none given
	credID := *idFlag
	if credID == "" {
		credID = time.Now().UTC().Format("20060102T150405Z")
	}

	// Create, sign and save
	cred := credentials.NewCredential(credID, *didFlag, subj)
	if err := cred.SignCredential(privBytes, *didFlag+"#keys-1"); err != nil {
		log.Fatalf("Signing credential: %v", err)
	}
	credDir := filepath.Join(*store, credentialsDir)
	if err := os.MkdirAll(credDir, 0755); err != nil {
		log.Fatalf("Creating credentials dir: %v", err)
	}
	fsStore := &credentials.FileStore{Dir: credDir}
	if err := fsStore.Save(cred); err != nil {
		log.Fatalf("Saving credential: %v", err)
	}

	fmt.Println(credID)
}

func listCreds(args []string) {
	fs := flag.NewFlagSet("list-creds", flag.ExitOnError)
	store := fs.String("store", "./store", "Directory for storing data")
	fs.Parse(args)

	credDir := filepath.Join(*store, credentialsDir)
	files, err := os.ReadDir(credDir)
	if err != nil {
		log.Fatalf("Listing creds: %v", err)
	}
	for _, f := range files {
		if !f.IsDir() && filepath.Ext(f.Name()) == ".json" {
			fmt.Println(f.Name()[:len(f.Name())-5])
		}
	}
}

func getCred(args []string) {
	fs := flag.NewFlagSet("get-cred", flag.ExitOnError)
	store := fs.String("store", "./store", "Directory for storing data")
	idFlag := fs.String("id", "", "Credential ID")
	fs.Parse(args)

	if *idFlag == "" {
		fs.Usage()
		os.Exit(1)
	}
	credDir := filepath.Join(*store, credentialsDir)
	fsStore := &credentials.FileStore{Dir: credDir}
	cred, err := fsStore.Get(*idFlag)
	if err != nil {
		log.Fatalf("Getting credential: %v", err)
	}
	data, _ := json.MarshalIndent(cred, "", "  ")
	fmt.Println(string(data))
}

func newPresentation(args []string) {
	fs := flag.NewFlagSet("new-presentation", flag.ExitOnError)
	store := fs.String("store", "./store", "Directory for storing data")
	didFlag := fs.String("did", "", "Holder DID")
	credsFlag := fs.String("creds", "", "Comma-separated credential IDs")
	revealFlag := fs.String("reveal", "", "Fields to reveal (optional)")
	fs.Parse(args)

	if *didFlag == "" || *credsFlag == "" {
		fs.Usage()
		os.Exit(1)
	}
	if err := os.MkdirAll(*store, 0755); err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}

	ks := loadKeyStore(filepath.Join(*store, keystoreFile))
	privEnc, ok := ks[*didFlag]
	if !ok {
		log.Fatalf("Unknown DID: %s", *didFlag)
	}
	privBytes, err := base58.Decode(privEnc)
	if err != nil {
		log.Fatalf("Decoding private key: %v", err)
	}

	ids := strings.Split(*credsFlag, ",")
	var credsList []credentials.Credential
	for _, id := range ids {
		c, err := (&credentials.FileStore{Dir: filepath.Join(*store, credentialsDir)}).Get(id)
		if err != nil {
			log.Fatalf("Loading credential %s: %v", id, err)
		}
		credsList = append(credsList, *c)
	}

	// Apply selective disclosure if --reveal provided
	if *revealFlag != "" {
		fields := strings.Split(*revealFlag, ",")
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

	pres := credentials.NewPresentation(credsList, *didFlag)
	if err := pres.SignPresentation(privBytes, *didFlag+"#keys-1"); err != nil {
		log.Fatalf("Signing presentation: %v", err)
	}
	presDir := filepath.Join(*store, "presentations")
	if err := os.MkdirAll(presDir, 0755); err != nil {
		log.Fatalf("Creating presentations dir: %v", err)
	}
	presID := time.Now().UTC().Format("20060102T150405Z")
	outFile := filepath.Join(presDir, presID+".json")
	data, _ := pres.ToJSON()
	if err := os.WriteFile(outFile, data, 0644); err != nil {
		log.Fatalf("Saving presentation: %v", err)
	}
	fmt.Println(presID)
}

func listPresents(args []string) {
	fs := flag.NewFlagSet("list-presents", flag.ExitOnError)
	store := fs.String("store", "./store", "Directory for storing data")
	fs.Parse(args)

	presDir := filepath.Join(*store, "presentations")
	files, err := os.ReadDir(presDir)
	if err != nil {
		log.Fatalf("Listing presentations: %v", err)
	}
	for _, f := range files {
		if !f.IsDir() && filepath.Ext(f.Name()) == ".json" {
			fmt.Println(strings.TrimSuffix(f.Name(), ".json"))
		}
	}
}

func getPresentation(args []string) {
	fs := flag.NewFlagSet("get-presentation", flag.ExitOnError)
	store := fs.String("store", "./store", "Directory for storing data")
	idFlag := fs.String("id", "", "Presentation ID")
	fs.Parse(args)

	if *idFlag == "" {
		fs.Usage()
		os.Exit(1)
	}
	path := filepath.Join(*store, "presentations", *idFlag+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Getting presentation: %v", err)
	}
	fmt.Println(string(data))
}

func verifyCred(args []string) {
	fs := flag.NewFlagSet("verify-cred", flag.ExitOnError)
	file := fs.String("file", "", "Path to credential JSON file")
	fs.Parse(args)
	if *file == "" {
		fs.Usage()
		os.Exit(1)
	}
	data, err := os.ReadFile(*file)
	if err != nil {
		log.Fatalf("reading file: %v", err)
	}
	var cred credentials.Credential
	if err := json.Unmarshal(data, &cred); err != nil {
		log.Fatalf("invalid credential JSON: %v", err)
	}
	if err := credentials.VerifyCredential(&cred); err != nil {
		fmt.Printf("Credential verification failed: %v", err)
		os.Exit(1)
	}
	fmt.Println("Credential is valid ✅")
}

func verifyPresentation(args []string) {
	fs := flag.NewFlagSet("verify-presentation", flag.ExitOnError)
	file := fs.String("file", "", "Path to presentation JSON file")
	fs.Parse(args)
	if *file == "" {
		fs.Usage()
		os.Exit(1)
	}
	data, err := os.ReadFile(*file)
	if err != nil {
		log.Fatalf("reading file: %v", err)
	}
	var pres credentials.Presentation
	if err := json.Unmarshal(data, &pres); err != nil {
		log.Fatalf("invalid presentation JSON: %v", err)
	}
	if err := credentials.VerifyPresentation(&pres); err != nil {
		fmt.Printf("Presentation verification failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Presentation is valid ✅")
}

func revokeCred(args []string) {
	fs := flag.NewFlagSet("revoke-cred", flag.ExitOnError)
	store := fs.String("store", "./store", "Directory for storing data")
	idFlag := fs.String("id", "", "Credential ID to revoke")
	fs.Parse(args)
	if *idFlag == "" {
		fs.Usage()
		os.Exit(1)
	}
	if err := os.MkdirAll(*store, 0755); err != nil {
		log.Fatalf("creating store: %v", err)
	}
	path := filepath.Join(*store, revocationFile)
	rl, err := credentials.NewRevocationList(path)
	if err != nil {
		log.Fatalf("init revocation list: %v", err)
	}
	if err := rl.Revoke(*idFlag); err != nil {
		log.Fatalf("revoking credential: %v", err)
	}
	fmt.Printf("Credential %s revoked\n", *idFlag)
}

// list-revoked: list all revoked IDs
func listRevoked(args []string) {
	fs := flag.NewFlagSet("list-revoked", flag.ExitOnError)
	store := fs.String("store", "./store", "Directory for storing data")
	fs.Parse(args)
	path := filepath.Join(*store, revocationFile)
	rl, err := credentials.NewRevocationList(path)
	if err != nil {
		log.Fatalf("init revocation list: %v", err)
	}
	for _, id := range rl.List() {
		fmt.Println(id)
	}
}

func checkRevoked(args []string) {
	fs := flag.NewFlagSet("check-revoked", flag.ExitOnError)
	store := fs.String("store", "./store", "Directory for storing data")
	idFlag := fs.String("id", "", "Credential ID to check")
	fs.Parse(args)
	if *idFlag == "" {
		fs.Usage()
		os.Exit(1)
	}
	path := filepath.Join(*store, revocationFile)
	rl, err := credentials.NewRevocationList(path)
	if err != nil {
		log.Fatalf("init revocation list: %v", err)
	}
	if rl.IsRevoked(*idFlag) {
		fmt.Printf("Credential %s is revoked\n", *idFlag)
	} else {
		fmt.Printf("Credential %s is not revoked", *idFlag)
	}
}

func loadKeyStore(path string) map[string]string {
	ks := make(map[string]string)
	data, err := os.ReadFile(path)

	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("Failed to read keystore: %v", err)
	}

	if len(data) > 0 {
		if err := json.Unmarshal(data, &ks); err != nil {
			log.Fatalf("Invalid keystore format: %v", err)
		}
	}

	return ks
}

func saveKeyStore(path string, ks map[string]string) {
	data, _ := json.MarshalIndent(ks, "", "  ")

	if err := os.WriteFile(path, data, 0644); err != nil {
		log.Fatalf("Failed to save keystore: %v", err)
	}
}
