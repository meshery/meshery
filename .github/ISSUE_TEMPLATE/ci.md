---
name: 🛠 Continuous Integration / DevOps
about: Improve or update workflows or other automation
title: '[CI]'
labels: 'area/ci'
assignees: ''
---
#### Current Behavior
<!-- A brief description of what the problem is. (e.g. I need to be able to...) -->

#### Desired Behavior
<!-- A brief description of what you expected to happen. -->

#### Implementation
<!-- Specifics on the approach to fulfilling the feature request. -->

#### Acceptance Tests
<!-- Stipulations of functional behavior or non-functional items that must be in-place in order for the issue to be closed. -->

---
#### Contributor [Guides](https://docs.meshery.io/project/contributing) and [Handbook](https://layer5.io/community/handbook)
- 🛠 [Meshery Build & Release Strategy](https://docs.meshery.io/project/contributing/build-and-release)
- 📚 [Instructions for contributing to documentation](https://github.com/meshery/meshery/blob/master/CONTRIBUTING.md#documentation-contribution-flow)
- Meshery documentation [site](https://docs.meshery.io/) and [source](https://github.com/meshery/meshery/tree/master/docs)
- 🎨 Wireframes and designs for Meshery UI in [Figma](https://www.figma.com/file/SMP3zxOjZztdOLtgN4dS2W/Meshery-UI)
- 🙋🏾🙋🏼 Questions: [Discussion Forum](https://discuss.layer5.io) and [Community Slack](http://slack.layer5.io)
- 📚[Meshery Test Plan - v0.7.0 sheet](https://docs.google.com/spreadsheets/d/13Ir4gfaKoAX9r8qYjAFFl_U9ntke4X5ndREY1T7bnVs/edit#gid=0)
- 📚 [Meshery Test Plan - Automated Tests/Builds sheet ](https://docs.google.com/spreadsheets/d/13Ir4gfaKoAX9r8qYjAFFl_U9ntke4X5ndREY1T7bnVs/edit#gid=1214153102)


---
### Contributors are encouraged to follow these guidelines to maintain an up-to-date and accurate Test plan sheet for test scenarios.

- Adding New Rows:
  - Create a new row in the v0.7.0 sheet when a new test scenario is implemented or merged into the **'master'** branch.
    Include the combination of Test Group, Component Under Test, Platform, Service Mesh, Action, and Expected Outcome in the new row.
    Specify if the new test scenario is related to recently started existing features or upcoming new features/enhancements.

- Updating Existing Rows:
  - Modify an existing row in the v0.7.0 sheet if a test scenario has changed in a significant or relevant way.
    Update the relevant fields (such as Component) in the existing row to reflect the changes.
    Updating **"Actual Outcome"** and/or **"Additional Remarks"**:

When using Meshery UI, Mesheryctl CLI, Meshery Docker Extension, specific Meshery adapters, or the Meshmap extension, update the "Actual Outcome" and/or "Additional Remarks" fields if any issues, abnormal behavior, or unexpected behavior are encountered.
Update the fields based on the user's experience with a recent release or when building from source.


