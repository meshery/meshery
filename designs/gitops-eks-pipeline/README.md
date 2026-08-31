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
in-cluster components and leave the `ack-system` namespace out — nothing else in
the design depends on those resources at apply time.

## Before you deploy

The design applies cleanly only where the CRDs it references are installed:

| Components | Requires |
|---|---|
| `ec2/eks/iam.services.k8s.aws` | ACK controllers for EC2, EKS and IAM |
| `argoproj.io` | Argo CD Operator |
| `gateway.networking.k8s.io` | Gateway API CRDs + AWS Load Balancer Controller ≥ 2.9 |
| `monitoring.coreos.com` | kube-prometheus-stack (Prometheus Operator) |
| `*.k8s.elastic.co` | ECK operator |

Placeholders to replace — every one is a `111122223333`-style dummy, so nothing
here resolves to a real account:

- AWS account id `111122223333` in every role ARN
- the ACM certificate ARN on the `boutique-gateway` Gateway
- the OIDC provider URL and `EXAMPLED539D4633E53DE1B716D3041E` id in the
  `external-dns-irsa` trust policy
- `shop.example.com` and the `--domain-filter` on `external-dns`
- `https://github.com/example-org/boutique-gitops.git` and
  `ghcr.io/example-org/boutique-app` on the Argo CD Application
- `ami-0c02fb55956c7d316` on the bastion (region-specific)
- `bastion-sg` allows SSH from `10.0.0.0/16`; narrow it to your own range

Three Secrets are referenced but deliberately not included — create them
yourself rather than committing them:

```bash
kubectl -n monitoring create secret generic alertmanager-slack \
  --from-literal=webhook-url='https://hooks.slack.com/services/...'
kubectl -n monitoring create secret generic grafana-admin --from-literal=password='...'
kubectl -n boutique-app create secret generic postgres-orders \
  --from-literal=username=orders --from-literal=password='...'
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
- **Model versions need not match your registry exactly.** Component lookup
  prefers an exact `model.version` match and falls back to any version of the
  same model, so the design keeps resolving as the registry moves on.
- **Every annotation carries a real apiVersion**, which is why the diagram
  furniture is drawn from `meshery-shapes` and `meshery-flowchart` rather than
  `meshery-dev-icons`. The validation stage resolves *every* component —
  annotations are only set aside later, during provisioning — and it rejects an
  empty apiVersion outright. A single dev-icon would therefore abort the whole
  deploy, which is a fine trade for a diagram-only design but not for this one.
