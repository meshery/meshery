# Deploy Azure resources with Meshery

> Install Azure Service Operator in your cluster and deploy Azure resources, such as Storage Accounts and Key Vaults, from a Meshery Design.

---

LLMS index: [llms.txt](/pr-preview/pr-22673/llms.txt)

---

In this tutorial, you deploy Azure resources from a Meshery Design. Meshery manages Azure resources through [Azure Service Operator (ASO)](https://azure.github.io/azure-service-operator/), which represents each Azure resource, such as a Storage Account, Key Vault, or SQL Server, as a Kubernetes custom resource. You install ASO in your cluster, give it an Azure identity, and then design and deploy Azure resources from Meshery.

## Before you begin

- Meshery, deployed by following the [installation guides](/pr-preview/pr-22673/installation/), with a Kubernetes cluster connected to it. In Meshery UI, the cluster appears as **Connected** under **Lifecycle → Connections**.
- `kubectl`, configured for the same cluster.
- An Azure subscription to create resources in.
- The [Azure CLI](https://learn.microsoft.com/cli/azure/install-azure-cli), signed in with `az login`.

## 1. Create a service principal

ASO signs in to Azure as a service principal. Create one with the Contributor role on your subscription:

```bash
az ad sp create-for-rbac -n azure-service-operator --role contributor --scopes /subscriptions/<AZURE_SUBSCRIPTION_ID>
```

The command prints `appId`, `displayName`, `password`, and `tenant`. Export them, with your subscription ID, for the next steps:

```bash
export AZURE_CLIENT_ID=<appId>
export AZURE_CLIENT_SECRET=<password>
export AZURE_TENANT_ID=<tenant>
export AZURE_SUBSCRIPTION_ID=<subscription-id>
```

## 2. Install cert-manager

ASO needs [cert-manager](https://cert-manager.io/) for its webhook certificates. Install it, and wait for its Pods to be ready:

```bash
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.14.1/cert-manager.yaml
kubectl wait --for=condition=Available deployment --all -n cert-manager --timeout=180s
```

## 3. Install Azure Service Operator

1. Apply the ASO operator manifest:

    ```bash
    kubectl apply --server-side=true -f https://github.com/Azure/azure-service-operator/releases/download/v2.13.0/azureserviceoperator_v2.13.0.yaml
    ```

    This manifest installs the operator but no resource CRDs, because its `--crd-pattern` is empty. Meshery installs the CRDs that your design's Azure components need when you deploy the design. See [Auto-Deployment of CRDs and Operators](/pr-preview/pr-22673/guides/infrastructure-management/overview/#auto-deployment-of-crds-and-operators).

2. Create the secret ASO reads its Azure identity from. ASO looks for a secret named `aso-controller-settings` in the `azureserviceoperator-system` namespace:

    ```bash
    kubectl create secret generic aso-controller-settings \
      --namespace azureserviceoperator-system \
      --from-literal=AZURE_SUBSCRIPTION_ID="$AZURE_SUBSCRIPTION_ID" \
      --from-literal=AZURE_TENANT_ID="$AZURE_TENANT_ID" \
      --from-literal=AZURE_CLIENT_ID="$AZURE_CLIENT_ID" \
      --from-literal=AZURE_CLIENT_SECRET="$AZURE_CLIENT_SECRET"
    ```

3. Confirm that the operator is running:

    ```bash
    kubectl get pods -n azureserviceoperator-system
    ```

## 4. Deploy an Azure resource

1. In Meshery UI, open **Configuration → Designs** and click **Create Design**.
2. Choose the Azure model, then add and configure the component for the resource you want, such as a Storage Account. Set its resource group, location, and name.
3. Save the design.
4. Deploy the design to your cluster. See [Deploying a Design](/pr-preview/pr-22673/guides/infrastructure-management/deploying-a-design/).
5. In the Azure portal, open the resource group and confirm that the new resource is there.

To start from a ready-made design instead, open the [Meshery Catalog](https://meshery.io/catalog), filter by **Azure**, and clone a design.

## Clean up

Undeploy the design from Meshery. ASO deletes the Azure resources that the design created. To remove the service principal, run `az ad sp delete --id "$AZURE_CLIENT_ID"`.

## What you learned

- ASO turns Kubernetes custom resources into Azure resources, using the identity in `aso-controller-settings`.
- Meshery deploys Azure resources the same way it deploys any other design.

For a walkthrough of one resource, see [Deploy Azure Storage Account with Meshery](/pr-preview/pr-22673/guides/tutorials/azure/deploy-azure-storage-account-with-meshery/). For more about ASO, see the [ASO documentation](https://azure.github.io/azure-service-operator/).
