# Meshery vision hypotheticals

This record preserves the complete review set, its steelman, and the author's verdicts.

## Round 1 - 2026-08-26

### H-1 - A generic resource importer without a model

Proposal: Add an import path that lets any Kubernetes or cloud resource appear in Meshery without a model, relationship behavior, or adapter capability.
Proposal: Teams can view and manually apply those resources while an integration is still being built.
Proposal: The importer would make new infrastructure visible sooner than a model release can.
Tests: "A Meshery design represents components and relationships that the platform can evaluate, render, and operate."
Why the answer is not obvious: Teams need fast coverage for long-tail infrastructure and a useful stopgap before a full integration exists.
Why the answer is not obvious: Rendering objects that Meshery cannot evaluate or safely operate dilutes the owned surface and makes capability ambiguous.
Verdict: In vision.
Reasoning: Users can use annotation-only components and annotation-only relationships (like those offered in the meshery-core and meshery-shapes models) to visually depict their full architecture or to visually and collaboratively convey a concept with the acknowledgement that those annotation-only components and relationships are not semantically meaningful; they do not bear configuration properties that can be orchestrated or deployed.
Reasoning: That annotation-only components are excluded from deployment operations.

### H-2 - A marketplace extension with a trust banner

Proposal: Allow any signed third-party UI extension to install from a marketplace after an administrator reads and accepts a trust warning.
Proposal: The extension still runs in the Meshery browser origin because sandboxing would break current extension capabilities.
Proposal: The marketplace would centralize discovery without adding a new isolation model.
Tests: "Meshery treats an enabled extension as code inside the deployment's trust boundary."
Why the answer is not obvious: A trust banner and administrator consent may be enough for an ecosystem where integrations need a practical distribution path.
Why the answer is not obvious: Centralizing installation can make a known trust boundary feel safer than it is and normalize unsafe defaults.
Verdict: In vision.

### H-3 - A temporary local API endpoint

Proposal: Let the UI add a hand-written endpoint and response type for a schema-owned API while a generated schemas release is delayed.
Proposal: The local endpoint can carry a cache tag that a generated hook currently lacks.
Proposal: The implementation would be deleted after the next schemas release.
Tests: "Meshery uses a generated schema client and type when that shared contract already exists."
Why the answer is not obvious: It can unblock a time-sensitive product fix and preserve precise cache behavior.
Why the answer is not obvious: Temporary duplicate contracts routinely outlive their deadline and silently diverge across Meshery and its providers.
Verdict: Off mission.
Reasoning: We need to ratchet away from the local API endpoints.
Reasoning: Meshery is a schema-driven project.
Reasoning: The source of truth for all API endpoints shall be sourced from https://github.com/meshery/schemas.

### H-4 - Canonicalize credentials on the next save

Proposal: When a user edits any legacy credential, rewrite its secret into the canonical schema before saving.
Proposal: The read path remains tolerant, but the fleet converges gradually without a separate migration job.
Proposal: The UI would show only a normal successful save message.
Tests: "Meshery preserves legacy persisted data during a canonical form transition through tolerant reads rather than silent rewrites."
Why the answer is not obvious: Each ordinary edit becomes a low-cost opportunity to reduce format diversity.
Why the answer is not obvious: A normal save is not informed migration consent, and a write-side change can alter data the user did not intend to migrate.
Verdict: In vision.

### H-5 - Keep a permission-denied design entry visible

Proposal: Always show the Design Configurator entry to organization members, even when they lack View Designs permission.
Proposal: Selecting it opens the standard permission-denied page with a link to request access.
Proposal: This would make capability discovery and support diagnosis easier for custom roles.
Tests: "Meshery gates protected content as well as protected controls and prevents denied content from loading."
Why the answer is not obvious: A visible route can explain why access is unavailable and help users request the right role.
Why the answer is not obvious: Controls that knowingly lead to denial are misleading and normalize a gap between navigation and authorization.
Verdict: Off mission.

### H-6 - Allow an unverified explicit chart pin

Proposal: Permit an air-gapped administrator to install an explicitly pinned operator chart that Meshery cannot verify against the public repository.
Proposal: Require a typed confirmation that the operator accepts responsibility for the version.
Proposal: This would support disconnected enterprise deployments without changing derived-version behavior.
Tests: "Meshery honors an operator's explicit configuration or refuses it with an actionable diagnostic."
Why the answer is not obvious: An explicit operator choice in a disconnected environment may be safer than blocking deployment outright.
Why the answer is not obvious: A confirmation can turn an unsupported state into an apparently supported one and make later diagnosis unreliable.
Verdict: In vision.

### H-7 - Skip missing protected-branch credentials

