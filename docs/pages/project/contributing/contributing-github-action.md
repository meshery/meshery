---
layout: page
title: Contributing to Github Actions
permalink: project/contributing/contributing-github-actions
abstract: How to contribute to Meshery Github Actions.
language: en
type: project
category: contributing
list: include
---

This documentation will provide all information about working with github actions.

## Source of Code
[Events that trigger workflows - GitHub Docs](https://docs.github.com/en/actions/writing-workflows/choosing-when-your-workflow-runs/events-that-trigger-workflows#about-events-that-trigger-workflows)
This documentation provides insight into the state of information available to a workflow during execution, such as the commit SHA (`GITHUB_SHA`) and the branch reference (`GITHUB_REF`). These variables allow you to identify the exact commit or branch that triggered the workflow.

### Types of GitHub Refs
1. **Default branch**: This refers to your default branch configured within GitHub, such as `master` or `main` within the.
2. **PR merge branch**: This refers to a PR coming from either the same repository or a forked repository.
3. **PR base branch**: This refers to a a pull request target the PR wanted to it to get merged.

### Understanding `GITHUB_SHA`
1. **Last commit**: This is basically the last commit with respect to the `GitHub Ref`.
2. **Merge commit**: This is a commit after it has been merged into the default branch, which might not always reflect your last commit. This is typically the default behavior for certain workflows.

### Github Event specific to workflow
[Webhook events and payloads - GitHub Docs](https://docs.github.com/en/webhooks/webhook-events-and-payloads)
This documentation will provide you with all the information about the context embedded into specific workflows, such as PR numbers and more.

