This post includes the proposed roadmap for Meshery v0.8.0 release. Use this thread to refine and hardened these plans and assignments. If you are interested in contributing to or leading any of these items, please comment below.

## Legend

> **Theme** or **Meshery Component**
> :meshery: [@Theme Champion]
> *[?] = help wanted*
> :x: incomplete/no implementation
> :warning: partial implementation
> :white_check_mark: done

## Resources
1. See [v0.8.0 Milestone](https://github.com/meshery/meshery/milestone/5) in `meshery/meshery`.
1. See the [Meshery Architecture deck](https://docs.google.com/presentation/d/1SQMfyu5shjpGKlYONdVzOtd7UYTgLWBcgUvHMLCZ2tY/edit#slide=id.g1044af767ce_0_127) for additional visual designs and functional specifications.

## Lifecycle Management
- **Support AWS and GCP** ([spec](https://docs.google.com/presentation/d/1SQMfyu5shjpGKlYONdVzOtd7UYTgLWBcgUvHMLCZ2tY/edit#slide=id.gcb74201a11_0_229)) [@MUzairS15]
    * :white_check_mark:Static Model Generation sources AWS, GCP as K8s CRDs. [@MUzairS15]
    * :white_check_mark: Relationships: AWS [@SAHU-01]
    * :warning: Relationships: GCP
    * :x: Connections, Credentials for AWS, GCP
    * :x: CLI: Refactor `system config` for AKS ([spec](https://docs.google.com/document/d/1XfIvMwKGKBS5_ielGWuYFdqJByEcOeloCKpUL9bBxlw/edit#heading=h.blih70a9hxp)) [@hexxdump]
  * :x: CLI: Refactor `system config` for EKS ([spec](https://docs.google.com/document/d/1XfIvMwKGKBS5_ielGWuYFdqJByEcOeloCKpUL9bBxlw/edit#heading=h.blih70a9hxp)) [@hexxdump]

- :warning: GitOps: Expand Flux Integration with Helm Repo [?]


* :x:MeshSync: Configurable and Tiered Discovery [:meshery:@Daniel.Kiptoon]
  - :warning: Support for runtime registration of new fingerprints. [@MUzairS15]
  - :x: Support for composite fingerprints where keys span multiple entities (or not just the entity being discovered). [?]
[details="Roadmap v0.9.0 (stretch goal)"]
- CLI:  Meshconfig: Support for multi-cluster ([spec](https://docs.google.com/document/d/1r_Yopt904qdqO6lutzZn8mfqjKlI-MQgD22GG8x4UYM/edit#)) [?]
[/details]

## Registry & Models [:meshery:@MUzairS15]
**Models**
* :white_check_mark: Refactor Model as first-class resource. [@Shlok_Mishra]
  * :white_check_mark: Model import/export  [@Shlok_Mishra]
* :white_check_mark: Server: [Model Generator](https://docs.google.com/document/d/1rUwDANT6-fk-d_8TUh73wFBeTRTOCE-tuDEQeHRchZA/edit#heading=h.1rjqe5qvw5j9) for direct Helm Chart or K8s Manifest Retrieval
  * :white_check_mark: Dynamic Generation via URL (Chart, Manifest) upon Design import. [@Shlok_Mishra ]
* :x: Support for Sub-Categories [?]

**Relationships**
- :white_check_mark: Refactor Relationship as first-class resource. [@Lee]
   - :white_check_mark: Relationship Coverage for All K8s Components ([spec](https://docs.google.com/spreadsheets/d/1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw/edit#gid=0)) [@Lee][@Ripul.Handoo][@SAHU-01]
   - Evaluation: Enhancing the mechanism used to reference policies on disk or in registry. [?]
   - :white_check_mark: Classification of Policy Evaluation into client-side, server-side, and cluster-side moments of evaluation and their distinctions. [@Lee]
   - :white_check_mark: SelectorSets for Relationships [@MUzairS15]

**Registry** [:meshery:@Yash.Sharma]
- :white_check_mark: UI: (Capabilities) Registry [:meshery:@Yash.Sharma]
- :warning: Registration of connections and credentials [?]
- :white_check_mark: CLI: Support for commands: models, relationships, components, connections, credentials [@alphaX86] [@Matthieu.EVRIN]

### Open Container Image (OCI) Support [:meshery:@theBeginner86]
- :white_check_mark: Adopt OCI as the packaging and distribution format for Meshery Designs, Patterns, Filters, Models ([#6447](https://github.com/meshery/meshery/issues/6447)) ([Design Spec](https://docs.google.com/document/d/1AmNnwtMbYSVQ00TZRYtEZBJQTrPtH_S51ZG2UiwBGqs/edit))
- :white_check_mark: CLI: Pushing and pulling Models and Designs from OCI registries [?] ([Design Spec](https://docs.google.com/document/d/1AmNnwtMbYSVQ00TZRYtEZBJQTrPtH_S51ZG2UiwBGqs/edit))

[details="Roadmap v0.9.0 (stretch goal)"]
- Connection: Supporting OCI Registries ([#8855](https://github.com/meshery/meshery/issues/8855)) ([Design Spec](https://docs.google.com/document/d/1AmNnwtMbYSVQ00TZRYtEZBJQTrPtH_S51ZG2UiwBGqs/edit))
- Signing OCI exported images ([#8855](https://github.com/meshery/meshery/issues/8855))
- Export a Design to GitHub [?]
- Export a Model to Docker Hub [?]
[/details]

## Spaces [:meshery:@theBeginner86]
* :white_check_mark: Environments [@senali]
* :white_check_mark: Workspaces [@amit]

## Extensibility & Extensions [?]

- :white_check_mark: UI: Extensible Authorization [@Yash.Sharma]
- :white_check_mark: Expand Meshery's Integration with Artifact Hub [#9966](https://github.com/meshery/meshery/issues/9966) [@Chris.Carrier]
  - :white_check_mark: [Artifact Hub] Meshery Designs as a new Artifact Hub kind [#9967](https://github.com/meshery/meshery/issues/9967) [@Chris.Carrier]

[details="Roadmap v0.9.0 (stretch goal)"]
Stretch
  - :warning: Helm Kanvas Snapshot [@Ijeoma.Eti @theBeginner86]
  - Kubectl Meshsync Snapshot [?]
  - Kubectl Kanvas Snapshot [?]
[/details]

## Adapters & Performance Management [:meshery:@Xin_Huang] [@Kunyue.Xing]
1. :warning: Adapter: [Evolve meshes.proto to adapters.proto](https://github.com/meshery/meshery/issues/9968)
1. :warning: :nighthawk: Externalize Nighthawk ([spec](https://docs.google.com/document/d/1Qy_BpHXibTvL9daLTc5L4fjmr1F051ya0UYeghNUq_I/edit?usp=sharing)) [@Lee][@Xin_Huang][@MUzairS15]
    1. :warning: :nighthawk:  Nighthawk as a Meshery Adapter (`meshery-nighthawk`) [@Xin_Huang]
    1. :warning: :nighthawk: Server: Lifecycle Management of `meshery-nighthawk` in Meshery Server [@MUzairS15]
    1. :warning: :nighthawk: CLI: `mesheryctl perf` support for load generator / adapter selection [?]
    1. :warning: :nighthawk: UI: `meshery-nighthawk` adapter chip [@Yash.Sharma]

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


## Docs [:meshery: @iArchitSharma][@Awani_Alero]

* :white_check_mark: [Playground] Learning Paths: Cloud Native tutorials with Meshery Playground [#9832](https://github.com/meshery/meshery/issues/9832)
[@vishalvivekm]
    * Creating three tutorials per category, allowing learners to perform hands-on lab interactively in the live cluster environment in Meshery Playground. [@sandramsc]
* :white_check_mark: Automated publication of security key reference [@Yash.Sharma]
* :white_check_mark: Automated publication of Model Relationships and Components [@theBeginner86][@vishalvivekm]
* :x: Dynamic Menu: Sidebar Table of Contents [@iArchitSharma]

## UI [:meshery: @Antonette.Caldwell @SAHU-01]
* :white_check_mark: UI: Complete incorporation of RTK [@Sudhanshu_Dasgupta]
   * UI: Complete elimination of useEffect.
* UI: Switch to Sistent as the theme provider[@Sudhanshu_Dasgupta ]
   * :white_check_mark:UI: All components sourced from Sistent .
* :white_check_mark: UI: Sistent Design System
   * packaging and build/release [@Antonette.Caldwell]
   * design [@Rex_Joshua]
* :warning: UI: Support for Material UI v5 [@SAHU-01]

## CLI [:meshery: @alphaX86 @hexxdump]
- :warning: CLI: code coverage goal: 50% [?]
- :white_check_mark:  CLI: Full migration from Apps to Designs  [@theBeginner86]
[details="Roadmap v0.9.0 (stretch goal)"]
-  CLI: `system report` - diagnostics reporting *(stretch goal)* [?]
[/details]

## Build, Test, Release [:meshery: ]
-  :white_check_mark: Server: End-to-End Testing coverage goal: 25% [:meshery: @Jerens_Lensun]
- :white_check_mark: Transition from Cypress to Playwright [:meshery: @Aabid.Sofi @Jerens_Lensun]

### Playground [@Sangram.Rath]
- :white_check_mark: Migration from Docker Compose to Kubernetes deployment
- :white_check_mark: Remote Provider migration from AWS to Equinix

## Meshery Catalog [:meshery: @vishalvivekm]
* :white_check_mark: 100% design-type coverage
* :white_check_mark: 50% technology-type coverage
* :warning: Models in Catalog

<hr>

See the [Meshery Architecture deck](https://docs.google.com/presentation/d/1SQMfyu5shjpGKlYONdVzOtd7UYTgLWBcgUvHMLCZ2tY/edit#slide=id.g1044af767ce_0_127) for additional visual designs and functional specifications.
