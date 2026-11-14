# meshery-broker

![Version: 0.5.0](https://img.shields.io/badge/Version-0.5.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 1.0.6](https://img.shields.io/badge/AppVersion-1.0.6-informational?style=flat-square)

Meshery Broker

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Meshery Authors | <maintainers@meshery.io> |  |
| darrenlau | <panyuenlau@gmail.com> |  |
| maintainers | <maintainers@meshery.io> |  |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| name | string | `"meshery-broker"` | Name of the Broker custom resource. Meshery Server addresses it as `meshery-broker`. |
| replica | int | `1` | NATS server replicas (Broker `spec.size`, 1-10). |
| service.annotations | object | `{}` | Annotations merged onto the client Service (cloud load-balancer hints, MetalLB pools). |
| service.externalEndpointOverride | string | `""` | Advertised external endpoint (`host:port`) when it cannot be derived from the Service (ingress, gateway, NAT). |
| service.loadBalancerClass | string | `""` | Load-balancer implementation to use. Requires `type: LoadBalancer`. |
| service.loadBalancerSourceRanges | list | `[]` | CIDRs allowed to reach the Broker. Requires `type: LoadBalancer`. |
| service.type | string | `""` | `ClusterIP`, `NodePort` or `LoadBalancer`. Empty: the operator default, `ClusterIP`. |
| version | string | `""` | NATS server image tag (Broker `spec.version`). Empty: the operator's pinned NATS release. |
