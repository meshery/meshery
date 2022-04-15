# meshery

![Version: 0.5.0](https://img.shields.io/badge/Version-0.5.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square)

Meshery chart for deploying Meshery and Meshery's adapters.

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Layer5 Authors | <community@layer5.io> |  |
| aisuko | <urakiny@gmail.com> |  |
| leecalcote | <lee.calcote@layer5.io> |  |
| darrenlau | <panyuenlau@gmail.com> |  |

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
| service.target_port | int | `8080` |  |
| service.type | string | `"ClusterIP"` |  |
| serviceAccount.name | string | `"meshery-server"` |  |
| testCase.enabled | bool | `false` |  |
| tolerations | list | `[]` |  |

