# Files to Commit - AWS Secrets Manager Relationships Fix

## Summary
**Total Files**: 79 (1 script + 78 relationship files)
**Issue Fixed**: #17179
**Previous PR**: #17189 (had issues with zero UUIDs and missing workloads)

---

## 1. Script File (1 file)

### scripts/fix_awssm_k8s_relationships.py
- **Status**: NEW
- **Lines**: 199
- **Purpose**: Generate/regenerate AWS Secrets Manager relationships with unique UUIDs
- **Features**:
  - Generates deterministic UUID v5 for each workload-version combination
  - Creates relationship JSON files for 6 workload types
  - Processes all controller versions
  - Can be run to regenerate relationships if needed

---

## 2. Relationship Files (78 files total)

### Structure
```
server/meshmodel/aws-secretsmanager-controller/{VERSION}/v1.0.0/relationships/
```

### Naming Pattern
- `edge-non-binding-network-deplo.json` → Deployment
- `edge-non-binding-network-podds.json` → Pod
- `edge-non-binding-network-stsfs.json` → StatefulSet
- `edge-non-binding-network-dmnst.json` → DaemonSet
- `edge-non-binding-network-jobbs.json` → Job
- `edge-non-binding-network-crnjb.json` → CronJob

### Versions Covered
1. v1.0.7 - (6 files)
2. v1.0.8 - (6 files)
3. v1.0.9 - (6 files)
4. v1.0.10 - (6 files)
5. v1.0.11 - (6 files)
6. v1.0.12 - (6 files)
7. v1.0.13 - (6 files)
8. v1.0.14 - (6 files)
9. v1.1.0 - (6 files)
10. v1.1.1 - (6 files)
11. v1.1.2 - (6 files)
12. v1.2.0 - (6 files)
13. v1.2.1 - (6 files)

### Full File List
```
server/meshmodel/aws-secretsmanager-controller/v1.0.7/v1.0.0/relationships/edge-non-binding-network-deplo.json
server/meshmodel/aws-secretsmanager-controller/v1.0.7/v1.0.0/relationships/edge-non-binding-network-podds.json
server/meshmodel/aws-secretsmanager-controller/v1.0.7/v1.0.0/relationships/edge-non-binding-network-stsfs.json
server/meshmodel/aws-secretsmanager-controller/v1.0.7/v1.0.0/relationships/edge-non-binding-network-dmnst.json
server/meshmodel/aws-secretsmanager-controller/v1.0.7/v1.0.0/relationships/edge-non-binding-network-jobbs.json
server/meshmodel/aws-secretsmanager-controller/v1.0.7/v1.0.0/relationships/edge-non-binding-network-crnjb.json

[Same 6 files for each version: v1.0.8 through v1.2.1]

server/meshmodel/aws-secretsmanager-controller/v1.2.1/v1.0.0/relationships/edge-non-binding-network-deplo.json
server/meshmodel/aws-secretsmanager-controller/v1.2.1/v1.0.0/relationships/edge-non-binding-network-podds.json
server/meshmodel/aws-secretsmanager-controller/v1.2.1/v1.0.0/relationships/edge-non-binding-network-stsfs.json
server/meshmodel/aws-secretsmanager-controller/v1.2.1/v1.0.0/relationships/edge-non-binding-network-dmnst.json
server/meshmodel/aws-secretsmanager-controller/v1.2.1/v1.0.0/relationships/edge-non-binding-network-jobbs.json
server/meshmodel/aws-secretsmanager-controller/v1.2.1/v1.0.0/relationships/edge-non-binding-network-crnjb.json
```

---

## How to Add Files to Git

### Option 1: Add Script Separately
```bash
git add scripts/fix_awssm_k8s_relationships.py
git add server/meshmodel/aws-secretsmanager-controller/*/v1.0.0/relationships/edge-non-binding-network-deplo.json
git add server/meshmodel/aws-secretsmanager-controller/*/v1.0.0/relationships/edge-non-binding-network-podds.json
git add server/meshmodel/aws-secretsmanager-controller/*/v1.0.0/relationships/edge-non-binding-network-stsfs.json
git add server/meshmodel/aws-secretsmanager-controller/*/v1.0.0/relationships/edge-non-binding-network-dmnst.json
git add server/meshmodel/aws-secretsmanager-controller/*/v1.0.0/relationships/edge-non-binding-network-jobbs.json
git add server/meshmodel/aws-secretsmanager-controller/*/v1.0.0/relationships/edge-non-binding-network-crnjb.json
```

### Option 2: Add All at Once
```bash
git add scripts/fix_awssm_k8s_relationships.py
git add server/meshmodel/aws-secretsmanager-controller/*/v1.0.0/relationships/edge-non-binding-network-*.json
```

