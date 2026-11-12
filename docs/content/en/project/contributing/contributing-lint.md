---
title: Go Lint Rules
description: The repo-specific Go lint gates Meshery enforces in CI, what each one is protecting, and what to do when one fires.
categories: [contributing]
---

Beyond the stock linters, Meshery enforces a small number of repo-specific Go
rules. Each one exists because a convention was being relied on to hold a
security or contract boundary, and a convention is not a control: it holds until
the first contributor who has not read the doc adds a call site.

All of them run in the `golangci-lint-server` job of
[`.github/workflows/go-testing-ci.yml`](https://github.com/meshery/meshery/blob/master/.github/workflows/go-testing-ci.yml),
which gates the unit-test and build jobs downstream, and all of them run locally
from:

```bash
make golangci
```

`mesheryctl` is part of the same root Go module and is covered by the same rules.
Its own target runs the identical pair scoped to the CLI:

```bash
cd mesheryctl && make lint
```

## ORDER BY must be sanitized

`gorm`'s `(*gorm.DB).Order` interpolates a string argument into the generated SQL
verbatim - it is a raw SQL sink, not a bound parameter. Every one of Meshery's
published CVEs came from that sink reached by a request-controlled `?order=`
value.

`models.SanitizeOrderInput` in
[`server/models/sql-utils.go`](https://github.com/meshery/meshery/blob/master/server/models/sql-utils.go)
closes it. It takes the requested order and the caller's allow-list of database
columns, and returns either a column drawn from that allow-list or the empty
string. The returned value is never echoed back from user input, which is what
makes it safe to interpolate.

Roughly two dozen call sites - the `server/models/*_persister.go` family,
`default_local_provider.go`, `database_handlers.go`, `meshsync_handler.go` -
route through it. The `orderby` analyzer
([`server/internal/lint/orderby`](https://github.com/meshery/meshery/tree/master/server/internal/lint/orderby))
is what keeps the next one from skipping it.

### What the rule accepts

For every call to `(*gorm.DB).Order` whose argument is string-typed, the
analyzer walks the argument backwards through SSA and requires **every**
definition that can reach it to be either a constant or the result of
`models.SanitizeOrderInput`. Because the walk is over SSA rather than syntax it
is flow-sensitive, so the prevailing call-site shape is accepted as-is:

```go
order = models.SanitizeOrderInput(order, []string{"created_at", "updated_at", "name"})
if order == "" {
    order = defaultOrderUpdatedAtDesc
}
query = query.Order(order)
```

The value reaching `Order` is a phi of `{sanitizer result, constant}`, and both
edges are safe. A function that ordered on its raw parameter and sanitized
afterwards is *not* accepted - which is the point, and is why the rule is not a
grep for `SanitizeOrderInput` somewhere in the function.

Calls that reach `Order` through an embedding wrapper are covered too: meshkit's
`database.Handler` embeds `*gorm.DB`, and every persister orders through it. So
are calls dispatched through an interface that declares
`Order(any) *gorm.DB` - `*gorm.DB` satisfies such an interface, so the call
still lands on the same raw sink.

### What to do when it fires

The diagnostic opens with:

```text
ORDER BY built from an unsanitized value: pass the argument through
models.SanitizeOrderInput(order, []string{...}) with this query's allow-list of
columns, or use a constant.
```

An argument whose static type is an interface produces a variant of that
message - "this argument's static type is an interface, so the pass cannot
prove it does not hold a string" - which asks for a concretely-typed value
rather than a sanitizer call. See *What the rule does not cover* below for that
case.

There are three correct fixes, in order of preference:

1. **Sanitize with this query's own allow-list.** Add
   `order = models.SanitizeOrderInput(order, []string{...})` before the `Order`
   call, listing the **snake_case database columns** this query may legitimately
   sort on. Do not copy another persister's allow-list without checking it
   against your table - the allow-list is the security boundary, and a column
   that does not exist on your table produces an empty order rather than an
   error.
2. **Use a constant** if the query's order is fixed. `defaultOrderUpdatedAtDesc`
   in `server/models/vars.go` already covers `updated_at desc`.
3. **Use a `clause` value** - `clause.OrderByColumn{Column: clause.Column{Name: col}}` -
   if you need gorm to build the clause. gorm renders those through its quoting
   clause builder rather than interpolating them. Keep the value concretely
   typed: an `any` holding a clause value is reported, for the reason below.

What is *not* a fix: sanitizing after the `Order` call, moving the
interpolation into a `fmt.Sprintf`, or silencing the rule.

### The one false positive: captured variables

If a closure captures the order variable, the analyzer reports the call even
when the value is sanitized:

```go
order = models.SanitizeOrderInput(order, validColumns)
return func() *gorm.DB {
    return db.Order(order) // reported, despite being safe
}
```

Capturing a local moves it out of an SSA register and into a memory cell, and
proving which store reaches a given load - including stores a closure performs
at a point the analyzer cannot order - needs a memory-flow analysis the rule
deliberately does not carry. Guessing in the other direction would let a closure
overwrite a sanitized value unnoticed, which is exactly the CVE.

The fix is to read the sanitized value out before anything captures it:

```go
sanitized := models.SanitizeOrderInput(order, validColumns)
captured := sanitized
return db.Order(sanitized), func() string { return captured }
```

Taking the variable's address (`overwrite(&order)`) has the same effect, for the
same reason. No call site in the repository is affected today.

### What the rule does not cover

An argument boxed from a **concrete** non-string type is out of scope, for the
reason in fix 3 above.

That exemption stops where the analyzer stops being able to see the type.
`Order` takes `any` and type-switches on it at runtime, and its `case string:`
branch builds a `clause.Column{Raw: true}` - the verbatim interpolation again.
So an argument whose static type is already an interface is **reported**, not
skipped:

```go
func applyOrder(db *gorm.DB, order any) *gorm.DB {
    return db.Order(order) // reported: this may be a string
}

query.Order(filters["order"]) // reported: map[string]any read
```

The pass cannot prove such a value is not a string, and assuming it is a clause
value would be a silent false negative - the CVE. Give the value a concrete
type: a `string` you have sanitized, or a `clause.OrderByColumn`.

A **type parameter** is treated the same way, for the same reason:

```go
func orderBy[T ~string](db *gorm.DB, order T) *gorm.DB {
    return db.Order(order) // reported: instantiated at string, this is the sink
}
```

SSA builds a generic body once with the parameter intact, so what reaches `Order`
is a type parameter rather than a concrete type. Instantiated at `string` it
lands on gorm's `case string:` branch like any other string.

In all of these cases the pass still accepts the value if it can trace where it
came from - `db.Order(T(models.SanitizeOrderInput(order, cols)))` passes, because
every definition reaching it is a sanitizer result. Being unable to see a value's
*type* is only fatal when its *origin* is also invisible.

That leaves two escape hatches inside gorm's own clause builder - `clause.Expr`,
and `clause.Column{Raw: true}` - which are raw SQL again. Neither is used in this
repository; if you reach for one, you are back to owning the sanitization
yourself.

The rule also covers only `Order`. Other gorm methods that take raw SQL
(`Raw`, `Select` with an expression, `Joins`) are not analyzed.

### Suppression

A diagnostic can be suppressed with `//nolint:orderby` trailing the offending
line, or standing alone on the line above it. Every use must carry an inline
justification, and that requirement is **enforced, not just documented**:

```go
//nolint:orderby // <why this value cannot reach user input>
```

A directive with no reason after it does not suppress anything. Instead the call
site reports a different diagnostic naming the missing justification, because a
bare directive reads as "handled" to the next reviewer - which is how a security
lint gets switched off one call site at a time. Both spellings golangci-lint
accepts work: `//nolint:orderby // reason` and `//nolint:orderby//reason`.

The name must appear in the directive's own list - a bare `//nolint` does not
silence this rule, and a *trailing* directive covers only its own line, never
the statement below it.

Suppressions are expected to stay in the low single digits. A growing list means
the rule is mis-specified - narrow the rule rather than muting the call sites.

### Running and testing it

```bash
# the whole root module (server and mesheryctl)
go run ./server/internal/lint/orderby/cmd/orderbylint ./...

# the analyzer's own tests
go test ./server/internal/lint/orderby/...
```

`go run` recompiles the analyzer on every invocation, which is fine in CI (the
build cache is warm) but noticeable locally after a clean cache. For iterative
local use, build it once:

```bash
go build -o /tmp/orderbylint ./server/internal/lint/orderby/cmd/orderbylint
/tmp/orderbylint ./...
```

The analyzer is pinned by an `analysistest` fixture at
`server/internal/lint/orderby/testdata/`, which asserts both directions: every
rejected shape reports, and every accepted shape stays silent. The second half
matters as much as the first - a rule that fired on the existing call sites
would be turned off rather than obeyed. Add a case there for any shape you teach
the rule about.

### Why it is not a golangci-lint linter

`forbidigo` and `gocritic` match on the *called function*, not on its arguments,
so neither can express "non-constant argument". Hosting a custom analyzer inside
golangci-lint requires building a bespoke `golangci-lint` binary through its
module plugin system, which would mean dropping
`golangci/golangci-lint-action` - and with it `only-new-issues` - from two jobs.
A single `go run` step in the same lint job is less machinery for the same gate.
The trade-off is that golangci-lint's own `//nolint` processing does not apply,
so the analyzer implements the directive itself.

## No `http.Error` in server handlers

Enforced by a `forbidigo` pattern in
[`.github/.golangci.yml`](https://github.com/meshery/meshery/blob/master/.github/.golangci.yml).
`http.Error` writes `Content-Type: text/plain`, which breaks RTK Query's default
`baseQuery`. Use `writeMeshkitError(w, err, status)` or
`writeJSONError(w, msg, status)` from `server/handlers/utils.go` instead.

The rule, its file-level allowlist, and the legitimate non-JSON responders are
documented with their subject in
[HTTP Error Response Contract]({{< ref "project/contributing/error-contract.md" >}}).

## Adding a rule

Prefer configuration. If an existing linter expresses the rule cleanly, add it to
`.github/.golangci.yml` with a `msg` that names the correct alternative and links
the doc explaining why - contributors meet the rule through that message, not
through this page.

Reach for a `go/analysis` analyzer only when configuration cannot express the
rule, as with the argument-shape requirement above. When you do, it needs: the
analyzer package, a `cmd/` wrapper using `singlechecker`, an `analysistest`
fixture covering accepted *and* rejected shapes, a step in the
`golangci-lint-server` job, a line in `make golangci`, and a section here.
