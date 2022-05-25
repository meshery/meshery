# meshery-cpx

![Version: 0.5.0](https://img.shields.io/badge/Version-0.5.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: stable-latest](https://img.shields.io/badge/AppVersion-stable--latest-informational?style=flat-square)

Meshery Adapter for Citrix chart.

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| aisuko | <urakiny@gmail.com> |  |
| leecalcote | <leecalcote@gmail.com> |  |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` |  |
| annotations | object | `{}` |  |
| env | object | `{}` |  |
| fullnameOverride | string | `""` |  |
| image.pullPolicy | string | `"Always"` |  |
| image.repository | string | `"layer5/meshery-cpx:stable-latest"` |  |
| imagePullSecrets | list | `[]` |  |
| ingress.annotations | string | `nil` |  |
| ingress.enabled | bool | `false` |  |
| ingress.hosts[0].host | string | `"chart-example.local"` |  |
| ingress.hosts[0].paths | list | `[]` |  |
| ingress.tls | list | `[]` |  |
| nodeSelector | object | `{}` |  |
| podSecurityContext | object | `{}` |  |
| probe.livenessProbe.enabled | bool | `false` |  |
| probe.readinessProbe.enabled | bool | `false` |  |
| replicaCount | int | `1` |  |
| resources | object | `{}` |  |
| securityContext | object | `{}` |  |
| service.annotations | object | `{}` |  |
| service.port | int | `10008` |  |
| service.type | string | `"ClusterIP"` |  |
| serviceAccount.create | bool | `false` |  |
| serviceAccount.name | string | `""` |  |
| serviceAccountNameOverride | string | `""` |  |
| testCase.enabled | bool | `false` |  |
| tolerations | list | `[]` |  |

