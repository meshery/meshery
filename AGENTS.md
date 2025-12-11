# GitHub Copilot Agents for Meshery

This repository utilizes specialized GitHub Copilot agents to assist contributors with various aspects of the Meshery project. These agents are custom-configured AI assistants designed to understand Meshery's architecture, conventions, and best practices.

## What are GitHub Copilot Agents?

GitHub Copilot agents are specialized AI assistants configured with domain-specific knowledge and instructions. Each agent in this repository is tailored to help with specific types of contributions to Meshery, from code development to documentation writing and CI/CD workflow management.

These agents provide:
- **Context-aware assistance**: Understanding Meshery's tech stack, architecture, and conventions
- **Consistent guidance**: Following established patterns and best practices
- **Specialized expertise**: Domain-specific knowledge for different areas of the project
- **Accelerated onboarding**: Helping new contributors understand project structure and workflows

## Available Agents

### 1. Meshery Code Contributor

**Location**: `.github/agents/meshery-code-contributor.md`

**Purpose**: Expert-level software engineering agent specialized in contributing to Meshery's cloud native infrastructure and application management platform.

**Capabilities**:
- Contributing to Meshery Server (Go/Golang backend)
- Developing Meshery UI (React/Next.js frontend)
- Building mesheryctl CLI (Go/Golang with Cobra framework)
- Understanding Meshery's design principles and architecture
- Following code style conventions and best practices
- Implementing schema-driven development patterns

**When to use**: When working on any code contributions to Meshery, including bug fixes, new features, or refactoring existing functionality.

### 2. Meshery Docs Contributor

**Location**: `.github/agents/meshery-docs-contributor.md`

**Purpose**: Specialized agent for creating, updating, and maintaining Meshery's documentation.

**Capabilities**:
- Writing and editing documentation in Jekyll/Markdown format
- Understanding Meshery's documentation structure and navigation
- Following documentation style guide and conventions
- Creating tutorials, guides, and reference documentation
- Ensuring consistency across documentation
- Maintaining links and cross-references

**When to use**: When creating new documentation, updating existing docs, fixing documentation issues, or improving clarity of explanations.

### 3. GitHub Actions Engineer

**Location**: `.github/agents/github-actions-engineer.md`

**Purpose**: Expert in GitHub Actions, CI/CD pipelines, and workflow automation.

**Capabilities**:
- Designing and debugging GitHub Actions workflows
- Implementing cross-repository orchestration
- Managing secrets and permissions securely
- Optimizing workflow performance
- Setting up automated testing and deployment
- Implementing failure handling and notifications

**When to use**: When working on CI/CD workflows, GitHub Actions, build automation, or deployment pipelines.

## Using the Agents

### Prerequisites

To use these agents, you need:
- Access to GitHub Copilot (GitHub Copilot Chat or Copilot Workspace)
- A development environment with the Meshery repository cloned
- Familiarity with basic Meshery concepts (see [README.md](./README.md))

### How to Interact with Agents

1. **In GitHub Copilot Chat**:
   - Reference the agent by mentioning its context: "Using the Meshery Code Contributor agent guidelines..."
   - Provide specific context about what you're working on
   - Ask for guidance, code examples, or reviews

2. **In Pull Requests**:
   - GitHub Copilot may automatically use these agent definitions when reviewing or suggesting changes
   - The agents help ensure contributions align with Meshery's standards

3. **During Development**:
   - Keep the agent guidelines in mind while coding
   - Use them as a reference for best practices
   - Consult them when unsure about conventions

### Best Practices

- **Be specific**: Provide clear context about what you're trying to accomplish
- **Reference existing code**: Point to similar implementations when asking for guidance
- **Iterate**: Work with the agent iteratively, refining solutions based on feedback
- **Validate**: Always validate agent suggestions against project standards and test thoroughly

## Agent Structure

Each agent definition file follows a consistent structure:

```markdown
---
name: Agent Name
description: Brief description of the agent's purpose
tools: [list of tools the agent can use]
---

# Detailed Instructions

[Comprehensive instructions including:
- Role and responsibilities
- Technical expertise areas
- Code style and conventions
- Best practices
- Common patterns
- Quality standards
- Resources and references]
```

## Contributing to Agents

### Improving Existing Agents

If you notice gaps, inaccuracies, or areas for improvement in existing agents:

1. Review the agent definition file in `.github/agents/`
2. Make your proposed changes
3. Test the changes by using the modified agent for actual tasks
4. Submit a pull request with:
   - Clear description of the changes
   - Rationale for the improvements
   - Examples of how the changes improve agent performance

### Creating New Agents

To create a new specialized agent for Meshery:

1. **Identify the need**: Determine if a new agent is necessary or if an existing agent should be expanded
2. **Define the scope**: Clearly outline the agent's area of expertise
3. **Structure the agent**: Follow the existing agent template structure
4. **Document capabilities**: List specific skills, knowledge areas, and use cases
5. **Include examples**: Provide concrete examples of the agent's expected behavior
6. **Test thoroughly**: Use the agent for real tasks before submitting
7. **Submit for review**: Create a pull request to `.github/agents/`

### Agent Maintenance

Agents should be updated when:
- Project conventions or best practices change
- New technologies or frameworks are adopted
- Common issues or gaps are identified
- The tech stack or architecture evolves

## Relationship to Copilot Instructions

The `.github/copilot-instructions.md` file provides general instructions for GitHub Copilot across the entire Meshery repository. It covers:
- Project overview and structure
- Tech stack information
- Build and test commands
- Common development tasks
- Links to documentation

The specialized agents in `.github/agents/` build upon this foundation with deeper, domain-specific expertise for particular types of contributions.

## Additional Resources

- **Contributing Guide**: [CONTRIBUTING.md](./CONTRIBUTING.md)
- **Meshery Documentation**: https://docs.meshery.io/project/contributing
- **Community Slack**: https://slack.meshery.io
- **Code of Conduct**: [CODE_OF_CONDUCT.md](./CODE_OF_CONDUCT.md)
- **Governance**: [GOVERNANCE.md](./GOVERNANCE.md)

## Questions or Feedback?

If you have questions about using the agents or suggestions for improvements:
- Open an issue in this repository
- Ask in the [Meshery Slack community](https://slack.meshery.io)
- Bring it up in community meetings

---

**Note**: These agents are continuously evolving. Your feedback and contributions help make them more effective for the entire Meshery community.
