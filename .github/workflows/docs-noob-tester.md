---
name: Meshery Docs Noob Tester
description: Tests Meshery Docs as a new user and first-time contributor, identifying confusing or broken onboarding, navigation, and Hugo/Docsy rendering
on:
  schedule: daily
  workflow_dispatch:
permissions:
  contents: read
features:
  copilot-requests: true
engine: copilot
timeout-minutes: 30
tools:
  playwright:
  edit:
  bash:
    - "*"
safe-outputs:
  upload-asset:
  create-discussion:
    category: "audits"
    close-older-discussions: true
    fallback-to-issue: false
  noop: false
  missing-tool: false
  missing-data: false
  report-incomplete: false
  report-failure-as-issue: false

network:
  allowed:
    - defaults
    - node

imports:
  - shared/mood.md
  - shared/docs-server-lifecycle.md
---

# Meshery Docs noob testing

You are evaluating **Meshery Docs**, a Hugo-based documentation website that uses the **Docsy** theme. Your task is to experience the site as both:

1. A new Meshery user trying to understand what Meshery is and how to install it.
2. A first-time docs contributor trying to build, preview, and navigate the documentation locally.

Identify anything confusing, broken, unclear, or inconsistent in the site experience, the contributor flow, and the rendered Hugo/Docsy pages.

## Context

- Repository: ${{ github.repository }}
- Working directory: ${{ github.workspace }}
- Documentation directory: ${{ github.workspace }}/docs

## Your Mission

Act as a complete beginner who is new to Meshery and new to contributing to Meshery Docs. Build and navigate the site, follow installation and contributor guides step-by-step, and document any issues you encounter.

## Step 1: Build and serve Meshery Docs

Meshery Docs is not a generic static site. It is built with Hugo Extended, uses the Docsy theme as a Hugo module, and relies on both Node.js and Go during local preview and build workflows.

From the docs directory, verify the toolchain and use the repo's real docs build steps:

```bash
cd ${{ github.workspace }}/docs
node --version
go version
npx hugo version
npm ci
npm run build:preview
```

Then start a local preview server on port `1313` using the repo-native preview command:

```bash
cd ${{ github.workspace }}/docs
npm run serve -- --bind 0.0.0.0 --disableFastRender
```

If the npm script fails, fall back to:

```bash
cd ${{ github.workspace }}/docs
npx hugo server -D -F --bind 0.0.0.0 --baseURL http://localhost:1313/
```

Follow the shared **Documentation Server Lifecycle Management** instructions:
1. Start the preview server.
2. Wait for server readiness.

If build or serve fails, treat that as a **critical issue** and capture the exact command and error message.

## Step 2: Navigate the site as a new Meshery user

Using Playwright, navigate through the docs as if you're a complete beginner who wants to install and understand Meshery:

1. **Visit the home page** at http://localhost:1313
   - Take a screenshot
   - Note: Is it immediately clear what Meshery does?
   - Note: Can you quickly find installation, quick start, and contributor guidance?
   - Note: Do the homepage sections and navigation feel scannable to a newcomer?

2. **Visit the installation overview** at http://localhost:1313/installation/
   - Take a screenshot
   - Note: Is it clear which installation path to choose?
   - Note: Are prerequisites and supported platforms explained clearly?

3. **Follow the Quick Start Guide** at http://localhost:1313/installation/quick-start/
   - Take screenshots of each major section
   - Try to understand each step from a beginner's perspective
   - Questions to consider:
     - Are prerequisites clearly listed?
     - Are installation instructions clear and complete?
     - Are there knowledge gaps or missing context?
     - Do alerts, callouts, images, and code blocks render correctly?
     - Are follow-up actions after installation obvious?

4. **Check the mesheryctl installation page** at http://localhost:1313/installation/mesheryctl/
   - Take a screenshot
   - Note: Are the supported install methods easy to compare?
   - Note: Do the embedded install sections and shortcodes render correctly?

5. **Visit the Meshery UI access guide** at http://localhost:1313/installation/accessing-meshery-ui/
   - Take a screenshot
   - Note: Can a new user figure out where Meshery UI will be available after installation?
   - Note: Are deployment-specific access patterns explained clearly?

6. **Use search** by visiting http://localhost:1313/search/?q=mesheryctl
   - Take a screenshot
   - Note: Does search load successfully?
   - Note: Are results relevant, readable, and easy to act on?
   - Note: Do result links open the expected pages?

7. **Check one generated reference page** at http://localhost:1313/reference/mesheryctl/system/check/
   - Take a screenshot
   - Note: Does the reference page render cleanly?
   - Note: Are commands, descriptions, headings, and navigation easy to follow?

## Step 3: Navigate the site as a first-time docs contributor

Using Playwright, now switch personas and explore the site as someone trying to contribute to Meshery Docs for the first time:

1. **Visit the contributing landing page** at http://localhost:1313/project/contributing/
   - Take a screenshot
   - Note: Is it obvious how a newcomer should begin contributing?
   - Note: Are DCO/signoff expectations and contribution flow understandable?

2. **Follow the docs contributor guide** at http://localhost:1313/project/contributing/contributing-docs/
   - Take screenshots of important sections
   - Questions to consider:
     - Are Node.js, Go, Hugo, and local preview requirements clear?
     - Are `make` and `npm` commands explained consistently?
     - Is it obvious how to run the site locally and verify changes?
     - Are images, alerts, code blocks, and examples rendered correctly?

