---
layout: page
title: Docs Migration Checklist (Jekyll to Hugo)
permalink: project/contributing/docs-migration-checklist
abstract: A checklist to preserve URLs, structure, and contributor workflows during the docs migration.
language: en
type: project
category: contributing
list: include
---

This checklist supports the planned migration of docs.meshery.io from Jekyll to Hugo. Use it to track and verify migration readiness while keeping URLs stable and contributor workflows intact.

## URL Preservation

- Ensure every page keeps its current `permalink`.
- Preserve `redirect_from` entries when moving or renaming pages.
- Verify that canonical URLs match existing docs.meshery.io paths.

## Frontmatter and Metadata

- Map Jekyll frontmatter fields to Hugo frontmatter equivalents.
- Preserve `title`, `abstract`, `language`, and navigation metadata.

## Includes and Shortcodes

- Inventory `{% include %}` usage (for example, `alert.html`, `code.html`).
- Create Hugo shortcodes that match behavior and naming.

## Collections and Content Structure

- Map Jekyll collections (for example, `_adapters`, `_releases`) to Hugo sections.
- Preserve any auto-generated content flow.

## Assets and Paths

- Keep asset paths stable (images, SVGs, downloads).
- Confirm base URL handling in templates and links.

## Navigation and IA

- Record current sidebar hierarchy (parent, children, grandchildren).
- Align new structure to the Diataxis framework.

## Contributor Workflow

- Document the new local dev command(s) and requirements.
- Update the docs contribution guide to point to the Hugo flow.

## Validation

- Spot check a sample of high-traffic pages for broken links.
- Run a link checker or CI validation where possible.

