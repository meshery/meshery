---
layout: default
title: Credentials
permalink: concepts/logical/credentials
type: concepts
abstract: Meshery uses one or more Credentials when authenticating to a managed or unmanaged Connection.
language: en
list: include
redirect_from:
- concepts/credentials
---
Meshery uses one or more Credentials when authenticating to a managed or unmanaged Connection. Credentials are based on the Meshery's [Credential Schema](https://github.com/meshery/schemas/blob/master/openapi/schemas/credentials.yml) with defined attributes.

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

- Credentials are encrypted at rest
- Access is controlled through fine-grained permissions
- Credentials are never exposed in logs or API responses
- Support for secret management integration

### Using Credentials with Connections

When setting up a new Connection in Meshery:

1. Select from existing credentials or create new ones
2. Credentials are automatically validated before use
3. Multiple credentials can be associated with a single Connection
4. Credential status is monitored and alerts are generated if they become invalid

For detailed instructions on managing credentials, see the [Credential Management Guide](/guides/credential-management).


