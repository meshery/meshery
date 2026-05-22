---
title: "Extensibility: Authorization"
description: Meshery architecture is extensible. Meshery provides several extension points for working with different cloud native projects via authorization, adapters, load generators and providers.
aliases:
- /extensibility/authorization
---

Meshery features an extensible authorization system that offers the ability to deliver fine-grained access control across it's web-based user interface, [Meshery UI]({{< ref "concepts/architecture/_index.md" >}}).

## Authorization Keys

The extensible authorization system consists of a large set of keys. Each key uniquely represents a specific capability, for example, the ability to view a [Connection]({{< ref "concepts/logical/connections/index.md" >}}), edit or delete a Connection. With the help of these keys, the system evaluates the permissions during runtime and renders UI both helping offer a secure management system and a customizable user experience.

{{% alert color="info" title="Note" %}}
The extensible authorization system is available to both Local and Remote Providers. Depending upon your chosen [Remote Provider]({{< ref "reference/extensibility/providers/index.md" >}}), keys, clustering of them, assigning them to user groups, not just individual users or to user roles may be offered.
{{% /alert %}}

### Adding a New Permission Key

Permission keys are the atomic units of this authorization system. Each key is a UUID that maps to a specific capability—for example, creating a workspace, deleting a connection, or evaluating a design. Keys are authored in [`meshery/schemas`](https://github.com/meshery/schemas) and then wired into Meshery UI and the Local Provider database through a **two-repository workflow**.

#### How Permission Checks Work at Runtime

Meshery UI does **not** import the generated `typescript/permissions.ts` from schemas. Instead:

1. The provider returns the user's assigned keys (each with a UUID `id` and a `function` name).
2. [`ui/rtk-query/ability.tsx`](https://github.com/meshery/meshery/blob/master/ui/rtk-query/ability.tsx) maps those keys into CASL abilities.
3. Components call `CAN(keys.<NAME>.action, keys.<NAME>.subject)` using entries from the manually maintained [`ui/utils/permission_constants.ts`](https://github.com/meshery/meshery/blob/master/ui/utils/permission_constants.ts).

Meshery Server does not call a `HasPermission` helper today. UI authorization is enforced client-side via CASL; Remote Providers enforce authorization on their own backend using the generated Go constants from schemas.

#### Prerequisites

Clone both repositories locally and install Node.js and Go:

{{< code code=`git clone https://github.com/meshery/schemas.git
git clone https://github.com/meshery/meshery.git
cd schemas && make setup` >}}

#### Step-by-Step Guide

#### Step 1: Generate a UUID

Each permission key requires a globally unique UUID v4:

{{< code code=`uuidgen | tr '[:upper:]' '[:lower:]'` >}}

#### Step 2: Add a Row to `permissions.csv`

Open [`build/permissions.csv`](https://github.com/meshery/schemas/blob/master/build/permissions.csv) and **copy an existing row in the same category**, then modify it. The spreadsheet has many columns; do not create a sparse row with only the fields below.

| Column | Description | Example |
|--------|-------------|---------|
| **Theme** | High-level category (also used as keychain name). | `Workspace Management` |
| **Category** | Specific feature category. | `Workspace` |
| **Function** | Permission name shown in UI and provider responses. | `Create Workspace` |
| **Feature** | Human-readable description of the capability. | `Create new workspace` |
| **Key ID** | The UUID from Step 1. | `eb42ac41-a883-465e-843c-d64e962a3a0e` |
| **Local Provider** | Set to `TRUE` if the key should be seeded for Local Provider development. | `TRUE` |
| **Inserted** | Leave as `FALSE` for new keys. | `FALSE` |

Place the new row near logically related entries under the same Theme/Category.

{{% alert color="info" title="Local Provider Requirement" %}}
If `Local Provider` is `FALSE`, the key will not appear in [`server/permissions/keys.csv`](https://github.com/meshery/meshery/blob/master/server/permissions/keys.csv) seed data and cannot be tested against the Local Provider until that column is set to `TRUE` and the CSV sync completes.
{{% /alert %}}

#### Step 3: Run the Generators

From the `meshery/schemas` root:

{{< code code=`make generate-permissions` >}}

To generate only one at a time:

{{< code code=`make generate-permissions-go` >}}

{{< code code=`make generate-permissions-ts` >}}

This regenerates:

- [`models/permissions/permissions.go`](https://github.com/meshery/schemas/blob/master/models/permissions/permissions.go) — Go constants for Remote Providers and other downstream Go consumers
- [`typescript/permissions.ts`](https://github.com/meshery/schemas/blob/master/typescript/permissions.ts) — TypeScript constants for downstream consumers

{{% alert color="warning" title="Do Not Edit Generated Files" %}}
Never edit `permissions.go` or `permissions.ts` directly. All changes must go through `permissions.csv` and the generators.
{{% /alert %}}

{{% alert color="info" title="Meshery UI Uses a Separate Map" %}}
Meshery UI reads from [`ui/utils/permission_constants.ts`](https://github.com/meshery/meshery/blob/master/ui/utils/permission_constants.ts), not from `typescript/permissions.ts`. You must add the key there manually in Step 6.
{{% /alert %}}

Go constant names combine **Theme + Function** in PascalCase with spaces removed. For example, `Workspace Management` + `Create Workspace` → `WorkspaceManagementCreateWorkspace`.

#### Step 4: Open a Pull Request in `meshery/schemas`

Commit the CSV and generated files, then validate before pushing:

{{< code code=`make validate-schemas && make consumer-audit` >}}

Open a PR against [`meshery/schemas`](https://github.com/meshery/schemas). After it merges and a new schemas release is tagged, proceed to the Meshery repository changes.

#### Step 5: Wait for `keys.csv` Sync in `meshery/meshery`

[`server/permissions/keys.csv`](https://github.com/meshery/meshery/blob/master/server/permissions/keys.csv) is updated automatically by the [`Import Keys`](https://github.com/meshery/meshery/blob/master/.github/workflows/generate_keys.yml) workflow, which pulls from the permissions spreadsheet daily. Keys with `Local Provider = TRUE` in the spreadsheet appear here and are seeded into the Local Provider database on server startup via [`SeedKeys`](https://github.com/meshery/meshery/blob/master/server/models/keys_helper.go).

If you need the key immediately for local testing, coordinate with maintainers to trigger the workflow or temporarily add the row to `server/permissions/keys.csv` in your PR (following the column layout of existing rows).

#### Step 6: Add the Key to Meshery UI

Add an entry to the `keys` object in [`ui/utils/permission_constants.ts`](https://github.com/meshery/meshery/blob/master/ui/utils/permission_constants.ts):

{{< code code=`<CONSTANT_NAME>: {
  subject: '<Function>',
  action: '<Generated UUID>',
},` >}}

Example for the existing `Create Workspace` permission:

{{< code code=`CREATE_WORKSPACE: {
  subject: 'Create Workspace',
  action: 'eb42ac41-a883-465e-843c-d64e962a3a0e',
},` >}}

*   `<CONSTANT_NAME>`: A `SCREAMING_SNAKE_CASE` label used in UI code (for example, `CREATE_WORKSPACE`).
*   `subject`: Must **exactly** match the `Function` column in `permissions.csv` (including capitalization and punctuation).
*   `action`: The UUID from Step 1.

Also mirror the same entry in [`docs/static/js/permission_constants.js`](https://github.com/meshery/meshery/blob/master/docs/static/js/permission_constants.js) so the [Permission Keys Reference](/reference/permissions) page displays the correct Moniker column.

#### Step 7: Gate UI with `CAN`

Import `CAN` and `keys`, then conditionally render or disable the control:

{{< code code=`import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';

const canCreate = CAN(keys.CREATE_WORKSPACE.action, keys.CREATE_WORKSPACE.subject);

return (
  <Button disabled={!canCreate} onClick={handleCreate}>
    Create Workspace
  </Button>
);` >}}

This matches production usage in components such as [`ui/components/workspaces/index.tsx`](https://github.com/meshery/meshery/blob/master/ui/components/workspaces/index.tsx).

#### Step 8: Verify End-to-End

From `meshery/schemas`:

{{< code code=`make validate-schemas && make consumer-audit` >}}

From `meshery/meshery`:

{{< code code=`make ui-lint
make docs-build` >}}

Then run Meshery locally, assign a role that includes your new key, and confirm the gated UI element appears only when the user holds the permission. For Local Provider testing, restart Meshery Server (or reset the database) after `keys.csv` contains your key with `Local Provider = TRUE`.

#### Using Generated Go Constants (Remote Providers)

Remote Providers and other Go services import constants from `github.com/meshery/schemas/models/permissions`. The constant for a key is `<Theme><Function>` in PascalCase—for example, `permissions.WorkspaceManagementCreateWorkspace`.

Meshery Server itself does not import this package for request-handler checks today. If you are building authorization logic in a Remote Provider, consume the generated constant UUID there using your provider's own enforcement mechanism.

#### Using Permission Keys in Meshery UI

Import the `CAN` utility and the `keys` map:

{{< code code=`import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';

const hasAccess = CAN(keys.<CONSTANT_NAME>.action, keys.<CONSTANT_NAME>.subject);
return hasAccess ? <YourComponent /> : <FallbackComponent />;` >}}

Example with a key that exists today:

{{< code code=`import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';

const canCreate = CAN(keys.CREATE_WORKSPACE.action, keys.CREATE_WORKSPACE.subject);
return canCreate ? <CreateWorkspaceButton /> : null;` >}}

`CAN` compares the UUID (`action`) and lowercased function name (`subject`) against the abilities loaded from the provider session.

## Authorization Framework

Meshery utilizes CASL (JS-based permission framework) to evaluate any given user's set of session keys against the built-in keyhooks populated through each invidual Meshery UI page. This allows for granular control over the UI, empowering you to tailor your Meshery experience to your organization's needs by limiting access to specific features and functionalities based on the user's assigned keys.

<a href="./images/permission-in-UI.png">
  <img style="width:min(100%,800px)" src="./images/permission-in-UI.png" />
</a>

### Introduction to CASL.js

[CASL.js](https://casl.js.org) is an isomorphic authorization JavaScript library which restricts what resources a given client is allowed to access. It's designed to be incrementally adoptable and can easily scale between a simple claim based and fully featured subject and attribute based authorization. It makes it easy to manage and share permissions/keys across UI components, API services, and database queries.

An example of how CASL evaluate permission in UI.
{{< code code=`<React.Fragment>
	{!CAN(keys.DELETE_CONNECTION.action, keys.DELETE_CONNECTION.subject) && (
		<Button id="delete-connection">Delete<Button/>
	)}
</React.Fragment>` >}}

 Once a user has logged in, the backend will send a response with the permissions that the user has. Those permissions will be used to create abilities on the frontend, CASL gets updated with those abilities. The UI maintains a constant file containing all allowed permissions, referred to as keys. With the help of these keys, CAN function evaluates the permissions during runtime and renders UI.

{{% alert color="dark" title="Note" %}}
It's important to understand not all pages uses CASL authorization, means even if you are not assigned with any role within organization you might access preferences page and Meshery UI dashboard.
{{% /alert %}}

## Authorization using Local Provider

Meshery's built-in identity provider, "Local" Provider, operates with a large set of predefined keys interspersed throughout Meshery UI and persisted in [Meshery Database]({{< ref "concepts/architecture/database/index.md" >}}). These keys are used to evaluate the permissions of a given user and render the UI accordingly. The keys are grouped into three categories: `action`, `subject`, and `object`.

{{< discuss >}}
