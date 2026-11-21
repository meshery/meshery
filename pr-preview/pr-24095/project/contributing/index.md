# Contributing

> How to contribute to the Meshery project and any of its components.

---

LLMS index: [llms.txt](/pr-preview/pr-24095/llms.txt)

---

# Contributing 

Please do! Thanks for your help! 🎈 Meshery is community-built and welcomes collaboration. Contributors are expected to adhere to the [CNCF's Code of Conduct](https://github.com/meshery/meshery/blob/master/CODE_OF_CONDUCT.md).

## General Contribution Flow

Meshery and it's various architectural components are written in different languages, including Golang, Javascript (React.js and Next.js) To make building, testing, and the experience of contributing consistent across all Meshery components, a `Makefile` is included in the every repository. These `make` targets are what you will use to build, run, test, and document.

To contribute to Meshery, please follow this basic fork-and-pull request [gitflow](/pr-preview/pr-24095/project/contributing/contributing-gitflow/).

<details>
<summary>Adding your sign-off on commits (Developer Certificate of Origin)</summary>
<ul>
<li>
To contribute to this project, you must agree to the Developer Certificate of
Origin (DCO) for each commit you make. The DCO is a simple statement that you,
as a contributor, have the legal right to make the contribution.
</li>

<li>
See the <a href="https://developercertificate.org"> DCO </a> file for the full text of what you must agree to
and how it works <a href="https://github.com/probot/dco#how-it-works">here</a>.
To signify that you agree to the DCO for contributions, you simply add a line to each of your
git commit messages:

<pre><code>
Signed-off-by: Jane Smith <jane.smith@example.com>
</code></pre></li>

<li>
In most cases, you can add this signoff to your commit automatically with the
<code>-s</code> or <code>--signoff</code> flag to <code>git commit</code>. You must use your real name and a reachable email
address (sorry, no pseudonyms or anonymous contributions). An example of signing off on a commit:

<pre><code>$ git commit -s -m “my commit message w/signoff”</code></pre>
</li>

<li>
To ensure all your commits are signed, you may choose to add this alias to your global <code>.gitconfig</code>:

~/.gitconfig

<pre><code>
[alias]
  amend = commit -s --amend
  cm = commit -s -m
  commit = commit -s
</code></pre>

Or you may configure your IDE, for example, VSCode to automatically sign-off commits for you:<a href="https://user-images.githubusercontent.com/7570704/64490167-98906400-d25a-11e9-8b8a-5f465b854d49.png" ><img src="https://user-images.githubusercontent.com/7570704/64490167-98906400-d25a-11e9-8b8a-5f465b854d49.png" width="50%"/></a>

</li>

<li>
Should an unsigned commit slip through anyway, the repository's <code>commit-msg</code> hook
(installed with the UI dependencies, see <code>ui/.husky/</code>) rejects it locally rather than
letting the DCO check fail in CI. The hook applies the same rules the DCO check does, so a
message it accepts is one CI accepts: merge commits are exempt, and at least one
<code>Signed-off-by</code> trailer has to name the author or the committer. When it rejects a
commit, nothing was committed - your message is waiting in <code>.git/COMMIT_EDITMSG</code>, so
re-run the commit with <code>-s</code>.
</li>

<li>
Treat that hook as a safety net rather than a guarantee. Git runs whichever hooks
<code>core.hooksPath</code> resolves to, so a clone where <code>make ui-setup</code> was never
run has no hook at all, and any tool that repoints <code>core.hooksPath</code> at a directory of
its own - other hook managers, and worktree-based automation that commits on your behalf -
replaces it without a warning. Commit with <code>-s</code> rather than relying on the hook to
notice.
</li>

<li>
To sign off a commit you have already written but have <em>not</em> pushed:

<pre><code>$ git commit --amend -s --no-edit</code></pre>

Once the branch is pushed, amending reaches only its tip, and a DCO failure usually names more
than one commit. Sign off every commit the pull request carries and replace the pushed branch:

<pre><code>$ git rebase &lt;base-branch&gt; --signoff
$ git push --force-with-lease</code></pre>

That rewrites the branch, so settle it first with anyone else working on it, and with any
automation that is tracking the branch by commit hash - a force push moves the tip out from
under both.
</li>
</ul>
</details>

<details id="no-mistakes-quality-gate">
<summary>Running the pre-landing Quality Gate (No-Mistakes)</summary>
<div>
<p>
Meshery enforces the <a href="https://github.com/kunchenguid/no-mistakes">no-mistakes</a> quality gate on all code contributions (UI, server, CLI, schemas, and configurations). The gate validates code review, tests, documentation parity, and linting in isolated, disposable worktrees before changes reach GitHub, satisfying the <code>PR must be raised via no-mistakes</code> CI status check.
</p>

<h4>1. Install No-Mistakes</h4>
Install the CLI tool on your local development system:
<ul>
<li><strong>macOS / Linux:</strong>
<pre><code>curl -fsSL https://raw.githubusercontent.com/kunchenguid/no-mistakes/main/docs/install.sh | sh</code></pre>
</li>
<li><strong>Windows (PowerShell):</strong>
<pre><code>irm https://raw.githubusercontent.com/kunchenguid/no-mistakes/main/docs/install.ps1 | iex</code></pre>
</li>
<li><strong>Via Go:</strong>
<pre><code>go install github.com/kunchenguid/no-mistakes/cmd/no-mistakes@latest</code></pre>
</li>
</ul>
<p>For other installation options, refer to the <a href="https://kunchenguid.github.io/no-mistakes/start-here/installation/">official installation guide</a>.</p>

<h4>2. Check Prerequisites</h4>
Verify your local environment with the built-in diagnostic tool:
<pre><code>no-mistakes doctor</code></pre>
<p>
This checks for:
</p>
<ul>
<li><strong>Git:</strong> Required for repository and branch operations.</li>
<li><strong>Supported Agent Runner:</strong> Claude Code, Codex, Copilot, Antigravity (<code>agy</code>), etc.</li>
<li><strong>GitHub CLI (<code>gh</code>):</strong> Required for opening pull requests and monitoring CI. Run <code>gh auth login</code>, then confirm with <code>gh auth status</code>. Doctor checks only that <code>gh</code> is installed.</li>
</ul>

<h4>3. Initialize the Meshery Repository</h4>
Navigate to your local Meshery repository and initialize the local quality gate with your GitHub fork. Point <code>origin</code> at the parent repository first: if you cloned your fork (the gitflow guide's layout), <code>no-mistakes init</code> refuses rather than opening pull requests against your own fork. Run <code>git remote set-url origin git@github.com:meshery/meshery.git</code>, then pass your fork with <code>--fork-url</code>:
<pre><code>cd meshery
git remote set-url origin git@github.com:meshery/meshery.git
no-mistakes init --fork-url git@github.com:YOUR_USERNAME/meshery.git</code></pre>
<p>
This registers the local <code>no-mistakes</code> remote, prepares isolated worktrees, and starts the background daemon.
</p>

<h4>4. Commit and Run the Quality Gate</h4>
<ol>
<li>Create your feature branch:
<pre><code>git checkout -b feature/your-feature</code></pre>
</li>
<li>Make your changes and commit them with DCO sign-off:
<pre><code>git commit -s -m "type(scope): descriptive commit message"</code></pre>
</li>
<li>Run the pipeline on your branch:
<pre><code>no-mistakes axi run --intent "&lt;brief description of your changes&gt;"</code></pre>
</li>
</ol>
<p>
You can also push the branch through the gate directly with <code>git push no-mistakes feature/your-feature</code>.
</p>

<h4>5. Monitor Pipeline Progress</h4>
<p>
Run <code>no-mistakes</code> to watch the pipeline and respond to approval requests; <code>no-mistakes status</code> shows a summary.
</p>
<p>
The automated pipeline executes 9 structured stages:
</p>
<ol>
<li><strong>Intent:</strong> Analyzes the purpose and scope of your changes.</li>
<li><strong>Rebase:</strong> Rebases onto the latest upstream branch.</li>
<li><strong>Review:</strong> Evaluates code styling and safety rules.</li>
<li><strong>Test:</strong> Scopes and runs baseline tests.</li>
<li><strong>Document:</strong> Checks documentation parity and identifies gaps.</li>
<li><strong>Lint:</strong> Runs repository linters.</li>
<li><strong>Push:</strong> Pushes validated commits to your remote fork.</li>
<li><strong>PR:</strong> Creates or updates the pull request.</li>
<li><strong>CI:</strong> Monitors remote GitHub Actions workflows.</li>
</ol>

<h4>6. Pull Request Attestation &amp; Requirements</h4>
<ul>
<li><strong>Attestation Matching:</strong> The CI check <code>PR must be raised via no-mistakes</code> verifies that the pull request description contains a valid attestation matching the current <code>head_sha</code> and that <code>review</code>, <code>test</code>, and <code>document</code> steps are marked <code>completed</code>. After pushing new commits, re-run the gate so the attestation covers the new head; an attestation for an older commit fails the check. Before a follow-up commit, run <code>no-mistakes axi status</code> and follow its <code>branch_sync.next_action</code>: when it says <code>sync</code>, run <code>no-mistakes axi sync</code> first, since the gate may have advanced the published branch past your local copy.</li>
<li><strong>Separate DCO Check:</strong> All commits must include a sign-off trailer (<code>git commit -s</code>) to satisfy GitHub's independent <code>DCO</code> status check.</li>
<li><strong>Exempt Paths:</strong> Pull requests whose changed paths are strictly documentation (<code>docs/**</code>, <code>CNAME</code>) or CI workflows (<code>.github/workflows/**</code>) are automatically exempt from the pipeline requirement.</li>
</ul>

<p>
For agent skill instructions and CI merge gate specifications, see the <a href="/pr-preview/pr-24095/project/contributing/build-and-release/#no-mistakes-merge-gate">No-mistakes merge gate technical specification</a> and <code>.agents/skills/no-mistakes/SKILL.md</code>.
</p>
</div></details>

## Cloning the Repository

The `meshery/meshery` repository is large - a full clone is on the order of tens of gigabytes and keeps growing. Two directories account for most of that weight, and neither is needed for the vast majority of contributions:

- **`models/`** - the generated [Meshery Model](/pr-preview/pr-24095/concepts/logical/models/) registry (400+ models). A typical build of the server, UI, or CLI only needs the `meshery-core` and `kubernetes` models.
- **`docs/static/v*/`** - archived snapshots of previously released documentation. Only the latest snapshot is needed to build and preview the current docs.

Unless you are specifically working on the model registry or on an archived documentation version, **clone sparsely**. You still get the full commit history and everything required to build the server, UI, CLI, and current docs, while skipping gigabytes of generated content.

### Recommended: sparse clone

First [fork](https://github.com/meshery/meshery/fork) the repository, then replace `USERNAME` and `FORKED-PROJECT` with your GitHub username and fork name:

```bash
# 1. Blobless partial clone - fetches history metadata, not every file's contents
git clone --no-checkout --filter=blob:none https://github.com/USERNAME/FORKED-PROJECT.git
cd meshery
git remote add upstream https://github.com/meshery/meshery.git

# 2. Check out everything EXCEPT the bulky generated directories:
#      - every model except meshery-core and kubernetes
#      - archived docs snapshots (keep only the latest, docs/static/v0.9)
git sparse-checkout set --no-cone \
  '/*' \
  '!/models/*' \
  '/models/meshery-core/' \
  '/models/kubernetes/' \
  '!/docs/static/v0.8/'

# 3. Populate the working tree
git checkout master
```

This yields a working tree a fraction of the size of a full clone. `--filter=blob:none` makes it a [partial clone](https://git-scm.com/docs/partial-clone), so Git transparently fetches any excluded file on demand if you ever check one out - nothing is permanently lost, and `git log`/`git blame` keep working across the full history.

The `!/models/*` pattern excludes the *contents* of `models/` (not the directory itself), which is what lets the two following lines re-include only the `meshery-core` and `kubernetes` models. For the docs, exclude the older archived snapshots under `docs/static/` and keep the newest (currently `docs/static/v0.9`); add another `!` line for each older version as new snapshots are archived.

### Re-including an excluded directory

Working on the model registry, or on an older archived docs version? Re-include just what you need at any time:

```bash
# Bring back the full model registry
git sparse-checkout add '/models/'

# ...or a single archived docs version
git sparse-checkout add '/docs/static/v0.8/'
```

Or restore the full working tree entirely:

```bash
git sparse-checkout disable
```

<div class="alert alert-info" role="alert"><div class="h4 alert-heading" role="heading">Requires a recent Git</div>


`git sparse-checkout` requires Git 2.25+, and partial clone (`--filter=blob:none`) works best with Git 2.27+. Check with `git --version`.
</div>


## Not sure where to start?

<details>
<summary>Follow these steps and you'll be right at home.</summary>
<ol>
  <li>See the <a href="https://meshery.io/community">Newcomers Guide</a> for how, where, and why to contribute.</li>
  <li>Sign up for a <a href="https://meshery.io/community#meshmates"><em>MeshMate</em></a> to find the perfect Mentor to help you explore the projects and find your place in the community.</li>
  <li><strong>Familiarize</strong> yourself with the broader set of projects in Meshery's ecosystem, including the <a href="https://github.com/meshery-extensions">meshery-extensions</a> repositories (the <a href="https://meshery.io/community">community handbook</a> is a helpful resource): Spend time understanding each of the initiatives through high-level overviews available in the community drive and through discussions with your MeshMate.</li>
  <li><strong>Identify</strong> your area of interest: Use the time with your MeshMate to familiarize yourself with the architecture and technologies used in the projects. Inform your MeshMate of your current skills and what skills you aim to develop.</li>
  <li><strong><a href="https://play.meshery.io">Play with Meshery</a></strong>: Put on your user hat and walk-through all of Meshery’s features and functions as a user.</li>
  <li><strong>Build Meshery Server and UI</strong>: Confirm that you have a usable development environment. See <a href="#specific-contribution-guides">Guides</a> below.</li>
  <li><strong>Discuss</strong> with the community by engaging in the <a href="https://discuss.meshery.io/">discussion forum</a>.</li>
  <li><strong>Contribute</strong> by grabbing any open issue with the <a href="https://github.com/search?q=is%3Aopen+is%3Aissue+archived%3Afalse+org%3Ameshery+org%3Ameshery-extensions+org%3Aservice-mesh-performance+org%3Aservice-mesh-patterns+label%3A%22help+wanted%22&type=issues">help-wanted label</a> and jump in. If needed, create a <a href="https://github.com/meshery/meshery/issues/new/choose">new issue</a>. All <a href="https://github.com/meshery/meshery/pulls">pull requests</a> should reference an open issue. Include keywords in your pull request descriptions, as well as commit messages, to <a href="https://help.github.com/en/github/managing-your-work-on-github/closing-issues-using-keywords">automatically close issues in GitHub</a>.</li>
  <li><strong>Fill in</strong> a <a href="https://meshery.io/newcomers">community member form</a> to gain access to community resources.</li>
</ol>
</details>

## Specific Contribution Guides

Here is a complete list of all of Meshery's contributing guides from Server to UI to CLI to Extensions and so on.

---

Section pages:

- [Contributing to Operator, MeshSync & Broker Settings](/pr-preview/pr-24095/project/contributing/contributing-controllers-config/): The complete surface of settings that govern Meshery Operator, MeshSync and Meshery Broker behavior: where each one lives, what it writes, what reads it, and what happens when it changes.
- [Build & Release (CI)](/pr-preview/pr-24095/project/contributing/build-and-release/): Details of Meshery's build and release strategy.
- [Build Environment Gotchas](/pr-preview/pr-24095/project/contributing/contributing-build-environment/): Local build and dependency failures in Meshery that report a cause other than their own.
- [Contributing Relationships - Testing and Best Practices](/pr-preview/pr-24095/project/contributing/contributing-relationships-best-practices/): How to test, validate, and debug relationship definitions before contribution
- [Contributing to Meshery Adapters](/pr-preview/pr-24095/project/contributing/contributing-adapters/): How to contribute to Meshery Adapters
- [Contributing to Meshery CLI](/pr-preview/pr-24095/project/contributing/cli/): How to contribute to the Meshery CLI and its commands.
- [Contributing to Meshery Docker Extension](/pr-preview/pr-24095/project/contributing/contributing-docker-extension/): How to contribute to Meshery Docker Extension
- [Contributing to Meshery Docs](/pr-preview/pr-24095/project/contributing/contributing-docs/): How to contribute to Meshery Docs.
- [Contributing to Meshery Models](/pr-preview/pr-24095/project/contributing/models/): How to contribute to the Meshery models and their components.
- [Contributing to Meshery Operator](/pr-preview/pr-24095/project/contributing/contributing-operator/): How to build, test, and deploy Meshery Operator from source
- [Contributing to Meshery Policies](/pr-preview/pr-24095/project/contributing/contributing-policies/): How to contribute to Meshery Policies
- [Contributing to Meshery Schemas](/pr-preview/pr-24095/project/contributing/contributing-schemas/): How to contribute to Meshery Schemas
- [Contributing to Meshery Server](/pr-preview/pr-24095/project/contributing/contributing-server/): How to contribute to Meshery Server
- [Contributing to Meshery Server Events](/pr-preview/pr-24095/project/contributing/contributing-server-events/): Guide is to help backend contributors send server events using Golang.
- [Contributing to Meshery UI](/pr-preview/pr-24095/project/contributing/ui/): How to contribute to the Meshery UI and its components.
- [Contributing to Meshery using git](/pr-preview/pr-24095/project/contributing/contributing-gitflow/): How to contribute to Meshery using git
- [Contributing to MeshSync](/pr-preview/pr-24095/project/contributing/contributing-meshsync/): How to build, run, test, and contribute to MeshSync from source
- [Go Lint Rules](/pr-preview/pr-24095/project/contributing/contributing-lint/): The repo-specific Go lint gates Meshery enforces in CI, what each one is protecting, and what to do when one fires.
- [How to write MeshKit compatible errors](/pr-preview/pr-24095/project/contributing/contributing-error/): How to declare errors in Meshery components.
- [HTTP Error Response Contract](/pr-preview/pr-24095/project/contributing/error-contract/): How Meshery Server responds to error conditions over HTTP, and how clients should parse those responses.
- [Meshery Documentation Structure and Organization](/pr-preview/pr-24095/project/contributing/contributing-docs-structure/): Audience, high-Level outline & information architecture for Meshery Documentation
- [Setting up Meshery Development Environment on Windows](/pr-preview/pr-24095/project/contributing/meshery-windows/): How to set up Meshery Development Environment on Windows
