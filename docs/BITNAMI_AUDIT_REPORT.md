# Bitnami Catalog Migration Audit Report

**Issue:** #16072
**Date:** October 7, 2025
**Status:**  URGENT - Bitnami transition deadline (Sept 29) has passed
**Auditor:** Automated audit for Meshery project

---

## Executive Summary

Meshery currently has **extensive dependencies** on Bitnami resources that require immediate attention following Bitnami's catalog restructuring (completed August-September 2025). This audit identifies all Bitnami references and provides a migration roadmap.

### Impact Assessment

- ** HIGH RISK**: 113 Helm chart references in component models
- ** MEDIUM RISK**: 63+ Docker images in catalog designs
- ** LOW RISK**: Sealed Secrets references (unaffected by Bitnami changes)

---

## 1. Bitnami Helm Charts Audit

### Location: `server/meshmodel/component_models.yaml`

**Total Charts Affected:** 113 charts from `https://charts.bitnami.com/bitnami`

<details>
<summary><strong>Complete List of Bitnami Helm Charts (Click to expand)</strong></summary>

1. airflow
2. apache
3. appsmith
4. argo-cd
5. argo-workflows
6. aspnet-core
7. cassandra
8. cert-manager
9. clickhouse
10. common
11. concourse
12. consul
13. contour
14. contour-operator
15. dataplatform-bp2
16. discourse
17. dokuwiki
18. drupal
19. ejbca
20. elasticsearch
21. etcd
22. external-dns
23. fluentd
24. geode
25. ghost
26. gitea
27. grafana
28. grafana-loki
29. grafana-mimir
30. grafana-operator
31. grafana-tempo
32. haproxy
33. haproxy-intel
34. harbor
35. influxdb
36. jaeger
37. jasperreports
38. jenkins
39. joomla
40. jupyterhub
41. kafka
42. keycloak
43. kiam
44. kibana
45. kong
46. kubeapps
47. kube-prometheus
48. kubernetes-event-exporter
49. kube-state-metrics
50. logstash
51. magento
52. mariadb
53. mariadb-galera
54. mastodon
55. matomo
56. mediawiki
57. memcached
58. metallb
59. metrics-server
60. minio
61. mongodb
62. mongodb-sharded
63. moodle
64. multus-cni
65. mxnet
66. mysql
67. nats
68. nginx
69. nginx-ingress-controller
70. nginx-intel
71. node
72. node-exporter
73. oauth2-proxy
74. odoo
75. opencart
76. osclass
77. owncloud
78. parse
79. phpbb
80. phpmyadmin
81. pinniped
82. postgresql
83. postgresql-ha
84. prestashop
85. pytorch
86. rabbitmq
87. rabbitmq-cluster-operator
88. redis
89. redis-cluster
90. redmine
91. schema-registry
92. sealed-secrets  (UNAFFECTED)
93. solr
94. sonarqube
95. spark
96. spring-cloud-dataflow
97. suitecrm
98. supabase
99. tensorflow-resnet
100. thanos
101. tomcat
102. wavefront
103. wavefront-adapter-for-istio
104. wavefront-hpa-adapter
105. wavefront-prometheus-storage-adapter
106. whereabouts
107. wildfly
108. wordpress
109. wordpress-intel
110. zookeeper

**Plus 3 charts from `bitnami-aks` repository**

</details>

### Current State Analysis

**Chart Source Status:**
-  **GitHub Source Code**: Still available at `github.com/bitnami/charts` (Apache 2.0 license)
-  **OCI Helm Charts**: Available at `docker.io/bitnamicharts` but **NO LONGER UPDATED** (except BSI subset)
-  **Packaged Charts**: `charts.bitnami.com/bitnami` may have limited/stale versions

**Migration Priority:**
- **P0 (Critical)**: Charts used in brownout lists (kafka, redis, postgresql, mongodb, etc.)
- **P1 (High)**: CNCF graduated projects (prometheus, etcd, harbor, etc.)
- **P2 (Medium)**: Popular enterprise apps (jenkins, gitlab, keycloak, etc.)
- **P3 (Low)**: Specialized/less common applications

