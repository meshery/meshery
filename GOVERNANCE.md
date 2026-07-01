# Meshery Governance

This document defines the governance of the Meshery project. It describes the
project's values, its structure across the two GitHub organizations that make up
the Meshery ecosystem, the roles that people hold, how those roles are earned and
relinquished, how decisions are made, and the systems and permissions that back
each role. It is an open, living document that the project iterates on as the
community and the project change.

Meshery is a [Cloud Native Computing Foundation (CNCF)](https://www.cncf.io)
project. Where this document does not address a matter, the project defers to the
[CNCF Charter](https://github.com/cncf/foundation/blob/main/charter.md) and the
[CNCF Technical Oversight Committee (TOC)](https://github.com/cncf/toc).

## Contents

- [Values](#values)
- [Vendor Neutrality](#vendor-neutrality)
- [Project Structure](#project-structure)
- [Roles and the Contributor Ladder](#roles-and-the-contributor-ladder)
- [The Maintainer Council](#the-maintainer-council)
- [Decision-Making](#decision-making)
- [Maintainer Lifecycle](#maintainer-lifecycle)
- [Subproject Lifecycle](#subproject-lifecycle)
- [Systems Access and Repository Permissions](#systems-access-and-repository-permissions)
- [Communication Channels](#communication-channels)
- [Meetings](#meetings)
- [Code of Conduct](#code-of-conduct)
- [Security Response](#security-response)
- [Requesting CNCF Resources](#requesting-cncf-resources)
- [How to Contribute](#how-to-contribute)
- [Changes to This Governance](#changes-to-this-governance)
- [Governance Lineage](#governance-lineage)
- [Appendix A: Mapping to CNCF Incubation Governance Review Criteria](#appendix-a-mapping-to-cncf-incubation-governance-review-criteria)
- [Appendix B: MAINTAINERS.md Structure](#appendix-b-maintainersmd-structure)

## Values

The Meshery project and its leadership embrace the following values. These values
apply to everything the project does, and every person in a leadership role is
expected to uphold them.

- **Openness.** Work happens in the open. Discussions, decisions, code, and
  documentation are public by default, on neutral platforms, so that anyone can
  follow along and participate.
- **Fairness and meritocracy.** Influence is earned through sustained, good-faith
  contribution rather than employer, title, or tenure. Every leadership role has a
  visible, attainable path, and contributions are evaluated on their merits.
- **Vendor neutrality.** No single company controls the project's direction,
  features, or roadmap. See [Vendor Neutrality](#vendor-neutrality).
- **Inclusivity.** The project welcomes contributors of every background and
  experience level, and it recognizes both code and non-code contributions as
  first-class.
- **Extensibility.** Meshery is built to be extended. The project favors clear
  interfaces and extension points over monolithic design, both in its software and
  in how the community organizes around it.
- **Stewardship.** Maintainers are stewards, not owners. They are responsible for
  the long-term health of the project and the success of its contributors, and they
  take on unglamorous work in service of both.
- **End-user focus.** The project exists to make managing cloud native
  infrastructure simpler and more reliable for the people who operate it.

## Vendor Neutrality

Meshery is a vendor-neutral project.

- Project blogs and announcements do not favor or advertise any specific vendor,
  and participation rules across Slack and social media are applied equally to all.
- Community meetings, code, and documentation are hosted on neutral platforms such
  as CNCF Zoom and GitHub.
- No single vendor controls the project's features or roadmap. Architectural
  decisions are made openly with equal opportunity for community input, and no
  vendor may block reasonable contributions from other organizations.
- No single organization holds sole decision-making authority. Leadership is open
  to any contributor regardless of employer or affiliation.
- The project tracks the affiliation of its maintainers in
  [MAINTAINERS.md](./MAINTAINERS.md) so that organizational concentration is visible
  and can be managed.
- Unresolved conflicts, including any concern that the project's neutrality is
  compromised, may be escalated to the [CNCF TOC](https://github.com/cncf/toc).

## Project Structure

### The Meshery ecosystem and its two GitHub organizations

To serve a large and growing ecosystem, the Meshery project is partitioned into two
GitHub organizations:

- **[github.com/meshery](https://github.com/meshery)** hosts the **core platform**,
  including critical components such as
  [Meshery Server](https://github.com/meshery/meshery), Meshery UI,
  [mesheryctl](https://docs.meshery.io/installation/mesheryctl), the
  [Meshery Operator](https://docs.meshery.io/concepts/architecture/operator),
  [MeshSync](https://docs.meshery.io/concepts/architecture/meshsync),
  [schemas](https://github.com/meshery/schemas), policies, and the documentation
  site.
- **[github.com/meshery-extensions](https://github.com/meshery-extensions)** hosts
  community-developed **extensions and integrations**, such as adapters, providers,
  and tooling that build on the core platform's extension points.

This separation lets the core maintainers keep the core platform robust and focused while
enabling broad community participation around it. It mirrors the approach taken by
graduated CNCF projects, including Kubernetes (`kubernetes` and `kubernetes-sigs`),
Crossplane (`crossplane` and `crossplane-contrib`), and Argo (`argoproj` and
`argoproj-labs`).

Throughout this document, **"the project"** refers to the entire Meshery ecosystem
across both organizations.

### Core platform (`github.com/meshery`)

The core platform is governed by the Meshery core maintainers under this document. It
provides full support: regular updates, bug fixes, security response, and
comprehensive documentation. Changes follow a structured release cycle with
[stable and edge channels](https://docs.meshery.io/project/contributing/build-and-release)
and undergo rigorous review to sustain stability.

The core platform is organized into **subprojects** (also called domains), each
focused on a component or an area of concern. See
[Subprojects](#subprojects-and-domains) below.

### Extensions (`github.com/meshery-extensions`)

Extensions allow the ecosystem to scale beyond what the core maintainers can directly support, and they act as an incubator for new ideas that may or may not later migrate toward the core. Extensions can be initiated, developed, and maintained by members of the community rather than by the core maintainers.

Extensions operate under a lighter governance structure to encourage innovation:

- **Maintainers.** Each extension has its own extensions maintainers. They may be
  nominated by community members or self-nominated, and are confirmed by the core
  maintainers based on contribution history and technical expertise.
- **Autonomy.** Extension teams have autonomy over their own development processes
  and release cadence, provided they adhere to the project's
  [Code of Conduct](./CODE_OF_CONDUCT.md), to the core component frameworks, and to
  the project's [integration guidelines](https://docs.meshery.io/extensibility).
- **Support labeling.** Each extension's documentation indicates its support level
  so that users understand what to expect:
  - **Official.** Maintained by the core maintainers or designated maintainers, with robust
    support and compatibility testing against supported core releases.
  - **Community.** Maintained by community contributors, with support and release
    cadence determined by those maintainers.
- **Onboarding new extensions.** New extensions are admitted through the process in
  [Adding a subproject or extension](#adding-a-subproject-or-extension).

The following table summarizes the difference in governance between the two
organizations.

| Aspect | Core Platform | Extensions |
| :--- | :--- | :--- |
| Governance | Structured, led by the core maintainers (the Maintainer Council) under this document | Flexible, led by per-extension maintainers under lighter rules |
| Maintainer selection | Nomination and a 2/3 majority vote of existing core maintainers | Nomination or self-nomination, confirmed by the core maintainers |
| Decision-making | Lazy consensus, with voting where required | Extension-team consensus, with core oversight |
| Support | Full support and documentation | Variable; labeled Official or Community |
| Communication | Public meetings, Slack, forums | Public issues, Slack, optional meetings |

### Subprojects and domains

Within the core platform, work is organized into **subprojects**. A subproject is a
technically distinct area that has its own maintainers and its own day-to-day
operations, while remaining aligned with the overall project. Depending on context,
the project also refers to these as **domains** or, following the Kubernetes model
that inspired this structure, as the project's equivalent of Special Interest
Groups.

Representative subprojects include, but are not limited to:

- **Server** - the Go backend (REST and GraphQL APIs, Kubernetes management, adapter
  orchestration, persistence).
- **UI** - the Next.js and React user interface.
- **CLI (mesheryctl)** - the command-line interface.
- **Schemas** - the [schema definitions](https://github.com/meshery/schemas) shared
  across the ecosystem.
- **MeshSync and Operator** - cluster discovery, synchronization, and lifecycle.
- **CI/CD** - [project quality](https://qa.meshery.io), build workflows, release process, deployment and operation of project playground.
- **Documentation** - the documentation site, guides, and references.
- **Community and Governance** - community programs, events, onboarding, recognition,
  and the maintenance of this document and related governance artifacts.

Each subproject is governed by its maintainers as a maintainer team. The
authoritative, current list of subprojects, the repositories and directories each
owns, and the maintainers responsible for each is recorded in
[MAINTAINERS.md](./MAINTAINERS.md) (see
[Appendix B](#appendix-b-maintainersmd-structure)). Each subproject
designates one of its maintainers as a [Subproject Lead](#subproject--domain-lead).

To ensure global coordination and a coherent vision for the project as a whole, each
subproject and its lead participate in the [Maintainer Council](#the-maintainer-council).

## Roles and the Contributor Ladder

The project defines a ladder of roles. Anyone can enter at the bottom, and each rung
has a clear set of requirements, responsibilities, and privileges. Movement up the
ladder is based on demonstrated, sustained contribution. Both **code** and
**non-code** contributions count toward advancement, and non-code leadership roles
such as Community Manager and MeshMate are recognized first-class roles.

The rungs, in ascending order of responsibility, are: Community Participant,
Contributor, Organization Member, Maintainer, and Subproject Lead. Maintainers fall
into two classes that reflect the project's two organizations: **core maintainers**,
who maintain the core platform in `github.com/meshery` and collectively form the
project's governing body, and **extensions maintainers**, who maintain extensions in `github.com/meshery-extensions` under a
lighter governance. Community Manager and MeshMate are non-code leadership
distinctions that overlay this ladder.

### Community Participant

Anyone who uses Meshery, asks questions, files issues, joins the Slack or the
discussion forum, or attends a community meeting is a community participant. No
special permissions are required, and participation requires only agreement to the
[Code of Conduct](./CODE_OF_CONDUCT.md). Community participants have read access to
the project's public repositories.

### Contributor

A Contributor is anyone who contributes to the project, with code, documentation,
design, community work, or other means.

**How to become a Contributor.** Open a pull request, file or triage an issue,
improve documentation, help others in Slack or the forum, or otherwise contribute,
and follow the steps in [CONTRIBUTING.md](./CONTRIBUTING.md). There is no application
and no waiting period: making an accepted contribution makes you a Contributor.
First-time contributors are encouraged to start with issues labeled
[help wanted](https://github.com/meshery/meshery/labels/help%20wanted) and to find a
[MeshMate](#meshmates) for guidance.

**What a Contributor can do.**

- Submit pull requests from a fork and have them reviewed.
- Open, comment on, and participate in issues, pull requests, and discussions.
- Have their work merged once it meets the project's standards and receives the
  required approvals.
- Be credited in the project's contributor list and through the
  [recognition program](#recognition-and-certification).

A current list of contributors is published at
[the project's contributor list](https://github.com/meshery/meshery/graphs/contributors)
and celebrated through [meshery.io/community/members](https://meshery.io/community/members).

**Reviewing work in progress.** Everyone is welcome to review and comment on open
pull requests and issues. The project does not maintain a separate, gated "Reviewer"
role: review and feedback are open to all. The authority to approve a change for
merge is held by the maintainers responsible for the affected code or documentation,
as recorded in [MAINTAINERS.md](./MAINTAINERS.md) (see
[Systems Access and Repository Permissions](#systems-access-and-repository-permissions)).

### Organization Member

Organization Members have `triage` access across the repositories in the
organization. Triage access allows a person to manage issues and pull requests
(labeling, assigning, requesting reviewers, and closing or reopening) without write
access to code.

**How to become an Organization Member.** A Contributor may request organization
membership once they meet the following requirements, which are evaluated by the
maintainers of the area in which the person is active and, where needed, by the
[core maintainers](#the-maintainer-council):

- Have enabled two-factor authentication (2FA) on their GitHub account.
- Have joined the [Meshery Slack](https://slack.meshery.io).
- Are actively contributing to the project. Examples include opening issues,
  providing feedback, engaging in discussions on issues, pull requests, and Slack,
  and attending community meetings.
- Have a sponsoring Maintainer who agrees to sponsor the membership request.

To request membership, open an issue in the relevant repository stating that you
meet the requirements and naming your sponsoring Maintainer. A Maintainer adds
approved members to the appropriate GitHub team.

**What an Organization Member can do.** In addition to everything a Contributor can
do, Organization Members can triage issues and pull requests across the
organization, which materially helps the project keep its backlog healthy.

### Maintainer

Maintainers are individuals who have demonstrated dedication to the betterment of the
project and the sustained good judgment to steward an area of it. A Maintainer is
both a Project Maintainer for CNCF purposes and an approver for the part of the
project they maintain.

A Maintainer is not merely someone who can make changes. A Maintainer is someone who
has shown the ability to collaborate with the community, support contributors, bring
the most knowledgeable people in to review code and documentation, contribute
high-quality work, and follow through to fix issues in code or tests. Maintainers
uphold the project's standards of process, code convention, architectural principles,
and respectful treatment of contributors. No task is beneath a Maintainer when it
serves the project and its community.

Maintainership is held per subproject, and Maintainers fall into two classes that
correspond to the project's two organizations.

#### Core Maintainers

Core maintainers maintain the core platform in
[github.com/meshery](https://github.com/meshery). They are the project's CNCF Project
Maintainers, and collectively they form the
[Maintainer Council](#the-maintainer-council), the project's governing body. A core
maintainer maintains one or more core subprojects (Server, UI, CLI, Schemas, MeshSync
and Operator, Documentation, Community and Governance, and so on). The core
maintainers of a given subproject are responsible for adding and removing the core
maintainers of that subproject, which does not require approval from the Maintainer
Council, and for designating that subproject's [Lead](#subproject--domain-lead).
Beyond their own subproject, core maintainers carry project-level governance
responsibility, including the oversight of extensions and the decisions described
throughout this document.

#### Extensions Maintainers

Extensions maintainers, maintain individual
extensions in [github.com/meshery-extensions](https://github.com/meshery-extensions).
They have autonomy over their extension's development and release cadence under the
lighter governance described in
[Extensions](#extensions-githubcommeshery-extensions), provided they adhere to the
[Code of Conduct](./CODE_OF_CONDUCT.md), the core component frameworks, and the
integration guidelines. Extensions maintainers are confirmed by the core maintainers,
they are not members of the Maintainer Council, and their authority is scoped to the
extensions they maintain. Sustained, high-impact work as an extensions maintainer is
one of the paths toward becoming a core maintainer.

The complete, current list of maintainers, the class of each, the repositories and
directories each is responsible for, and each maintainer's contact information and
affiliation, is maintained in [MAINTAINERS.md](./MAINTAINERS.md). See the
[Maintainer Lifecycle](#maintainer-lifecycle) for how these roles are earned,
exercised, and relinquished, and
[Appendix B](#appendix-b-maintainersmd-structure) for the structure of
`MAINTAINERS.md`.

### Subproject / Domain Lead

Each subproject designates one of its Maintainers as its Lead. The Lead is a
Maintainer who additionally:

- Coordinates the subproject's roadmap and release work and ensures the subproject
  stays aligned with the overall project.
- Represents the subproject in the [Maintainer Council](#the-maintainer-council) and
  brings cross-subproject matters to it.
- Serves as the first point of contact for questions about the subproject.

A subproject's Maintainers select and may replace their Lead by simple majority. The
current Lead of each subproject is recorded in [MAINTAINERS.md](./MAINTAINERS.md).

### Community Managers

Community Managers are responsible for the health and growth of the community. A
Community Manager serves as a link between the project and its community, working
across teams to resolve obstacles before and as they arise. Their work includes:

- Cultivating an environment that attracts and retains members, in part by ensuring
  questions receive timely and complete responses.
- Aspects of project management, such as organizing meetings and triaging issues.
- Aspects of community marketing, such as member and project promotion.
- Helping existing members stay engaged on an ongoing basis.

Community Manager is a non-code leadership role. Sustained, high-impact community
work is a recognized path up the contributor ladder, including toward Maintainership
of the Community and Governance subproject. 

### MeshMates

MeshMates are experienced community members who help others become successful
contributors. MeshMates help newcomers identify areas of the project to engage with,
working groups to join, and ways to grow their open source and cloud native
knowledge, often by connecting one-on-one and sharing tips for the best community
experience.

MeshMate is a distinction awarded to select members of the community who embody the
project's culture of helping others, paying it forward, and sharing knowledge.
MeshMates are community ambassadors, not employees. As with other roles, the project
recognizes Emeritus MeshMates who have previously served with distinction. Learn more
and find a MeshMate at [meshery.io/community/meshmates](https://meshery.io/community/meshmates).

### Recognition and Certification

The project recognizes contributions, both code and non-code, through a public
[recognition program](https://meshery.io/community/recognition) that awards badges for milestones, skills, and roles. 
Recognition is one of the ways the project makes the contributor ladder visible and celebrates
progress along it. Badges fall into several categories: achievement badges for user
milestones such as a first deployment, project badges for contribution milestones,
community badges for community and no-code contributions such as the MeshMate badge,
and certification badges earned by passing a Meshery certification such as the
[Certified Meshery Contributor](#certified-meshery-contributor-cmc).

#### Certified Meshery Contributor (CMC)

The CMC is the first certification in the broader [Meshery Certification Program](https://meshery.io/community/certifications), a
structured pathway for validating skills across the ecosystem. The program is
organized into a Developers track (Certified Meshery Contributor and Certified Meshery
Developer) and an Administrators track (Certified Meshery Associate, Certified Meshery
Professional, and Certified Meshery Expert). It complements the project's
[recognition program](#recognition-and-certification) by adding a certification badge
category.

The Certified Meshery Contributor (CMC) is the project's contributor certification and
the first contributor certification of its kind in the CNCF. It validates a
contributor's technical proficiency across Meshery's major architectural domains
through written assessments, giving contributors a way to demonstrate and validate
their expertise as part of the project's onboarding and recognition experience.

The CMC comprises five exams, one for each of Meshery's major architectural domains:
Meshery Server, Meshery CLI, Meshery UI, Meshery Models, and Meshery Extensibility.
The exams are created and reviewed by the maintainers of each domain, which keeps the
certification aligned with how the project actually works. It is aimed at developers
with intermediate proficiency in Go, React, and OpenAPI schemas and hands-on
experience in the codebase, and it is open to developers, technical writers, and
community members. The exam is a free, online, multiple-choice assessment, two hours
in duration, with a 70% passing score and a two-year validity. No minimum number of
contributions is required to sit for it. See the [recognition program](https://meshery.io/community/recognition).

## The Maintainer Council

The **Maintainer Council** is the top-level governing body of the Meshery project.
It is the collective body of the project's **core maintainers**, with each core
subproject represented through its core maintainers and its
[Lead](#subproject--domain-lead). The members of the Maintainer Council are the
project's owners as far as the CNCF is concerned. Extensions maintainers are not 
members of the Maintainer Council; they govern their extensions
under lighter rules and coordinate with the core maintainers.

Meshery deliberately governs itself through a Maintainer Council rather than through
a separately elected steering committee. This reflects how the project actually
operates: the people who approve changes to critical parts of the project are the
same people who handle its overall direction. The Maintainer Council therefore
performs the steering and oversight functions for the project, and this section
supersedes earlier references in project documents to a separate "steering
committee."

The Maintainer Council is responsible for:

- Stewarding the vision, values, scope, and roadmap of the overall project.
- **Cross-organization coordination** between the core platform
  (`github.com/meshery`) and extensions (`github.com/meshery-extensions`), including
  approving new subprojects and extensions and resolving conflicts that cannot be
  handled within a single subproject.
- Acting as the final point of escalation for disputes in any project repository or
  group, before any escalation to the CNCF TOC.
- Defining and evolving project and subproject governance, including this document.
- Appointing and reviewing delegated teams, such as the
  [Security Response Team](#security-response), and, where the project chooses, a
  Code of Conduct response group.
- Making requests of the CNCF on behalf of the project (see
  [Requesting CNCF Resources](#requesting-cncf-resources)).

The Maintainer Council may delegate any of its powers to an individual, a subproject,
or a smaller team. Day-to-day technical decisions are handled within each subproject;
the Council exists for what cannot be, or should not be, decided within a single
subproject.

## Decision-Making

### Lazy consensus

Most business in Meshery is conducted by
[lazy consensus](https://community.apache.org/committers/lazyConsensus.html).
Proposals move forward when there is no sustained objection. The project does not
take a formal tally for routine decisions, which are made in issues, pull requests,
discussions, and meetings.

### Voting

A formal vote is taken when any of the following is true:

1. There is expressed disagreement on a matter that lazy consensus did not resolve.
2. The action is major or effectively irreversible, so that the project wants to be
   certain every relevant Maintainer is aware of it.
3. This document requires a vote, such as adding or removing a Maintainer, or
   changing this governance.

Votes are held in a GitHub issue labeled `vote` in the relevant repository, or on a
Maintainer communication channel, and remain open long enough for Maintainers in
different time zones to participate (normally one to two weeks for nominations).

Votes are counted against the total number of eligible Maintainers, not against the
number who happen to participate. This lets the project be flexible about where and
how a vote is held, because there is little risk of deliberate exclusion. The
applicable thresholds are:

- **Simple majority of all eligible Maintainers** for routine matters that require a
  vote, such as approving a request for CNCF resources or selecting a Subproject
  Lead.
- **Two-thirds (2/3) majority** for matters that warrant deliberation, such as
  adding or removing a Maintainer, admitting or removing a subproject, or changing
  this document.

For a vote scoped to a single subproject, "eligible Maintainers" means the
Maintainers of that subproject. For a project-wide vote, it means the Maintainer
Council.

### How specific kinds of decisions are made

- **Leadership roles.** Maintainers are added and removed through the
  [Maintainer Lifecycle](#maintainer-lifecycle). Organization membership is granted
  per the [Organization Member](#organization-member) requirements. Subproject Leads
  are selected by their subproject's Maintainers.
- **Acceptance of contributions.** A change is accepted when it meets the project's
  quality standards and receives approval from the maintainers responsible for the
  affected code or documentation, as recorded in `MAINTAINERS.md`. Anyone may
  review and comment; approval authority rests with those maintainers.
- **Requests to the CNCF.** Handled per
  [Requesting CNCF Resources](#requesting-cncf-resources).
- **Changes to governance or project goals.** Handled per
  [Changes to This Governance](#changes-to-this-governance), by a 2/3 vote of the
  Maintainer Council.

## Maintainer Lifecycle

This section defines how a person becomes a core maintainer, is onboarded, what they
are responsible for and entitled to, and how the role ends, whether voluntarily,
through inactivity, or into emeritus status. A working maintainer lifecycle,
including removal, is essential to prevent the project from stalling when people move
on.

Unless stated otherwise, this lifecycle describes **core maintainers**. **Extensions
maintainers** follow the lighter path described under
[Extensions](#extensions-githubcommeshery-extensions): they are nominated or
self-nominated and confirmed by the core maintainers, and the onboarding, inactivity,
and emeritus practices below apply to them within the extensions they maintain.

### Becoming a Core Maintainer

New core maintainers are selected through a transparent, merit-based process:

1. **Nomination.** A community member, typically an existing core maintainer, nominates an
   individual for the role. The nomination highlights the nominee's contributions and
   explains why they are a good fit. Nominees are usually already active Contributors
   or Organization Members with a track record in the relevant subproject.
2. **Announcement.** The nomination is announced to the community, for example on
   Slack or a mailing list, with information about the nominee and their
   contributions, and it describes the voting process.
3. **Voting period.** A voting period is opened, normally one to two weeks, so that
   everyone has time to consider the nomination.
4. **Voting.** Existing core maintainers cast binding votes. Other community members may
   express support or concerns through non-binding votes.
5. **Decision.** At the end of the period, votes are tallied. The nominee is approved
   on a **2/3 majority of the existing core maintainers** of the relevant subproject. If
   the nomination does not reach that threshold, it is not approved.
6. **Announcement of result.** The outcome is announced to the community. A
   successful nominee is welcomed and granted the responsibilities and privileges of
   the role.

### Onboarding and access

When a person becomes a core maintainer:

- They are added to [MAINTAINERS.md](./MAINTAINERS.md) with their name, contact
  information, the repositories and directories they are responsible for, and their
  affiliation.
- They are added to the relevant GitHub team or teams for their subproject, which
  grants the appropriate repository permission (normally `write`, and `maintain`
  where required). See
  [Systems Access and Repository Permissions](#systems-access-and-repository-permissions).
- They are added to the
  [CNCF Project Maintainers list](https://github.com/cncf/foundation/blob/main/project-maintainers.csv)
  so that the CNCF's records match the project's records.
- They are added to the private Maintainer communication channel and, where relevant,
  given access to other project systems for their role.

### Maintainer responsibilities

Maintainership is largely more of what a consistently engaged and impactful
contributor already does. Responsibilities include:

- Reviewing, approving, merging, and closing issues and pull requests in their area,
  and cutting releases.
- Regular attendance at project meetings and good-faith stewardship of the project
  and its external representation.
- Encouraging and mentoring new contributors.
- The janitorial work that keeps a project healthy, such as tests, release
  coordination, meeting preparation and hosting, and following up with contributors
  to complete work.
- Upholding the project's standards, this governance, and the
  [Code of Conduct](./CODE_OF_CONDUCT.md).

### Maintainer privileges

These responsibilities come with the following:

1. A voice and a vote on architectural, roadmap, and general project decisions.
2. Approval rights on GitHub pull requests that count toward branch protection
   requirements for the code they own.
3. Commit and write access to the repositories they maintain, which greatly improves
   their ability to collaborate with contributors on pull requests.
4. Depending on the project's level in the CNCF, speaking spots at KubeCon, booth
   hosting, webinars, and similar opportunities.
5. Advance insight into unannounced project vulnerabilities and a role in
   coordinating the patch and public response.
6. Authorship rights on the project blog.

### Stepping down

Maintainers may step down voluntarily at any time. If you decide to step down, please
inform the other Maintainers and explain your reasons. Before stepping down, please
help identify and train your replacement and ensure a smooth handover. When ready,
you or another Maintainer submits a pull request moving your name from the active
list to the [Emeritus Maintainers](#emeritus-maintainers) list in
[MAINTAINERS.md](./MAINTAINERS.md), and another Maintainer adjusts your permissions
and access accordingly, including removal from the
[CNCF Project Maintainers list](https://github.com/cncf/foundation/blob/main/project-maintainers.csv).

### Inactivity and removal

Continuity depends on active Maintainers, so the project removes Maintainers who have
become inactive.

- Any Maintainer who is unresponsive or non-participatory across the project's
  communication channels for **two months** forfeits Maintainer privileges and
  returns to a Contributor role.
- Because much communication is asynchronous, Maintainers are nonetheless required to
  attend at least **one meeting per month**, which enables the high-fidelity,
  synchronous communication needed to convey complex topics and advance the project
  in real time.
- A Maintainer may also be removed for cause by a **2/3 vote** of the relevant
  subproject's maintainers, or of the core maintainers for a project-wide concern.
- When a Maintainer is removed, the change is recorded in
  [MAINTAINERS.md](./MAINTAINERS.md), their access is revoked, and they are removed
  from the CNCF Project Maintainers list.

Removal is not a judgment of a person's past contributions. A removed Maintainer is
welcome to re-earn the role later through the normal process.

### Emeritus Maintainers

Emeritus Maintainers are former Maintainers who have stepped down or become inactive
but whose past contributions the project continues to recognize. Emeritus Maintainers
are listed in [MAINTAINERS.md](./MAINTAINERS.md). They hold no special permissions
and are not counted when computing vote thresholds, but they are welcome to
participate as Contributors and may return to active Maintainership through the normal
nomination and voting process.

## Subproject Lifecycle

### Adding a subproject or extension

A new subproject or extension may be proposed by any Maintainer or by a group of
contributors who intend to maintain it. A proposal:

- States the scope and purpose of the subproject, why it belongs in the Meshery
  ecosystem, and which organization it should live in
  (`github.com/meshery` for core areas, `github.com/meshery-extensions` for
  community extensions).
- Identifies the initial maintainers and how they meet the maintainer requirements.
- Confirms that it will adopt this governance, the
  [Code of Conduct](./CODE_OF_CONDUCT.md), the project's security practices, and the
  relevant integration guidelines.

Core subprojects are admitted by a **2/3 vote** of the core maintainers. Extensions
in `github.com/meshery-extensions` are admitted by **confirmation from the core
maintainers**, based on the maintainers' contribution history and the extension's quality
and compatibility, consistent with the lighter governance described in
[Extensions](#extensions-githubcommeshery-extensions).

### Experimental subprojects

The project may admit an **experimental** subproject or extension that shows promise
but does not yet have production-ready code or an established pool of maintainers.
Experimental subprojects are clearly labeled as such, are given a route to join the
ecosystem without the same expectations and authority as established subprojects, and
are expected either to graduate to full status or to be archived within a reasonable
period.

### Removing or archiving a subproject

A subproject or extension may be archived when it is no longer maintained, no longer
aligns with the project, or has been superseded. Core subprojects are archived by a
**2/3 vote** of the core maintainers; extensions are archived by the core maintainers in
coordination with their maintainers. Archived repositories are moved to the project's
archive (attic) and marked read-only, with a note pointing to any replacement.

## Systems Access and Repository Permissions

This section documents the mechanics of the project's roles: which systems and
permissions back each role, and how access is granted and revoked. Code and
documentation ownership in GitHub and elsewhere is kept consistent with the roles
defined in this document.

### GitHub repository permissions by role

The project uses GitHub's permission levels and grants them through GitHub teams, so
that access follows role rather than individual. The mapping is:

| Role | GitHub permission | How it is granted |
| :--- | :--- | :--- |
| Community Participant | Read (public) | Default for public repositories |
| Contributor | Read; contributes via forks and pull requests | No grant required |
| Organization Member | Triage across org repositories | Added to the org members team after meeting requirements |
| Core maintainer | Write (and Maintain where required) on the core repositories and directories they own in `github.com/meshery` | Added to the subproject's GitHub team |
| Extensions maintainer | Write (and Maintain where required) on the extension repositories they maintain in `github.com/meshery-extensions` | Added to the extension's GitHub team |
| Subproject Lead | Maintain on the subproject's repositories | Added to the subproject's lead team |
| Organization administration | Admin | Held by a small number of core maintainers, designated by the Maintainer Council |

`Admin` access is limited to the smallest practical number of maintainers needed to
administer each organization. Granting or revoking `Admin` is a decision of the core
maintainers.

### GitHub teams and MAINTAINERS.md

- **GitHub teams.** Permissions are assigned through teams such as
  `@meshery/maintainers` and per-subproject teams (for example, a server team, a UI
  team, a docs team). Adding a person to the appropriate team is how their role's
  access is granted; removing them is how it is revoked.
- **MAINTAINERS.md.** The single source of truth for who maintains what.
  [MAINTAINERS.md](./MAINTAINERS.md) lists every Maintainer with their name, contact
  information, **the specific repositories and directories they are responsible for**
  (for example `meshery/meshery`, `meshery/schemas`, `meshery/meshsync`, and
  directories within them), their subproject and Lead status, and their affiliation.
  It also lists Emeritus Maintainers. This list is kept in sync with the
  [CNCF Project Maintainers list](https://github.com/cncf/foundation/blob/main/project-maintainers.csv).

Branch protection on the default branch requires review and approval from the
maintainers responsible for the affected area before merge, which is how the project
ensures that ownership in GitHub matches the governance roles in this document.

### Non-GitHub systems

Access to the project's other systems also follows role and is granted at onboarding
and revoked at offboarding:

- **CNCF Zoom** for community and maintainer meetings.
- **Slack** ([slack.meshery.io](https://slack.meshery.io)); workspace administration
  is held by a small number of Maintainers and Community Managers.
- **Discussion forum** ([community forum](https://discuss.meshery.io)).
- **Community calendar** ([meshery.io/calendar](https://meshery.io/calendar)).
- **Documentation and website** publishing.
- **Release and registry systems**, including container registries and package
  registries used to publish releases.
- **Recognition and badges** system ([badges.layer5.io](https://badges.layer5.io)).

When a person changes or leaves a role, the Maintainers and Community Managers
responsible for each system update access accordingly.

## Communication Channels

The project conducts its work on public, vendor-neutral channels. All channels,
including those used by subprojects, are listed here. Any non-public channel and its
purpose is also listed.

**Public channels:**

- **Slack:** [slack.meshery.io](https://slack.meshery.io), including
  subproject-specific and topic-specific channels.
- **Discussion forum:** [discuss.meshery.io](https://discuss.meshery.io).
- **GitHub issues, pull requests, and discussions** across
  [github.com/meshery](https://github.com/meshery) and
  [github.com/meshery-extensions](https://github.com/meshery-extensions).
- **Community meetings** on CNCF Zoom, published on the
  [community calendar](https://meshery.io/calendar), with
  [recordings on YouTube](https://www.youtube.com/@mesheryio).
- **Project blog** for announcements.

**Non-public channels:**

- **Private Maintainer channel.** Used only for sensitive matters that cannot be
  discussed in public, such as Code of Conduct reports, security vulnerabilities
  prior to disclosure, and individual personnel matters such as maintainer nominations
  where privacy is warranted. Routine project decisions are not made here.
- **Security reporting.** Vulnerability reports are received privately per the
  [Security Response](#security-response) process.

## Meetings

The project holds regular, public meetings, and their schedule is published on the
[community calendar](https://meshery.io/calendar) and integrated with the CNCF
calendar. These include weekly community and newcomer meetings, where official
project decisions may be made, alongside discussion in issues, pull requests, and on
the forum. Anyone is welcome to join; see [meshery.io/community](https://meshery.io/community)
and the [community handbook](https://meshery.io/community#handbook). Meetings are
recorded and posted publicly.

## Code of Conduct

The project adopts and adheres to the
[Meshery Code of Conduct](./CODE_OF_CONDUCT.md), which is based on, and not in
conflict with, the [CNCF Code of Conduct](https://github.com/cncf/foundation/blob/main/code-of-conduct.md).
All community members, in both organizations and across all channels, are expected to
follow it.

Code of Conduct reports are received confidentially by the core maintainers. If a
report concerns a core maintainer, the core maintainers designate two uninvolved core
maintainers to handle it. Reports that require a full investigation, or that involve a
core maintainer, are forwarded to the
[CNCF Code of Conduct Committee](https://github.com/cncf/foundation/blob/main/code-of-conduct.md)
for handling, and the core maintainers appoint a non-involved contributor to work with
the Committee as needed.

## Security Response

The core maintainers appoint a **Security Response Team** to handle reports of
security vulnerabilities. This team may consist of the core maintainers themselves; if
the responsibility is delegated, the core maintainers appoint a team of at least
**two** contributors. The core maintainers review the composition of the Security
Response Team at least **once a year**.

The Security Response Team handles all reports according to the project's
[Security Policy](./SECURITY.md). Vulnerabilities are reported privately to
[security@meshery.dev](mailto:security@meshery.dev) and are acknowledged within the
timeframe stated in the security policy. Maintainers receive advance notice of
unannounced vulnerabilities and help coordinate the fix and public disclosure.

## Requesting CNCF Resources

Any Maintainer may suggest a request for CNCF resources by filing an issue in the
project's [community repository](https://github.com/meshery/meshery). A **simple
majority of the core maintainers** approves the request. The core maintainers may also
delegate working with the CNCF to non-maintainer community members, such as Community
Managers, who are then added to the
[CNCF Project Maintainers list](https://github.com/cncf/foundation/blob/main/project-maintainers.csv)
for that purpose.

## How to Contribute

Everyone is welcome to contribute. Start with the
[Contributor Guides](https://docs.meshery.io/project/contributing) and
[CONTRIBUTING.md](./CONTRIBUTING.md), the [Newcomers' Guide](https://meshery.io/community/newcomers),
and the [community handbook](https://meshery.io/community#handbook). Good first issues
are labeled [help wanted](https://github.com/meshery/meshery/labels/help%20wanted),
and a [MeshMate](#meshmates) can help you get started.

Contribution expectations include:

- **Sign-off.** All commits must be signed off under the Developer Certificate of
  Origin (`git commit -s`).
- **Commit messages.** Use the project's `[component] message` convention and
  reference related issues and pull requests.
- **Quality gates.** Changes must pass the project's required checks and reviews
  before merge, and significant features should include corresponding documentation
  updates.
- **AI-assisted contributions.** Contributions created with the help of AI tooling
  are welcome, and are held to the same standards as any other contribution: the
  contributor is responsible for the correctness, licensing, and quality of what they
  submit, must sign off under the DCO, and must comply with the project's AI
  contribution policy as described in [CONTRIBUTING.md](./CONTRIBUTING.md).

The project's direction and upcoming work are described in the
[Roadmap](./ROADMAP.md).

## Changes to This Governance

This document is expected to evolve as the project does. Changes are proposed through
a pull request to this file and are adopted by a **2/3 vote** of the Maintainer
Council. Material changes are announced to the community. The project keeps the
history of this document in version control so that the evolution of its governance
remains visible alongside the evolution of the project itself.

## Governance Lineage

This governance builds on the project's prior governance and continues to iterate on
it as the project has grown, including the
[2026 reorganization](https://www.cncf.io/blog/2026/03/04/scaling-organizational-structure-with-mesherys-expanding-ecosystem/)
that partitioned the ecosystem into the `meshery` and `meshery-extensions`
organizations. It draws on the CNCF Contributor Strategy
[Maintainer Council](https://contribute.cncf.io/projects/best-practices/governance/templates/governance-maintainer/)
and [Subprojects](https://contribute.cncf.io/projects/best-practices/governance/templates/governance-subprojects/)
governance templates, and on the practices of graduated CNCF projects, including the
[Kubernetes governance](https://github.com/kubernetes/community/blob/master/governance.md)
and [OWNERS](https://www.kubernetes.dev/docs/guide/owners/) model.

---

## Appendix A: Mapping to CNCF Incubation Governance Review Criteria

This appendix maps each consideration in the CNCF TOC
[Governance Review Template](https://github.com/cncf/toc/blob/main/toc_subprojects/project-reviews-subproject/governance-review-template.md)
to where it is addressed, so the project and reviewers can confirm coverage during
Incubation diligence. It is a self-assessment companion to this document and may be
maintained as part of the project's Incubation self-assessment rather than published
as part of `GOVERNANCE.md`. Levels shown are the template's: **R** = Required at
Incubation, **S** = Suggested at Incubation (and typically Required at Graduation).

| # | Criterion | Level | Where addressed |
| :-- | :--- | :--- | :--- |
| 1 | Governance evolution / history | S | [Governance Lineage](#governance-lineage); version history of this file |
| 2 | Clear and discoverable governance documentation | S | This `GOVERNANCE.md`, linked from the README and contributing guides |
| 3 | Governance up to date with actual activities (meetings, elections, leadership, approvals) | S | [Meetings](#meetings), [Decision-Making](#decision-making), [Maintainer Lifecycle](#maintainer-lifecycle) |
| 4 | Vendor-neutrality of project direction documented | S | [Vendor Neutrality](#vendor-neutrality) |
| 5 | How the project decides leadership roles, contribution acceptance, requests to CNCF, and governance/goal changes | S | [Decision-Making](#how-specific-kinds-of-decisions-are-made), [Requesting CNCF Resources](#requesting-cncf-resources), [Changes to This Governance](#changes-to-this-governance) |
| 6 | How role/function-based members or sub-teams are assigned, onboarded, and removed (e.g., Security Response) | S | [Security Response](#security-response), [Subproject Lifecycle](#subproject-lifecycle), [Systems Access](#systems-access-and-repository-permissions) |
| 7 | Complete maintainer lifecycle (roles, onboarding, offboarding, emeritus) | S | [Maintainer Lifecycle](#maintainer-lifecycle) |
| 8 | Demonstrated use of the maintainer lifecycle (additions/replacements) | S | Record outcomes in [MAINTAINERS.md](./MAINTAINERS.md) history (active and emeritus changes) |
| 9 | Complete list of current maintainers with names, contact, domain of responsibility, and affiliation | R | [MAINTAINERS.md](./MAINTAINERS.md); structure in [Appendix B](#appendix-b-maintainersmd-structure) |
| 10 | A number of active maintainers appropriate to size and scope | R | [MAINTAINERS.md](./MAINTAINERS.md) (maintain active count across subprojects) |
| 11 | Maintainers from at least 2 organizations (survivability) | N/A at Incubation | Affiliation column in [MAINTAINERS.md](./MAINTAINERS.md) |
| 12 | Code and doc ownership in GitHub matches documented governance roles | R | [Systems Access and Repository Permissions](#systems-access-and-repository-permissions) |
| 13 | Adoption and adherence to CNCF CoC, or a project CoC based on it | R | [Code of Conduct](#code-of-conduct), [CODE_OF_CONDUCT.md](./CODE_OF_CONDUCT.md) |
| 14 | CNCF CoC cross-linked from other governance documents | R | [Code of Conduct](#code-of-conduct) links the CNCF CoC |
| 15 | All subprojects listed | R | [Subprojects and domains](#subprojects-and-domains); [MAINTAINERS.md](./MAINTAINERS.md) for the authoritative list |
| 16 | Subproject leadership, contribution, maturity status, and add/remove process documented | S | [Subprojects](#subprojects-and-domains), [Subproject Lifecycle](#subproject-lifecycle) |
| 17 | Contributor ladder with multiple roles | S | [Roles and the Contributor Ladder](#roles-and-the-contributor-ladder) |
| 18 | Clearly defined and discoverable process to submit issues or changes | R | [How to Contribute](#how-to-contribute), [CONTRIBUTING.md](./CONTRIBUTING.md) |
| 19 | At least one public communications channel documented | R | [Communication Channels](#communication-channels) |
| 20 | All communication channels listed, including subprojects and any non-public channels and their purpose | R | [Communication Channels](#communication-channels) |
| 21 | Up-to-date public meeting schedule / CNCF calendar integration | R | [Meetings](#meetings); [community calendar](https://meshery.io/calendar) |
| 22 | Documentation of how to contribute, with increasing detail as the project matures | R | [How to Contribute](#how-to-contribute); [contributing guides](https://docs.meshery.io/project/contributing) |
| 23 | Demonstrated contributor activity and recruitment | R | [contributor list](https://github.com/meshery/meshery/graphs/contributors), [MeshMates](#meshmates), [recognition program](#recognition-and-certification), newcomer meetings |

**Resolution of the open governance issues.** This document also resolves the two
open governance issues that prompted this work:

- **[#19447](https://github.com/meshery/meshery/issues/19447) (describe maintainer
  assignments explicitly).** Addressed by
  [Systems Access and Repository Permissions](#systems-access-and-repository-permissions)
  and [Appendix B](#appendix-b-maintainersmd-structure): maintainers
  are mapped to specific repositories and directories in `MAINTAINERS.md`; the
  document states how maintainer rights are granted (GitHub teams) and revoked,
  including the two-month inactivity rule.
- **[#19460](https://github.com/meshery/meshery/issues/19460) (describe the process
  for becoming a Contributor).** Addressed by the
  [Contributor](#contributor) section: it describes the benefits and permissions of
  Contributors and the process to become one, links to a list of contributors, and
  replaces the former standalone "Reviewer" role with a clear statement that anyone
  is welcome to review and comment on work in progress.

## Appendix B: MAINTAINERS.md Structure

This appendix specifies the structure that satisfies criterion 9 and issue #19447. It
is a template, not a list of people; the live data lives in `MAINTAINERS.md`.

### MAINTAINERS.md

`MAINTAINERS.md` lists, for every Maintainer, at least the following, grouped by
subproject:

| Column | Description |
| :--- | :--- |
| Name | The maintainer's name |
| GitHub handle | Used to map to GitHub teams |
| Class | Whether the person is a **core maintainer** or an **extensions maintainer** |
| Contact | A reachable contact, such as email or Slack handle |
| Subproject | The subproject, domain, or extension, and whether the person is its Lead |
| Repositories and directories | The specific repositories and paths the person is responsible for, for example `meshery/meshery` (`/server`, `/mesheryctl`), `meshery/schemas`, `meshery/meshsync`, or an extension under `meshery-extensions` |
| Affiliation | The person's employer or "independent", for vendor-neutrality tracking |

`MAINTAINERS.md` also contains an **Emeritus Maintainers** section and a note that the
list is kept in sync with the
[CNCF Project Maintainers list](https://github.com/cncf/foundation/blob/main/project-maintainers.csv).

If a repository in either organization lacks its own `MAINTAINERS.md`, the
authoritative mapping for it is the entry in
[meshery/meshery/MAINTAINERS.md](./MAINTAINERS.md).
