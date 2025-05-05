# Documentation Roadmap for `docs/ego`

Below is a proposed set of documentation files and directories to include under `docs/ego/`. Each file serves a distinct purpose to guide users and developers through installation, usage, extension, and integration.

## 1. `functional.md`

**Purpose:** Describe the functional specification of the Ego Wallet browser extension (already created).

* Overview of components
* User flows (challenge interception, native messaging)
* Security and configuration

## 2. `usage.md`

**Purpose:** End-user guide for the Ego CLI and extension.

### Outline

1. **Quickstart**

   * Install Ego CLI and browser extension
   * Initialize a vault
   * Issue a credential
   * Present a credential
2. **CLI Commands**

   * `ego init`, `ego use`, `ego set`, `ego issue`, `ego present`, `ego revoke`, `ego auth-request`, `ego auth-callback`, `ego auth-verify`
   * Flags, examples, common workflows
3. **Browser Extension**

   * Installation (Chrome/Firefox)
   * Configuring native messaging host
   * Performing an OIDC4VP login

## 3. `api.md`

**Purpose:** Technical reference for the Native Messaging protocol and JSON schemas.

### Outline

1. **Native Messaging Endpoint**

   * Host registration name (`com.minervaid.ego`)
   * Message envelope (`{ cmd: [...], params?: {...} }`)
   * Response envelope (`{ stdout: string, stderr: string, code: int }`)
2. **Authentication Challenge Schema**

   * `AuthenticationChallenge` JSON structure
3. **Authentication Response Schema**

   * Signed proof object format

## 4. `architecture.md`

**Purpose:** High-level diagrams and component interactions.

### Outline

1. **Component Diagram** (Extension, Native Host, CLI, AS, RP)
2. **Sequence Diagrams**

   * OIDC4VP login flow
   * Credential issuance and presentation

## 5. `tutorial.md`

**Purpose:** Step-by-step end-to-end tutorial for a sample application.

### Outline

1. **Setup**: Ego CLI + extension + sample AS stub
2. **Vault & Credential**
3. **OIDC4VP Login**
4. **Verifying on the server side**

## 6. `examples/` directory

**Purpose:** Store concrete JSON examples and scripts.

* `challenge.json`
* `response.json`
* `presentation_definition.json`
* Shell scripts showcasing full flows (`login.sh`, `verify.sh`)

## 7. `configuration.md`

**Purpose:** Document configuration options and environment variables.

### Outline

1. **CLI config (`~/.ego/config.json`)**
2. **Extension settings (chrome.storage)**
3. **Native Messaging host path**
4. **AS endpoints and trust anchors**

---

*This plan defines the minimum documentation set to support users and integrators. Let me know if you’d like to add or adjust any of these items!*
