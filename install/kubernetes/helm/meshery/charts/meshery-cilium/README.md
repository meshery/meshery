# meshery-cilium

![Version: 0.1.0](https://img.shields.io/badge/Version-0.1.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.1.0](https://img.shields.io/badge/AppVersion-0.1.0-informational?style=flat-square)

Meshery Adapter for Cilium chart.

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| aisuko | urakiny@gmail.com |  |
| leecalcote | leecalcote@gmail.com |  |
| Layer5 Authors | community@layer5.io |  |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` |  |
| annotations | object | `{}` |  |
| containerPort.port | int | `10012` |  |
| env | object | `{}` |  |
| fullnameOverride | string | `""` |  |
| image.pullPolicy | string | `"Always"` |  |
| image.repository | string | `"layer5/meshery-cilium"` |  |
| imagePullSecrets | list | `[]` |  |
| ingress.annotations | string | `nil` |  |
| ingress.enabled | bool | `false` |  |
| ingress.hosts | list | `[{"host":"chart-example.local","paths":[]}]` |  kubernetes.io/tls-acme: "true" |
| ingress.tls | list | `[]` |  |
| nameOverride | string | `""` |  |
| nodeSelector | object | `{}` |  |
| podSecurityContext | object | `{}` |  |
| probe.livenessProbe.enabled | bool | `false` |  |
| probe.readinessProbe.enabled | bool | `false` |  |
| replicaCount | int | `1` |  |
| resources | object | `{}` |  |
| securityContext | object | `{}` |  |
| service.annotations | object | `{}` |  |
| service.port | int | `10012` |  |
| service.type | string | `"ClusterIP"` |  |
| serviceAccount.create | bool | `false` |  |
| serviceAccount.name | string | `""` |  If not set and create is true, a name is generated using the fullname template |
| serviceAccountNameOverride | string | `""` |  |
| testCase.enabled | bool | `false` |  |
| tolerations | list | `[]` |  |

