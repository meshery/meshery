---
name: relationship-engineering
description: Guide for creating and validating Meshery relationship definitions; use for relationship schema/taxonomy/selector requests.
---

# Meshery Relationship Engineering Skill

## Overview
This skill enables Claude to create, validate, and refine Meshery relationship definitions for cloud native infrastructure models. It provides expertise in the CNCF Meshery v1alpha3 relationship schema, taxonomy, and selector specifications.

## Triggers
Use this skill when the user:
- Mentions "Meshery", "relationship definition", or "component relationships"
- Provides Kubernetes CRDs and asks about relationships
- Requests validation of Meshery JSON definitions
- Asks about cloud service dependencies (AWS, Azure, GCP) in Meshery context
- Needs help with relationship taxonomy, selectors, or schema compliance
- References "v1alpha3", "relationship schema", or "Meshery models"

## Core Competencies

### 1. Relationship Taxonomy Mastery

**Hierarchical Relationships (Containment & Ownership)**
- **Inventory Subtype**: Component membership or containment
  - Examples: Container in Pod, Subnet in VPC, Node Pool in Cluster
  - Indicators: "contains", "member of", "inside", "part of"
  - Use when: One component is physically/logically contained within another

- **Parent/Child Subtype**: Lifecycle ownership and management
  - Examples: Deployment manages ReplicaSet, Operator creates Custom Resources
  - Indicators: "manages", "creates", "owns", "controls lifecycle"
  - Use when: Controller/operator pattern where parent creates and manages child

**Edge Relationships (Functional Connections)**
- **Network Subtype**: Traffic flow and network connectivity
  - Examples: Service → Pod, Ingress → Service, SecurityGroup allows Instance
  - Indicators: "routes to", "connects to", "allows traffic", "exposes"
  - Use when: Network communication or traffic flow exists

- **Mount Subtype**: Storage or configuration attachment
  - Examples: Pod mounts PVC, Lambda mounts EFS, ConfigMap mounted to Pod
  - Indicators: "mounts", "attaches", "volumes", "filesystem"
  - Use when: Storage or configuration is attached/mounted

- **Binding Subtype**: Access control, policy, or configuration association
  - Examples: ClusterRoleBinding → ServiceAccount + ClusterRole, ServiceMonitor → Service
  - Indicators: "binds", "grants access", "applies policy", "associates", "selects for monitoring"
  - Use when: RBAC, policy application, or configuration selection

**Sibling Relationships (Peer Components)**
- **Default Subtype**: General peer relationships
  - Examples: Components sharing parent, representing same entity in different states
  - Indicators: "related to", "peer of", "associated with"
  - Use when: Logical relationship without dependency or hierarchy

### 2. Validation Framework

**Schema Compliance Checks:**
```
Required Fields:
- kind (must be "hierarchical", "edge", or "sibling")
- subType (must match allowed values for the kind)
- version (must be "v1alpha3")
- model (from and to components)
- selectors (from and to with proper structure)
- metadata (name, description)

Type Validation:
- Ensure all field types match schema (strings, objects, arrays)
- Verify enum values are from allowed lists
- Check nested object structures
```

**Selector Integrity Checks:**
```
Precision Analysis:
- Are selectors too broad? (e.g., matching all pods instead of specific ones)
- Are selectors too narrow? (e.g., missing valid matches)
- Do label selectors use proper matchLabels/matchExpressions syntax?
- Are model and kind constraints present and accurate?

Selector Components to Validate:
- from.model, from.kind
- to.model, to.kind  
- from.match (labels, annotations, identifiers)
- to.match (labels, annotations, identifiers)
- deny selectors (if present)
```

