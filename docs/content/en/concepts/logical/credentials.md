---
title: Credentials
description: Meshery uses one or more Credentials when authenticating to a managed or unmanaged Connection.
aliases:
- /concepts/credentials/
---
Meshery uses one or more Credentials when authenticating to a managed or unmanaged Connection. Credentials are based on the Meshery's [Credential Schema](https://github.com/meshery/schemas/blob/master/schemas/constructs/v1beta2/credential/api.yml) with defined attributes.

## Understanding Credentials in Meshery

Credentials in Meshery provide secure authentication to your infrastructure and cloud native services. They are a foundational component that enables Meshery to establish and maintain secure connections with your resources.

### Types of Credentials

Meshery supports several types of credentials:

- **API Keys/Tokens** - For services that use token-based authentication
- **Username/Password** - Basic authentication credentials
- **Certificates** - TLS/SSL certificates for secure connections
- **Cloud Provider Credentials** - Authentication for various cloud platforms
- **Service Account Tokens** - For Kubernetes authentication

### Managing Credentials

Credentials can be:

1. Created and managed through the Meshery UI
2. Imported from existing configuration
3. Auto-discovered from your environment
4. Associated with one or more Connections
5. Shared across team members (with appropriate permissions)

### Credential Security

Meshery takes several measures to protect your credentials:

- Credential secrets are [encrypted at rest](#encryption-at-rest) in Meshery's own datastore
- Access is controlled through fine-grained permissions
- Credential secrets are not written to Meshery Server logs
- Support for secret management integration

#### Encryption at rest

When Meshery Server persists a credential to its own datastore, the credential's
secret is sealed with AES-256-GCM before it is written. The row holds a
ciphertext envelope rather than the API key, token or service-account credential
you entered. Decryption happens inside the server as it reads the credential, so
this is invisible in the UI and in the API.

The encryption key is derived, through HKDF-SHA256, from the secret that is built
into the Meshery Server binary at release time. **The key therefore ships inside
the binary.** Be precise about what that buys you:

- **It protects a datastore separated from the binary that wrote it** - a stolen
  or exfiltrated database file, a copied `~/.meshery` directory, a filesystem
  backup or volume snapshot, a support bundle. In all of those the credentials
  are ciphertext.
- **It does not protect against someone who has the Meshery Server binary or
  container image**, who can recover the key from it. This is not key-managed
  encryption, and it is not a substitute for keeping your datastore, your
  backups and your images access-controlled.
- **The same secret is also the `X-API-Key` Meshery Server sends when submitting
  anonymous performance results**, so it leaves the binary on the network too,
  and anyone who captures it can derive the same key.

Two consequences are worth knowing before you meet them:

- **Upgrading takes no action to keep working, but one action to finish the
  job.** Credentials written before this shipped are plaintext and keep reading;
  each becomes ciphertext the next time it is written, and there is no migration
  to run. That conversion is not a scrub, though. SQLite does not zero a page
  when it frees one, so on an install that predates this change the earlier
  plaintext secret can remain recoverable from the raw datastore file even after
  the credential has been rewritten as ciphertext - and rewriting more rows will
  not clear it, which is why a bulk migration would not help either. Reclaim
  those pages by compacting the datastore with Meshery Server stopped:

  ```bash
  sqlite3 ~/.meshery/config/mesherydb.sql 'VACUUM;'
  ```

  A fresh install is unaffected: every credential in it was sealed on its first
  write, so there is no earlier plaintext to reclaim.
- **Credentials are readable only by a build carrying the same secret.** A
  locally built Meshery Server (`make server`) links no release secret and falls
  back to a development default, so credentials it wrote cannot be read by a
  released image, and vice versa, even against the same `~/.meshery`. Every
  ciphertext carries a short identifier of the key that sealed it, so Meshery
  reports this as a key mismatch ([`meshery-server-1484`]({{< ref
  "reference/references/error-codes.md" >}})) rather than as a corrupt
  credential. The remedy is to run the build that wrote them, or to re-enter the
  affected credentials.

When Meshery is configured with a [Remote Provider]({{< ref
"reference/extensibility/providers/index.md" >}}), that provider persists
credentials in its own datastore and is responsible for protecting them there;
Meshery sends them to it over TLS.

### Using Credentials with Connections

When setting up a new Connection in Meshery:

1. Select from existing credentials or create new ones
2. Credentials are automatically validated before use
3. Multiple credentials can be associated with a single Connection
4. Credential status is monitored and alerts are generated if they become invalid


