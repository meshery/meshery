---
title: Troubleshooting Errors while running Meshery
description: Troubleshooting Meshery errors when running make server / mesheryctl system start
aliases:
  - /guides/troubleshooting/running
categories: [troubleshooting]
---

## Meshery Server startup troubleshooting

If Meshery Server fails to start, first identify how Meshery is being run and then follow the troubleshooting steps for that startup or deployment method.

Meshery Server can be started in different ways. The troubleshooting steps differ depending on whether you are running Meshery locally from source, using `mesheryctl`, in Docker, or on Kubernetes.

### Common symptoms

| Symptom | Possible cause |
| --- | --- |
| Meshery Server does not start | Server, container, or dependency failure |
| Meshery Server is not reachable | Server/container is not running or the required port is unavailable |
| `mesheryctl system start` fails | Deployment, configuration, or platform-specific issue |
| `make server` fails | Local development environment, dependency, or database issue |

### Basic diagnostics

Before troubleshooting, identify how Meshery Server is being run and use the corresponding section below.

#### Local development with `make server`

Use `make server` when running Meshery Server locally from source.

```bash
make server
```

If the command fails, inspect the terminal output for errors related to:

- missing or incompatible dependencies
- local configuration
- database initialization or migration
- build failures
- port conflicts

See the [`make server`](#make-server) troubleshooting example below for a known database-related startup failure.

#### Meshery CLI with `mesheryctl`

When starting Meshery using the Meshery CLI, run:

```bash
mesheryctl system start
```

If the command fails, inspect the CLI output to determine whether the failure is related to the configured deployment platform, existing resources, configuration, or another startup dependency.

See the [`mesheryctl system start`](#mesheryctl-system-start) troubleshooting example below for a known Kubernetes resource ownership issue.

#### Docker

If Meshery Server is running in Docker, first verify that Docker is running and list the active containers:

```bash
docker ps
```

Identify the Meshery Server container from the output, then inspect its logs:

```bash
docker logs <container-name-or-id>
```

Look for errors related to:

- container startup or unexpected restarts
- port binding
- configuration
- database or dependency connections
- application startup failures

If the Meshery container is not listed among the running containers, also check stopped containers:

```bash
docker ps -a
```

#### Kubernetes

If Meshery Server is running on Kubernetes, first check the status of the Meshery pods:

```bash
kubectl get pods -n meshery
```

Identify the affected Meshery Server pod from the output, then inspect its logs:

```bash
kubectl logs -n meshery <pod-name>
```

If the pod contains multiple containers, list their names:

```bash
kubectl get pod <pod-name> -n meshery -o jsonpath='{.spec.containers[*].name}'
```

Then inspect the appropriate container:

```bash
kubectl logs -n meshery <pod-name> -c <container-name>
```

When checking the pod status and logs, look for:

- `CrashLoopBackOff`
- image pull failures
- configuration errors
- readiness or liveness probe failures
- application startup errors

## mesheryctl system start

**Error:**

```text
mesheryctl system start : : cannot start Meshery: rendered manifests contain a resource that already exists.
Unable to continue with install: ServiceAccount "meshery-operator" in namespace "meshery" exists and cannot
be imported into the current release: invalid ownership metadata; label validation error: missing key
"app.kubernetes.io/managed-by": must be set to "Helm"; annotation validation error: missing key
"meta.helm.sh/release-name": must be set to "meshery"; annotation validation error: missing key
"meta.helm.sh/release-namespace": must be set to "meshery"
```

**Fix: Clean the cluster using:**

<pre class="codeblock-pre"><div class="codeblock">
<div class="clipboardjs">
kubectl delete ns meshery
kubectl delete clusterroles.rbac.authorization.k8s.io meshery-controller-role meshery-operator-role meshery-proxy-role meshery-metrics-reader
kubectl delete clusterrolebindings.rbac.authorization.k8s.io meshery-controller-rolebinding meshery-operator-rolebinding meshery-proxy-rolebinding
</div></div>
</pre>

_Issue Reference: [#4578](https://github.com/meshery/meshery/issues/4578)_

### make server

**Error:**

```text
FATA[0000] constraints not implemented on sqlite, consider using DisableForeignKeyConstraintWhenMigrating, more details https://github.com/go-gorm/gorm/wiki/GORM-V2-Release-Note-Draft#all-new-migrator
exit status 1
make: *** [Makefile:76: server] Error 1
```

**Fix:**

1. Flush the database by deleting the `.meshery/config`.
2. Run:

```bash
make server
```

## Additional Resources

- [Error Code Reference]({{< ref "reference/references/error-codes.md" >}})