// Package orderby implements a go/analysis pass that keeps every gorm
// ORDER BY clause fed from a value the server can prove is safe.
//
// gorm's (*gorm.DB).Order interpolates a string argument into the generated
// SQL verbatim - it is a raw SQL sink, not a bound parameter. Every one of
// Meshery's published CVEs came from that sink reached by a request-controlled
// `?order=` value. models.SanitizeOrderInput closes it by echoing back only a
// column drawn from the caller's allow-list, and roughly two dozen call sites
// route through it.
//
// Until this pass existed that convergence was a convention: a new persister
// that built an ORDER BY from a non-constant string compiled, passed CI and
// shipped. This analyzer turns the convention into a control - it fails the
// build unless the value reaching Order is a constant or comes out of the
// sanitizer.
//
// # What it accepts
//
// For each call to (*gorm.io/gorm.DB).Order whose argument is string-typed,
// the pass walks the argument backwards through SSA and requires every
// definition that can reach it to be one of:
//
//   - a constant (a string literal, or a named constant such as
//     models.defaultOrderUpdatedAtDesc);
//   - the result of models.SanitizeOrderInput.
//
// Because the walk is over SSA rather than syntax it is flow-sensitive: the
// prevailing call-site shape
//
//	order = models.SanitizeOrderInput(order, []string{"created_at", "name"})
//	if order == "" {
//	    order = defaultOrderUpdatedAtDesc
//	}
//	query = query.Order(order)
//
// reaches Order as a phi of {sanitizer result, constant} and is accepted,
// while a function that ordered on its raw parameter first and sanitized
// afterwards is not.
//
// # What it deliberately does not cover
//
// An argument the pass can see is boxed from a concrete non-string type - the
// clause.OrderByColumn and clause.OrderBy values gorm also accepts - is out of
// scope. gorm renders those through its clause builder, which quotes the
// identifier rather than interpolating it, so they are not the
// string-interpolation sink this rule guards.
//
// That skip is narrow on purpose: it holds only where the boxed type is
// statically visible and not string-underlying. gorm's Order takes `any` and
// type-switches on it at runtime, and its `case string:` branch builds a
// clause.Column{Raw: true} - the verbatim interpolation again. An argument
// whose static type is already an interface (an `any` parameter, a
// map[string]any read, an `any` struct field) is therefore reported rather than
// skipped: the pass cannot prove what such a value holds, and assuming a clause
// value would be the silent false negative that is the CVE.
//
// A type parameter is treated the same way and for the same reason: SSA builds a
// generic body once with the parameter intact, so `db.Order(order)` inside
// `func f[T ~string]` boxes a type parameter rather than a concrete type, and
// instantiating it at string reaches that same `case string:` branch.
//
// In every one of those cases the pass still accepts the value when it can
// trace the value's definitions to a constant or the sanitizer, whatever the
// static type says. Not seeing a value's type only decides the outcome when its
// origin is invisible too.
//
// clause.Expr and clause.Column{Raw: true} are escape hatches within gorm's own
// clause builder and are not checked here; neither is used in this repository.
//
// # Captured variables
//
// The walk covers SSA registers. Once a local's address is taken, or a closure
// captures it, SSA holds it in a memory cell instead and every read becomes a
// load - and proving which store reaches a given load needs a memory-flow
// analysis, including stores a closure performs at a time this pass cannot
// order. Rather than guess, the pass reports the load.
//
// That is deliberately the conservative direction: an order value that *is*
// sanitized but happens to be captured produces a false positive, which is
// loud and fixable, instead of a false negative, which is the CVE. No call site
// in this repository is affected. If you hit it, either read the sanitized
// value outside the closure and capture the result, or suppress the diagnostic
// with a justification.
//
// # Suppression
//
// A diagnostic can be suppressed with a `//nolint:orderby` comment trailing the
// offending line, or standing alone on the line above it, matching the spelling
// contributors already use for golangci-lint. Both spellings golangci-lint
// accepts work: `//nolint:orderby // reason` and `//nolint:orderby//reason`.
//
// Each one must carry an inline justification, and that is enforced rather than
// merely documented: a directive with no reason after it suppresses nothing and
// instead reports the missing justification. A bare directive reads as
// "handled" to the next reviewer, which is how a security lint gets switched
// off one call site at a time. A growing suppression list means the rule is
// mis-specified - narrow the rule instead of muting the call site.
package orderby

