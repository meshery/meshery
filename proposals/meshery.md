# Meshery CNCF Proposal

Name of project: Meshery

Description: Meshery is an open source multi-service mesh management plane.

Alignment with CNCF - Why does CNCF need a multi-service mesh management plane?
-------------------------------------------------------------------------------
The CNCF has an impressive portfolio of projects that can be leveraged to build and run complex distributed systems; our team believes Meshery to be a great fit for the CNCF. Meshery's core mission aligns well with Kubernetes, Linkerd and the container and networking ecosystem. 

The CNCF's mission is to "create and drive the adoption of a new computing paradigm that is optimized for modern distributed systems environments capable of scaling to tens of thousands of self-healing multi-tenant nodes." Modern distributed systems rely on networking, connectivity, observability, security, uniform control. As a result, service meshes are an essential piece of this architecture.

Meshery Overview
================

Features
--------

-   Meshery allows you to benchmark service meshes

-   Meshery helps you to quickly try out different services meshes

-   Meshery allows to persist and compare benchmark results

-   Meshery supports Service Mesh Interface

To learn more about Meshery's features, read and view the following resources:

-   <https://layer5.io/meshery/>

-   <https://meshery.io/blog/a-standard-interface-for-service-meshes>

Project Timeline and Snapshot
-----------------------------

-   Layer5 launched and open sourced Meshery in January 2019 as a multi-service mesh manager

Users
-----

-   Ziglu, TicketMaster, HPE, HashiCorp

In-Flight Features
------------------

The Meshery team is currently working on the following feature improvements:

-   Exporting load test results to a common accepted Service Mesh Performance (SMP)

The direction of the project has generally been guided by our open source community and users. There are a plethora of GitHub issues requesting various features that we prioritize based on popularity of user requests and engineering capacity.

A roadmap for future features, including those listed above, can be found in GitHub at <https://github.com/meshery/meshery>.

The project welcomes contributions of any kind: code, documentation, bug reporting via issues, and project management to help track and prioritize workstreams.

Use Cases 
----------

The following is a list of common use-cases for Meshery users:

-   Service Meshes installation with default configuration

-   Trying out sample applications on the service meshes

-   Ability to run performance tests on applications

-   Persist and compare results of different performance tests

CNCF Donation Details
---------------------

-   Preferred Maturity Level: Sandbox

-   Sponsors: [@lcalcote](https://twitter.com/lcalcote), [@justincormack](https://twitter.com/justincormack), [@michellenoorali](https://twitter.com/michellenoorali)

-   License: Apache 2

-   Source control repositories / issue tracker:  <https://github.com/meshery/meshery>.

-   Infrastructure Required: N/A

-   Website:  [https://meshery.io](https://meshery.io/)

-   Release Methodology and Mechanics: Documented at <https://meshery.io/releases>

Social Media Accounts
---------------------

-   Twitter:  <https://twitter.com/mesheryio>

-   Slack:  [http://slack.layer5.io](http://slack.layer5.io/)

Contributor Statistics (valid as of Mar 2020)
---------------------------------------------

-   230 GitHub Stars

-   1,900 Issues Opened

-   94 Contributors

-   42 releases

Asks from CNCF
--------------

-   A vendor-neutral home for Meshery to facilitate growth and community involvement

-   Logistics -- General access to resource and staff to provide advice, and help optimize our growth

-   Infrastructure for CI / CD

-   Integration with CNCF devstat

Appendix
--------

### Architecture 

Meshery is architected as a client of the Kubernetes API, leveraging Envoy. You can learn more about its architecture at <https://layer5.io/meshery/#architecture>

### Landscape

There are numerous service meshes available for developers and platform architecture teams to leverage. 

An analysis of the various options will be performed at a future time.