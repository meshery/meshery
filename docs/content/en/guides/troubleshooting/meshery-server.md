---
title: Troubleshooting Errors while running Meshery
description: Troubleshoot errors from mesheryctl system start and from running Meshery Server locally with make server.
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

```bash
mesheryctl system check --preflight
```

```bash
mesheryctl system logs
```

If the command fails, inspect the CLI output to determine whether the failure is related to the configured deployment platform, existing resources, configuration, or another startup dependency.

See the [`mesheryctl system start`](#mesheryctl-system-start) troubleshooting example below for a known Kubernetes resource ownership issue.

#### Docker

If Meshery Server is running in Docker, first verify that Docker is running and list the active containers:

```bash
docker ps -a --filter "name=meshery"
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
docker ps -a --filter "name=meshery"
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
- image pull failures
- configuration errors
- readiness or liveness probe failures
- application startup errors

## mesheryctl system start

### Resources already exist and cannot be imported into the Helm release

**Error:**

```text
mesheryctl system start : : cannot start Meshery: rendered manifests contain a resource that already exists.
Unable to continue with install: ServiceAccount "meshery-operator" in namespace "meshery" exists and cannot
be imported into the current release: invalid ownership metadata; label validation error: missing key
"app.kubernetes.io/managed-by": must be set to "Helm"; annotation validation error: missing key
"meta.helm.sh/release-name": must be set to "meshery"; annotation validation error: missing key
"meta.helm.sh/release-namespace": must be set to "meshery"
```

**Cause:** Meshery resources from an earlier installation are still in the cluster, but Helm does not manage them, so the new Helm release cannot take them over.

**Fix:** Remove the leftover resources, then start Meshery again.

1. Let `mesheryctl` remove every Meshery resource it finds, including cluster roles and cluster role bindings whose names contain `meshery`:

    ```bash
    mesheryctl system stop --force
    ```

2. If the error persists, find and delete the leftovers yourself. Check what the `meshery` namespace contains first (`kubectl get all -n meshery`), and delete it only if it holds nothing but Meshery. `mesheryctl system stop` deletes this namespace too unless you pass `--keep-namespace`.

    ```bash
    kubectl delete namespace meshery
    kubectl get clusterroles,clusterrolebindings -o name | grep meshery
    kubectl delete <each resource listed above>
    ```

3. Run `mesheryctl system start`.

_Issue reference: [meshery/meshery#4578](https://github.com/meshery/meshery/issues/4578)_

## make server

### Database migration errors

**Error:** Meshery Server exits during startup with a database error, for example:

```text
FATA[0000] constraints not implemented on sqlite, consider using DisableForeignKeyConstraintWhenMigrating
exit status 1
make: *** [Makefile: server] Error 1
```

**Cause:** Meshery Server keeps its local SQLite database in `~/.meshery/config` (the `USER_DATA_FOLDER` environment variable overrides this location). A database created by a different Meshery version can fail to migrate.

**Fix:**

1. Delete the local database file and its write-ahead log files. This removes your locally stored Meshery data but leaves the rest of the folder in place:

    ```bash
    DB_DIR="${USER_DATA_FOLDER:-$HOME/.meshery/config}"
    rm -f "$DB_DIR/mesherydb.sql" "$DB_DIR/mesherydb.sql-wal" "$DB_DIR/mesherydb.sql-shm"
    ```

2. Run `make server` again.

## Additional resources

- [Error Code Reference]({{< ref "reference/error-codes.md" >}})
