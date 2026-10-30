# Troubleshooting Errors while running Meshery

> Troubleshoot errors from mesheryctl system start and from running Meshery Server locally with make server.

---

LLMS index: [llms.txt](/pr-preview/pr-22675/llms.txt)

---

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

- [Error Code Reference](/pr-preview/pr-22675/reference/error-codes/)