**Taxonomic Correctness:**
```
Validation Questions:
- Does the subtype match the actual behavior?
  ❌ Wrong: Using Inventory for Service → Pod (should be Network)
  ✅ Right: Using Inventory for Container → Pod
  
- Is the relationship direction correct?
  - Parent/Child: Parent in "from", Child in "to"
  - Network: Source in "from", Destination in "to"
  - Mount: Consumer in "from", Resource in "to"
  
- Are we using a supported subtype?
  - Never invent new subtypes
  - Only use supported Meshery subTypes (examples: inventory, parent, network, mount, binding, sibling, reference, permission, firewall, alias, wallet, matchlabels)
```

**Logic & Cyclic Checks:**
```
- Does this create infinite loops in graph visualization?
- Are there conflicting relationships (A→B and B→A with wrong types)?
- Does the relationship make semantic sense?
```

**Model Packaging Validation (Required for new/updated relationships):**
```
Follow: Using mesheryctl to modify and verify a relationship.md

From the Meshery repo root, package the model using mesheryctl:
mesheryctl model build <model-name>/<model-version> --path server/meshmodel

Example:
mesheryctl model build aws-s3-controller/v1.3.0 --path server/meshmodel

This step should fail fast if the model folder structure and/or model artifacts are invalid, and serves as a functional validation that the relationship definitions are packageable into a model artifact.
```

### 3. Generation Workflow

**Step 1: Analysis**
```
Given a component definition or CRD:
1. Identify the resource's purpose and function
2. Examine spec fields for references to other resources
3. Check status fields for runtime relationships
4. Review metadata for labels, annotations, owner references
5. Understand the resource lifecycle and dependencies
```

**Step 2: Relationship Identification**
```
Ask critical questions:
- What does this component contain? → Inventory
- What does this component create/manage? → Parent/Child
- What does this component connect to? → Network
- What does this component mount? → Mount
- What does this component bind/select? → Binding
- What are its peers? → Sibling
```

**Step 3: Selector Design**
```
Create precise selectors:
- Use specific model and kind constraints
- Add label/annotation matchers when available
- Consider deny selectors for exclusions
- Test selector scope mentally (what would match?)
```

**Step 4: JSON Generation**
```json
{
  "schemaVersion": "relationships.meshery.io/v1alpha3",
  "version": "v1.0.0",
  "kind": "hierarchical|edge|sibling",
  "type": "non-binding",
  "subType": "inventory|parent|network|mount|binding|sibling",
  "metadata": {
    "description": "Clear description of what this relationship represents"
  },
  "model": {
    "name": "kubernetes",
    "version": "v1.25.0",
    "model": {
      "version": "v1.0.0"
    }
  },
  "selectors": [
    {
      "allow": {
        "from": [
          {
            "kind": "SpecificKind",
            "model": {
              "name": "specific-model"
            },
            "match": {
              "refs": [
                // Reference patterns
              ]
            }
          }
        ],
        "to": [
          {
            "kind": "TargetKind",
            "model": {
              "name": "target-model"
            },
            "match": {
              "refs": [
                // Reference patterns
              ]
            }
          }
        ]
      },
      "deny": {
        "from": [],
        "to": []
      }
    }
  ]
}
```

**Step 5: Validation**
```
Run through all validation checks:
✓ Schema compliance
✓ Selector integrity  
✓ Taxonomic correctness
✓ Logic verification
✓ No cycles
```

### 4. Common Patterns Reference

**Kubernetes Common Relationships:**
```
Service → Pod (Edge/Network)
- Service selects pods via label selector
- Network traffic routing relationship

Deployment → ReplicaSet (Hierarchical/Parent)
- Deployment creates and manages ReplicaSet
- Lifecycle ownership

Pod → PersistentVolumeClaim (Edge/Mount)
- Pod mounts PVC for storage
- Volume attachment relationship

ServiceAccount → Role via RoleBinding (Edge/Binding)
- RoleBinding grants ServiceAccount permissions
- RBAC policy application

Namespace → Pod (Hierarchical/Inventory)
- Pod exists within Namespace
- Logical containment
```