3. **Review the docs structure guide** at http://localhost:1313/project/contributing/contributing-docs-structure/
   - Take a screenshot
   - Note: Does the information architecture make sense?
   - Note: Can a first-time contributor understand where to add or edit content?

4. **Review the contributor gitflow guide** at http://localhost:1313/project/contributing/contributing-gitflow/
   - Take a screenshot
   - Note: Is the fork, branch, sign-off, and pull request flow understandable?
   - Note: Are links to related contribution guidance valid and useful?

5. **Check an alias/redirect path** by visiting http://localhost:1313/platforms/
   - Note whether it correctly resolves to the canonical installation section
   - If the redirect is broken, slow, or confusing, record it as an issue

6. **Pay attention to Docsy/Hugo behavior across these pages**
   - Verify that shortcodes render as UI components instead of raw template syntax
   - Verify that code snippets, clipboard blocks, alerts, images, and internal links work
   - Verify that the sidebar, breadcrumbs, and page hierarchy feel coherent
   - Note any layout regressions, overflow, or styling issues that make the content harder to use

## Step 4: Validate documentation quality

For each page you visit, check the following:

- The local preview renders correctly at `http://localhost:1313`
- Navigation is coherent: homepage, section landing pages, sidebar, breadcrumbs, and next steps
- Search can find relevant content by page title or key terms
- Internal links and redirects resolve to the expected canonical page
- Images, code blocks, clipboard blocks, alerts, and embedded content render correctly
- Page titles, descriptions, and terminology feel consistent with Meshery docs

Only when relevant to the page being tested, also check:

- Contributor instructions are consistent across `make`, `npm`, and direct Hugo commands
- Alias/redirect behavior is sensible when pages have moved
- Search results surface the right page and snippet
- Hugo/Docsy shortcodes render as components rather than raw syntax
- Layout remains readable across long code blocks, tables, alerts, and images

## Step 5: Identify pain points

As you navigate, specifically look for:

### 🔴 Critical Issues (Block getting started)
- Build or preview server fails to start
- Missing Node.js / Go / Hugo / Docsy prerequisites in contributor docs
- Broken links or 404 pages
- Broken redirects or aliases
- Incomplete or incorrect commands
- Missing critical information
- Search page fails or returns unusable results
- Raw Hugo shortcode syntax visible in rendered pages
- Missing images, broken assets, or layout-breaking rendering bugs
- Confusing navigation structure
- Steps that don't work as described

### 🟡 Confusing Areas (Slow down learning)
- Unclear explanations
- Too much Meshery, Hugo, or Docsy jargon without definitions
- Inconsistent `make`, `npm`, and `hugo` guidance
- Lack of examples or context
- Inconsistent terminology
- Assumptions about prior knowledge
- Unclear contributor workflow or local preview instructions
- Layout or formatting issues that make content hard to read

### 🟢 Good Stuff (What works well)
- Clear, helpful examples
- Good explanations
- Useful screenshots or diagrams
- Logical flow
- Helpful search results or navigation
- Contributor guidance that is easy to follow

## Step 6: Take screenshots

For each confusing or broken area:
- Take a screenshot showing the issue
- Name the screenshot descriptively (for example: `confusing-quick-start-prereqs.png` or `broken-search-results-page.png`)
- Note the page URL and specific section

## Step 7: Create discussion report

Create a GitHub discussion titled **"📚 Meshery Docs newcomer test report - [Date]"** with:

### Summary
- Date of test: [Today's date]
- Branch / commit tested: [fill from `git rev-parse --abbrev-ref HEAD` and `git rev-parse HEAD`]
- Toolchain used: [Node version, Go version, Hugo version]
- Pages visited: [List URLs]
- Overall impression as a new Meshery user: [1-2 sentences]
- Overall impression as a first-time docs contributor: [1-2 sentences]

### Critical Issues Found
[List any blocking issues with screenshots]

### Confusing Areas
[List confusing sections with explanations and screenshots]

### Hugo / Docsy rendering and navigation findings
[List issues related to shortcodes, search, redirects, sidebar structure, breadcrumbs, images, code blocks, and layout]

### What Worked Well
[Positive feedback on clear sections]

### Recommendations
- Prioritized suggestions for improving the newcomer and contributor experience
- Quick wins that would help docs contributors immediately
- Longer-term Meshery Docs improvements

### Screenshots
[Embed all relevant screenshots showing issues or confusing areas]

Start the discussion body with:

`Focus areas: docs, user-experience, hugo, docsy, newcomer-onboarding`

## Step 8: Cleanup

Follow the shared **Documentation Server Lifecycle Management** instructions for cleanup and stop the local Hugo preview server cleanly.

## Guidelines

- **Be genuinely naive**: Don't assume prior knowledge of Meshery, Hugo, Docsy, DCO, or contributor workflows
- **Document everything**: Even minor confusion points matter
- **Be specific**: "This is confusing" is less helpful than "I don't understand what `make site` does or why Go is required"
- **Be constructive**: Focus on helping improve the docs, not just criticizing
- **Check rendered output, not just prose**: Hugo/Docsy issues often show up as broken rendering, search, redirects, or missing assets
- **Be thorough but efficient**: Cover the key installation and docs-contributor journeys without testing every single page
- **Take good screenshots**: Make sure they clearly show the issue

## Success Criteria

You've successfully completed this task if you:
- Built and previewed the Hugo/Docsy site locally
- Navigated the key newcomer and docs-contributor pages
- Checked search, redirects, and rendered Hugo/Docsy components
- Identified specific pain points with examples
- Provided actionable recommendations
- Created a discussion with clear findings and screenshots
