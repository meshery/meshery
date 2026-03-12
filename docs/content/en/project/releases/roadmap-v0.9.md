This post includes the proposed roadmap for Meshery v0.9.0 release. Use this thread to refine and hardened these plans and assignments. If you are interested in contributing to or leading any of these items, please comment below.

_Resources_
1. [v0.9.0 Milestone](https://github.com/meshery/meshery/milestone/6) in `meshery/meshery`.
1. [Meshery Architecture deck](https://docs.google.com/presentation/d/1SQMfyu5shjpGKlYONdVzOtd7UYTgLWBcgUvHMLCZ2tY/edit#slide=id.g1044af767ce_0_127) for additional visual designs and functional specifications.
1. [Meshery API Endpoints](https://docs.google.com/spreadsheets/d/1ABJCbfQRi0uN_YoP2kmHZ-lTI4S4QuvRk1_unomKNRE/edit)
2. [Meshery CLI Global Design Spec](https://docs.google.com/document/d/1xRlFpElRmybJ3WacgPKXgCSiQ2poJl3iCCV1dAalf0k/edit)

## Infrastructure Lifecycle and Configuration Management
- **Support Azure and ASO** ([spec](https://docs.google.com/presentation/d/1SQMfyu5shjpGKlYONdVzOtd7UYTgLWBcgUvHMLCZ2tY/edit#slide=id.gcb74201a11_0_229)) [?]
    * Static [Model Generattion](https://docs.google.com/document/d/1rUwDANT6-fk-d_8TUh73wFBeTRTOCE-tuDEQeHRchZA/edit#heading=h.1rjqe5qvw5j9)  sources ASO [?]
    * Relationships for ASO [?]
    * Connections, Credentials for Azure
    * CLI: Refactor `system config` for AKS ([spec](https://docs.google.com/document/d/1XfIvMwKGKBS5_ielGWuYFdqJByEcOeloCKpUL9bBxlw/edit#heading=h.blih70a9hxp)) [@hexxdump???]
  * CLI: Refactor `system config` for EKS ([spec](https://docs.google.com/document/d/1XfIvMwKGKBS5_ielGWuYFdqJByEcOeloCKpUL9bBxlw/edit#heading=h.blih70a9hxp)) [@hexxdump ???]

[details="Roadmap v0.9.0 (stretch goal)"]
- :warning: GitOps: Expand Flux Integration with Helm Repo [?]
[/details]

**MeshSync**
* Connection Registration and Enrichment [?]
* Azure, AWS, GCP: Discovery and Enrichment of Public Cloud resources [?]
* Make meshsync apis realtime [@Aabid.Sofi ]

[details="Roadmap v1.0 (stretch goal)"]
- MeshSync: Configurable and Tiered Discovery [?]
  - Support for runtime registration of new fingerprints. [?]
  - Support for composite fingerprints where keys span multiple entities (or not just the entity being discovered). [?]
- CLI:  Meshconfig: Support for multi-cluster ([spec](https://docs.google.com/document/d/1r_Yopt904qdqO6lutzZn8mfqjKlI-MQgD22GG8x4UYM/edit#)) [?]
[/details]

## Registry & Models [:meshery: ?]
**Models**
* Support for Sub-Categories [?]

**Policies**
- Refactor Policies as first-class resource. [@Lee]
- Policies within Models: import/export, OCI [?]
   - Evaluation: Enhancing the mechanism used to reference policies on disk or in registry. [?]
- Registry support in Meshery UI [?]
- `mesheryctl policy` command [?]

**Relationships**
- Basic relationship coverage for all Models ([spec](https://docs.google.com/spreadsheets/d/1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw/edit#gid=0)) [?]

**Connections and Credentials**
- Registration of connections and credentials [?]

[details="Roadmap v0.9.0 (stretch goal)"]
- Connection: Supporting OCI Registries ([#8855](https://github.com/meshery/meshery/issues/8855)) ([Design Spec](https://docs.google.com/document/d/1AmNnwtMbYSVQ00TZRYtEZBJQTrPtH_S51ZG2UiwBGqs/edit))
- Signing OCI exported images ([#8855](https://github.com/meshery/meshery/issues/8855))
- Export a Design to GitHub [?]
- Export a Model to Docker Hub [?]
[/details]

## Spaces [:meshery:@amit ]
* Integration of Workspaces [@amit]

[details="Roadmap v0.9.0 (stretch goal)"]
Stretch
* Integration of Environments
[/details]

## Extensibility & Extensions [:meshery: @theBeginner86]
  - Helm Kanvas Snapshot [@Ijeoma.Eti @theBeginner86]
  - Kubectl Meshsync Snapshot [ @devhindo ]

[details="Roadmap v0.9.0 (stretch goal)"]
Stretch
  - Kubectl Kanvas Snapshot [?]
[/details]

## Adapters & Performance Management [:meshery:?]

* Add and manage Prometheus connections [@codeknight03]
* Add and manage Grafana connections [@codeknight03]
* Working JSON board with chart display [@codeknight03]

[details="Roadmap v0.9.0 (stretch goal)"]
1. :warning: Adapter: [Evolve meshes.proto to adapters.proto](https://github.com/meshery/meshery/issues/9968)
1. :warning: :nighthawk: Externalize Nighthawk ([spec](https://docs.google.com/document/d/1Qy_BpHXibTvL9daLTc5L4fjmr1F051ya0UYeghNUq_I/edit?usp=sharing)) [@Lee][@Xin_Huang][@MUzairS15]
    1. :warning: :nighthawk:  Nighthawk as a Meshery Adapter (`meshery-nighthawk`) [@Xin_Huang]
    1. :warning: :nighthawk: Server: Lifecycle Management of `meshery-nighthawk` in Meshery Server [@MUzairS15]
    1. :warning: :nighthawk: CLI: `mesheryctl perf` support for load generator / adapter selection [?]
    1. :warning: :nighthawk: UI: `meshery-nighthawk` adapter chip [@Yash.Sharma]
[/details]

#### Basic Adaptive Load Controller
[details="Roadmap v0.9.0 (stretch goal)"]
**Basic Adaptive Load Controller** ([spec](https://docs.google.com/presentation/d/1arewKYZuhkv7NGRZKx0w-7bQdS-ho-YSYJF0X32qfHo/edit#slide=id.gfbc0c4fb3f_0_8)) (*v0.9.0*)
  * Exposure of “Adaptive Test” in Meshery UI. [?]
  * Meshery Server endpoints. [@Xin_Huang]
  * Storage of n result sets in Provider for a given Performance Profile.  [@Xin_Huang], [@theBeginner86]
[/details]

#### Custom Telemetry Provider
[details="Roadmap v0.9.0 (stretch goal)"]
- Multiple Telemetry Providers and Custom Boards ([spec](https://docs.google.com/presentation/d/1SQMfyu5shjpGKlYONdVzOtd7UYTgLWBcgUvHMLCZ2tY/edit#slide=id.g1044af767ce_5_21)) [?]
- User-defined Dashboards and Metrics ([spec](https://docs.google.com/presentation/d/1SQMfyu5shjpGKlYONdVzOtd7UYTgLWBcgUvHMLCZ2tY/edit#slide=id.gcb74201a11_0_119)) [?]
[/details]

## Docs [:meshery: ?]

* [Playground] Learning Paths: Cloud Native tutorials with Meshery Playground [#9832](https://github.com/meshery/meshery/issues/9832)
[@vishalvivekm]
    * Creating three tutorials per category, allowing learners to perform hands-on lab interactively in the live cluster environment in Meshery Playground. [?]
* Automated publication of End-to-End test results [?]
* Dynamic Menu: Docs Sidebar Table of Contents [?]

## :warning: UI [:meshery: @Antonette.Caldwell @SAHU-01]
* :warning: UI: Switch to Sistent as the theme provider [@Sudhanshu_Dasgupta]
   * UI: All components sourced from Sistent.
   * Creation of new Components in Sistent Design System
* :warning: UI: Support for Material UI v5 [@SAHU-01 @Sudhanshu_Dasgupta]

## CLI [:meshery: @Matthieu.EVRIN @alphaX86 @hexxdump]
- :warning: CLI: code coverage goal: 50% [?]
[details="Roadmap v0.9.0 (stretch goal)"]
-  CLI: `system report` - diagnostics reporting *(stretch goal)* [?]
[/details]
- CLI: Support for commands: policies, models, relationships, components, connections, credentials [@alphaX86] [@Matthieu.EVRIN][@Riya]

## Build, Test, Release [:meshery: ??? ]
-  Server: End-to-End Testing coverage goal: 50% [:meshery: Ian Whitney ?][@palSagnik ]

### Playground [:meshery: @Sangram.Rath]
- Remote Provider: Deployment using Meshery (as a Meshery Design)
- Remote Provider: Multi-master Postgres

## Meshery Catalog [:meshery: @vishalvivekm]
* :warning: Models in Catalog
* Addition of user-defined, custom tags as filters
* Robust category/type filters

<hr>

## Legend

> **Theme** or **Meshery Component**
> :meshery: [\@Theme Champion]
> *[?] = unassigned/help wanted*
> :x: incomplete/no implementation
> :warning: partial implementation
> :white_check_mark: done
