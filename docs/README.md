# Meshery Docs

![assets/img/docs-screenshot.png](assets/img/docs-screenshot.png)
_Check out the [Meshery Docs](https://docs.meshery.io/)!_
Detailed documentation on contributing to Meshery docs is available here - [https://docs.meshery.io/project/contributing/contributing-docs](https://docs.meshery.io/project/contributing/contributing-docs).

## Contributing to the Meshery Documentation

Before contributing, please review the [Documentation Contribution Flow](https://github.com/meshery/meshery/blob/master/CONTRIBUTING.md#documentation-contribution-flow). In the following steps you will set up your development environment, fork and clone the repository, run the site locally, and finally commit, sign-off, and push any changes made for review.

### 1. Set up your development environment

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

### 3. Serve the site

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

# Meshery Documentation Stucture

## High-Level Outline & Information Architecture for Meshery Documentation

**Goal:** Offer comprehensive, organized, and accessible documentation for diverse audiences, from new users to expert contributors.

**Target Audience:**

- **Personas:** Beginners, developers, admins, operators, security specialists, contributors, users of all experience levels.
- **Needs:** Varied - learning fundamentals, managing tasks, understanding advanced concepts, contributing code.

**Information Architecture:**

**1. Getting Started (Beginner & User Persona):**

- **Concepts:** Meshery basics (clusters, pods, deployments, services), terminology glossary.
- **Tutorials:** Quick start guides for setting up different environments (local, cloud, minikube).
- **Use Cases:** Demonstrations of common scenarios (web app deployment, data processing pipeline).
- **Best Practices:** Recommendations for securing the cluster, monitoring performance, managing versions.

**2. Concepts & Architecture (All Personas):**

- **Deep Dives:** Detailed explanations of core Meshery components (pods, deployments, networking, storage).
- **Relationships:** Diagrams illustrating interaction between components, resource dependencies.
- **Advanced Concepts:** StatefulSets, DaemonSets, Ingress controllers, Service Meshes.
- **Troubleshooting Guides:** Identifying and resolving common errors, debugging techniques.

**3. Tasks & Operations (Admin & User Persona):**

- **Command References:** Detailed explanations and examples for kubectl commands, API resources.
- **Management Guides:** Configuring resources (limits, requests, health checks), scaling applications.
- **Workflows:** Step-by-step procedures for common tasks (rolling updates, blue-green deployments).
- **Security Best Practices:** Authentication, authorization, secrets management, vulnerability scanning.

**4. Advanced Topics (Expert & Operator Persona):**

- **Multi-Cluster Management:** Federation, cluster federation, GitOps for configuration management.
- **Performance Optimization:** Resource usage analysis, profiling tools, tuning techniques.
- **DevOps & CI/CD integration:** Integrating Meshery with continuous integration and deployment pipelines.
- **Custom Resource Definitions:** Extending Meshery with custom APIs and controllers.

**5. Reference & Resources (All Personas):**

- **API Documentation:** Comprehensive reference for Meshery API objects and fields.
- **Release Notes & Changelogs:** Detailed information about version changes, new features, deprecations.
- **External Resources:** Links to blogs, community forums, case studies, books, training materials.
- **Contributing Guide:** How to contribute documentation, code, and other resources to the project.

**Additional Considerations:**

- **Versioning:** Separate documentation for different Meshery versions, clear versioning markings.
- **Search & Navigation:** Robust search functionality, intuitive navigation structure, breadcrumbs.
- **Feedback & Improvement:** Mechanisms for user feedback, reporting errors, suggesting improvements.
- **Community Integration:** Highlight community forums, events, contributor guidelines, recognition.
- **Personalization:** Allow users to personalize their documentation experience (filters, preferences).

This high-level outline provides a comprehensive framework for structuring the Meshery documentation, catering to diverse user needs while ensuring information is readily accessible and actionable. By further refining each section with specific content recommendations and considering the needs of specific personas, the documentation can effectively serve as a valuable resource for everyone interacting with Meshery.

