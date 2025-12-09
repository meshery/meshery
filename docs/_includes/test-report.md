### END-TO-END TESTS

- Testing started at: December 9th 2025, 12:40:06 pm

**📦 Test Result Summary**

- ✅ 88 passed
- ❌ 2 failed
- ⚠️ 0 flaked
- ⏩ 10 skipped

⌛ _Duration: 6 minutes and 34 seconds_

**Overall Result**: 👎 Some tests failed.



<details>
    <summary>[Show/Hide] Test Result Details</summary>
    <div markdown="1">

| Test | Browser | Test Case | Tags | Result |
| :---: | :---: | :--- | :---: | :---: |
| 1 | chromium-meshery-provider | Add a cluster connection by uploading kubeconfig file |  | ❌ |
| 2 | chromium-meshery-provider | Transition to disconnected state and then back to connected state |  | ➖ |
| 3 | chromium-meshery-provider | Transition to ignored state and then back to connected state |  | ➖ |
| 4 | chromium-meshery-provider | Transition to not found state and then back to connected state |  | ➖ |
| 5 | chromium-meshery-provider | Delete Kubernetes cluster connections |  | ➖ |
| 6 | chromium-meshery-provider | Configure Existing Istio adapter through Mesh Adapter URL from Management page | unstable | ⚠️ |
| 7 | chromium-meshery-provider | Ping Istio Adapter | unstable | ⚠️ |
| 8 | chromium-meshery-provider | Connect to Meshery Istio Adapter and configure it |  | ➖ |
| 9 | chromium-local-provider | displays public design card correctly |  | ➖ |
| 10 | chromium-local-provider | Transition to disconnected state and then back to connected state |  | ❌ |
| 11 | chromium-local-provider | Transition to ignored state and then back to connected state |  | ➖ |
| 12 | chromium-local-provider | Transition to not found state and then back to connected state |  | ➖ |
| 13 | chromium-local-provider | Delete Kubernetes cluster connections |  | ➖ |
| 14 | chromium-local-provider | Configure Existing Istio adapter through Mesh Adapter URL from Management page | unstable | ⚠️ |
| 15 | chromium-local-provider | Ping Istio Adapter | unstable | ⚠️ |
| 16 | chromium-local-provider | Connect to Meshery Istio Adapter and configure it |  | ➖ |

</div>
</details>




**🔗 Relationship Tests**

<details>
    <summary>[Show/Hide] Relationship Test Details (13 tests)</summary>
    <div markdown="1">

| Kind | Type | SubType | From | To | Model | Design Name | Status |
| :---: | :---: | :---: | :---: | :---: | :---: | :--- | :---: |
| edge | non-binding | network | Service | Deployment | kubernetes | service-to-deployment-network | ✅ |
| edge | non-binding | network | Service | Deployment | kubernetes | meshery-design | ✅ |
| edge | non-binding | reference | ClusterRoleBinding | ClusterRole | kubernetes | meshery-design | ✅ |
| edge | non-binding | reference | ConfigMap | Deployment | kubernetes | deployment-configmap-reference-relationship | ✅ |
| hierarchical | parent | alias | Container | Deployment | kubernetes | service-to-deployment-network | ✅ |
| hierarchical | parent | alias | Container | Deployment | kubernetes | container-hierarchical-parent-alias-relationship | ✅ |
| hierarchical | parent | alias | Container | Deployment | kubernetes | meshery-design | ✅ |
| hierarchical | parent | alias | Container | Pod | kubernetes | container-hierarchical-parent-alias-relationship | ✅ |
| hierarchical | parent | inventory | * | Namespace | kubernetes | service-to-deployment-network | ✅ |
| hierarchical | parent | inventory | * | Namespace | kubernetes | hierarchical-parent-namespace-relationship | ✅ |
| hierarchical | sibling | matchlabels | ClusterRole | ClusterRole | kubernetes | meshery-design | ✅ |
| hierarchical | sibling | matchlabels | Deployment | Deployment | kubernetes | meshery-design | ✅ |
| hierarchical | sibling | matchlabels | Service | Service | kubernetes | meshery-design | ✅ |

</div>
</details>


<!-- To see the full report, please visit our CI/CD pipeline with reporter. -->