**AWS Common Relationships:**
```
VPC → Subnet (Hierarchical/Inventory)
- Subnet is contained within VPC
- Network containment

EC2 Instance → SecurityGroup (Edge/Network)
- SecurityGroup controls instance traffic
- Network policy application

Lambda → VPC (Edge/Network)
- Lambda executes within VPC
- Network context association

RDS Instance → Subnet (Hierarchical/Inventory)
- RDS exists in specific subnet
- Network placement
```

**Prometheus Operator Patterns:**
```
ServiceMonitor → Service (Edge/Binding)
- ServiceMonitor selects Service for scraping
- Monitoring configuration binding

PodMonitor → Pod (Edge/Binding)
- PodMonitor selects Pods for metrics
- Monitoring target selection
```

### 5. Error Patterns to Avoid

**❌ Common Mistakes:**
```
1. Using Inventory for network connections
   Wrong: Service → Pod as Inventory
   Right: Service → Pod as Edge/Network

2. Using Parent/Child without lifecycle ownership
   Wrong: Service → Pod as Parent/Child
   Right: Deployment → ReplicaSet as Parent/Child

3. Overly broad selectors
   Wrong: match all pods in cluster
   Right: match specific label selector

4. Missing model constraints
   Wrong: kind: "Service" without model
   Right: kind: "Service", model: {name: "kubernetes"}

5. Incorrect relationship direction
   Wrong: Pod → Service (backward)
   Right: Service → Pod (service selects pods)

6. Inventing new subtypes
   Wrong: subType: "reference"
   Right: Use existing subtypes only
```

## Operation Modes

### Mode 1: Generation
**Input:** Component definition, CRD, or service description
**Output:** Complete, validated relationship definition JSON
**Process:** Analysis → Classification → Generation → Validation

### Mode 2: Validation
**Input:** Existing relationship definition JSON
**Output:** Validation report with errors, warnings, suggestions
**Process:** Schema check → Selector check → Taxonomy check → Logic check

### Mode 3: Refinement
**Input:** Existing relationship with issues
**Output:** Corrected definition with explanation of changes
**Process:** Identify issues → Apply corrections → Validate → Explain

### Mode 4: Explanation
**Input:** Question about taxonomy or relationships
**Output:** Clear explanation with examples
**Process:** Clarify concept → Provide examples → Contrast alternatives

## Response Format

### For Generation Requests:
```
1. Brief analysis of component relationships
2. Valid JSON definition(s)
3. Optional: Reasoning if requested
```

### For Validation Requests:
```
1. Validation status (✓ Valid / ✗ Invalid)
2. Specific errors with locations
3. Warnings or suggestions
4. Corrected version if needed
```

### For Explanation Requests:
```
1. Concept explanation
2. Concrete examples
3. Common mistakes to avoid
4. When to use each approach
```

## Key Principles

1. **Never invent subtypes** - Only use supported v1alpha3 subtypes
2. **Precision over breadth** - Selectors should be as specific as needed
3. **Validate rigorously** - Check schema, selectors, taxonomy, and logic
4. **Cross-reference schemas** - Always verify against v1alpha3 specifications
5. **Semantic accuracy** - Subtype must match actual component behavior
6. **Direction matters** - Ensure from/to direction is semantically correct
7. **Avoid cycles** - Check for potential infinite loops in graphs

## Domain Knowledge Requirements

- Kubernetes APIs, controllers, and operator patterns
- Cloud provider service architectures (AWS, Azure, GCP)
- Network, storage, and security relationship patterns
- Meshery model constructs and component schemas
- CNCF ecosystem tools (Prometheus, service meshes, etc.)

## Success Criteria

A successful relationship definition:
- ✓ Validates against v1alpha3 schema
- ✓ Uses correct Kind → Type → Subtype taxonomy
- ✓ Has precise, appropriate selectors
- ✓ Makes semantic sense for the components involved
- ✓ Won't create false positive matches
- ✓ Won't miss valid relationship instances
- ✓ Has clear, accurate metadata
- ✓ Follows correct directional conventions
