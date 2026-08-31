# GitOps Pipeline on Amazon EKS

A Meshery design that recreates an end-to-end GitOps architecture — commit to
production on a private EKS cluster — as a Kanvas canvas that is also
deployable.

Import `design.yaml` into Kanvas (**Designs → Import → Upload File**), or from
the CLI:

```bash
mesheryctl design import -f designs/gitops-eks-pipeline/design.yaml
```

## What is in it

87 components: 57 that Meshery applies to a cluster, and 30 annotations that
carry the diagram's structure (region/VPC/subnet boundaries, the CI lane, the
actors) without ever being deployed.

| Region of the canvas | Components |
|---|---|
| Continuous Integration | GitHub, GitHub Actions, checkout → build → Trivy scan → push, GHCR — all annotations; GitHub Actions has no cluster-side resource |
| AWS network | `VPC`, four `Subnet`s (2 public / 2 private), `InternetGateway`, `NATGateway`, `ElasticIPAddress`, two `RouteTable`s, `SecurityGroup`, bastion `Instance` — AWS Controllers for Kubernetes (ACK) |
| Cluster | `Cluster` (private endpoint, public access disabled) and `Nodegroup` in the private subnets |
| Continuous Delivery | `ArgoCD` control plane, an `Application` carrying the Image Updater annotations, and the `argocd-image-updater` Deployment |
| Ingress | `GatewayClass` → `Gateway` → `HTTPRoute`, the AWS Load Balancer Controller, and `external-dns` |
| Application | four service `Deployment`s, `redis-cart` and `postgres-orders` `StatefulSet`s, their `Service`s, a load generator, and a `HorizontalPodAutoscaler` |
| Logging | ECK — `Elasticsearch`, `Kibana`, `Beat` (filebeat), and the operator `StatefulSet` |
| Monitoring | `Prometheus`, `Alertmanager`, an `AlertmanagerConfig` that routes to Slack, `ServiceMonitor`, `PrometheusRule`, and Grafana |
| Identity | an IAM `Role` for IRSA, an EKS `PodIdentityAssociation`, the `eks-pod-identity-agent` DaemonSet, and IRSA-annotated `ServiceAccount`s |

### Two provisioning paths, deliberately

The source diagram provisions AWS with Terraform and keeps state in S3. Those
two remain annotations, because Terraform state is not a cluster resource.

The AWS layer is *also* modelled as ACK custom resources, so the design can
provision it from inside a management cluster if you would rather not run
Terraform. **Pick one.** If your VPC and cluster already exist, deploy only the
in-cluster components and leave the `ack-system` namespace out.

That split is clean for networking, but **not for identity**. Dropping
`ack-system` also drops the one IAM role the design creates, and several
workloads carry a `role-arn` annotation regardless of which path you take.
Kubernetes will happily create a ServiceAccount whose annotation points at a
role that does not exist — nothing fails at apply time, and the workload then
cannot reach AWS at runtime. Read the next section before you deploy either
way.

### Identity is mostly a prerequisite, not part of the design

The design references six IAM roles and creates exactly one. This is
deliberate — a role's trust policy is specific to your account and cluster OIDC
provider, so the design cannot supply a working one — but it means the other
five must exist before the workloads that name them will function:

| Role | Named by | Supplied by |
|---|---|---|
| `external-dns-irsa` | `external-dns` ServiceAccount | the design, in `ack-system` |
| `aws-load-balancer-controller-irsa` | `aws-load-balancer-controller` ServiceAccount | **you** |
| `argocd-image-updater-irsa` | `argocd-image-updater` ServiceAccount | **you** |
| `boutique-app-pod-identity` | the `PodIdentityAssociation` | **you** |
| `eks-cluster-role` | the EKS `Cluster` | **you**, before the cluster |
| `eks-node-role` | the `Nodegroup` | **you**, before the node group |

`external-dns-irsa` shows the shape to copy: an `iam.services.k8s.aws` `Role`
with a web-identity trust policy naming your OIDC provider and the
`system:serviceaccount:<ns>:<name>` subject, plus a least-privilege inline
policy. The cluster and node roles are the standard EKS service roles and have
to exist before the resources that reference them, so they are outside what
this design can bootstrap.

## Before you deploy

The design applies cleanly only where the CRDs it references are installed:

| Components | Requires |
|---|---|
| `ec2/eks/iam.services.k8s.aws` | ACK controllers for EC2, EKS and IAM |
| `argoproj.io` | Argo CD Operator |
| `gateway.networking.k8s.io` | Gateway API CRDs + AWS Load Balancer Controller ≥ 2.9 |
| `monitoring.coreos.com` | kube-prometheus-stack (Prometheus Operator) |
| `*.k8s.elastic.co` | ECK operator |

The design is seeded with dummy values throughout, and no ARN in it points at a
real AWS account — but they are not all `111122223333`-style, so work the list
rather than grepping for one pattern. Replace every one of these:

- AWS account ID `111122223333`, in every role ARN
- the ACM certificate ARN on the `boutique-gateway` Gateway
- the OIDC provider URL and `EXAMPLED539D4633E53DE1B716D3041E` ID in the
  `external-dns-irsa` trust policy
