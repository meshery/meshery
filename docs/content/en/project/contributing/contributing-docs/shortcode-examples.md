---
title: Shortcode Formatting Examples
description: A visual reference guide for shortcodes available in the Meshery documentation theme.
categories: [contributing]
aliases: [/project/contributing/shortcode-examples]
---

The Meshery documentation uses a variety of shortcodes to create rich, formatted content such as callouts, tabs, and generated tables.

This page is a visual reference for documentation contributors. It shows the markdown syntax required to invoke each shortcode, followed by its rendered output, so you can see exactly how a shortcode behaves before using it.

Shortcodes come from two places:

- [`docs/layouts/shortcodes/`](https://github.com/meshery/meshery/tree/master/docs/layouts/shortcodes) — defined in this repository and documented in full below.
- The [Docsy theme](https://www.docsy.dev/docs/adding-content/shortcodes/) — inherited, and summarized under [Inherited Shortcodes](#inherited-shortcodes).

The Meshery documentation also imports `github.com/layer5io/docs` as a Hugo module, but that import is mounted to a single partial in `hugo.toml` and contributes no shortcodes.

---

## Shortcode Delimiters

Hugo shortcodes are invoked with one of two delimiter pairs, and the choice changes how the shortcode's body is treated:

| Delimiter | Body treatment |
|---|---|
| `{{</* name */>}}` | The body is passed to the shortcode as raw text. |
| `{{%/* name */%}}` | The body is rendered as Markdown before the shortcode receives it. |

Shortcodes defined in this repository that take a body render it themselves via `.Page.RenderString`, and therefore expect the `{{</* */>}}` form. The inherited Docsy callouts expect `{{%/* */%}}`. Each entry states which form to use.

The form a shortcode expects is determined by its implementation, not by convention. When an entry here disagrees with what you see on the site, the file in `docs/layouts/shortcodes/` is the source of truth, and the entry needs fixing.

---

## Meshery Shortcodes

Listed alphabetically.

### code

Wraps a snippet in a styled code block with a copy-to-clipboard button, without needing a manually assigned ID per snippet.

**Use when** the reader is expected to copy and run the snippet. For code that is only being read — illustrative output, config fragments under discussion — a standard fenced block is lighter and easier to maintain.

| Parameter | Required | Default | Description |
|---|---|---|---|
| `code` | Yes | — | The snippet to render. HTML-escaped on output. Multi-line values are supported; wrap the value in backticks rather than quotes when the snippet itself contains quotes. |

Takes no body. Use the `{{</* */>}}` form.

#### Syntax

```markdown
{{</* code code="make site" */>}}
```

#### Rendered output

{{< code code="make site" >}}

{{% alert color="warning" title="Note" %}}
Don't use the `code` shortcode when the snippet appears between items of an ordered list — it breaks the numbering, and the item following the snippet restarts at 1.
{{% /alert %}}

### discuss

Renders a static callout linking to the Meshery Discussion Forum.

**Use when** closing a page where readers are likely to have follow-up questions. The text is fixed, so one per page at most.

Takes no parameters and no body.

#### Syntax

```markdown
{{</* discuss */>}}
```

#### Rendered output

{{< discuss >}}

### extension-guide

Renders a responsive card grid from a data file in `docs/data/`. Each card shows an image, a title, and a description.

**Use when** presenting a set of parallel options the reader chooses between — component types, extension points — where a visual grid aids scanning. For a linear list of links, use `section-pages`.

| Parameter | Required | Default | Description |
|---|---|---|---|
| `data_file` | Yes | — | Base name of the file in `docs/data/` supplying the card entries, without extension. |
| `guide_title` | Yes | — | Key in each data entry whose value is used as the card title, and as the image `alt` text. |
| `guide_description` | Yes | — | Key in each data entry whose value is used as the card description. Entries missing this key render a card with no description. |
| `guide_svg` | Yes | — | Key in each data entry whose value is the image filename. Entries missing this key render a card with no image. |
| `guide_assests_folder` | Yes | — | Folder name under `images/`. Lowercased before use, producing `images/<folder>/<filename>`. Note the spelling — `assests` is the parameter name in the implementation. |

Takes no body. Use the `{{</* */>}}` form.

#### Syntax

```markdown
{{</* extension-guide
  data_file="edges"
  guide_title="Edge"
  guide_description="Description"
  guide_svg="SVG"
  guide_assests_folder="shapes"
*/>}}
```

#### Rendered output

{{< extension-guide
  data_file="edges"
  guide_title="Edge"
  guide_description="Description"
  guide_svg="SVG"
  guide_assests_folder="shapes"
>}}

### latest-release

Prints the tag of the most recent non-prerelease entry under the `releases` content type.

**Use when** a page needs to show the current release and would otherwise go stale. Note that it emits an `<h2>`, so the release tag appears in the page's table of contents.

Takes no parameters and no body.

#### Syntax

```markdown
{{</* latest-release */>}}
```

#### Rendered output

{{< latest-release >}}

### model-categories

Renders one `<li>` per model category, each with a parenthetical count of the models in that category, ordered by count descending.

**Use when** listing categories inside an existing list. This shortcode emits list items only — **you must supply the surrounding `<ul>` or `<ol>` yourself**, or the items will render unwrapped.

Counts are derived from the `integrations-category` param on each model page, so they stay accurate as models are added or recategorized. Category names are plain labels rather than links, because the models index is a single alphabetized page with no per-category filtering.

Takes no parameters and no body.

#### Syntax

```markdown
<ul>
  {{</* model-categories */>}}
</ul>
```

#### Rendered output

<ul>
  {{< model-categories >}}
</ul>

### model-count

Prints the current total number of Meshery models, counted from the pages under `extensions/models/`.

**Use when** a sentence quotes the model total and would otherwise need manual updating as models are added.

Takes no parameters and no body. Emits a bare number, so it composes inline.

#### Syntax

```markdown
Meshery currently supports {{</* model-count */>}} models.
```

#### Rendered output

Meshery currently supports {{< model-count >}} models.

### related-discussions

Lists up to eight forum topics matching a tag, plus a link to view all discussions with that tag.

**Use when** an active forum tag maps closely to the page topic. Topics are read from `site.Data.discuss.<tag>`, not fetched at request time — if no data file exists for the tag, only the "view all" link renders, with no error. Output changes as that data is refreshed, so don't write prose that depends on which topics appear.

| Parameter | Required | Default | Description |
|---|---|---|---|
| `tag` | No | `meshery` | Forum tag to match. Must correspond to a key under `site.Data.discuss`. |

Takes no body. Use the `{{</* */>}}` form.

#### Syntax

```markdown
{{</* related-discussions tag="meshery" */>}}
```

#### Rendered output

{{< related-discussions tag="meshery" >}}

### relationships

Inserts a predefined, collapsible set of visual examples of component relationships, grouped by kind and subtype, each linking to the design in the playground.

**Use when** explaining how components relate structurally. The content is fixed — it takes no parameters and cannot be scoped to a subset of relationship kinds.

Takes no parameters and no body.

#### Syntax

```markdown
{{</* relationships */>}}
```

#### Rendered output

{{< relationships >}}

### section-pages

Renders one `<li>` per immediate child page of a section, each linking to the page and appending its `description` from front matter when set.

**Use when** a landing page should index its children and stay correct as pages are added or removed. This shortcode emits list items only — **you must supply the surrounding `<ul>` or `<ol>` yourself**.

| Parameter | Required | Default | Description |
|---|---|---|---|
| `section` | No | The page's own Hugo section (`.Page.Section`) | Section whose children are listed. For nested paths the default is rarely what you want — set it explicitly. |

Pages are sorted by title. A page is excluded when its front matter sets `hide_homepage`, and the section's own landing page is always excluded. The final path segments `quick-start` and `faq` are excluded unconditionally by the implementation. When nothing matches, a single italic "No pages found in this section." item renders in place of the list.

Takes no body. Use the `{{</* */>}}` form.

#### Syntax

```markdown
<ul>
  {{</* section-pages section="concepts/architecture" */>}}
</ul>
```

#### Rendered output

<ul>
  {{< section-pages section="concepts/architecture" >}}
</ul>

### static

Resolves a path into a site-relative URL.

**Use when** referencing an asset from a page that may be moved, or where `baseURL` differs between local and production builds. Most commonly nested inside an `<img>` tag or a Markdown image link.

| Parameter | Required | Default | Description |
|---|---|---|---|
| _positional_ | Yes | — | Path relative to the site root, passed as the first unnamed argument. Emitted as a bare URL with no surrounding markup. |

Takes no body. Use the `{{</* */>}}` form.

#### Syntax

```markdown
![Meshery Architecture]({{</* static "images/meshery-architecture.webp" */>}})
```

#### Rendered output

![Meshery Architecture]({{< static "images/meshery-architecture.webp" >}})

### tabs

Organizes content into selectable tabs.

**Use when** presenting the same task across mutually exclusive contexts — operating systems, installation methods — where the reader needs exactly one branch. Don't use it for sequential content; readers miss what's behind an unselected tab.

| Parameter | Required | Default | Description |
|---|---|---|---|
| `id` | Yes | — | Identifier for the tab group. Omitting it raises a build error naming the file and position. Each invocation is additionally given an instance number, so repeating an `id` on one page does not cause tab groups to interfere. |
| `level` | No | `1` | Nesting level. Values above `1` add a nested-container class, for tabs inside tabs. |

Takes a body. Use the `{{</* */>}}` form. Within the body:

- Separate tabs with an `<!-- tab -->` HTML comment.
- The **first line of each block is the tab label**; the content begins after the first blank line.
- A label may carry an icon by appending `| ` and an icon value. A value starting with `fa` is treated as a Font Awesome class; anything else is inserted as raw markup.
- The first non-empty block is selected by default.

#### Syntax

```markdown
{{</* tabs id="example-tab-group" */>}}
Tab 1 Title

This is the content for the first tab.
<!-- tab -->
Tab 2 Title | fas fa-cog

This is the content for the second tab.
{{</* /tabs */>}}
```

#### Rendered output

{{< tabs id="example-tab-group" >}}
Tab 1 Title

This is the content for the first tab.
<!-- tab -->
Tab 2 Title | fas fa-cog

This is the content for the second tab.
{{< /tabs >}}

### Snippet shortcodes

Some shortcodes are Markdown files in subdirectories of `docs/layouts/shortcodes/` rather than templates, and are invoked by path. They take no parameters and no body, and exist to keep a reusable block of prose in one place:

| Shortcode | Content |
|---|---|
| `{{</* mesheryctl/installation-bash */>}}` | `mesheryctl` installation via bash |
| `{{</* mesheryctl/installation-brew */>}}` | `mesheryctl` installation via Homebrew |
| `{{</* mesheryctl/installation-scoop */>}}` | `mesheryctl` installation via Scoop |
| `{{</* mesheryctl/system-dashboard */>}}` | `mesheryctl system dashboard` usage |
| `{{</* compatibility/adapter-status */>}}` | Adapter status block |
| `{{</* compatibility/compatibilityMatrix */>}}` | Compatibility matrix prose |
| `{{</* compatibility/rancher */>}}` | Rancher-specific notes |

To add one, drop a `.md` file into the relevant subdirectory — no template is required.

---

## Data-Driven Shortcodes

The shortcodes below each render a full, live dataset — error codes, vulnerabilities, `mesheryctl` commands, permissions, and so on — pulled from `docs/data/` or the site's own content. Because their real output can run to hundreds of rows and changes as that underlying data changes, this section shows the invocation syntax and links to the page where each one is used, rather than re-rendering the full table here.

All of them take no body, and all take no parameters except `mesheryctl-command-table`, which requires a `command` naming the command whose subcommands are tabulated.

| Shortcode | Syntax | See it rendered live |
|---|---|---|
| `adapters-table` | `{{</* adapters-table */>}}` | [Adapters]({{< ref "concepts/architecture/adapters.md" >}}) |
| `compatibility-matrix-kubernetes` | `{{</* compatibility-matrix-kubernetes */>}}` | [Compatibility Matrix]({{< ref "project/compatibility-matrix/compatibility-matrix.md" >}}) |
| `integration-tests` | `{{</* integration-tests */>}}` | [Compatibility Matrix]({{< ref "project/compatibility-matrix/compatibility-matrix.md" >}}) |
| `error-codes-index` | `{{</* error-codes-index */>}}` | [Error Codes Reference]({{< ref "reference/references/error-codes.md" >}}) |
| `error-codes-detail` | `{{</* error-codes-detail */>}}` | [Error Codes Reference]({{< ref "reference/references/error-codes.md" >}}) |
| `troubleshooting-guides-list` | `{{</* troubleshooting-guides-list */>}}` | [Error Codes Reference]({{< ref "reference/references/error-codes.md" >}}) |
| `mesheryctl-command-table` | `{{</* mesheryctl-command-table command="adapter" */>}}` | [mesheryctl Reference]({{< ref "reference/references/mesheryctl/_index.md" >}}) |
| `mesheryctl-guides-list` | `{{</* mesheryctl-guides-list */>}}` | [mesheryctl Installation]({{< ref "installation/mesheryctl/_index.md" >}}) |
| `network-ports` | `{{</* network-ports */>}}` | [Architecture]({{< ref "concepts/architecture/_index.md" >}}) |
| `permissions` | `{{</* permissions */>}}` | [Permissions Reference]({{< ref "reference/references/permissions.md" >}}) |
| `tutorials-list` | `{{</* tutorials-list */>}}` | [Quick Start]({{< ref "installation/quick-start/index.md" >}}) |
| `vulnerabilities-table` | `{{</* vulnerabilities-table */>}}` | [Security Vulnerabilities]({{< ref "project/security-vulnerabilities.md" >}}) |

---

## Inherited Shortcodes

These come from the Docsy theme rather than this repository. They are listed here because they appear throughout the Meshery documentation and contributors will encounter them; the [Docsy shortcode reference](https://www.docsy.dev/docs/adding-content/shortcodes/) is authoritative for their full behavior.

### alert

Renders a coloured callout box with a heading.

**Use when** flagging a warning, prerequisite, or caveat that would interrupt the flow of the surrounding paragraph. For inline emphasis, use bold instead — a page dense with callouts trains readers to skip them.

| Parameter | Required | Default | Description |
|---|---|---|---|
| `color` | No | `primary` | Bootstrap contextual class: `primary`, `secondary`, `success`, `info`, `warning`, `danger`. |
| `title` | No | — | Bold heading above the body. |

Takes a body. Use the `{{%/* */%}}` form, or the body's Markdown will render as literal asterisks and backticks.

#### Syntax

```markdown
{{%/* alert color="warning" title="Note" %}}
Run `make site` locally before opening your PR.
{{% /alert */%}}
```

#### Rendered output

{{% alert color="warning" title="Note" %}}
Run `make site` locally before opening your PR.
{{% /alert %}}

### pageinfo

Renders a page-scoped information banner.

**Use when** the note applies to the whole page — deprecation, beta status, "applies to version X" — and belongs above the content. For a note scoped to one section, use `alert`.

| Parameter | Required | Default | Description |
|---|---|---|---|
| `color` | No | `primary` | Bootstrap contextual class, as for `alert`. |

Takes a body. Use the `{{%/* */%}}` form.

#### Syntax

```markdown
{{%/* pageinfo color="info" %}}
This page describes a **beta** feature.
{{% /pageinfo */%}}
```

#### Rendered output

{{% pageinfo color="info" %}}
This page describes a **beta** feature.
{{% /pageinfo %}}

Docsy also provides `youtube`, `imgproc`, `blocks/*`, `cardpane`, and others. See the Docsy reference for those.

---

## Adding a New Shortcode

Shortcode implementations live in `docs/layouts/shortcodes/`. The filename becomes the shortcode name: `docs/layouts/shortcodes/model-count.html` is invoked as `{{</* model-count */>}}`, and a file in a subdirectory is invoked by path, as `{{</* mesheryctl/installation-bash */>}}`.

A template in this repository overrides a theme template of the same name, so adding `docs/layouts/shortcodes/alert.html` would replace Docsy's `alert` site-wide. Do that deliberately, if at all.

1. **Create the template.** Open with a comment block stating what the shortcode outputs, any assumptions it makes about page location or front matter, and a usage line.

2. **Guard your assumptions.** Shortcodes are invoked from pages you did not write. Check that the values you depend on exist before indexing into them.

3. **Choose a failure mode deliberately.** Use `errorf` when a missing argument means the author made a mistake that must be fixed, as `tabs` does when `id` is absent. Use `warnf`, or emit nothing, when the page still renders usefully without the output. Both take `.Position` so the message names the file and line.

4. **Control whitespace.** Use `{{- -}}` trim markers so the shortcode does not emit stray newlines into the surrounding Markdown, which can break list numbering and inline usage.

5. **Say what wrapping the caller must supply.** A shortcode that emits bare `<li>` elements, as `section-pages` and `model-categories` do, is unusable without that fact being written down.

6. **Document it here.** Add an entry to the catalog above with the parts every entry carries: a description, a "use when" line, a parameter table, and a syntax and rendered-output pair. A shortcode that isn't on this page is a shortcode contributors won't find.

**Example.** `model-count` derives its number from the site's own content rather than a hardcoded value, so the figure never goes stale:

```go-html-template
{{- $models := where site.RegularPages ".File.Dir" "like" `^extensions/models/` -}}
{{- len $models -}}
```

Note: When a shortcode can produce empty or zero output without failing, its entry above must say so, or contributors will read the result as a bug and file it as one.

{{% alert color="warning" title="Register your shortcode" %}}
If a new shortcode ships without an entry on this page, contributors have no way to discover it short of reading `docs/layouts/shortcodes/`, which is the problem this page exists to solve. Treat the catalog entry as part of the shortcode, not as follow-up work.
{{% /alert %}}