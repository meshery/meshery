# meshery

![Version: 2.1.2](https://img.shields.io/badge/Version-2.1.2-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: v0.6.0](https://img.shields.io/badge/AppVersion-v0.6.0-informational?style=flat-square)

Meshery chart for deploying Meshery and Meshery's adapters.

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Layer5 Authors | community@layer5.io |  |
| aisuko | urakiny@gmail.com |  |
| leecalcote | lee.calcote@layer5.io |  |
| darrenlau | panyuenlau@gmail.com |  |

## Requirements

| Repository | Name | Version |
|------------|------|---------|
|  | meshery-app-mesh | 2.1.2 |
|  | meshery-broker | 0.1.0 |
|  | meshery-cilium | 0.1.0 |
|  | meshery-consul | 2.1.2 |
|  | meshery-cpx | 2.1.2 |
|  | meshery-istio | 2.1.2 |
|  | meshery-kuma | 2.1.2 |
|  | meshery-linkerd | 2.1.2 |
|  | meshery-meshsync | 0.1.0 |
|  | meshery-nginx-sm | 2.1.2 |
|  | meshery-nsm | 2.1.2 |
|  | meshery-operator | 2.1.2 |
|  | meshery-osm | 2.1.2 |
|  | meshery-traefik-mesh | 2.1.2 |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` |  |
| annotations | object | `{}` |  |
| env.ADAPTER_URLS | string | `"meshery-istio:10000 meshery-linkerd:10001 meshery-consul:10002 meshery-kuma:10007 meshery-osm:10009 meshery-nginx-sm:10010 meshery-nsm:10004 meshery-app-mesh:10005 meshery-traefik-mesh:10006 meshery-cpx:10008"` |  |
| env.EVENT | string | `"mesheryLocal"` |  |
| env.PROVIDER_BASE_URLS | string | `"https://meshery.layer5.io"` |  |
| env.SAAS_BASE_URL | string | `"https://meshery.layer5.io"` |  |
| fullnameOverride | string | `""` |  |
| image.pullPolicy | string | `"Always"` |  |
| image.repository | string | `"layer5/meshery"` |  |
| image.tag | string | `"stable-latest"` |  |
| imagePullSecrets | list | `[]` |  |
| ingress.annotations | object | `{}` |  |
| ingress.enabled | bool | `false` |  |
| ingress.hosts[0].host | string | `"chart-example.local"` |  |
| ingress.hosts[0].paths | list | `[]` |  |
| ingress.tls | list | `[]` |  |
| meshery-app-mesh | object | `{"enabled":true,"fullnameOverride":"meshery-app-mesh","serviceAccountNameOverride":"meshery-server"}` |  If need to use a different service accounts for the adapters, modify the serviceAccountNameOverride for the corresponding adapters here, or set the serviceAccountNameOverride to empty while setting the serviceAccount information in the sub-chart accordingly |
| meshery-cilium.enabled | bool | `false` |  |
| meshery-cilium.fullnameOverride | string | `"meshery-cilium"` |  |
| meshery-consul.enabled | bool | `true` |  |
| meshery-consul.fullnameOverride | string | `"meshery-consul"` |  |
| meshery-consul.serviceAccountNameOverride | string | `"meshery-server"` |  |
| meshery-cpx.enabled | bool | `true` |  |
| meshery-cpx.fullnameOverride | string | `"meshery-cpx"` |  |
| meshery-cpx.serviceAccountNameOverride | string | `"meshery-server"` |  |
| meshery-istio.enabled | bool | `true` |  |
| meshery-istio.fullnameOverride | string | `"meshery-istio"` |  |
| meshery-istio.serviceAccountNameOverride | string | `"meshery-server"` |  |
| meshery-kuma.enabled | bool | `true` |  |
| meshery-kuma.fullnameOverride | string | `"meshery-kuma"` |  |
| meshery-kuma.serviceAccountNameOverride | string | `"meshery-server"` |  |
| meshery-linkerd.enabled | bool | `true` |  |
| meshery-linkerd.fullnameOverride | string | `"meshery-linkerd"` |  |
| meshery-linkerd.serviceAccountNameOverride | string | `"meshery-server"` |  |
| meshery-nginx-sm.enabled | bool | `true` |  |
| meshery-nginx-sm.fullnameOverride | string | `"meshery-nginx-sm"` |  |
| meshery-nginx-sm.serviceAccountNameOverride | string | `"meshery-server"` |  |
| meshery-nsm.enabled | bool | `true` |  |
| meshery-nsm.fullnameOverride | string | `"meshery-nsm"` |  |
| meshery-nsm.serviceAccountNameOverride | string | `"meshery-server"` |  |
| meshery-operator.enabled | bool | `true` |  |
| meshery-operator.fullnameOverride | string | `"meshery-operator"` |  |
| meshery-osm.enabled | bool | `true` |  |
| meshery-osm.fullnameOverride | string | `"meshery-osm"` |  |
| meshery-osm.serviceAccountNameOverride | string | `"meshery-server"` |  |
| meshery-traefik-mesh.enabled | bool | `true` |  |
| meshery-traefik-mesh.fullnameOverride | string | `"meshery-traefik-mesh"` |  |
| meshery-traefik-mesh.serviceAccountNameOverride | string | `"meshery-server"` |  |
| mesherygateway | object | `{"enabled":false,"selector":{"istio":"ingressgateway"}}` |  There will be an issue when deploying Meshery before Istio and this could make the deploying fail. meshery-gateway |
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
| service.target_port | int | `8080` |  |
| service.type | string | `"LoadBalancer"` |  |
| serviceAccount.name | string | `"meshery-server"` |  If not set and create is true, a name is generated using the fullname template |
| testCase.enabled | bool | `false` |  |
| tolerations | list | `[]` |  |