import (
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/ssa"
)

const (
	// analyzerName is the name the rule is reported and suppressed under.
	analyzerName = "orderby"

	// gormPkgPath and gormDBTypeName identify the raw SQL sink being guarded.
	gormPkgPath    = "gorm.io/gorm"
	gormDBTypeName = "DB"
	orderMethod    = "Order"

	// sanitizerPkgPath and sanitizerFuncName identify the only function whose
	// result is trusted to be interpolated into an ORDER BY clause.
	sanitizerPkgPath  = "github.com/meshery/meshery/server/models"
	sanitizerFuncName = "SanitizeOrderInput"

	// sanitizerFileName is the one file exempt from the rule: the sanitizer's
	// own, which is where an ORDER BY fragment is legitimately built from an
	// unsanitized string.
	sanitizerFileName = "sql-utils.go"

	// nolintPrefix introduces the directive that suppresses a single
	// diagnostic - `//nolint:orderby`. See the package doc.
	nolintPrefix = "//nolint:"
)

const diagnosticMsg = "ORDER BY built from an unsanitized value: pass the argument through " +
	"models.SanitizeOrderInput(order, []string{...}) with this query's allow-list of columns, " +
	"or use a constant. gorm's Order() interpolates a string into the SQL verbatim - this is the " +
	"sink behind every Meshery CVE. See docs/content/en/project/contributing/contributing-lint.md."

// opaqueArgMsg covers the argument whose static type is an interface. The
// remedy differs from diagnosticMsg's: the value has to be given a concrete
// type before it can be proven either way.
const opaqueArgMsg = "ORDER BY built from an unsanitized value: this argument's static type is an " +
	"interface, so the pass cannot prove it does not hold a string. gorm's Order() type-switches " +
	"at runtime and interpolates a string into the SQL verbatim - this is the sink behind every " +
	"Meshery CVE. Pass a string sanitized with models.SanitizeOrderInput(order, []string{...}), " +
	"or a concretely-typed clause value. " +
	"See docs/content/en/project/contributing/contributing-lint.md."

// unjustifiedNolintMsg replaces the diagnostic when the call site carries a
// `//nolint:orderby` with no reason after it. Reporting the missing
// justification rather than the original problem is what makes the requirement
// enforceable: a bare directive reads as "handled" to the next reviewer, so it
// has to fail loudly and say what is missing.
const unjustifiedNolintMsg = "//nolint:orderby without a justification does not suppress this rule. " +
	"Write `//nolint:orderby // <why this value cannot reach user input>`, or fix the call site by " +
	"passing the argument through models.SanitizeOrderInput(order, []string{...}). " +
	"See docs/content/en/project/contributing/contributing-lint.md."

// Analyzer reports gorm ORDER BY clauses built from values that are neither
// constant nor sanitized.
var Analyzer = &analysis.Analyzer{
	Name:     analyzerName,
	Doc:      "check that every gorm ORDER BY clause is built from a constant or from models.SanitizeOrderInput",
	URL:      "https://github.com/meshery/meshery/issues/21443",
	Requires: []*analysis.Analyzer{buildssa.Analyzer},
	Run:      run,
}

