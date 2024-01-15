---
layout: default
title: Troubleshooting Errors while running Meshery 
abstract: Troubleshooting Meshery errors when running make run-fast / meshery system start 
permalink: guides/troubleshooting/meshery-server
redirect_from: guides/troubleshooting/running
type: guides
category: troubleshooting
language: en
---

## mesheryctl system start 
**Error:**
```
mesheryctl system start : : cannot start Meshery: rendered manifests contain a resource that already exists. 
Unable to continue with install: ServiceAccount "meshery-operator" in namespace "meshery" exists and cannot 
be imported into the current release: invalid ownership metadata; label validation error: missing key 
"app.kubernetes.io/managed-by": must be set to "Helm"; annotation validation error: missing key 
"meta.helm.sh/release-name": must be set to "meshery"; annotation validation error: missing key 
"meta.helm.sh/release-namespace": must be set to "meshery"
```
 
**(Fix) Clean the cluster using :**
 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
kubectl delete ns meshery
kubectl delete clusterroles.rbac.authorization.k8s.io meshery-controller-role meshery-operator-role meshery-proxy-role meshery-metrics-reader
kubectl delete clusterrolebindings.rbac.authorization.k8s.io meshery-controller-rolebinding meshery-operator-rolebinding meshery-proxy-rolebinding
 </div></div>
 </pre>
 *Issue Reference : [https://github.com/meshery/meshery/issues/4578](https://github.com/meshery/meshery/issues/4578)*

### make server
**Error:**
```
FATA[0000] constraints not implemented on sqlite, consider using DisableForeignKeyConstraintWhenMigrating, more details https://github.com/go-gorm/gorm/wiki/GORM-V2-Release-Note-Draft#all-new-migrator 
exit status 1
make: *** [Makefile:76: server] Error 1
```

**Fix:**
1. Flush the  database by deleting the `.meshery/config`
2. `make server`

#### See Also

- [Error Code Reference](/reference/error-codes)