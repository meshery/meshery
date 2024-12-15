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

| Plugin Name | GitHub Repository | Description |
|------------|------------------|-------------|
| helm-kanvas-snapshot | [GitHub Repo](https://github.com/meshery/helm-kanvas-snapshot) | Create visual snapshots of Helm charts |
| kubectl-kanvas-snapshot | [GitHub Repo](https://github.com/meshery/kubectl-kanvas-snapshot) | Native kubectl plugin for visual snapshots of Kubernetes manifests |
| kubectl-meshsync-snapshot | [GitHub Repo](https://github.com/meshery/kubectl-meshsync-snapshot) | Create visual snapshots of Kubernetes cluster state using MeshSync |

## Available CLI Plugins

{% assign plugin_collections = site.pages | where_exp: "item", "item.path contains '/pages/extensions/'" | where: "type", "extensions" %}

{% if plugin_collections.size > 0 %}
### Current Plugins

{% for plugin in plugin_collections limit:5 %}
- [{{ plugin.title }}]({{ plugin.url }})
  - {{ plugin.abstract | default: 'No description available' }}
{% endfor %}

{% if plugin_collections.size > 5 %}
#### More Plugins

{% for plugin in plugin_collections offset:5 %}
- [{{ plugin.title }}]({{ plugin.url }})
{% endfor %}
{% endif %}

{% else %}
### Current Plugins
No CLI plugins currently available.
{% endif %}

## Designing Custom CLI Plugins for Meshery

Meshery CLI plugins are designed to extend Meshery's functionality and provide custom use cases. Developers can create new CLI plugins to add features or integrate with other tools.

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

For more detailed information, refer to the CONTRIBUTING.md file in each plugin's repository.

## Additional Resources

- [Meshery Documentation](https://docs.meshery.io/)
- [Meshery Community Forum](https://discuss.layer5.io/)
- [Layer5 Blog](https://layer5.io/blog)