### Option 3: Bulk Add with Pattern
```bash
git add -A -- scripts/fix_awssm_k8s_relationships.py server/meshmodel/aws-secretsmanager-controller/
```

---

## Verification Before Commit

### 1. Count Files
```bash
git status --short | grep "aws-secretsmanager-controller" | wc -l
# Should be 78 (or more if including existing relationship files)

git status --short | grep "scripts/fix_awssm" | wc -l
# Should be 1
```

### 2. Verify No Zero UUIDs
```bash
grep -r "00000000-0000-0000-0000-000000000000" server/meshmodel/aws-secretsmanager-controller/v1.0.7/v1.0.0/relationships/edge-non-binding-network-*.json
# Should output nothing
```

### 3. Check Specific Files
```bash
# Deployment file should exist
ls server/meshmodel/aws-secretsmanager-controller/v1.2.0/v1.0.0/relationships/edge-non-binding-network-deplo.json
# Output: server/meshmodel/aws-secretsmanager-controller/v1.2.0/v1.0.0/relationships/edge-non-binding-network-deplo.json

# Pod file should exist  
ls server/meshmodel/aws-secretsmanager-controller/v1.2.1/v1.0.0/relationships/edge-non-binding-network-podds.json
# Output: server/meshmodel/aws-secretsmanager-controller/v1.2.1/v1.0.0/relationships/edge-non-binding-network-podds.json
```

---

## Commit Command

```bash
git commit -s -m "[Models] Fix AWS Secrets Manager relationships with unique UUIDs and add Deployment/Pod support

- Fixed zero UUID issue from PR #17189 by generating unique UUID v5
- Added AWS Secrets Manager Secret → Deployment relationships (all 13 versions)
- Added AWS Secrets Manager Secret → Pod relationships (all 13 versions)
- Updated StatefulSet, DaemonSet, Job, CronJob with unique UUIDs
- Applied to aws-secretsmanager-controller v1.0.7 through v1.2.1
- All files follow relationships.meshery.io/v1alpha3 schema
- Non-binding network edges for visual Kanvas connections

Fixes #17179"
```

---

## DO NOT COMMIT

The following files are documentation only - DO NOT add to git:
- `AWS_SECRETS_MANAGER_RELATIONSHIPS_SUMMARY.md`
- `QUICK_REFERENCE.md`
- `IMPLEMENTATION_COMPLETE.md`
- `CHANGES_OVERVIEW.txt`
- `FILES_TO_COMMIT.md` (this file)

These are helper documentation to understand what was done.

---

## Files Changed Summary

| Category | Count | Status |
|----------|-------|--------|
| Script Files | 1 | NEW |
| Relationship Files | 78 | NEW |
| Documentation (not to commit) | 5 | NEW |
| **TOTAL TO COMMIT** | **79** | **NEW** |

---

## What Each File Contains

### Script: fix_awssm_k8s_relationships.py
```python
# Features:
- Imports: os, json, uuid, typing
- Function: generate_unique_uuid_for_workload()
- Function: create_relationship()
- Function: generate_filename_suffix()
- Function: main()
- Supports: Deployment, Pod, StatefulSet, DaemonSet, Job, CronJob
- Versions: v1.0.7 through v1.2.1
```

### Each Relationship File
```json
{
  "id": "unique-uuid-v5-per-workload",
  "evaluationQuery": "$.spec.template.spec.volumes[*].secret.secretName" OR "$.spec.volumes[*].secret.secretName",
  "kind": "edge",
  "type": "non-binding",
  "subType": "network",
  "status": "enabled",
  "schemaVersion": "relationships.meshery.io/v1alpha3",
  "model": { "version": "v1.x.x", ... },
  "selectors": [ { "allow": { "from": [...], "to": [...] } } ]
}
```

---

## Quality Assurance Checklist

Before creating PR:
- [ ] All 79 files exist
- [ ] No zero UUIDs found
- [ ] All JSON files are valid
- [ ] Script runs without errors
- [ ] Relationships can be regenerated
- [ ] Files follow v1alpha3 schema
- [ ] Commit message is clear and references #17179
- [ ] DCO sign-off included (-s flag)

---

## References

**Issue**: GitHub Issue #17179 - Add AWS Secrets Manager to Deployment/Pod relationships  
**Previous PR**: #17189 - Had zero UUID issue and missing workload types  
**Schema**: relationships.meshery.io/v1alpha3  
**Models**: aws-secretsmanager-controller (v1.0.7 through v1.2.1)  

---

## Next Steps After Commit

1. Push branch: `git push origin <branch-name>`
2. Create PR on GitHub
3. Add description referencing the issue and previous PR
4. Link to this repository's contributing guidelines
5. Request review from Meshery maintainers
6. Address any feedback
7. Merge to master once approved

---

END OF FILE LIST
