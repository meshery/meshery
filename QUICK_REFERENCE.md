# Quick Reference: AWS Secrets Manager Relationships Fix

## Summary
✅ Fixed UUID issue from PR #17189
✅ Added Deployment and Pod relationships to AWS Secrets Manager
✅ Generated 78 relationship files across 13 controller versions
✅ All files have unique UUIDs and follow v1alpha3 schema

## Files to Commit

### New Script
```
scripts/fix_awssm_k8s_relationships.py
```

### Relationship Files (78 files)
All files under:
```
server/meshmodel/aws-secretsmanager-controller/v*/v1.0.0/relationships/
```

Specifically:
- `edge-non-binding-network-deplo.json` (Deployment) - NEW
- `edge-non-binding-network-podds.json` (Pod) - NEW  
- `edge-non-binding-network-stsfs.json` (StatefulSet)
- `edge-non-binding-network-dmnst.json` (DaemonSet)
- `edge-non-binding-network-jobbs.json` (Job)
- `edge-non-binding-network-crnjb.json` (CronJob)

For versions: v1.0.7, v1.0.8, v1.0.9, v1.0.10, v1.0.11, v1.0.12, v1.0.13, v1.0.14, v1.1.0, v1.1.1, v1.1.2, v1.2.0, v1.2.1

## Key Differences from PR #17189

| Aspect | PR #17189 (Previous) | This Fix |
|--------|-------------------|----------|
| UUID Format | `00000000-0000-0000-0000-000000000000` | Unique UUID v5 |
| Deployments | ❌ Missing | ✅ Added |
| Pods | ❌ Missing | ✅ Added |
| StatefulSets | ✅ Present | ✅ Updated with unique UUID |
| DaemonSets | ✅ Present | ✅ Updated with unique UUID |
| Jobs | ✅ Present | ✅ Updated with unique UUID |
| CronJobs | ✅ Present | ✅ Updated with unique UUID |
| Versions | 2 (v1.2.0, v1.2.1) | 13 versions |

## How the UUID Generation Works

**Deterministic Method**: UUID v5 (namespace-based hashing)
```python
uuid.uuid5(uuid.NAMESPACE_DNS, f"{workload_kind}-{model_version}-aws-secretsmanager.meshery.io")
```

**Example**:
- Deployment + v1.2.0 → `dc589c39-343f-5a56-ae72-edd0fc9446a9`
- Pod + v1.2.0 → `f81b6d2f-5f09-5065-9cf9-b98fcd00104e`
- StatefulSet + v1.2.0 → Similar unique UUID

**Benefits**:
- Reproducible (same input = same UUID every time)
- Unique across all workloads and versions
- No conflicts with existing relationships
- No external dependencies needed

## Relationship Structure (v1alpha3)

Each relationship file contains:
```json
{
  "id": "unique-uuid-v5",
  "evaluationQuery": "JSONPath to find secrets in workload spec",
  "kind": "edge",
  "type": "non-binding",
  "subType": "network",
  "status": "enabled",
  "selectors": [
    {
      "allow": {
        "from": [{ "kind": "Deployment|Pod|StatefulSet|...", ... }],
        "to": [{ "kind": "Secret", ... }]
      }
    }
  ]
}
```

## Commit Message Template

```
[Models] Fix AWS Secrets Manager relationships with unique UUIDs and add Deployment/Pod support

- Fixed zero UUID issue from PR #17189 by generating unique UUID v5 for each relationship
- Added AWS Secrets Manager Secret → Deployment relationships for all 13 controller versions
- Added AWS Secrets Manager Secret → Pod relationships for all 13 controller versions  
- Updated existing StatefulSet, DaemonSet, Job, CronJob relationships with unique UUIDs
- Applied to aws-secretsmanager-controller v1.0.7 through v1.2.1
- All files follow relationships.meshery.io/v1alpha3 schema
- Non-binding network edge relationships for visual secret-to-workload connections in Meshery Kanvas

Fixes #17179
```

## Testing the Changes

### 1. Verify File Generation
```bash
# Count relationship files created
find server/meshmodel/aws-secretsmanager-controller -name "edge-non-binding-network-*.json" | wc -l
# Should output: 78 (13 versions × 6 workloads)
```

### 2. Verify UUID Format
```bash
# Check that no zero UUIDs exist
grep -r "00000000-0000-0000-0000-000000000000" server/meshmodel/aws-secretsmanager-controller/
# Should output: nothing (no matches)
```

### 3. Verify Schema Compliance
```bash
# Check valid JSON structure
find server/meshmodel/aws-secretsmanager-controller -name "edge-non-binding-network-*.json" -exec jq . {} \; > /dev/null
# If no error output, all files are valid JSON
```

### 4. Manual Verification
Open any file and verify:
- ✅ Has valid unique UUID in `id` field
- ✅ Has correct `evaluationQuery` path
- ✅ Has `"status": "enabled"`
- ✅ Has `"type": "non-binding"` and `"subType": "network"`
- ✅ Correct workload `kind` in selectors

## Integration with Meshery Kanvas

Once merged:
1. Deploy AWS Secrets Manager controller
2. Deploy a Kubernetes workload (any of the 6 supported types)
3. Meshery will automatically detect and visualize the relationship
4. Secrets will show connected to workloads in Kanvas

## Files for Review

When submitting the PR, make sure to include:

✅ `scripts/fix_awssm_k8s_relationships.py` - Generation script
✅ All 78 relationship JSON files from `server/meshmodel/aws-secretsmanager-controller/*/v1.0.0/relationships/`
❌ Do NOT include the AWS_SECRETS_MANAGER_RELATIONSHIPS_SUMMARY.md (this is just reference)
