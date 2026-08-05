---
title: "Shortcode Formatting Examples"
description: "A visual reference guide for shortcodes available in the Meshery documentation theme."
categories: [contributing]
aliases: [/project/contributing/shortcode-examples]
---
The Meshery documentation utilizes a variety of shortcodes to create rich, formatted content such as callouts, alerts, and embedded media.

This page serves as a visual reference guide for documentation contributors. It allows you to see exactly how these shortcodes render on the live site before you use them. Below, you will find the required markdown syntax to invoke each custom Meshery shortcode, followed immediately by its rendered output.

---
## Meshery Custom Shortcodes

### 1. Discuss Shortcode

The `discuss` shortcode renders a static alert box linking to the Meshery Discussion Forum. It does not require any parameters.

**Syntax:**
```markdown
{{</* discuss */>}}
```
Rendered Output:

{{< discuss >}}

### 2. Relationships Shortcode

The `relationships` shortcode inserts a predefined, interactive hierarchical visual representation of components and bindings. It does not require any parameters.

**Syntax:**
```markdown
{{</* relationships */>}}
```
Rendered Output:

{{< relationships >}}

### 3. Tabs Shortcode

The `tabs` shortcode allows you to organize content into selectable tabs. It requires a unique `id` argument. Separate individual tabs using the `<!-- tab -->` HTML comment.

**Syntax:**
```markdown
{{</* tabs id="example-tab-group" */>}}
Tab 1 Title

This is the content for the first tab.
<!-- tab -->
Tab 2 Title

This is the content for the second tab.
{{</* /tabs */>}}
```
Rendered Output:

{{< tabs id="example-tab-group" >}}
Tab 1 Title

This is the content for the first tab.
<!-- tab -->
Tab 2 Title

This is the content for the second tab.
{{< /tabs >}}

### 4. Code Shortcode

The `code` shortcode wraps a snippet in a styled code block with a "copy to clipboard" button, without needing a manually assigned ID per snippet.

**Syntax:**
```markdown
{{</* code code="make site" */>}}
```
Rendered Output:

{{< code code="make site" >}}

{{% alert color="warning" title="Note" %}}
Don't use the `code` shortcode when the snippet appears between items of an ordered list — it breaks the numbering, and the item following the snippet restarts at 1.
{{% /alert %}}

### 5. Static Shortcode

The `static` shortcode resolves a relative path into a proper site-relative URL, most commonly used inside an `<img>` tag or Markdown image link so that asset paths keep working regardless of the page's location or `baseURL` setting.

**Syntax:**
```markdown
![Meshery Architecture]({{</* static "images/meshery-architecture.webp" */>}})
```
Rendered Output:
![Meshery Architecture]({{< static "images/meshery-architecture.webp" >}})

### 6. Extension Guide Shortcode

The `extension-guide` shortcode renders a responsive card grid from a data file in `docs/data/`. Each card shows an SVG/image, a title, and a description. It requires `data_file`, `guide_title`, `guide_description`, `guide_svg`, and `guide_assests_folder` (spelling intentional) parameters that map to keys/paths in the referenced data file.

**Syntax:**
```markdown
{{</* extension-guide
  data_file="edges"
  guide_title="Edge"
  guide_description="Description"
  guide_svg="SVG"
  guide_assests_folder="shapes"
*/>}}
```
Rendered Output:

{{< extension-guide
  data_file="edges"
  guide_title="Edge"
  guide_description="Description"
  guide_svg="SVG"
  guide_assests_folder="shapes"
>}}

### 7. Related Discussions Shortcode

The `related-discussions` shortcode pulls up to eight recent topics from the Meshery Discussion Forum matching a given `tag`, plus a link to view all discussions with that tag.

**Syntax:**
```markdown
{{</* related-discussions tag="meshery" */>}}
```
Rendered Output:

{{< related-discussions tag="meshery" >}}

### 8. Section Pages Shortcode

The `section-pages` shortcode auto-generates a linked list of the immediate child pages of a given `section`, along with each page's `description` (when set in front matter). If `section` is omitted, it defaults to the page’s top-level Hugo section (Hugo’s `.Page.Section`), so for nested paths you’ll usually want to set `section` explicitly.

**Syntax:**
```markdown
{{</* section-pages section="concepts/architecture" */>}}
```
Rendered Output:

{{< section-pages section="concepts/architecture" >}}

### 9. Latest Release Shortcode

The `latest-release` shortcode looks up all non-prerelease entries under the `releases` content type and prints the tag of the most recent one as an `<h2>` heading. It does not require any parameters.

**Syntax:**
```markdown
{{</* latest-release */>}}
```
Rendered Output:

{{< latest-release >}}

### 10. Model Count Shortcode

The `model-count` shortcode prints the current total number of Meshery models, useful for inline sentences that need to stay accurate as new models are added. It does not require any parameters.

**Syntax:**
```markdown
Meshery currently supports {{</* model-count */>}} models.
```
Rendered Output:
Meshery currently supports {{< model-count >}} models.

---
## Reference & Data Shortcodes

The shortcodes below each render a full, live dataset (error codes, vulnerabilities, mesheryctl commands, permissions, and so on) pulled directly from `docs/data/` or the site's own content. Because their real output can run to dozens or hundreds of rows and changes as that underlying data changes, this guide shows the invocation syntax and links to the page where each one is actually used, rather than re-rendering the full table here — that avoids maintaining the same large dataset in two places.

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

{{% alert color="info" title="Scope note for reviewers" %}}
This page covers every shortcode in `docs/layouts/shortcodes/` except the ones nested under `compatibility/`, `installation/`, `mesheryctl/`, and `svg/` — those are single-purpose content fragments authored for one specific page each, rather than general-purpose formatting shortcodes, so they aren't included here.
{{% /alert %}}