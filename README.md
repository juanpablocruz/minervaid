# CLI Usage Guide for MinervaID

MinervaID provides a command-line interface for decentralized SSI: DIDs, Verifiable Credentials (with optional zero-knowledge proofs), Presentations, Revocation, and Authentication.

## Installation

```bash
# From module root
go install github.com/juanpablocruz/minervaid/cmd/minervaid@latest
```

## Global Options

| Flag           | Description                         | Default   |
| -------------- | ----------------------------------- | --------- |
| `--store DIR`  | Directory for storing data          | `./store` |

## Commands

### 1. DID Management

- **new-did**  
  Generate a new DID and store its private key in `<store>/keystore.json`.  
  **Usage:**

  ```bash
  minervaid new-did --store ./store
  ```

  **Example output:**

  ```
  New DID: did:key:z123abc...
  ```

- **list-dids**  
  List all stored DIDs.  
  **Usage:**

  ```bash
  minervaid list-dids --store ./store
  ```

### 2. Verifiable Credentials

- **new-cred**  
  Create and sign a new Verifiable Credential.  
  **Flags:**

  | Flag               | Description                                                                  |
  | ------------------ | ---------------------------------------------------------------------------- |
  | `--did DID`        | Issuer DID (must exist via `new-did`)                                         |
  | `--subject JSON`   | Subject JSON inline or `@file.json` (e.g., `'{"id":"did:ex:456","age":30}'`) |
  | `--id ID`          | Credential ID (optional; defaults to timestamp)                               |
  | `--zkp-min-age N`  | Attach ZK proof that `age ≥ N`, removing cleartext `age`                      |
  
  **Usage:**

  ```bash
  minervaid new-cred \
    --did did:key:z123abc... \
    --subject '{"id":"did:example:456","age":30}' \
    --store ./store
  ```

  or with ZKP:

  ```bash
  minervaid new-cred \
    --did did:key:z123abc... \
    --subject @subject.json \
    --zkp-min-age 18 \
    --store ./store
  ```

- **list-creds**  
  List all issued credential IDs.  
  **Usage:**

  ```bash
  minervaid list-creds --store ./store
  ```

- **get-cred**  
  Retrieve and print a credential JSON by ID.  
  **Usage:**

  ```bash
  minervaid get-cred --id <credID> --store ./store
  ```

- **verify-cred**  
  Verify the Ed25519 signature and any embedded zero-knowledge proof.  
  **Usage:**

  ```bash
  minervaid verify-cred --file ./store/credentials/<credID>.json
  ```

  **Output:**

  ```
  Credential is valid ✅
  ```

### 3. Presentations

- **new-presentation**  
  Generate a Verifiable Presentation from one or more credentials.  
  **Flags:**

  | Flag              | Description                                                       |
  | ----------------- | ----------------------------------------------------------------- |
  | `--did DID`       | Holder DID for signing                                            |
  | `--creds IDs`     | Comma-separated list of credential IDs                            |
  | `--reveal fields` | Comma-separated list of JSON field names to selectively disclose  |
  
  **Usage:**

  ```bash
  minervaid new-presentation \
    --did did:key:z123abc... \
    --creds <credID> \
    --reveal name,age \
    --store ./store
  ```

- **list-presents**  
  List all stored presentation IDs.  
  **Usage:**

  ```bash
  minervaid list-presents --store ./store
  ```

- **get-presentation**  
  Retrieve and print a presentation JSON by ID.  
  **Usage:**

  ```bash
  minervaid get-presentation --id <presID> --store ./store
  ```

- **verify-presentation**  
  Verify the presentation signature and embedded credentials.  
  **Usage:**

  ```bash
  minervaid verify-presentation --file ./store/presentations/<presID>.json
  ```

  **Output:**

  ```
  Presentation is valid ✅
  ```

### 4. Revocation

- **revoke-cred**  
  Revoke a credential by ID.  
  **Usage:**

  ```bash
  minervaid revoke-cred --id <credID> --store ./store
  ```

- **list-revoked**  
  List all revoked credential IDs.  
  **Usage:**

  ```bash
  minervaid list-revoked --store ./store
  ```

- **check-revoked**  
  Check the revocation status of a credential.  
  **Usage:**

  ```bash
  minervaid check-revoked --id <credID> --store ./store
  ```

### 5. Authentication

- **auth-challenge**  
  Generate a cryptographic challenge (nonce) for a DID.  
  **Usage:**

  ```bash
  minervaid auth-challenge \
    --domain api.example.com \
    --expiry 5m \
    --store ./store
  ```

- **auth-respond**  
  Respond to an authentication challenge by signing it with your DID.  
  **Usage:**

  ```bash
  minervaid auth-respond \
    --did did:key:z123abc... \
    --file challenge.json \
    --store ./store
  ```

---

For more details, refer to the code in `internal/credentials` and `internal/identity`.
