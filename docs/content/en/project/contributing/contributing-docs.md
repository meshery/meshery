---
title: Contributing to Meshery Docs
description: How to contribute to Meshery Docs.
categories: [contributing]
---

{{% alert color="info" title="Important" %}}
Before contributing, please review the [Documentation Contribution Flow](https://github.com/meshery/meshery/blob/master/CONTRIBUTING.md#documentation-contribution-flow).
{{% /alert %}}

Use the following steps to set up your development environment depending on your Operating System of choice, fork the repository, install dependencies and run the site locally.  You can then make changes, test locally, sign-off and commit, and push the changes made for review.

{{% alert color="info" title="Meshery Documentation Design Specification" %}}
See the [Meshery Documentation Design Specification](https://docs.google.com/document/d/17guuaxb0xsfutBCzyj2CT6OZiFnMu9w4PzoILXhRXSo/edit#) which serves to provide an overview of the tooling and approach used to create Meshery's documentation and its information architecture.
{{% /alert %}}

## Documentation Framework

Meshery documentation is built using the following components:

- Framework - [Hugo](https://gohugo.io) (Extended)
- Themes and modules (imported as Hugo modules):
  - [Docsy](https://www.docsy.dev) — fully mounted
  - [Layer5io Docs](https://github.com/layer5io/docs) — manual mounting required
- Repo - [https://github.com/meshery/meshery/tree/master/docs](https://github.com/meshery/meshery/tree/master/docs)
- Site - [https://docs.meshery.io](https://docs.meshery.io)

## Set up your development environment

{{% alert color="info" title="Hugo" %}}
The Meshery Docs site is built using Hugo - a fast static site generator. Hugo is installed automatically as an npm dependency, so you only need Node.js and Go installed on your system.
{{% /alert %}}

### Prerequisites

The following tools are required to build and run the Meshery documentation site locally:

- [Node.js](https://nodejs.org/) (LTS recommended)
- [Go](https://go.dev/dl/) (required for Hugo modules)

{{% alert color="light" title="Note" %}}
In case of any installation issues, use the [discussion forum](https://discuss.meshery.io/).
{{% /alert %}}

### Get the code

- Fork and then clone the [Meshery repository](https://github.com/meshery/meshery)

  <pre class="codeblock-pre"><div class="codeblock">
  <code class="clipboardjs">git clone https://github.com/YOUR-USERNAME/meshery
  </code>
  </div></pre>

- Install dependencies

  <pre class="codeblock-pre"><div class="codeblock">
  <code class="clipboardjs">make docs:setup
  </code>
  </div></pre>

  This runs `npm install`, which installs Hugo Extended and all other required dependencies.

### Serve the site

- Serve the site locally.

  <pre class="codeblock-pre"><div class="codeblock">
  <code class="clipboardjs">make site
  </code>
  </div></pre>

  This runs `hugo --cleanDestinationDir -e dev server -DFE`, which serves the site with draft and future content enabled. The site will be available at `http://localhost:1313`.

- To build the site without serving:

  <pre class="codeblock-pre"><div class="codeblock">
  <code class="clipboardjs">make docs:build
  </code>
  </div></pre>

### Using Docker

If you have Docker and `make` installed on your system, then you can serve the site locally. This doesn't require the need to install _Hugo_, _Node.js_, or _Go_ on your system.

{{% alert color="info" %}}
This may not work in Windows.
{{% /alert %}}

Run the following command from the _docs_ folder.

```bash
make docker
```

This runs `docker compose watch`, which builds and serves the site inside a container.

### Using Gitpod
- Ensure you have an account on <a href="https://www.gitpod.io/" class="meshery-light">Gitpod</a> and add the <a href="https://chromewebstore.google.com/detail/gitpod/dodmmooeoklaejobgleioelladacbeki" class="meshery-light">Gitpod extension</a> to your browser.
- Open your forked Meshery repository on GitHub.
- Click on the "**Gitpod**" button in the top right corner of the repository page (it is only visible with the <a href="https://chromewebstore.google.com/detail/gitpod/dodmmooeoklaejobgleioelladacbeki" class="meshery-light">Gitpod browser extension</a> installed).

{{% alert color="dark" title="About Gitpod" %}}
Gitpod will automatically clone and open the repository for you in a browser-based version of Visual Studio Code. It will also automatically build the project for you on launch, comes with Docker and other tools pre-installed, making it one of the fastest ways to spin up a development environment for [Meshery](https://github.com/meshery/meshery).
{{% /alert %}}

- After opening the project on Gitpod, install dependencies

  <pre class="codeblock-pre"><div class="codeblock">
  <code class="clipboardjs">make docs:setup
  </code>
  </div></pre>

- Serve the site locally.

  <pre class="codeblock-pre"><div class="codeblock">
  <code class="clipboardjs">make site
  </code>
  </div></pre>

  This runs `hugo --cleanDestinationDir -e dev server -DFE`, which serves the site with draft and future content enabled. The site will be available at `http://localhost:1313`.

You should be able to access the site on port `1313`. If you want to access it in your localhost read the [docs for port-forwarding using ssh](https://www.gitpod.io/docs/configure/workspaces/ports#local-port-forwarding-via-ssh).

### Using GitHub Codespaces

- Ensure you have an account on <a href="https://www.github.com/" class="meshery-light">GitHub</a>
- Open your forked Meshery repository on GitHub.
- Click on the "**Code**" button of your forked repository it will give you the option to open the project in GitHub Codespaces.

{{% alert color="dark" title="About GitHub Codespaces" %}}
GitHub Codespaces will automatically clone and open the repository for you in a browser-based version of Visual Studio Code. It comes with pre-installed tools which are quite helpful to spin up the development environment for [Meshery](https://github.com/meshery/meshery).
{{% /alert %}}

- After opening the project on GitHub Codespaces, install dependencies

  <pre class="codeblock-pre"><div class="codeblock">
  <code class="clipboardjs">make docs:setup
  </code>
  </div></pre>

- Serve the site locally.

  <pre class="codeblock-pre"><div class="codeblock">
  <code class="clipboardjs">make site
  </code>
  </div></pre>

  This runs `hugo --cleanDestinationDir -e dev server -DFE`, which serves the site with draft and future content enabled. The site will be available at `http://localhost:1313`.

You should be able to access the site on port `1313`. If you want to access it in your localhost just click the `code` button on your forked repository and select open with Visual Studio Code under your GitHub Codespace this will launch a GitHub Codespace instance in your local machine and connects with the remote GitHub Codespace environment after that run `make docs:setup && make site` and it will start the development server on port `1313`.

### Make Necessary Changes
- Make changes as required based on the issue you are solving.
- Ensure to verify that your changes reflect correctly by serving the site locally using `make site`.

{{% alert color="light" title="Note" %}}
If the issue requires making new doc page that replaces the old page, please don't forget to add a redirect link on the old page. In Hugo, you can add an `aliases` field in the frontmatter of the new page to redirect from the old URL.
{{% /alert %}}


### Create a Pull Request

- After making changes, don't forget to commit with the sign-off flag (-s)!

  <pre class="codeblock-pre"><div class="codeblock">
  <code class="clipboardjs">git commit -s -m "my commit message w/signoff"
  </code>
  </div></pre>

- Once all changes have been committed, push the changes.

  <pre class="codeblock-pre"><div class="codeblock">
  <code class="clipboardjs">git push origin [branch-name]
  </code>
  </div></pre>

- Then on GitHub, navigate to the [Meshery repository](https://github.com/meshery/meshery) and create a pull request from your recently pushed changes!

---

- _See the [Meshery Documentation Google Doc](https://docs.google.com/document/d/17guuaxb0xsfutBCzyj2CT6OZiFnMu9w4PzoILXhRXSo/edit) for additional reference._

## Using the features of Meshery Docs

### Linking to Pages and Assets

**Internal content pages** use Hugo's `ref` shortcode:

```md
[About]({{</* ref "about.md" */>}})
```

**Page-bundle assets** use standard relative paths. For assets adjacent to `index.md` or `_index.md`, or inside that page bundle:

```md
![Architecture]({{< static "images/meshery-architecture.webp" >}})
<img src="{{< static "images/meshery-architecture.webp" >}}" />
```

**Global shared assets** under `static/` use the `static` shortcode:

```md
{{</* static "images/logo.png" */>}}
```

Do not use raw root-relative page links like `[About](/about/)`. External links (starting with `http://`, `https://`, `mailto:`, `tel:`, `#`) should be kept unchanged.

### Clipboard Feature

Most popular clipboard plugins like Clipboard JS require the manual creation of a new ID for each code snippet. A different approach is used here. For code snippets, we either use html tags or markdown in the following manner:

```
<pre class="codeblock-pre"><div class="codeblock">
<code class="clipboardjs">code_snippet_here</code>
</div></pre>
```

You can also use the `code` shortcode created for this feature to make it easy to use. It can be used by passing the code snippet to it.

If the code snippet does not contain any special characters then, it can be used directly as below:
```
{{</* code code="code_snippet_here" */>}}
```

If the code snippet has special characters then use backtick delimiters:

```
{{</* code code=`code_snippet_here` */>}}
```

Don't use `code` shortcode when the snippet is in between an ordered list, this breaks the order and next item in the list will start numbering from 1. Instead, use `<pre class="codeblock-pre">...</pre>` method described above.

**A full block:**
```
```code snippet```
```

**Inline formatting:**

{{< code code="`code snippet`: `code snippet`" >}}

**Language specific:**

{{< code code=`(language name)
code snippet
` >}}

Whenever the code tags are detected, the clipboard javascript file is automatically loaded. Each code element is given a custom id and a clipboard-copy icon to copy the content.

## Documentation Contribution Flow Summary

{{% alert color="light" title="Note" %}}
For contributing `mesheryctl` reference section, refer [Contributing CLI]({{< ref "project/contributing/cli" >}})
{{% /alert %}}


The following is a concise summary of the steps to contribute to Meshery documentation.

1. Create a fork, if you have not already, by following the steps described [here](CONTRIBUTING-gitflow.md)
2. In the local copy of your fork, navigate to the docs folder.
   `cd docs`
3. Create and checkout a new branch to make changes within
   `git checkout -b <my-changes>`
4. Edit/add documentation.
   `vi <specific page>.md`
5. Add redirect link on the old page (only when a new page is created that replaces the old page). In Hugo, use the `aliases` frontmatter field.
6. Run site locally to preview changes.
   `make site`
7. Commit, [sign-off](#commit-signing), and push changes to your remote branch.
   `git push origin <my-changes>`
8. Open a pull request (in your web browser) against the repo: https://github.com/meshery/meshery.

### Navigation Table of Contents in Sidebar (toc)

Sidebars use a data file at `data/toc.yml` to create a table of contents. It is written in the following manner:

```
- title: Group 1
  url: group1
  links:
    - title: Thing 1
      url: group1/thing1
    - title: Thing 2
      url: group1/thing2
    - title: Thing 3
      url: group1/thing3
```

The output of the code snippet would be:

```
    Group 1
      Thing 1
      Thing 2
      Thing 3
```

In this example, **Group 1** is a ``parent`` section within the Table of Contents, and **Thing 1**, **Thing 2**, and **Thing 3** are ``children`` of **Group 1**. This is a direct hierarchy with a single level of children under the parent.

``Grandchildren`` are added in the sidebars in the following manner:

```
- title: Group 1
  url: group1
  links:
    - title: Thing 1
      url: group1/thing1
      children:
        - title: Subthing 1.1
          url: group1/thing1/subthing1_1
        - title: Subthing 1.2
          url: group1/thing1/subthing1_2
    - title: Thing 2
      url: group1/thing2
    - title: Thing 3
      url: group1/thing3
```

The output of the code snippet would be:

```
    Group 1
      Thing 1
       Subthing 1.1
       Subthing 1.2
      Thing 2
      Thing 3
```

Here, **Subthing 1.1** and **Subthing 1.2** are the ``grandchildren`` of **Thing 1**.

The sidebar template in `layouts/partials/sidebar.html` reads from `data/toc.yml` and supports up to four levels of navigation hierarchy.

- **Parent**: It serves as a top level category for related content.
- **Children** (`links`): They are immediate subsections or topics that fall under the parent section.
- **Grandchildren** (`children`): They are nested under a link and provide a more detailed breakdown of information within the child section. Grandchildren are used to organize content further, offering a more detailed structure for a specific topic.
- **Great-grandchildren** (`grandchildren`): An additional level of nesting for even finer-grained organization.

These sections create a hierarchical and organized navigation experience for readers.

### Alerts

{{% alert color="info" title="What is an alert?" %}}
An alert is a box that can stand out to indicate important information. You can choose from levels success, warning, danger, info, and primary. This example is an info box, and the code for another might look like this:
{{% /alert %}}

{{< code code=`{{% alert color="info" title="Here is another!" %}}
{{% /alert %}}` >}}

Just for fun, here are all the types:

{{% alert color="warning" %}}
This is a warning
{{% /alert %}}
{{% alert color="danger" %}}
This alerts danger!
{{% /alert %}}
{{% alert color="success" %}}
This alerts success
{{% /alert %}}
{{% alert color="info" %}}
This is useful information.
{{% /alert %}}
{{% alert color="primary" %}}
This is a primary alert
{{% /alert %}}
{{% alert color="secondary" %}}
This is a secondary alert
{{% /alert %}}

{{% alert color="light" %}}
This is a light alert
{{% /alert %}}
{{% alert color="dark" %}}
This is a dark alert
{{% /alert %}}

#### alert shortcode

Meshery Docs uses the `alert` shortcode (provided by the Docsy theme) to provide consistent formatting for notes, warnings, and various informative callouts intended for the readers.

To use the `alert` shortcode feature in our documentation include the following code:

{{< code code=`{{% alert color="info" title="Here is another!" %}}
{{% /alert %}}` >}}


Other supported alert colors include `warning`, `danger`,`success`,`primary`, `secondary`, `light`, `dark` .


### Image Handling

Control image display in documentation pages using the following methods.

#### Default Image Size in Markdown Format

Using the Markdown image syntax:

{{< code code="[![Image Title](../../../assets/img/your-image.png)](../../../assets/img/your-image.png)" >}}

This renders as:

{{< code code=`<a href=\"/assets/img/your-image.png\">
    <img src=\"/assets/img/your-image.png\" alt=\"Image Title\">
</a>` >}}

**Effect:**
- Image is displayed at its original resolution, limited by global CSS (`max-width: 90vw; max-height: 90vh; height: auto; width: auto;`)
- On larger screens (`min-width: 1200px`), image width can be up to `1200px`
- Clickable, opening in Lightbox if the format is supported

#### Custom Image Size (for Control)

If you need to specify dimensions, use:

{{< code code=`<a href=\"/assets/img/your-image.png\">
    <img src=\"/assets/img/your-image.png\" style=\"width:500px; height:auto;\" alt=\"Image Title\">
</a>` >}}

**Effect:**
- Image width is fixed at `500px`, maintaining aspect ratio
- If `500px` exceeds `90vw`, it will be constrained to `90vw`
- Clickable, opening in Lightbox for full-size viewing

### Quotes

You can include block quotes to emphasize text.

> Here is an example. Isn't this much more prominent to the user?

## Development

### `if` conditional

In Hugo templates, conditionals are written using Go template syntax:

{{< code code=`{{ if eq .Params.title "Awesome Shoes" }}
  These shoes are awesome!
{{ end }}` >}}

If the condition is true, the output would be:

```
    These shoes are awesome!
```

### `range` loop

The `range` action iterates over a collection. It is written in the following manner:

{{< code code=`{{ range .Pages }}
  {{ .Title }}
{{ end }}` >}}

The output produced by the above code snippet would list the titles of all pages in the current section.

### Comment

In Hugo templates, comments are written using Go template comment syntax:

{{< code code=`{{/* This is a comment */}}` >}}

Any content inside `{{/* */}}` will not be rendered in the output.

### Partial

Partials are used to include reusable template fragments. They are written in the following manner:

{{< code code=`{{ partial "header.html" . }}` >}}

### Shortcode

Shortcodes are reusable content snippets that can be used in Markdown files. They are written in the following manner:

{{< code code=`{{</* shortcode-name param="value" */>}}` >}}

### Variable

In Hugo templates, variables are defined using the `:=` operator:

{{< code code=`{{ $variable1 := true }}` >}}

### Hugo modules

Meshery Docs is built from Hugo modules — self-contained, versioned bundles of layouts, assets, and configuration that the site imports from other repositories rather than keeping in its own tree. Because each module is pulled in by version, shared functionality is defined once upstream and reused here, and updates arrive by bumping a single version reference instead of editing copied files.

Modules are declared in `hugo.toml` under the `[module]` block. Each module is listed as a `[[module.imports]]` entry with its `path` set to the module's repository. When a module's files should map directly onto the site, that import is all Hugo needs. When only certain files are wanted, the import adds `[[module.imports.mounts]]` entries, each pairing a `source` path inside the module with a `target` path in the site, so only the mounted paths are brought in.

In a mount, `source` is the location of the file *within the imported module*, and `target` is the location it should appear at *within the site*. Hugo does not copy these files into the project; instead it composes every module, including the site itself, into a single union file system — one virtual tree where each module contributes its files at their target paths. A mount is what places a module's `source` file at a given `target` in that shared tree, and from Hugo's perspective the file then sits there as if it were part of the site. Because all modules share one tree, targets can overlap: when more than one module contributes a file to the same target, precedence decides which one is used, with the site's own files taking priority over imports and earlier imports taking priority over later ones.

This precedence is what makes overriding possible. To replace an imported file, the site adds its own file at the same target path; because the site's files outrank imports in the union tree, its version wins and Hugo uses it in place of the module's. This is especially useful with fully imported modules like Docsy, where the entire theme is mapped onto the site: rather than forking the theme to change a single layout, partial, or stylesheet, the site can supply a file of the same name at the matching path and override just that one piece, while everything else continues to come from the module and stays in sync with upstream updates.

The site imports three modules:

- **`github.com/google/docsy`** — the Docsy theme, which provides the site's base look and feel: its layouts, partials, and styling. It is imported directly, so everything the theme defines is available to the site.
- **`github.com/google/docsy/dependencies`** — a companion module that supplies the front-end dependencies the Docsy theme relies on, such as Bootstrap and its other bundled assets. It too is imported directly.
- **`github.com/layer5io/docs`** — the shared Layer5 docs repository. Its own configuration is ignored, and instead of importing the module wholesale, selected partials and assets are mounted individually into the site. These bring in shared pieces of markup, scripts, and styles that layer onto the existing structure to add functionality without altering the site's own layouts.

Because these assets are shared, prefer using the upstream ones over creating equivalent assets inside this repository.

{{% alert color="warning" title="Contribute shared assets upstream" %}}
Any change that involves the layer5io docs assets — whether adding new functionality to the docs or making quality-of-life improvements to the shared assets — should be made in the [`layer5io/docs`](https://github.com/layer5io/docs) repository rather than here, and then imported downstream. Making the change upstream keeps a single source of truth and ensures every site that consumes the module, both Meshery Docs and Layer5's own docs, receives the same improvement instead of the work living in one repo and drifting from the other.
{{% /alert %}}

### Layer5io Mounts

#### Related Reading

The `related-reading.html` partial renders a list of links to other content covering similar material, shown at the end of a page. It is mounted in from the [`layer5io/docs`](https://github.com/layer5io/docs) module and invoked from the default layout, so the block appears across the site without any per-page configuration:

```
{{ partial "related-reading.html" . -}}
```

The context (`.`) passed in the call is the current page, which the partial uses as the basis for its recommendations. These are drawn from [Hugo's related-content engine](https://gohugo.io/content-management/related/), which scores the site's content against the current page and returns the closest matches. That result is optionally constrained to the types and sections listed under `Site.Params.relatedReading`, filtered to exclude video entries, and capped at six.

When the engine yields no matches, the partial resolves recommendations from taxonomy instead: it gathers content sharing the current `categories` or `tags`, applies the same type, section, and video filters, and again caps the result at six. The block is omitted entirely when neither path produces a result.

This partial is the single source of recommendations across the site. Hardcoded lists of links should not be written into a page to stand in for it; recommendations are meant to be resolved from taxonomy so they stay consistent site-wide and update automatically as content is added. To have a page participate, assign it relevant `categories` and `tags` in its frontmatter:

```
---
title: Managing Meshery Deployments
categories: [deployment]
tags: [kubernetes, operator]
---
```

Pages sharing these values become candidates for one another's recommendations, so aligning taxonomy across related content is what surfaces it here.

#### Resizable layouts

Resizable layouts allow a section of the page to be dragged to a new width, with the size retained across visits. The behavior arrives through the [`layer5io/docs`](https://github.com/layer5io/docs) module, which mounts the markup, script, and styles it depends on so it drops into a layout as a complete unit.

A region is wrapped between `resizable-start.html` and `resizable-end.html`; the content placed between them becomes adjustable. The opening call configures it through a parameter dictionary:

```
{{ partial "resizable-start.html" (dict "side" "right" "key" "sidebar" "min" 200 "max" 480 "default" 260) }}
  <aside class="td-sidebar d-print-none">
    {{ partial "sidebar.html" . }}
  </aside>
{{ partial "resizable-end.html" . }}
```

The dictionary accepts the following:

- `side` — the edge the drag handle sits on, either `right` or `left`.
- `key` — a unique identifier under which the width is stored and restored; each region must supply its own.
- `min` / `max` — the lower and upper bounds of the width, in pixels.
- `default` — the width applied when nothing has been stored yet.
- `class` — optional additional classes placed on the wrapper.

{{% alert color="info" title="Companion partials and assets" %}}
Alongside the wrappers, two companion partials are included once in the layout: `resizable-head-script.html` in the `<head>` to apply saved widths before the first paint, and `resizable-script.html` near the end of the `<body>` to load the behavior. For either to work, the mounted assets — `assets/js/resizable.js` and `assets/scss/_resizable.scss` — must also be imported into the project so the script and styles are built into the site.
{{% /alert %}}

Sizing is governed by the component through the `--resizable-w` property, which it sets directly on the element as the region is dragged. Because the value lives inline on the element rather than in a stylesheet, anything that needs to react to the current width has to observe the element to read it, rather than depending on a fixed CSS rule.

The component is built to be reused across sites, so its own rules are not meant to be rewritten. Overriding them locally couples the region to one layout and works against the shared design the component exists to provide.

{{% alert color="info" title="Style through the wrapper, not the component" %}}
To apply site-specific styling, pass a class through the `class` parameter and style that instead. The class lands on the wrapper, so it surrounds the component and layers styling on top of it without touching the component's own rules — keeping the shared asset intact and reusable.
{{% /alert %}}

If an override is truly unavoidable — and it should be treated as a last resort rather than a first option — it belongs in a dedicated `resizer_override` file imported into the project's styles, rather than edited inline against the component. Isolating overrides this way keeps changes to the width variable and the component's rules tracked in one place and maintainable over time, instead of scattering specificity fixes across the layout. It also contains the awkwardness that comes with promoting the component's class to raise its specificity when an override has to win, keeping that complexity out of the shared asset.

### Adding Integration Specific Information to Individual Integration Pages

Integration pages ([example]({{< ref "extensions/models/aws/index.md" >}})) are automatically generated, however, integration specific documentation is often needed.

[modelscustominfo](https://github.com/meshery/meshery/tree/master/docs/data/modelscustominfo) collection holds custom markdown files. Follow these steps:

1. Create a file inside the `modelscustominfo` collection.
2. Ensure the file includes frontmatter with the `title` key, set to the title of the integration.

The content that follows the frontmatter in this file will be automatically parsed and rendered on the integration page

Example:

file: `modelscustominfo/aad.md`

```
---
title: Azure Active Directory (AAD)
---
<h2>Azure Active Directory (AAD)</h2>

```

In this example, the heading "<b>Azure Active Directory</b>" will be displayed on the integration page:
[Azure Active Directory Integration Page]({{< ref "extensions/models/aad-pod-identity/index.md" >}})

### Editable Intra-page Table of Contents Toolbar

Control the display of this intra-page navigator with either page level or layout level frontmatter variables:

`display-toolbar`

Set to `true` (make "editable" toolbar visible) or `false` (hide "editable" toolbar)

Two helpful resources:

1. Hugo Docs - [https://gohugo.io/documentation/](https://gohugo.io/documentation/)
2. Go Template Docs - [https://pkg.go.dev/text/template](https://pkg.go.dev/text/template)