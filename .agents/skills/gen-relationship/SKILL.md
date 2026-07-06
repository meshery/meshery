---
name: gen-relationship
description: Generate schema-backed relationship definitions for Meshery models between given components.
tools: ['search/changes', 'search/codebase', 'edit/editFiles', 'vscode/extensions', 'web/fetch', 'web/githubRepo', 'vscode/getProjectSetupInfo', 'vscode/installExtension', 'vscode/newWorkspace', 'vscode/runCommand', 'vscode/openSimpleBrowser', 'read/problems', 'execute/getTerminalOutput', 'execute/runInTerminal', 'read/terminalLastCommand', 'read/terminalSelection', 'execute/createAndRunTask', 'execute', 'execute/runTask', 'execute/runTests', 'search', 'search/searchResults', 'execute/testFailure', 'search/usages', 'vscode/vscodeAPI', 'github/*', 'memory']
---

# Skill: gen-relationship

Generate schema-backed relationship definitions for Meshery models between given components. This skill ensures that relationships are mathematically accurate, logically sound, and strictly compliant with the canonical `v1beta1` relationship schema.

## Usage

Invoke this skill when you need to define a new relationship between two or more components:
- `/gen-relationship Define a reference relationship between a Deployment and a PersistentVolumeClaim in the kubernetes model.`
- `/gen-relationship Create a hierarchical parent-inventory relationship between a ServiceMesh and its ControlPlane.`

## Instructions

1.  **Identify the Components**: Determine the source (`from`) and target (`to`) components. Note their `kind` and the `model` they belong to.
2.  **Determine the Relationship Type**: Select the appropriate `kind`, `type`, and `subType` based on the desired interaction.
3.  **Formulate Selectors**: Construct the `selectors` array with `allow` (and optionally `deny`) blocks.
4.  **Define Patches**: If the relationship involves data flow or configuration binding, define `mutatorRef` and `mutatedRef` as nested arrays of strings.
5.  **Assign Metadata**: Include UI capabilities and styles to ensure the relationship is rendered correctly in MeshMap.
6.  **Validate against Schema**: Ensure the output matches the `v1beta1` schema structure rigorously.

## Technical Context

### Relationship Terminologies

- **Kind**:
    - `hierarchical`: For parent-child or ownership interactions.
    - `edge`: For peer-to-peer or connective interactions.
    - `sibling`: For components that occupy the same logical level or group.
