---
title: Meshery Documentation Design Specification
categories: [project]
description: Overview of Meshery Docs tooling, architecture, and release strategy.
---

## Purpose

This document serves to provide an overview of the tooling and approach used to create Meshery's documentation.

## Contribution

Please contribute. The Meshery documentation uses GitHub Pages to host the docs site. The process of contributing to documentation looks something like this:

1. [Fork the repo](https://help.github.com/en/github/getting-started-with-github/fork-a-repo) and get a local copy of the documentation.

	```bash
	git clone https://github.com/YOUR-USERNAME/meshery
	```

2. Navigate to the docs folder.

	```bash
	cd docs
	```

3. Create and checkout a new branch to make changes within.

	```bash
	git checkout -b bug/<my-username>/<issue#>
	```

4. Edit or add documentation.

	```bash
	vi <specific page>.md
	```

5. Run the site locally to preview changes.

	```bash
	make site
	```

6. Commit, sign off, and push changes to your remote branch.

	```bash
	git push origin <my-changes>
	```

7. Open a pull request against the `master` branch on [https://github.com/meshery/meshery](https://github.com/meshery/meshery).

## Framework

Meshery documentation is made of these components:

- Framework: [Hugo](https://gohugo.io)
- Theme: [https://github.com/google/docsy](https://github.com/google/docsy)
- Repo: [https://github.com/meshery/meshery/tree/master/docs](https://github.com/meshery/meshery/tree/master/docs)
- DNS: [https://docs.meshery.io/](https://docs.meshery.io/)

## Types of Documentation

There is a spectrum of documentation produced in the process of managing Layer5's open source projects. In general, we automate the top of this spectrum to the extent possible. Commonly, our project needs more of the bottom, content-rich, of the spectrum.

### Spectrum

- Changelogs: Ad nauseum list every change made from prior release. Essentially what ReleaseDrafter does.
- Release Notes: A curated, bulleted list of highlights, summarized and categorized in human-readable with some engineering-speak and some issue references. Hopefully, we can get here with ReleaseDrafter.
- Release Announcement: A human-written summary, the bubbles the most significant items to the top, highlights any caveats, for example incompatibility on upgrade, points people to other sources of information, for example upgrade guide, feature blogs, full list of all bug fixed. Should include a graphic or two and links to new docs covering things in-depth. Each significant feature should have a hyperlink to something explaining it.
- Feature Overview: Feature-functionality blog, or net new sections in user docs, or upgrade guide. In-depth review of new functionality or significantly augmented behavior, explaining what the functionality is, why it was added, and how to use it. May incorporate video.
- Full Documentation: User-facing documentation that details the functionality offered by the project.
- Design Specifications: Exhaustive description of the behaviors of the system and functionality.

## Theme

The theme is based on the Docsy Hugo Theme found below:

- [https://github.com/google/docsy](https://github.com/google/docsy)
- [https://www.docsy.dev/docs/](https://www.docsy.dev/docs/)

## Clipboard Feature

Most popular clipboard plugins like Clipboard JS require the manual creation of a new ID for each code snippet. Meshery documentation now uses the `code` shortcode to render clipboard-enabled code blocks consistently.

Use the shortcode in Markdown like this:

<pre><code>&#123;&#123;&lt; code code=`&#123;&#123;&lt; code code =&quot;mesheryctl system login&quot;&#125;&#125;` &gt;&#125;&#125;</code></pre>

For multi-line snippets, pass the block as the `code` parameter:

<pre><code>&#123;&#123;&lt; code code=`
apiVersion: v1
kind: Namespace
metadata:
name: meshery
` &gt;&#125;&#125;</code></pre>

The shortcode renders a wrapped `<code>` element with the `clipboardjs` class, which is what powers the copy-to-clipboard behavior in the docs site.

## Automagic Features Based on Headings

Meshery documentation uses Docsy's built-in heading behavior to generate anchors and the table of contents from page headings.

This theme includes two features based on headings:

- Permalinks: headings get clickable hover anchors, currently enabled for h1 through h5
- TOC: generated automatically from the page's heading structure

The permalinks feature automatically creates hidden anchor links for each heading that display a clickable link on hover. Clicking the link provides you a URL that jumps straight to this section of the docs.

The TOC feature automatically generates a clickable table of contents based on the headings rendered on the page.

## Getting Started

### Features

- User Interaction
On the right side of any page, you'll notice links to edit the page, or open an issue. This ensures that any time you have a question or want to suggest or request a change, you can do so immediately and link directly to the section of interest. The sections on the page also have permalinks so you can link directly to them.

The TOC feature automatically generates a clickable table of contents based on the headings.

<a href="/project/images/docsy-toc.png" class="lightbox-image">
<img src="/project/images/docsy-toc.png" width="50%" /></a>

## Versioning

### Approaches to Solve the Issue

- We build the site and have the folders in build time, which would inflate the repo and is not scalable.
- We use deployment technique to deploy release tags using deployment services like Vercel and Netlify and have a webview into those so simply changing URLs on dropdown would change or replace the URL with its version alternative.

### Research

- The Meshery Docs site is now built with Hugo and the Docsy theme.
- The current docs workflow keeps the latest documentation in the main `docs` source tree and tracks older releases as versioned static output in `/static` directory.

### Approach 1: Build and Publish Versioned Static Output

Process:

1. Build the Hugo docs site for a release.
2. Publish the generated static output under a versioned path such as `v0.x`.
3. Keep the latest docs available at the main docs site while older versions remain accessible by their versioned paths.

For more info: [https://gohugo.io/](https://gohugo.io/)

Pros:

- Easier to keep a stable published snapshot for each release.
- Versioning can be achieved without changing the source content structure.

Cons:

- Versioned static output still needs to be generated and published for each release.
- Keeping old versions available increases storage and maintenance overhead over time.

### Approach 2: Manage Release Using Preview Deployment Managers like Vercel or Netlify

For more info: [https://docs.netlify.com/site-deploys/overview/#deploy-contexts](https://docs.netlify.com/site-deploys/overview/#deploy-contexts)

Pros:

- No change in structure of how everything is done.
- Versioning can be achieved.

Cons:

- Requires configuration for branch-based or release-based deploy previews.
- Custom-domain version switching still needs a clear routing strategy.
- Release branches need to be kept intact and tied to the published version they represent.

### CNAME

The `docs.meshery.io` CNAME points the docs host to GitHub Pages - [https://meshery.github.io/meshery/](https://meshery.github.io/meshery/)