// run walks every gorm Order call in the package and reports the ones whose
// argument the pass cannot show is safe.
func run(pass *analysis.Pass) (any, error) {
	ssaResult, ok := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)
	if !ok || ssaResult == nil {
		return nil, nil
	}

	suppressed, unjustified := suppressedLines(pass)

	for _, fn := range ssaResult.SrcFuncs {
		if fn == nil {
			continue
		}
		for _, block := range fn.Blocks {
			for _, instr := range block.Instrs {
				call, ok := instr.(ssa.CallInstruction)
				if !ok {
					continue
				}
				common := call.Common()
				if !isGormOrderCall(common) {
					continue
				}
				arg, kind := classifyOrderArg(common)
				if kind == argClause {
					// A clause.OrderByColumn / clause.OrderBy value, boxed from
					// a concrete non-string type; gorm quotes those rather than
					// interpolating them.
					continue
				}
				// A value whose definitions are all constants or sanitizer
				// results is safe whatever its static type says: if the pass
				// can see where it came from, it does not need to see what it
				// is. That matters for a type parameter, where the sanitized
				// form is `T(models.SanitizeOrderInput(...))` and the static
				// type is opaque by construction. A genuinely opaque value -
				// an `any` parameter, a map[string]any read - has no traceable
				// definition and still reports.
				if isSafeOrderValue(arg, map[ssa.Value]bool{}) {
					continue
				}
				msg := diagnosticMsg
				if kind == argOpaque {
					msg = opaqueArgMsg
				}

				pos := common.Pos()
				if !pos.IsValid() {
					pos = instr.Pos()
				}
				if pos.IsValid() {
					at := pass.Fset.Position(pos)
					if isSanitizerFile(pass, at) || suppressed[sourceLine{at.Filename, at.Line}] {
						continue
					}
					if unjustified[sourceLine{at.Filename, at.Line}] {
						msg = unjustifiedNolintMsg
					}
				}
				pass.Reportf(pos, "%s", msg)
			}
		}
	}

	return nil, nil
}

// isGormOrderCall reports whether common calls (*gorm.io/gorm.DB).Order.
//
// The callee is matched on its signature - `Order(any) *gorm.DB` - rather than
// on its receiver type, so a call that reaches Order through an embedding
// wrapper is caught too. meshkit's database.Handler embeds *gorm.DB and every
// persister in this repository orders through it.
//
// A dynamic dispatch through an interface declaring that method is matched on
// the same signature: *gorm.DB satisfies such an interface, so the call lands on
// the same raw sink, and waving it through would leave "wrap the handler in an
// interface" as a way to opt out of the rule.
func isGormOrderCall(common *ssa.CallCommon) bool {
	if common == nil {
		return false
	}
	var sig *types.Signature
	if common.IsInvoke() {
		if common.Method == nil || common.Method.Name() != orderMethod {
			return false
		}
		sig, _ = common.Method.Type().(*types.Signature)
	} else {
		callee := common.StaticCallee()
		if callee == nil || callee.Name() != orderMethod {
			return false
		}
		sig = callee.Signature
	}
	if sig == nil || sig.Recv() == nil || sig.Params().Len() != 1 || sig.Results().Len() != 1 {
		return false
	}
	named := namedOf(sig.Results().At(0).Type())
	if named == nil {
		return false
	}
	obj := named.Obj()
	return obj != nil &&
		obj.Name() == gormDBTypeName &&
		obj.Pkg() != nil &&
		obj.Pkg().Path() == gormPkgPath
}

// orderArgKind says what the pass can prove about the value handed to gorm's
// `Order(value any)` parameter. Skipping a call is only sound when the dynamic
// type is provably not a string, so the unknown case is a diagnostic rather
// than a third silent exit.
type orderArgKind int

const (
	// argClause: the argument is boxed from a concrete non-string type, so
	// gorm's runtime type switch cannot reach its raw `case string:` branch.
	argClause orderArgKind = iota
	// argString: a string, the one value gorm interpolates verbatim. Its
	// definitions decide whether the call is reported.
	argString
	// argOpaque: the argument is already interface-typed, so what it holds at
	// runtime is not visible here - and a string is one of the things it can
	// hold.
	argOpaque
)

// classifyOrderArg returns the value passed to Order together with what the pass
// can prove about it, unwrapping the interface boxing gorm's `value any`
// parameter forces on a concrete argument.
func classifyOrderArg(common *ssa.CallCommon) (ssa.Value, orderArgKind) {
	args := common.Args
	if !common.IsInvoke() {
		// Args[0] is the receiver of the static method call; an invoke keeps
		// its receiver in common.Value instead.
		if len(args) == 0 {
			return nil, argClause
		}
		args = args[1:]
	}
	if len(args) != 1 || args[0] == nil {
		return nil, argClause
	}

	v := args[0]
	// A concrete argument arrives boxed, and the boxed type is the dynamic type
	// gorm switches on - visible right here. Both the boxed and the unboxed
	// form run through the same classifier: keeping two copies of it is what
	// let the boxed one fall behind and miss type parameters.
	if mi, ok := v.(*ssa.MakeInterface); ok {
		if mi.X == nil {
			return nil, argClause
		}
		v = mi.X
	}

	kind := classifyOrderArgType(v.Type())
	if kind == argClause {
		return nil, argClause
	}
	return v, kind
}

