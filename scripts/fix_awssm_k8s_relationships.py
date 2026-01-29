#!/usr/bin/env python3
"""
Script to add AWS Secrets Manager to Kubernetes workload relationships with unique UUIDs.

This script adds/fixes relationships between AWS Secrets Manager Secret and the following
Kubernetes workloads:
- Deployment (NEW)
- StatefulSet
- DaemonSet
- Job
- CronJob
- Pod (NEW)

These relationships enable visual secret management in Meshery Kanvas.
"""

import os
import json
import uuid as uuid_module
from typing import Dict, Any
from pathlib import Path

# Get the absolute path relative to this script's location
BASE_PATH = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", "server/meshmodel/aws-secretsmanager-controller"))

def generate_unique_uuid_for_workload(workload_kind: str, model_version: str) -> str:
    """
    Generate a stable, unique UUID for a given workload and model version.
    Uses a deterministic approach based on workload kind and model version.
    """
    namespace = uuid_module.NAMESPACE_DNS
    name = f"{workload_kind.lower()}-{model_version}-aws-secretsmanager.meshery.io"
    return str(uuid_module.uuid5(namespace, name))

def create_relationship(
    workload_kind: str,
    model_version: str,
    is_pod: bool = False
) -> Dict[str, Any]:
    """
    Create a relationship JSON for a specific workload kind.
    
    Args:
        workload_kind: The Kubernetes workload kind (Deployment, StatefulSet, etc.)
        model_version: The AWS Secrets Manager controller version
        is_pod: Whether this is for a Pod (which has different path structure)
    """
    
    # Generate unique UUIDs
    relationship_id = generate_unique_uuid_for_workload(workload_kind, model_version)
    model_id = generate_unique_uuid_for_workload(f"model-{model_version}", model_version)
    k8s_model_id = generate_unique_uuid_for_workload(f"kubernetes-model-{model_version}", model_version)
    secret_model_id = generate_unique_uuid_for_workload(f"secret-model-{model_version}", model_version)
    
    # Pod has spec.volumes (not spec.template.spec.volumes)
    if is_pod:
        volumes_path = ["spec", "volumes"]
        evaluation_query = "$.spec.volumes[*].secret.secretName"
        description = f"AWS Secrets Manager Secret provides secrets to {workload_kind}"
    else:
        volumes_path = ["spec", "template", "spec", "volumes"]
        evaluation_query = "$.spec.template.spec.volumes[*].secret.secretName"
        description = f"AWS Secrets Manager Secret provides secrets to {workload_kind}"
    
    relationship = {
        "id": relationship_id,
        "evaluationQuery": evaluation_query,
        "kind": "edge",
        "metadata": {
            "description": description,
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
            "id": model_id,
            "registrant": {
                "kind": ""
            },
            "model": {
                "version": model_version
            }
        },
        "schemaVersion": "relationships.meshery.io/v1alpha3",
        "selectors": [
            {
                "allow": {
                    "from": [
                        {
                            "id": None,
                            "kind": workload_kind,
                            "match": {},
                            "match_strategy_matrix": None,
                            "model": {
                                "version": "",
                                "name": "kubernetes",
                                "displayName": "",
                                "id": k8s_model_id,
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
                                    ["name"],
                                    ["configuration"] + volumes_path
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
                                "id": secret_model_id,
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
                                    ["name"],
                                    ["configuration", "arn"]
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
    
    return relationship

def generate_filename_suffix(workload_kind: str) -> str:
    """Generate a 5-character suffix for the filename based on workload kind."""
    # Map workloads to consistent suffixes for easy identification
    suffix_map = {
        "Deployment": "deplo",
        "StatefulSet": "stsfs",
        "DaemonSet": "dmnst",
        "Job": "jobbs",
        "CronJob": "crnjb",
        "Pod": "podds"
    }
    return suffix_map.get(workload_kind, workload_kind[:5].lower())

def main():
    """Main function to add/fix relationships to all versions."""
    # Workloads to add relationships for (order matters for consistency)
    WORKLOADS = ["Deployment", "StatefulSet", "DaemonSet", "Job", "CronJob", "Pod"]
    
    # Find all version directories
    versions = []
    for item in os.listdir(BASE_PATH):
        version_path = os.path.join(BASE_PATH, item)
        if os.path.isdir(version_path) and item.startswith("v"):
            versions.append(item)
    
    print(f"Found {len(versions)} versions: {sorted(versions)}")
    
    files_created = 0
    files_updated = 0
    
    for version in sorted(versions):
        relationships_path = os.path.join(BASE_PATH, version, "v1.0.0", "relationships")
        
        if not os.path.exists(relationships_path):
            print(f"Warning: {relationships_path} does not exist, skipping...")
            continue
        
        for workload in WORKLOADS:
            # Generate filename suffix
            suffix = generate_filename_suffix(workload)
            filename = f"edge-non-binding-network-{suffix}.json"
            filepath = os.path.join(relationships_path, filename)
            
            # Create the relationship
            is_pod = (workload == "Pod")
            relationship = create_relationship(workload, version, is_pod=is_pod)
            
            # Write the file
            file_exists = os.path.exists(filepath)
            with open(filepath, 'w') as f:
                json.dump(relationship, f, indent=2)
            
            if file_exists:
                print(f"Updated: {filepath}")
                files_updated += 1
            else:
                print(f"Created: {filepath}")
                files_created += 1
    
    print(f"\n{'='*60}")
    print(f"Total files created: {files_created}")
    print(f"Total files updated: {files_updated}")
    print(f"Total files changed: {files_created + files_updated}")
    print(f"{'='*60}")

if __name__ == "__main__":
    main()