Proposal: Change CI so a missing remote-provider credential always skips the provider project, including protected branches.
Proposal: Forks and external contributors would keep a green run when they cannot access secrets.
Proposal: A mainline secret regression would no longer fail the workflow before tests begin.
Tests: "Meshery fails a required remote-provider setup in CI rather than passing a suite that did not run."
Why the answer is not obvious: Contributors should not receive a red result for credentials they cannot access or repair.
Why the answer is not obvious: Protected branches need an unmistakable signal when their required suite has not exercised its promised coverage.
Verdict: In vision.

### H-8 - Run relationship integration tests nightly

Proposal: Run only fast Rego lint and unit tests on relationship-policy pull requests, and move complete design integration scenarios to a nightly workflow.
Proposal: The pull-request cycle becomes faster and the nightly suite can use more expensive fixtures.
Proposal: A failed nightly run is triaged after merge.
Tests: "Meshery tests relationship policy against complete designs and real registrant shapes."
Why the answer is not obvious: Policy contributors get faster feedback and a richer nightly environment may uncover more cases.
Why the answer is not obvious: Merging before realistic relationship evaluation runs makes a passing pull request less meaningful and delays the failure to users or maintainers.
Verdict: Conditional.
Reasoning: Not now.
Reasoning: Perhaps, in the future.

### H-9 - Make annotation edges operational by default

Proposal: Treat every user-created annotation relationship as a lifecycle constraint so that users do not have to choose between semantic and non-semantic edges.
Proposal: The canvas becomes more uniform and policy authors can use every connection.
Proposal: Existing annotations would gain behavior unless explicitly opted out.
Tests: "Meshery distinguishes relationships that guide lifecycle behavior from annotations that only aid human understanding."
Why the answer is not obvious: One relationship model is simpler to teach and expands the policy surface.
Why the answer is not obvious: Visual grouping and human notes would accidentally acquire deployment authority, violating the distinction that protects operator intent.
Verdict: Off mission.

### H-10 - Relax a conservative SQL lint rule globally

Proposal: Accept any interface-typed ORDER BY value without a diagnostic after a partner generic helper triggers a conservative false positive.
Proposal: The change removes friction for reusable query helpers and avoids per-call suppressions.
Proposal: It also makes an interface value carrying an untrusted string unprovable to the analyzer.
Tests: "Meshery favors a documented, narrow false positive over a silent security false negative."
Why the answer is not obvious: Excessive false positives can discourage adoption of a safety rule and block valid abstraction.
Why the answer is not obvious: Broad relaxation recreates the exact unobserved path the rule exists to prevent, while a narrow suppression keeps the risk visible.
Verdict: Off mission.

## Round 2 - 2026-08-27

### H-11 - Curate the extension marketplace

Proposal: Require Meshery maintainers to manually review each marketplace publisher before the first installation can be approved.
Proposal: The marketplace would still show signature and trust-boundary information, but review would signal an additional project-level judgment.
Proposal: This slows the long tail of integrations.
Tests: "Meshery may distribute signed third-party extensions through a marketplace when administrator activation makes the trust boundary explicit."
Why the answer is not obvious: Manual review can prevent a marketplace from looking like a list of unvetted code and protects users who equate discovery with endorsement.
Why the answer is not obvious: Centralized review can bottleneck ecosystem growth and blur the administrator responsibility that the explicit trust boundary is meant to preserve.
Verdict: In vision.

### H-12 - Auto-update marketplace extensions

Proposal: Let an enabled marketplace extension update automatically to the latest verified version during the next Meshery restart.
Proposal: Administrators receive a changelog after the update rather than approving each version.
Proposal: Vulnerability repairs reach deployments faster, but a plugin behavior change can arrive without a chosen maintenance window.
Tests: "Meshery treats an enabled extension as code inside the deployment trust boundary."
Why the answer is not obvious: Automatic updates reduce known-vulnerable extension time and operational work.
Why the answer is not obvious: Code running inside a deployment trust boundary should not change without a deliberate operator choice.
Verdict: In vision.

### H-13 - Consume a schemas prerelease

Proposal: Permit a Meshery provider to depend on a signed prerelease of meshery/schemas when a contract must ship before the normal schemas release.
Proposal: The endpoint remains generated from schemas, but repositories can coordinate through an expiring prerelease version.
Proposal: The prerelease tag would be removed after release.
Tests: "Every HTTP API endpoint sources its contract from the meshery/schemas repository."
Why the answer is not obvious: The source of truth remains schemas while urgent cross-repository work is not blocked on release timing.
Why the answer is not obvious: Prerelease ranges can become an untracked second release channel and weaken contract discipline.
Verdict: In vision.

### H-14 - Drop unknown fields during credential canonicalization

