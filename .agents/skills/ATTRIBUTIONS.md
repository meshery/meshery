# Third-party skill attributions

Skills in this directory that are adapted from third-party sources. Each adapted
skill records its upstream URL and version in its `SKILL.md` frontmatter.

## armory

Source: <https://github.com/Mathews-Tom/armory> (MIT)

| Skill | Upstream skill | Changes |
|---|---|---|
| `pre-landing-review` | `skills/pre-landing-review` 1.0.1 | Checklist rewritten around Meshery's AGENTS.md rules (gorm `.Order()` sanitization, provider resolution at boot, page access gates, schemas-first APIs, MeshKit error codes); repo gates added; evals dropped |
| `migration-risk-analyzer` | `skills/migration-risk-analyzer` 1.1.1 | Rewritten for gorm `AutoMigrate` on Meshery's SQLite database (upgrade/downgrade/reset paths, column renames from schemas bumps) in place of PostgreSQL/MySQL lock analysis; reference files and evals dropped |

```
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
