# meshery

![Version: 0.6.0](https://img.shields.io/badge/Version-0.6.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square)

Meshery chart for deploying Meshery

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Meshery Authors | <maintainers@meshery.io> |  |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` |  |
| annotations | object | `{}` |  |
| env.ADAPTER_URLS | string | `"meshery-istio:10000 meshery-linkerd:10001 meshery-consul:10002 meshery-kuma:10007 meshery-osm:10009 meshery-nginx-sm:10010 meshery-nsm:10004 meshery-app-mesh:10005 meshery-traefik-mesh:10006 meshery-cilium:10012"` |  |
| env.EVENT | string | `"mesheryLocal"` |  |
| env.PROVIDER | string | `""` |  |
| env.MESHERY_SERVER_CALLBACK_URL | string | `""` | The OAuth callback URL to use when Meshery Server is not directly exposed. Refer [docs](https://docs.meshery.io/extensibility/providers#configurable-oauth-callback-url) for more details. |
| env.PROVIDER_BASE_URLS | string | `"https://meshery.layer5.io"` |  |
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
| meshery-app-mesh.enabled | bool | `false` |  |
| meshery-app-mesh.fullnameOverride | string | `"meshery-app-mesh"` |  |
| meshery-app-mesh.serviceAccountNameOverride | string | `"meshery-server"` |  |
| meshery-cilium.enabled | bool | `false` |  |
| meshery-cilium.fullnameOverride | string | `"meshery-cilium"` |  |
| meshery-consul.enabled | bool | `false` |  |
| meshery-consul.fullnameOverride | string | `"meshery-consul"` |  |
| meshery-consul.serviceAccountNameOverride | string | `"meshery-server"` |  |
| meshery-istio.enabled | bool | `false` |  |
| meshery-istio.fullnameOverride | string | `"meshery-istio"` |  |
| meshery-istio.serviceAccountNameOverride | string | `"meshery-server"` |  |
| meshery-kuma.enabled | bool | `false` |  |
| meshery-kuma.fullnameOverride | string | `"meshery-kuma"` |  |
| meshery-kuma.serviceAccountNameOverride | string | `"meshery-server"` |  |
| meshery-linkerd.enabled | bool | `false` |  |
| meshery-linkerd.fullnameOverride | string | `"meshery-linkerd"` |  |
| meshery-linkerd.serviceAccountNameOverride | string | `"meshery-server"` |  |
| meshery-nginx-sm.enabled | bool | `false` |  |
| meshery-nginx-sm.fullnameOverride | string | `"meshery-nginx-sm"` |  |
| meshery-nginx-sm.serviceAccountNameOverride | string | `"meshery-server"` |  |
| meshery-nsm.enabled | bool | `false` |  |
| meshery-nsm.fullnameOverride | string | `"meshery-nsm"` |  |
| meshery-nsm.serviceAccountNameOverride | string | `"meshery-server"` |  |
| meshery-operator.enabled | bool | `true` |  |
| meshery-operator.fullnameOverride | string | `"meshery-operator"` |  |
| meshery-osm.enabled | bool | `false` |  |
| meshery-osm.fullnameOverride | string | `"meshery-osm"` |  |
| meshery-osm.serviceAccountNameOverride | string | `"meshery-server"` |  |
| meshery-perf.enabled | bool | `false` |  |
| meshery-perf.fullnameOverride | string | `"meshery-perf"` |  |
| meshery-perf.serviceAccountNameOverride | string | `"meshery-server"` |  |
| meshery-traefik-mesh.enabled | bool | `false` |  |
| meshery-traefik-mesh.fullnameOverride | string | `"meshery-traefik-mesh"` |  |
| meshery-traefik-mesh.serviceAccountNameOverride | string | `"meshery-server"` |  |
| mesherygateway.enabled | bool | `false` |  |
| mesherygateway.selector.istio | string | `"ingressgateway"` |  |
| metadata.name | string | `"meshery"` |  |
| metadata.namespace | string | `"meshery"` |  |
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
| serviceAccount.name | string | `"meshery-server"` |  |
| testCase.enabled | bool | `false` |  |
| tolerations | list | `[]` |  |

## Setup Repo Info

```console
helm repo add meshery meshery https://meshery.io/charts/
helm repo update
```

_See [helm repo](https://helm.sh/docs/helm/helm_repo/) for command documentation._

## Installing the Chart

To install the chart with the release name `meshery`:

```console
kubectl create namespace meshery
helm install meshery meshery/meshery
```

## Uninstalling the Chart

To uninstall/delete the `meshery` deployment:

```console
helm delete meshery
```

## Installing the Chart with a custom namespace

```console
kubectl create namespace meshery
helm install meshery meshery/meshery --namespace meshery
```

## Installing the Chart with a custom Meshery Adapters

Eg: For [Meshery Adapter for Istio](https://github.com/meshery/meshery-istio)
```console
kubectl create namespace meshery
helm install meshery meshery/meshery --set meshery-istio.enabled=true
```
