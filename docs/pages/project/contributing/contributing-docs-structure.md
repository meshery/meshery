---
layout: page
title: "Meshery Docs: Information Architecture"
permalink: project/contributing/contributing-docs-structure
abstract: Audience, high-Level outline & information architecture for Meshery Documentation
language: en
type: project
category: contributing
list: include
display-title: false
---

# High-Level Outline & Information Architecture for Meshery Documentation

**Goal:** Offer comprehensive, organized, and accessible documentation for diverse audiences, from new users to expert users.
**Target Audience:**

1. Persona: Users
   1. Types: New Users, Advanced Users
2. Persona: Contributors and System Integrators (people extending Meshery)
   1. Types: Open Source Contributors, System Integrators professionally extending Meshery
3. Persona: Administrators
   1. Types: Operators, DevOps, SREs, IT professionals, system administrators, security specialists.

## High-Level Outline

Three documentation sets to cater to different personas. Each documentation set is tailored to the specific needs and expertise levels of the target audience.

### Overview and Installation (User Persona)

- **Getting Started:** Overview of Meshery, installation options, prerequisites, and setup instructions.
- **Installation Guides:** Step-by-step instructions for installing Meshery on different platforms (local, cloud, minikube).
  - **Configuration Guides:** Configuring Meshery for different environments (local, cloud, minikube).
- **Tasks & Operations** (User Guide)
- **Concepts:** Meshery basics (clusters, pods, deployments, services), terminology glossary.
- **Use Cases:** Demonstrations of common scenarios (web app deployment, data processing pipeline).

### Getting Started (New Users)

1. Introduction to Meshery:
   1. Overview, key features, benefits, use cases.
   2. "What is Meshery?" and "Why use Meshery?"
   3. Terminology glossary (basic terms).
2. Installation and Setup:
   1. Prerequisites, system requirements.
   2. Platform-specific installation guides (local, cloud, air-gapped).
   3. Configuration guides.
   4. Verifying Installation.
3. First Steps:
   1. Basic Meshery UI/CLI walkthrough.
   2. Deploying a simple application.
   3. Connecting to a cluster.

### Concepts (All Personas)

A concept page explains some aspect of Meshery. For example, a concept page might describe the Meshery Models object and explain the role it plays as an application once it is deployed, scaled, and updated. Typically, concept pages don't include sequences of steps, but instead provide links to tasks or tutorials.

- **Architectural Concepts:** Meshery architecture, design, and implementation details. Diagrams illustrating interaction between components, resource dependencies.
- **Logical Concepts:** Meshery components, resources, and relationships. Diagrams illustrating interaction between components, resource dependencies.
<!-- - **Deep Dives:** Detailed explanations of core Meshery components. -->

### Tasks & Operations (User Guide)

A task page shows how to do a single thing, typically by giving a short sequence of steps. Task pages have minimal explanation, but often provide links to conceptual topics that provide related background and knowledge.

- **Task Guides:** Step-by-step instructions for common tasks (deploying applications, managing resources).
- **Configuration Management:** Designing your infrastructure, managing configuration files.
- **Lifecycle Management:** Discoverying, registering, configuring infrastructure
  - Discovery (MeshSync)
    - Greenfield
    - Brownfield
  - Managing Connections
    - Registering, updating, and deleting connections.
  - Managing Credentials
    - Registering, updating, and deleting credentials.
- **Performance Management:** Load testing, performance monitoring, resource usage analysis.
<!-- - **Workflows:** Step-by-step procedures for common tasks (rolling updates, blue-green deployments). -->

### Administrator's Guide (Platform Engineer & Operator Persona) (Advanced Topics)

Platform Engineers, Operators, DevOps, SREs, IT Professionals, System Administrators, Security Specialists

- **Performance Optimization:** Resource usage analysis, profiling tools, tuning techniques.
- **GitOps** DevOps & CI/CD integration - Integrating Meshery with continuous integration and deployment pipelines.
- **Best Practices:** Recommendations for securing the Meshery, monitoring performance, managing versions.
- **Advanced Concepts:** Advanced Meshery concepts, features, and capabilities.
- **Air-gapped Environments:** Deploying Meshery in air-gapped environments.
- **Troubleshooting Guides:** Identifying and resolving common errors, debugging techniques.
- **Multi-Meshery Management:** Federation, cluster federation, GitOps for configuration management.
- **Security and Compliance:** Security best practices, compliance requirements, vulnerability management.
- **Backup and Recovery:** Backup strategies, disaster recovery planning, data restoration.
- **Scaling and High Availability:** Scaling Meshery, high availability, load balancing, failover strategies.
- **Monitoring and Logging:** Monitoring Meshery, logging, alerting, observability tools.
- **Upgrading Meshery:** Version compatibility, upgrade paths, release notes, rollback procedures.

### Integrations and Extensions (All Personas)

- **Extensibility:** Customizing Meshery with plugins, adapters, and extensions.
  - **APIs** Simple summary of Meshery's APIs.
  - **Extension Points** Meshery extension points for different capabilities.
    - Providers, plugins, adapters, and modules.
- **Integrations:** Meshery Models. Integrating Meshery with different infrastructure and systems.
- **Extensions** List of all of Meshery adapters, plugins, and modules for different extennsion points.
  - **Adapters** Integrating Meshery with different infrastructure and extended capabilities.
  - **Plugins** Meshery plugins for different capabilities.
  - **Remote Providers** Meshery remote providers for different capabilities.
  - **Security and Identity:** Authentication, authorization, secrets management, vulnerability scanning.

### Tutorials (All Personas)

A tutorial page shows how to accomplish a goal that is larger than a single task. Typically a tutorial page has several sections, each of which has a sequence of steps. For example, a tutorial might provide a walkthrough of a code sample that illustrates a certain feature of Kubernetes. Tutorials can include surface-level explanations, but should link to related concept topics for deep explanations.

- **Tutorials:** Dedicated walk-throughs with labs and step-by-step instructions using Meshery's features.

### Reference & Resources (All Personas)

A component tool reference page shows the description and flag options output for a Meshery component. For example, a component tool reference page might describe the Meshery CLI and explain the role it plays as an application once it is deployed, scaled, and updated. Typically, component tool reference pages don't include sequences of steps, but instead provide links to tasks or tutorials.

- **Command References:** Detailed explanations and examples for mesheryctl commands, API resources.
- **API Documentation:** Comprehensive reference for Meshery API objects and fields.
- **Custom Resource Definition Reference:** Comprehensive reference for Meshery CRDs.
- **Release Notes:** Detailed information about version changes, new features, deprecations.
- **Glossary:** Definitions of common terms, acronyms, and abbreviations.
- **Vulnerability Reports:** Security advisories, CVEs, and vulnerability reports.

### Contributing and Community (All Personas)

- **External Resources:** Links to blogs, community forums, case studies, books, training materials.
- **Contributing Guide:** How to contribute documentation, code, and other resources to the project.
- **Community:** Highlight community forums, events, contributor guidelines, recognition.
- **FAQ:** Answers to frequently asked questions.

This high-level outline provides a comprehensive framework for structuring the Meshery documentation, catering to diverse user needs while ensuring information is readily accessible and actionable. By further refining each section with specific content recommendations and considering the needs of specific personas, the documentation can effectively serve as a valuable resource for everyone interacting with Meshery.
</details>
