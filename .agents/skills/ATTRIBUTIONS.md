# Skill Attributions and Provenance

This document records the provenance, modifications, and upstream licensing for third-party skills adapted into `.agents/skills/`.

---

## 1. Mathews-Tom/armory

* **Source Repository**: [https://github.com/Mathews-Tom/armory](https://github.com/Mathews-Tom/armory)
* **License**: MIT License
* **Upstream Components Adapted**:
  * `pre-landing-review`
  * `migration-risk-analyzer`

### Summary of Adaptations for Meshery

1. **`pre-landing-review`**:
   * Rewritten to enforce the architectural rules and failure classes specified in Meshery's root `AGENTS.md`.
   * Encodes blocking checks for:
     * GORM `.Order(...)` SQL injection prevention via `models.SanitizeOrderInput`.
     * Safe boot-time provider lookup rules avoiding nil pointer panics on pinned provider deployments.
     * UI permission gating (`useHasPermission(Keys.X)`) and query suppression via `{ skip: !canAccess }` in RTK Query.
     * Schema-first API discipline requiring changes to route through `meshery/schemas`.
     * Proper casing conventions (DB snake_case, wire camelCase, Go PascalCase).
     * MeshKit structured error requirements and error reference documentation updates.
     * DCO sign-off requirements.
   * Configured verification to run Meshery's linter targets: `make golangci` and `make ui-lint`.

2. **`migration-risk-analyzer`**:
   * Rewritten specifically for Meshery Server's GORM `AutoMigrate` lifecycle on SQLite (`mesherydb.sql`).
   * Analyzes SQLite-specific migration hazards including:
     * Column renames from `@meshery/schemas` version bumps (which add new columns rather than renaming in place, risking silent data orphaning).
     * `NOT NULL` constraints without defaults breaking existing SQLite databases.
     * Coverage verification in `server/models/system_migration.go:SystemDatabaseModels()` for clean database resets.
     * Safe snapshotting and WAL journal handling during reset and backup operations.

---

### Upstream MIT License Text

```text
MIT License

Copyright (c) 2026 Tom

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
