---
layout: default
title: "Extensibility: Providers"
permalink: extensibility/providers
type: Extensibility
#redirect_from: architecture/adapters
abstract: "Meshery offers support for more adapters than any other project or product in the world. Meshery uses adapters for managing the various service meshes."
language: en
list: include
---

Meshery offers Providers as a point of extensibility. With a built-in Local Provider (named "None"), Meshery Remote Providers are designed to be pluggable. Remote Providers offer points of extension to users / integrators to deliver enhanced functionality, using Meshery as a platform. A specific provider can be enforced in a Meshery instance by passing the name of the provider with the env variable PROVIDER.

1. **Extensibility points offer clean separation of Meshery's core functionality versus plugin functionality.**
   - Meshmap is an example of a feature to be delivered via Remote Provider.
1. **Remote Providers should be able to offer custom RBAC, custom UI components, and custom backend components**
   - Dynamically loadable frameworks need to be identified or created to serve each of these purposes.

### Design Principles: Meshery Remote Provider Framework

Meshery's Remote Provider extensibility framework is designed to enable:

1. **Pluggable UI Functionality:**

   - Out-of-tree custom UI components with seamless user experience.
   - A system of remote retrieval of extension packages (ReactJS components and Golang binaries).

1. **Pluggable Backend Functionality:**

   - Remote Providers have any number of capabilities unbeknownst to Meshery.

1. **Pluggable AuthZ**
   - Design an extensible role based access control system such that Remote Providers can determine their own set of controls. Remote Providers to return JWTs with custom roles, permission keys and permission keychains.

![Providers](/assets/img/providers/provider_screenshot.png)

### What functionality do Providers perform?

What a given Remote Provider offers might vary broadly between providers. Meshery offers extension points that Remote Providers are able to use to inject different functionality - functionality specific to that provider.

- **Authentication and Authorization**
  - Examples: session management, two factor authentication, LDAP integration.
- **Long-Term Persistence**
  - Examples: Storage and retrieval of performance test results.
  - Examples: Storage and retrieval of user preferences.
- **Enhanced Visualization**
  - Examples: Creation of a visual service mesh topology.
  - Examples: Different charts (metrics), debug (log viewer), distributed trace explorers.
- **Reporting**
  - Examples: Using Meshery's GraphQL server to compose new dashboards.

## Types of providers

Two types of providers are defined in Meshery: `local` and `remote`. The Local provider is built-into Meshery. Remote providers are may be implemented by anyone or organization that wishes to integrate with Meshery. Any number of Remote providers may be available in your Meshery deployment.

### Remote Providers

Use of a Remote Provider, puts Meshery into multi-user mode and requires user authentication. Use a Remote provider when your use of Meshery is ongoing or used in a team environment (used by multiple people).

Name: **“Meshery”** (default)

- Enforces user authentication.
- Long-term term persistence of test results.
- Save environment setup.
- Retrieve performance test results.
- Retrieve conformance test results.
- Free to use.

### Local Provider

Use of the Local Provider, "None", puts Meshery into single-user mode and does not require authentication. Use the Local provider when your use of Meshery is intended to be shortlived.

Name: **“None”**

- No user authentication.
- Container-local storage of test results. Ephemeral.
- Environment setup not saved.
- No performance test result history.
- No conformance test result history.
- Free to use.

## Building a Provider

Meshery interfaces with Providers through a Go interface. The Provider implementations have to be placed in the code and compiled together today. A Provider instance will have to be injected into Meshery when the program starts.

Meshery keeps the implementation of Remote Providers separate so that they are brought in through a separate process and injected into Meshery at runtime (OR) change the way the code works to make the Providers invoke Meshery.

### Remote Provider Extension Points

Interwoven into Meshery’s web-based, user interface are a variety of extension points. Each extension point is carefully carved out to afford a seamless user experience. Each extension point is identified by a name and type. The following Meshery UI extension points are available:

- **Name:** "navigator"
  **Type:** Menu Items  
  **Description:** This is supposed to be a full page extension which will get a dedicated endpoint in the meshery UI. And will be listed in the meshery UI’s navigator/sidebar. Menu items may refer to full page extensions.

- **Name:** "user_prefs"
  **Type:** Single Component  
  **Description:** This is supposed to be remote react components which will get placed in a pre-existing page and will not have a dedicated endpoint. As of now, the only place where this extension can be loaded is the “User Preference” section under meshery settings.

- **Name:** "account"
  **Type:** Full Page
  **Description:** Remote Reactjs components (or other) are placed in a pre-existing page and will have dedicated endpoint: `/extension/account`.

