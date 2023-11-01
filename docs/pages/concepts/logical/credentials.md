---
layout: default
title: Credentials
permalink: concepts/credentials
type: concepts
abstract: Meshery uses one or more Credentials when authenticating to a managed or unmanaged Connection.
language: en
list: include
---
Meshery uses one or more Credentials when authenticating to a managed or unmanaged Connection.
Credentials are pivotal in Meshery as they serve as the authentication mechanism for connections registered with the Meshery server. They ensure secure interactions between Meshery and the various services and systems Meshery is designed to work with. <br>
For more info, see [Meshery Integrations](/integrations/)


## Key Features

- **Authentication for Connections**
These credentials serve to authenticate the Meshery connection when it communicates with external services.
For example, when Meshery scrapes metrics from a Prometheus instance, it uses the credentials associated with that connection to authenticate itself to the Prometheus server.
<br>
For more information, see "[Connections](/concepts/connections)".

- **Reusability Across Connections**
Credentials are not tied to a specific connection but can be shared between different connections within the same workspace. This reusability ensures that you don't have to create separate credentials for each connection.
For example, if you have two Kubernetes clusters, you can use the same credentials for both clusters. This makes it easier to manage and maintain your credentials.
<br>
For more information, see "[Workspaces](/concepts/workspaces)".

- **Connectivity Testing**
Credentials play a crucial role in connectivity testing within Meshery. When you initiate ad-hoc connectivity tests, the system uses the appropriate credential associated with the connection being tested. If the credential is valid and the target system is reachable, the test succeeds. This ensures the connections are actively monitored and validated.

- **Granular Access Control**
Meshery provides fine-grained access control through credentials. You can specify which credentials are associated with particular connections and workspace, giving you control over who can access and manage specific resources.
<br>
For more information, see "[Workspaces](/concepts/workspaces)".


## Examples

### Kubernetes Connections
In the case of Kubernetes connections, your credentials are typically the certificates or authentication tokens required to access the Kubernetes cluster. When you perform an ad-hoc connectivity test by clicking on the Kubernetes chip, the specific credential associated with that connection is used to reach the Kubernetes API server. If the credential is still valid and the cluster is reachable, the test passes. However, if there's an issue, Meshery provides you with a set of remediation actions to fix the connectivity. This can involve either registering the Kubernetes connection or re-uploading your kubeconfig file.


### GitHub Connections

When you authorize the Meshery GitHub app for services like the MeshMap Snapshot, Meshery stores your specific installation ID. This installation ID is then used for various processes, such as verifying Meshery's access to the repositories you assign during the setup process. Additionally, it's used for tasks like installing GitHub workflows directly into the selected repositories and fetching the chosen repositories from your GitHub organization.
