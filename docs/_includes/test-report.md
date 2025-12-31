### END-TO-END TESTS

- Testing started at: December 31st 2025, 7:05:39 am

**📦 Test Result Summary**

- ✅ 7 passed
- ❌ 37 failed
- ⚠️ 0 flaked
- ⏩ 10 skipped

⌛ _Duration: 4 minutes and 34 seconds_

**Overall Result**: 👎 Some tests failed.



<details>
    <summary>[Show/Hide] Test Result Details</summary>
    <div markdown="1">

| Test | Browser | Test Case | Tags | Result |
| :---: | :---: | :--- | :---: | :---: |
| 1 | remote-setup | authenticate as Meshery provider |  | ❌ |
| 2 | chromium-local-provider | Verify that UI components are displayed |  | ❌ |
| 3 | chromium-local-provider | Add a cluster connection by uploading kubeconfig file |  | ➖ |
| 4 | chromium-local-provider | Transition to disconnected state and then back to connected state |  | ➖ |
| 5 | chromium-local-provider | Transition to ignored state and then back to connected state |  | ➖ |
| 6 | chromium-local-provider | Transition to not found state and then back to connected state |  | ➖ |
| 7 | chromium-local-provider | Delete Kubernetes cluster connections |  | ➖ |
| 8 | chromium-local-provider | renders design page UI |  | ❌ |
| 9 | chromium-local-provider | Verify Kanvas Snapshot using data-testid |  | ❌ |
| 10 | chromium-local-provider | Test if Left Navigation Panel is displayed |  | ❌ |
| 11 | chromium-local-provider | Verify Performance Analysis Details |  | ❌ |
| 12 | chromium-local-provider | Test if Notification button is displayed |  | ❌ |
| 13 | chromium-local-provider | Verify Kanvas Details |  | ❌ |
| 14 | chromium-local-provider | should verify Design Configurator page elements |  | ❌ |
| 15 | chromium-local-provider | displays published design card correctly |  | ❌ |
| 16 | chromium-local-provider | displays public design card correctly |  | ❌ |
| 17 | chromium-local-provider | should edit design in Design Configurator |  | ❌ |
| 18 | chromium-local-provider | Logout from current user session |  | ❌ |
| 19 | chromium-local-provider | Create a Model |  | ❌ |
| 20 | chromium-local-provider | Search a Model and Export it |  | ➖ |
| 21 | chromium-local-provider | Import a Model via File Import |  | ➖ |
| 22 | chromium-local-provider | Import a Model via Url Import |  | ➖ |
| 23 | chromium-local-provider | Import a Model via CSV Import |  | ➖ |
| 24 | chromium-local-provider | Common UI elements |  | ❌ |
| 25 | chromium-local-provider | Verify Meshery Docker Extension Details |  | ❌ |
| 26 | chromium-local-provider | Test if Profile button is displayed |  | ❌ |
| 27 | chromium-local-provider | should identify relationships for All Relationships | relationship | ❌ |
| 28 | chromium-local-provider | should identify relationships for Namespace-Namespace-Deny-Relationship-Test | relationship | ❌ |
| 29 | chromium-local-provider | Configure Existing Istio adapter through Mesh Adapter URL from Management page | unstable | ⚠️ |
| 30 | chromium-local-provider | Ping Istio Adapter | unstable | ⚠️ |
| 31 | chromium-local-provider | Add performance profile with load generator fortio |  | ❌ |
| 32 | chromium-local-provider | View detailed result of a performance profile (Graph Visualiser) with load generator fortio |  | ❌ |
| 33 | chromium-local-provider | Edit the configuration of a performance profile with load generator fortio and service mesh None |  | ❌ |
| 34 | chromium-local-provider | Compare test of a performance profile with load generator fortio |  | ❌ |
| 35 | chromium-local-provider | Delete a performance profile with load generator fortio |  | ❌ |
| 36 | chromium-local-provider | imports design via File |  | ❌ |
| 37 | chromium-local-provider | imports design via URL |  | ❌ |
| 38 | chromium-local-provider | deletes a published design from the list |  | ❌ |
| 39 | chromium-local-provider | deploys a published design to a connected cluster |  | ❌ |
| 40 | chromium-local-provider | Verify Meshery Design Embed Details |  | ❌ |
| 41 | chromium-local-provider | All settings tabs |  | ❌ |
| 42 | chromium-local-provider | Aggregation Charts are displayed |  | ❌ |
| 43 | chromium-local-provider | Connect to Meshery Istio Adapter and configure it |  | ➖ |
| 44 | chromium-local-provider | Toggle &quot;Send Anonymous Usage Statistics&quot; |  | ❌ |
| 45 | chromium-local-provider | Toggle &quot;Send Anonymous Performance Results&quot; |  | ❌ |
| 46 | chromium-local-provider | Verify Meshery Catalog Section Details |  | ❌ |
| 47 | chromium-local-provider | Action buttons on adapters tab |  | ❌ |
| 48 | chromium-local-provider | Grafana elements on metrics tab |  | ❌ |
| 49 | chromium-local-provider | Verify Meshery Adapter for Istio Section |  | ❌ |

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