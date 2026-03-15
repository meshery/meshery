---
layout: page
title: Meshery Documentation Structure and Organization
permalink: project/contributing/contributing-docs-structure
abstract: Audience, high-Level outline & information architecture for Meshery Documentation
language: en
type: project
category: contributing
list: include
display-title: false
---

Providing well-structured, easy-to-navigate documentation tailored for a diverse audience—from beginners to experts.


## Target Audience

We first identify **who** this documentation serves before diving into how it’s organized. By clarifying different user personas, we can better direct them to relevant sections:

- **Users**  
  - **New Users**: Individuals just getting started with Meshery, needing a straightforward, guided setup.  
  - **Advanced Users**: Those who already have some experience and want deeper insights or specialized use cases.

- **Contributors & System Integrators**  
  - **Open Source Contributors**: Developers who extend Meshery’s features, fix bugs, or contribute documentation.  
  - **System Integrators**: Professionals adding Meshery into existing ecosystems or building custom integrations.

- **Administrators**  
  - **Operators, DevOps, SREs, IT professionals, system administrators, security specialists**: Responsible for large-scale deployments, performance tuning, security, and overall upkeep of Meshery in production environments.

## High-Level Outline

After defining our audience, we introduce the **overall structure** of Meshery’s documentation. Each section below addresses a different need, level of expertise, or user focus.

1. [Overview and Installation (User Persona)](#overview-and-installation-user-persona)
2. [Tasks & Operations (User Guide)](#tasks--operations-user-guide)
3. [Administrator’s Guide (Advanced Topics)](#administrators-guide-advanced-topics)
4. [Integrations and Extensions (All Personas)](#integrations-and-extensions-all-personas)
5. [Tutorials (All Personas)](#tutorials-all-personas)
6. [Reference & Resources (All Personas)](#reference--resources-all-personas)
7. [Contributing and Community (All Personas)](#contributing-and-community-all-personas)

## Overview and Installation (User Persona)

This section focuses on giving new and experienced users an overarching understanding of Meshery, including basic setup, installation options, and general operations.

- **Getting Started:** Overview of Meshery, installation options, prerequisites, and setup instructions.
- **Installation Guides:** Step-by-step instructions for installing Meshery on different platforms (local, cloud, minikube).
  - **Configuration Guides:** How to configure Meshery in these environments.
- **Tasks & Operations** (User Guide)
- **Concepts:** Basic Meshery concepts (clusters, pods, deployments, services), terminology glossary.
- **Use Cases:** Demonstrations of common scenarios (web app deployment, data processing pipeline).

### Getting Started (New Users)

This section provides a concise, hands-on introduction for newcomers, referring back to “Overview and Installation” for more details.

1. **Introduction to Meshery**:
   - What is Meshery? Why use it? Key features, benefits, use cases.  
   - Basic terminology glossary.  
2. **Installation and Setup**
   - Prerequisites, system requirements.
   - Platform-specific installation guides (local, cloud, air-gapped).
   - Configuration guides, Verifying Installation.
3. **First Steps**
   - Basic Meshery UI/CLI walkthrough.
   - Deploying a simple application.
   - Connecting to a cluster.

### Concepts (All Personas)

Concept pages describe **what** a particular aspect of Meshery is, but not **how** to do it step-by-step. They provide background knowledge and link out to related tasks or tutorials.

- **Architectural Concepts:** Meshery architecture, design, and implementation details. Diagrams illustrating interaction between components, resource dependencies.
- **Logical Concepts:** Meshery components, resources, and relationships. Diagrams illustrating interaction between components, resource dependencies.
<!-- - **Deep Dives:** Detailed explanations of core Meshery components. -->

### Tasks & Operations (User Guide)

A task page shows how to do a **single** thing, typically by giving a short sequence of steps. Task pages have minimal explanation, but often provide links to conceptual topics that provide related background and knowledge.

- **Task Guides:** Step-by-step instructions for common tasks (deploying applications, managing resources).
- **Configuration Management:** Designing infrastructure, managing configuration files.
- **Lifecycle Management:** Discoverying, registering, configuring infrastructure
  - Discovery (MeshSync): Greenfield/Brownfield scenarios.
  - Managing Connections: Registering, updating, and deleting connections.
  - Managing Credentials: Registering, updating, and deleting credentials.
- **Performance Management:** Load testing, performance monitoring, resource usage analysis.
<!-- - **Workflows:** Step-by-step procedures for common tasks (rolling updates, blue-green deployments). -->

### Administrator's Guide (Platform Engineer & Operator Persona, Advanced Topics)

Designed for platform engineers, operators, DevOps, and other advanced users who require deeper insights into Meshery’s configuration, security, and maintenance.

- **Performance Optimization:** Resource usage analysis, profiling tools, tuning techniques.
- **GitOps / DevOps & CI/CD**: Integrating Meshery with continuous integration and deployment pipelines.  
- **Best Practices:** Recommendations for securing the Meshery, monitoring performance, managing versions.
- **Advanced Concepts:** Advanced Meshery concepts, features, and capabilities.
- **Air-gapped Environments:** Deploying Meshery in air-gapped environments.
- **Troubleshooting Guides:** Common errors, debugging, tips.  
- **Multi-Meshery Management:** Federation, cluster federation, GitOps for configuration management.
- **Security and Compliance:** Security best practices, compliance requirements, vulnerability management.
- **Backup and Recovery:** Backup strategies, disaster recovery planning, data restoration.
- **Scaling and High Availability:** Scaling Meshery, high availability, load balancing, failover strategies.
- **Monitoring and Logging:** Monitoring Meshery, logging, alerting, observability tools.
- **Upgrading Meshery:** Version compatibility, upgrade paths, release notes, rollback procedures.

### Integrations and Extensions (All Personas)

Meshery’s extensibility allows for seamless integration with existing systems and the addition of new capabilities through adapters, plugins, and extensions.

- **Extensibility:** Customizing Meshery with plugins, adapters, and extensions.
  - **APIs:** Simple summary of Meshery's APIs.
  - **Extension Points:** Providers, plugins, adapters, modules.
- **Integrations:** Meshery Models. Integrating Meshery with different infrastructure and systems.
- **Extensions:** List of all of Meshery adapters, plugins, and modules for different extension points.
  - **Adapters:** Integrating Meshery with different infrastructure and extended capabilities.
  - **Plugins:** Extended capabilities for Meshery.
  - **Remote Providers:** Remote providers for additional functionality.
  - **Security and Identity:** Authentication, authorization, secrets management, vulnerability scanning.

### Tutorials (All Personas)

A tutorial page shows how to accomplish a goal that is larger than a single task. Typically a tutorial page has several sections, each of which has a sequence of steps. For example, a tutorial might provide a walkthrough of a code sample that illustrates a certain feature of Kubernetes. Tutorials can include surface-level explanations, but should link to related concept topics for deep explanations.

- **Tutorials:** Dedicated walk-throughs with labs and step-by-step instructions using Meshery's features.

### Reference & Resources (All Personas)

A component tool reference page shows the description and flag options output for a Meshery component. For example, a component tool reference page might describe the Meshery CLI and explain the role it plays as an application once it is deployed, scaled, and updated. Typically, component tool reference pages don't include sequences of steps, but instead provide links to tasks or tutorials.

- **Command References:** Detailed explanations and examples for `mesheryctl` commands, API resources.
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
