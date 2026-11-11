<!--
  HAND-MAINTAINED - do NOT regenerate with `make helm-operator-docs` (helm-docs).
  values.yaml carries no `# --` description comments and there is no
  README.md.gotmpl here, so the generator would delete the extended chart
  description, the "CRD lifecycle" section, and every Description cell below.
  Regenerating is safe only after that migration (curated prose -> README.md.gotmpl,
  descriptions -> `# --` comments in values.yaml); until then edit by hand and keep
  the badges, Requirements and Values defaults in step with Chart.yaml/values.yaml.
  The two subchart READMEs under charts/ are pure helm-docs output and are safe.
-->

# meshery-operator

![Version: 1.0.6](https://img.shields.io/badge/Version-1.0.6-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 1.0.6](https://img.shields.io/badge/AppVersion-1.0.6-informational?style=flat-square)

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
- **Uninstall**: the CRDs deliberately survive `helm uninstall` - they are
  never owned by the release. The `meshery-broker` and `meshery-meshsync`
  custom resources the subcharts render *are* release-owned, and each carries
  a finalizer only the operator removes, so the `customResourceCleanup`
  pre-delete hook deletes them while the operator is still running (see
  [Uninstall](#uninstall)). To remove the CRDs permanently, including all
  custom resources of those types:

  ```console
  kubectl delete crd brokers.meshery.io meshsyncs.meshery.io
  ```

- **Conversion**: both CRDs serve `v1alpha1` and `v1alpha2` (storage) with
  conversion strategy `None`, which is exact while the two schemas are
  field-identical — no webhook and no cert-manager required. When the schemas
  diverge upstream, enable `webhook.enabled` (chart-generated self-signed
  certificate) or `webhook.certManager.enabled` (cert-manager v1) and the CRD
  update Job patches webhook conversion onto the CRDs.

## MeshSync identity and RBAC

The operator's own ServiceAccount (`meshery-operator`, bound to
`meshery-operator-role`) is least-privilege: it can manage the Broker and
MeshSync workloads but cannot read the cluster. MeshSync, which the operator
deploys from the `MeshSync` custom resource, needs the opposite: read access
to everything it discovers, and no writes. The chart therefore creates a
separate identity for it:

| Object | Name |
| --- | --- |
| ServiceAccount | `meshery-meshsync` (`meshsync.serviceAccount.name`) |
| ClusterRole | `meshery-meshsync-role` |
| ClusterRoleBinding | `meshery-meshsync-rolebinding` |

and passes the ServiceAccount name to the manager as
`MESHSYNC_SERVICE_ACCOUNT`, which the operator uses for the MeshSync
Deployment. `meshery-meshsync-role` grants `list`/`watch` on every resource -
MeshSync's watch-list is user-configurable and includes CRDs and
cluster-scoped kinds - plus `get` on its own `MeshSync` resource, Namespaces
and `pods/log`. It does not wildcard `get`, because `get` on `*` also
authorizes the exec, attach, port-forward and proxy subresources. Add rules
(for example `create` on `pods/exec` for Meshery's interactive terminal) with
`meshsync.rbac.extraRules`.

**Compatibility with operator 1.0.6 and earlier.** Those releases ignore
`MESHSYNC_SERVICE_ACCOUNT` and run MeshSync as `meshery-operator`. While
`meshsync.rbac.bindOperatorServiceAccount` is `true` (the default), the chart
also binds the read-only `meshery-meshsync-role` to `meshery-operator`
(`meshery-meshsync-compat-rolebinding`), so MeshSync can read the cluster under
either operator. The operator ServiceAccount gains no write access. The
default flips to `false` in the change that moves `image.tag` to the first
operator release that honors `MESHSYNC_SERVICE_ACCOUNT`; set it to `false`
earlier only if you run such an operator.

## Broker exposure

The `meshery-broker` subchart renders the `Broker` custom resource and can set
how the Broker is exposed. An out-of-cluster Meshery Server must reach the
Broker on `4222/tcp`:

```console
helm upgrade --install meshery-operator meshery/meshery-operator -n meshery \
  --set meshery-broker.service.type=LoadBalancer \
  --set 'meshery-broker.service.loadBalancerSourceRanges={10.0.0.0/8}'
```

Unset `meshery-broker.service.*` and `meshery-broker.version` fields are
omitted from the resource, keeping the operator defaults (`ClusterIP`, the
pinned NATS release) and leaving them free for Meshery Server, which
server-side-applies `spec.service` and `spec.version` from its controllers
configuration under the field manager `meshery-server`. Set each field in one
place only: when both set it, Helm 4 (server-side apply) fails the upgrade with
a field-ownership conflict, and Helm 3 lets whichever wrote last win.
`loadBalancerClass` and `loadBalancerSourceRanges` require
`type: LoadBalancer`; the chart refuses to render them otherwise. The resource
is rendered as `meshery.io/v1alpha1`: both served versions are field-identical,
and the storage version `v1alpha2` would break `helm upgrade` from charts whose
CRDs predate it, because Helm maps a manifest's apiVersion before the CRD
update Job runs. See [charts/meshery-broker](charts/meshery-broker/README.md)
for every field.

## Uninstall

`helm uninstall` would delete the operator and the Broker and MeshSync custom
resources in the same pass, leaving the resources stuck `Terminating` on
finalizers no controller is left to remove. The `customResourceCleanup`
pre-delete hook Job runs first, while the operator is still up: it deletes the
custom resources this release owns (identified by the release's Helm ownership
annotation, so a `helm template | kubectl apply` render of the Job deletes
nothing) and waits up to 90 seconds for the operator to process the
finalizers. If the operator is not reconciling, it removes the finalizers
instead; everything the operator created for those resources carries a
controller ownerReference and is garbage-collected. The Job uses
`crds.updateJob.image` and a Role limited to `brokers` and `meshsyncs` in the
release namespace, and fails only on API errors, in which case it is kept for
inspection and `helm uninstall` stops. Set `customResourceCleanup.enabled=false`
to opt out.

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
| customResourceCleanup.enabled | bool | `true` | Pre-delete hook Job that deletes this release's Broker and MeshSync resources while the operator still runs, so their finalizers are processed (see [Uninstall](#uninstall)) |
| env | object | `{}` | Extra environment variables for the manager container. `ENABLE_WEBHOOKS` and `MESHSYNC_SERVICE_ACCOUNT` are chart-managed and rejected here |
| fullnameOverride | string | `"meshery-operator"` |  |
| image.pullPolicy | string | `"IfNotPresent"` |  |
| image.repository | string | `"meshery/meshery-operator"` |  |
| image.tag | string | `"1.0.6"` | Pinned operator release, stamped by the sync workflow. Kept explicit because server-release chart publishing re-stamps appVersion with the server tag; empty falls back to the chart appVersion |
| imagePullSecrets | list | `[]` |  |
| ingress.annotations | object | `{}` |  |
| ingress.enabled | bool | `false` |  |
| ingress.hosts[0].host | string | `"chart-example.local"` |  |
| ingress.hosts[0].paths | list | `[]` |  |
| ingress.tls | list | `[]` |  |
| leaderElection.enabled | bool | `true` | Leader election (adds the leases Role/RoleBinding and `--enable-leader-election`) |
| meshery-broker.enabled | bool | `true` |  |
| meshery-broker.fullnameOverride | string | `"meshery-broker"` |  |
| meshery-meshsync.enabled | bool | `true` |  |
| meshery-meshsync.fullnameOverride | string | `"meshery-meshsync"` |  |
| meshsync.rbac.bindOperatorServiceAccount | bool | `true` | Also bind the read-only `meshery-meshsync-role` to `meshery-operator`, which operator releases up to 1.0.6 run MeshSync as; flips to `false` with the bump to the first operator release that honors `MESHSYNC_SERVICE_ACCOUNT` (see [MeshSync identity and RBAC](#meshsync-identity-and-rbac)) |
| meshsync.rbac.extraRules | list | `[]` | Rules appended to the otherwise read-only `meshery-meshsync-role`, e.g. `create` on `pods/exec` for the interactive terminal |
| meshsync.serviceAccount.create | bool | `true` | Create the MeshSync ServiceAccount; when `false` it must already exist in the release namespace |
| meshsync.serviceAccount.name | string | `"meshery-meshsync"` | ServiceAccount the operator runs MeshSync under, passed to the manager as `MESHSYNC_SERVICE_ACCOUNT` |
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
