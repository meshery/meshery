---
name: Model Relationships Test Agent
description: How to create new relationship test fixture designs for Meshery
tools: ['search/changes', 'search/codebase', 'edit/editFiles', 'vscode/extensions', 'web/fetch', 'web/githubRepo', 'vscode/getProjectSetupInfo', 'vscode/installExtension', 'vscode/newWorkspace', 'vscode/runCommand', 'vscode/openSimpleBrowser', 'read/problems', 'execute/getTerminalOutput', 'execute/runInTerminal', 'read/terminalLastCommand', 'read/terminalSelection', 'execute/createAndRunTask', 'execute', 'execute/runTask', 'execute/runTests', 'search', 'search/searchResults', 'execute/testFailure', 'search/usages', 'vscode/vscodeAPI', 'github/*', 'memory']
---

# Relationship Test Fixture Design Guide

This document describes the process to create new relationship test fixtures for the Meshery relationship evaluation integration tests.

## Overview

Relationship test fixtures are JSON files that define a **Meshery Design** containing:
1. **Components** - Kubernetes resources (e.g., Deployment, Service, ConfigMap)
2. **Relationships** - The expected relationships that should be identified between components

These fixtures are used by `ui/tests/e2e/relationship_evaluation.spec.js` to verify that the relationship evaluation API correctly identifies relationships between components.

---

## Step 1: Identify a Relationship Definition to Test

Browse the relationship definitions in:
```
/server/meshmodel/kubernetes/v1.35.0-beta.0/v1.0.0/relationships/
```

Each file defines a relationship with:
- **kind**: `edge` or `hierarchical`
- **type**: `binding`, `non-binding`, or `parent`
- **subType**: `network`, `reference`, `permission`, `inventory`, `alias`, `wallet`, etc.
- **selectors**: Defines `from` and `to` components with `mutatorRef` and `mutatedRef` paths

### Key Fields to Understand

| Field | Description |
|-------|-------------|
| `mutatorRef` | Path in the source component providing the value |
| `mutatedRef` | Path in the target component receiving/matching the value |

**Example** (Secret → Deployment reference):
```json
{
  "kind": "edge",
  "type": "non-binding", 
  "subType": "reference",
  "selectors": [{
    "allow": {
      "from": [{ "kind": "Secret", "patch": { "mutatorRef": [["displayName"]] }}],
      "to": [{ "kind": "Deployment", "patch": { "mutatedRef": [["configuration","spec","template","spec","containers","0","envFrom","0","secretRef","name"]] }}]
    }
  }]
}
```

---

## Step 2: Create the Fixture JSON File

### File Naming Convention
```
{component1}-{component2}-{kind}-{type}-{subType}-fixture.json
```

**Examples:**
- `secret-deployment-edge-non-binding-reference-fixture.json`
- `namespace-hierarchical-parent-inventory-fixture.json`
- `role-rolebinding-serviceaccount-edge-binding-permission-fixture.json`

### File Location
```
ui/tests/e2e/fixtures/relationships/
```

---

## Step 3: Define the Fixture Structure

### Basic Structure

```json
{
    "id": "<unique-uuid>",
    "name": "<descriptive-fixture-name>",
    "schemaVersion": "designs.meshery.io/v1beta1",
    "version": "0.0.1",
    "metadata": {},
    "components": [
        // Component definitions go here
    ],
    "relationships": [
        // Expected relationship definitions go here
    ]
}
```

### Component Definition Template

```json
{
    "id": "<unique-component-uuid>",
    "schemaVersion": "components.meshery.io/v1beta1",
    "version": "v1.0.0",
    "displayName": "<component-display-name>",
    "description": "<component-description>",
    "format": "JSON",
    "modelReference": {
        "version": "v1.0.0",
        "name": "kubernetes",
        "displayName": "Kubernetes",
        "id": "00000000-0000-0000-0000-000000000000",
        "registrant": { "kind": "github" },
        "model": { "version": "v1.35.0-rc.0" }
    },
    "styles": {
        "background-opacity": 1,
        "primaryColor": "#326CE5",
        "secondaryColor": "#7aa1f0",
        "shape": "round-rectangle"
    },
    "capabilities": null,
    "status": "enabled",
    "metadata": {
        "isAnnotation": false,
        "isNamespaced": true
    },
    "configuration": {
        // Kubernetes resource spec goes here
    },
    "component": {
        "version": "<api-version>",
        "kind": "<Kind>",
        "schema": ""
    }
}
```

---

## Step 4: Configure Components to Match Relationship Selectors

**Critical**: The component configurations MUST satisfy the `mutatorRef` and `mutatedRef` paths defined in the relationship selector.

### Example: Secret → Deployment Reference

