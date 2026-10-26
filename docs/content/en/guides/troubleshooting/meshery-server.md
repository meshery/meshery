---
title: Troubleshooting Errors while running Meshery
description: Troubleshoot errors from mesheryctl system start and from running Meshery Server locally with make server.
aliases:
  - /guides/troubleshooting/running
categories: [troubleshooting]
---

If Meshery Server fails to start, first identify how Meshery is being run and then follow the troubleshooting steps for that startup or deployment method.

Meshery Server can be started in different ways. The troubleshooting steps differ depending on whether you are running Meshery locally from source, using `mesheryctl`, in Docker, or on Kubernetes.

## Common symptoms

| Symptom | Possible cause |
| --- | --- |
| Meshery Server does not start | Server, container, or dependency failure |
| Meshery Server is not reachable | Server/container is not running or the required port is unavailable |
| `mesheryctl system start` fails | Deployment, configuration, or platform-specific issue |
| `make server` fails | Local development environment, dependency, or database issue |

## Where to look for logs

Before troubleshooting, identify how Meshery Server is being run and use the corresponding section below.

### Local development with `make server`

Use `make server` when running Meshery Server locally from source.

```bash
make server
```

If the command fails, inspect the terminal output for errors related to:

- Missing or incompatible dependencies
- Local configuration
- Database initialization or migration
- Build failures
- Port conflicts

See the [`make server`](#make-server) troubleshooting example below for a known database-related startup failure.

### Meshery CLI with `mesheryctl`

When starting Meshery using the Meshery CLI, use the following commands to diagnose startup problems.

**1. Start Meshery**

Run the following command to start Meshery:

```bash
mesheryctl system start
```

**2. Check system prerequisites**

Run the preflight check to identify potential environment or configuration issues. For details, see [Running system checks]({{< ref "guides/mesheryctl/running-system-checks-using-mesheryctl.md" >}}).

```bash
mesheryctl system check --preflight
```

**3. Collect startup logs**

If Meshery fails to start or behaves unexpectedly, collect its logs for troubleshooting:

```bash
mesheryctl system logs
```

If preflight checks fail and the CLI cannot collect logs, inspect the deployment directly using the Docker or Kubernetes log commands described below.

If a command fails, inspect its output to determine whether the issue is related to the configured deployment platform, existing resources, configuration, or another startup dependency.

For installation-related issues, including preflight checks and Kubernetes connectivity, see the [installation troubleshooting guide]({{< ref "guides/troubleshooting/installation/index.md" >}}).

See the [`mesheryctl system start`](#mesheryctl-system-start) troubleshooting example below for a known Kubernetes resource ownership issue.

### Docker

If Meshery Server is running in Docker, list all Meshery containers, including stopped containers, to check their status:

```bash
docker ps -a --filter "name=meshery"
```

Identify the Meshery Server container from the output, then inspect its logs:

```bash
docker logs <container-name-or-id>
```

Look for errors related to:

- Container startup or unexpected restarts
- Port binding
- Configuration
- Database or dependency connections
- Application startup failures

### Kubernetes

If Meshery Server is running on Kubernetes, first check the status of the Meshery pods:

```bash
kubectl get pods -n meshery
```

Identify the affected Meshery Server pod from the output, then inspect its logs:

```bash
kubectl logs -n meshery <pod-name>
```

Inspect the pod's configuration, events, and status to identify scheduling, image-pull, or startup issues:

```bash
kubectl describe pod <pod-name> -n meshery
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
- Image pull failures
- Configuration errors
- Readiness or liveness probe failures
- Application startup errors

## `mesheryctl system start`

### Resources already exist and cannot be imported into the Helm release

**Error:**

```text
mesheryctl system start: cannot start Meshery: rendered manifests contain a resource that already exists.
Unable to continue with install: ServiceAccount "meshery-operator" in namespace "meshery" exists and cannot
be imported into the current release: invalid ownership metadata; label validation error: missing key
"app.kubernetes.io/managed-by": must be set to "Helm"; annotation validation error: missing key
"meta.helm.sh/release-name": must be set to "meshery"; annotation validation error: missing key
"meta.helm.sh/release-namespace": must be set to "meshery"
```

**Cause:** Meshery resources from an earlier installation are still in the cluster, but Helm does not manage them, so the new Helm release cannot take them over.

**Fix:** Stop Meshery and check for leftover resources before deleting anything.

1. Stop Meshery and allow the CLI to clean up its resources:

   ```bash
   mesheryctl system stop --force
   ```

2. If the error persists, inspect the namespace:

   ```bash
   kubectl get all -n meshery
   kubectl get serviceaccount -n meshery
   kubectl get clusterroles,clusterrolebindings | grep meshery
   ```

   Check the namespace and cluster-level resources for leftovers associated with Meshery. Delete only resources that you have confirmed belong to the obsolete Meshery installation. Do not delete a shared namespace or resources used by other applications.

   For example, if you have confirmed that these resources are obsolete Meshery resources, delete them explicitly:

   ```bash
   kubectl delete ns meshery
   kubectl delete clusterroles.rbac.authorization.k8s.io meshery-controller-role meshery-operator-role meshery-proxy-role meshery-metrics-reader
   kubectl delete clusterrolebindings.rbac.authorization.k8s.io meshery-controller-rolebinding meshery-operator-rolebinding meshery-proxy-rolebinding
   ```

   See [Meshery issue #4578](https://github.com/meshery/meshery/issues/4578) for the related issue.

3. Start Meshery again:

   ```bash
   mesheryctl system start
   ```

## `make server`

### Database migration errors

**Error:** Meshery Server exits during startup with a database error:

```text
FATA[0000] constraints not implemented on sqlite, consider using DisableForeignKeyConstraintWhenMigrating
exit status 1
make: *** [Makefile: server] Error 1
```

**Cause:** Meshery Server keeps its local SQLite database in `~/.meshery/config` by default. The `USER_DATA_FOLDER` environment variable overrides this location. A database created by a different Meshery version can fail to migrate.

**Fix:** Back up your local Meshery data before removing the affected database files. The following commands remove the database file and its SQLite write-ahead log files, which may contain locally stored Meshery data.

```bash
DB_DIR="${USER_DATA_FOLDER:-$HOME/.meshery/config}"
rm -f "$DB_DIR/mesherydb.sql" "$DB_DIR/mesherydb.sql-wal" "$DB_DIR/mesherydb.sql-shm"
```

Run the server again:

```bash
make server
```

For more information, see the [Error Code Reference]({{< ref "reference/references/error-codes.md" >}}).
