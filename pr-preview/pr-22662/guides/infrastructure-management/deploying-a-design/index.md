# Deploying a Design

> Validate, dry-run, deploy, and undeploy a Meshery Design on the Kubernetes clusters in your Environments, and read the deployment results.

---

LLMS index: [llms.txt](/pr-preview/pr-22662/llms.txt)

---

Use this guide to put a [Design](/pr-preview/pr-22662/concepts/logical/designs/) onto one or more of your Kubernetes clusters, check what happened, and remove it again. Meshery deploys a design to the clusters you pick from your [Environments](/pr-preview/pr-22662/concepts/logical/environments/), so you can deploy the same design to a staging cluster and a production cluster in one pass.

## Before you begin

- A design in Meshery. Create one with the [Design Configurator](/pr-preview/pr-22662/guides/configuration-management/creating-a-meshery-design/) or [import](/pr-preview/pr-22662/guides/configuration-management/import-export-designs/) your Kubernetes manifests, Helm chart, Docker Compose file, or Kustomize package.
- At least one Kubernetes [Connection](/pr-preview/pr-22662/guides/infrastructure-management/registering-a-connection/) in the **Connected** state.
- That Connection assigned to an Environment. The deployment wizard only offers clusters that belong to an Environment. To assign one, see [Sharing a Cluster with Your Team](/pr-preview/pr-22662/guides/infrastructure-management/sharing-a-cluster-with-your-team/#create-an-environment-and-add-the-cluster), or use the **Environments** column in **Lifecycle → Connections**.
- Permission to deploy designs. If the **Deploy** action is disabled, ask your organization admin.

## Deploy a design

1. Open **Configuration → Designs**.
2. Find your design, open its actions, and click **Deploy**. The deployment wizard opens.
3. **Validate Design.** Meshery checks every component against its schema and lists any errors by component. Fix errors in the design and start again, or click **Next** to continue.
4. **Identify Environments.** Each of your Environments appears with its Kubernetes Connections. Select a whole Environment, or select individual clusters inside it. The counter next to each Environment shows how many of its clusters you selected. Select at least one cluster, then click **Next**.

    If you see "No environments found", create an Environment and assign your cluster to it first. If an Environment shows "No connections configured", click **Add a connection**.

5. **Dry Run.** Meshery sends the design to each selected cluster's API server as a server-side dry run, without changing anything, and lists the errors it would hit, by component and cluster. Common errors are a missing namespace, a missing CRD, or a field the cluster's Kubernetes version rejects.
    - **Next** stays disabled while the dry run reports errors. Fix the design, or select **Bypass errors and initiate deployment** to deploy anyway.
    - For components from Artifact Hub-sourced models, Meshery tries to install the CRDs and operators they need before it applies them, on every deployment from Meshery UI. The **Include Dependencies** checkbox on this step does not currently change that. See [Auto-Deployment of CRDs and Operators](/pr-preview/pr-22662/guides/infrastructure-management/overview/#auto-deployment-of-crds-and-operators).
6. **Finalize Deployment.** Review the design and the Environments you selected, then click **Deploy**.
7. Read the deployment summary on the last step. It lists each component's result on each cluster, with error details for any component that failed.

Meshery also posts the result to the [Notification Center](/pr-preview/pr-22662/guides/infrastructure-management/notification-management/), so you can come back to it later.

## Understand partial results

Meshery deploys components in the order their `dependsOn` entries describe, and deploys independent components in parallel. When a component fails, Meshery withholds every component that depends on it and reports each one as withheld, along with the failure that caused it. Components that don't depend on the failed one deploy normally. When you select several clusters, Meshery applies each component to every selected cluster. If a component fails on any one cluster, Meshery treats it as failed and withholds its dependents on every cluster, including the clusters where it succeeded. Components outside that dependency chain still deploy on every cluster.

To fix a partial deployment, correct the failed component and deploy the design again. Meshery applies the design declaratively, so components that already exist are updated in place.

For the full rules, see the [Deployment Engine](/pr-preview/pr-22662/concepts/architecture/deployment-engine/).

## Check a design without deploying it

From the design's actions, click **Validate Design** to run only the schema check, or **Dry Run** to run only the dry run against the clusters you select. Neither changes your clusters.

## Undeploy a design

1. Open **Configuration → Designs**.
2. Open the design's actions and click **Undeploy**.
3. Step through the same wizard. Select the clusters to remove the design from, review the dry run, and confirm.

Undeploying deletes the resources that the design describes from the selected clusters, removing dependents before what they depend on. If a component can't be removed, Meshery leaves what it depends on in place, so part of the design can keep running. Check the summary on the last step for failed and withheld components before you assume the cluster is clean. The design itself stays in Meshery, so you can deploy it again later.

<div class="alert alert-warning" role="alert"><div class="h4 alert-heading" role="heading">Deploying from mesheryctl</div>


`mesheryctl design deploy`, `mesheryctl design apply`, and `mesheryctl design undeploy` do not currently tell Meshery Server which clusters to target, and Meshery Server rejects a deployment that names no cluster. Deploy and undeploy designs from Meshery UI.
</div>


## Related

- [Creating a Meshery Design](/pr-preview/pr-22662/guides/configuration-management/creating-a-meshery-design/)
- [Importing and Exporting Designs](/pr-preview/pr-22662/guides/configuration-management/import-export-designs/)
- [Deployment Engine](/pr-preview/pr-22662/concepts/architecture/deployment-engine/)
- [Sharing a Cluster with Your Team](/pr-preview/pr-22662/guides/infrastructure-management/sharing-a-cluster-with-your-team/)

<div class="alert alert-dark" role="alert">
  <h4 class="alert-heading">Discussion Forum</h4>
  <p>Don't find an answer to your question here? Ask on the <a href="https://discuss.meshery.io/">Discussion Forum</a>.</p>
</div>
