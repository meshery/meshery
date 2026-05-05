---
name: First Time Contributor
description: Welcomes first-time contributors to Meshery by posting a warm welcome comment and adding the first-time-contributor label to their pull request
on:
  pull_request_target:
    types: [opened]
  roles: all
if: github.repository == 'meshery/meshery' && github.event.pull_request.author_association == 'FIRST_TIME_CONTRIBUTOR'
permissions:
  contents: read
  pull-requests: read
engine: copilot
timeout-minutes: 10
tools:
  github:
    toolsets: [context, pull_requests]
safe-outputs:
  add-comment:
    max: 1
    pull-requests: true
    issues: false
    discussions: false
  add-labels:
    allowed: [first-time-contributor]
    max: 1
---

# Welcome First-Time Contributor

Pull request #${{ github.event.pull_request.number }} was just opened by a first-time contributor to this repository.

Complete both tasks below.

## Task 1: Post a welcome comment

1. Fetch pull request #${{ github.event.pull_request.number }} to get the author's GitHub login.
2. Post a comment on the pull request with the following content, replacing `{username}` with the author's login:

Welcome, @{username}! Thank you for your first contribution! 🎉 A contributor will be by to give feedback soon. In the meantime, please review the [Newcomers' Guide](https://meshery.io/community) and be sure to join the [community Slack](https://slack.meshery.io/).

<p align="center" width="100%">
<img src="https://github.com/user-attachments/assets/ba4699dc-18b2-4884-9dce-36ed47c38e93" width="30%" />
</p>

Be sure to double-check that you have signed your commits. Here are instructions for [making signing an implicit activity while performing a commit](https://docs.meshery.io/project/contributing#general-contribution-flow).

## Task 2: Add a label

Add the `first-time-contributor` label to pull request #${{ github.event.pull_request.number }}.
