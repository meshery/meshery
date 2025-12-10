### END-TO-END TESTS

- Testing started at: December 10th 2025, 4:59:20 pm

**📦 Test Result Summary**

- ✅ 86 passed
- ❌ 8 failed
- ⚠️ 0 flaked
- ⏩ 10 skipped

⌛ _Duration: 8 minutes and 3 seconds_

**Overall Result**: 👎 Some tests failed.



<details>
    <summary>[Show/Hide] Test Result Details</summary>
    <div markdown="1">

| Test | Browser | Test Case | Tags | Result |
| :---: | :---: | :--- | :---: | :---: |
| 1 | chromium-meshery-provider | Verify Performance Analysis Details |  | ❌ |
| 2 | chromium-meshery-provider | displays published design card correctly |  | ❌ |
| 3 | chromium-meshery-provider | Add a cluster connection by uploading kubeconfig file |  | ❌ |
| 4 | chromium-meshery-provider | Transition to disconnected state and then back to connected state |  | ➖ |
| 5 | chromium-meshery-provider | Transition to ignored state and then back to connected state |  | ➖ |
| 6 | chromium-meshery-provider | Transition to not found state and then back to connected state |  | ➖ |
| 7 | chromium-meshery-provider | Delete Kubernetes cluster connections |  | ➖ |
| 8 | chromium-meshery-provider | should edit design in Design Configurator |  | ❌ |
| 9 | chromium-meshery-provider | Verify Kanvas Details |  | ❌ |
| 10 | chromium-meshery-provider | Configure Existing Istio adapter through Mesh Adapter URL from Management page | unstable | ⚠️ |
| 11 | chromium-meshery-provider | Ping Istio Adapter | unstable | ⚠️ |
| 12 | chromium-meshery-provider | Connect to Meshery Istio Adapter and configure it |  | ➖ |
| 13 | chromium-meshery-provider | should identify relationships for All Relationships | relationship | ❌ |
| 14 | chromium-local-provider | displays public design card correctly |  | ➖ |
| 15 | chromium-local-provider | Transition to disconnected state and then back to connected state |  | ❌ |
| 16 | chromium-local-provider | Transition to ignored state and then back to connected state |  | ➖ |
| 17 | chromium-local-provider | Transition to not found state and then back to connected state |  | ➖ |
| 18 | chromium-local-provider | Delete Kubernetes cluster connections |  | ➖ |
| 19 | chromium-local-provider | Configure Existing Istio adapter through Mesh Adapter URL from Management page | unstable | ⚠️ |
| 20 | chromium-local-provider | Ping Istio Adapter | unstable | ⚠️ |
| 21 | chromium-local-provider | Connect to Meshery Istio Adapter and configure it |  | ➖ |
| 22 | chromium-local-provider | should identify relationships for All Relationships | relationship | ❌ |

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