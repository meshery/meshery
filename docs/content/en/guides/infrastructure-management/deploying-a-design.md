---
title: Deploying a Design
categories: [infrastructure]
weight: -3
description: Validate, dry-run, deploy, and undeploy a Meshery Design on the Kubernetes clusters in your Environments, and read the deployment results.
---

Use this guide to put a [Design]({{< ref "concepts/logical/designs.md" >}}) onto one or more of your Kubernetes clusters, check what happened, and remove it again. Meshery deploys a design to the clusters you pick from your [Environments]({{< ref "concepts/logical/environments.md" >}}), so you can deploy the same design to a staging cluster and a production cluster in one pass.

## Before you begin

- A design in Meshery. Create one with the [Design Configurator]({{< ref "guides/configuration-management/creating-a-meshery-design.md" >}}) or [import]({{< ref "guides/configuration-management/import-export-designs.md" >}}) your Kubernetes manifests, Helm chart, Docker Compose file, or Kustomize package.
- At least one Kubernetes [Connection]({{< ref "guides/infrastructure-management/registering-a-connection.md" >}}) in the **Connected** state.
- That Connection assigned to an Environment. The deployment wizard only offers clusters that belong to an Environment. To assign one, see [Sharing a Cluster with Your Team]({{< ref "guides/infrastructure-management/sharing-a-cluster-with-your-team.md#create-an-environment-and-add-the-cluster" >}}), or use the **Environments** column in **Lifecycle → Connections**.
- Permission to deploy designs. If the **Deploy** action is disabled, ask your organization admin.

## Deploy a design

1. Open **Configuration → Designs**.
2. Find your design, open its actions, and click **Deploy**. The deployment wizard opens.
3. **Validate Design.** Meshery checks every component against its schema and lists any errors by component. Fix errors in the design and start again, or click **Next** to continue.
4. **Identify Environments.** Each of your Environments appears with its Kubernetes Connections. Select a whole Environment, or select individual clusters inside it. The counter next to each Environment shows how many of its clusters you selected. Select at least one cluster, then click **Next**.

    If you see "No environments found", create an Environment and assign your cluster to it first. If an Environment shows "No connections configured", click **Add a connection**.

