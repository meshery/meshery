# meshery-operator

![Version: 1.0.0](https://img.shields.io/badge/Version-1.0.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 1.0.0](https://img.shields.io/badge/AppVersion-1.0.0-informational?style=flat-square)

Meshery Operator chart. Deploys the [meshery-operator](https://github.com/meshery/meshery-operator) manager, which reconciles the `Broker` (NATS) and `MeshSync` custom resources.

The chart's `version`/`appVersion` and the CRD bundles under `crds/` and
`files/` are kept in lockstep with meshery-operator releases by that repo's
`sync-downstream` workflow — do not hand-edit them here.

## CRD lifecycle

- **Install**: Helm applies `crds/crds.yaml` before anything else, and the
  `crds.updateJob` pre-install hook server-side-applies the same bundle.
- **Upgrade**: Helm [never touches `crds/` on upgrade](https://helm.sh/docs/chart_best_practices/custom_resource_definitions/);
  the `crds.updateJob` pre-upgrade hook is what refreshes the CRDs, so
  `helm upgrade` delivers schema updates to live clusters.
- **Uninstall**: the CRDs (and every `Broker`/`MeshSync` object) deliberately
  survive `helm uninstall` — they are never owned by the release. To remove
  them permanently, including all custom resources of those types:

  ```console
  kubectl delete crd brokers.meshery.io meshsyncs.meshery.io
  ```

- **Conversion**: both CRDs serve `v1alpha1` and `v1alpha2` (storage) with
  conversion strategy `None`, which is exact while the two schemas are
  field-identical — no webhook and no cert-manager required. When the schemas
  diverge upstream, enable `webhook.enabled` (chart-generated self-signed
  certificate) or `webhook.certManager.enabled` (cert-manager v1) and the CRD
  update Job patches webhook conversion onto the CRDs.

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Meshery Authors | <developers@meshery.io> |  |
| aisuko | <urakiny@gmail.com> |  |
| maintainers | <maintainers@meshery.io> |  |

## Requirements

| Repository | Name | Version |
|------------|------|---------|
|  | meshery-broker | 0.5.0 |
|  | meshery-meshsync | 0.5.0 |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` |  |
| annotations | object | `{}` | Extra annotations for the manager Deployment |
| crds.updateJob.enabled | bool | `true` | Pre-install/pre-upgrade hook Job that server-side-applies `files/crds.yaml`; disabling means upgrades will NOT refresh CRDs |
| crds.updateJob.image.pullPolicy | string | `"IfNotPresent"` |  |
| crds.updateJob.image.repository | string | `"alpine/k8s"` | kubectl-capable image for the CRD update Job |
| crds.updateJob.image.tag | string | `"1.35.6"` |  |
| env | object | `{}` | Extra environment variables for the manager container |
| fullnameOverride | string | `"meshery-operator"` |  |
| image.pullPolicy | string | `"IfNotPresent"` |  |
| image.repository | string | `"meshery/meshery-operator"` |  |
| image.tag | string | `"1.0.0"` | Pinned operator release, stamped by the sync workflow. Kept explicit because server-release chart publishing re-stamps appVersion with the server tag; empty falls back to the chart appVersion |
| imagePullSecrets | list | `[]` |  |
| ingress.annotations | object | `{}` |  |
| ingress.enabled | bool | `false` |  |
| ingress.hosts[0].host | string | `"chart-example.local"` |  |
| ingress.hosts[0].paths | list | `[]` |  |
| ingress.tls | list | `[]` |  |
| leaderElection.enabled | bool | `true` | Leader election (adds the leases Role/RoleBinding and `--enable-leader-election`) |
| meshery-broker.enabled | bool | `true` |  |
| meshery-broker.fullnameOverride | string | `"meshery-broker"` |  |
| meshery-broker.serviceAccountNameOverride | string | `"meshery-server"` |  |
| meshery-meshsync.enabled | bool | `true` |  |
| meshery-meshsync.fullnameOverride | string | `"meshery-meshsync"` |  |
| meshery-meshsync.serviceAccountNameOverride | string | `"meshery-server"` |  |
| nameOverride | string | `""` |  |
| nodeSelector | object | `{}` |  |
| podAnnotations | object | `{}` |  |
| podSecurityContext | object | runAsNonRoot 65532, RuntimeDefault seccomp | Parity with the operator's own config/manager |
| replicaCount | int | `1` |  |
| resources | object | limits 500m/256Mi, requests 100m/64Mi | Parity with the operator's own config/manager |
| securityContext | object | no privilege escalation, read-only rootfs, drop ALL |  |
| service.annotations | object | `{}` |  |
| service.port | int | `8443` | TLS metrics endpoint (authn/authz-filtered; bind scrapers to the `meshery-metrics-reader` ClusterRole) |
| service.type | string | `"ClusterIP"` |  |
| serviceAccount.create | bool | `true` |  |
| serviceAccount.name | string | `"meshery-operator"` |  |
| tolerations | list | `[]` |  |
| webhook.certManager.enabled | bool | `false` | Issue the webhook serving cert with cert-manager (cert-manager.io/v1) + CA injection instead of the chart-generated self-signed certificate |
| webhook.enabled | bool | `false` | Serve v1alpha1<->v1alpha2 CRD conversion through the operator's webhook; not required while the schemas are field-identical |
