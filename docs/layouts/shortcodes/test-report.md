### END-TO-END TESTS

- Testing started at: April 13th 2026, 6:17:29 am

**📦 Test Result Summary**

- ✅ 0 passed
- ❌ 2 failed
- ⚠️ 0 flaked
- ⏩ 0 skipped

⌛ _Duration: 0 minutes and 2 seconds_

**Overall Result**: 👎 Some tests failed.



<details>
    <summary>[Show/Hide] Test Result Details</summary>
    <div markdown="1">

| Test | Provider | Browser | Test Case | Tags | Result |
| :---: | :---: | :---: | :--- | :---: | :---: |
| 1 | None | local-setup | authenticate as None provider |  | ❌ |
| 2 | None | remote-setup | authenticate with Remote Provider |  | ❌ |

</div>
</details>




**❌ Failed Test Details**

The following tests failed. Click on each to see the failure reason:


<details>
<summary>❌ authenticate with Remote Provider (None - remote-setup)</summary>

**File Location:** `/home/runner/work/meshery/meshery/ui/tests/e2e/pages/LoginPage.js`

**Error Message:**
```
Error: page.goto: net::ERR_CONNECTION_REFUSED at http://localhost:9081/
Call log:
[2m  - navigating to "http://localhost:9081/", waiting until "load"[22m

```

**Code Snippet:**
```
   at pages/LoginPage.js:44

[0m [90m 42 |[39m
 [90m 43 |[39m     [90m// Navigate to the baseURL after setting cookies[39m
[31m[1m>[22m[39m[90m 44 |[39m     [36mawait[39m [36mthis[39m[33m.[39mpage[33m.[39mgoto(baseURL)[33m;[39m
 [90m    |[39m                     [31m[1m^[22m[39m
 [90m 45 |[39m   }
 [90m 46 |[39m
 [90m 47 |[39m   [36masync[39m loginWithEmail(email[33m,[39m password) {[0m
```

</details>

<details>
<summary>❌ authenticate as None provider (None - local-setup)</summary>

**File Location:** `/home/runner/work/meshery/meshery/ui/tests/e2e/pages/ProviderSelectionPage.js`

**Error Message:**
```
Error: page.goto: net::ERR_CONNECTION_REFUSED at http://localhost:9081/provider
Call log:
[2m  - navigating to "http://localhost:9081/provider", waiting until "load"[22m

```

**Code Snippet:**
```
   at pages/ProviderSelectionPage.js:12

[0m [90m 10 |[39m
 [90m 11 |[39m   [36masync[39m navigateToProviderSelection() {
[31m[1m>[22m[39m[90m 12 |[39m     [36mawait[39m [36mthis[39m[33m.[39mpage[33m.[39mgoto([32m'/provider'[39m)[33m;[39m
 [90m    |[39m                     [31m[1m^[22m[39m
 [90m 13 |[39m   }
 [90m 14 |[39m
 [90m 15 |[39m   [36masync[39m selectProvider(providerName) {[0m
```

</details>






<!-- To see the full report, please visit our CI/CD pipeline with reporter. -->
