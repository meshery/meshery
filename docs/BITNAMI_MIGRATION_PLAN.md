# Bitnami Migration Implementation Plan

**Related Issue:** #16072  
**Started:** October 7, 2025  
**Target Completion:** November 7, 2025 (4 weeks)

---

## Current Status

**Validation Results (Oct 7, 2025):**
- 63 deprecated Docker image references in catalog designs
- 110 Bitnami Helm chart references in component_models.yaml
- 15 unique image types need migration

---

## Phase 1: Critical Images (Week 1)

### P0 Images - From Brownout Lists

| Image | Occurrences | Target Replacement | Status | Assignee |
|-------|-------------|-------------------|---------|----------|
| `bitnami/redis` | 3 | `redis:7-alpine` | ⏳ TODO | - |
| `bitnami/postgresql` | 1 | `postgres:16-alpine` | ⏳ TODO | - |
| `bitnami/mongodb` | 2 | `mongo:8.0` | ⏳ TODO | - |
| `bitnami/etcd` | 1 | `quay.io/coreos/etcd:v3.5` | ⏳ TODO | - |
| `bitnami/os-shell` | 6 | `busybox:latest` or `alpine:3.19` | ⏳ TODO | - |

**Affected Files:**
- `docs/catalog/5275aa84-91e0-4df1-883d-8df61397f951/0.0.1/design.yml` (Appsmith)
- `docs/catalog/df7815d8-10dd-4001-90b9-c37e79136133/0.0.1/design.yml` (Airflow)
- `docs/catalog/e17af2e2-b48d-42d9-98f2-01b9305c1b0f/0.0.1/design.yml` (APISIX)
- `docs/catalog/cfde9483-de24-498f-869f-d8d74353a2af/0.0.4/design.yml` (Clickhouse)

---

## Phase 2: High Priority Images (Week 2)

### P1 Images - Popular Applications

| Image | Occurrences | Target Replacement | Status | Assignee |
|-------|-------------|-------------------|---------|----------|
| `bitnami/airflow` | 16 | `apache/airflow:2.8` | ⏳ TODO | - |
| `bitnami/apache` | 2 | `httpd:2.4-alpine` | ⏳ TODO | - |
| `bitnami/git` | 2 | `alpine/git:latest` | ⏳ TODO | - |
| `bitnami/haproxy` | 1 | `haproxy:2.9-alpine` | ⏳ TODO | - |
| `bitnami/jenkins` | 1 | `jenkins/jenkins:lts` | ⏳ TODO | - |
| `bitnami/apisix` | 3 | `apache/apisix:3.8.0` | ⏳ TODO | - |
| `bitnami/apisix-ingress-controller` | 1 | `apache/apisix-ingress-controller:1.8` | ⏳ TODO | - |

**Affected Files:**
- `docs/catalog/df7815d8-10dd-4001-90b9-c37e79136133/0.0.1/design.yml` (Airflow - 16 refs)
- `docs/catalog/a67af210-22af-4115-a788-c64aadcd67ac/0.0.5/design.yml` (Apache - 2 refs)
- `docs/catalog/4c918765-c145-4607-a0ad-af01a7fa696e/0.0.1/design.yml` (Jenkins - 1 ref)
- `docs/catalog/e17af2e2-b48d-42d9-98f2-01b9305c1b0f/0.0.1/design.yml` (APISIX - 4 refs)

---

## Phase 3: Medium Priority Images (Week 3)

### P2 Images - Specialized Applications

| Image | Occurrences | Target Replacement | Status | Assignee |
|-------|-------------|-------------------|---------|----------|
| `bitnami/appsmith` | 10 | Official appsmith image or build from source | ⏳ TODO | - |
| `bitnami/aspnet-core` | 2 | `mcr.microsoft.com/dotnet/aspnet:9.0` | ⏳ TODO | - |
| `bitnami/dotnet-sdk` | 2 | `mcr.microsoft.com/dotnet/sdk:9.0` | ⏳ TODO | - |

**Affected Files:**
- `docs/catalog/5275aa84-91e0-4df1-883d-8df61397f951/0.0.1/design.yml` (Appsmith - 10 refs)
- `docs/catalog/1012970d-bbf5-4bc5-b5e2-f6cb1da64b84/0.0.1/design.yml` (ASP.NET - 3 refs)
- `docs/catalog/8cff4bae-818f-4c9b-8511-a6feaabf6b1e/0.0.1/design.yml` (ASP.NET - 3 refs)

---

## Phase 4: Helm Charts Migration (Week 4)

### Component Models Update

**File:** `server/meshmodel/component_models.yaml`

**Strategy:**
1. Identify charts still available from Artifact Hub
2. Find official project Helm repositories
3. Remove charts no longer maintained
4. Add verification for chart availability

**Top Priority Charts (CNCF/Popular):**
- [ ] cert-manager → Official cert-manager.io chart
- [ ] external-dns → Official Kubernetes SIGs
- [ ] harbor → Official goharbor.io chart
- [ ] prometheus → Official prometheus-community
- [ ] grafana → Official grafana.com chart
- [ ] kafka → Official Apache/Strimzi charts
- [ ] redis → Official Redis charts
- [ ] postgresql → Official PostgreSQL charts

