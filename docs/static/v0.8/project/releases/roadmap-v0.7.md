The proposed Meshery v0.7.0 Roadmap is listed by component. Use this thread to refine and hardened these plans and assignments. If you are interested in contributing to or leading any of these items, please comment below.

See the [Meshery Architecture deck](https://docs.google.com/presentation/d/1SQMfyu5shjpGKlYONdVzOtd7UYTgLWBcgUvHMLCZ2tY/edit#slide=id.g1044af767ce_0_164) for additional visual designs and functional specifications.

**Legend:**
> **Component**
> [Feature Champion] Feature Description
> *[?] = help wanted*
> :x: incomplete/no implementation
> :warning: partial implementation
> :white_check_mark: done

**Meshery Server**

* :white_check_mark: Design Deployment Dry-run
* :white_check_mark: Enhance Kubernetes multi-cluster support with Environments
* :white_check_mark: Messaging Framework and Enhanced Notification Center ([spec](https://docs.google.com/document/d/1_zpEWBcC6ngOLen_E_nc-octiUdBbpmdwwIX6qzo70M/edit#))
* :white_check_mark: [@revoly] Support for Provider Enforcement
* :white_check_mark: [@Nithish_karthik] Scheduled Workflow: Static component generation for Helm Charts
* :white_check_mark: [@Nithish_karthik] Realtime CLI/UI: Static component generation for Helm Charts
* :warning: [@Philip-21] Unit coverage goal: 15%
* :warning: [Mario Arriaga] Functional test coverage: 25%

**Meshery Adapters**

* :x: [@Antonette.Caldwell] Adapter for App Mesh → Stable
* :x: [@Ruturaj] MeshOps v2 support for NSM Adapter
* :x: [@MUzairS15] Implementation of Messaging Framework
* :warning: [?] Unit coverage goal: 25%
* :white_check_mark: [?] Functional coverage goal: 70%

**Meshery CLI**

* [?] :white_check_mark: `mesh deploy` with MeshOps v2
* [?] :white_check_mark: Meshconfig: Support for Provider Enforcement
* [?] :x: Meshconfig: Support for multi-cluster ([spec](https://docs.google.com/document/d/1r_Yopt904qdqO6lutzZn8mfqjKlI-MQgD22GG8x4UYM/edit#))
* :white_check_mark: [@harkiratsm] `system dashboard`
* :x: [@hexxdump] Refactor `system config` for EKS ([spec](https://docs.google.com/document/d/1XfIvMwKGKBS5_ielGWuYFdqJByEcOeloCKpUL9bBxlw/edit#heading=h.blih70a9hxp))
* :warning: [@Philip-21] Code coverage goal: 80% ([spec](https://docs.google.com/document/d/1xRlFpElRmybJ3WacgPKXgCSiQ2poJl3iCCV1dAalf0k/edit#heading=h.rzpmb66db1sq))

**Meshery Operator**

* :white_check_mark: [@Daniel.Kiptoon] MeshSync: Expanded scope of object synchronization
* :x: [?] MeshSync: Dynamically-loaded fingerprints
* :x: [?] MeshSync: Discovery pipeline stages
* :warning: Code coverage goal: 25%

**Meshery Database**
* :warning: Support for Postgres as a deployment choice
:white_check_mark: [@revoly] SQLite → Postgres ([spec](https://docs.google.com/presentation/d/1SQMfyu5shjpGKlYONdVzOtd7UYTgLWBcgUvHMLCZ2tY/edit#slide=id.gcb74201a11_0_229))

:x: **Meshery Perf**
*  :x:Externalize Nighthawk ([spec](https://docs.google.com/document/d/1Qy_BpHXibTvL9daLTc5L4fjmr1F051ya0UYeghNUq_I/edit#heading=h.w6zppv50p453))
  * :x: [@Xin_Huang] Nighthawk as a Meshery Component
  * :x: [@Antonette.Caldwell] Lifecycle Management of `meshery-perf`
  * :white_check_mark: [@Abhi] Allow user to identify which Kubernetes Cluster

* Basic Adaptive Load Controller ([spec](https://docs.google.com/presentation/d/1arewKYZuhkv7NGRZKx0w-7bQdS-ho-YSYJF0X32qfHo/edit#slide=id.gfbc0c4fb3f_0_8))
  * [@Abhi] Exposure of “Adaptive Test” in Meshery UI.
  * [@Xin_Huang] Meshery Server endpoints.
  * [?] Storage of n result sets in Provider for a given Performance Profile.

**Meshery UI**

* :warning: [@Aabid.Sofi] State Management with Redux ([spec](https://docs.google.com/document/d/1sb7XZIczx9FqKA3Hpd1XgtpEjfj8qjqMbgKY68t9Sko/edit#heading=h.blih70a9hxp))
  * Refactor of all components ([tracker](https://docs.google.com/spreadsheets/u/0/d/1RlUcr-iLtCFCsV1VZvWD0jBijPJAx-Rk-1GLLsL9sSI/edit))
* :white_check_mark:WASM Envoy Filter Management
* :white_check_mark: [@abhi] Global: Cluster Selector
*  :white_check_mark:[@Abhi] Global: MeshSync Pulse (design)
* :white_check_mark: React 17, Nextjs 13
* [Mario Arriaga] Functional test coverage: 25%

**Meshery.io**

* :white_check_mark: [?] Meshery Catalog ([spec](https://docs.google.com/presentation/d/1SQMfyu5shjpGKlYONdVzOtd7UYTgLWBcgUvHMLCZ2tY/edit#slide=id.gfc8b2b2554_0_0))

----
See the [Meshery Architecture deck](https://docs.google.com/presentation/d/1SQMfyu5shjpGKlYONdVzOtd7UYTgLWBcgUvHMLCZ2tY/edit#slide=id.g1044af767ce_0_164) for additional visual designs and functional specifications like this one -

![Meshery v0.6.5 - Data Persistence|690x387](upload://mkFFy7ERkWPzUI2OOhsGLfNW2Sh.jpeg)
