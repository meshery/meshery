### END-TO-END TESTS

- Testing started at: December 30th 2025, 11:54:00 pm

**📦 Test Result Summary**

- ✅ 34 passed
- ❌ 13 failed
- ⚠️ 0 flaked
- ⏩ 6 skipped

⌛ _Duration: 6 minutes and 50 seconds_

**Overall Result**: 👎 Some tests failed.



<details>
    <summary>[Show/Hide] Test Result Details</summary>
    <div markdown="1">

| Test | Browser | Test Case | Tags | Result |
| :---: | :---: | :--- | :---: | :---: |
| 1 | remote-setup | authenticate as Meshery provider |  | ❌ |
| 2 | chromium-local-provider | Transition to disconnected state and then back to connected state |  | ❌ |
| 3 | chromium-local-provider | Transition to ignored state and then back to connected state |  | ➖ |
| 4 | chromium-local-provider | Transition to not found state and then back to connected state |  | ➖ |
| 5 | chromium-local-provider | Delete Kubernetes cluster connections |  | ➖ |
| 6 | chromium-local-provider | Logout from current user session |  | ❌ |
| 7 | chromium-local-provider | displays published design card correctly |  | ❌ |
| 8 | chromium-local-provider | displays public design card correctly |  | ➖ |
| 9 | chromium-local-provider | Verify Meshery Docker Extension Details |  | ❌ |
| 10 | chromium-local-provider | Verify Meshery Design Embed Details |  | ❌ |
| 11 | chromium-local-provider | should identify relationships for All Relationships | relationship | ❌ |
| 12 | chromium-local-provider | Configure Existing Istio adapter through Mesh Adapter URL from Management page | unstable | ⚠️ |
| 13 | chromium-local-provider | Ping Istio Adapter | unstable | ⚠️ |
| 14 | chromium-local-provider | Import a Model via Url Import |  | ❌ |
| 15 | chromium-local-provider | Import a Model via CSV Import |  | ➖ |
| 16 | chromium-local-provider | Connect to Meshery Istio Adapter and configure it |  | ➖ |
| 17 | chromium-local-provider | imports design via File |  | ❌ |
| 18 | chromium-local-provider | Verify Meshery Catalog Section Details |  | ❌ |
| 19 | chromium-local-provider | Verify Meshery Adapter for Istio Section |  | ❌ |
| 20 | chromium-local-provider | Add performance profile with load generator fortio |  | ❌ |
| 21 | chromium-local-provider | Edit the configuration of a performance profile with load generator fortio and service mesh None |  | ❌ |

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