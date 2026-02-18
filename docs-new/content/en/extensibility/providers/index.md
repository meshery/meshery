---
title: "Extensibility: Providers"
description: Meshery uses providers to enrich the level of depth by which it manages cloud native infrastructure.
---

Meshery offers _Providers_ as a point of extensibility. It has a built-in Local Provider, named **"None"** and supports _Remote Providers_ that are designed to be pluggable. _Remote Providers_ offer a point of extension to users / integrators to deliver enhanced functionality such as authentication and authorization, using Meshery as a platform.

1. **Extensibility points offer clean separation of Meshery's core functionality versus plugin functionality.**
   - See a list of [Meshery's extensions](https://meshery.io/extensions).
1. **Remote Providers should be able to offer custom RBAC, custom UI components, and custom backend components**
   - Dynamically loadable frameworks need to be identified or created to serve each of these purposes.

### What functionality do Providers offer?

The functionalities offerred by a _Remote Provider_ may vary broadly between providers. Meshery offers extension points that Remote Providers are able to use to inject various functionalities, functionalities that are specific to that provider.
Some examples include:

- **Authentication and Authorization**
  - Examples: session management, two factor authentication, LDAP integration.
- **Long-Term Persistence**
  - Examples: Storage and retrieval of performance test results.
  - Examples: Storage and retrieval of user preferences.
- **Enhanced Visualization**
  - Examples: Creation of a visual cloud native infrastructure topology.
  - Examples: Different charts (metrics), debug (log viewer), distributed trace explorers.
- **Reporting**
  - Examples: Using Meshery's GraphQL server to compose new dashboards.
- **Server Event Management**
  - Examples: Local event storage in database
  - Examples: Remote event synchronization (Remote providers only)

<a href="./images/providers.png">
<img src="./images/providers.png" width="50%" /></a>

  <figcaption>Figure: Selecting a provider in Meshery</figcaption>

## Types of providers

There are two types of providers defined in Meshery, `local` and `remote`.

- The **Local** provider is built-into Meshery.
- **Remote providers** can be implemented by anyone or an organization that wishes to integrate with Meshery. Any number of Remote providers may be available in your Meshery deployment.

### Remote Providers

The use of a Remote Provider, puts Meshery into multi-user mode and requires user authentication. This provides security for the public-facing Meshery UI as the remote provider enforces identity with authentication and authorization. You should also use a remote provider when your use of Meshery is ongoing or used in a team environment (used by multiple people). This can be seen when using Meshery Playground, where a user is prompted to login through the _Layer5 Cloud_ remote provider. Visit [Meshery Playground](https://playground.meshery.io/) to experience this.

A specific provider can be enforced in a Meshery instance by passing the name of the provider with the env variable `PROVIDER`. This applies to both remote and local providers.

Name: **"Layer5"** (default)

- Enforces user authentication.
- Long-term term persistence.
- Save environment setup.
- Retrieve performance test results.
- Events are stored locally and can be published to remote provider. [Read more about server events](/project/contributing/contributing-server-events)

#### Example remote provider

Name: **"Acme"** (default)

- Enforces user authentication.
- Long-term term persistence.
- Save environment setup.
- Retrieve performance test results.
- Events are stored locally and can be published to remote provider. [Read more about server events](/project/contributing/contributing-server-events)
- Free to use.

### Local Provider

The use of the Local Provider, **"None"**, puts Meshery into a single-user mode and does not require authentication. Use the Local provider when your use of Meshery is intended to be shortlived.

Name: **"None"**

- No user authentication.
- Immediate login - users are redirected directly to the dashboard without authentication prompts.
- Container-local storage of test results. Ephemeral.
- Environment setup not saved.
- No performance test result history.
- Server events are stored locally in database. [Read more about server events](/project/contributing/contributing-server-events)
- Free to use.

#### Login Behavior

When the Local Provider is selected, users are immediately redirected to the Meshery dashboard without any authentication challenges. The login flow for **"None"** provider:

1. User selects **"None"** provider from the provider selection UI (or it's enforced via `PROVIDER` environment variable)
2. The provider cookie is set and user is redirected to `/user/login`
3. `InitiateLogin` immediately redirects to `/` (dashboard) or to the originally requested page if a deep-link was preserved
4. User begins working with Meshery without any authentication barriers

This streamlined flow makes the Local Provider ideal for quick evaluations, demos, and single-user scenarios where authentication overhead is unnecessary.

## Provider Login Flow

Understanding how provider selection leads to authentication and dashboard access is important for both users and integrators.

### Interactive (Browser) Flow

1. **Provider Selection**: User visits Meshery and is presented with the provider selection UI at `/provider` (unless a provider is enforced via the `PROVIDER` environment variable)
2. **Provider Activation**: User selects a provider → Meshery sets a `meshery-provider` cookie and redirects to `/user/login?provider=<name>`
3. **Login Initiation**: The `/user/login` route resolves the provider and calls the provider's `InitiateLogin` method:
   - **For "None" (Local)**: Immediately redirects to `/` (dashboard) or to the deep-link target if one was preserved
   - **For Remote Providers**: Redirects to the remote provider's OAuth login page (e.g., GitHub, Google)
4. **Post-Authentication**: After successful authentication, user is redirected to the dashboard or originally requested page

### Enforced Provider Flow

When the `PROVIDER` environment variable is set (e.g., `PROVIDER=None` or `PROVIDER=Meshery`):

1. Provider selection UI is bypassed
2. The specified provider is automatically activated and cookie is set
3. User is redirected directly to `/user/login`
4. Login flow proceeds as described above based on provider type

### Deep-Link Preservation

Meshery preserves the originally requested URL when authentication is required, enabling seamless navigation after login:

- When an unauthenticated user attempts to access a protected page, the URL is base64-encoded into a `ref` query parameter
- After successful login, the provider decodes `ref` and redirects to the original destination
- The Local Provider supports this functionality, validating `ref` values to prevent open redirects (absolute URLs, protocol-relative URLs, and URLs with schemes/hosts are rejected)
- If `ref` validation fails or is absent, users are redirected to the dashboard (`/`)

**Example**: User visits `https://meshery.example.com/extension/meshmap` while unauthenticated → redirected to login with `ref` parameter → after login, automatically returned to MeshMap extension.

{{< alert type="info" title="Deep-Link Security" >}}
Deep-link targets are validated to prevent open redirect vulnerabilities. Only relative paths within the Meshery application are accepted.
{{< /alert >}}

## Runtime Configuration Options

Meshery provides runtime configuration options to control provider behavior:

### PROVIDER

This environment variable enforces a specific provider, bypassing the provider selection UI. This is useful for:
- Dedicated deployments where only one provider should be available
- Automated environments and CI/CD pipelines
- Simplified user experience when provider choice is predetermined

Accepted values:
- `None` - Enforces the Local Provider (no authentication)
- `Meshery` - Enforces the Meshery (Layer5 Cloud) Remote Provider
- Any other registered remote provider name

Example: `PROVIDER=None`

When set, users are automatically directed to the specified provider's login flow upon accessing Meshery.

### PROVIDER_CAPABILITIES_FILEPATH

This environment variable allows you to specify a local file path to load provider capabilities from a static JSON file instead of fetching them from the remote provider's `/capabilities` endpoint. This is useful for:
- Offline development and testing
- Environments with restricted network access
- Ensuring consistent provider capabilities across deployments

Example: `PROVIDER_CAPABILITIES_FILEPATH=/path/to/capabilities.json`

### SKIP_DOWNLOAD_EXTENSIONS

This boolean environment variable controls whether Meshery downloads and refreshes provider extension packages. When set to `true`, Meshery will skip downloading extension packages even when new versions are available. This is particularly useful for:
- Development environments where you want to use locally modified extensions
- Deployments where extensions are pre-packaged or managed separately
- Reducing startup time and network bandwidth usage
- Preventing automatic updates to extension packages

Default: `false` (extensions are downloaded/refreshed)

Example: `SKIP_DOWNLOAD_EXTENSIONS=true`

**Note:** Extension downloads occur during:
- User login (via the TokenHandler)
- Provider capability refresh operations
- Release channel updates

When `SKIP_DOWNLOAD_EXTENSIONS` is enabled, existing extension packages will still be loaded if present, but no new versions will be retrieved.

## Design Principles: Meshery Remote Provider Framework

Meshery's Remote Provider extensibility framework is designed to enable the following functionalities:

1. **Pluggable UI Functionality**
   - Out-of-tree custom UI components with seamless user experience.
   - A system of remote retrieval of extension packages (ReactJS components and Golang binaries).

1. **Pluggable Backend Functionality**
   - Remote Providers have any number of capabilities unbeknownst to Meshery.

1. **Pluggable AuthZ**
   - Design an extensible role based access control system such that Remote Providers can determine their own set of controls. Remote Providers to return JWTs with custom roles, permission keys and permission keychains.

## Building a Provider

Meshery interfaces with providers through a Go interface. The Provider implementations have to be placed in the code and compiled together today. A Provider instance will have to be injected into Meshery when the program starts.

Meshery keeps the implementation of Remote Providers separate so that they are brought in through a separate process and injected into Meshery at runtime (OR) change the way the code works to make the Providers invoke Meshery.

### Verifying Compatibility With Golang Version Update

When Meshery is updated to a newer version of Golang, extension providers need to ensure their integrations remain compatible with the updated version. Changes in the Golang version can lead to compatibility issues, so it’s important to update your extension to align with Meshery’s new environment. For a detailed guide on how to verify and address any compatibility issues, refer to this [guide on verifying compatibility](./verify-compatibility).

### Remote Provider Extension Points

Interwoven into Meshery’s web-based user interface are a variety of extension points. Each extension point is carefully carved out to afford a seamless user experience. Each extension point is identified by a name and type. The following Meshery UI extension points are available:

|**Name**   |**Type**   | **Description**  |
|---|---|---|
|navigator |Menu Items   |This is supposed to be a full page extension which will get a dedicated endpoint in the meshery UI. And will be listed in the meshery UI’s navigator/sidebar. Menu items may refer to full page extensions.   |
|user_prefs   |Single Component   |This is supposed to be remote react components which will get placed in a pre-existing page and will not have a dedicated endpoint. As of now, the only place where this extension can be loaded is the “User Preference” section under meshery settings.   |
|account   |Full Page   |Remote Reactjs components (or other) are placed in a pre-existing page and will have dedicated endpoint: `/extension/account`.   |
|collaborator |Single Component|This is supposed to be remote react components which will get placed in a pre-existing page and will not have a dedicated endpoint. Currently, placed at the Header component of Mehery UI. Its work is to show active Meshery users under the same remote provider.|
  |/extension/\<your name here>|Full Page|The Provider package is unzipped into Meshery server filesystem under `/app/provider-pkg/<package-name>`.|

Remote Providers must fulfill the following endpoints:

1. `/login` - return valid token
1. `/logout` - invalidating token
1. `/capabilities` - return capabilities.json

## UI Extension Points

All UI extensions will be hosted under the endpoint `<mesheryserver:port/provider>`.

### UserPrefs

The UserPrefs extension point expects and loads a component to be displayed into `/userpreferences` page.

### Navigator

The Navigator extension point loads a set of menu items to be displayed in the menu bar on the left hand side of the Meshery UI.

## Provider-Defined Redirects

Remote Providers can define custom URL redirects through the `redirects` field in the capabilities response. This allows providers to dynamically control how certain URL paths are handled by the Meshery server, enabling seamless navigation experiences for users.

### How Redirects Work

When Meshery server receives a request for a specific URL path, it checks if the connected Remote Provider has defined a redirect for that path in its capabilities. If a matching redirect is found, the server issues an HTTP permanent redirect (308) to the target URL.

### Configuring Redirects

Redirects are defined as a key-value map in the provider's capabilities response, where:
- **Key**: The source URL path to intercept (e.g., `/`)
- **Value**: The destination URL to redirect to (e.g., `/extension/meshmap`)

Example configuration in the capabilities response:

{{< code code=`"redirects": {
  "/": "/extension/meshmap",
  "/dashboard": "/extension/meshmap/dashboard"
}` >}}

### Use Cases

1. **Default Landing Page**: Redirect users from the root path (`/`) to a custom extension page.
2. **Custom Navigation**: Guide users to specific extension pages based on their workflow.
3. **URL Aliasing**: Create shorter or more memorable URLs that redirect to extension endpoints.

{{< alert type="info" title="Note" >}}
Redirects are evaluated before serving UI content. Only exact path matches are redirected.
{{< /alert >}}

## Capabilities Endpoint Example

Meshery Server will proxy all requests to remote provider endpoints. Endpoints are dynamically determined and identified in the _**capabilities**_ section of the `/capabilities` endpoint. Providers as an object have the following attributes (this must be returned as a response to `/capabilities` endpoint):
<details>
<summary>Capabilities Endpoint Example</summary>
{{< code code=`{
  "provider_type": "remote",
  "package_version": "v0.1.0",
  "package_url": "https://layer5labs.github.io/meshery-extensions-packages/provider.tar.gz",
  "provider_name": "Meshery",
  "provider_description": [
    "Persistent sessions",
    "Save environment setup",
    "Retrieve performance test results",
    "Free use"
  ],
  "extensions": {
    "navigator": [
      {
        "title": "MeshMap",
        "href": {
          "uri": "/meshmap",
          "external": false
        },
        "component": "/provider/navigator/meshmap/index.js",
        "icon": "/provider/navigator/img/meshmap-icon.svg",
        "link:": true,
        "show": true,
        "children": [
          {
            "title": "View: Single Mesh",
            "href": {
              "uri": "/meshmap/mesh/all",
              "external": false
            },
            "component": "/provider/navigator/meshmap/index.js",
            "icon": "/provider/navigator/img/singlemesh-icon.svg",
            "link": false,
            "show": true
          }
        ]
      }
    ],
    "account": [
      {
          "title": "Overview",
          "on_click_callback": 1,
          "href": {
              "uri": "/account/overview",
              "external": false
          },
          "component": "/provider/account/profile/overview.js",
          "link": true,
          "show": true,
          "type": "full_page"
      },
      {
          "title": "Profile",
          "on_click_callback": 1,
          "href": {
              "uri": "/account/profile",
              "external": false
          },
          "component": "/provider/account/profile/profile.js",
          "link": true,
          "show": true,
          "type": "full_page"
      },
      {
          "title": "API Tokens",
          "on_click_callback": 1,
          "href": {
              "uri": "/account/tokens",
              "external": false
          },
          "component": "/provider/account/profile/tokens.js",
          "link": true,
          "show": true,
          "type": "full_page"
      }
    ],
    "user_prefs": [
      {
        "component": "/provider/userprefs/meshmap-preferences.js"
      }
    ],
    "collaborator": [
      {
        "component": "/provider/collaborator/avatar.js"
      }
    ]
  },
  "redirects": {
    "/": "/extension/meshmap"
  },
  "capabilities": [
    {
        "feature": "sync-prefs",
        "endpoint": "/user/preferences"
    },
    {
        "feature": "persist-results",
        "endpoint": "/results"
    },
    {
        "feature": "persist-result",
        "endpoint": "/result"
    },
    {
        "feature": "persist-metrics",
        "endpoint": "/result/metrics"
    },
    {
        "feature": "persist-smp-test-profile",
        "endpoint": "/user/test-config"
    },
    {
        "feature": "persist-performance-profiles",
        "endpoint": "/user/performance/profiles"
    },
    {
        "feature": "persist-schedules",
        "endpoint": "/user/schedules"
    },
    {
        "feature": "persist-meshery-patterns",
        "endpoint": "/patterns"
    },
    {
        "feature": "persist-meshery-filters",
        "endpoint": "/filters"
    },
    {
        "feature": "persist-meshery-applications",
        "endpoint": "/applications"
    },
    {
        "feature": "persist-meshery-pattern-resources",
        "endpoint": "/patterns/resource"
    },
    {
        "feature": "meshery-patterns-catalog",
        "endpoint": "/patterns/catalog"
    },
    {
        "feature": "meshery-filters-catalog",
        "endpoint": "/filters/catalog"
    },
    {
        "feature": "clone-meshery-patterns",
        "endpoint": "/patterns/clone"
    },
    {
        "feature": "clone-meshery-filters",
        "endpoint": "/filters/clone"
    },
    {
        "feature": "share-designs",
        "endpoint": "/api/content/design/share"
    },
    {
        "feature": "persist-connection",
        "endpoint": "/api/connection"
    }
  ]
}` >}}
</details>

#### Meshery Server Registration

Every Meshery server is capable of registering itself with the remote provider, considering that remote provider supports this feature as a capability.
On successful authentication with the remote provider, Meshery server registers itself by sending a POST request to the remote provider through the `persist-connection` capability. The body of the request should include information so as to uniquely indentify Meshery server and its status.

Example of the request body:

{{< code code=`{
  "server_id": "xxxx-xxxxx-xxxx-xxxx",
  "server_version": "vx.x.x",
  "server_build-sha": "xxxx-xxxxx",
  "server_location": "<protocol>://<hostname>:<port>"
}` >}}

## Configurable OAuth Callback URL

As of the release `v0.5.39`, Meshery Server has extended its capability to configure the callback URL for connected _Remote Providers_. This is helpful when you are deploying Meshery behind an ingress gateway or reverse proxy. You can specify a custom, redirect endpoint for the connected Remote Provider.
Let's understand this through an example deployment scenario below.

### Example Deployment Scenario: Meshery Cloud, Istio, and Azure Kubernetes Service

User has deployed the Meshery behind a Istio Ingress Gateway and the Istio is also behind an Application Gateway (e.g. AKS Application Gateway). Generally, when you use a GitHub Remote Provider for authentication, it redirect the request to the Istio Ingress Gateway FQDN. In this setup, redirection won't be successful because the Ingress Gateway is behind an additional Application Gateway. In this case, you have to define where the request should be redirected once it is authenticated from GitHub.

**_SOLUTION_**

You can define a custom callback URL by setting up the `MESHERY_SERVER_CALLBACK_URL` environment variable while installing Meshery.

**Example using Helm**

First, construct the URL in the format
```https://[CUSTOM_URL]/api/user/token```

- Where `[CUSTOM_URL]` is the URL where **Meshery server** will be or is installed. For example, `https://k8s-staging.test.io/`.
- And `api/user/token` is the **Auth Endpoint** and is append at the end of your custom URL.

So, the final URL would look similar to
```https://k8s-staging.test.io/api/user/token```

{{< alert type="info" title="Note" >}}
Callback URL is not the same as Provider URL. In scenarios where Meshery server and Provider is installed on same server, pay attention to paths or subdomains.
{{< /alert >}}

Next, set the `MESHERY_SERVER_CALLBACK_URL` variable when running the `helm install`. Below is an example:

{{< code code="helm install meshery meshery/meshery --namespace meshery --set env.MESHERY_SERVER_CALLBACK_URL=https://k8s-staging.test.io/api/user/token" >}}

**NOTE**

If Meshery server is accessible in a path of the URL such as `https://k8s-staging.test.io/meshery`, then the callback URL will be `https://k8s-staging.test.io/meshery/api/user/token`.

## Managing your Remote Provider Extension Code

Remote Provider extensions are kept out-of-tree from Meshery (server and UI). You might need to build your extensions under the same environment and set of dependencies as Meshery. The Meshery framework of extensibility has been designed such that in-tree extensions can be safely avoided while still providing a robust platform from which to extend Meshery’s functionality. Often, herein lies the delineation of open vs. closed functionality within Meshery. Remote Providers can bring (plugin) what functionality that they want behind this extensible interface (more about Meshery extensibility), at least that is up to the point that Meshery has provided a way to plug that feature in.

Offering out-of-tree support for Meshery extensions means that:

1. source code to your Meshery extensions are not required to be open source,
1. liability to Meshery’s stability is significantly reduced, avoiding potential bugs in extended components.

Through clearly defined extension points, Meshery extensions may be offered as closed source capabilities that plug into open source Meshery code. To facilitate integration of your Meshery extensions, you might automate the building and releasing of your separate, but interdependent code repositories. You will be responsible for sustaining both your ReactJS and Golang-based extensions.
