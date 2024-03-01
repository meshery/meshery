# Meshery Proposal for CNCF Incubation

## About Meshery

[Meshery](https://meshery.io/) is a self-service engineering platform, designed to manage the complete lifecycle of modern infrastructure.

The project encompasses tools that can be utilized through a command-line interface, integrated into clusters, or incorporated into various other tools. It enables users to assess cloud native configurations during development, integration, and post-deployment phases.

## Meshery in the CNCF Sandbox

Since its entry into the CNCF Sandbox on June 22, 2021, Meshery has undergone significant enhancements and expansions. 

* Supports every CNCF project and public cloud services that have native resource / service integration with Kubernetes.
* Introduction of advanced load generation and performance characterization features.
* Enhanced support for multi-cluster and hybrid cloud environments.
* Meshery has evolved into a full-blown self-service engineering platform with a robust extensibility framework for easy integration of new extensions whether they are backend adapters, frontend plugins, or separately hosted providers of additional functionality (e.g. identity providers).

### Community & Growth

Since joining the CNCF Sandbox, Meshery has experienced notable community growth and engagement:

* [Is the 10th fastest growing project in the CNCF](https://www.cncf.io/blog/2023/10/27/october-2023-where-we-are-with-velocity-of-cncf-lf-and-top-30-open-source-projects/).
* Accumulated a total of 4,000 new GitHub stars
* Closed 9,000 [Pull Requests (PRs)](https://github.com/pulls?q=is%3Aissue+created%3A%3E%3D2021-06-21+org%3Ameshery), showcasing an active contributor community.
* Addressed and resolved 3,000 [issues](https://github.com/issues?q=is%3Aissue+created%3A%3E%3D2021-06-21+org%3Ameshery+), ensuring ongoing project improvement.
* Participated heavily in mentoring initiatives and is the #1 most popular LFX internship.

### Neutrality

Meshery, committed to maintaining vendor neutrality, ensuring a diverse collection of maintainers representing 10 different organizations. As an extensible platform, Meshery goes to great lengths to support third-party plugins, elevating them within its ecosystem of integrations. Meshery places emphasis on fostering an open and inclusive community of contributors, maintainers, and integrators. 

## Incubation stage requirements

### Used successfully in production by at least three independent direct adopters

Meshery's non-exhaustive, public list of [adaptors](https://github.com/meshery/meshery/blob/master/ADOPTERS.md) is available for reference. With thousands of users already leveraging its features, Meshery has proven its production-readiness and garnered positive feedback from adopters and the broader community.

### Have a healthy number of committers
Meshery boasts a diverse and active community of maintainers, with [15 maintainers](https://github.com/meshery/meshery/blob/main/MAINTAINERS.md) overseeing the ~30 repositories included in the `meshery` GitHub organization.  and [extension modules](https://github.com/meshery/meshery/blob/main/EXTENSIONS.md). 

### Demonstrates a substantial ongoing flow of commits and merged contributions

Meshery is the 10th fastest growing CNCF project.

Over the past year, over [300 committers have had PRs merged](https://meshery.devstats.cncf.io/d/66/developer-activity-counts-by-companies?orgId=1&var-period_name=Last%20year&var-metric=merged_prs&var-repogroup_name=All&var-country_name=All&var-companies=All), and [800 contributors actively participating](https://meshery.devstats.cncf.io/d/66/developer-activity-counts-by-companies?orgId=1&var-period_name=Last%20year&var-metric=contributions&var-repogroup_name=All&var-country_name=All&var-companies=All).

### Clear versioning scheme

Meshery follows a clear and [well-documented](https://docs.meshery.io/project/contributing/build-and-release) build and release process. All components adhere to the principles of [semantic versioning](https://semver.org/). The project releases page documents [Meshery Server, UI, and CLI releases](https://docs.meshery.io/project/releases) with Meshery Adapters, Meshery Operator and custom controllers, centralized API schemas, and so on all releasing independently. See each repo for a release history.

### Clearly documented security processes explaining how to report security issues to the project, and describing how the project provides updated releases or patches to resolve security vulnerabilities

Meshery's [security reporting process](https://docs.meshery.io/project/security-vulnerabilities) is well-documented and has [previous and current CVEs](https://docs.meshery.io/project/security-vulnerabilities) published along with the reporting process and expectation setting of how reports are handled. This process is followed successfully with GitHub engineers being the most recent to report vulneraabilities and a Meshery maintainer (from Intel) to provide patches. Mainatiners strive to acknowledge reports within 24 hours of being received and offer an analysis and plan for a full response withing 10 days.

## Alignment with the CNCF Mission

Meshery aligns closely with the CNCF mission by contributing to the security and efficiency of cloud-native computing. Focused solely on cloud native management, Meshery integrates deeply and seamlessly with _many_ CNCF projects. The full list of project integrations can be seen on https://meshery.io/integrations. Meshery's commitment to interoperability and collaboration is one of a number of reasons that it is a valuable asset to the CNCF ecosystem.

## Future plans

Looking ahead to 2024, Meshery's roadmap includes the introduction of advanced features in collaborative workflows of developer-defined infrasturcture and the engineering teams that operate this infrastructure. These additions will be crafted through a combination of new AI-centric developments, integration of workflow engine, and potential launch as the [CNCF's official project playground](https://docs.google.com/document/d/1Cr0MxlOxWq70d-BUisfKXF_VPGss0hK8XqRxy4xm4YE/edit#heading=h.58lqw93jp55u). Meshery remains dedicated to evolving as a key player in the cloud native landscape, responding to user needs and industry trends.

