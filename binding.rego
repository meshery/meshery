---------
package kubernetes_policy

# Verify that each ClusterRoleBinding is bounded to a valid ClusterRole
valid_cluster_role_binding {
    clusterRoleBinding := input.clusterRoleBindings[_]
    clusterRole := input.clusterRoles[clusterRoleBinding.roleRef.name]

    clusterRole != null
    clusterRoleBinding.roleRef.kind == "ClusterRole"
}
------------

------------
package k8s_pvc_policies

import data.kubernetes

# Define the main rule to evaluate the relationship between Pod, PV, and PVC
deny_pod_without_claim[pod] {
    # Get the PVC name from the pod's spec
    pvc_name = pod.spec.volumes[_].persistentVolumeClaim.claimName

    # Check if the PVC exists in the cluster
    not kubernetes.persistent_volume_claims[pvc_name]

    # Ensure that the pod has at least one PVC
    count(pvc_name) == 0
}

# Main rule for evaluating relationships between PV and PVC
deny_pv_without_claim[pv] {
    # Get the PVC name from the PV's claimRef
    pvc_name = pv.spec.claimRef.name

    # Check if the PVC exists in the cluster
    not kubernetes.persistent_volume_claims[pvc_name]
}

# Define the Kubernetes data schema
package kubernetes {
    # Define data structures for Kubernetes resources
    pod_manifests = {
        "example-pod": {
            "apiVersion": "v1",
            "kind": "Pod",
            "metadata": {
                "name": "example-pod",
            },
            "spec": {
                "volumes": [
                    {
                        "name": "data-volume",
                        "persistentVolumeClaim": {
                            "claimName": "example-claim",
                        }
                    }
                ],
                "containers": [
                    {
                        "name": "app-container",
                        "image": "nginx:latest",
                    }
                ],
            },
        },
    }

    persistent_volume_claims = {
        "example-claim": {
            "apiVersion": "v1",
            "kind": "PersistentVolumeClaim",
            "metadata": {
                "name": "example-claim",
            },
            "spec": {
                "accessModes": ["ReadWriteOnce"],
                "resources": {
                    "requests": {
                        "storage": "1Gi",
                    },
                },
            },
        },
    }
}

------------
------------
package kubernetes_policy

# Verify that each service is associated with a valid deployment
valid_service_deployment {
    service := input.services[_]
    deployment := input.deployments[service.metadata.name]

    deployment != null
    deployment.spec.selector.matchLabels == service.spec.selector
}
-----------

----------
package kubernetes_policy

# Helper function to check if two sets are equal
is_set_equal(s1, s2) {
    s1 == s2
}

# Verify that the service selectors match the deployment labels
service_deployment_match {
    service := input.services[_]
    deployment := input.deployments[service.deployment]

    is_set_equal(service.selectors, deployment.labels)
}

# Verify that the deployment selectors match the pod labels
deployment_pod_match {
    deployment := input.deployments[_]
    pod := input.pods[deployment.pod]

    is_set_equal(deployment.selectors, pod.labels)
}

# Verify that each service is associated with a valid deployment
valid_service_deployment {
    service := input.services[_]
    deployment := input.deployments[service.deployment]

    deployment != null
}

# Verify that each deployment is associated with a valid pod
valid_deployment_pod {
    deployment := input.deployments[_]
    pod := input.pods[deployment.pod]

    pod != null
}

# Verify that each pod is associated with a valid service
valid_pod_service {
    pod := input.pods[_]
    deployment := input.deployments[_]
    service := input.services[deployment.service]

    service != null
}
--------------