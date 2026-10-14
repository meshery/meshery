# Pre-Landing Review Checklist (Meshery)

## Pass 1 - CRITICAL (blocking)

### SQL & data safety

- [ ] No dynamic string reaches gorm `.Order(...)` unless it went through
      `models.SanitizeOrderInput(order, []string{...})` **before** the call,
      with that query's own snake_case allow-list (AGENTS.md, Go; the
      `orderby` analyzer)
- [ ] No other raw SQL built by string concatenation - `Where`, `Raw`, `Exec`,
      `Joins`, `Select` use `?` placeholders
- [ ] Bulk updates/deletes have a verified `WHERE` clause (gorm refuses a bare
      `Delete` only when `AllowGlobalUpdate` is off - check it is)
- [ ] gorm errors are propagated, not discarded - a schemas rename silently
      renames the column (AGENTS.md, Testing)
- [ ] Raw column references (`Select`, `Where`, `Order`, `Joins`, `Scan`) match
      the column gorm derives from the Go field name, not just the `db:` tag
- [ ] New system tables are added to `SystemDatabaseModels()` in
      `server/models/system_migration.go`, not migrated at a call site
- [ ] Credential `secret` maps are read only through
      `server/models/credential_secret.go` / `ui/utils/credentialSecret.ts`

### Concurrency

- [ ] No read-then-write without a lock or transaction
- [ ] Uniqueness enforced by a DB constraint, not only application code
- [ ] Shared state guarded (`go test -race` on the touched package where practical)
- [ ] A boot-time goroutine spawned inside `models.RunSeedStage` defers its own
      `recover` at the spawn site
- [ ] Event handlers / NATS subscribers are idempotent (safe to redeliver)

### Trust boundaries

- [ ] Boot-time code never resolves a provider via
      `HandlerConfig.Providers[name]`; request-path lookups use comma-ok;
      events outside a request go through `HandlerConfig.SystemEventPersister`
- [ ] Every new content-bearing UI page has an access gate:
      `useHasPermission(Keys.X)`, `<DefaultError permissionKey={Keys.X} />`
      (inside `MesheryPage` when in `ui/pages/`), and `skip: !canX` on its RTK
      hooks - with a test pinning the **deny** path
- [ ] New server endpoints check authentication and authorization
- [ ] User input validated before privileged operations (Kubernetes apply,
      shell/exec, file paths)
- [ ] Remote-provider and external API responses decoded into the schemas type,
      not a local redeclaration
- [ ] No secrets hardcoded or logged

### Meshery contracts

- [ ] New or changed HTTP API is defined in `meshery/schemas` and consumed via
      the generated client - no hand-rolled `builder.query`/`builder.mutation`,
      response types or Go structs duplicating a schemas type
- [ ] Wire casing: JSON tags and URL params camelCase (`Id`, never `ID`), DB
      snake_case, `json:` tag never equals `db:` tag
- [ ] No new import of `v1beta1` schemas packages
- [ ] Cache tags attached with the callback form of `enhanceEndpoints`
- [ ] New MeshKit error codes: `make error` run, `next_error_code` bumped in the
      same commit, `docs/data/errorref/` regenerated
- [ ] No hand-edited lock file; `ui/tsconfig.tsbuildinfo` not staged
- [ ] No test file under `ui/pages/`
- [ ] Commits signed off (`git log --format=%b | grep -c Signed-off-by`)

## Pass 2 - INFORMATIONAL (non-blocking)

### Hidden side effects

- [ ] Side effects not buried in conditional branches or defers
- [ ] Errors not silently swallowed; MeshKit errors logged with
      `utils.Log.Error(err)` so code/cause/remediation render

### Magic numbers

- [ ] Literals, thresholds and timeouts named and justified

### Dead code

- [ ] No unreachable branches, unused imports, or unexplained commented-out code

### Test gaps

- [ ] New paths covered; modified behavior has updated tests
- [ ] Error paths and deny paths tested
- [ ] Golden files regenerated only where they still encode intended behavior
- [ ] mesheryctl tests do not point viper at `mesheryctl/pkg/utils/TestConfig.yaml`

### Crypto & entropy

- [ ] `crypto/rand` for anything security-relevant; TLS for sensitive data

### Time handling

- [ ] Timezone-aware comparisons; `time.Time` stored as UTC

### Type coercion

- [ ] No lossy numeric conversions across Go/JSON/TS boundaries (int64 in JS)

## Gate classification

| Category | Severity | Blocking |
|---|---|---|
| SQL & data safety | CRITICAL | Yes |
| Concurrency | CRITICAL | Yes |
| Trust boundaries | CRITICAL | Yes |
| Meshery contracts | CRITICAL | Yes |
| All Pass 2 categories | INFORMATIONAL | No |
