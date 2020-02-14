== Meshery Proposal

*Name of project:* Meshery

*Description:* Meshery is an open source, multi-service mesh management plane.

=== Alignment with CNCF - Why does CNCF need a multi-service mesh management plane?

The CNCF has an impressive portfolio of projects that can be leveraged to build and run complex distributed systems; our team believes Meshery to be a great fit for the CNCF. Meshery's core mission aligns well with Kubernetes, Linkerd, Envoy, Network Service Mesh, and the container and networking ecosystem. 

The CNCF's mission is to “create and drive the adoption of a new computing paradigm that is optimized for modern distributed systems environments capable of scaling to tens of thousands of self-healing multi-tenant nodes.” Meshery furthers the mission of facilitating adoption of cloud native infrastructure. Modern distributed systems rely on networking and connectivity. Commonly, developers and operators focused on building and managing cloud native infrastructure are unfamiliar, and therefore, uncomfortable with virtualized networking. Meshery aims to empower them to operate with ease. 

=== Meshery Overview

==== Features

 * Envoy, Linkerd and Network Service Mesh Inside - Meshery is built as the management plane for these projects.
 * 

To learn more about Meshery's features, read and view the following resources:

 * [Meshery Project](https://meshery.io)
 * [Meshery Documentation](https://docs.meshery.io)
 
=== Project Timeline and Snapshot
 * Layer5 launched and open sourced Contour in February 2019 as an multi-mesh manager.
 * Meshery's one year anniversity is marked by a healthy number and diverse set of contributors and users.
 
== Production Users
 * TBD

== In-Flight Features

The Meshery team is currently working on the following feature improvements:
 * 
 * 

The direction of the project has generally been guided by our open source community and users. There are a plethora of GitHub issues requesting various features that we prioritize based on popularity of user requests and engineering capacity. 

A roadmap for future features, including those listed above, can be found in GitHub at https://github.com/meshery/meshery/ROADMAP.md
The project welcomes contributions of any kind: code, documentation, bug reporting via issues, and project management to help track and prioritize workstreams.

== Use Cases
The following is a list of common use-cases for Meshery users:  

 * Performance management of service meshes and applications running on top of them.
 * Lifecycle management of the most popular service meshes.
 * 

== CNCF Donation Details
 * *Preferred Maturity Level:* Sandbox
 * *Sponsors:* @leecalcote, 
 * *License:* Apache v2
 * *Source control repositories / issue tracker:* https://github.com/layer5io/meshery.
 * *Infrastructure Required:* N/A
 * *Website:* https://meshery.io/
 * *Release Methodology and Mechanics:* Documented at https://docs.meshery.io/releases

== Social Media Accounts:

 * *Twitter:* https://twitter.com/mesheryio
 * *Youtube:* https://www.youtube.com/channel/UCFL1af7_wdnhHXL1InzaMvA/playlists
  * *Slack:* http://slack.layer5.io

== Contributor Statistics (valid as of Feb 2020)
 * 150 GitHub Stars
 * 1,605 Commits
 * 322 PRs
 * 625 Issues Opened
 * 44 Contributors
 * 40 releases

== Asks from CNCF
 * A vendor-neutral home for Meshery to facilitate growth and community involvement
 * Logistics – General access to resource and staff to provide advice, and help optimize our growth
 * Infrastructure for CI / CD
 * Integration with CNCF devstat

== Appendix

=== Architecture
Meshery is architected as a client of the Kubernetes API, leveraging Envoy. You can learn more about its architecture at https://docs.meshery.io/architecture

== Landscape
There are numerous service meshes controllers available for developers and platform architecture teams to leverage. An analysis of the various options is presented on https://layer5.io/landscape