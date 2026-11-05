# Security Policy

## Reporting a Vulnerability

We are very grateful to the security researchers and users that report back Meshery security vulnerabilities. We investigate every report thoroughly.

To report a vulnerability and for details of published vulnerabilities, see [security vulnerabilities](https://docs.meshery.io/project/security-vulnerabilities) in Meshery Docs.

## Supported Versions

See [release cadence](https://docs.meshery.io/project/contributing/build-and-release#release-cadence) in Meshery Docs for the support schedule.

## Privileged Workflows

The following workflows run in a privileged context (`pull_request_target`, `workflow_run`, or requiring repository write access) to perform automated triage, labeling, documentation previews, or testing notifications:

| Workflow File | Purpose | Security Controls & Isolation Boundary |
| :--- | :--- | :--- |
| [`auto-merge-meeting-minutes.yml`](https://github.com/meshery/meshery/blob/master/.github/workflows/auto-merge-meeting-minutes.yml) | Automates merging of scheduled community notes. | Scoped strictly to bot-generated meeting minute files; does not execute code from fork PRs. |
| [`comment-e2e-test.yml`](https://github.com/meshery/meshery/blob/master/.github/workflows/comment-e2e-test.yml) | Posts end-to-end test execution status to PR threads. | Operates solely on completed run metadata; does not check out or execute code from untrusted forks. |
| [`copilot-pr-handler.yml`](https://github.com/meshery/meshery/blob/master/.github/workflows/copilot-pr-handler.yml) | Triage and context handler for PR contributions. | Performs metadata and issue operations; uses no untrusted build steps. |
| [`deploy-docs-preview.yml`](https://github.com/meshery/meshery/blob/master/.github/workflows/deploy-docs-preview.yml) | Generates isolated documentation site previews. | Evaluates doc changes; executes inside isolated GitHub-hosted environments without persistent cache pollution risks. |
| [`image-scan.yml`](https://github.com/meshery/meshery/blob/master/.github/workflows/image-scan.yml) | Runs container security & vulnerability scans. | Scans published image targets; does not mount elevated credentials during inspection of untrusted layers. |
| [`label-commenter.yml`](https://github.com/meshery/meshery/blob/master/.github/workflows/label-commenter.yml) | Comments on issues/PRs based on triage labels. | Pure metadata automation; does not perform code checkouts. |
| [`labeler.yml`](https://github.com/meshery/meshery/blob/master/.github/workflows/labeler.yml) | Applies file-based labels on incoming contributions. | Native triage action using GitHub API; does not checkout or build untrusted code. |
| [`welcome-first-time-contributors.yml`](https://github.com/meshery/meshery/blob/master/.github/workflows/welcome-first-time-contributors.yml) | Welcomes new community contributors. | Emits standard welcome message via GitHub API; does not access repository source files or credentials. |

### Isolation Policies & Review Scope
- **Checkout Restrictions:** Automated AST analysis confirms privileged workflows do not check out untrusted pull request code heads (`pull_request_target` checkouts are restricted to base branches or omitted).
- **Artifact Handling:** Workflows consuming artifacts from upstream pull request runs are restricted to non-executable data parsing to mitigate artifact poisoning.
- **Privilege Separation:** Automated triage and commenting tasks are isolated from deployment, release, and signing pipelines.