- **Name:** "collaborator"
  **Type:** Single Component
  **Description:** This is supposed to be remote react components which will get placed in a pre-existing page and will not have a dedicated endpoint. Currently, placed at the Header component of Mehery UI.
  Its work is to show active Meshery users under the same remote provider.

- **Name:** /extension/\<your name here>  
  **Type:** Full Page  
  **Description:** The Provider package is unzipped into Meshery server filesystem under `/app/provider-pkg/<package-name>`.

Remote Providers must fulfill the following endpoints:

1. `/login` - return valid token
1. `/logout` - invalidating token
1. `/capabilities` - return capabilities.json

## UI Extension Points

All UI extensions will be hosted under the endpoint <mesheryserver:port/provider>

### UserPrefs

The UserPrefs extension point expects and loads a component to be displayed into /userpreferences page.

### Navigator

The Navigator extension point loads a set of menu items to be displayed in the menu bar on the left hand side of the Meshery UI.

## Capabilities Endpoint Example

Meshery Server will proxy all requests to remote provider endpoints. Endpoints are dynamically determined and identified in the "capabilities" section of the `/capabilities` endpoint. Providers as an object have the following attributes (this must be returned as a response to `/capabilities` endpoint):

{% capture code_content %}{
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
        "feature": "persist-smi-results",
        "endpoint": "/smi/results"
    },
    {
        "feature": "persist-smi-result",
        "endpoint": "/smi/result"
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
}{% endcapture %}
{% include code.html code=code_content %}

##### Meshery Server Registration

Every Meshery server is capable of registering itself with the remote provider, considering that remote provider supports this feature as a capability.
On successful authentication with the remote provider, Meshery server registers itself by sending a POST request to the remote provider through the `persist-connection` capability. The body of the request should include information so as to uniquely indentify Meshery server and its status.

Example of the request body:

{% capture code_content %}
  {
    "server_id": "xxxx-xxxxx-xxxx-xxxx",
    "server_version": "vx.x.x",
    "server_build-sha": "xxxx-xxxxx",
    "server_location": "<protocol>://<hostname>:<port>”
  }
{% endcapture %}
{% include code.html code=code_content %}

## Configurable OAuth Callback URL

As of release v0.5.39, Meshery Server has extended its capability to configure the callback URL for connected Remote Providers. Helpful when you are deploying Meshery behind an ingress gateway or reverse proxy. You can specify a custom, redirect endpoints for the connected Remote Provider.

### Example Deployment Scenario: Meshery Cloud, Istio, and Azure Kubernetes Service

User has deployed the Meshery behind a Istio Ingress Gateway and the Istio is also behind an Application Gateway (e.g. AKS Application Gateway). Generally, when you use a GitHub Remote Provider for authentication, it redirect the request to the Istio Ingress Gateway FQDN. In this setup, redirection won't be successful because the Ingress Gateway is behind an additional Application Gateway. In this case, you have to define where the request should be redirected once it is authenticated from GitHub.

**_Solution:_**

You can define your custom callback URL by setting up `MESHERY_SERVER_CALLBACK_URL` env variable while installing meshery in your K8s cluster.

Example:

If you are deploying Meshery using Helm, you can configure the MESHERY_SERVER_CALLBACK_URL as shown in the following example.

- **Custom URL:** `https://k8s-staging.test.io/`
- **Auth Endpoint:** `api/user/token` (append at the end of your custom URL)

{% capture code_content %}helm install meshery meshery/meshery --namespace meshery --set env.MESHERY_SERVER_CALLBACK_URL=https://k8s-staging.test.io/api/user/token{% endcapture %}
{% include code.html code=code_content %}

##### Note

With a path like, `https://k8s-staging.test.io/meshery`, your callback URL will be `https://k8s-staging.test.io/meshery/api/user/token`.

## Managing your Remote Provider Extension Code

Remote Provider extensions are kept out-of-tree from Meshery (server and UI). You might need to build your extensions under the same environment and set of dependencies as Meshery. The Meshery framework of extensibility has been designed such that in-tree extensions can be safely avoided while still providing a robust platform from which to extend Meshery’s functionality. Often, herein lies the delineation of open vs. closed functionality within Meshery. Remote Providers can bring (plugin) what functionality that they want behind this extensible interface (more about Meshery extensibility), at least that is up to the point that Meshery has provided a way to plug that feature in.

Offering out-of-tree support for Meshery extensions means that:

1. source code to your Meshery extensions are not required to be open source,
1. liability to Meshery’s stability is significantly reduced, avoiding potential bugs in extended components.

Through clearly defined extension points, Meshery extensions may be offered as closed source capabilities that plug into open source Meshery code. To facilitate integration of your Meshery extensions, you might automate the building and releasing of your separate, but interdependent code repositories. You will be responsible for sustaining both your ReactJS and Golang-based extensions.