// classifyOrderArgType reports what the pass can prove about a value of type t
// reaching Order.
func classifyOrderArgType(t types.Type) orderArgKind {
	if t == nil {
		return argClause
	}
	if isStringUnderlying(t) {
		return argString
	}
	// A type parameter is not a concrete type: SSA builds a generic body once,
	// with the parameter intact, so `db.Order(order)` in `func f[T ~string]`
	// boxes a *types.TypeParam. Instantiating that at string reaches gorm's
	// `case string:` branch like any other string, so it cannot be skipped as
	// though it were a clause value. Its Underlying is the constraint
	// interface, which the check below already treats as opaque - naming the
	// case explicitly keeps it from being "simplified" away again.
	if _, ok := types.Unalias(t).(*types.TypeParam); ok {
		return argOpaque
	}
	if _, ok := t.Underlying().(*types.Interface); ok {
		// An `any` parameter, a map[string]any read, an `any` struct field, a
		// widened interface. Nothing here rules out a string.
		return argOpaque
	}
	return argClause
}

// isStringUnderlying reports whether t is a string or a named type over one.
// gorm's type switch matches the unnamed `string` exactly; keeping the named
// form in scope costs a false positive at worst, and which spelling that switch
// happens to match is an implementation detail to stay clear of.
func isStringUnderlying(t types.Type) bool {
	basic, ok := t.Underlying().(*types.Basic)
	return ok && basic.Info()&types.IsString != 0
}

// isSafeOrderValue reports whether every definition reaching v is a constant or
// the result of the sanitizer. seen breaks phi cycles in loops: a value already
// on the current path contributes no new obligation, so it is treated as
// satisfied and the remaining edges decide.
func isSafeOrderValue(v ssa.Value, seen map[ssa.Value]bool) bool {
	if v == nil {
		return false
	}
	if seen[v] {
		return true
	}
	seen[v] = true

	switch value := v.(type) {
	case *ssa.Const:
		return true
	case *ssa.Call:
		return isSanitizerCall(value.Common())
	case *ssa.Phi:
		for _, edge := range value.Edges {
			if !isSafeOrderValue(edge, seen) {
				return false
			}
		}
		return true
	case *ssa.Convert:
		return isSafeOrderValue(value.X, seen)
	case *ssa.ChangeType:
		return isSafeOrderValue(value.X, seen)
	case *ssa.MultiConvert:
		// The conversion SSA emits when either side is a type parameter whose
		// type set has more than one term - `T(s)` under `T ~string | ~[]byte`.
		// Like Convert and ChangeType it only reshapes the value, so the
		// obligation passes through to its operand.
		return isSafeOrderValue(value.X, seen)
	case *ssa.ChangeInterface:
		// Widening one interface type to another, e.g. a Stringer assigned to
		// an `any`. Also value-preserving.
		return isSafeOrderValue(value.X, seen)
	case *ssa.MakeInterface:
		// Boxing a concrete value. classifyOrderArg unwraps this at the call
		// itself; it reappears here when the boxing happened further back, as
		// in a value boxed into one interface and then widened into another.
		return isSafeOrderValue(value.X, seen)
	}

	// Parameters, field and map reads, string concatenation, fmt.Sprintf
	// results - anything else is request-reachable as far as this pass can
	// tell. So is a load from memory: SSA spills a local into a memory cell
	// once a closure captures it, and proving which store reaches a given load
	// needs a memory-flow analysis this pass deliberately does not carry. See
	// "Captured variables" in the package doc.
	return false
}

// isSanitizerCall reports whether common is a static call to
// models.SanitizeOrderInput.
func isSanitizerCall(common *ssa.CallCommon) bool {
	if common == nil {
		return false
	}
	callee := common.StaticCallee()
	if callee == nil || callee.Name() != sanitizerFuncName || callee.Pkg == nil {
		return false
	}
	return callee.Pkg.Pkg != nil && callee.Pkg.Pkg.Path() == sanitizerPkgPath
}

// sourceLine identifies one line of one file. Suppression is tracked per file:
// a package compiles several, and a directive in one of them must not silence
// the same line number in another.
type sourceLine struct {
	filename string
	line     int
}

