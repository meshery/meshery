# Meshery Docs

<h2>Contributing to the Meshery Documentation</h2>

Please do! Thank you for your help in improving [Meshery Docs](https://docs.meshery.io)! :balloon:

<details>
<summary>
Find the complete set of Meshery Docs contributor guides at https://docs.meshery.io/project/contributing/contributing-docs</summary>

Before contributing, please review the [Documentation Contribution Flow](https://github.com/meshery/meshery/blob/master/CONTRIBUTING.md#documentation-contribution-flow). In the following steps you will set up your development environment, fork and clone the repository, run the site locally, and finally commit, sign-off, and push any changes made for review.

<h3>1. Set up your development environment</h3>

- _The Meshery Docs site is built using Jekyll - a simple static site generator! You can learn more about Jekyll and setting up your development environment in the [Jekyll Docs](https://jekyllrb.com/docs/)._

- First [install Ruby](https://jekyllrb.com/docs/installation/), then install Jekyll and Bundler.

**Note:** Windows users can run Jekyll by following the [Windows Installation Guide](https://jekyllrb.com/docs/installation/windows/) and also installing Ruby Version Manager [RVM](https://rvm.io). RVM is a command-line tool which allows you to work with multiple Ruby environments on your local machine. Alternatively, if you're running Windows 10 version 1903 Build 18362 or higher, you can upgrade to Windows Subsystem for Linux [WSL](https://docs.microsoft.com/en-us/windows/wsl/install-win10) and run Jekyll in Linux instead.

Alternatively, if you are running Windows 10, you may install the Windows Subsystem for Linux:

- [WSL1](https://docs.microsoft.com/en-us/windows/wsl/install-win10) for Windows build version 1607 or higher.

### 2. Get the code

- Fork and then clone the [Meshery repository](https://github.com/meshery/meshery)
  ```bash
  $ git clone https://github.com/YOUR-USERNAME/meshery
  ```
- Change to the docs directory
  ```bash
  $ cd docs
  ```
- Install any Ruby dependencies
  ```bash
  $ bundle install
  ```

<h3>3. Serve the site</h3>

- Serve the code locally
  ```bash
  $ make docs
  ```
  _Note: From the Makefile, this command is actually running `$ bundle exec jekyll serve --drafts --livereload`. There are two Jekyll configuration, `jekyll serve` for developing locally and `jekyll build` when you need to generate the site artifacts for production._

### 4. Create a Pull Request

- After making changes, don't forget to commit with the sign-off flag (-s)!
  ```bash
  $ git commit -s -m “my commit message w/signoff”
  ```
- Once all changes have been committed, push the changes.
  ```bash
  $ git push origin <branch-name>
  ```
- Then on Github, navigate to the [Meshery repository](https://github.com/meshery/meshery) and create a pull request from your recently pushed changes!

---

- _See the [Meshery Documentation Google Doc](https://docs.google.com/document/d/17guuaxb0xsfutBCzyj2CT6OZiFnMu9w4PzoILXhRXSo/edit) for additional reference._
- Theme - https://github.com/vsoch/docsy-jekyll
</details>

# Meshery Documentation Structure

<details>
<summary>High-Level Outline & Information Architecture for Meshery Documentation</summary>

**Goal:** Offer comprehensive, organized, and accessible documentation for diverse audiences, from new users to expert contributors.
**Target Audience:**
- **Personas:** Beginners, developers, admins, operators, security specialists, contributors, users of all experience levels.
- **Needs:** Varied - learning fundamentals, managing tasks, understanding advanced concepts, contributing code.


## High-Level Outline

### Overview and Installation (User Persona)

- **Getting Started:** Overview of Meshery, installation options, prerequisites, and setup instructions.
- **Installation Guides:** Step-by-step instructions for installing Meshery on different platforms (local, cloud, minikube).
- **Configuration Guides:** Configuring Meshery for different environments (local, cloud, minikube).
- **Concepts:** Meshery basics (clusters, pods, deployments, services), terminology glossary.
- **Use Cases:** Demonstrations of common scenarios (web app deployment, data processing pipeline).

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

### Administrator's Guide (Expert & Operator Persona) (Advanced Topics)

- **Multi-Meshery Management:** Federation, cluster federation, GitOps for configuration management.
- **Performance Optimization:** Resource usage analysis, profiling tools, tuning techniques.
- **GitOps** DevOps & CI/CD integration - Integrating Meshery with continuous integration and deployment pipelines.
- **Best Practices:** Recommendations for securing the Meshery, monitoring performance, managing versions.
- **Advanced Concepts:** Advanced Meshery concepts, features, and capabilities.
- **Air-gapped Environments:** Deploying Meshery in air-gapped environments.
- **Troubleshooting Guides:** Identifying and resolving common errors, debugging techniques.

### Integrations and Extensions (All Personas)

- **Integrations:** Integrating Meshery with different infrastructure and systems.
- **Extensibility:** Customizing Meshery with plugins, adapters, and extensions.
- API Reference: Comprehensive reference for Meshery API objects and fields.
  - **Extension Points** Meshery extension points for different capabilities.
    - Providers, plugins, adapters, and modules.
- **Extensions** Meshery adapters, plugins, and modules for different extionsion points.
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
