---
title: Relationship Definition Guide
description: Learn how to define relationships between cloud-native components with practical examples
weight: 4
---

# Relationship Definition Guide

Relationships in Meshery define how components interact and depend on each other. This guide teaches you to create meaningful, maintainable relationships.

## Quick Start

A relationship is a JSON file that defines:
- **From & To**: Which components connect
- **Kind**: Type of relationship (binding, hierarchical, reference)
- **Selectors**: Conditions that trigger the relationship

## Example: IAM Role to ServiceAccount (IRSA)

This real-world example shows the AWS IAM Roles for Service Accounts pattern:

```json
{
  "name": "IAM Role to Kubernetes ServiceAccount (IRSA)",
  "description": "AWS IAM permissions for pod authentication",
  "kind": "binding",
  "type": "edge",
  "subType": "security",
  "model": "aws-iam-controller",
  "modelVersion": "v1.0.0",
  "componentVersions": {
    "from": {
      "model": "aws-iam-controller",
      "apiVersion": "v1.0.0",
      "name": "Role"
    },
    "to": {
      "model": "kubernetes",
      "apiVersion": "v1",
      "name": "ServiceAccount"
    }
  },
  "evalProperties": {
    "direction": "bidirectional",
    "allowMultipleConnections": false
  },
  "selectors": {
    "allow": [
      {
        "from": {
          "kind": "Role",
          "metadata": {
            "annotations": {
              "key": "eks.amazonaws.com/role-arn"
            }
          }
        },
        "to": {
          "kind": "ServiceAccount",
          "metadata": {
            "annotations": {
              "key": "eks.amazonaws.com/role-arn"
            }
          }
        }
      }
    ]
  },
  "metadata": {
    "isNamespaced": true,
    "supportsWildcard": false
  }
}
```

## Relationship Types

### 1. Binding Relationships
**Use when:** One component depends on another for functionality

**Example:** IAM Role depends on ServiceAccount for permission mapping

```json
{
  "kind": "binding",
  "type": "edge",
  "direction": "bidirectional"
}
```

### 2. Hierarchical Relationships
**Use when:** Components have parent-child or containment relationships

**Example:** VPC contains Subnets

```json
{
  "kind": "hierarchical",
  "type": "inventory",
  "subType": "parent",
  "direction": "directed"
}
```

### 3. Reference Relationships (Visual Only)
**Use when:** Components interact but don't have binding constraints

**Example:** Lambda function reads from S3 bucket

```json
{
  "kind": "reference",
  "type": "edge",
  "subType": "visual",
  "direction": "directed"
}
```

## Field Reference

| Field | Purpose | Example |
|-------|---------|---------|
| `kind` | Type of relationship | `"binding"`, `"hierarchical"`, `"reference"` |
| `type` | Edge or inventory | `"edge"`, `"inventory"` |
| `subType` | Category | `"security"`, `"network"`, `"compute"`, `"storage"` |
| `direction` | Flow direction | `"directed"`, `"bidirectional"` |
| `model` | Source component model | `"aws-iam-controller"` |
| `selectors` | Matching conditions | See examples below |

## Selector Patterns

### Match by Kind
```json
"selectors": {
  "allow": [
    {
      "from": { "kind": "Role" },
      "to": { "kind": "ServiceAccount" }
    }
  ]
}
```

### Match by Annotation
```json
"selectors": {
  "allow": [
    {
      "from": {
        "kind": "Role",
        "metadata": {
          "annotations": {
            "key": "eks.amazonaws.com/role-arn"
          }
        }
      },
      "to": { "kind": "ServiceAccount" }
    }
  ]
}
```

### Match by Label
```json
"selectors": {
  "allow": [
    {
      "from": { "kind": "Pod" },
      "to": {
        "kind": "PersistentVolume",
        "metadata": {
          "labels": {
            "key": "storage-type",
            "value": "ebs"
          }
        }
      }
    }
  ]
}
```

## Best Practices

✅ **DO:**
- **Be specific with versions**: Always specify exact API versions
- **Use meaningful subTypes**: Choose from: security, network, compute, storage, identity, etc.
- **Document purpose**: Include clear description of what the relationship represents
- **Test in Meshery UI**: Verify relationship appears as expected

❌ **DON'T:**
- **Mix multiple relationships in one file**: Each relationship = one file
- **Use overly broad selectors**: Be specific about matching conditions
- **Ignore namespace context**: Set `isNamespaced` correctly

## File Naming Convention

Use descriptive filenames that clearly indicate:
- relationship type
- source component
- target component

Example naming pattern:

```text
edge-{kind}-{from}-to-{to}.json
```

Examples:

```text
edge-binding-role-to-serviceaccount.json
edge-binding-sg-to-networkpolicy.json
edge-reference-lambda-to-s3.json
```

## Common AWS-Kubernetes Relationships

These are frequently needed patterns:

| From | To | Kind | SubType | Use Case |
|------|----|----|---------|----------|
| IAM Role | ServiceAccount | binding | security | IRSA pattern |
| SecurityGroup | NetworkPolicy | binding | network | Network governance |
| ALB | Ingress | binding | network | Load balancing |
| EBS Volume | PersistentVolume | binding | storage | Persistent storage |
| Secrets Manager | Secret | binding | security | External secrets |
| EC2 Instance | StatefulSet | binding | compute | Node-to-workload |

## Contributing a New Relationship

1. Research existing relationships to avoid duplicates
2. Define selectors and relationship behavior
3. Validate the JSON structure
4. Add a clear description and use case
5. Create a sample architecture in Meshery Playground or Kanvas
6. Attach screenshots of the relationship flow in your PR
7. Reference Issue #17096 when applicable

## Resources

- [Contributing to Models](../contributing-models/)
- [Meshery Relationships Docs](https://docs.meshery.io/concepts/logical/relationships)
- [Issue #17096: AWS Relationships Epic](https://github.com/meshery/meshery/issues/17096)