// isSanitizerFile reports whether at is inside models/sql-utils.go, the one
// file allowed to assemble an ORDER BY fragment from unsanitized input.
func isSanitizerFile(pass *analysis.Pass, at token.Position) bool {
	if pass.Pkg == nil || pass.Pkg.Path() != sanitizerPkgPath {
		return false
	}
	return filepath.Base(at.Filename) == sanitizerFileName
}

// suppressedLines collects the 1-based lines a `//nolint:orderby` directive
// covers: the comment's own line always, and the line below it only when the
// comment stands alone on its line. That second restriction matters - without
// it a trailing directive would silently cover the statement underneath it,
// which is a false negative in a rule whose whole job is to have none.
func suppressedLines(pass *analysis.Pass) (suppressed, unjustified map[sourceLine]bool) {
	suppressed = map[sourceLine]bool{}
	unjustified = map[sourceLine]bool{}
	sources := map[string][]byte{}
	for _, file := range pass.Files {
		for _, group := range file.Comments {
			for _, comment := range group.List {
				names, justified := hasNolintDirective(comment)
				if !names {
					continue
				}
				// An unjustified directive does not suppress. The docs have
				// always required a reason; letting the bare form work anyway
				// is how a security lint gets quietly switched off one call
				// site at a time, so the requirement is enforced here rather
				// than merely documented.
				target := suppressed
				if !justified {
					target = unjustified
				}
				pos := pass.Fset.Position(comment.Slash)
				target[sourceLine{pos.Filename, pos.Line}] = true
				if startsItsLine(pass, sources, pos) {
					target[sourceLine{pos.Filename, pos.Line + 1}] = true
				}
			}
		}
	}
	return suppressed, unjustified
}

// startsItsLine reports whether only whitespace precedes pos on its line. It
// reads each file at most once, and answers conservatively - a source it cannot
// read yields no extra suppression rather than an unchecked line.
func startsItsLine(pass *analysis.Pass, sources map[string][]byte, pos token.Position) bool {
	if pos.Column <= 1 {
		return true
	}
	src, cached := sources[pos.Filename]
	if !cached {
		if pass.ReadFile != nil {
			src, _ = pass.ReadFile(pos.Filename)
		}
		sources[pos.Filename] = src
	}
	// token.Position.Column counts bytes from the start of the line, so the
	// line begins Column-1 bytes before the comment.
	start := pos.Offset - (pos.Column - 1)
	if src == nil || start < 0 || pos.Offset > len(src) {
		return false
	}
	return strings.TrimSpace(string(src[start:pos.Offset])) == ""
}

// hasNolintDirective reports whether comment names this analyzer in a nolint
// directive, in any position of the comma-separated linter list. Bare
// `//nolint` is deliberately not honoured, and neither is the space-separated
// `// nolint` spelling golangci-lint also rejects: silencing this rule must be
// explicit about what is being silenced.
func hasNolintDirective(comment *ast.Comment) (names, justified bool) {
	text := strings.TrimSpace(comment.Text)
	if !strings.HasPrefix(text, nolintPrefix) {
		return false, false
	}

	// Split the linter list from the explanation. golangci-lint accepts both
	// `//nolint:orderby // reason` and `//nolint:orderby//reason`, so the
	// separator is the next `//` rather than a space-prefixed one - a linter
	// name can never contain a slash.
	rest := text[len(nolintPrefix):]
	list, explanation := rest, ""
	if i := strings.Index(rest, "//"); i >= 0 {
		list, explanation = rest[:i], strings.TrimSpace(rest[i+2:])
	}

	// Match the linter list exactly so `//nolint:gocritic,orderby` counts and
	// `//nolint:orderbyish` does not.
	for name := range strings.SplitSeq(list, ",") {
		if strings.TrimSpace(name) == analyzerName {
			return true, explanation != ""
		}
	}
	return false, false
}

// namedOf unwraps a pointer to the named type underneath, seeing through any
// alias on the way.
func namedOf(t types.Type) *types.Named {
	if ptr, ok := types.Unalias(t).(*types.Pointer); ok {
		t = ptr.Elem()
	}
	named, _ := types.Unalias(t).(*types.Named)
	return named
}
