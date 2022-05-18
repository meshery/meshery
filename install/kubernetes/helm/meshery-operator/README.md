# meshery-operator

![Version: 0.6.0](https://img.shields.io/badge/Version-0.6.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square)

Meshery Operator chart.

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Layer5 Authors | <community@layer5.io> |  |
| aisuko | <urakiny@gmail.com> |  |
| leecalcote | <leecalcote@gmail.com> |  |

## Requirements

| Repository | Name | Version |
|------------|------|---------|
|  | meshery | 0.5.0 |
|  | meshery-app-mesh | 0.5.0 |
|  | meshery-broker | 0.5.0 |
|  | meshery-cilium | 0.5.0 |
|  | meshery-consul | 0.5.0 |
|  | meshery-cpx | 0.5.0 |
|  | meshery-istio | 0.5.0 |
|  | meshery-kuma | 0.5.0 |
|  | meshery-linkerd | 0.5.0 |
|  | meshery-meshsync | 0.5.0 |
|  | meshery-nginx-sm | 0.5.0 |
|  | meshery-nsm | 0.5.0 |
|  | meshery-osm | 0.5.0 |
|  | meshery-perf | 0.5.0 |
|  | meshery-traefik-mesh | 0.5.0 |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` |  |
| annotations | object | `{}` |  |
| env | object | `{}` |  |
| fullnameOverride | string | `"meshery-operator"` |  |
| imagePullSecrets | list | `[]` |  |
| ingress.annotations | object | `{}` |  |
| ingress.enabled | bool | `false` |  |
| ingress.hosts[0].host | string | `"chart-example.local"` |  |
| ingress.hosts[0].paths | list | `[]` |  |
| ingress.tls | list | `[]` |  |
| kubeRbac.args[0] | string | `"--secure-listen-address=0.0.0.0:8443"` |  |
| kubeRbac.args[1] | string | `"--upstream=http://127.0.0.1:8080/"` |  |
| kubeRbac.args[2] | string | `"--logtostderr=false"` |  |
| kubeRbac.args[3] | string | `"--v=10"` |  |
| kubeRbac.image.pullPolicy | string | `"IfNotPresent"` |  |
| kubeRbac.image.repository | string | `"gcr.io/kubebuilder/kube-rbac-proxy:v0.5.0"` |  |
| kubeRbac.name | string | `"kube-rbac-proxy"` |  |
| meshery-app-mesh.enabled | bool | `false` |  |
| meshery-app-mesh.fullnameOverride | string | `"meshery-app-mesh"` |  |
| meshery-app-mesh.serviceAccountNameOverride | string | `"meshery-server"` |  |
| meshery-broker.enabled | bool | `true` |  |
| meshery-broker.fullnameOverride | string | `"meshery-broker"` |  |
| meshery-broker.serviceAccountNameOverride | string | `"meshery-server"` |  |
| meshery-cilium.enabled | bool | `false` |  |
| meshery-cilium.fullnameOverride | string | `"meshery-cilium"` |  |
| meshery-consul.enabled | bool | `true` |  |
| meshery-consul.fullnameOverride | string | `"meshery-consul"` |  |
| meshery-consul.serviceAccountNameOverride | string | `"meshery-server"` |  |
| meshery-cpx.enabled | bool | `false` |  |
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
| meshery-meshsync.enabled | bool | `true` |  |
| meshery-meshsync.fullnameOverride | string | `"meshery-meshsync"` |  |
| meshery-meshsync.serviceAccountNameOverride | string | `"meshery-server"` |  |
| meshery-nginx-sm.enabled | bool | `false` |  |
| meshery-nginx-sm.fullnameOverride | string | `"meshery-nginx-sm"` |  |
| meshery-nginx-sm.serviceAccountNameOverride | string | `"meshery-server"` |  |
| meshery-nsm.enabled | bool | `false` |  |
| meshery-nsm.fullnameOverride | string | `"meshery-nsm"` |  |
| meshery-nsm.serviceAccountNameOverride | string | `"meshery-server"` |  |
| meshery-osm.enabled | bool | `true` |  |
| meshery-osm.fullnameOverride | string | `"meshery-osm"` |  |
| meshery-osm.serviceAccountNameOverride | string | `"meshery-server"` |  |
| meshery-perf.enabled | bool | `false` |  |
| meshery-perf.fullnameOverride | string | `"meshery-perf"` |  |
| meshery-perf.serviceAccountNameOverride | string | `"meshery-server"` |  |
| meshery-traefik-mesh.enabled | bool | `false` |  |
| meshery-traefik-mesh.fullnameOverride | string | `"meshery-traefik-mesh"` |  |
| meshery-traefik-mesh.serviceAccountNameOverride | string | `"meshery-server"` |  |
| meshery.enabled | bool | `true` |  |
| meshery.fullnameOverride | string | `"meshery-server"` |  |
| meshery.serviceAccountNameOverride | string | `"meshery-server"` |  |
| mesheryOperator.args[0] | string | `"--metrics-addr=127.0.0.1:8080"` |  |
| mesheryOperator.args[1] | string | `"--enable-leader-election"` |  |
| mesheryOperator.command[0] | string | `"/manager"` |  |
| mesheryOperator.image.pullPolicy | string | `"IfNotPresent"` |  |
| mesheryOperator.image.repository | string | `"layer5/meshery-operator:stable-latest"` |  |
| mesheryOperator.name | string | `"manager"` |  |
| nameOverride | string | `""` |  |
| nodeSelector | object | `{}` |  |
| podSecurityContext | object | `{}` |  |
| probe.livenessProbe.enabled | bool | `false` |  |
| probe.readinessProbe.enabled | bool | `false` |  |
| replicaCount | int | `1` |  |
| resources | object | `{}` |  |
| securityContext | object | `{}` |  |
| service.annotations | object | `{}` |  |
| service.port | int | `10000` |  |
| service.type | string | `"ClusterIP"` |  |
| serviceAccount.create | string | `"create"` |  |
| serviceAccount.name | string | `"meshery-operator"` |  |
| testCase.enabled | bool | `false` |  |
| tolerations | list | `[]` |  |

