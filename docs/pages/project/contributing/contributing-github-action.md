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

This documentation covers everything you need to know about working with GitHub Actions.
## Commit and Branch
For a detailed explanation of events that trigger workflows, see the [GitHub Docs](https://docs.github.com/en/actions/writing-workflows/choosing-when-your-workflow-runs/events-that-trigger-workflows#about-events-that-trigger-workflows).

This documentation provides insight into the state of information available to a workflow during execution, such as the commit SHA (`GITHUB_SHA`) and the branch reference (`GITHUB_REF`). These variables allow you to identify the exact commit or branch that triggered the workflow.

### Types of GitHub Refs
1. **Default branch**: This refers to your default branch configured within GitHub, such as `master` or `main` within the.
2. **PR merge branch**: Refers to a pull request coming from either the same repository or a forked repository.
3. **PR base branch**: The target branch for a pull request that the PR aims to be merged into.

Developers often build from pull requests, whether they originate from the same repository or a fork. There are two common workflows you can use:
1. pull_request_target
	This workflow is popular because it has direct access to secrets. However, it always checks out from the targeted base branch, meaning it does not check code coming from the PR branch. For example, if your action needs an event like adding a comment (which requires a PR number), keep in mind that the event will be updated to the base branch. If it can’t find your PR live on the base branch, it won't do anything.
1. pull_request
	 This workflow checks out from your PR merge branch but does not have access to secrets.

### Understanding GITHUB_SHA
1. **Last commit**: This is basically the last commit head with respect to the `GitHub Ref`.
2. **Merge commit**: A hypothetical merge with the base branch, which might not always reflect your last commit. This is typically the default behavior for certain workflows.
To check out directly from your PR head branch, specify the [head SHA](https://github.com/actions/checkout?tab=readme-ov-file#checkout-pull-request-head-commit-instead-of-merge-commit) of your commit.

### Fork Repo
When working with a forked repository, the pull request you create will come from your own fork. This means GitHub, for security reasons, uses the context of the base branch and might be unable to find your PR because it lives on your fork branch.

### Secrets
GitHub secrets will be unavailable if you run the workflow outside of the base branch, especially if your pull request is coming from a forked repository. This is to protect against potential malicious activity.

### Github Events and Payloads
 For a comprehensive understanding of GitHub events and payloads, see the [Webhook events and payloads documentation](https://docs.github.com/en/webhooks/webhook-events-and-payloads).  
 For certain workflow trigger between forking and base branch github will start dropping certain event and payload which can be unpredictable, if this the case you can utilize github artifact for passing certain events