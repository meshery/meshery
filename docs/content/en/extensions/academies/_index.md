---
title: Academies
description: Transform Meshery into an integrated learning platform.
type: extension
---
Academy extensions enable Meshery as an integrated learning platform. Academy extensions provide a hands-on, structured learning experience for beginners entering the cloud native ecosystem and advanced practitioners seeking deeper expertise. By embedding Meshery [designs](/concepts/logical/designs) directly into courses, learners gain immediate, interactive practice with real Meshery-powered workflows—turning passive study into active, contextual application.

### Learning Paths, Challenges, and Certifications

The platform organizes educational material into three primary formats:

- **Learning Paths**: Guided sequences of courses that break down complex topics (e.g., Mastering Meshery, Kubernetes fundamentals, service mesh patterns) into manageable modules.
- **Challenges**: Practical, scenario-based exercises where learners solve specific problems, often involving deploying or configuring infrastructure with Meshery.
- **Certifications**: Formal assessments that validate skills, such as the [Certified Meshery Contributor](https://meshery.io/community/certifications) (CMC) program, complete with badges and shareable credentials.

### Content Authoring

Built on a modular Hugo-based static site framework (with a multi-repository architecture separating theme, content, and build concerns), academies is highly extensible. Organizations and contributors can:

- Plug in custom academies or content repositories tailored to specific technologies, teams, or use cases.
- Embed live Meshery designs, visualizations, and interactions for immersive labs.
- Use Git-native workflows for content creation—manage everything in Markdown with Hugo shortcodes, custom CSS, and HTML mixing, while leveraging version control, pull requests, and automated builds/releases.

This architecture supports multi-tenancy, white-labeling (for branded experiences), and seamless integration with the broader Meshery ecosystem.

### Roles Supported

- **Learners**: Progress through paths, tackle challenges, earn certifications, and interact with Meshery designs.
- **Content Creators/Instructors**: Build and publish courses using the intuitive, Git-based tools and instructor console.
- **Platform Developers/Contributors**: Extend the academy framework, theme, integrations, and Meshery-specific features via the open source repositories (e.g., meshery-academy, academy-theme).

Meshery Academies exemplifies Meshery’s philosophy of extensibility—empowering the community to democratize cloud native knowledge through practical, visual, and interactive learning experiences. It serves as both an official learning hub for Meshery (with paths like “Mastering Meshery”) and a framework for anyone to create their own specialized academies.

To find a complete list of academies available, explore the repositories under the [meshery-academy topic](https://github.com/topics/meshery-academy), like [meshery-academy](https://github.com/meshery-extensions/meshery-academy). Contributions and extensions are welcome in the Meshery Extensions organization.