- **Type**:
    - `Parent`: Defines a strong hierarchical bond where the child belongs to the parent (child is usually discovered within the parent's inventory).
    - `Binding`: Defines a functional connection between components.
    - `non-binding`: Defines a logical or informational connection without functional enforcement.
- **SubType**:
    - `inventory`: Hierarchical mapping of resources (e.g., used in `Parent` type).
    - `reference`: Binding based on a property reference (e.g., a Deployment referencing a ConfigMap).
    - `network`: Connectivity mapping (L3/L4/L7 interactions).
    - `permission`: Authorization or access control mapping.
    - `wallet`: Cloud billing or cost mapping.

### Canonical Schema (v1beta1)

- **schemaVersion**: `relationships.meshery.io/v1beta1`
- **Selectors**:
    - `allow`: Contains `from` and `to` arrays of component selectors.
    - `deny`: (Optional) Excludes specific components from the relationship.
- **Selector Item**:
    - `kind`: The kind of the component (e.g., `Deployment`).
    - `model`: The model name (e.g., `kubernetes`).
    - `patch`: (Optional) Specifies how properties are mapping.
- **Path References**:
    - `mutatorRef`: Source path, formatted as `[["path", "to", "field"]]`.
    - `mutatedRef`: Target path, formatted as `[["path", "to", "field"]]`.
    - **CRITICAL**: Both must be nested arrays of strings (e.g., `[["settings", "name"]]`) to pass upstream evaluation.
- **UI Metadata**:
    - `metadata.capabilities.designer.edit`: Must be `true` for a relationship to be editable in MeshMap.
    - `description`: A clear, human-readable description of the relationship's purpose.

## Examples

### 1. Peer-to-Peer Reference (Edge)
**Reasoning**: The Deployment is the originator of the reference, containing the field that must match the source PVC. The PVC provides its identity (mutator) to the Deployment's volume configuration (mutated).

```json
{
  "schemaVersion": "relationships.meshery.io/v1beta1",
  "kind": "edge",
  "type": "non-binding",
  "subType": "reference",
  "metadata": {
    "description": "Deployment referencing a PVC via claimName",
    "capabilities": {
      "designer": { "edit": true }
    }
  },
  "selectors": [
    {
      "allow": {
        "from": [
          {
            "kind": "Deployment",
            "model": "kubernetes",
            "patch": {
              "mutatedRef": [["configuration", "spec", "template", "spec", "volumes", "0", "persistentVolumeClaim", "claimName"]]
            }
          }
        ],
        "to": [
          {
            "kind": "PersistentVolumeClaim",
            "model": "kubernetes",
            "patch": {
             "mutatorRef": [["configuration", "metadata", "name"]]
            }
          }
        ]
      }
    }
  ]
}
```

### 2. Hierarchical Inventory (Parent)
**Reasoning**: The Namespace is the authoritative parent defining the logical scope. It provides its name (mutator) to be injected into the namespace field of all child components (mutated).

```json
{
  "schemaVersion": "relationships.meshery.io/v1beta1",
  "kind": "hierarchical",
  "type": "parent",
  "subType": "inventory",
  "metadata": {
    "description": "Namespace-level resource containment",
    "capabilities": {
      "designer": { "edit": true }
    }
  },
  "selectors": [
    {
      "allow": {
        "from": [
          {
            "kind": "Namespace",
            "model": "kubernetes",
            "patch": {
             "mutatorRef": [["configuration", "metadata", "name"]]
             
            }
          }
        ],
        "to": [
          {
            "kind": "*",
            "model": "kubernetes",
            "patch": {
              "mutatedRef": [["configuration", "metadata", "namespace"]]
            }
          }
        ]
      }
    }
  ]
}
```

### 3. Network Connectivity (Edge)
**Reasoning**: The Service is the originator of the network connection, defining selection criteria. It provides its selector labels (mutator) to match or be mirrored by the Deployment's pod labels (mutated).

```json
{
  "schemaVersion": "relationships.meshery.io/v1beta1",
  "kind": "edge",
  "type": "binding",
  "subType": "network",
  "metadata": {
    "description": "Service targeting deployment pods",
    "capabilities": {
      "designer": { "edit": true }
    }
  },
  "selectors": [
    {
      "allow": {
        "from": [
          {
            "kind": "Service",
            "model": "kubernetes",
            "patch": {
              "mutatorRef": [["configuration", "spec", "selector"]]
            }
          }
        ],
        "to": [
          {
            "kind": "Deployment",
            "model": "kubernetes",
            "patch": {
              "mutatedRef": [["configuration", "spec", "template", "metadata", "labels"]]
            }
          }
        ]
      }
    }
  ]
}
```

## Guidelines

1.  **Source of Truth**: Always prioritize `v1beta1` schema rules from `meshery/schemas`. Do not copy existing `v1beta2` or older files blindly as they may contain legacy non-compliant structures.
2.  **Strict Pathing**: Ensure that the `mutatorRef` and `mutatedRef` paths correspond to the actual JSON schema of the components being linked.
3.  **Logical Consistency**: Match `Kinds`, `Types`, and `subTypes` that corresponds logically to the components being linked (e.g. `inventory` for `Parent`, `network` for `Service` connections).
4.  **No Guessing**: If component schemas are unavailable, use `grep_search` to find component definitions in `models` to verify valid paths for `patch` operations.
