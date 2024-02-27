Meshery Design Document: Meshery Model Generation  
Status: Draft | Under Review | Approved

Meshery Model Generation
A factory of component generators

---

Prologue 3
Design Goals and Objectives 3
Actors Glossary 4
Generating Models and Components 5
Flow: Statically Generated Components 5
Flow: Dynamically Generated Components 6
Flow: Importing Designs 6
Model Lifecycle 6
Providing granular control for Models 7
Functional Architecture 8
Architecture Diagram 8
Sequence Diagram 8
System Flows 9
Flow: <short title> 9
Flow: <short title> 9
Action Items 10
ActionItem 1: Integration Spreadsheet 10
ActionItem 2: Generate components from multiple sources. 10
ActionItem 3: 11
ActionItem 4: Doc generation being supported by all Model Generators 11
ActionItem 5: Generate components from different formats 11
Sources of Interaction with the spreadsheet: 11

---

Model Factory

See also: Meshery Models

    Prologue

Entities under Meshery’s management must be represented in Meshery’s native Model format. In order for Meshery to support a broad variety of platforms, a uniform, but extensible factory for the creation of Models and their components must be implemented. Meshery needs to facilitate the generation of Models and Components from a variety of sources (e.g. a URL or an OCI registry) and of a variety of formats (e.g. a Kubernetes Manifest or a Helm Chart or Operator Lifecycle Manager (OLM)).

1. Models are distinct from, but related to Designs
   Designs should benefit from Model generation.
   Model generation should consider the needs of Designs and of Patterns.
2. Versioning of all things is necessary
   Versioning of schemas and of models and of designs and so on.
3. Plan for mandatory user involvement
   Example: During Connection Registration
   Example: During importing of or opening of a Design
   Understand Designs may or may not be valid. A particular Meshery Server may or may not have all components represented in a design.
   Design Goals and Objectives
   The designs in this specification should result in upholding the following goals and should achieve these specific functions:

Goal 1: Translation between infrastructure and visual design tools
Infrastructure

1. Support for Kubernetes CRDs
2. Support for Terraform Providers
3. Support for GCP Providers
4. Support for AWS Providers
   Drawing Tools
5. Mermaid.js
6. Draw.io (LucidCharts)
7. Excalidraw (Google Cloud Architecture Tool, too)
8. Nice to have: eraser.io

Goal 2: Make Models versatile by supporting a broad set of actions.

1. Import and export of Models in OCI format.
2. Check for Duplicate components upon Model import. Track “DuplicatedBy” as an array.
3. Use `Status` field to store labels such as Preferred and Ignored.
4. Change versioning Schemas to semver and start making releases, so that the Golang package can be imported. (meshery/schemas/issues/43)
5. Ensure that we’re including the schema version in Designs (meshery/meshery/issues/9773)
   Actors Glossary
   Registry - a component within Meshery that contains a database of known capabilities.
   Registrar - the Meshery Server process responsible for the process of record keeping.
   Entity (Registree) - an capability (an entity) that is found in the Meshery Registry; e.g. a model, component, relationship, policy.Sometimes referred to as a capability.
   Entity Source - an entity’s original location from which it was sourced; e.g. (source_uri is used as the flag by th pattern engine to assess whether additional support). The Entity Source should have all the information that Meshery needs to generate the components.  
   Registrant (host) - A Connection responsible for the sourcing and registration of entities. Sometimes Registrants perform registration for their own components, while other times they act on behalf of the component source.
   Supports pulling packages from Artifact Hub and other sources like Docker Hub.
   Should envelope Meshery Application importer.

mesheryctl model repo add meshery/meshery

Tracks a list of packages (e.g. helm chart names) to import or not import from the given Package Importer.

Package Manager will be responsible for tracking and managing the list of Packages.

So, A `PackageManager` knows how to:

1. Get the packages

The PackageManager interface will look something like this:
type PackageManager interface {
GetPackages() (Package, error)
}
These importers are mechanical in nature; machines to configured to interact with the repo-specific APIs
E.g. Artifact Hub, Docker Hub
Environment - a collection of Connections. Sometimes the target of a deployment.
Connection - an individual instance of one type of Meshery’s supported integrations. Sometimes the target of a deployment.
Generating Models and Components
Component Generator (Meshery Server) using MeshKit’s functions
Meshery Server is scheduled to run on schedule, generate components, and persist components in the local database.

As a user of MeshMap, I would like to mark some categories and components as displayed or not.
Prioritize component generation for “Official” repos first.

Generating Models does not require Meshery Server

1. Model and Component generation logic is MeshKit. `mesheryctl` and Meshery Server both utilize MeshKit’s libraries for ongoing programmatic generation of models and components.
   Flow: Statically Generated Components
   Goal: As a user, I would like to have a large number of integrations (models) available in Meshery out-of-the-box without having to find and generate my own, so that I can quickly move along with either creating new Designs or using my existing Designs.
   Actors:

- GitHub Actions, mesheryctl, Meshery Integrations spreadsheet
  Preconditions:

1. List of model sources is defined in Meshery Integrations spreadsheet.
2. Scheduled workflow is configured to run daily.
   Assumptions:
3. Meshery Users: Statically generated models are available in Meshery container image.
4. Meshery Developers: Statically generated models are available in meshery/meshery repo.
   Behavior:
5. Schedule: Component generation as part of nightly workflow and mesheryctl script updates components based on the attribute value from the spreadsheet, pushing these components to meshery/meshery and updating spreadsheet.
6. The Meshery project’s scheduled workflow is only additive in the list of statically-generated models. In other words, if a scheduled workflow fails or if a particular model source is unavailable or a particular model within a given source is unavailable, the Meshery project’s scheduled workflow will not remove support for that model.
7. Meshery Server on boot registers Models generated from supported registrants (e.g. Artifact Hub).
   Implementation:
8. One Helm connection will be available in every Meshery deployment. This Helm connection will have many repos configured.[a]
9. Label all duplicative components that are not of the “Preferred” model with “Ignored”.
   1. “Preferred”
   2. “Ignored” - Not shown in the UI.
   3. “DuplicatedBy”: Shown in the UI, with (below “Preferred” model), the preference can be changed from Settings/
10. Error Handling:
11. Acceptance Criteria:
12. Helm is the Registrant (and Connection) with many repos configured.
    Flow: Dynamically Generated Components
    Import and Export of Models
13. Dynamic generation and registration occurs upon initial connection[b] to each supported platform (e.g. a Kubernetes context transitions to Discovered state)
14. [Adapter] pushes its specific capabilities to Meshery Server. In order to register capabilities, Adapters will send a POST request to/api/meshmodel/components/register.
15. [Meshery Server] Components can be registered with the server by interfacing with MeshModel APIs, mesheryctl and Meshery UI are two clients (will be) capable of invoking generation and registration. User imports a Meshery Application as a single file reference, folder reference (to be recursively searched), or a git repository to be walked.
    Flow: Importing Designs
    Importing Designs requires Meshery Server, so that your Design can be associated with your account.
16. Improve failure experience when importing Designs that have components that are registered.
17. Cloud only supports importing of Designs and Patterns. It does not support import and translation of charts, compose, manifests, and so on.
    Registry Lifecycle
    Users should have control to enable/disable Categories, Models and Components, this can be an individual/team-level preference. As a result, the meshmodel APIs should be intelligent enough to expose only those Models that are enabled by the user.
    Providing granular control for Models
    What is “preferred-models.yaml”?
    [c]
    A manually curated list of models, referred during component generation (static/dynamic) which controls what is shown in the UI to the user when designing infrastructure.
    The Registry UI still shows all the models irrespective of their preference in “preferred-models.yaml”.

18. Meshery tracks a list of preferred models and components in preferred-models.yaml. The Model Generator generates and registers models and components for all imported or discovered entities.
    1. If a model has duplicative components of those found in a model in the preferred list, that model and its components will still be generated and will still be registered, however, they will have an “Ignored” label assigned.
19. Having a “preferred” list means inviting duplicative components across models.
20. AI: Move from whitelist to preferred list, while still generating duplicative components. Label all duplicative components that are not of the “Preferred” model with “Ignored”.
    1. “Preferred”
    2. “Ignored” - Not shown in the UI.
    3. “DuplicatedBy”: Shown in the UI, with (below “Preferred” model), the preference can be changed from Settings/

The default behavior will be to expose all components, with the exception of a few, such as those components which have invalid schema or components which cannot be configured on MeshMap, such components will be marked as disabled in the Integrations spreadsheet.

The presence of schema will take precedence over the user’s preference, i.e. if a component is enabled by the user but has an invalid/empty schema such components will not be available (unless explicitly requested by the user).

Considering the current situation of Meshmap, Kubernetes CRD is not configurable and we experience a crash. Adding granular control for MeshModels will ensure that, by default such components are disabled. This in no sense means that UI should be allowed to crash, instead, UI should handle it gracefully with appropriate fallback mechanisms.

Categories, Models and Components each will have an attribute which specifies whether it is enabled / disabled.[d]

Categories will take precedence over models, which in turn takes precedence over components, in case of “disabled” behavior, i.e. If a model is marked as disabled, all components for that model will be disabled irrespective of their individual behavior. Similarly for Categories and models.

1. Dynamic Registration / Adapter registering itself:
   1. Component already registered statically: Dynamic registration will still occur inheriting property from static component.
   2. Component is not registered statically: Component will be marked as enabled.
