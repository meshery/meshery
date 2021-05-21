---
layout: default
title: Service Mesh Pattern Management
permalink: functionality/pattern-management
type: functionality
language: en
list: include
---

## Use Meshery to Deploy a Service Mesh Pattern

You can use Meshery to deploy a service mesh pattern. Service mesh patterns are detailed as a YAML file. The format of the YAML is based on the [Service Mesh Performance](https://smp-spec.io) (SMP) specification.

### Pattern Repository

Patterns are available in the [Service Mesh Patterns](https://github.com/service-mesh-patterns) repository. In this repository, you will find each of Meshery's patterns. Each service mesh pattern carries a unique identifier for reference. The patterns in this repository serve in an educational capacity, facilitating learning, and also serve in an operational capacity, facilitating implementation and validation of your service mesh deployment’s adherence to a pattern.

#### To deploy a pattern, execute this command:

```
$ mesheryctl pattern apply -f <path to the pattern>
Deploying pattern…
Pattern successfully deployed.
```

From here, output and behavior will vary based upon the specific pattern you’re deploying. Should you find that your pattern is not successfully deployed, refer to the pattern troubleshooting guide in the Meshery documentation.

As you step through each pattern, you might choose to modify the pattern to suit your specific environment and workload, exploring in-context of your specific situation. Optionally, you may choose to use a plugin like MeshMap, or another visual service mesh topology tool, to facilitate your comprehension of the patterns and to literally see the patterns in-action.

Take time to explore. Try deploying one service mesh’s sample application onto a different service mesh and compare differences in behavior and each service mesh’s ability to manage it. If using Meshery to do so, execute the following commands as an example to deploy the sample application, BookCatalog, onto Open Service Mesh:

```
$ mesheryctl pattern apply -f book-catalog
Deploying application “BookCatalog”...
Deployed. Endpoint(s) available at: http://localhost:8000/catalog
```

<!-- ## Related Reading -->
