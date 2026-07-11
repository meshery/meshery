---
title: "Extensibility: Authorization"
description: Meshery architecture is extensible. Meshery provides several extension points for working with different cloud native projects via authorization, adapters, load generators and providers.
aliases:
- /extensibility/authorization
---

Meshery features an extensible authorization system that offers the ability to deliver fine-grained access control across its web-based user interface, [Meshery UI]({{< ref "concepts/architecture/_index.md" >}}).

## Authorization Keys

The extensible authorization system consists of a large set of keys. Each key uniquely represents a specific capability, for example, the ability to view, edit, or delete a [Connection]({{< ref "concepts/logical/connections/index.md" >}}). With the help of these keys, the system evaluates permissions at runtime to render the UI, offering both a secure management system and a customizable user experience.

### Key Naming and Scoping Conventions

Permission keys in `@meshery/schemas` are generated dynamically using the following formula:
`PascalCase(Theme) + PascalCase(Function)`

*   **Theme (Domain Capability)**: Represents the high-level business domain/capability that governs the action (e.g. `Catalog Management`, `Lifecycle Management`).
*   **Function (Operation)**: Represents the specific operation being performed (e.g. `Unpublish Design`, `Evaluate Relationships`).

#### Domain Scoping vs. UI Layout

Keys are categorized by their **domain of authority** rather than the UI component or file path where they are used. 

For example, the component `MesheryPatternCard` (which resides in the designs/patterns directory) uses the key `CatalogManagementUnpublishDesign` because publishing/unpublishing a design is a **Catalog Management** action (making it public or private in the catalog), even though the card itself represents a pattern/design UI element.

{{% alert color="info" title="Note" %}}
The extensible authorization system is available to both Local and Remote Providers. Depending on your chosen [Remote Provider]({{< ref "reference/extensibility/providers/index.md" >}}), you may be offered features such as grouping keys or assigning them to user groups or roles, rather than just individual users.
{{% /alert %}}

### Adding a New Permission Key

Permission keys are defined and managed centrally via an automated Google Spreadsheet workflow. The spreadsheet serves as the authoritative source of truth, while `@meshery/schemas` compiles and publishes the generated Go and TypeScript artifacts consumed by downstream projects:

1. **Spreadsheet Registration**: Keys are registered in the central permissions spreadsheet.
2. **Meshery Schemas Automation**: A GitHub Actions workflow automatically imports the spreadsheet, updates `build/permissions.csv`, and compiles the definitions into Go and TypeScript libraries.
3. **Dynamic Consumption**: Downstream projects (`meshery/meshery` backend and React UI) consume the generated models and constants dynamically, removing the need for duplicate, hardcoded constants.

---

### Step-by-Step Guide

Follow these steps to generate, register, sync, and wire up a new permission key.

#### Phase 1: Define Key in Spreadsheet

##### Step 1: Generate a UUID v4
Generate a unique, lowercase UUID v4 for the new permission key:
{{< code code=`uuidgen | tr '[:upper:]' '[:lower:]'` >}}

##### Step 2: Add to the Permissions Spreadsheet
Add a new row to the authoritative **Permissions Spreadsheet**. Ensure the following columns are set:
*   **Theme**: The high-level category (e.g. `Catalog Management`).
*   **Category**: Specific feature area (e.g. `Designs`).
*   **Function**: Human-readable operation name (e.g. `Evaluate Relationships`). This value becomes `key.function` in the Provider API.
*   **Feature**: Feature description (e.g. `Evaluate relationships inside a design`).
*   **Key ID**: The generated UUID from Step 1.
*   **Local Provider**: Set to `TRUE` if this key should be seeded for the Local Provider database.

