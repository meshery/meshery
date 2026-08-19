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

- Support for application-layer AES-256-GCM encryption at rest
- Access is controlled through fine-grained permissions
- Credentials are never exposed in logs or API responses
- Support for secret management integration

### At-Rest Encryption

Meshery provides optional application-layer encryption at rest for sensitive data (API keys, authentication tokens, and Kubernetes contexts stored in the datastore) using AES-256-GCM.

#### Enabling Encryption

At-rest encryption is activated by setting an environment variable before starting Meshery Server:

- `MESHERY_ENCRYPTION_KEY`: A 32-byte key in 64-character hexadecimal format (generate with `openssl rand -hex 32`) or 44-character base64 format (generate with `openssl rand -base64 32`).
- `MESHERY_ENCRYPTION_KEY_FILE`: Path to a file containing the 32-byte key (recommended for Kubernetes Secrets mounted as files).

When enabled:
- Credential secrets in the `credentials` table and kubeconfig `auth` and `cluster` sections in the `k8s_contexts` table are encrypted using AES-256-GCM before writing to the datastore.
- Each value uses a unique cryptographic nonce generated at encryption time.
- Key material is never stored in the database.
- If neither environment variable is set, encryption is disabled and existing plaintext credentials continue to function with zero migration required.

#### Migrating Existing Datastores

Existing unencrypted credentials and kubeconfigs can be encrypted in-place using `mesheryctl`:

```bash
# Ensure Meshery server is stopped first
mesheryctl system encrypt-datastore --key <your-32-byte-hex-key>

# Or with a key file
mesheryctl system encrypt-datastore --key-file /path/to/key.txt

# Preview changes without modifying the database
mesheryctl system encrypt-datastore --dry-run --key <your-32-byte-hex-key>

# To decrypt existing data back to plaintext:
mesheryctl system encrypt-datastore --decrypt --key <your-32-byte-hex-key>
```

### Using Credentials with Connections

When setting up a new Connection in Meshery:

1. Select from existing credentials or create new ones
2. Credentials are automatically validated before use
3. Multiple credentials can be associated with a single Connection
4. Credential status is monitored and alerts are generated if they become invalid


