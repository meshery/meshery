---
layout: default
title: "Extensibility: Providers"
permalink: extensibility/providers
type: Reference
#redirect_from: architecture/adapters
abstract: "Meshery offers support for more adapters than any other project or product in the world. Meshery uses adapters for managing the various service meshes."
language: en
list: include
---
Meshery offers Providers as a point of extensibility. 

![Providers](/assets/img/providers/provider_screenshot.png)

### What functionality do Providers perform? 

What a given Remote Provider offers might vary broadly between providers. Meshery offers extension points that Remote Providers are able to use to inject different functionality - functionality specific to that provider.

- Authentication and Authorization
  - Examples: session management, two factor authentication, LDAP integration
- Long-Term Persistence
  - Examples: Storage and retrieval of performance test results
  - Examples: Storage and retrieval of user preferences

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

Providers as an object have the following attributes:

```
{
 id: aower2-1234xd-1234 [guid] 	# future consideration
 type: local, [ local | remote ]
 display-name: "None", [ string ]
 description: "Default Provider
   - feature 1 
   - feature 2", [ multi-line string ] 
 capabilities: [ 			# future consideration
    {featureName:"perf_results", present: true}, 
    {featureName:"two_factor_auth", present: true}
   ]
}
```

Meshery enables you as a service mesh owner to customize your service mesh deployment.