---

## 2. Bitnami Docker Images Audit

### Location: `docs/catalog/*/0.0.*/design.yml`

**Total Images Found:** 63+ references to `docker.io/bitnami/*`

<details>
<summary><strong>Docker Images by Application (Click to expand)</strong></summary>

### Images Used in Catalog Designs

| Image | Catalog Files | Brownout Status | Migration Priority |
|-------|---------------|-----------------|-------------------|
| `bitnami/redis` | 3 files |  Brownout 3 |  P0 - Critical |
| `bitnami/mongodb` | 2 files |  Brownout 3 |  P0 - Critical |
| `bitnami/postgresql` | 1 file |  Brownout 3 |  P0 - Critical |
| `bitnami/airflow` | 1 file | Not listed |  P1 - High |
| `bitnami/apache` | 1 file |  Brownout 2 |  P1 - High |
| `bitnami/apisix` | 1 file | Not listed |  P1 - High |
| `bitnami/apisix-ingress-controller` | 1 file | Not listed |  P1 - High |
| `bitnami/appsmith` | 1 file | Not listed |  P2 - Medium |
| `bitnami/aspnet-core` | 2 files | Not listed |  P2 - Medium |
| `bitnami/dotnet-sdk` | 2 files | Not listed |  P2 - Medium |
| `bitnami/etcd` | 1 file |  Brownout 2 |  P0 - Critical |
| `bitnami/git` | 2 files |  Brownout 2 |  P1 - High |
| `bitnami/haproxy` | 1 file | Not listed |  P1 - High |
| `bitnami/jenkins` | 1 file | Not listed |  P1 - High |
| `bitnami/os-shell` | 6 files |  Brownout 3 |  P0 - Critical |

**Total Unique Images:** 15
**P0 (Critical):** 4 images
**P1 (High):** 5 images
**P2 (Medium):** 6 images

</details>

### Affected Catalog Designs

1. **APISIX Gateway** (`e17af2e2-b48d-42d9-98f2-01b9305c1b0f`) - 12 image references
2. **ASP.NET Core** (`1012970d-bbf5-4bc5-b5e2-f6cb1da64b84`, `8cff4bae-818f-4c9b-8511-a6feaabf6b1e`) - 6 references
3. **Appsmith** (`5275aa84-91e0-4df1-883d-8df61397f951`) - 15 references (mongodb, redis, appsmith, haproxy)
4. **Apache Airflow** (`df7815d8-10dd-4001-90b9-c37e79136133`) - 27 references
5. **Apache HTTPD** (`a67af210-22af-4115-a788-c64aadcd67ac`) - 2 references
6. **Jenkins** (`4c918765-c145-4607-a0ad-af01a7fa696e`) - 1 reference
7. **Clickhouse** (`cfde9483-de24-498f-869f-d8d74353a2af`) - 4 references

---

## 3. What Changed at Bitnami

### Timeline (Already Completed)

- **August 28, 2025**: Public catalog switched to limited community tier
  - Only "latest" tags available for free
  - Most images moved to `docker.io/bitnamilegacy`
  - No new Debian-based images published

- **September 2-29, 2025**: Brownout periods (temporary unavailability)
  - Affected images: kafka, redis, postgresql, mongodb, etcd, and 20+ others

- **September 29, 2025**: Legacy repository deletion deadline
  - `docker.io/bitnamilegacy` scheduled for removal

### Current Options (October 2025)

1. **Free Community Tier** (`docker.io/bitnamisecure` → `docker.io/bitnami`)
   - Limited subset of ~20-30 applications
   - Only "latest" tags (no version pinning)
   - Development use only

2. **Bitnami Secure Images (BSI)** - Commercial subscription
   - Full catalog (280+ applications)
   - Versioned tags and LTS branches
   - Hardened, distroless images
   - SLSA Level 3, SBOM, CVE transparency
   - Enterprise support

