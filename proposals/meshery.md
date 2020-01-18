== Meshery Proposal

*Name of project:* Meshery

*Description:* Meshery is an open source multi-service mesh management plane.

=== Alignment with CNCF - Why does CNCF need a multi-service mesh management plane?

The CNCF has an impressive portfolio of projects that can be leveraged to build and run complex distributed systems; our team believes Meshery to be a great fit for the CNCF. Meshery's core mission aligns well with Kubernetes, Linkerd Envoy and the container and networking ecosystem. The CNCF's mission is to “create and drive the adoption of a new computing paradigm that is optimized for modern distributed systems environments capable of scaling to tens of thousands of self-healing multi-tenant nodes.” Modern distributed systems rely on networking and connectivity. As a result, ingress controllers for the compute platform, Kubernetes, are an essential piece of this architecture. 

=== Meshery Overview

==== Features

 * Envoy Inside - Meshery is built as the control plane for Envoy, the high performance L7 proxy and load balancer
 * 

To learn more about Meshery's features, read and view the following resources:

 * https://projectcontour.io/announcing-contour-1.0/[Intro to Contour]
 * https://projectcontour.io/routing-traffic-to-applications-in-kubernetes-with-contour/[Routing Traffic to Applications in Kubernetes with Contour]
 * https://projectcontour.io/httpproxy-in-action/[HTTPProxy in Action]
 * https://projectcontour.io/resources/[Contour Resources]

=== Project Timeline and Snapshot
 * Heptio launched and open sourced Contour in October 2017 as an ingress controller for Kubernetes based on Envoy
 * Contour 1.0 ships in November 2019, signaling to the community a commitment to a stable set of APIs 
 
== Production Users
 * TBD

== In-Flight Features

The Contour team is currently working on the following feature improvements:
 * Using Contour to redirect requests to locations other than Pods, rewriting the ExternalName in the host header for proxied requests
 * Migration tooling for IngressRoute to HTTPProxy

The direction of the project has generally been guided by our open source community and users. There are a plethora of GitHub issues requesting various features that we prioritize based on popularity of user requests and engineering capacity. 

A roadmap for future features, including those listed above, can be found in GitHub at https://github.com/projectcontour/contour#workspaces/contour-5bc5116124028a7e4bf2ef81/board?repos=108462822. 
The project welcomes contributions of any kind: code, documentation, bug reporting via issues, and project management to help track and prioritize workstreams.

== Use Cases
The following is a list of common use-cases for Contour users:  

 * High performance ingress controller for Kubernetes 
 * Multi-team and multi-tenant ingress controller for Kubernetes 
 * Blue/Green and canary deployments

== CNCF Donation Details
 * *Preferred Maturity Level:* Sandbox
 * *Sponsors:* @leecalcote, 
 * *License:* Apache 2
 * *Source control repositories / issue tracker:* https://github.com/layer5io/meshery.
 * *Infrastructure Required:* N/A
 * *Website:* https://meshery.io/
 * *Release Methodology and Mechanics:* Documented at https://projectcontour.io/resources/release-process/

== Social Media Accounts:

 * *Twitter:* https://twitter.com/mesheryio
 * *Google Groups:* https://groups.google.com/forum/#!forum/projectcontour-announce
 * *Slack:* http://slack.layer5.io

== Contributor Statistics (valid as of Jan 2020)
 * 2k GitHub Stars
 * 1120 PRs
 * 689 Issues Opened
 * 77 Contributors
 * 40 releases

== Asks from CNCF
 * A vendor-neutral home for Meshery to facilitate growth and community involvement
 * Logistics – General access to resource and staff to provide advice, and help optimize our growth
 * Infrastructure for CI / CD
 * Integration with CNCF devstat

== Appendix

=== Architecture
Meshery is architected as a client of the Kubernetes API, leveraging Envoy. You can learn more about its architecture at https://projectcontour.io/docs/v1.0.1/architecture/

== Landscape
There are numerous ingress controllers available for developers and platform architecture teams to leverage. An analysis of the various options will be performed at a future time.