**Notes for Reviewers**

- This PR fixes `mesheryctl filter view` accepting unlimited positional arguments without validation.
- Scope is intentionally limited to `mesheryctl/internal/cli/root/filter/view.go`.

**What changed**

- Replaced `cobra.ArbitraryArgs` with `cobra.MaximumNArgs(1)` to reject extra positional arguments.
- Moved `--all` flag conflict check and missing name/ID check from `RunE` to `PreRunE` for early validation.
- Removed now-unused `parseQuotedArg` function (no longer needed since `MaximumNArgs(1)` limits args to one).
- Removed unused `subCmdUsed` variable from `RunE`.

**Before (broken)**

```bash
# Silently accepted — no error
mesheryctl filter view arg1 arg2 arg3
```

**After (fixed)**

```bash
# Properly rejected
mesheryctl filter view arg1 arg2
Error: accepts at most 1 arg(s), received 2
```

**Validation**

- Ran:
  ```
  go build ./mesheryctl/internal/cli/root/filter/...
  go test -v ./mesheryctl/internal/cli/root/filter/...
  ```
- All 7 existing tests pass:
  - `TestDeleteCmd/Delete_Kuma-Test` ✅
  - `TestDeleteCmd/Delete_RolloutAndIstio` ✅
  - `TestFilterCmd/Filter_command_with_no_args` ✅
  - `TestListCmd/Fetch_Filter_List` ✅
  - `TestViewCmd/Fetch_Filter_View` ✅
  - `TestViewCmd/Fetch_Kuma_Filter_View_with_ID` ✅
  - `TestViewCmd/Fetch_Filter_View_for_non_existing_filter` ✅

**[Signed commits](https://github.com/meshery/meshery/blob/master/CONTRIBUTING.md#signing-off-on-commits-developer-certificate-of-origin)**
- [x] Yes, I signed my commits.