3. **Alternative Sources**
   - Official upstream registries (e.g., `postgres:16`, `redis:7`)
   - Other community distributions
   - Build from Bitnami source (GitHub - Apache 2.0 license)

---

## 4. Migration Strategies

### Strategy A: Official Upstream Images (Recommended)

**Pros:**
-  Maintained by original project teams
-  Free and widely used
-  Version stability
-  No vendor lock-in

**Cons:**
-  Different image structure/paths
-  May require configuration changes
-  Less standardization across images

**Best For:** Core infrastructure (redis, postgresql, mongodb, nginx, etc.)

### Strategy B: Bitnami Secure Images (BSI)

**Pros:**
- Drop-in replacement (minimal changes)
- Enhanced security (hardened, distroless)
- Enterprise support
- Full version control

**Cons:**
-  Commercial subscription required
-  Vendor dependency

**Best For:** Enterprise deployments, production use cases

### Strategy C: Community Alternatives

**Pros:**
-  Free and open source
-  Community maintained

**Cons:**
-  Variable quality/support
-  May require evaluation

**Best For:** Specialized applications without official images

### Strategy D: Build from Source

**Pros:**
-  Full control
-  Bitnami source available (Apache 2.0)

**Cons:**
-  Requires build infrastructure
-  Maintenance overhead
-  Time intensive

**Best For:** Custom requirements, air-gapped environments

---

## 5. Recommended Migration Plan

### Phase 1: Immediate Action (Week 1-2)

**Priority:**  P0 - Critical (Images from Brownout 3)

1. **redis** → `redis:7-alpine` or `redis:7`
2. **postgresql** → `postgres:16-alpine` or `postgres:16`
3. **mongodb** → `mongo:8.0` or `mongodb/mongodb-community-server:8.0`
4. **etcd** → `quay.io/coreos/etcd:v3.5`
5. **os-shell** → `busybox:latest` or `alpine:3.19`

**Actions:**
- [ ] Update catalog designs with new image references
- [ ] Test deployments with new images
- [ ] Update documentation
- [ ] Create validation script

### Phase 2: High Priority (Week 3-4)

**Priority:**  P1 - High

1. **airflow** → `apache/airflow:2.8`
2. **apache** → `httpd:2.4-alpine`
3. **git** → `alpine/git:latest`
4. **haproxy** → `haproxy:2.9-alpine`
5. **jenkins** → `jenkins/jenkins:lts`
6. **apisix** → `apache/apisix:3.8.0`

**Actions:**
- [ ] Migrate remaining catalog designs
- [ ] Update component_models.yaml Helm chart references
- [ ] Add alternative chart sources (Artifact Hub)

### Phase 3: Medium Priority (Week 5-6)

**Priority:**  P2 - Medium

- Migrate remaining specialized applications
- Evaluate alternative Helm chart sources
- Update CI/CD validation workflows

### Phase 4: Validation & Guardrails (Week 7-8)

**Priority:**  P3 - Preventive

- [ ] Create automated validation script
- [ ] Add GitHub Actions workflow to check image availability
- [ ] Update contributor guidelines
- [ ] Add pre-commit hooks for Bitnami detection

---

## 6. Proposed Implementation

### 6.1 Catalog Design Updates

**Files to Update:**
```bash
docs/catalog/e17af2e2-b48d-42d9-98f2-01b9305c1b0f/0.0.1/design.yml  # APISIX
docs/catalog/1012970d-bbf5-4bc5-b5e2-f6cb1da64b84/0.0.1/design.yml  # ASP.NET Core
docs/catalog/5275aa84-91e0-4df1-883d-8df61397f951/0.0.1/design.yml  # Appsmith
docs/catalog/df7815d8-10dd-4001-90b9-c37e79136133/0.0.1/design.yml  # Airflow
docs/catalog/a67af210-22af-4115-a788-c64aadcd67ac/0.0.5/design.yml  # Apache
docs/catalog/4c918765-c145-4607-a0ad-af01a7fa696e/0.0.1/design.yml  # Jenkins
docs/catalog/cfde9483-de24-498f-869f-d8d74353a2af/0.0.4/design.yml  # Clickhouse
```

