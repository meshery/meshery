### END-TO-END TESTS

- Testing started at: December 11th 2025, 5:16:13 am

**📦 Test Result Summary**

- ✅ 84 passed
- ❌ 7 failed
- ⚠️ 0 flaked
- ⏩ 13 skipped

⌛ _Duration: 8 minutes and 20 seconds_

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
| 6 | chromium-meshery-provider | Search a Model and Export it |  | ❌ |
| 7 | chromium-meshery-provider | Import a Model via File Import |  | ➖ |
| 8 | chromium-meshery-provider | Import a Model via Url Import |  | ➖ |
| 9 | chromium-meshery-provider | Import a Model via CSV Import |  | ➖ |
| 10 | chromium-meshery-provider | Configure Existing Istio adapter through Mesh Adapter URL from Management page | unstable | ⚠️ |
| 11 | chromium-meshery-provider | Ping Istio Adapter | unstable | ⚠️ |
| 12 | chromium-meshery-provider | Action buttons on adapters tab |  | ❌ |
| 13 | chromium-meshery-provider | deletes a published design from the list |  | ❌ |
| 14 | chromium-meshery-provider | Connect to Meshery Istio Adapter and configure it |  | ➖ |
| 15 | chromium-meshery-provider | should identify relationships for All Relationships | relationship | ❌ |
| 16 | chromium-local-provider | displays public design card correctly |  | ➖ |
| 17 | chromium-local-provider | Transition to disconnected state and then back to connected state |  | ❌ |
| 18 | chromium-local-provider | Transition to ignored state and then back to connected state |  | ➖ |
| 19 | chromium-local-provider | Transition to not found state and then back to connected state |  | ➖ |
| 20 | chromium-local-provider | Delete Kubernetes cluster connections |  | ➖ |
| 21 | chromium-local-provider | Configure Existing Istio adapter through Mesh Adapter URL from Management page | unstable | ⚠️ |
| 22 | chromium-local-provider | Ping Istio Adapter | unstable | ⚠️ |
| 23 | chromium-local-provider | Connect to Meshery Istio Adapter and configure it |  | ➖ |
| 24 | chromium-local-provider | should identify relationships for All Relationships | relationship | ❌ |

</div>
</details>




**🔗 Relationship Tests**

<details>
    <summary>[Show/Hide] Relationship Test Details (18 tests)</summary>
    <div markdown="1">

| Kind | Type | SubType | From | To | Model | Design Name | Status |
| :---: | :---: | :---: | :---: | :---: | :---: | :--- | :---: |
| edge | binding | permission | ClusterRole | ServiceAccount | kubernetes | Understanding Relationships | ❌ |
| edge | binding | permission | Role | ServiceAccount | kubernetes | Understanding Relationships | ❌ |
| edge | non-binding | network | Service | Deployment | kubernetes | service-to-deployment-network | ✅ |
| edge | non-binding | network | Service | Deployment | kubernetes | meshery-design-fixture.json | ✅ |
| edge | non-binding | reference | ClusterRoleBinding | ClusterRole | kubernetes | meshery-design-fixture.json | ✅ |
| hierarchical | parent | alias | Container | Deployment | kubernetes | container-hierarchical-parent-alias-relationship | ✅ |
| hierarchical | parent | alias | Container | Deployment | kubernetes | deployment-configmap-reference-relationship | ✅ |
| hierarchical | parent | alias | Container | Deployment | kubernetes | service-to-deployment-network | ✅ |
| hierarchical | parent | alias | Container | Deployment | kubernetes | meshery-design-fixture.json | ✅ |
| hierarchical | parent | alias | Container | Deployment | kubernetes | Understanding Relationships | ✅ |
| hierarchical | parent | alias | Container | Pod | kubernetes | container-hierarchical-parent-alias-relationship | ✅ |
| hierarchical | parent | inventory | * | Namespace | kubernetes | hierarchical-parent-namespace-relationship | ✅ |
| hierarchical | parent | inventory | * | Namespace | kubernetes | service-to-deployment-network | ✅ |
| hierarchical | parent | inventory | * | Namespace | kubernetes | pv-pvc-edge-non-binding-reference-relationship | ✅ |
| hierarchical | parent | inventory | * | Namespace | kubernetes | Understanding Relationships | ✅ |
| hierarchical | sibling | matchlabels | ClusterRole | ClusterRole | kubernetes | meshery-design-fixture.json | ✅ |
| hierarchical | sibling | matchlabels | Deployment | Deployment | kubernetes | meshery-design-fixture.json | ✅ |
| hierarchical | sibling | matchlabels | Service | Service | kubernetes | meshery-design-fixture.json | ✅ |

</div>
</details>


<!-- To see the full report, please visit our CI/CD pipeline with reporter. -->