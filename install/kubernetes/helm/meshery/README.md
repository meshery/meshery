# meshery

![Version: 0.6.0](https://img.shields.io/badge/Version-0.6.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square)

Meshery chart for deploying Meshery

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Meshery Authors | <maintainers@meshery.io> |  |

## Values

| Key                                             | Type | Default                                                                                                                                                                                               | Description |
|-------------------------------------------------|------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-------------|
| affinity                                        | object | `{}`                                                                                                                                                                                                  |  |
| annotations                                     | object | `{}`                                                                                                                                                                                                  |  |
| env.ADAPTER_URLS                                | string | `"meshery-istio:10000 meshery-linkerd:10001 meshery-consul:10002 meshery-kuma:10007 meshery-nginx-sm:10010 meshery-nsm:10004 meshery-app-mesh:10005 meshery-traefik-mesh:10006 meshery-cilium:10012"` | Optionally, pre-configure Meshery Server with the set of Meshery Adapters used in the deployment. |
| env.EVENT                                       | string | `"mesheryLocal"`                                                                                                                                                                                      |  |
| env.PROVIDER                                    | string | `"Local"`                                                                                                                                                                                             | Use this security-related setting to enforce selection of one and only one Provider. In this way, your Meshery deployment will only trust and only allow users to authenticate using the Provider you have configured in this setting. See the [Remote Provider documentation](https://docs.meshery.io/extensibility/providers) for a description of what a Provider is.  |
| env.MESHERY_SERVER_CALLBACK_URL                 | string | `""`                                                                                                                                                                                                  | Configure an OAuth callback URL for Meshery Server to use when signing into a Remote Provider and your Meshery Server instance is not directly reachable by that Remote Provider. See the [Remote Provider documentation](https://docs.meshery.io/extensibility/providers#configurable-oauth-callback-url) for more details. |
| env.PROVIDER_BASE_URLS                          | string | `"https://meshery.layer5.io"`                                                                                                                                                                         | Configure your Remote Provider of choice. See the [Remote Provider documentation](https://docs.meshery.io/extensibility/providers) for a description of what a Provider is. |
| fullnameOverride                                | string | `""`                                                                                                                                                                                                  |  |
| image.pullPolicy                                | string | `"Always"`                                                                                                                                                                                            |  |
| image.repository                                | string | `"layer5/meshery"`                                                                                                                                                                                    |  |
| image.tag                                       | string | `"stable-latest"`                                                                                                                                                                                     |  |
| imagePullSecrets                                | list | `[]`                                                                                                                                                                                                  |  |
| ingress.annotations                             | object | `{}`                                                                                                                                                                                                  |  |
| ingress.enabled                                 | bool | `false`                                                                                                                                                                                               |  |
| ingress.hosts[0].host                           | string | `"chart-example.local"`                                                                                                                                                                               |  |
| ingress.hosts[0].paths                          | list | `[]`                                                                                                                                                                                                  |  |
| ingress.tls                                     | list | `[]`                                                                                                                                                                                                  |  |
| meshery-app-mesh.enabled                        | bool | `false`                                                                                                                                                                                               | Enable to deploy this Meshery Adapter upon initial deployment. Meshery Adapters can be deployed post-installation using either Meshery CLI or UI. |
| meshery-app-mesh.fullnameOverride               | string | `"meshery-app-mesh"`                                                                                                                                                                                  |  |
| meshery-app-mesh.serviceAccountNameOverride     | string | `"meshery-server"`                                                                                                                                                                                    |  |
| meshery-cilium.enabled                          | bool | `false`                                                                                                                                                                                               | Enable to deploy this Meshery Adapter upon initial deployment. Meshery Adapters can be deployed post-installation using either Meshery CLI or UI. |
| meshery-cilium.fullnameOverride                 | string | `"meshery-cilium"`                                                                                                                                                                                    |  |
| meshery-consul.enabled                          | bool | `false`                                                                                                                                                                                               | Enable to deploy this Meshery Adapter upon initial deployment. Meshery Adapters can be deployed post-installation using either Meshery CLI or UI. |
| meshery-consul.fullnameOverride                 | string | `"meshery-consul"`                                                                                                                                                                                    |  |
| meshery-consul.serviceAccountNameOverride       | string | `"meshery-server"`                                                                                                                                                                                    |  |
| meshery-istio.enabled                           | bool | `false`                                                                                                                                                                                               | Enable to deploy this Meshery Adapter upon initial deployment. Meshery Adapters can be deployed post-installation using either Meshery CLI or UI. |
| meshery-istio.fullnameOverride                  | string | `"meshery-istio"`                                                                                                                                                                                     |  |
| meshery-istio.serviceAccountNameOverride        | string | `"meshery-server"`                                                                                                                                                                                    |  |
| meshery-kuma.enabled                            | bool | `false`                                                                                                                                                                                               | Enable to deploy this Meshery Adapter upon initial deployment. Meshery Adapters can be deployed post-installation using either Meshery CLI or UI. |
| meshery-kuma.fullnameOverride                   | string | `"meshery-kuma"`                                                                                                                                                                                      |  |
| meshery-kuma.serviceAccountNameOverride         | string | `"meshery-server"`                                                                                                                                                                                    |  |
| meshery-linkerd.enabled                         | bool | `false`                                                                                                                                                                                               | Enable to deploy this Meshery Adapter upon initial deployment. Meshery Adapters can be deployed post-installation using either Meshery CLI or UI. |
| meshery-linkerd.fullnameOverride                | string | `"meshery-linkerd"`                                                                                                                                                                                   |  |
| meshery-linkerd.serviceAccountNameOverride      | string | `"meshery-server"`                                                                                                                                                                                    |  |
| meshery-nginx-sm.enabled                        | bool | `false`                                                                                                                                                                                               | Enable to deploy this Meshery Adapter upon initial deployment. Meshery Adapters can be deployed post-installation using either Meshery CLI or UI. |
| meshery-nginx-sm.fullnameOverride               | string | `"meshery-nginx-sm"`                                                                                                                                                                                  |  |
| meshery-nginx-sm.serviceAccountNameOverride     | string | `"meshery-server"`                                                                                                                                                                                    |  |
| meshery-nsm.enabled                             | bool | `false`                                                                                                                                                                                               | Enable to deploy this Meshery Adapter upon initial deployment. Meshery Adapters can be deployed post-installation using either Meshery CLI or UI. |
| meshery-nsm.fullnameOverride                    | string | `"meshery-nsm"`                                                                                                                                                                                       |  |
| meshery-nsm.serviceAccountNameOverride          | string | `"meshery-server"`                                                                                                                                                                                    |  |
| meshery-nighthawk.enabled                       | bool | `false`                                                                                                                                                                                               | Enable to deploy this Meshery Adapter upon initial deployment. Meshery Adapters can be deployed post-installation using either Meshery CLI or UI. |
| meshery-nighthawk.fullnameOverride              | string | `"meshery-nighthawk"`                                                                                                                                                                                 |  |
| meshery-nighthawk.serviceAccountNameOverride    | string | `"meshery-nighthawk"`                                                                                                                                                                                    |  |
| meshery-operator.enabled                        | bool | `true`                                                                                                                                                                                                | Enable to deploy this Meshery Operator upon initial deploymeent. Meshery Operator can be deployed post-installation using Meshery UI. |
| meshery-operator.fullnameOverride               | string | `"meshery-operator"`                                                                                                                                                                                  |  |
| meshery-osm.enabled                             | bool | `false`                                                                                                                                                                                               | OSM is an archived project. |
| meshery-osm.fullnameOverride                    | string | `"meshery-osm"`                                                                                                                                                                                       |  |
| meshery-osm.serviceAccountNameOverride          | string | `"meshery-server"`                                                                                                                                                                                    |  |
| meshery-traefik-mesh.enabled                    | bool | `false`                                                                                                                                                                                               | Enable to deploy this Meshery Adapter upon initial deployment. Meshery Adapters can be deployed post-installation using either Meshery CLI or UI. |
| meshery-traefik-mesh.fullnameOverride           | string | `"meshery-traefik-mesh"`                                                                                                                                                                              |  |
| meshery-traefik-mesh.serviceAccountNameOverride | string | `"meshery-server"`                                                                                                                                                                                    |  |
| mesherygateway.enabled                          | bool | `false`                                                                                                                                                                                               |  |
| mesherygateway.selector.istio                   | string | `"ingressgateway"`                                                                                                                                                                                    |  |
| metadata.name                                   | string | `"meshery"`                                                                                                                                                                                           |  |
| metadata.namespace                              | string | `"meshery"`                                                                                                                                                                                           |  |
| nameOverride                                    | string | `""`                                                                                                                                                                                                  |  |
| nodeSelector                                    | object | `{}`                                                                                                                                                                                                  |  |
| podSecurityContext                              | object | `{}`                                                                                                                                                                                                  |  |
| probe.livenessProbe.enabled                     | bool | `false`                                                                                                                                                                                               |  |
| probe.readinessProbe.enabled                    | bool | `false`                                                                                                                                                                                               |  |
| rbac.nodes                                      | bool | `false`                                                                                                                                                                                               |  |
| replicaCount                                    | int | `1`                                                                                                                                                                                                   |  |
| resources                                       | object | `{}`                                                                                                                                                                                                  |  |
| restartPolicy                                   | string | `"Always"`                                                                                                                                                                                            |  |
| securityContext                                 | object | `{}`                                                                                                                                                                                                  |  |
| service.annotations                             | object | `{}`                                                                                                                                                                                                  |  |
| service.port                                    | int | `9081`                                                                                                                                                                                                |  |
| service.target_port                             | int | `8080`                                                                                                                                                                                                |  |
| service.type                                    | string | `"LoadBalancer"`                                                                                                                                                                                      |  |
| serviceAccount.name                             | string | `"meshery-server"`                                                                                                                                                                                    |  |
| testCase.enabled                                | bool | `false`                                                                                                                                                                                               |  |
| tolerations                                     | list | `[]`                                                                                                                                                                                                  |  |

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
