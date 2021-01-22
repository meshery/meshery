# meshery

![Version: 2.1.2](https://img.shields.io/badge/Version-2.1.2-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: stable-latest](https://img.shields.io/badge/AppVersion-stable--latest-informational?style=flat-square)

Meshery chart for deploying Meshery and Meshery's adapters.

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Layer5 Authors | community@layer5.io |  |
| aisuko | urakiny@gmail.com |  |
| leecalcote | lee.calcote@layer5.io |  |

## Requirements

| Repository | Name | Version |
|------------|------|---------|
|  | meshery-consul | 2.0.1 |
|  | meshery-istio | 2.0.1 |
|  | meshery-kuma | 1.0.1 |
|  | meshery-linkerd | 2.0.1 |
|  | meshery-nsm | 1.0.1 |
|  | meshery-osm | 1.0.1 |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` |  |
| annotations | object | `{}` |  |
| env.ADAPTER_URLS | string | `"meshery-istio:10000 meshery-linkerd:10001 meshery-consul:10002 meshery-kuma:10007 meshery-osm:10009 meshery-nsm:10004"` |  |
| env.EVENT | string | `"mesheryLocal"` |  |
| env.PROVIDER_BASE_URLS | string | `"https://meshery.layer5.io"` |  |
| env.SAAS_BASE_URL | string | `"https://meshery.layer5.io"` |  |
| fullnameOverride | string | `""` |  |
| image.pullPolicy | string | `"IfNotPresent"` |  |
| image.repository | string | `"layer5/meshery:stable-latest"` |  |
| imagePullSecrets | list | `[]` |  |
| ingress.annotations | object | `{}` |  |
| ingress.enabled | bool | `false` |  |
| ingress.hosts[0].host | string | `"chart-example.local"` |  |
| ingress.hosts[0].paths | list | `[]` |  |
| ingress.tls | list | `[]` |  |
| meshery-consul.enabled | bool | `true` |  |
| meshery-consul.fullnameOverride | string | `"meshery-consul"` |  |
| meshery-istio.enabled | bool | `true` |  |
| meshery-istio.fullnameOverride | string | `"meshery-istio"` |  |
| meshery-kuma.enabled | bool | `true` |  |
| meshery-kuma.fullnameOverride | string | `"meshery-kuma"` |  |
| meshery-linkerd.enabled | bool | `true` |  |
| meshery-linkerd.fullnameOverride | string | `"meshery-linkerd"` |  |
| meshery-nsm.enabled | bool | `true` |  |
| meshery-nsm.fullnameOverride | string | `"meshery-nsm"` |  |
| meshery-osm.enabled | bool | `true` |  |
| meshery-osm.fullnameOverride | string | `"meshery-osm"` |  |
| mesherygateway.enabled | bool | `false` |  |
| mesherygateway.selector.istio | string | `"ingressgateway"` |  |
| nameOverride | string | `""` |  |
| nodeSelector | object | `{}` |  |
| podSecurityContext | object | `{}` |  |
| probe.livenessProbe.enabled | bool | `false` |  |
| probe.readinessProbe.enabled | bool | `false` |  |
| rbac.nodes | bool | `false` |  |
| replicaCount | int | `1` |  |
| resources | object | `{}` |  |
| restartPolicy | string | `"Always"` |  |
| securityContext | object | `{}` |  |
| service.annotations | object | `{}` |  |
| service.port | int | `8080` |  |
| service.type | string | `"ClusterIP"` |  |
| serviceAccount.name | string | `"meshery-service-account"` |  |
| testCase.enabled | bool | `false` |  |
| tolerations | list | `[]` |  |

