---
title: Troubleshooting Errors while running Meshery
description: Troubleshooting Meshery errors when running make run-fast / meshery system start
aliases: 
- /guides/troubleshooting/running
categories: [troubleshooting]
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

_Issue Reference : [https://github.com/meshery/meshery/issues/4578](https://github.com/meshery/meshery/issues/4578)_

### make server

**Error:**

```
FATA[0000] constraints not implemented on sqlite, consider using DisableForeignKeyConstraintWhenMigrating, more details https://github.com/go-gorm/gorm/wiki/GORM-V2-Release-Note-Draft#all-new-migrator
exit status 1
make: *** [Makefile:76: server] Error 1
```

**Fix:**

1. Flush the database by deleting the `.meshery/config`
2. `make server`

#### See Also

- [Error Code Reference]({{< ref "reference/references/error-codes.md" >}})

## Meshery Server startup checklist

If Meshery Server does not become ready, first confirm that the runtime is
available and then inspect the server logs:

### Docker

```bash
docker info
docker ps -a --filter name=meshery
docker logs meshery --tail 200
```

Check that the configured port is free and that the container has not exited.
Use `docker inspect meshery` to review its environment and exit status.

### Kubernetes

```bash
kubectl get pods -n meshery
kubectl get events -n meshery --sort-by=.lastTimestamp
kubectl logs -n meshery deploy/meshery --all-containers --tail=200
kubectl describe pod -n meshery -l app=meshery
```

If the pod is pending, check scheduling, image-pull, and volume events. If it
restarts, inspect the container's previous logs with `--previous`. A port
forward can confirm that the service is reachable from the local machine:

```bash
kubectl port-forward -n meshery svc/meshery 9081:9081
```

These checks identify runtime, connectivity, configuration, and dependency
failures without requiring access to production systems.
