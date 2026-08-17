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
// Non-string arguments - the clause.OrderByColumn and clause.OrderBy values
// gorm also accepts - are out of scope. gorm renders those through its clause
// builder, which quotes the identifier rather than interpolating it, so they
// are not the string-interpolation sink this rule guards. clause.Expr and
// clause.Column{Raw: true} are escape hatches within that builder and are not
// checked here; neither is used in this repository.
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
// contributors already use for golangci-lint. Each one must carry an inline
// justification. A growing suppression list means the rule is mis-specified -
// narrow the rule instead of muting the call site.
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

// Analyzer reports gorm ORDER BY clauses built from values that are neither
// constant nor sanitized.
var Analyzer = &analysis.Analyzer{
	Name:     analyzerName,
	Doc:      "check that every gorm ORDER BY clause is built from a constant or from models.SanitizeOrderInput",
	URL:      "https://github.com/meshery/meshery/issues/21443",
	Requires: []*analysis.Analyzer{buildssa.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (any, error) {
	ssaResult, ok := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)
	if !ok || ssaResult == nil {
		return nil, nil
	}

	suppressed := suppressedLines(pass)

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
				arg := stringOrderArg(common)
				if arg == nil {
					// A clause.OrderByColumn / clause.OrderBy value; gorm
					// quotes those rather than interpolating them.
					continue
				}
				if isSafeOrderValue(arg, map[ssa.Value]bool{}) {
					continue
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
				}
				pass.Reportf(pos, "%s", diagnosticMsg)
			}
		}
	}

	return nil, nil
}

// isGormOrderCall reports whether common is a static call to
// (*gorm.io/gorm.DB).Order. gorm's *DB is a concrete type, so an interface
// dispatch can never land there and is rejected outright.
//
// The callee is matched on its signature - `Order(any) *gorm.DB` - rather than
// on its receiver type, so a call that reaches Order through an embedding
// wrapper is caught too. meshkit's database.Handler embeds *gorm.DB and every
// persister in this repository orders through it.
func isGormOrderCall(common *ssa.CallCommon) bool {
	if common == nil || common.IsInvoke() {
		return false
	}
	callee := common.StaticCallee()
	if callee == nil || callee.Name() != orderMethod {
		return false
	}
	sig := callee.Signature
	if sig.Recv() == nil || sig.Params().Len() != 1 || sig.Results().Len() != 1 {
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

// stringOrderArg returns the string-typed value passed to Order, unwrapping the
// interface boxing gorm's `value any` parameter forces. It returns nil when the
// argument is not string-typed.
func stringOrderArg(common *ssa.CallCommon) ssa.Value {
	// Args[0] is the receiver of the static method call.
	if len(common.Args) != 2 {
		return nil
	}
	v := common.Args[1]
	if mi, ok := v.(*ssa.MakeInterface); ok {
		v = mi.X
	}
	if v == nil {
		return nil
	}
	basic, ok := v.Type().Underlying().(*types.Basic)
	if !ok || basic.Info()&types.IsString == 0 {
		return nil
	}
	return v
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
func suppressedLines(pass *analysis.Pass) map[sourceLine]bool {
	lines := map[sourceLine]bool{}
	sources := map[string][]byte{}
	for _, file := range pass.Files {
		for _, group := range file.Comments {
			for _, comment := range group.List {
				if !hasNolintDirective(comment) {
					continue
				}
				pos := pass.Fset.Position(comment.Slash)
				lines[sourceLine{pos.Filename, pos.Line}] = true
				if startsItsLine(pass, sources, pos) {
					lines[sourceLine{pos.Filename, pos.Line + 1}] = true
				}
			}
		}
	}
	return lines
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
func hasNolintDirective(comment *ast.Comment) bool {
	text := strings.TrimSpace(comment.Text)
	if !strings.HasPrefix(text, nolintPrefix) {
		return false
	}
	// Drop the trailing ` // explanation`, then match the linter list exactly so
	// `//nolint:gocritic,orderby` counts and `//nolint:orderbyish` does not.
	names, _, _ := strings.Cut(text[len(nolintPrefix):], " //")
	for name := range strings.SplitSeq(names, ",") {
		if strings.TrimSpace(name) == analyzerName {
			return true
		}
	}
	return false
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
