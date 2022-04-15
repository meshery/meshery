# meshery-nginx-sm

![Version: 0.5.0](https://img.shields.io/badge/Version-0.5.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: stable-latest](https://img.shields.io/badge/AppVersion-stable--latest-informational?style=flat-square)

Meshery Adapter for NGINX Service Mesh.

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Layer5 Authors | <community@layer5.io> |  |
| darrenlau | <panyuenlau@gmail.com> |  |
| leecalcote | <lee.calcote@layer5.io> |  |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` |  |
| annotations."kompose.cmd" | string | `"kompose convert -f ../docker-compose.yaml"` |  |
| annotations."kompose.version" | string | `"1.17.0 ()"` |  |
| env | object | `{}` |  |
| fullnameOverride | string | `""` |  |
| image.pullPolicy | string | `"Always"` |  |
| image.repository | string | `"layer5/meshery-nginx-sm:stable-latest"` |  |
| imagePullSecrets | list | `[]` |  |
| ingress.annotations | object | `{}` |  |
| ingress.enabled | bool | `false` |  |
| ingress.hosts[0].host | string | `"chart-example.local"` |  |
| ingress.hosts[0].paths | list | `[]` |  |
| ingress.tls | list | `[]` |  |
| nameOverride | string | `""` |  |
| nodeSelector | object | `{}` |  |
| podSecurityContext | object | `{}` |  |
| ports.http | int | `10010` |  |
| probe.livenessProbe.enabled | bool | `false` |  |
| probe.readinessProbe.enabled | bool | `false` |  |
| replicaCount | int | `1` |  |
| resources | object | `{}` |  |
| securityContext | object | `{}` |  |
| service.annotations | object | `{}` |  |
| service.port | int | `10010` |  |
| service.targetPort | int | `10010` |  |
| service.type | string | `"ClusterIP"` |  |
| serviceAccount.create | bool | `false` |  |
| serviceAccount.name | string | `""` |  |
| serviceAccountNameOverride | string | `""` |  |
| testCase.enabled | bool | `false` |  |
| tolerations | list | `[]` |  |