##### Step 3: Run the Schema Sync & Generation Workflow
Once added to the spreadsheet, the GitHub Actions workflow [`generate-artifacts-from-schemas.yml`](https://github.com/meshery/schemas/blob/master/.github/workflows/generate-artifacts-from-schemas.yml) in `meshery/schemas` runs automatically on a daily schedule (and can also be triggered manually) to sync the spreadsheet keys to the local [`build/permissions.csv`](https://github.com/meshery/schemas/blob/master/build/permissions.csv).

This workflow automatically executes the generators to produce:
*   Go constants: [`models/permissions/permissions.go`](https://github.com/meshery/schemas/blob/master/models/permissions/permissions.go)
*   TypeScript definitions: [`typescript/permissions.ts`](https://github.com/meshery/schemas/blob/master/typescript/permissions.ts)

These generated files are committed to the master branch and published as part of the `@meshery/schemas` package.

---

#### Phase 2: Backend Sync

##### Step 4: Wait for `keys.csv` Sync in Meshery
The local database seeds are populated via `server/permissions/keys.csv` in the `meshery/meshery` repository. This file is automatically kept in sync with the spreadsheet by the [`Import Keys`](https://github.com/meshery/meshery/blob/master/.github/workflows/generate_keys.yml) workflow, which runs daily. 

On startup, Meshery Server's [`SeedKeys`](https://github.com/meshery/meshery/blob/master/server/models/keys_helper.go) seeds these keys into the database.

---

#### Phase 3: Wire Key in the UI

##### Step 5: Direct Import from Schemas
Because Meshery UI depends on `@meshery/schemas`, you can gate new UI behavior by importing `Keys` directly from `@meshery/schemas/permissions` instead of adding a manual entry to `ui/utils/permission_constants.ts`. Many existing components still use `permission_constants.ts`, so only those legacy call sites require updates until they are migrated.

Simply import the `Keys` object directly from `@meshery/schemas/permissions` in your component:
{{< code code=`import { Keys } from '@meshery/schemas/permissions';` >}}

##### Step 6: Gate the Component using `CAN`
Pass the schema key's `id` (action) and `function` (subject) to the `CAN` utility to check permissions:
{{< code code=`import CAN from '@/utils/can';
import { Keys } from '@meshery/schemas/permissions';

const key = Keys.CatalogManagementEvaluateRelationships;
const canEvaluate = CAN(key.id, key.function);

return (
  <Button disabled={!canEvaluate} onClick={handleEvaluate}>
    Evaluate Relationships
  </Button>
);` >}}

##### Step 7: Verify End-to-End
1. Run `make ui-lint` to verify that there are no formatting or typescript errors.
2. Restart Meshery Server (or reset database settings) to seed the local provider with the new key.
3. Login and verify that your UI component gates correctly.

---

## Reference: End-to-End Key Mapping

This example shows how the **Evaluate Relationships** key is wired across each layer of the application:

1. **Provider API key object (excerpt from API output)**:
{{< code code=`{
  "id": "c7752be7-5c0f-465d-a8ba-5594acd08b93",
  "function": "Evaluate Relationships",
  "category": "Catalog Management",
  "subcategory": "Designs"
}` >}}

2. **TypeScript constant from `@meshery/schemas`**:
{{< code code=`export const Keys = {
  CatalogManagementEvaluateRelationships: {
    id: "c7752be7-5c0f-465d-a8ba-5594acd08b93",
    function: "Evaluate Relationships",
    category: "Catalog Management",
    subcategory: "Designs"
  }
}` >}}

3. **React Component Gating**:
{{< code code=`import CAN from '@/utils/can';
import { Keys } from '@meshery/schemas/permissions';

const key = Keys.CatalogManagementEvaluateRelationships;
const canEvaluate = CAN(key.id, key.function);
return canEvaluate ? <EvaluateRelationshipsButton /> : null;` >}}

---

## Where Permission Keys Are Stored in the Browser

When testing permission keys locally, the main client-side caches are stored under **`sessionStorage`** and **Cookies** (not `localStorage`):

| Location | Key / object | Contents |
|----------|--------------|----------|
| **`Cookies`** | `token` | The session authorization token. Sent automatically with requests to authenticate the user and authorize access to org-specific permission keys. |
| **`Cookies`** | `meshery-provider` | The active provider (e.g., `Local` or `Layer5`). |
| **`sessionStorage`** | `keys` | JSON array of key objects from the provider (`id`, `function`, `category`, …). Written by [`setKeys`](https://github.com/meshery/meshery/blob/master/ui/store/slices/mesheryUi.ts) and read on startup by [`loadAbility`](https://github.com/meshery/meshery/blob/master/ui/pages/_app.tsx). |
| **`sessionStorage`** | `currentOrg` | Selected organization. Keys are fetched per org via `GET /api/identity/orgs/{orgId}/users/keys`. |
| **In-memory (CASL)** | `ability` in [`ui/utils/can.ts`](https://github.com/meshery/meshery/blob/master/ui/utils/can.ts) | Runtime rules: `{ action: key.id, subject: lowerCase(key.function) }`. Updated by [`ability.update(...)`](https://github.com/meshery/meshery/blob/master/ui/rtk-query/ability.tsx). |
| **Redux store** | `state.ui.keys` | Same array as `sessionStorage.keys`. |
| **RTK Query cache** | `getUserKeys` | Cached response for `/api/identity/orgs/{orgId}/users/keys`. |

On login, Meshery either reuses `sessionStorage.keys` or refetches from the provider, then updates CASL. The static map in [`permission_constants.ts`](https://github.com/meshery/meshery/blob/master/ui/utils/permission_constants.ts) is source code—not browser storage—but `CAN(...)` compares those constants against the CASL rules built from the stored provider keys.

Inspect keys in DevTools (**Application → Session Storage**):
{{< code code=`JSON.parse(sessionStorage.getItem('keys'))
JSON.parse(sessionStorage.getItem('currentOrg'))` >}}

You can also check cookies under DevTools (**Application → Storage → Cookies**).

---

## Troubleshooting Permissions

Use this checklist when a gated button does not appear, permissions look stale after a role change, or a newly added key does not work end-to-end.

##### 1. Confirm the provider returned the key
In DevTools **Network**, check the response for:
`GET /api/identity/orgs/{orgId}/users/keys`
Verify your UUID is in the `keys` array. If you're gating via `permission_constants.ts`, `_.lowerCase(function)` from the API should match `_.lowerCase(subject)` passed to `CAN(...)`.

##### 2. Clear stale browser cache
Meshery reuses `sessionStorage.keys` until the org changes or keys are refetched:
{{< code code=`sessionStorage.removeItem('keys');
location.reload();` >}}

##### 3. Verify the UI constant map
In [`permission_constants.ts`](https://github.com/meshery/meshery/blob/master/ui/utils/permission_constants.ts):
*   `action` = key UUID (`id` from the API)
*   `subject` = key `function` from the API (capitalization/camelCase differences are normalized by `_.lowerCase` in `CAN`)

##### 4. Confirm CASL loaded the rules
*   **Redux DevTools**: check `state.ui.keys` for the expected UUID.
*   **Session Storage**: confirm `sessionStorage.keys` is populated.
*   **Loader**: [`LoadSessionGuard`](https://github.com/meshery/meshery/blob/master/ui/rtk-query/ability.tsx) must finish loading before `CAN(...)` returns meaningful results.

##### 5. Local Provider: confirm database seeding
The key must be in [`server/permissions/keys.csv`](https://github.com/meshery/meshery/blob/master/server/permissions/keys.csv) with **`Local Provider = TRUE`**. Restart Meshery Server (or reset the local DB) after the CSV updates.

##### 6. Remote Provider: confirm role assignment
Keys come from roles assigned in the Remote Provider admin UI. An empty API response usually means a role/keychain issue—not a missing `permission_constants.ts` entry alone.

##### 7. Verify browser cookies and session validity
Verify that the `token` cookie is set and not expired:
*   Open DevTools **Application** -> **Storage** -> **Cookies** and select the Meshery site URL.
*   Confirm that the `token` cookie is present. If it is missing or has expired, you will be redirected to the login page or receive a `401 Unauthorized` response when fetching keys.
*   Check that the `meshery-provider` cookie matches the intended provider (e.g., `Local` or `Layer5`).

##### Common symptoms

| Symptom | Likely cause |
|---------|----------------|
| Button never appears | User lacks the key; if gating via `permission_constants.ts`, missing entry; or `CAN(...)` not wired |
| New key not visible after merge | Stale `sessionStorage.keys`; missing Local Provider seed row; or only schemas PR merged |
| Works in one org, not another | Keys are org-scoped—check `currentOrg` and refetch keys |
| API returns `401` or `403` when fetching keys | Expired or missing `token` cookie; verify browser cookie store |

---

## Authorization Framework

Meshery utilizes CASL (JS-based permission framework) to evaluate any given user's set of session keys against the built-in keyhooks populated through each individual Meshery UI page. This allows for granular control over the UI, empowering you to tailor your Meshery experience to your organization's needs by limiting access to specific features and functionalities based on the user's assigned keys.

<a href="./images/permission-in-UI.png">
  <img style="width:min(100%,800px)" src="./images/permission-in-UI.png" />
</a>

### Introduction to CASL.js

[CASL.js](https://casl.js.org) is an isomorphic authorization JavaScript library which restricts what resources a given client is allowed to access. It's designed to be incrementally adoptable and can easily scale between a simple claim based and fully featured subject and attribute based authorization. It makes it easy to manage and share permissions/keys across UI components, API services, and database queries.

An example of how CASL evaluates permissions in the UI:
{{< code code=`<React.Fragment>
	{CAN(keys.DELETE_CONNECTION.action, keys.DELETE_CONNECTION.subject) && (
		<Button id="delete-connection">Delete</Button>
	)}
</React.Fragment>` >}}

Upon user login, the backend returns the list of authorized permissions. These permissions are then used to build and update the CASL ability rules on the frontend. The UI maintains a constant file containing all allowed permissions (referred to as keys). With the help of these keys, the `CAN` function evaluates permissions at runtime and renders the UI accordingly.

{{% alert color="dark" title="Note" %}}
It's important to understand not all pages uses CASL authorization, means even if you are not assigned with any role within organization you might access preferences page and Meshery UI dashboard.
{{% /alert %}}

## Authorization using Local Provider

Meshery's built-in identity provider, "Local" Provider, operates with a large set of predefined keys interspersed throughout Meshery UI and persisted in [Meshery Database]({{< ref "concepts/architecture/database/index.md" >}}). These keys are used to evaluate the permissions of a given user and render the UI accordingly. The keys are grouped into three categories: `action`, `subject`, and `object`.

{{< discuss >}}
