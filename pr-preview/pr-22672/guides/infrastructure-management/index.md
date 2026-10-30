# Infrastructure Management

> Guides for understanding Meshery's ability to manage infrastructure.

---

LLMS index: [llms.txt](/pr-preview/pr-22672/llms.txt)

---

Guides for understanding Meshery's ability to manage infrastructure.

---

Section pages:

- [Overview](/pr-preview/pr-22672/guides/infrastructure-management/overview/): Meshery provides the ability to manage infrastructure for agility, maintainability, diversity, reliability and isolation, security, and speed.
- [Registering a Connection](/pr-preview/pr-22672/guides/infrastructure-management/registering-a-connection/): Use the Connection Wizard to create and update Connections - Kubernetes clusters, Grafana, Prometheus, and more - in your Meshery deployment.
- [Configuring Meshery Operator, MeshSync, and Broker](/pr-preview/pr-22672/guides/infrastructure-management/configuring-operator-meshsync-broker/): Tune how Meshery deploys and operates Meshery Operator, MeshSync, and Meshery Broker on your clusters - deployment modes, discovery scope, secret handling, service networking, and the behavioral impact of every setting.
- [Bringing Existing Infrastructure Under Meshery Management](/pr-preview/pr-22672/guides/infrastructure-management/managing-existing-infrastructure/): What happens when you point Meshery at a cluster that is already running workloads - how brownfield discovery works, what to evaluate before you connect, the end-to-end process, and the caveats.
- [Deploying a Design](/pr-preview/pr-22672/guides/infrastructure-management/deploying-a-design/): Validate, dry-run, deploy, and undeploy a Meshery Design on the Kubernetes clusters in your Environments, and read the deployment results.
- [Sharing a Cluster with Your Team](/pr-preview/pr-22672/guides/infrastructure-management/sharing-a-cluster-with-your-team/): Give your teammates access to a Kubernetes cluster you manage in Meshery by grouping it into an Environment, attaching the Environment to a Workspace, and assigning teams to the Workspace.
- [Kubernetes Connection Lifecycle](/pr-preview/pr-22672/guides/infrastructure-management/kubernetes-connection-lifecycle/): What happens end-to-end when Meshery manages a Kubernetes cluster - the connection states, the in-cluster components (Operator, Broker/NATS, MeshSync), MeshSync deployment modes, how Meshery reaches and authenticates to the Broker, the built-in Diagnostics, and how to operate and troubleshoot it all.
- [Envoy WASM Filter Management](/pr-preview/pr-22672/guides/infrastructure-management/filter-management/): Import, inspect, publish, and remove Envoy WebAssembly (WASM) filters with Meshery UI and mesheryctl.
- [GitOps with Meshery](/pr-preview/pr-22672/guides/infrastructure-management/gitops-with-meshery/): Add Meshery's GitHub Actions to your pipelines to render a visual snapshot of your infrastructure on each pull request and to catch performance regressions before you merge.
- [Managing Connections](/pr-preview/pr-22672/guides/infrastructure-management/lifecycle-management/): Find, inspect, transition, reconfigure, and remove the Connections Meshery manages, from Meshery UI or mesheryctl.
- [Managing Events with Notification Center](/pr-preview/pr-22672/guides/infrastructure-management/notification-management/): Meshery tracks operations performed on your infrastructure and workloads, and provides notification of environment issues, application conflicts with infrastructure configuration, policy violations, and so on.
