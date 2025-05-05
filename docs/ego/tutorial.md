# End-to-End Tutorial: SSI Login Flow with Ego

This tutorial demonstrates a complete end-to-end authentication flow using the Ego CLI and a minimal Authorization Server (AS) that issues and verifies cryptographic challenges.

## Prerequisites

* **Ego CLI** installed (`go install github.com/juanpablocruz/minervaid/cmd/ego@latest`).
* **minervaid** server binaries (`auth-challenge`, `auth-respond`, `auth-verify`).
* A working directory `./store` for vault data.

---

## 1. Initialize Your Vault

1. Create a vault named `alice`:

   ```bash
   ego init --name alice --out ./store
   ```

2. Select it as the active vault:

   ```bash
   ego use alice --store ./store
   ```

After these commands, `./store/alice` contains:

* `did.json`: your DID Document
* `keystore.json`: your encoded private key

---

## 2. Issue a Verifiable Credential (VC)

1. Add attributes in your vault:

   ```bash
   ego set email alice@example.com --store ./store
   ego set role admin --store ./store
   ```

2. Issue a credential capturing those attributes:

   ```bash
   ego issue --id vc1 --store ./store
   ```

Result: `./store/alice/credentials/vc1.json` is signed by your DID.

---

## 3. AS Generates an Authentication Challenge

On your Authorization Server, run:

```bash
minervaid auth-challenge \
  --domain api.example.com \
  --expiry 5m \
  --store ./store/alice \
> challenge.json
```

Sample `challenge.json`:

```json
{
  "type": "AuthenticationChallenge",
  "challenge": "HJUcjqHjAYT4QQ8239q57K",
  "domain": "api.example.com",
  "issuedAt": "2025-05-05T06:18:43Z",
  "expiresAt": "2025-05-05T06:23:43Z"
}
```

**Server must persist** the `challenge` nonce for later verification.

---

## 4. Client Signs the Challenge

In the client (e.g., user’s CLI or wallet), produce a signed response:

```bash
ego auth-respond \
  --did $(ego did list --store ./store/alice | head -1) \
  --file challenge.json \
  --store ./store/alice \
> response.json
```

Sample `response.json`:

```json
{
  "challenge": { /* original challenge */ },
  "proof": {
    "type": "Ed25519Signature2018",
    "created": "2025-05-05T06:20:33Z",
    "proofPurpose": "authentication",
    "verificationMethod": "did:key:z…#keys-1",
    "jws": "2b065b9f86fc…"
  }
}
```

This proof signs the raw challenge bytes with your private key.

---

## 5. AS Verifies the Response and Issues Session

Client POSTs to the AS endpoint:

```bash
curl -X POST https://api.example.com/verify \
  -H "Content-Type: application/json" \
  --data @response.json
```

On the server:

1. **Load** the original challenge from storage by its `nonce`.
2. **Compare** the stored object with `response.challenge`.
3. **Verify** the Ed25519 signature in `response.proof` against the exact challenge bytes.
4. Upon success, **create** a session (cookie or JWT) for the authenticated DID.

Alternatively, run locally:

```bash
minervaid auth-verify --file response.json --store ./store/alice
```

---

## 6. Automated Example Script

```bash
#!/usr/bin/env bash
set -e

VAULT=./store/alice
AS_URL=https://api.example.com

# 1. Generate challenge
minervaid auth-challenge --domain api.example.com --expiry 5m --store $VAULT > challenge.json

# 2. Respond to challenge
ego auth-respond --did $(ego did list --store $VAULT | head -1) \
  --file challenge.json --store $VAULT > response.json

# 3. Verify and login
curl -s -X POST $AS_URL/verify \
  -H "Content-Type: application/json" \
  --data @response.json && echo "✅ Authentication succeeded"
```

Save this as `login.sh` and run `bash login.sh` to exercise the full flow.

---
