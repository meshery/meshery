# meshery-meshsync

![Version: 0.5.0](https://img.shields.io/badge/Version-0.5.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 1.0.5](https://img.shields.io/badge/AppVersion-1.0.5-informational?style=flat-square)

Meshery MeshSync

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Meshery Authors | <maintainers@meshery.io> |  |
| darrenlau | <panyuenlau@gmail.com> |  |
| maintainers | <maintainers@meshery.io> |  |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| broker.name | string | `"meshery-broker"` |  |
| broker.namespace | string | `"meshery"` |  |
| name | string | `"meshery-meshsync"` |  |
| replica | int | `1` |  |
| watchConfig.whitelist | string | `"\"[{\\\"Resource\\\":\\\"grafanas.v1beta1.grafana.integreatly.org\\\",\\\"Events\\\":[\\\"ADDED\\\",\\\"MODIFIED\\\",\\\"DELETED\\\"]},{\\\"Resource\\\":\\\"prometheuses.v1.monitoring.coreos.com\\\",\\\"Events\\\":[\\\"ADDED\\\",\\\"MODIFIED\\\",\\\"DELETED\\\"]},{\\\"Resource\\\":\\\"namespaces.v1.\\\",\\\"Events\\\":[\\\"ADDED\\\",\\\"MODIFIED\\\",\\\"DELETED\\\"]},{\\\"Resource\\\":\\\"configmaps.v1.\\\",\\\"Events\\\":[\\\"ADDED\\\",\\\"MODIFIED\\\",\\\"DELETED\\\"]},{\\\"Resource\\\":\\\"nodes.v1.\\\",\\\"Events\\\":[\\\"ADDED\\\",\\\"MODIFIED\\\",\\\"DELETED\\\"]},{\\\"Resource\\\":\\\"secrets.v1.\\\",\\\"Events\\\":[\\\"ADDED\\\",\\\"MODIFIED\\\",\\\"DELETED\\\"]},{\\\"Resource\\\":\\\"persistentvolumes.v1.\\\",\\\"Events\\\":[\\\"ADDED\\\",\\\"MODIFIED\\\",\\\"DELETED\\\"]},{\\\"Resource\\\":\\\"persistentvolumeclaims.v1.\\\",\\\"Events\\\":[\\\"ADDED\\\",\\\"MODIFIED\\\",\\\"DELETED\\\"]},{\\\"Resource\\\":\\\"replicasets.v1.apps\\\",\\\"Events\\\":[\\\"ADDED\\\",\\\"MODIFIED\\\",\\\"DELETED\\\"]},{\\\"Resource\\\":\\\"pods.v1.\\\",\\\"Events\\\":[\\\"ADDED\\\",\\\"MODIFIED\\\",\\\"DELETED\\\"]},{\\\"Resource\\\":\\\"services.v1.\\\",\\\"Events\\\":[\\\"ADDED\\\",\\\"MODIFIED\\\",\\\"DELETED\\\"]},{\\\"Resource\\\":\\\"deployments.v1.apps\\\",\\\"Events\\\":[\\\"ADDED\\\",\\\"MODIFIED\\\",\\\"DELETED\\\"]},{\\\"Resource\\\":\\\"statefulsets.v1.apps\\\",\\\"Events\\\":[\\\"ADDED\\\",\\\"MODIFIED\\\",\\\"DELETED\\\"]},{\\\"Resource\\\":\\\"daemonsets.v1.apps\\\",\\\"Events\\\":[\\\"ADDED\\\",\\\"MODIFIED\\\",\\\"DELETED\\\"]},{\\\"Resource\\\":\\\"ingresses.v1.networking.k8s.io\\\",\\\"Events\\\":[\\\"ADDED\\\",\\\"MODIFIED\\\",\\\"DELETED\\\"]},{\\\"Resource\\\":\\\"endpoints.v1.\\\",\\\"Events\\\":[\\\"ADDED\\\",\\\"MODIFIED\\\",\\\"DELETED\\\"]},{\\\"Resource\\\":\\\"endpointslices.v1.discovery.k8s.io\\\",\\\"Events\\\":[\\\"ADDED\\\",\\\"MODIFIED\\\",\\\"DELETED\\\"]},{\\\"Resource\\\":\\\"cronjobs.v1.batch\\\",\\\"Events\\\":[\\\"ADDED\\\",\\\"MODIFIED\\\",\\\"DELETED\\\"]},{\\\"Resource\\\":\\\"replicationcontrollers.v1.\\\",\\\"Events\\\":[\\\"ADDED\\\",\\\"MODIFIED\\\",\\\"DELETED\\\"]},{\\\"Resource\\\":\\\"storageclasses.v1.storage.k8s.io\\\",\\\"Events\\\":[\\\"ADDED\\\",\\\"MODIFIED\\\",\\\"DELETED\\\"]},{\\\"Resource\\\":\\\"clusterroles.v1.rbac.authorization.k8s.io\\\",\\\"Events\\\":[\\\"ADDED\\\",\\\"MODIFIED\\\",\\\"DELETED\\\"]},{\\\"Resource\\\":\\\"volumeattachments.v1.storage.k8s.io\\\",\\\"Events\\\":[\\\"ADDED\\\",\\\"MODIFIED\\\",\\\"DELETED\\\"]},{\\\"Resource\\\":\\\"apiservices.v1.apiregistration.k8s.io\\\",\\\"Events\\\":[\\\"ADDED\\\",\\\"MODIFIED\\\",\\\"DELETED\\\"]}]\"\n  "` |  |

