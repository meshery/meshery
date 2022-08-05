# meshery-operator

![Version: 0.5.0](https://img.shields.io/badge/Version-0.5.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: stable-latest](https://img.shields.io/badge/AppVersion-stable--latest-informational?style=flat-square)

Meshery Operator

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
| imagePullSecrets | list | `[]` |  |
| ingress.annotations | string | `nil` |  |
| ingress.enabled | bool | `false` |  |
| ingress.hosts[0].host | string | `"chart-example.local"` |  |
| ingress.hosts[0].paths | list | `[]` |  |
| ingress.tls | list | `[]` |  |
| kubeRbac.args[0] | string | `"--secure-listen-address=0.0.0.0:8443"` |  |
| kubeRbac.args[1] | string | `"--upstream=http://127.0.0.1:8080/"` |  |
| kubeRbac.args[2] | string | `"--logtostderr=true"` |  |
| kubeRbac.args[3] | string | `"--v=10"` |  |
| kubeRbac.image.pullPolicy | string | `"Always"` |  |
| kubeRbac.image.repository | string | `"gcr.io/kubebuilder/kube-rbac-proxy:v0.5.0"` |  |
| kubeRbac.name | string | `"kube-rbac-proxy"` |  |
| mesheryOperator.args[0] | string | `"--metrics-addr=127.0.0.1:8080"` |  |
| mesheryOperator.args[1] | string | `"--enable-leader-election"` |  |
| mesheryOperator.command[0] | string | `"/manager"` |  |
| mesheryOperator.image.pullPolicy | string | `"Always"` |  |
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
| serviceAccount.create | bool | `true` |  |
| serviceAccount.name | string | `"meshery-operator"` |  |
| testCase.enabled | bool | `false` |  |
| tolerations | list | `[]` |  |