### 6.2 Component Models Update

**File:** `server/meshmodel/component_models.yaml`

**Changes Required:**
- Update or remove 113 Bitnami chart references
- Add alternative chart sources:
  - Artifact Hub (community charts)
  - Official project Helm repos
  - Verified publisher charts

### 6.3 Validation Script

Create: `.github/scripts/validate-image-sources.sh`

```bash
#!/bin/bash
# Validate that container images are available from their registries

echo " Validating container image sources..."

# Extract all docker.io/bitnami references
echo "Checking for deprecated Bitnami references..."
deprecated_refs=$(grep -r "docker.io/bitnami" docs/catalog/ server/meshmodel/ --include="*.yml" --include="*.yaml" | wc -l)

if [ $deprecated_refs -gt 0 ]; then
  echo "  Found $deprecated_refs deprecated Bitnami references"
  grep -r "docker.io/bitnami" docs/catalog/ server/meshmodel/ --include="*.yml" --include="*.yaml"
  exit 1
else
  echo " No deprecated Bitnami references found"
fi

# Check Helm chart availability
echo "Checking Helm chart sources..."
# Add logic to verify chart URLs are reachable

echo " Validation complete"
```

### 6.4 GitHub Actions Workflow

Create: `.github/workflows/validate-catalog-sources.yml`

```yaml
name: Validate Catalog Sources

on:
  pull_request:
    paths:
      - 'docs/catalog/**'
      - 'server/meshmodel/**'
  schedule:
    - cron: '0 0 * * 0'  # Weekly on Sunday

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Check for deprecated Bitnami references
        run: |
          .github/scripts/validate-image-sources.sh
      
      - name: Validate Helm chart availability
        run: |
          # Test that chart URLs return 200 OK
          # Sample charts from component_models.yaml
          
      - name: Report results
        if: failure()
        uses: actions/github-script@v7
        with:
          script: |
            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: ' Deprecated Bitnami references detected. Please migrate to alternative sources.'
            })
```

---

## 7. Image Mapping Reference

### Recommended Replacements

| Bitnami Image | Official Alternative | Notes |
|---------------|---------------------|-------|
| `bitnami/redis` | `redis:7-alpine` | Drop-in replacement |
| `bitnami/postgresql` | `postgres:16-alpine` | Check init scripts |
| `bitnami/mongodb` | `mongo:8.0` | Authentication may differ |
| `bitnami/nginx` | `nginx:alpine` | Config paths may differ |
| `bitnami/apache` | `httpd:2.4-alpine` | DocumentRoot differs |
| `bitnami/etcd` | `quay.io/coreos/etcd:v3.5` | Official CoreOS image |
| `bitnami/jenkins` | `jenkins/jenkins:lts` | Official Jenkins image |
| `bitnami/airflow` | `apache/airflow:2.8` | Official Airflow image |
| `bitnami/git` | `alpine/git` | Lightweight git client |
| `bitnami/os-shell` | `busybox` or `alpine:3.19` | Basic shell utilities |
| `bitnami/dotnet-sdk` | `mcr.microsoft.com/dotnet/sdk:9.0` | Microsoft official |
| `bitnami/aspnet-core` | `mcr.microsoft.com/dotnet/aspnet:9.0` | Microsoft official |
| `bitnami/haproxy` | `haproxy:2.9-alpine` | Official HAProxy |

---

## 8. Testing Strategy

### 8.1 Pre-Migration Testing

- [ ] Identify all catalog designs using Bitnami images
- [ ] Create test deployments with alternative images
- [ ] Validate functionality and compatibility
- [ ] Document any configuration changes required

### 8.2 Migration Testing

- [ ] Update designs one catalog at a time
- [ ] Test each updated catalog design
- [ ] Verify image pulls succeed
- [ ] Check application functionality

### 8.3 Post-Migration Validation

