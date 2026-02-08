---
title: Logical
description: Concepts for understanding Meshery's various features and components.
---

As an extensible platform, Meshery empowers you with a wide range of logical constructs that provide support for the majority of the systems in the cloud and cloud native ecosystems. Meshery abstracts away the system specific requirements and help you focus on getting things done.

## Logical Concepts

<ul>
  {{< concepts-list section="concepts" subsection="logical" lang="en" >}}
</ul>

The logical concepts included in Meshery establish a set of foundational constructs. Each logical construct is:

1. Versioned (see [Schemas](https://github.com/meshery/schemas))
2. Extensible (see [Extension Points](/extensibility))
3. Composable (see [Patterns](/concepts/logical/patterns))
4. Portable (see Export/Import of [Designs](/concepts/logical/designs) and [Models](/concepts/logical/models))
5. Interoperable (see [Compatibility Matrix](/installation/compatibility-matrix))
6. Configurable (see [Lifecycle Management](/guides/infrastructure-management/lifecycle-management))
7. Documented (_you are here_)
8. Testable
9. Maintainable
10. Secure (v0.9.0)
11. Observable (v0.1.0)

Every construct is represented in each of the following forms:

- **Schema** (static) - the skeletal structure representing a logical view of the size, shape, characteristics of a construct.
  - *Example: Component schema found in github.com/meshery/schemas*
- **Definition** (static) - An implementation of the Schema containing specific configuration for the construct at-hand.
  - *Example: Component definition generically describing a Kubernetes Pod*
- **Declaration** (static) - A defined construct; A specific deof the Definition.
  - *Example: Component configuration of an NGINX container as a Kubernetes Pod*
- **Instance** (dynamic) - A realized construct (deployed/discovered); An instantiation of the Declaration.
  - *Example: NGINX-as234z2 pod running in cluster*