# Model Templates

This folder contains JSON files defining Meshery model templates.  
Each template describes a component or technology (e.g., Kubernetes) for use within Meshery, specifying its visual appearance, metadata, and supported capabilities.

## Structure

A typical template includes:
- **model**: Name of the technology/component (e.g., "kubernetes")
- **primaryColor**, **secondaryColor**, **shape**: Visual properties for UI representation
- **capabilities**: List of actions, configuration options, and views supported by the component

## Example

See `kubernetes.json` for a complete example.

## Usage

To add or update a model template in Meshery, place your JSON file in this directory.  
