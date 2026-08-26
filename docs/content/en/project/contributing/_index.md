---
title: Contributing
display_title: false
categories: [contributing]
display-suggested-reading: false
description: How to contribute to the Meshery project and any of its components.
---
# Contributing 

Please do! Thanks for your help! 🎈 Meshery is community-built and welcomes collaboration. Contributors are expected to adhere to the [CNCF's Code of Conduct](https://github.com/meshery/meshery/blob/master/CODE_OF_CONDUCT.md).

## General Contribution Flow

Meshery and it's various architectural components are written in different languages, including Golang, Javascript (React.js and Next.js) To make building, testing, and the experience of contributing consistent across all Meshery components, a `Makefile` is included in the every repository. These `make` targets are what you will use to build, run, test, and document.

To contribute to Meshery, please follow this basic fork-and-pull request [gitflow]({{< ref "project/contributing/contributing-gitflow.md" >}}).

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

## Cloning the Repository

The `meshery/meshery` repository is large - a full clone is on the order of tens of gigabytes and keeps growing. Two directories account for most of that weight, and neither is needed for the vast majority of contributions:

- **`models/`** - the generated [Meshery Model]({{< ref "concepts/logical/models/index.md" >}}) registry (400+ models). A typical build of the server, UI, or CLI only needs the `meshery-core` and `kubernetes` models.
- **`docs/static/v*/`** - archived snapshots of previously released documentation. Only the latest snapshot is needed to build and preview the current docs.

Unless you are specifically working on the model registry or on an archived documentation version, **clone sparsely**. You still get the full commit history and everything required to build the server, UI, CLI, and current docs, while skipping gigabytes of generated content.

### Recommended: sparse clone

```bash
# 1. Blobless partial clone - fetches history metadata, not every file's contents
git clone --no-checkout --filter=blob:none https://github.com/meshery/meshery.git
cd meshery

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

{{% alert color="info" title="Requires a recent Git" %}}
`git sparse-checkout` requires Git 2.25+, and partial clone (`--filter=blob:none`) works best with Git 2.27+. Check with `git --version`.
{{% /alert %}}

## Not sure where to start?

<details>
<summary>Follow these steps and you'll be right at home.</summary>
<ol>
  <li>See the <a href="https://meshery.io/community">Newcomers Guide</a> for how, where, and why to contribute.</li>
  <li>Sign up for a <a href="https://meshery.io/community#meshmates"><em>MeshMate</em></a> to find the perfect Mentor to help you explore the projects and find your place in the community.</li>
  <li><strong>Familiarize</strong> yourself with the broader set of projects in Meshery's ecosystem, including the <a href="https://github.com/meshery-extensions">meshery-extensions</a> repositories (the <a href="https://meshery.io/community">community handbook</a> is a helpful resource): Spend time understanding each of the initiatives through high-level overviews available in the community drive and through discussions with your MeshMate.</li>
  <li><strong>Identify</strong> your area of interest: Use the time with your MeshMate to familiarize yourself with the architecture and technologies used in the projects. Inform your MeshMate of your current skills and what skills you aim to develop.</li>
  <li><strong><a href="https://play.meshery.io">Play with Meshery</a></strong>: Put on your user hat and walk-through all of Meshery’s features and functions as a user.</li>
  <li><strong>Build Meshery Server and UI</strong>: Confirm that you have a usable development environment. See <a href="#contributing-guides">Guides</a> below.</li>
  <li><strong>Discuss</strong> with the community by engaging in the <a href="https://discuss.meshery.io/">discussion forum</a>.</li>
  <li><strong>Contribute</strong> by grabbing any open issue with the <a href="https://github.com/issues?q=is%3Aopen+is%3Aissue+archived%3Afalse+org%3Ameshery+org%3Ameshery-extensions+org%3Aservice-mesh-performance+org%3Aservice-mesh-patterns+label%3A%22help+wanted%22+">help-wanted label</a> and jump in. If needed, create a <a href="https://github.com/meshery/meshery/issues/new/choose">new issue</a>. All <a href="https://github.com/meshery/meshery/pulls">pull requests</a> should reference an open issue. Include keywords in your pull request descriptions, as well as commit messages, to <a href="https://help.github.com/en/github/managing-your-work-on-github/closing-issues-using-keywords">automatically close issues in GitHub</a>.</li>
  <li><strong>Fill in</strong> a <a href="https://meshery.io/newcomers">community member form</a> to gain access to community resources.</li>
</ol>
</details>

## Specific Contribution Guides

Here is a complete list of all of Meshery's contributing guides from Server to UI to CLI to Extensions and so on.
