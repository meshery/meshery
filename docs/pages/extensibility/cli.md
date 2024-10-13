---
layout: default
title: "Extensibility: CLI Plugins"
permalink: extensibility/cli-plugins
type: Extensibility
abstract: "Meshery CLI plugins extend the functionality of the Meshery command-line interface, allowing users to interact with various Meshery features and integrations directly from their terminal."
language: en
list: include
---

# Meshery CLI Plugins

Meshery CLI plugins extend the functionality of the Meshery command-line interface, allowing users to interact with various Meshery features and integrations directly from their terminal.

## Available CLI Plugins

Meshery currently supports the following CLI plugin:

- helm-kanvas-snapshot

The following plugins are planned but not yet implemented:

- kubectl-kanvas-snapshot
- kubectl-meshsync-snapshot

### helm-kanvas-snapshot

The Helm Kanvas Snapshot plugin allows users to create visual snapshots of Helm charts.

#### Installation

{% capture code_content %}
helm plugin install https://github.com/meshery/helm-kanvas-snapshot
{% endcapture %}
{% include code.html code=code_content %}

#### Usage

{% capture code_content %}
helm snapshot -f <chart-URI> [-n <snapshot-name>] [-e <email>]
{% endcapture %}
{% include code.html code=code_content %}

For more information and detailed usage instructions, please refer to the [helm-kanvas-snapshot repository](https://github.com/meshery/helm-kanvas-snapshot).

#### Design Specification

The helm-kanvas-snapshot plugin is designed with the following goals and objectives:

1. Allow Helm users to readily visualize their Helm charts
   - Create visual snapshots of packaged Helm charts
   - Create visual snapshots of unpackaged Helm charts
2. Offer choice of snapshot delivery
   - Synchronous: Return snapshot URL in the terminal
   - Asynchronous: Send snapshot URL via email
3. Provide engagement beyond the snapshot
   - Highlight Meshery's ability to deploy the Helm chart

For a detailed understanding of the plugin's architecture, system flows, and implementation details, please refer to the [helm-kanvas-snapshot design specification](https://docs.google.com/document/d/1NdgLoOS3Xy1Z3vwB1dnZ9ETyKCrSinVgAzL8UO3UeQk/edit).

#### Design Overview

The helm-kanvas-snapshot plugin is designed with the following key aspects:

1. **Purpose**: Generates visual snapshots of Helm charts as Meshery Snapshots, integrating with Meshery Cloud and GitHub Actions.

2. **Core Functionality**: 
   - Implemented in Go (version 1.21.8)
   - Uses the github.com/layer5io/meshkit library for core functionality
   - Leverages github.com/spf13/viper for configuration management
   - Utilizes github.com/sirupsen/logrus for logging

3. **Plugin Structure**: Follows the Helm plugin structure with a plugin.yaml file defining metadata and commands. Main logic is implemented in cmd/kanvas-snapshot/main.go.

4. **Features**:
   - Supports both synchronous and asynchronous snapshot delivery
   - Can handle packaged Helm charts
   - Integrates with Meshery's snapshot functionality

5. **Command-line Interface**: Invoked as a Helm plugin (`helm snapshot`) with various flags for customization.

6. **Configuration**: Uses Viper for configuration management, allowing for flexible configuration options.

### kubectl-kanvas-snapshot (Planned)

**Note: This plugin is currently in development and has not yet been created. The installation instructions provided here are provisional and subject to change upon release.**

The kubectl Kanvas Snapshot plugin is planned to be a native kubectl plugin for creating visual snapshots of Kubernetes manifest files.

#### Planned Installation

{% capture code_content %}
kubectl krew index add kanvas-snapshot https://github.com/meshery/kubectl-kanvas-snapshot.git
kubectl krew install kanvas-snapshot
{% endcapture %}
{% include code.html code=code_content %}

### kubectl-meshsync-snapshot (Planned)

**Note: This plugin is currently in development and has not yet been created. The installation instructions provided here are provisional and subject to change upon release.**

The kubectl MeshSync Snapshot plugin is planned to work with Meshery's MeshSync component, allowing users to create visual snapshots of the current state of their Kubernetes cluster as seen by MeshSync.

#### Planned Installation

{% capture code_content %}
kubectl krew index add meshsync-snapshot https://github.com/meshery/kubectl-meshsync-snapshot.git
kubectl krew install meshsync-snapshot
{% endcapture %}
{% include code.html code=code_content %}

## Designing Custom CLI Plugins for Meshery

Meshery CLI plugins are designed to extend Meshery's functionality and provide custom use-cases. Developers can create new CLI plugins to add features or integrate with other tools.

### Plugin Development Guidelines

When developing a new CLI plugin for Meshery, consider the following:

- Use a consistent naming convention (e.g., mesheryctl-plugin-name or kubectl-meshery-plugin-name).
- Implement proper error handling and user-friendly output.
- Follow Meshery's coding standards and best practices.
- Provide clear documentation on installation, usage, and any dependencies.

## Contributing

Contributions to Meshery CLI plugins are welcome. To contribute:

1. Fork the respective plugin repository.
2. Create a feature branch (`git checkout -b feature/Feature`).
3. Commit your changes (`git commit -m 'Add some Feature'`).
4. Push to the branch (`git push origin feature/Feature`).
5. Open a Pull Request.

For more detailed information, please refer to the CONTRIBUTING.md file in each plugin's repository.

## Additional Resources

- [Meshery Documentation](https://docs.meshery.io/)
- [Meshery Community Forum](https://discuss.layer5.io/)
- [Layer5 Blog](https://layer5.io/blog)