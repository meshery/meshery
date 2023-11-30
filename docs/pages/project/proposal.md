# Meshery Proposal for CNCF Incubation

## About Meshery

[Meshery](https://meshery.io/) is a cloud native management plane designed to facilitate the management, configuration, and performance testing of cloud native in cloud-native environments. 

The project encompasses tools that can be utilized through a command-line interface, integrated into clusters, or incorporated into various other tools. It enables users to assess cloud native configurations during development, integration, and post-deployment phases.

Meshery stands out for its commitment to providing a comprehensive approach to cloud native security, adopting and automating best practices. The project supports various cloud native frameworks, including Istio, Linkerd, and Consul.

## Meshery in the CNCF Sandbox

Meshery entered the CNCF Sandbox on June 22, 2021, and has since undergone substantial enhancements and expansions:

* Inclusion of support for additional cloud native frameworks such as Envoy and Kuma.
* Introduction of advanced traffic management features.
* Integration with popular observability tools like Grafana and Jaeger.
* Enhanced support for multi-cluster and hybrid cloud environments.
* Implementation of a robust extensibility framework for easy integration of new features and extensions.

### Community & Growth

Since joining the CNCF Sandbox, Meshery has experienced notable community growth and engagement:

* [Is the 10th fastest growing project in the CNCF](https://www.cncf.io/blog/2023/10/27/october-2023-where-we-are-with-velocity-of-cncf-lf-and-top-30-open-source-projects/).
* Accumulated a total of 4,000 new GitHub stars
* Created 8,000 [Pull Requests (PRs)](https://github.com/pulls?q=is%3Aissue+created%3A%3E%3D2021-06-21+org%3Ameshery), showcasing an active contributor community.
* Addressed and resolved 2,600 [issues](https://github.com/issues?q=is%3Aissue+created%3A%3E%3D2021-06-21+org%3Ameshery+), ensuring ongoing project improvement.
* Participated in mentoring initiatives, guiding hundreds students through the Linux Foundation (LFX) program.

### Neutrality

Meshery, committed to maintaining vendor neutrality, ensures its tools and services are accessible to users without bias. Any commercial services associated with Meshery are transparently presented and separated from the core project. Meshery places emphasis on fostering an open and inclusive ecosystem.

## Incubation stage requirements

### Used successfully in production by at least three independent direct adopters

Meshery has demonstrated robust adoption with thousands of users leveraging its features directly or through integration with commercial solutions. Feedback from adopters and the broader community has been positive, validating its production-readiness.

### Have a healthy number of committers

Meshery boasts a diverse and active community of contributors, with 10 maintainers overseeing the [main repository](https://github.com/meshery/meshery/blob/main/MAINTAINERS.md) and [extension modules](https://github.com/meshery/meshery/blob/main/EXTENSIONS.md). Over the past year, [85 committers have had PRs merged](https://meshery.devstats.cncf.io/d/66/developer-activity-counts-by-companies?orgId=1&var-period_name=Last%20year&var-metric=merged_prs&var-repogroup_name=All&var-country_name=All&var-companies=All), and [280 contributors actively participated](https://meshery.devstats.cncf.io/d/66/developer-activity-counts-by-companies?orgId=1&var-period_name=Last%20year&var-metric=contributions&var-repogroup_name=All&var-country_name=All&var-companies=All).

### Demonstrates a substantial ongoing flow of commits and merged contributions

Meshery has maintained a consistent and active development pipeline, as evidenced by [1800 merged PRs](https://meshery.devstats.cncf.io/d/66/developer-activity-counts-by-companies?orgId=1&var-period_name=Last%20year&var-metric=merged_prs&var-repogroup_name=All&var-country_name=All&var-companies=All) over the past 12 months.

### Clear versioning scheme

Meshery follows a clear and well-documented versioning scheme for its components:

* [Meshery CLI and core services](https://github.com/meshery/meshery/releases)
* [Meshery Adapters](https://github.com/meshery/adapters/releases)
* [Meshery Operator](https://github.com/meshery/operator/releases)

All components adhere to the principles of [semantic versioning](https://semver.org/).

### Clearly documented security processes explaining how to report security issues to the project, and describing how the project provides updated releases or patches to resolve security vulnerabilities

Meshery employs the [GitHub security reporting process](https://github.com/meshery/meshery/blob/main/SECURITY.md), ensuring a swift response within 7 days for contact and 90 days for disclosure. The project adheres to a regular release schedule, with security patches integrated into the latest version.

## Alignment with the CNCF Mission

Meshery aligns closely with the CNCF mission by contributing to the security and efficiency of cloud-native computing. Focused solely on cloud native management, Meshery integrates seamlessly with various CNCF projects, including:

* Istio: primary data source for Istio configurations
* Envoy: supported cloud native framework
* Prometheus: metrics and observability
* Jaeger: distributed tracing
* OpenTelemetry: telemetry data collection
* Helm: installation and configuration

Meshery's commitment to interoperability and collaboration makes it a valuable asset to the CNCF ecosystem.

## Future plans

Looking ahead to 2024, Meshery's roadmap includes the introduction of advanced features in network security and enhanced runtime/threat detection. These additions will be crafted through a combination of new developments and integration with best-of-breed open-source engines. Meshery remains dedicated to evolving as a key player in the cloud-native landscape, responding to user needs and industry trends.