- `Z0EXAMPLEHOSTEDZONE` in the `external-dns-irsa` inline policy — the hosted
  zone External DNS is allowed to write to
- `shop.example.com`, on both the Gateway listener and the `--domain-filter`
  argument to `external-dns`
- `https://github.com/example-org/boutique-gitops.git` and
  `ghcr.io/example-org/boutique-app` on the Argo CD Application
- `us-east-1`, and the `us-east-1a` / `us-east-1b` availability zones on the
  subnets and node group — the region is baked into the ACK resources, the role
  ARNs and the OIDC URL alike
- `ami-0c02fb55956c7d316` on the bastion. Unlike the rest of this list it is a
  real Amazon Linux 2 image, but AMI IDs are per-region: it resolves only in
  `us-east-1` and must change with the region above.
- the `10.0.0.0/16` VPC CIDR and its four subnet CIDRs, if they collide with a
  network you already peer with
- `bastion-sg` allows SSH from `10.0.0.0/16`; narrow it to your own range

Four Secrets are referenced but deliberately not included — create them
yourself rather than committing them. Put each *secret* value in a file first:
a `--from-literal` lands in your shell history and is briefly visible in `ps`
while the command runs. Non-sensitive values are fine as literals, which is why
the postgres username below still is one.

```bash
# Slack webhook for the slack-alerts AlertmanagerConfig
kubectl -n monitoring create secret generic alertmanager-slack \
  --from-file=webhook-url=./slack-webhook.txt

# Grafana admin password
kubectl -n monitoring create secret generic grafana-admin \
  --from-file=password=./grafana-password.txt

# postgres-orders credentials (the username is not sensitive)
kubectl -n boutique-app create secret generic postgres-orders \
  --from-literal=username=orders \
  --from-file=password=./postgres-password.txt

# GHCR credentials for Argo CD Image Updater. Required: the Application's
# `app.pull-secret: pullsecret:argocd/ghcr-creds` annotation resolves to this,
# and without it Image Updater cannot read tags from GHCR — the CD half of the
# pipeline then does nothing, without failing loudly. Needs a PAT with
# `read:packages`, and must be a `kubernetes.io/dockerconfigjson` secret.
# `create secret docker-registry` cannot read the password from a file, and
# `--docker-password="$(< file)"` is expanded by the shell before kubectl starts,
# so the token would land in argv. Log in against a throwaway DOCKER_CONFIG
# instead and build the Secret from the config file it writes. The temporary
# config also keeps other registries' credentials, and any credsStore
# indirection, out of the Secret.
# Runs in a subshell so the EXIT trap fires when the block ends rather than when
# your shell does, and so the DOCKER_CONFIG override never escapes into the
# calling session. set -e aborts on the first failure; the trap still runs, so
# neither the token nor the config it lives in survives an interrupted setup.
(
  set -e
  export DOCKER_CONFIG="$(mktemp -d)"
  trap 'rm -rf "$DOCKER_CONFIG"; rm -f ./ghcr-pat.txt' EXIT
  docker login ghcr.io --username '<github-username>' --password-stdin < ./ghcr-pat.txt
  kubectl -n argocd create secret generic ghcr-creds \
    --type=kubernetes.io/dockerconfigjson \
    --from-file=.dockerconfigjson="$DOCKER_CONFIG/config.json"
  docker logout ghcr.io
)

# Remove the remaining credential files (the GHCR PAT is already gone, removed
# by the subshell's trap above). shred is GNU coreutils and is absent on macOS,
# where this would otherwise fail and quietly leave the files behind. Note that
# neither command reliably destroys data on a copy-on-write or wear-levelled
# filesystem — write the files to a tmpfs mount if that matters to you.
if command -v shred >/dev/null 2>&1; then
  shred -u ./slack-webhook.txt ./grafana-password.txt ./postgres-password.txt
else
  rm -f ./slack-webhook.txt ./grafana-password.txt ./postgres-password.txt
fi
```

## Notes on the modelling

- **Relationships are left empty.** Kanvas's relationship evaluation engine
  derives hierarchical and network edges on load, so hand-authoring them would
  only go stale.
- **Namespace comes from `configuration.metadata.namespace`**, which is what
  `getNamespaceForComponent` reads when Meshery applies a component.
- **Gateway API components use the `gko` model**, the only model in the
  registry that carries the full `gateway.networking.k8s.io/v1` set. The
  apiVersion, kind and configuration are what get applied, so the manifests are
  correct regardless of which model supplies the icon.
- **Every `model.version` here matches the registry exactly, deliberately.**
  `FindCompDefinitionWithVersion` breaks on an exact `model.version` match but
  otherwise assigns `match` on each iteration, so with no exact match you get
  whichever entity the registry returned last. That is a best-effort fallback,
  not a stable one — it can change as the registry gains versions. Keep the
  versions aligned when you want predictable validation and provisioning.
- **Every annotation carries a real apiVersion**, which is why the diagram
  furniture is drawn from `meshery-shapes` and `meshery-flowchart` rather than
  `meshery-dev-icons`. The validation stage resolves *every* component —
  annotations are only set aside later, during provisioning — and it rejects an
  empty apiVersion outright. A single dev-icon would therefore abort the whole
  deploy, which is a fine trade for a diagram-only design but not for this one.
