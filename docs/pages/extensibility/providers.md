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

Meshery interfaces with Providers through a Go interface. The Provider implementations have to be placed in the code and compiled together today. A Provider instance will have to be injected into Meshery when the program starts.

![Providers](/assets/img/providers/provider_screenshot.png)

Eventually, Meshery will keep the implementation of Providers separate so that they are brought in through a separate process and injected into Meshery at runtime (OR) change the way the code works to make the Providers invoke Meshery.

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

### What functionality do Providers perform? 
- Authentication and Authorization
 - Examples: session management, two factor authentication, LDAP integration
- Long-Term Persistence
 - Storage and retrieval of performance test results
 - Storage and retrieval of user preferences

### Types of providers
Two types of providers are defined in Meshery: local and remote.

**Remote Providers**
Name: “Meshery” (default)
- Authentication and Authorization
- Results long term persistence
- Save environment setup
- Retrieve performance test results
- Free to use

**Local Provider**
Name: “None”
- No Authentication
- Local storage of results. Mainly ephemeral.
- Environment setup not saved.
- No performance test result history.
- Free to use.

Meshery provides the ability for you as a service mesh manager to customize your service mesh deployment.