2. Meshery Server (Mesheryctl command / interfacing with MeshModel APIs)
3. Component already registered: Registration will still occur inheriting property from static component.
4. Component not registered: Component will be marked as enabled.

FAQ

- How does it solve the problem at hand for CRD?
  The “hasSchema” attribute present in the spreadsheet needs to be renamed to be more meaningful, something that depicts the above-discussed behavior. This attribute will decide whether the component is enabled/disabled for the end user. CRD components will be marked as disabled, by default.[e][f]

- What happens if during a collaborative session user A has model A and model B, but user B has model A disabled.
  Union of enabled categories/modes/components of user A and user B will be taken. For each user we have the required information, regarding enabled resources, in situations like this, the user will be informed indicating some of the disabled resources have been enabled, as part of a collaborative session. The synchronization of enabled/disabled resources can happen at the time of initializing a collaboration session or at runtime.

Union of Enabled Resources: Meshery could synchronize the registered Models between user A and user B. In this case, Meshery would combine the registered Models from both users, resulting in model A being enabled for the collaborative session. User B would be informed that some of their disabled resources have been temporarily enabled due to the collaboration.
Functional Architecture
Although right now the only source we will be using is CRDs, this will change in the future, so the architecture should take that into account.

Our sources can be anything. It can be ArtifactHub packages, DockerHub etc. Depending on the source we use, the data type varies along with the logic.
Architecture Diagram
<here>
Sequence Diagram
<here>
System Flows
Flow: <short title>
Goal: As a [developer/integrator/mesheryctl/operator] user,
I would like to …….. ,
so that ………..
Actors:

- Preconditions:

3. Assumptions:
4. Behavior:
5. Implementation:
6. Error Handling:
7. Acceptance Criteria:
8. Flow: <short title>
   Goal: As a [developer/integrator/mesheryctl/operator] user,
   I would like to …….. ,
   so that ………..
   Actors:
9. Preconditions:
10. Assumptions:
11. Behavior:
12. Implementation:
13. Error Handling:
14. Acceptance Criteria:

Action Items

1. Remove hardcoded consideration of Artifact Hub in pattern engine. But has special hard-coded exceptions for Meshery Server… but why when we have a registry and should capture this in the registry?
   There is no hardcoded assumption for ArtifactHub there’s an interface defined and the artifact hub source implements this interface. Based on the source from which comp is generated accordingly the method is called.
2. Making Model first-class;
3. Support for Relationships in Models other than Kubernetes
4. Support for sources other than Artifact Hub
   1. Vihas - URL as source for dynamic generation fo
5. Import and export of Model.
6. OCI support -
   ActionItem 1: Integration Spreadsheet
   Integration spreadsheet needs to specify the
7. source: Helm repostiories/OCI/url(github)/file
8. Format: Helmchart/manifests/Terraform/compose
9. Link: l[g][h]ink to the source targz file in case of helm, repository link for github or a remote URL
   ActionItem 2: Generate components from multiple sources.
   The format of the source components will either be Helm Charts/K8s manifests/(compose file ?)
10. Helm repositories
11. OCI regsitries
12. URL/filesystem

Each source should implement the interface:

GetComponents
UpdateComponents
Process // Any specific conversion
Generate
HandleLifecycle // Used at the time of deployment, to provisioning/de-provisioning of required Operators/Controllers

Currently there’s an interface defined IHost (see the definition of Host/Registrant above) (https://github.com/meshery/meshkit/blob/9b066f8219d796da5b77ad695c6a802c41339e60/models/meshmodel/registry/host.go#L70C1-L70C1)
ActionItem 3:

1. The functions should be written in meshkit repo.
2. Create new commands in mesheryctl to source and generate components.
3. Enhance/Create APIs in Meshery Server.

Refer the notes shared in slack during implementation. (Support for optionally updating spreadsheet is also mentioned there)
ActionItem 4: Doc generation being supported by all Model Generators

ActionItem 5: Generate components from different formats

1. Terraform
2. OLM

Any other formats (and not specifically K8S manifests)

Key points:
The factory should be such that after it has been implemented (all 5 Action Items), we have the ability to generate components with different combinations
eg: Helm Charts from OCI/URL/file
Terraforms provider from OCI/URL/file
Sources of Interaction with the spreadsheet:
At the time of generation:
The spreadsheet will be referenced and can optionally updated some columns. (hasSchema column)

At the of component update:
To update the SVGs, shape, colour and other metadata.
During this step spreadsheet is referenced as read only

// Supporting Connection and credential definitions and relationships.
// Publishing to Targets

1.  Docs site
2.  Providers
    // Making models first class
    // Importing relationships
    // Export of whatever is generated (All sources needs to implement the export functionality).
    // Plugin/adapter model to handle deployments of those components (i.e. MeshKit should not have the logic for managing lifecycle, [HandleLifcecycle implemented by that specific plugin/adatper]

