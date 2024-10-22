# Meshery Governance

This document defines governance policies for the Meshery project.

Anyone can become a Meshery contributor simply by contributing to the project, with code, documentation or other means. As with all Meshery community members, contributors are expected to follow the [Code of Conduct](./CODE_OF_CONDUCT.md).

## Contributors

GitHub organization members are people who have `triage` access to all repos in the organization. Community members who wish to become members of the organization should meet the following requirements, which are open to the discretion of the steering committee:

## GitHub Project Administration

- Have enabled 2FA on their GitHub account.
- Have joined the Meshery slack.
- Are actively contributing to the project. Examples include:
  - opening issues
  - providing feedback on the project
  - engaging in discussions on issues, pull requests, Slack, etc.
  - attending community meetings
  - have a sponsoring Maintainer, who has agreed to sponsor their membership request

## Maintainership

The Meshery project consists of multiple repositories. Each repository is subject to the same governance model, but different maintainers and reviewers. See [Maintainer Teams](#maintainer teams).

### Who are Maintainers?

The current maintainers of a repo can be found in [MAINTAINERS.md](./MAINTAINERS.md) file of the repo.

Maintainers are individuals that have demonstrated their dedication the betterment of the project and dedication to Maintainers are empowered to review, approve, merge, close issues and pull requests. Maintainers are empowered to make releases. These individuals are most experience with the given project (or specific repo) and are expected to lead its growth and improvement. Adding and removing maintainers for a given repo is the responsibility of the existing maintainer team for that repo and therefore does not require approval from the steering committee.

This privilege is granted with some expectation of responsibility: maintainers are people who care about the Meshery project and want to help it grow and improve.
A maintainer is not just someone who can make changes, but someone who has demonstrated his or her ability to collaborate with the community, support contributors, perform menial tasks - no task is beneath them as it pertains to the betterment of the project and its community. Maintainers uphold standards of process, code convention, architectural guiding principles, and treatment of contributors. Maintainers are the most knowledgeable individuals to review code, contribute high quality code, and follow through to fix issues (in code or tests).

#### Maintainer Teams

A core team of maintainers steward Meshery, however, Meshery is comprised of any number of components (adapters, operators, ...) and any number of logical areas of concern (e.g. community, governance, ...). The size of Meshery's project vision and the pace of its development, requires sub-teams to supplement the core team. Each sub-team is focused on a specific area/component e.g., adapters, docs, website, UI, CLI, and so on. To ensure global coordination and a strong, coherent vision for the project as a whole, each sub-team is led by a member of the core team.

### Becoming a Maintainer

The process of nominating and voting for maintainers involves the following steps:

1. **Nomination:** A community member, typically an existing maintainer, nominates another individual for the role of maintainer. The nomination usually highlights the nominee's contributions to the project and community, and explains why they are a good fit for the role.
2. **Announcement:** The nomination is announced to the entire community, often through an email or a post on a communication platform like Slack or Discord. The announcement includes information about the nominee and their contributions, and outlines the voting process.
3. **Voting Period:** A voting period is set, during which community members can cast their votes. The voting period typically lasts for a week to two weeks, to allow everyone time to consider the nomination and make an informed decision.
4. **Voting:** Community members cast their votes using a designated method, such as replying to an email or reacting to a message. Votes can be either binding or non-binding. Binding votes, cast by existing maintainers, directly impact the outcome of the vote. Non-binding votes, cast by other community members, express their support or concerns, but don't directly determine the result.
5. **Decision:** At the end of the voting period, the votes are tallied. If the nominee receives sufficienct support (a 2/3rds majority from existing maintainers), they are officially granted the role of maintainer. If the nominee does not receive enough support, the nomination is rejected.
6. **Announcement of the Result:** The outcome of the vote is announced to the community. If the nominee is successful, they are welcomed as a new maintainer and granted the associated responsibilities and privileges.

This process ensures that new maintainers are selected in a transparent and democratic way, with input from the entire community. It also helps maintain the health and sustainability of the project by ensuring that those who take on leadership roles are committed to its success and have the support of their peers.

Maintainers will be added to the GitHub @meshery/maintainers team, and made a GitHub maintainer of that team. They will be given write permission to the Meshery GitHub repository <https://github.com/meshery/meshery> repo.

### Maintainer Responsibilities

Maintainership essentially involves more of the same of what you've already doing as a consistently engaged and impactful contributor, including regular attendance to project meetings, good-faith stewardship of the project and it’s external representation, encouragement of new contributors and so forth. Maintainers are generally saddled with the janitorial work that others are either not entrusted to do or don’t want to do (e.g. non-glamorous work like tests, release coordination, meeting prep and hosting, chasing down contributors to complete work).

#### Maintainer Privileges

The aforementioned responsibilities come with the following benefits:

1. Maintainers' voice and vote count on architectural, roadmap, and general decisions.
1. Maintainers' approval on GitHub pull requests counts toward branch restrictions (e.g.the requirement of approvals prior to pull request merge).
1. Maintainers have commit/write access to project repositories, which among other things, greatly enhances your ability to collaborate with contributors on their pull requests as you perform reviews. You are able to checkout their branch, make commits, and push them (or do so conveniently from the github.com web UI by using CMD+G in a single or multi-line comment).
1. Depending upon project level in the CNCF, speaking spots at KubeCon, booth host, webinar and other opportunities.
1. Maintainers have insight to unannounced project vulnerabilities and help coordinate the patch and the public response.
1. Maintainers have authorship rights on project blog.

#### Relinquishment of Role

Maintainers can voluntarily step down from their role at any time.  If you decide to step down, please inform the other maintainers of your decision and explain your reasons.  Before stepping down, please help identify and train your replacement(s), and ensure a seamless transition of responsibilities to prevent any disruption of work. When ready, you or another maintainer will submit a pull request removing your name from the MAINTAINERS.md file and into the list below. Another maintainer will adjust your various permissions / access levels  to various project systems accordingly.

#### Abandonment of Role

Any maintainer who is unresponsive or non-participatory (across the various forms of communication used by the project) for a period of two months forfeits maintainer privileges and resumes a contributor role. While much communication is asynchronous, maintainers are required to attend a minimum of one meeting a month in order facilitate high-fidelity, synchronous communication, which is necessary for the quality collaboration, conveyance of complex topics, and real-time advancement of the project.

##### Emeritus Maintainers

## Reviewers

Everyone is welcome to offer review on open pull requests. Maintainer approval or change requests count toward the requirement of approvals prior to pull request merge.