**Relationship Definition requires:**
- `mutatorRef`: `["displayName"]` on Secret
- `mutatedRef`: `["configuration","spec","template","spec","containers","0","envFrom","0","secretRef","name"]` on Deployment

**Secret Component:**
```json
{
    "id": "44444444-4444-4444-4444-444444444444",
    "displayName": "my-app-secret",  // <- matches mutatorRef
    "component": { "version": "v1", "kind": "Secret" },
    "configuration": {
        "metadata": { "name": "my-app-secret", "namespace": "default" },
        "type": "Opaque",
        "data": { "username": "YWRtaW4=" }
    }
}
```

**Deployment Component:**
```json
{
    "id": "55555555-5555-5555-5555-555555555555",
    "displayName": "my-app-deployment",
    "component": { "version": "apps/v1", "kind": "Deployment" },
    "configuration": {
        "spec": {
            "template": {
                "spec": {
                    "containers": [{
                        "name": "my-app-container",
                        "envFrom": [{
                            "secretRef": {
                                "name": "my-app-secret"  // <- matches mutatedRef, value equals displayName
                            }
                        }]
                    }]
                }
            }
        }
    }
}
```

---

## Step 5: Define the Expected Relationship

```json
{
    "id": "<unique-relationship-uuid>",
    "evaluationQuery": null,
    "kind": "<edge|hierarchical>",
    "metadata": {
        "description": "<relationship-description>",
        "styles": { "primaryColor": "", "svgColor": "", "svgWhite": "" },
        "isAnnotation": false
    },
    "model": {
        "version": "v1.0.0",
        "name": "kubernetes",
        "displayName": "Kubernetes",
        "id": "00000000-0000-0000-0000-000000000000",
        "registrant": { "kind": "github" },
        "model": { "version": "v1.35.0-rc.0" }
    },
    "schemaVersion": "relationships.meshery.io/v1alpha3",
    "selectors": [{
        "allow": {
            "from": [{
                "id": "<from-component-id>",
                "kind": "<FromKind>",
                "model": { "name": "kubernetes", "registrant": { "kind": "github" }},
                "patch": { "patchStrategy": "replace", "mutatorRef": [[...]] }
            }],
            "to": [{
                "id": "<to-component-id>",
                "kind": "<ToKind>",
                "model": { "name": "kubernetes", "registrant": { "kind": "github" }},
                "patch": { "patchStrategy": "replace", "mutatedRef": [[...]] }
            }]
        }
    }],
    "subType": "<subtype>",
    "status": "pending",
    "type": "<type>",
    "version": ""
}
```

> [!IMPORTANT]
> The `selectors.allow.from[].id` and `selectors.allow.to[].id` MUST match the component IDs defined in your fixture.

---

## Step 6: Register the Fixture

Update `ui/tests/e2e/fixtures/relationships/index.js`:

```javascript
import myNewFixture from './my-new-fixture.json';

export const RelationshipTestFixtures = [
    // ... existing fixtures
    myNewFixture
];
```

---

## Step 7: Add to Test Designs

Update `ui/tests/e2e/relationship_evaluation.spec.js` to include the fixture:

```javascript
const DESIGNS_TO_TEST = [
    // ... existing designs
    {
        id: "<fixture-id>",
        name: "<fixture-name>"
    }
]
```

---

## Common Relationship Types Reference

| Kind | Type | SubType | Example Use Case |
|------|------|---------|------------------|
| `edge` | `non-binding` | `reference` | ConfigMap/Secret → Pod/Deployment |
| `edge` | `non-binding` | `network` | Service → Deployment (label selector) |
| `edge` | `binding` | `permission` | Role → RoleBinding → ServiceAccount |
| `edge` | `binding` | `mount` | PV → PVC mount |
| `hierarchical` | `parent` | `inventory` | Namespace → namespaced resources |
| `hierarchical` | `parent` | `alias` | Container alias within Deployment |
| `hierarchical` | `parent` | `wallet` | EndpointSlice → Service |

---

## Validation Checklist

- [ ] Unique UUIDs for design, components, and relationships
- [ ] Component IDs in relationship selectors match component definitions
- [ ] `displayName` values match where `mutatorRef` references `["displayName"]`
- [ ] Configuration paths exist and match `mutatedRef` paths exactly
- [ ] Relationship `status` is set to `"pending"` or `"approved"`
- [ ] File is valid JSON (use `jq .` to validate)
- [ ] Fixture is registered in `index.js`
- [ ] Design is added to `DESIGNS_TO_TEST` in the spec file

---

## Running Tests

```bash
// turbo
cd ui && npx playwright test relationship_evaluation.spec.js --grep "@relationship"
```

This will test all relationship fixtures and report which relationships were correctly identified by the evaluation API.