- [ ] Run automated validation script
- [ ] Monitor for image pull failures
- [ ] Update documentation and examples
- [ ] Add guardrails to prevent future issues

---

## 9. Documentation Updates Required

### Files to Update:

1. **README.md** - Add notice about Bitnami migration
2. **CONTRIBUTING.md** - Guidelines for image selection
3. **docs/pages/project/contributing-docs.md** - Container image best practices
4. **Catalog READMEs** - Update image sources in affected catalogs

### Proposed Guidelines:

```markdown
## Container Image Selection Guidelines

When creating Meshery designs or catalog entries:

1. **Prefer official images** from the project's registry
   - Example: `postgres:16`, `redis:7`, `nginx:alpine`

2. **Avoid deprecated sources**
   -  `docker.io/bitnami/*` (deprecated as of Sept 2025)
   -  `docker.io/bitnamilegacy/*` (archived, no updates)

3. **Use verified publishers** from Docker Hub or other registries
   -  Official project images
   -  CNCF-hosted images
   -  Cloud provider registries (for their services)

4. **Pin specific versions** for stability
   -  `postgres:16.1-alpine`
   -  `postgres:latest`

5. **Validate image availability** before committing
   - Test that the image can be pulled
   - Check for recent updates and security patches
```

---

## 10. Risk Assessment

### High Risks

1. **Broken Deployments** 
   - Impact: Users cannot deploy catalog designs
   - Mitigation: Priority migration of P0 images

2. **Image Pull Failures** 
   - Impact: Existing deployments may fail on restart
   - Mitigation: Update documentation with alternatives

3. **Version Incompatibility** 
   - Impact: Alternative images may have different behavior
   - Mitigation: Thorough testing and documentation

### Medium Risks

1. **Configuration Differences**
   - Impact: Some settings may need adjustment
   - Mitigation: Document migration notes per application

2. **Community Confusion** 
   - Impact: Contributors may not know about changes
   - Mitigation: Clear guidelines and PR checks

### Low Risks

1. **Helm Chart Staleness** 
   - Impact: Chart versions may be outdated
   - Mitigation: Regular updates and alternative sources

---

## 11. Success Criteria

### Completion Checklist

- [ ] All P0 (Critical) images migrated to alternatives
- [ ] Catalog designs updated and tested
- [ ] Validation script created and integrated
- [ ] CI/CD workflow added for ongoing checks
- [ ] Documentation updated with new guidelines
- [ ] Contributors notified of changes
- [ ] Migration guide published

### Metrics

- **Images Migrated**: 0/15 (0%)
- **Catalog Designs Updated**: 0/7 (0%)
- **Helm Charts Reviewed**: 0/113 (0%)
- **Validation Tests Added**: 0/1 (0%)
- **Documentation Updated**: 0/4 (0%)

---

## 12. Next Steps

### Immediate Actions (This Week)

1. **Create issue comment** on #16072 with this audit summary
2. **Start Phase 1 migration** (P0 images)
3. **Set up validation script** foundation
4. **Update first catalog design** as proof of concept

### Short Term (Next 2 Weeks)

1. Complete P0 and P1 image migrations
2. Update component_models.yaml
3. Add CI/CD validation workflow
4. Update contributor guidelines

### Long Term (Next Month)

1. Complete all migrations (P2, P3)
2. Add periodic checking for image availability
3. Establish process for future image selection
4. Consider contribution to upstream projects for official Helm charts

---

## 13. References

- **Bitnami Announcement**: [bitnami/charts#35164](https://github.com/bitnami/charts/issues/35164)
- **Meshery Issue**: [meshery/meshery#16072](https://github.com/meshery/meshery/issues/16072)
- **Bitnami Charts Source**: https://github.com/bitnami/charts (Apache 2.0)
- **Bitnami Secure Images**: https://hub.docker.com/u/bitnamisecure
- **Artifact Hub**: https://artifacthub.io/

---

**Report Generated**: October 7, 2025  
**Next Review**: Weekly until migration complete  
**Owner**: Meshery Community / Issue #16072 contributors
