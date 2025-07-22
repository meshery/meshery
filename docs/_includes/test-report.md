### END-TO-END TESTS

- Testing started at: July 22nd 2025, 5:55:38 am

**📦 Test Result Summary**

- ✅ 16 passed
- ❌ 42 failed
- ⚠️ 10 flaked
- ⏩ 15 skipped

⌛ _Duration: 9 minutes and 4 seconds_

**Overall Result**: 👎 Some tests failed.



<details>
    <summary>[Show/Hide] Test Result Details</summary>
    <div markdown="1">

| Test | Browser | Test Case | Tags | Result |
| :---: | :---: | :--- | :---: | :---: |
| 1 | chromium-meshery-provider | Create a Model |  | ⚠️ |
| 2 | chromium-meshery-provider | Search a Model and Export it |  | ❌ |
| 3 | chromium-meshery-provider | Import a Model via File Import |  | ➖ |
| 4 | chromium-meshery-provider | Import a Model via Url Import |  | ➖ |
| 5 | chromium-meshery-provider | Import a Model via CSV Import |  | ➖ |
| 6 | chromium-meshery-provider | Verify that UI components are displayed |  | ⚠️ |
| 7 | chromium-meshery-provider | Transition to disconnected state and then back to connected state |  | ❌ |
| 8 | chromium-meshery-provider | Transition to ignored state and then back to connected state |  | ➖ |
| 9 | chromium-meshery-provider | Transition to not found state and then back to connected state |  | ➖ |
| 10 | chromium-meshery-provider | Delete Kubernetes cluster connections |  | ➖ |
| 11 | chromium-meshery-provider | Configure Existing Istio adapter through Mesh Adapter URL from Management page | unstable | ⚠️ |
| 12 | chromium-meshery-provider | Add performance profile with load generator fortio |  | ❌ |
| 13 | chromium-meshery-provider | Ping Istio Adapter | unstable | ⚠️ |
| 14 | chromium-meshery-provider | View detailed result of a performance profile (Graph Visualiser) with load generator fortio |  | ❌ |
| 15 | chromium-meshery-provider | Aggregation Charts are displayed |  | ❌ |
| 16 | chromium-meshery-provider | Edit the configuration of a performance profile with load generator fortio and service mesh None |  | ❌ |
| 17 | chromium-meshery-provider | Connect to Meshery Istio Adapter and configure it |  | ❌ |
| 18 | chromium-meshery-provider | Compare test of a performance profile with load generator fortio |  | ❌ |
| 19 | chromium-meshery-provider | Toggle &quot;Send Anonymous Usage Statistics&quot; |  | ❌ |
| 20 | chromium-meshery-provider | Delete a performance profile with load generator fortio |  | ❌ |
| 21 | chromium-meshery-provider | Toggle &quot;Send Anonymous Performance Results&quot; |  | ❌ |
| 22 | chromium-local-provider | Verify that UI components are displayed |  | ❌ |
| 23 | chromium-local-provider | Add a cluster connection by uploading kubeconfig file |  | ➖ |
| 24 | chromium-local-provider | Transition to disconnected state and then back to connected state |  | ➖ |
| 25 | chromium-local-provider | Transition to ignored state and then back to connected state |  | ➖ |
| 26 | chromium-local-provider | Transition to not found state and then back to connected state |  | ➖ |
| 27 | chromium-local-provider | Delete Kubernetes cluster connections |  | ➖ |
| 28 | chromium-meshery-provider | Verify Meshery Design Embed Details |  | ❌ |
| 29 | chromium-local-provider | Verify Kanvas Snapshot using data-testid |  | ❌ |
| 30 | chromium-local-provider | Test if Left Navigation Panel is displayed |  | ❌ |
| 31 | chromium-meshery-provider | Verify Meshery Catalog Section Details |  | ❌ |
| 32 | chromium-local-provider | Verify Performance Analysis Details |  | ❌ |
| 33 | chromium-local-provider | Test if Notification button is displayed |  | ❌ |
| 34 | chromium-meshery-provider | Verify Meshery Adapter for Istio Section |  | ❌ |
| 35 | chromium-local-provider | Verify Kanvas Details |  | ❌ |
| 36 | chromium-local-provider | Test if Profile button is displayed |  | ❌ |
| 37 | chromium-local-provider | Logout from current user session |  | ❌ |
| 38 | chromium-local-provider | Create a Model |  | ❌ |
| 39 | chromium-local-provider | Search a Model and Export it |  | ➖ |
| 40 | chromium-local-provider | Import a Model via File Import |  | ➖ |
| 41 | chromium-local-provider | Import a Model via Url Import |  | ➖ |
| 42 | chromium-local-provider | Import a Model via CSV Import |  | ➖ |
| 43 | chromium-local-provider | Verify Meshery Docker Extension Details |  | ❌ |
| 44 | chromium-local-provider | Common UI elements |  | ❌ |
| 45 | chromium-local-provider | Configure Existing Istio adapter through Mesh Adapter URL from Management page | unstable | ⚠️ |
| 46 | chromium-local-provider | Verify Meshery Design Embed Details |  | ❌ |
| 47 | chromium-local-provider | Ping Istio Adapter | unstable | ⚠️ |
| 48 | chromium-local-provider | Verify Meshery Catalog Section Details |  | ❌ |
| 49 | chromium-local-provider | Verify Meshery Adapter for Istio Section |  | ❌ |
| 50 | chromium-local-provider | Add performance profile with load generator fortio |  | ❌ |
| 51 | chromium-local-provider | Aggregation Charts are displayed |  | ❌ |
| 52 | chromium-local-provider | View detailed result of a performance profile (Graph Visualiser) with load generator fortio |  | ❌ |
| 53 | chromium-local-provider | Connect to Meshery Istio Adapter and configure it |  | ❌ |
| 54 | chromium-local-provider | Edit the configuration of a performance profile with load generator fortio and service mesh None |  | ❌ |
| 55 | chromium-local-provider | Toggle &quot;Send Anonymous Usage Statistics&quot; |  | ❌ |
| 56 | chromium-local-provider | Compare test of a performance profile with load generator fortio |  | ❌ |
| 57 | chromium-local-provider | Toggle &quot;Send Anonymous Performance Results&quot; |  | ❌ |
| 58 | chromium-local-provider | Delete a performance profile with load generator fortio |  | ❌ |
| 59 | chromium-meshery-provider | Action buttons on adapters tab |  | ❌ |
| 60 | chromium-local-provider | All settings tabs |  | ❌ |
| 61 | chromium-meshery-provider | Grafana elements on metrics tab |  | ❌ |
| 62 | chromium-local-provider | Action buttons on adapters tab |  | ❌ |
| 63 | chromium-local-provider | Grafana elements on metrics tab |  | ❌ |

</div>
</details>


<!-- To see the full report, please visit our CI/CD pipeline with reporter. -->