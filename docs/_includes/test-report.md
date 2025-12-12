### END-TO-END TESTS

- Testing started at: December 12th 2025, 6:27:30 pm

**📦 Test Result Summary**

- ✅ 73 passed
- ❌ 17 failed
- ⚠️ 0 flaked
- ⏩ 14 skipped

⌛ _Duration: 8 minutes and 25 seconds_

**Overall Result**: 👎 Some tests failed.



<details>
    <summary>[Show/Hide] Test Result Details</summary>
    <div markdown="1">

| Test | Browser | Test Case | Tags | Result |
| :---: | :---: | :--- | :---: | :---: |
| 1 | chromium-meshery-provider | Verify Meshery Docker Extension Details |  | ❌ |
| 2 | chromium-meshery-provider | Add a cluster connection by uploading kubeconfig file |  | ❌ |
| 3 | chromium-meshery-provider | Transition to disconnected state and then back to connected state |  | ➖ |
| 4 | chromium-meshery-provider | Transition to ignored state and then back to connected state |  | ➖ |
| 5 | chromium-meshery-provider | Transition to not found state and then back to connected state |  | ➖ |
| 6 | chromium-meshery-provider | Delete Kubernetes cluster connections |  | ➖ |
| 7 | chromium-meshery-provider | Logout from current user session |  | ❌ |
| 8 | chromium-meshery-provider | Create a Model |  | ❌ |
| 9 | chromium-meshery-provider | Search a Model and Export it |  | ➖ |
| 10 | chromium-meshery-provider | Import a Model via File Import |  | ➖ |
| 11 | chromium-meshery-provider | Import a Model via Url Import |  | ➖ |
| 12 | chromium-meshery-provider | Import a Model via CSV Import |  | ➖ |
| 13 | chromium-meshery-provider | Common UI elements |  | ❌ |
| 14 | chromium-meshery-provider | displays public design card correctly |  | ❌ |
| 15 | chromium-meshery-provider | Test if Left Navigation Panel is displayed |  | ❌ |
| 16 | chromium-meshery-provider | Verify Meshery Design Embed Details |  | ❌ |
| 17 | chromium-meshery-provider | Configure Existing Istio adapter through Mesh Adapter URL from Management page | unstable | ⚠️ |
| 18 | chromium-meshery-provider | Ping Istio Adapter | unstable | ⚠️ |
| 19 | chromium-meshery-provider | All settings tabs |  | ❌ |
| 20 | chromium-meshery-provider | Connect to Meshery Istio Adapter and configure it |  | ➖ |
| 21 | chromium-meshery-provider | Add performance profile with load generator fortio |  | ❌ |
| 22 | chromium-meshery-provider | should identify relationships for All Relationships | relationship | ❌ |
| 23 | chromium-meshery-provider | deploys a published design to a connected cluster |  | ❌ |
| 24 | chromium-local-provider | displays public design card correctly |  | ➖ |
| 25 | chromium-local-provider | Transition to disconnected state and then back to connected state |  | ❌ |
| 26 | chromium-local-provider | Transition to ignored state and then back to connected state |  | ➖ |
| 27 | chromium-local-provider | Transition to not found state and then back to connected state |  | ➖ |
| 28 | chromium-local-provider | Delete Kubernetes cluster connections |  | ➖ |
| 29 | chromium-local-provider | Configure Existing Istio adapter through Mesh Adapter URL from Management page | unstable | ⚠️ |
| 30 | chromium-local-provider | Ping Istio Adapter | unstable | ⚠️ |
| 31 | chromium-local-provider | Connect to Meshery Istio Adapter and configure it |  | ➖ |
| 32 | chromium-local-provider | should identify relationships for All Relationships | relationship | ❌ |
| 33 | chromium-local-provider | Import a Model via CSV Import |  | ❌ |
| 34 | chromium-local-provider | Edit the configuration of a performance profile with load generator fortio and service mesh None |  | ❌ |
| 35 | chromium-local-provider | Toggle &quot;Send Anonymous Performance Results&quot; |  | ❌ |

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