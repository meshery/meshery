# Logical

> Logical Concepts for understanding Meshery's various features and components.

---

LLMS index: [llms.txt](/pr-preview/pr-22676/llms.txt)

---

As an extensible platform, Meshery empowers you with a wide range of logical constructs that provide support for the majority of the systems in the cloud and cloud native ecosystems. Meshery abstracts away the system specific requirements and help you focus on getting things done.

The logical concepts included in Meshery establish a set of foundational constructs. Each logical construct is:

1. Versioned (see [Schemas](https://github.com/meshery/schemas))
2. Extensible (see [Extension Points](/pr-preview/pr-22676/reference/extensibility/))
3. Composable (see [Patterns](/pr-preview/pr-22676/concepts/logical/patterns/))
4. Portable (see Export/Import of [Designs](/pr-preview/pr-22676/concepts/logical/designs/) and [Models](/pr-preview/pr-22676/concepts/logical/models/))
5. Interoperable (see [Compatibility Matrix](/pr-preview/pr-22676/project/compatibility-matrix/compatibility-matrix/))
6. Configurable (see [Lifecycle Management](/pr-preview/pr-22676/guides/infrastructure-management/lifecycle-management/))
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

## Logical Concepts

---

Section pages:

- [Components](/pr-preview/pr-22676/concepts/logical/components/): Meshery Components identify and characterize infrastructure under management.
- [Connections](/pr-preview/pr-22676/concepts/logical/connections/): Meshery Connections are managed and unmanaged resources that either through discovery or manual entry are managed by a state machine and used within one or more Environments.
- [Credentials](/pr-preview/pr-22676/concepts/logical/credentials/): Meshery uses one or more Credentials when authenticating to a managed or unmanaged Connection.
- [Designs](/pr-preview/pr-22676/concepts/logical/designs/): Meshery Designs are descriptive, declarative characterizations of how your Kubernetes infrastructure should be configured.
- [Environments](/pr-preview/pr-22676/concepts/logical/environments/): Environments are how you organize your deployment targets (whether on-premises servers or cloud services) into resource groups.
- [Models](/pr-preview/pr-22676/concepts/logical/models/): Meshery uses a set of resource models to define concrete boundaries to ensure extensible and sustainable management.
- [Organizations](/pr-preview/pr-22676/concepts/logical/organizations/): Organizations are the unit of tenancy in Meshery — they group users together and own all of the resources those users create.
- [Patterns](/pr-preview/pr-22676/concepts/logical/patterns/): Meshery Patterns are descriptive, declarative characterizations of how your Kubernetes infrastructure should be configured.
- [Policies](/pr-preview/pr-22676/concepts/logical/policies/): Meshery Policies enable you with a broad set of controls and governance of the behavior of systems under Meshery's management.
- [Registry](/pr-preview/pr-22676/concepts/logical/registry/): Meshery Registry is a database acting as the central repository for all capabilities known to Meshery. These capabilities encompass various entities, including models, components, relationships, and policies.
- [Relationships](/pr-preview/pr-22676/concepts/logical/relationships/): Meshery Relationships identify and facilitate genealogy between Components.
- [Workspaces](/pr-preview/pr-22676/concepts/logical/workspaces/): Meshery Workspaces act as central collaboration point for teams.
