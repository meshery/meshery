# Cleanup: When meshkit#1007 merges

Once `github.com/meshery/meshkit/retry` is available on the upstream module:

```bash
rm -rf server/models/retryutil/
```

Then in `server/models/remote_provider.go`, flip the import:

```diff
-	retry "github.com/meshery/meshery/server/models/retryutil"
+	"github.com/meshery/meshkit/retry"
```

In `go.mod`, either uncomment and adjust the `replace` line, or just
bump `github.com/meshery/meshkit` in `require` to the version that
includes the `retry` package — whichever is appropriate at the time.

```diff
-// TODO: replace github.com/meshery/meshkit => ../meshkit
-// Once meshery/meshkit#1007 merges, switch the import in
-// server/models/remote_provider.go back to github.com/meshery/meshkit/retry
-// and delete server/models/retryutil/.
-//
-// Until then, server/models/retryutil/ provides the same retry API by
-// wrapping github.com/cenkalti/backoff/v5 directly.
+replace github.com/meshery/meshkit => ../meshkit
```
