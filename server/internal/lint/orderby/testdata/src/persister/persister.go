// Package persister exercises the orderby analyzer against the ORDER BY shapes
// that appear in - or could plausibly be added to - server/models and
// server/handlers.
package persister

import (
	"fmt"

	"github.com/meshery/meshery/server/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const defaultOrder = "updated_at desc"

var validColumns = []string{"created_at", "updated_at", "name"}

//
// Accepted: the value reaching Order is provably constant or sanitized.
//

func literal(db *gorm.DB) *gorm.DB {
	return db.Order("updated_at desc")
}

func namedConstant(db *gorm.DB) *gorm.DB {
	return db.Order(defaultOrder)
}

func sanitizedDirectly(db *gorm.DB, order string) *gorm.DB {
	return db.Order(models.SanitizeOrderInput(order, validColumns))
}

// sanitizedThenDefaulted is the shape every persister in this repository uses:
// the sanitizer empties an unknown column, and a constant fills the gap.
func sanitizedThenDefaulted(db *gorm.DB, order string) *gorm.DB {
	order = models.SanitizeOrderInput(order, validColumns)
	if order == "" {
		order = defaultOrder
	}
	return db.Order(order)
}

// sanitizedComposition mirrors default_local_provider.go's events query, where
// two request fields are joined before sanitization rather than after.
func sanitizedComposition(db *gorm.DB, sortOn, direction string) *gorm.DB {
	order := models.SanitizeOrderInput(fmt.Sprintf("%s %s", sortOn, direction), validColumns)
	return db.Order(order)
}

// sanitizedInLoop keeps a phi cycle in play: every edge reaching Order is still
// either the constant seed or a sanitizer result.
func sanitizedInLoop(db *gorm.DB, candidates []string) *gorm.DB {
	order := defaultOrder
	for _, candidate := range candidates {
		order = models.SanitizeOrderInput(candidate, validColumns)
	}
	return db.Order(order)
}

// clauseValue is out of scope: gorm renders clause values through its quoting
// clause builder rather than interpolating them.
func clauseValue(db *gorm.DB, order string) *gorm.DB {
	order = models.SanitizeOrderInput(order, validColumns)
	return db.Order(clause.OrderByColumn{Column: clause.Column{Name: order}, Desc: true})
}

// clauseVariable pins the same skip when the clause value reaches Order through
// a variable: what takes it out of scope is the boxed type being concrete and
// not a string, not the argument being a composite literal in the call.
func clauseVariable(db *gorm.DB, order string) *gorm.DB {
	order = models.SanitizeOrderInput(order, validColumns)
	column := clause.OrderByColumn{Column: clause.Column{Name: order}, Desc: true}
	return db.Order(column)
}

// boxedSanitized is the other half of that contract: the boxed type is visible
// and it *is* a string, so the value's definitions still have to be safe.
func boxedSanitized(db *gorm.DB, order string) *gorm.DB {
	var value any = models.SanitizeOrderInput(order, validColumns)
	return db.Order(value)
}

// orderer is the shape a persister would reach Order through if it abstracted
// the handler behind an interface. *gorm.DB satisfies it, so the call lands on
// the same sink and the argument is checked the same way.
type orderer interface {
	Order(value any) *gorm.DB
}

func throughInterfaceSanitized(o orderer, order string) *gorm.DB {
	return o.Order(models.SanitizeOrderInput(order, validColumns))
}

// promotedThroughEmbedding mirrors meshkit's database.Handler, which embeds
// *gorm.DB so persisters call Order on the handler.
type handler struct {
	*gorm.DB
}

func promotedSanitized(h *handler, order string) *gorm.DB {
	order = models.SanitizeOrderInput(order, validColumns)
	return h.Order(order)
}

func suppressed(db *gorm.DB, order string) *gorm.DB {
	//nolint:orderby // fixture: proves the directive is honoured.
	return db.Order(order)
}

func suppressedTrailing(db *gorm.DB, order string) *gorm.DB {
	return db.Order(order) //nolint:orderby // fixture: trailing form.
}

func suppressedAlongsideAnotherLinter(db *gorm.DB, order string) *gorm.DB {
	return db.Order(order) //nolint:gocritic,orderby // fixture: not the first name in the list.
}

// suppressedNoSpaceBeforeReason pins the spelling golangci-lint also tolerates:
// the explanation separator is the next `//`, not a space-prefixed one.
func suppressedNoSpaceBeforeReason(db *gorm.DB, order string) *gorm.DB {
	return db.Order(order) //nolint:orderby//fixture: no space before the reason.
}

//
// Rejected: the directive names this analyzer but carries no justification.
// The docs have always required a reason; a bare directive reads as "handled"
// to the next reviewer, so it has to fail loudly rather than silently work.
//

// These use the standalone directive form because a trailing `//nolint:orderby`
// cannot be tested here: analysistest keys its expectation annotation to the
// diagnostic's own line, and that annotation would itself be read as the
// justification. The analyzer takes the same path either way - what it reads is
// the text after the linter list, wherever the comment sits.

func unjustifiedStandalone(db *gorm.DB, order string) *gorm.DB {
	//nolint:orderby
	return db.Order(order) // want "without a justification does not suppress"
}

func unjustifiedWithOtherLinter(db *gorm.DB, order string) *gorm.DB {
	//nolint:gocritic,orderby
	return db.Order(order) // want "without a justification does not suppress"
}

//
// Rejected: the directive does not actually name this analyzer.
//

func notSuppressedByOtherLinter(db *gorm.DB, order string) *gorm.DB {
	return db.Order(order) //nolint:gocritic // want "ORDER BY built from an unsanitized value"
}

func notSuppressedByBareNolint(db *gorm.DB, order string) *gorm.DB {
	//nolint
	return db.Order(order) // want "ORDER BY built from an unsanitized value"
}

func notSuppressedByPrefixMatch(db *gorm.DB, order string) *gorm.DB {
	return db.Order(order) //nolint:orderbyish // want "ORDER BY built from an unsanitized value"
}

// notSuppressedByTrailingDirectiveAbove pins the line+1 rule to standalone
// directives: a trailing one excuses its own call, never the next.
func notSuppressedByTrailingDirectiveAbove(db *gorm.DB, first, second string) *gorm.DB {
	db.Order(first)         //nolint:orderby // fixture: excuses this line only.
	return db.Order(second) // want "ORDER BY built from an unsanitized value"
}

//
// Rejected: the value reaching Order is request-reachable.
//

func rawParameter(db *gorm.DB, order string) *gorm.DB {
	return db.Order(order) // want "ORDER BY built from an unsanitized value"
}

func concatenated(db *gorm.DB, column string) *gorm.DB {
	return db.Order(column + " desc") // want "ORDER BY built from an unsanitized value"
}

func formatted(db *gorm.DB, column, direction string) *gorm.DB {
	return db.Order(fmt.Sprintf("%s %s", column, direction)) // want "ORDER BY built from an unsanitized value"
}

// orderedBeforeSanitizing is the case a syntactic "does this function call the
// sanitizer?" check would wave through.
func orderedBeforeSanitizing(db *gorm.DB, order string) *gorm.DB {
	query := db.Order(order) // want "ORDER BY built from an unsanitized value"
	order = models.SanitizeOrderInput(order, validColumns)
	_ = order
	return query
}

// partiallySanitized falls back to a caller-supplied string rather than a
// constant, so one edge into Order is still request-reachable.
func partiallySanitized(db *gorm.DB, order, fallback string) *gorm.DB {
	order = models.SanitizeOrderInput(order, validColumns)
	if order == "" {
		order = fallback
	}
	return db.Order(order) // want "ORDER BY built from an unsanitized value"
}

type filter struct {
	SortOn string
}

func structField(db *gorm.DB, f filter) *gorm.DB {
	return db.Order(f.SortOn) // want "ORDER BY built from an unsanitized value"
}

func mapValue(db *gorm.DB, params map[string]string) *gorm.DB {
	return db.Order(params["order"]) // want "ORDER BY built from an unsanitized value"
}

//
// Rejected: the argument's static type is an interface, so the pass cannot see
// what it holds. gorm's Order type-switches at runtime and its `case string:`
// branch builds a clause.Column{Raw: true} - the verbatim interpolation - so
// assuming these are clause values would be a silent false negative.
//

// anyMapValue is the reachable shape: the `filter` query parameter decodes into
// a map[string]any, and the order is read straight back out of it.
func anyMapValue(db *gorm.DB, filters map[string]any) *gorm.DB {
	return db.Order(filters["order"]) // want "ORDER BY built from an unsanitized value"
}

// anyParameter is the shared-helper shape: one function forwarding both strings
// and clause values keeps the string case invisible to every caller.
func anyParameter(db *gorm.DB, order any) *gorm.DB {
	return db.Order(order) // want "ORDER BY built from an unsanitized value"
}

type anyFilter struct {
	SortOn any
}

func anyStructField(db *gorm.DB, f anyFilter) *gorm.DB {
	return db.Order(f.SortOn) // want "ORDER BY built from an unsanitized value"
}

// namedInterface widens one interface into another rather than boxing a
// concrete value, and a named string type can satisfy it.
type orderStringer interface {
	String() string
}

func namedInterface(db *gorm.DB, order orderStringer) *gorm.DB {
	return db.Order(order) // want "ORDER BY built from an unsanitized value"
}

// throughInterfaceUnsanitized pins that abstracting the handler behind an
// interface is not a way out of the rule.
func throughInterfaceUnsanitized(o orderer, order string) *gorm.DB {
	return o.Order(order) // want "ORDER BY built from an unsanitized value"
}

// unsanitizedInLoop exercises the same phi cycle as sanitizedInLoop, with an
// edge the pass cannot vouch for.
func unsanitizedInLoop(db *gorm.DB, candidates []string) *gorm.DB {
	order := defaultOrder
	for _, candidate := range candidates {
		if candidate != "" {
			order = candidate
		}
	}
	return db.Order(order) // want "ORDER BY built from an unsanitized value"
}

func promotedUnsanitized(h *handler, order string) *gorm.DB {
	return h.Order(order) // want "ORDER BY built from an unsanitized value"
}

// inClosure pins that the pass descends into anonymous functions - gorm scopes
// are routinely built inside one.
func inClosure(db *gorm.DB, order string) func() *gorm.DB {
	return func() *gorm.DB {
		return db.Order(order) // want "ORDER BY built from an unsanitized value"
	}
}

// sanitizedInClosure pins the pass's one known false positive. Capturing
// `order` moves it out of an SSA register and into a memory cell, and the pass
// will not guess which store reaches the load - see "Captured variables" in the
// package doc. Reporting a sanitized value here is the deliberate trade: it is
// loud and fixable, where the opposite error is the CVE.
func sanitizedInClosure(db *gorm.DB, order string) func() *gorm.DB {
	order = models.SanitizeOrderInput(order, validColumns)
	return func() *gorm.DB {
		return db.Order(order) // want "ORDER BY built from an unsanitized value"
	}
}

// capturedElsewhere is the same limitation seen from outside: the Order call is
// not in the closure, but capturing `order` anywhere in the function is what
// put it in memory.
func capturedElsewhere(db *gorm.DB, order string) *gorm.DB {
	order = models.SanitizeOrderInput(order, validColumns)
	if order == "" {
		order = defaultOrder
	}
	logOrder := func() string { return order }
	_ = logOrder
	return db.Order(order) // want "ORDER BY built from an unsanitized value"
}

// capturedThenRead is the fix for both of the above: read the sanitized value
// out before anything captures it, and capture the copy.
func capturedThenRead(db *gorm.DB, order string) (*gorm.DB, func() string) {
	sanitized := models.SanitizeOrderInput(order, validColumns)
	captured := sanitized
	logOrder := func() string { return captured }
	return db.Order(sanitized), logOrder
}

// rewrittenInsideClosure is the case the conservative treatment buys: the
// sanitizer runs, and then a closure overwrites the same variable with request
// input at a point no register-level walk would see.
func rewrittenInsideClosure(db *gorm.DB, order, raw string) *gorm.DB {
	order = models.SanitizeOrderInput(order, validColumns)
	override := func() { order = raw }
	override()
	return db.Order(order) // want "ORDER BY built from an unsanitized value"
}

// addressEscapes hands the variable to a function that could write anything
// into it, so the pass refuses to vouch for what comes back out.
func addressEscapes(db *gorm.DB, order string) *gorm.DB {
	order = models.SanitizeOrderInput(order, validColumns)
	overwrite(&order)
	return db.Order(order) // want "ORDER BY built from an unsanitized value"
}

func overwrite(order *string) { *order = "anything" }

//
// Rejected: the argument is a type parameter. SSA builds a generic body once
// with the parameter intact, so the boxed type is a *types.TypeParam rather
// than a concrete one. Instantiated at string it reaches gorm's `case string:`
// branch like any other string, so skipping it as a clause value would be a
// silent false negative.
//

// genericOrder is the reachable shape: a helper written once over ~string and
// called with a request value.
func genericOrder[T ~string](db *gorm.DB, order T) *gorm.DB {
	return db.Order(order) // want "ORDER BY built from an unsanitized value"
}

// genericOrderAny is the same hole through an unconstrained parameter.
func genericOrderAny[T any](db *gorm.DB, order T) *gorm.DB {
	return db.Order(order) // want "ORDER BY built from an unsanitized value"
}

// instantiateGenericOrder pins that the generic helpers above are instantiated
// at string by a real caller, so the diagnostic is about a reachable path and
// not an abstract one.
func instantiateGenericOrder(db *gorm.DB, order string) *gorm.DB {
	db = genericOrder(db, order)
	return genericOrderAny(db, order)
}

// genericSanitized still has to satisfy the rule the ordinary way - the type
// parameter does not exempt it, and a sanitized value passes.
func genericSanitized[T ~string](db *gorm.DB, order string) *gorm.DB {
	return db.Order(T(models.SanitizeOrderInput(order, validColumns)))
}

// multiTermConstraint pins the conversion SSA emits when a type parameter's
// type set has more than one term: `T(s)` under `T ~string | ~[]byte` becomes a
// MultiConvert rather than a ChangeType. It reshapes the value like any other
// conversion, so a sanitized operand still has to be accepted.
func multiTermConstraint[T ~string | ~[]byte](db *gorm.DB, order string) *gorm.DB {
	return db.Order(T(models.SanitizeOrderInput(order, validColumns)))
}

func instantiateMultiTermConstraint(db *gorm.DB, order string) *gorm.DB {
	return multiTermConstraint[string](db, order)
}

// sortKey is a named string type carrying a method, so a value of it can be
// held in a narrower interface and then widened.
type sortKey string

func (s sortKey) String() string { return string(s) }

type stringish interface{ String() string }

// widenedInterface pins ChangeInterface: the sanitized value is boxed into
// `stringish` and then widened to `any` on the way into Order. That widening is
// value-preserving, so the sanitizer result still has to be traceable through
// it.
func widenedInterface(db *gorm.DB, order string) *gorm.DB {
	var narrow stringish = sortKey(models.SanitizeOrderInput(order, validColumns))
	var widened any = narrow
	return db.Order(widened)
}

// lookalike is not gorm's Order: the rule keys on the `Order(any) *gorm.DB`
// signature, so an unrelated method of the same name is left alone.
type sorter struct{}

func (s sorter) Order(by string) string { return by }

func lookalike(s sorter, by string) string {
	return s.Order(by)
}
