I'm lucky singh (Mann) i'm a first year  Computer Science student with a deep passion for **cloud-native systems**, **distributed architecture**, and **backend engineering**. Over the past half years, I have been systematically studying Go, Kubernetes internals, and the design patterns that power production-grade infrastructure platforms.

My approach to open source is deliberate and methodical:

- I read the entire relevant codebase before opening a single issue
- I reproduce bugs locally and document them with minimal reproducible examples
- I submit small, focused PRs that are easy to review, building trust incrementally
- I respond to reviewer feedback — both from human maintainers and automated tools like Gemini AI review — quickly and professionally
- I never claim ownership of a task until I have demonstrated familiarity with the code it touches


I discovered Meshery while exploring the CNCF landscape and was immediately struck by the quality of its engineering and the warmth of its community. The codebase presented real challenges — not tutorial-level problems — and that drew me in.

what i do till now i make some pr 
## Pull Request #17663 — Fix Variable Shadowing in `VerifyAndConvertToDesign`

**Link:** https://github.com/meshery/meshery/pull/17663

**Branch:** `fix/variable-shadowing-VerifyAndConvertToDesign`

## Issue Filed — Nil Pointer Dereference in `handleFilterPOST`

**File:** `server/handlers/meshery_filter_handler.go`

**Function:** `handleFilterPOST` (~lines 106–121)

**Bug:** `parsedBody` is declared as a pointer and initialized to `nil`. When `json.Decode` fails, the error handler immediately accesses `parsedBody.FilterData.Name` — dereferencing a nil pointer — causing a **guaranteed server panic / goroutine crash** on any malformed filter POST request.

## Pull Request #17751 — Fix Missing `return` After Error in Catalog Publish/Unpublish Handlers

**Link:** https://github.com/meshery/meshery/pull/17751

**Closes:** Issue #17746

**Problem Solved:** In 4 HTTP handler functions across `meshery_pattern_handler.go` and `meshery_filter_handler.go`, after a `json.Unmarshal` failure, the code called `http.Error(...)` but was **missing a `return` statement** — causing execution to continue into code that accesses `respBody.ContentName` on a nil pointer, resulting in a **guaranteed server panic**.
