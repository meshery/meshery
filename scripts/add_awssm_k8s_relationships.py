#!/usr/bin/env python3
"""
Script to add AWS Secrets Manager to Kubernetes workload relationships.

This script adds relationships between AWS Secrets Manager Secret and the following
Kubernetes workloads:
- StatefulSet
- DaemonSet
- Job
- CronJob

These relationships enable visual secret management in Meshery Kanvas.
"""

import os
import json
import hashlib

BASE_PATH = "/home/sarika/meshery/server/meshmodel/aws-secretsmanager-controller"

# Template for the relationship files
# The key difference from Pod is that these workloads have spec.template.spec.volumes
RELATIONSHIP_TEMPLATE = {
  "id": "00000000-0000-0000-0000-000000000000",
  "evaluationQuery": "",
  "kind": "edge",
  "metadata": {
    "description": "",
    "styles": {
      "primaryColor": "",
      "svgColor": "",
      "svgWhite": ""
    },
    "isAnnotation": False
  },
  "model": {
    "version": "",
    "name": "aws-secretsmanager-controller",
    "displayName": "",
    "id": "00000000-0000-0000-0000-000000000000",
    "registrant": {
      "kind": ""
    },
    "model": {
      "version": ""  # Will be set dynamically
    }
  },
  "schemaVersion": "relationships.meshery.io/v1alpha3",
  "selectors": [
    {
      "allow": {
        "from": [
          {
            "id": None,
            "kind": "",  # Will be set dynamically (StatefulSet, DaemonSet, Job, CronJob)
            "match": {},
            "match_strategy_matrix": None,
            "model": {
              "version": "",
              "name": "kubernetes",
              "displayName": "",
              "id": "00000000-0000-0000-0000-000000000000",
              "registrant": {
                "kind": "artifacthub"
              },
              "model": {
                "version": ""
              }
            },
            "patch": {
              "patchStrategy": "replace",
              "mutatorRef": [
                [
                  "name"
                ],
                [
                  "configuration",
                  "spec",
                  "template",
                  "spec",
                  "volumes"
                ]
              ]
            }
          }
        ],
        "to": [
          {
            "id": None,
            "kind": "Secret",
            "match": {},
            "match_strategy_matrix": None,
            "model": {
              "version": "",
              "name": "aws-secretsmanager-controller",
              "displayName": "",
              "id": "00000000-0000-0000-0000-000000000000",
              "registrant": {
                "kind": "artifacthub"
              },
              "model": {
                "version": ""
              }
            },
            "patch": {
              "patchStrategy": "replace",
              "mutatedRef": [
                [
                  "name"
                ],
                [
                  "configuration",
                  "arn"
                ]
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
  ],
  "subType": "network",
  "status": "enabled",
  "type": "non-binding",
  "version": "v1.0.0"
}

# Workloads to add relationships for
WORKLOADS = ["StatefulSet", "DaemonSet", "Job", "CronJob"]

def generate_unique_suffix(workload_kind, model_version):
    """Generate a unique 5-character suffix based on workload and version."""
    data = f"{workload_kind}-{model_version}-aws-secretsmanager"
    hash_obj = hashlib.md5(data.encode())
    return hash_obj.hexdigest()[:5]

def create_relationship(workload_kind, model_version):
    """Create a relationship JSON for a specific workload kind."""
    import copy
    rel = copy.deepcopy(RELATIONSHIP_TEMPLATE)
    
    # Set the model version
    rel["model"]["model"]["version"] = model_version
    
    # Set the workload kind
    rel["selectors"][0]["allow"]["from"][0]["kind"] = workload_kind
    
    return rel

def main():
    """Main function to add relationships to all versions."""
    # Find all version directories
    versions = []
    for item in os.listdir(BASE_PATH):
        version_path = os.path.join(BASE_PATH, item)
        if os.path.isdir(version_path) and item.startswith("v"):
            versions.append(item)
    
    print(f"Found {len(versions)} versions: {sorted(versions)}")
    
    files_created = 0
    
    for version in sorted(versions):
        relationships_path = os.path.join(BASE_PATH, version, "v1.0.0", "relationships")
        
        if not os.path.exists(relationships_path):
            print(f"Warning: {relationships_path} does not exist, skipping...")
            continue
        
        for workload in WORKLOADS:
            # Generate unique filename suffix
            suffix = generate_unique_suffix(workload, version)
            filename = f"edge-non-binding-network-{suffix}.json"
            filepath = os.path.join(relationships_path, filename)
            
            # Check if file already exists
            if os.path.exists(filepath):
                print(f"File already exists: {filepath}, skipping...")
                continue
            
            # Create the relationship
            relationship = create_relationship(workload, version)
            
            # Write the file
            with open(filepath, 'w') as f:
                json.dump(relationship, f, indent=2)
            
            print(f"Created: {filepath}")
            files_created += 1
    
    print(f"\nTotal files created: {files_created}")

if __name__ == "__main__":
    main()
