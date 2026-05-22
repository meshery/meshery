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

Permission keys are the atomic units of this authorization system. Each key is a UUID that maps to a specific action like deploying a design, validating a configuration, or managing a connection. Keys are defined in the [`meshery/schemas`](https://github.com/meshery/schemas) repository and consumed by both Meshery Server and Meshery UI.

#### Prerequisites

You will need access to these repositories [`meshery/schemas`](https://github.com/meshery/schemas) and [`meshery/meshery`](https://github.com/meshery/meshery). Ensure Node.js and Go are installed.

{{< code code=`git clone https://github.com/meshery/schemas.git
cd schemas
make setup` >}}

#### Step-by-Step Guide

##### Step 1: Generate a UUID

Each permission key requires a globally unique UUID v4.

{{< code code=`uuidgen | tr '[:upper:]' '[:lower:]'` >}}

##### Step 2: Add a Row to `permissions.csv`

Open [`build/permissions.csv`](https://github.com/meshery/schemas/blob/master/build/permissions.csv) and add a new row. Use an existing row in the same category as a reference. The key columns are:

| Column | Description | Example |
|--------|-------------|---------|
| **Theme** | The high-level category theme. | `Catalog Management` |
| **Category** | The specific feature category. | `Designs` |
| **Function** | The permission name (in Title Case). | `Evaluate Design` |
| **Feature** | Description of the permission capability. | `Evaluate relationships in a design by applying policies.` |
| **Primary Persona** | The default persona for this permission. | `Operator` |
| **Source** | The license source identifier. | `Open` |
| **Subscription Tier** | The subscription tier required. | `TeamDesigner` |
| **Key ID** | The newly generated UUID v4. | `5ad992cc-2288-42c4-8259-40fba1500224` |
| **Inserted** | Database seeding flag (set to `FALSE`). | `FALSE` |

Place the new row near logically related entries (e.g., next to existing permissions under the same Category).

##### Step 3: Run the Generators

From the `meshery/schemas` root, generate both Go and TypeScript permission key constants:

{{< code code=`make generate-permissions` >}}

To generate only one at a time:

{{< code code=`make generate-permissions-go` >}}

{{< code code=`make generate-permissions-ts` >}}

This regenerates two files from the CSV:
- [`models/permissions/permissions.go`](https://github.com/meshery/schemas/blob/master/models/permissions/permissions.go) — Go constants used by Meshery Server
- [`typescript/permissions.ts`](https://github.com/meshery/schemas/blob/master/typescript/permissions.ts) — TypeScript constants

{{% alert color="warning" title="Do Not Edit Generated Files" %}}
Never edit `permissions.go` or `permissions.ts` directly. All changes must go through `permissions.csv` and the generators.
{{% /alert %}}

##### Step 4: Add the Key to Meshery UI

The Meshery UI maintains a separate, manually-maintained permission mapping in [`ui/utils/permission_constants.ts`](https://github.com/meshery/meshery/blob/master/ui/utils/permission_constants.ts). Add your key to the `keys` object using the template below:

{{< code code=`<CONSTANT_NAME>: {
  subject: '<Function>',
  action: '<Generated UUID>',
},` >}}

For example:

{{< code code=`EVALUATE_DESIGN: {
  subject: 'Evaluate Design',
  action: '5ad992cc-2288-42c4-8259-40fba1500224',
},` >}}

*   `<CONSTANT_NAME>`: Typically a `SCREAMING_SNAKE_CASE` representation of your permission function.
*   `subject`: Must exactly match the value in the `Function` column of the CSV.
*   `action`: The UUID generated in Step 1.

##### Step 5: Verify

Run the schema audit from the `meshery/schemas` checkout to confirm correctness:

{{< code code=`make validate-schemas && make consumer-audit` >}}

#### Using Permission Keys

##### In Go Server Handlers

Import the generated constants from `github.com/meshery/schemas/models/permissions` and check the key before processing a request. The constant name is generated by combining the `Theme` and `Function` columns in `PascalCase` (with spaces removed):

{{< code code=`import (
    "github.com/meshery/schemas/models/permissions"
)

if !provider.HasPermission(r.Context(), permissions.<Theme><Function>) {
    http.Error(w, "Forbidden", http.StatusForbidden)
    return
}` >}}

For example, using the `Catalog Management` theme and `Evaluate Design` function:

{{< code code=`if !provider.HasPermission(r.Context(), permissions.CatalogManagementEvaluateDesign) {
    http.Error(w, "Forbidden", http.StatusForbidden)
    return
}` >}}

##### In React / TypeScript UI

Import the `CAN` utility and the permission `keys` object to evaluate permissions in UI components:

{{< code code=`import CAN from 'utils/can';
import { keys } from 'utils/permission_constants';

const hasAccess = CAN(keys.<CONSTANT_NAME>.action, keys.<CONSTANT_NAME>.subject);
return hasAccess ? <YourComponent /> : <FallbackComponent />;` >}}

For example:

{{< code code=`import CAN from 'utils/can';
import { keys } from 'utils/permission_constants';

const canEvaluate = CAN(keys.EVALUATE_DESIGN.action, keys.EVALUATE_DESIGN.subject);
return canEvaluate ? <EvaluateButton /> : null;` >}}


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

