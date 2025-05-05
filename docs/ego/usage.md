# Ego CLI and Extension Usage Guide

This guide walks you through installing and using the Ego CLI and browser extension for managing your decentralized identity (DID) and Verifiable Credentials (VCs) in SSI workflows.

---

## 1. Quickstart

### 1.1 Install Ego CLI

```bash
# Requires Go 1.18+
go install github.com/juanpablocruz/minervaid/cmd/ego@latest
```

### 1.2 Install Browser Extension

1. Clone or download the `extension/` folder from the Ego repo.
2. In Chrome/Edge: go to `chrome://extensions`, enable “Developer mode”, click “Load unpacked”, and select the `extension/` folder.
3. In Firefox: go to `about:debugging`, click “This Firefox”, then “Load Temporary Add-on” and choose the `manifest.json`.

### 1.3 Initialize Your Vault

Create a vault (a secure folder) containing your DID and key pair:

```bash
ego init --name alice --out ./store
```

Select it as active:

```bash
ego use alice --store ./store
```

### 1.4 Set Attributes and Issue a Credential

Store metadata as key/value pairs:

```bash
ego set email alice@example.com --store ./store
ego set role admin --store ./store
```

Issue a Verifiable Credential capturing all stored attributes:

```bash
ego issue --id vc-auth --store ./store
```

### 1.5 Generate a Verifiable Presentation

Present one or more credentials, optionally revealing specific fields:

```bash
ego present --creds vc-auth --reveal email --out ./store > vp.json
```

---

## 2. CLI Commands Reference

Below is a summary of all Ego CLI commands. Use `ego <command> --help` for full details and available flags.

| Command                 | Description                                                     |
| ----------------------- | --------------------------------------------------------------- |
| `ego init`              | Create a new vault with a fresh DID and keystore.               |
| `ego use`               | Select an active vault by name.                                 |
| `ego set <key> <value>` | Add or update a metadata attribute in the vault.                |
| `ego issue`             | Issue a new Verifiable Credential from stored attributes.       |
| `ego present`           | Create a Verifiable Presentation from existing credentials.     |
| `ego revoke`            | Revoke a credential in the vault.                               |
| `ego auth-request`      | Build the OIDC4VP authorization URL (challenge request).        |
| `ego auth-callback`     | Launch HTTP server to capture the `id_token` callback.          |
| `ego auth-verify`       | Verify an OIDC4VP `id_token` and extract the authenticated DID. |

---

## 3. Browser Extension Workflow

1. **Vault Setup**: ensure your vault is created and active via Ego CLI.
2. **Authorize**: when a website triggers an OIDC4VP flow, the extension intercepts the challenge URL.
3. **Consent Prompt**: extension UI displays requested fields and credentials.
4. **Native Messaging**: extension calls Ego CLI (`ego present`) to produce the VP.
5. **Response Delivery**: extension sends the VP to the application to complete authentication.

---
