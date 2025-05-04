# CLI Usage Guide for MinervaID

## Installation

```bash
# from module root
go install github.com/tuusuario/minervaid/cmd/minervaid@latest
```

## Commands

- **new-did**
  - Generates a new DID and stores the private key in `<store>/keystore.json`.
  - Usage:

    ```bash
    minervaid new-did --store ./store
    ```

  - Example output:

    ```bash
    New DID: did:key:z123abc...
    ```

- **list-dids**
  - Lists all stored DIDs.
  - Usage:

    ```bash
    minervaid list-dids --store ./store
    ```

- **new-cred**
  - Creates and signs a new Verifiable Credential.
  - **Subject** can be provided as:
    - Inline JSON string: e.g.

      ```bash
      --subject '{"id":"did:example:456","name":"Alice","age":30}'
      ```

    - Or via file reference:

      ```bash
      --subject @subject.json
      ```

  - **ID** flag is optional; defaults to timestamp.
  - Usage examples:

    ```bash
    minervaid new-cred --did did:key:z123abc... --subject '{"id":"did:example:456","age":30}' --store ./store
    ```

    Or:

    ```bash
    minervaid new-cred --did did:key:z123abc... --subject @subject.json --id cred123 --store ./store
    ```

- **list-creds**
  - Lists all credential IDs saved under `<store>/credentials`.
  - Usage:

    ```bash
    minervaid list-creds --store ./store
    ```

- **get-cred**
  - Retrieves and prints a VC by its ID.
  - Usage:

    ```bash
    minervaid get-cred --id cred123 --store ./store
    ```

### Example subject.json

```json
{
  "id": "did:example:456",
  "name": "Alice",
  "age": 30
}
```

- **new-presentation**
  - Generates a Verifiable Presentation from one or more credentials.
  - **Flags**:
    - `--did DID` (holder DID used to sign the presentation)
    - `--creds ids` (comma-separated list of credential IDs)
    - `--reveal fields` (comma-separated list of JSON field names to selectively disclose)
  - Usage:

    ```bash
    minervaid new-presentation --did did:key:z123abc... --creds id1,id2 --reveal name,age --store ./store
    ```

- **list-presents**
  - Lists all stored presentations under `<store>/presentations`.
  - Usage:

    ```bash
    minervaid list-presents --store ./store
    ```

- **get-presentation**
  - Retrieves and prints a Verifiable Presentation by its ID.
  - Usage:

    ```bash
    minervaid get-presentation --id pres1 --store ./store
    ```
