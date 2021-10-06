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
|  | meshery-consul | stable-latest |
|  | meshery-istio | stable-latest |
|  | meshery-kuma | stable-latest |
|  | meshery-linkerd | stable-latest |
|  | meshery-nsm | stable-latest |
|  | meshery-osm | stable-latest |
|  | meshery-nginx-sm | stable-latest |

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
| meshery-consul.serviceAccountNameOverride | string | `"meshery-server"` |  |
| meshery-istio.enabled | bool | `true` |  |
| meshery-istio.fullnameOverride | string | `"meshery-istio"` |  |
| meshery-istio.serviceAccountNameOverride | string | `"meshery-server"` |  |
| meshery-kuma.enabled | bool | `true` |  |
| meshery-kuma.fullnameOverride | string | `"meshery-kuma"` |  |
| meshery-kuma.serviceAccountNameOverride | string | `"meshery-server"` |  |
| meshery-linkerd.enabled | bool | `true` |  |
| meshery-linkerd.fullnameOverride | string | `"meshery-linkerd"` |  |
| meshery-linkerd.serviceAccountNameOverride | string | `"meshery-server"` |  |
| meshery-nsm.enabled | bool | `true` |  |
| meshery-nsm.fullnameOverride | string | `"meshery-nsm"` |  |
| meshery-nsm.serviceAccountNameOverride | string | `"meshery-server"` |  |
| meshery-osm.enabled | bool | `true` |  |
| meshery-osm.fullnameOverride | string | `"meshery-osm"` |  |
| meshery-osm.serviceAccountNameOverride | string | `"meshery-server"` |  |
| meshery-nginx-sm.fullnameOverride | string | `"meshery-ngnix-sm"` |  |
| meshery-nginx-sm.serviceAccountNameOverride | string | `"meshery-server"` |  |
| meshery-cpx.fullnameOverride | string | `"meshery-cpx"` |  |
| meshery-cpx.serviceAccountNameOverride | string | `"meshery-server"` |
| meshery-traefik-mesh.fullnameOverride | string | `"meshery-traefik-mesh"` |  |
| meshery-traefik-mesh.serviceAccountNameOverride | string | `"meshery-server"` |  |
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
| service.port | int | `9081` |  |
| service.type | string | `"ClusterIP"` |  |
| serviceAccount.name | string | `"meshery:server"` |  |
| testCase.enabled | bool | `false` |  |
| tolerations | list | `[]` |  |