5. **Dry Run.** Meshery sends the design to each selected cluster's API server as a server-side dry run, without changing anything, and lists the errors it would hit, by component and cluster. Common errors are a missing namespace, a missing CRD, or a field the cluster's Kubernetes version rejects.
    - **Next** stays disabled while the dry run reports errors. Fix the design, or select **Bypass errors and initiate deployment** to deploy anyway.
    - For components from Artifact Hub-sourced models, Meshery tries to install the CRDs and operators they need before it applies them, on every deployment from Meshery UI. The **Include Dependencies** checkbox on this step does not currently change that. See [Auto-Deployment of CRDs and Operators]({{< ref "guides/infrastructure-management/overview.md#auto-deployment-of-crds-and-operators" >}}).
    - **Export Report:** Click **Export Report** to download the dry run results as a standardized JSON report for auditing, sharing, or CI/CD integration. See [Exporting Dry Run Validation Reports](#exporting-dry-run-validation-reports) below.
6. **Finalize Deployment.** Review the design and the Environments you selected, then click **Deploy**.
7. Read the deployment summary on the last step. It lists each component's result on each cluster, with error details for any component that failed.

Meshery also posts the result to the [Notification Center]({{< ref "guides/infrastructure-management/notification-management/index.md" >}}), so you can come back to it later.

## Understand partial results

Meshery deploys components in the order their `dependsOn` entries describe, and deploys independent components in parallel. When a component fails, Meshery withholds every component that depends on it and reports each one as withheld, along with the failure that caused it. Components that don't depend on the failed one deploy normally. When you select several clusters, Meshery applies each component to every selected cluster. If a component fails on any one cluster, Meshery treats it as failed and withholds its dependents on every cluster, including the clusters where it succeeded. Components outside that dependency chain still deploy on every cluster.

To fix a partial deployment, correct the failed component and deploy the design again. Meshery applies the design declaratively, so components that already exist are updated in place.

For the full rules, see the [Deployment Engine]({{< ref "concepts/architecture/deployment-engine/index.md" >}}).

## Check a design without deploying it

From the design's actions, click **Validate Design** to run only the schema check, or **Dry Run** to run only the dry run against the clusters you select. Neither changes your clusters.

### Exporting Dry Run Validation Reports

After running a Dry Run, you can export the validation outcomes as a structured JSON report by clicking **Export Report** in the Dry Run results interface. The exported report makes it easy to archive results, share findings with teammates, and automate validation checks in CI/CD pipelines without copying details manually.

#### Example Report

```json
{
  "schemaVersion": "1.0.0",
  "reportType": "meshery.dry-run-validation",
  "design": {
    "id": "ecommerce-app",
    "name": "ecommerce-app"
  },
  "run": {
    "startTime": "2026-10-28T10:42:21Z",
    "completedTime": "2026-10-28T10:42:25Z",
    "status": "failed",
    "durationSeconds": 4.2
  },
  "summary": {
    "resourcesChecked": 4,
    "resourcesPassed": 3,
    "resourcesFailed": 1,
    "warnings": 0
  },
  "results": [
    {
      "name": "frontend",
      "type": "Deployment",
      "status": "valid",
      "errors": []
    },
    {
      "name": "service",
      "type": "Service",
      "status": "valid",
      "errors": []
    },
    {
      "name": "redis",
      "type": "StatefulSet",
      "status": "valid",
      "errors": []
    },
    {
      "name": "postgres",
      "type": "StatefulSet",
      "status": "failed",
      "errors": [
        {
          "message": "StorageClass \"fast-ssd\" not found",
          "errorType": "ComponentError",
          "fieldPath": "spec.storageClassName",
          "resourceKind": "StorageClass",
          "apiVersion": "apps/v1"
        }
      ]
    }
  ]
}
```

#### Report Fields

| Field | Type | Description |
|---|---|---|
| `schemaVersion` | string | Version of the report schema (currently `"1.0.0"`). |
| `reportType` | string | Identifies the report domain (`"meshery.dry-run-validation"`). |
| `design.id` | string | Unique identifier of the evaluated design. |
| `design.name` | string | Name of the evaluated design. |
| `run.startTime` | string | ISO-8601 timestamp when validation began. |
| `run.completedTime` | string | ISO-8601 timestamp when validation finished. |
| `run.status` | string | Overall validation outcome: `"success"`, `"failed"`, or `"incomplete"`. |
| `run.durationSeconds` | number | Validation duration in seconds. |
| `summary.resourcesChecked` | number | Total number of configurable resources analyzed. |
| `summary.resourcesPassed` | number | Count of resources that passed validation. |
| `summary.resourcesFailed` | number | Count of resources with validation errors. |
| `summary.warnings` | number | Count of validation warnings. |
| `results` | array | Resource-level outcomes. |
| `results[].name` | string | Name of the resource in the design. |
| `results[].type` | string | Component kind (e.g., `Deployment`, `Service`, `StatefulSet`). |
| `results[].status` | string | Resource status (`"valid"` or `"failed"`). |
| `results[].errors` | array | List of errors encountered for this resource. |
| `results[].errors[].message` | string | Error message or status string. |
| `results[].errors[].errorType` | string | Category of error (e.g., `"ComponentError"`). |
| `results[].errors[].fieldPath` | string | JSON/spec path associated with the issue. |
| `results[].errors[].resourceKind` | string | Kind of resource referenced by the error, when available. |
| `results[].errors[].apiVersion` | string | API version of the component, when available. |

## Undeploy a design

1. Open **Configuration → Designs**.
2. Open the design's actions and click **Undeploy**.
3. Step through the same wizard. Select the clusters to remove the design from, review the dry run, and confirm.

Undeploying deletes the resources that the design describes from the selected clusters, removing dependents before what they depend on. If a component can't be removed, Meshery leaves what it depends on in place, so part of the design can keep running. Check the summary on the last step for failed and withheld components before you assume the cluster is clean. The design itself stays in Meshery, so you can deploy it again later.

{{% alert color="warning" title="Deploying from mesheryctl" %}}
`mesheryctl design deploy`, `mesheryctl design apply`, and `mesheryctl design undeploy` do not currently tell Meshery Server which clusters to target, and Meshery Server rejects a deployment that names no cluster. Deploy and undeploy designs from Meshery UI.
{{% /alert %}}

## Related

- [Creating a Meshery Design]({{< ref "guides/configuration-management/creating-a-meshery-design.md" >}})
- [Importing and Exporting Designs]({{< ref "guides/configuration-management/import-export-designs.md" >}})
- [Deployment Engine]({{< ref "concepts/architecture/deployment-engine/index.md" >}})
- [Sharing a Cluster with Your Team]({{< ref "guides/infrastructure-management/sharing-a-cluster-with-your-team.md" >}})

{{< discuss >}}
