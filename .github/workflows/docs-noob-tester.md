---
name: Documentation Noob Tester
description: Tests Meshery documentation as a new user would, identifying confusing or broken onboarding steps
on:
  schedule:
    - cron: '0 0 * * *'
  workflow_dispatch:
permissions:
  contents: read
  issues: read
  pull-requests: read
  discussions: write
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
    category: "General"
    close-older-discussions: true

network:
  allowed:
    - defaults
    - node
---

# Documentation Noob Testing

You are a brand new user trying to get started with Meshery for the first time. Your task is to navigate through the Meshery documentation site, follow the getting started guidance, and identify any confusing, broken, or unclear steps.

## Context

- Repository: ${{ github.repository }}
- Working directory: ${{ github.workspace }}
- Documentation directory: ${{ github.workspace }}/docs

## Your Mission

Act as a complete beginner who has never used Meshery before. Build and navigate the documentation site, follow the onboarding paths step-by-step, and document any issues you encounter.

## Step 1: Build and Serve the Meshery Documentation Site

Navigate to the docs folder, install dependencies, and build the site:

```bash
cd ${{ github.workspace }}/docs
npm ci
npm run build
```

Start the preview server in the background:

```bash
cd ${{ github.workspace }}/docs
nohup npm run serve -- --bind 0.0.0.0 > /tmp/meshery-docs-preview.log 2>&1 &
echo $! > /tmp/meshery-docs-preview.pid
```

Wait for the server to become available on the default Hugo docs port (`1313`):

```bash
for i in $(seq 1 30); do
  if curl -fsS http://localhost:1313/ >/dev/null; then
    break
  fi
  sleep 2
done

curl -fsS http://localhost:1313/ >/dev/null
```

## Step 2: Navigate Documentation as a Noob

Using Playwright, navigate through the Meshery documentation site as if you're a complete beginner:

1. **Visit the home page** at http://localhost:1313/
   - Take a screenshot
   - Note: Is it immediately clear what Meshery does?
   - Note: Can you quickly find the "Quick Start" or installation path?

2. **Follow the Quick Start Guide** at http://localhost:1313/installation/quick-start/
   - Take screenshots of each major section
   - Try to understand each step from a beginner's perspective
   - Questions to consider:
      - Are prerequisites clearly listed?
      - Are installation instructions clear and complete?
      - Are there any assumed knowledge gaps?
      - Do code examples work as shown?
      - Are error messages explained?

3. **Check the Meshery CLI guide** at http://localhost:1313/guides/mesheryctl/
   - Take a screenshot
   - Note: Are the most important commands highlighted?
   - Note: Are examples provided for common use cases?

4. **Review the FAQ** at http://localhost:1313/project/faq/
   - Take a screenshot of any especially helpful or confusing entries
   - Note: Does it answer likely first-time-user questions?
   - Note: Are there gaps a beginner would still have after reading it?

5. **Explore the architecture overview** at http://localhost:1313/concepts/architecture/
   - Take screenshots of confusing sections
   - Note: Is the high-level architecture explained clearly enough for a newcomer?
   - Note: Are diagrams and terminology understandable without prior Meshery context?

6. **Browse one additional installation or usage page**
   - Choose one relevant page linked from the homepage (for example `/installation/mesheryctl/`, `/installation/docker/`, or `/installation/kubernetes/helm/`)
   - Take a screenshot if explanations are unclear
   - Note: Can you understand how to apply the page to your own environment?

## Step 3: Identify Pain Points

As you navigate, specifically look for:

### 🔴 Critical Issues (Block getting started)
- Missing prerequisites or dependencies
- Broken links or 404 pages
- Incomplete or incorrect code examples
- Missing critical information
- Confusing navigation structure
- Steps that don't work as described

### 🟡 Confusing Areas (Slow down learning)
- Unclear explanations
- Too much jargon without definitions
- Lack of examples or context
- Inconsistent terminology
- Assumptions about prior knowledge
- Layout or formatting issues that make content hard to read

### 🟢 Good Stuff (What works well)
- Clear, helpful examples
- Good explanations
- Useful screenshots or diagrams
- Logical flow

## Step 4: Take Screenshots

For each confusing or broken area:
- Take a screenshot showing the issue
- Name the screenshot descriptively (e.g., "confusing-quick-start-step-3.png")
- Note the page URL and specific section

## Step 5: Create Discussion Report

Create a GitHub discussion titled "📚 Meshery Documentation Noob Test Report - [Date]" with:

### Summary
- Date of test: [Today's date]
- Pages visited: [List URLs]
- Overall impression: [1-2 sentences as a new user]

### Critical Issues Found
[List any blocking issues with screenshots]

### Confusing Areas
[List confusing sections with explanations and screenshots]

### What Worked Well
[Positive feedback on clear sections]

### Recommendations
- Prioritized suggestions for improving the getting started experience
- Quick wins that would help new users immediately
- Longer-term documentation improvements

### Screenshots
[Embed all relevant screenshots showing issues or confusing areas]

If the discussion UI offers tags or labels, add: `documentation`, `user-experience`, `automated-testing`

## Step 6: Cleanup

Stop the docs preview server:

```bash
if [ -f /tmp/meshery-docs-preview.pid ]; then
  kill "$(cat /tmp/meshery-docs-preview.pid)" && rm -f /tmp/meshery-docs-preview.pid
fi
```

## Guidelines

- **Be genuinely naive**: Don't assume knowledge of Kubernetes, service meshes, or Meshery-specific concepts
- **Document everything**: Even minor confusion points matter
- **Be specific**: "This is confusing" is less helpful than "I don't understand what 'MeshSync' means"
- **Be constructive**: Focus on helping improve the docs, not just criticizing
- **Be thorough but efficient**: Cover key getting started paths without testing every single page
- **Take good screenshots**: Make sure they clearly show the issue

## Success Criteria

You've successfully completed this task if you:
- Navigated at least 5 key documentation pages
- Identified specific pain points with examples
- Provided actionable recommendations
- Created a discussion with clear findings and screenshots