Proposal: When intentional credential save encounters obsolete fields that do not exist in the canonical schema, drop them and write canonical state rather than preserving a legacy envelope.
Proposal: The save completes normally and logs a migration event.
Proposal: This makes later reads simple but can lose information no current schema recognizes.
Tests: "Meshery canonicalizes legacy credential data on an intentional save while retaining tolerant reads during the transition."
Why the answer is not obvious: Canonicalization should reduce persistent complexity instead of carrying obsolete fields forever.
Why the answer is not obvious: Unknown fields may carry provider-specific intent that the current schema has not yet modeled.
Verdict: In vision.

### H-15 - Allow skipped provider coverage to merge

Proposal: Allow a protected-branch pull request to merge when its required provider project is visibly skipped because credentials are unavailable.
Proposal: A required-check summary names the missing suite, but the skipped check does not block the merge.
Proposal: This turns visible skipped coverage into a delivery rule.
Tests: "Meshery marks a provider CI project as skipped when required credentials are unavailable, including on protected branches."
Why the answer is not obvious: A credential outage should not indefinitely block unrelated urgent fixes.
Why the answer is not obvious: A protected branch could merge a provider regression while its promised coverage is absent.
Verdict: In vision.

### H-16 - Scope relationship integration by changed path

Proposal: Run complete relationship integration scenarios only when a pull request changes policy, model, or deployment-engine paths.
Proposal: Documentation, UI styling, and unrelated API changes would retain unit and lint coverage but bypass relationship scenarios.
Proposal: This preserves pull-request coverage where change detection believes it matters.
Tests: "Meshery keeps relationship integration coverage on pull requests today rather than deferring it to a nightly-only job."
Why the answer is not obvious: Targeted execution shortens feedback for changes with no relationship behavior.
Why the answer is not obvious: Path rules miss indirect coupling and can make a green pull request less meaningful.
Verdict: In vision.
Reasoning: Skip on documentation-only changes and provider-ui-only changes, but execute on all other PRs.

### H-17 - Verify an air-gapped pin from local evidence

Proposal: Let an air-gapped administrator upload an operator chart and manifest checksum to a local Meshery registry before choosing an explicit chart pin.
Proposal: Meshery verifies the selected pin against that local evidence without contacting a public chart index.
Proposal: This gives an explicit pin a stronger audit trail while adding artifact lifecycle responsibility to Meshery.
Tests: "Meshery permits an air-gapped administrator to use an unverified explicit chart pin only after a typed acknowledgement of responsibility."
Why the answer is not obvious: Local checksums preserve operator control while making the chosen artifact auditable.
Why the answer is not obvious: A local artifact registry expands Meshery from control surface into a package distribution system.
Verdict: In vision.

### H-18 - Allow a generated-code safety baseline

Proposal: Allow a generated package to register a package-level SQL lint exception after its generator has a test proving it cannot emit an untrusted ORDER BY value.
Proposal: The generated code would no longer need line-by-line narrow suppressions.
Proposal: A future generator change could expand the exception beyond the proof.
Tests: "Meshery favors a documented, narrow safety exception over a global relaxation that creates a silent security false negative."
Why the answer is not obvious: Generated code should be assessed at its source and not burden every generated call site.
Why the answer is not obvious: A broad package exception can outlive the generator proof and hide a later unsafe path.
Verdict: In vision.

### H-19 - Let remote providers publish private schemas

Proposal: Allow a remote provider to publish a private schema construct that emits generated endpoints without merging that construct into the public meshery/schemas repository.
Proposal: The generated client retains a schema source, but contract governance becomes distributed.
Proposal: This could serve proprietary provider capabilities without local endpoint copies.
Tests: "Every HTTP API endpoint sources its contract from the meshery/schemas repository."
Why the answer is not obvious: Private providers need a way to evolve proprietary capabilities without exposing their contracts publicly.
Why the answer is not obvious: Distributed schema ownership erodes the single source of truth that prevents consumer drift.
Verdict: In vision.

### H-20 - Use encrypted production-derived test fixtures

Proposal: Allow relationship integration tests to load encrypted, redacted production-derived design fixtures from a controlled test store.
Proposal: This exercises real metadata patterns that synthetic fixtures miss.
Proposal: It adds credentials, retention, and reproducibility obligations to the test system.
Tests: "Meshery tests relationship policy against complete designs and real registrant shapes."
Why the answer is not obvious: Production-derived fixtures expose the shapes that have already defeated synthetic coverage.
Why the answer is not obvious: Sensitive data handling and external test dependencies can make a safety suite less portable and less trustworthy.
Verdict: In vision.

### A-1 - Approve the revised acceptance policy

Proposal: Adopt the full revised VISION.md as the current acceptance policy for Meshery.
Tests: The full revised draft is visible in the review board.
Why the answer is not obvious: Approval makes the recorded decisions govern future change review, while a final adjustment can name the sentence that still needs work.
Verdict: In vision.