---

## Catalog Design Migration Checklist

### For Each Catalog Design:

- [ ] Identify all Bitnami image references
- [ ] Research official upstream images
- [ ] Update image references in design.yml
- [ ] Test image availability (docker pull)
- [ ] Update any init container images
- [ ] Update sidecar container images
- [ ] Verify environment variables compatibility
- [ ] Check volume mount paths
- [ ] Update ConfigMap/Secret references if needed
- [ ] Test deployment (if possible)
- [ ] Update catalog README/documentation
- [ ] Create PR with changes
- [ ] Request review

---

## Migration Script Template

```bash
#!/bin/bash
# Example migration for a catalog design

FILE="docs/catalog/[UUID]/[VERSION]/design.yml"

# Backup
cp "$FILE" "$FILE.backup"

# Replace images
sed -i '.bak' 's|docker\.io/bitnami/redis:[0-9.-]*|redis:7-alpine|g' "$FILE"
sed -i '.bak' 's|docker\.io/bitnami/postgresql:[0-9.-]*|postgres:16-alpine|g' "$FILE"
sed -i '.bak' 's|docker\.io/bitnami/mongodb:[0-9.-]*|mongo:8.0|g' "$FILE"
sed -i '.bak' 's|docker\.io/bitnami/etcd:[0-9.-]*|quay.io/coreos/etcd:v3.5|g' "$FILE"
sed -i '.bak' 's|docker\.io/bitnami/os-shell:[0-9.-]*|busybox:latest|g' "$FILE"

# Verify changes
diff "$FILE.backup" "$FILE"

# Validate YAML syntax
yamllint "$FILE"

echo "Migration complete for $FILE"
```

---

## Testing Strategy

### Per-Image Testing

1. **Image Availability**
   ```bash
   docker pull redis:7-alpine
   docker pull postgres:16-alpine
   docker pull mongo:8.0
   # etc.
   ```

2. **Compatibility Check**
   - Compare environment variables
   - Check exposed ports
   - Verify volume paths
   - Test health check endpoints

3. **Deployment Test** (where possible)
   - Deploy updated catalog design
   - Verify pods start successfully
   - Check application functionality
   - Monitor logs for errors

### Validation Checklist

- [ ] All images pull successfully
- [ ] YAML syntax is valid
- [ ] No Bitnami references remain (except sealed-secrets)
- [ ] Documentation updated
- [ ] Validation script passes
- [ ] CI/CD checks pass

---

## Progress Tracking

### Week 1 (Oct 7-13)
- [x] Complete audit report
- [x] Create validation script
- [x] Set up CI/CD workflow
- [ ] Begin P0 image migrations
- [ ] Test first catalog update

### Week 2 (Oct 14-20)
- [ ] Complete P0 migrations
- [ ] Begin P1 migrations
- [ ] Update documentation
- [ ] Create migration guide for contributors

### Week 3 (Oct 21-27)
- [ ] Complete P1 migrations
- [ ] Begin P2 migrations
- [ ] Update component_models.yaml
- [ ] Test Helm chart alternatives

### Week 4 (Oct 28 - Nov 3)
- [ ] Complete all migrations
- [ ] Final validation pass
- [ ] Update all documentation
- [ ] Create summary report

---

## Risk Mitigation

### Potential Issues & Solutions

| Risk | Impact | Mitigation |
|------|--------|------------|
| Image incompatibility | High | Test each replacement thoroughly |
| Breaking changes | High | Document all changes, provide migration notes |
| Missing images | Medium | Research alternatives, build from source if needed |
| Configuration differences | Medium | Update ConfigMaps, provide examples |
| Volume path changes | Low | Document path mappings |

---

## Communication Plan

### Updates

- **Weekly**: Update this document with progress
- **Per Phase**: Create summary comment on #16072
- **Blockers**: Immediately flag in #16072

### Documentation

- [ ] Update README with Bitnami deprecation notice
- [ ] Update CONTRIBUTING.md with image guidelines
- [ ] Create migration guide for existing deployments
- [ ] Update catalog documentation

---

## Success Metrics

-  **0 deprecated Bitnami image references** (currently: 63)
-  **All catalog designs deployable** with new images
-  **Validation script passes** in CI/CD
-  **Documentation complete** and up-to-date
-  **Community informed** of changes

---

## Resources

- [Bitnami Audit Report](./BITNAMI_AUDIT_REPORT.md)
- [Bitnami Announcement #35164](https://github.com/bitnami/charts/issues/35164)
- [Meshery Issue #16072](https://github.com/meshery/meshery/issues/16072)
- [Validation Script](./.github/scripts/validate-image-sources.sh)
- [CI/CD Workflow](./.github/workflows/validate-catalog-sources.yml)

---

**Last Updated:** October 7, 2025  
**Next Review:** October 14, 2025
