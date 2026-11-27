package policies

import (
	"slices"
	"sort"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/pattern"
	"github.com/meshery/schemas/models/v1beta2/relationship"
)

// EdgeBindingPolicy handles edge binding relationships (3-party: from, binding, to).
type EdgeBindingPolicy struct{}

func (p *EdgeBindingPolicy) Identifier() string {
	return "edge_binding"
}

func (p *EdgeBindingPolicy) IsImplicatedBy(rel *relationship.RelationshipDefinition) bool {
	return strings.EqualFold(string(rel.Kind), "edge") &&
		strings.EqualFold(rel.RelationshipType, "binding")
}

func (p *EdgeBindingPolicy) IsInvalid(rel *relationship.RelationshipDefinition, design *pattern.PatternFile) bool {
	return fromOrToComponentsDontExist(rel, design)
}

func (p *EdgeBindingPolicy) AlreadyExists(rel *relationship.RelationshipDefinition, design *pattern.PatternFile) bool {
	return relationshipAlreadyExists(design, rel)
}

func (p *EdgeBindingPolicy) IdentifyRelationship(relDef *relationship.RelationshipDefinition, design *pattern.PatternFile) []*relationship.RelationshipDefinition {
	return identifyBindingRelationships(relDef, design)
}

func (p *EdgeBindingPolicy) SideEffects(rel *relationship.RelationshipDefinition, design *pattern.PatternFile) []PolicyAction {
	actions := patchMutatorsAction(rel, design)
	if len(actions) > 0 {
		return actions
	}
	return patchBindingMatchFields(rel, design)
}

// patchBindingMatchFields patches the binding component's configuration
// using values from the match fields of the binding relationship selectors.
// When the relationship is in StatusDeleted, it emits reverse patches (schema
// default when declared, otherwise remove) so deletion restores the
// pre-mutation state.
func patchBindingMatchFields(rel *relationship.RelationshipDefinition, design *pattern.PatternFile) []PolicyAction {
	if rel.Selectors == nil {
		return nil
	}
	reverse := getRelStatus(rel) == StatusDeleted
	var actions []PolicyAction

	for _, ss := range *rel.Selectors {
		for _, selArr := range [][]relationship.SelectorItem{ss.Allow.From, ss.Allow.To} {
			for _, sel := range selArr {
				if sel.Match == nil {
					continue
				}
				compID := selectorItemID(sel)
				comp := componentDeclarationByID(design, compID)
				if comp == nil {
					continue
				}

				actions = append(actions, patchBindingMatchFieldsForSelector(sel.Match, comp, design, reverse)...)
			}
		}
	}
	return actions
}

// patchBindingMatchFieldsForSelector processes match fields for a single selector.
// If reverse is true, emits cleanup patches (schema default or remove) for fields
// that still hold the mutator's value.
func patchBindingMatchFieldsForSelector(match *relationship.MatchSelector, _ *component.ComponentDefinition, design *pattern.PatternFile, reverse bool) []PolicyAction {
	var actions []PolicyAction

	type matchSide struct {
		self  *[]relationship.MatchSelectorItem
		other *[]relationship.MatchSelectorItem
	}
	sides := []matchSide{
		{match.From, match.To},
		{match.To, match.From},
	}
	for _, ms := range sides {
		if ms.self == nil || ms.other == nil || len(*ms.self) == 0 || len(*ms.other) == 0 {
			continue
		}
		selfLen := len(*ms.self)
		otherLen := len(*ms.other)
		pairCount := max(selfLen, otherLen)
		if selfLen > 1 && otherLen > 1 {
			pairCount = min(selfLen, otherLen)
		}
		for p := 0; p < pairCount; p++ {
			selfIdx := p
			if selfIdx >= selfLen {
				selfIdx = 0
			}
			otherIdx := p
			if otherIdx >= otherLen {
				otherIdx = 0
			}
			mSel := (*ms.self)[selfIdx]
			if mSel.MutatorRef == nil {
				continue
			}
			matchComp := findMatchSelectorComponent(mSel, design)
			if matchComp == nil {
				continue
			}

			otherSel := (*ms.other)[otherIdx]
			if otherSel.MutatedRef == nil {
				continue
			}
			otherComp := findMatchSelectorComponent(otherSel, design)
			if otherComp == nil {
				continue
			}

			count := min(len(*mSel.MutatorRef), len(*otherSel.MutatedRef))

			for i := range count {
				origMutatorPath := (*mSel.MutatorRef)[i]
				origMutatedPath := (*otherSel.MutatedRef)[i]

				resolvedMutatorPaths := resolveComponentPaths(origMutatorPath, matchComp, design)
				resolvedMutatedPaths := resolveComponentPaths(origMutatedPath, otherComp, design)

				for j, mutatedPath := range resolvedMutatedPaths {
					oldValue := configurationForComponentAtPath(mutatedPath, otherComp, design)

					if reverse {
						if !matchesAnyMutatorValue(oldValue, resolvedMutatorPaths, matchComp, design) {
							continue
						}
						actions = append(actions, cleanupActionForPath(otherComp.ID.String(), mutatedPath, otherComp))
						continue
					}
					mutatorValue := selectMutatorValue(resolvedMutatorPaths, j, len(resolvedMutatedPaths), matchComp, design)
					if mutatorValue == nil {
						continue
					}
					if deepEqual(mutatorValue, oldValue) {
						continue
					}
					actions = append(actions, newComponentUpdateAction(getComponentUpdateOp(mutatedPath), otherComp.ID.String(), mutatedPath, mutatorValue))
				}
			}
		}
	}
	if reverse {
		actions = orderCleanupActions(actions)
	}
	return actions
}

// findMatchSelectorComponent finds the component referenced by a MatchSelectorItem.
func findMatchSelectorComponent(mSel relationship.MatchSelectorItem, design *pattern.PatternFile) *component.ComponentDefinition {
	if mSel.ID == nil {
		return nil
	}
	return componentDeclarationByID(design, mSel.ID.String())
}

// identifyBindingRelationships identifies 3-party binding relationships.
func identifyBindingRelationships(relDef *relationship.RelationshipDefinition, design *pattern.PatternFile) []*relationship.RelationshipDefinition {
	if relDef.Selectors == nil {
		return nil
	}
	var identified []*relationship.RelationshipDefinition

	// Every component takes part in O(N) candidate pairs below. Share its
	// generic-map form across all of them instead of re-encoding it per pair.
	maps := componentMapCache{}

	for _, ss := range *relDef.Selectors {
		fromSelByKind := make(map[string]relationship.SelectorItem)
		for _, fromSel := range ss.Allow.From {
			fromSelByKind[selectorItemKind(fromSel)] = fromSel
		}

		bindingKinds := collectBindingKindsTyped(ss.Allow.From)
		toSelByPair := indexToSelectorsByKindPairTyped(ss.Allow.To)

		fromComps := filterComponentsByKindSetTyped(design.Components, fromSelByKind)
		toComps := filterComponentsByToSelectorsTyped(design.Components, ss.Allow.To)

		// Iterate kinds in sorted order so identification output (and thus
		// generated relationship IDs) does not depend on map iteration randomness.
		sortedBindingKinds := make([]string, 0, len(bindingKinds))
		for k := range bindingKinds {
			sortedBindingKinds = append(sortedBindingKinds, k)
		}
		sort.Strings(sortedBindingKinds)
		for _, bindingKind := range sortedBindingKinds {
			bindingDecls := filterComponentsByKindTyped(design.Components, bindingKind)
			if len(bindingDecls) == 0 {
				continue
			}

			// Precompute, per bindingDecl, the toDecls whose (kind, mutated value)
			// validate against this binding. Each (bindingDecl, toDecl) pair is
			// checked exactly once here instead of once per matching fromDecl in
			// the loop below, which drops the dominating O(N³) work to O(N²).
			// Lists are appended in design order so output stays deterministic.
			type toMatch struct {
				decl *component.ComponentDefinition
				sel  relationship.SelectorItem
			}
			validToDecls := make(map[uuid.UUID][]toMatch, len(bindingDecls))
			for _, bindingDecl := range bindingDecls {
				for _, toDecl := range toComps {
					if toDecl.ID == bindingDecl.ID {
						continue
					}
					toKind := toDecl.Component.Kind
					toSel, exists := toSelByPair[toKind+"#"+bindingKind]
					if !exists {
						continue
					}
					if !isValidBinding(bindingDecl, toDecl, toSel, design, maps) {
						continue
					}
					validToDecls[bindingDecl.ID] = append(validToDecls[bindingDecl.ID], toMatch{decl: toDecl, sel: toSel})
				}
			}

			for _, fromDecl := range fromComps {
				fromKind := fromDecl.Component.Kind
				fromSel, exists := fromSelByKind[fromKind]
				if !exists {
					continue
				}

				for _, bindingDecl := range bindingDecls {
					if fromDecl.ID == bindingDecl.ID {
						continue
					}
					if !isValidBinding(fromDecl, bindingDecl, fromSel, design, maps) {
						continue
					}

					for _, tm := range validToDecls[bindingDecl.ID] {
						toDecl, toSel := tm.decl, tm.sel

						if ss.Deny != nil && isRelationshipDenied(fromDecl, toDecl, ss.Deny) {
							continue
						}

						newFromSel := patchBindingSelectorIDsTyped(&fromSel, fromDecl.ID.String(), bindingDecl.ID.String())
						newFromSel.ID = &fromDecl.ID

						newToSel := patchBindingSelectorIDsTyped(&toSel, bindingDecl.ID.String(), toDecl.ID.String())
						newToSel.ID = &toDecl.ID

						selectorSetItem := relationship.SelectorSetItem{
							Allow: relationship.Selector{
								From: []relationship.SelectorItem{newFromSel},
								To:   []relationship.SelectorItem{newToSel},
							},
						}

						seed := map[string]any{
							"from":    fromDecl.ID.String(),
							"binding": bindingDecl.ID.String(),
							"to":      toDecl.ID.String(),
							"relId":   relDef.ID.String(),
						}

						decl := deepCopyRelDef(relDef)
						decl.ID = staticUUID(seed)
						selectors := relationship.SelectorSet{selectorSetItem}
						decl.Selectors = &selectors
						setRelStatus(decl, StatusApproved)

						identified = append(identified, decl)
					}
				}
			}
		}
	}
	return identified
}

// collectBindingKindsTyped extracts binding component kinds from from selectors' match fields.
func collectBindingKindsTyped(fromSels []relationship.SelectorItem) map[string]bool {
	kinds := make(map[string]bool)
	for _, fromSel := range fromSels {
		if fromSel.Match == nil {
			continue
		}
		for _, matchSels := range []*[]relationship.MatchSelectorItem{fromSel.Match.From, fromSel.Match.To} {
			if matchSels == nil {
				continue
			}
			for _, mSel := range *matchSels {
				if mSel.Kind != "" && !strings.EqualFold(mSel.Kind, "self") {
					kinds[mSel.Kind] = true
				}
			}
		}
	}
	return kinds
}

// indexToSelectorsByKindPairTyped indexes to selectors by "toKind#matchKind".
func indexToSelectorsByKindPairTyped(toSels []relationship.SelectorItem) map[string]relationship.SelectorItem {
	result := make(map[string]relationship.SelectorItem)
	for _, toSel := range toSels {
		toKind := selectorItemKind(toSel)
		if toSel.Match == nil {
			continue
		}
		for _, matchSels := range []*[]relationship.MatchSelectorItem{toSel.Match.From, toSel.Match.To} {
			if matchSels == nil {
				continue
			}
			for _, mSel := range *matchSels {
				if mSel.Kind != "" && !strings.EqualFold(mSel.Kind, "self") {
					result[toKind+"#"+mSel.Kind] = toSel
				}
			}
		}
	}
	return result
}

// copyMatchSelector creates a deep copy of a MatchSelector, isolating MatchSelectorItem
// slices and their inner ID pointers so that ID patching does not mutate shared
// definitions or corrupt sibling relationships in multi-binding designs.
func copyMatchSelector(m *relationship.MatchSelector) *relationship.MatchSelector {
	if m == nil {
		return nil
	}
	cp := &relationship.MatchSelector{}
	if m.From != nil {
		fromSlice := make([]relationship.MatchSelectorItem, len(*m.From))
		for i, item := range *m.From {
			fromSlice[i] = item
			if item.ID != nil {
				id := *item.ID
				fromSlice[i].ID = &id
			}
		}
		cp.From = &fromSlice
	}
	if m.To != nil {
		toSlice := make([]relationship.MatchSelectorItem, len(*m.To))
		for i, item := range *m.To {
			toSlice[i] = item
			if item.ID != nil {
				id := *item.ID
				toSlice[i].ID = &id
			}
		}
		cp.To = &toSlice
	}
	if m.Refs != nil {
		refsSlice := make([][]string, len(*m.Refs))
		for i, inner := range *m.Refs {
			innerCopy := make([]string, len(inner))
			copy(innerCopy, inner)
			refsSlice[i] = innerCopy
		}
		cp.Refs = &refsSlice
	}
	return cp
}

// patchBindingSelectorIDsTyped deep-copies the selector item and sets the IDs in its match field,
// ensuring the source template and previously identified relationships remain unmodified.
func patchBindingSelectorIDsTyped(sel *relationship.SelectorItem, matchFromID, matchToID string) relationship.SelectorItem {
	cp := *sel
	if cp.Match == nil {
		return cp
	}
	cp.Match = copyMatchSelector(cp.Match)
	fromUUID, _ := uuidFromString(matchFromID)
	toUUID, _ := uuidFromString(matchToID)
	if cp.Match.From != nil {
		for i := range *cp.Match.From {
			(*cp.Match.From)[i].ID = &fromUUID
		}
	}
	if cp.Match.To != nil {
		for i := range *cp.Match.To {
			(*cp.Match.To)[i].ID = &toUUID
		}
	}
	return cp
}

// componentMapCache memoizes toGenericMap per component for the lifetime of
// one identification pass. toGenericMap is a JSON round-trip, and the result
// is only read (by resolvePath), so a component's map can be shared by every
// candidate pair it appears in. A nil cache disables memoization.
type componentMapCache map[*component.ComponentDefinition]map[string]interface{}

func (c componentMapCache) get(comp *component.ComponentDefinition) map[string]interface{} {
	if m, ok := c[comp]; ok {
		return m
	}
	m, _ := toGenericMap(comp)
	if c != nil {
		c[comp] = m
	}
	return m
}

// resolveComponentPath resolves "_" wildcards in path against comp. The
// generic map is only built when the path actually contains a wildcard;
// otherwise resolvePath would return the path unchanged without reading it.
func resolveComponentPath(path []string, comp *component.ComponentDefinition, maps componentMapCache) []string {
	if !slices.Contains(path, "_") {
		return path
	}
	return resolvePath(path, maps.get(comp))
}

// valueForComponentAtPath reads the value at path in comp. Lookups outside
// "configuration" go through the pass cache: configurationForComponentAtPath
// converts the component with toGenericMap on every call for such paths, so
// routing them through maps keeps the one-conversion-per-pass guarantee for
// non-configuration refs too (e.g. ["displayName"]). Configuration paths
// are unchanged, except that an aliased component's parent is also read from
// the cache: getComponentConfiguration converts the parent on every call,
// so without this the parent would be converted once per candidate pair. A
// failed conversion still yields nil, as before.
func valueForComponentAtPath(path []string, comp *component.ComponentDefinition, design *pattern.PatternFile, maps componentMapCache) interface{} {
	if len(path) == 0 {
		return nil
	}
	if path[0] == "configuration" {
		if parent, refPath := resolveComponentAlias(comp.ID.String(), design); parent != nil {
			return objectGetNested(objectGetNested(maps.get(parent), refPath, nil), popFirst(path), nil)
		}
		return configurationForComponentAtPath(path, comp, design)
	}
	return objectGetNested(maps.get(comp), path, nil)
}

// isValidBindingTyped checks if two components match via a binding selector's match field.
func isValidBindingTyped(comp1, comp2 *component.ComponentDefinition, sel relationship.SelectorItem, design *pattern.PatternFile) bool {
	return isValidBinding(comp1, comp2, sel, design, nil)
}

// isValidBinding is isValidBindingTyped with a caller-owned componentMapCache,
// so a caller checking many pairs converts each component at most once.
func isValidBinding(comp1, comp2 *component.ComponentDefinition, sel relationship.SelectorItem, design *pattern.PatternFile, maps componentMapCache) bool {
	if sel.Match == nil {
		return false
	}

	mutatorComp, mutatorPaths := extractMutatorFromMatch(sel.Match, comp1, comp2)
	mutatedComp, mutatedPaths := extractMutatedFromMatch(sel.Match, comp1, comp2)
	if mutatorComp == nil || mutatedComp == nil || len(mutatorPaths) == 0 || len(mutatedPaths) == 0 {
		return false
	}

	count := min(len(mutatorPaths), len(mutatedPaths))

	for i := range count {
		resolvedFrom := resolveComponentPath(mutatorPaths[i], mutatorComp, maps)
		resolvedTo := resolveComponentPath(mutatedPaths[i], mutatedComp, maps)
		val1 := valueForComponentAtPath(resolvedFrom, mutatorComp, design, maps)
		val2 := valueForComponentAtPath(resolvedTo, mutatedComp, design, maps)
		if val1 != nil && val2 != nil && !deepEqual(val1, val2) {
			return false
		}
	}
	return true
}

// extractMutatorFromMatch finds the mutator component and ref paths from a match field.
func extractMutatorFromMatch(match *relationship.MatchSelector, comp1, comp2 *component.ComponentDefinition) (*component.ComponentDefinition, [][]string) {
	if match.From != nil {
		for _, item := range *match.From {
			if item.MutatorRef != nil && len(*item.MutatorRef) > 0 {
				return comp1, *item.MutatorRef
			}
		}
	}
	if match.To != nil {
		for _, item := range *match.To {
			if item.MutatorRef != nil && len(*item.MutatorRef) > 0 {
				return comp2, *item.MutatorRef
			}
		}
	}
	return nil, nil
}

// extractMutatedFromMatch finds the mutated component and ref paths from a match field.
func extractMutatedFromMatch(match *relationship.MatchSelector, comp1, comp2 *component.ComponentDefinition) (*component.ComponentDefinition, [][]string) {
	if match.From != nil {
		for _, item := range *match.From {
			if item.MutatedRef != nil && len(*item.MutatedRef) > 0 {
				return comp1, *item.MutatedRef
			}
		}
	}
	if match.To != nil {
		for _, item := range *match.To {
			if item.MutatedRef != nil && len(*item.MutatedRef) > 0 {
				return comp2, *item.MutatedRef
			}
		}
	}
	return nil, nil
}

// filterComponentsByKindTyped returns components matching the given kind.
func filterComponentsByKindTyped(comps []*component.ComponentDefinition, kind string) []*component.ComponentDefinition {
	var result []*component.ComponentDefinition
	for _, comp := range comps {
		if matchName(comp.Component.Kind, kind) {
			result = append(result, comp)
		}
	}
	return result
}

// filterComponentsByKindSetTyped returns components whose kind is in the given set.
func filterComponentsByKindSetTyped(comps []*component.ComponentDefinition, kindSet map[string]relationship.SelectorItem) []*component.ComponentDefinition {
	var result []*component.ComponentDefinition
	for _, comp := range comps {
		if _, exists := kindSet[comp.Component.Kind]; exists {
			result = append(result, comp)
		}
	}
	return result
}

// filterComponentsByToSelectorsTyped returns components matching any to selector kind.
func filterComponentsByToSelectorsTyped(comps []*component.ComponentDefinition, toSels []relationship.SelectorItem) []*component.ComponentDefinition {
	toKinds := make(map[string]bool)
	for _, toSel := range toSels {
		toKinds[selectorItemKind(toSel)] = true
	}
	var result []*component.ComponentDefinition
	for _, comp := range comps {
		if toKinds[comp.Component.Kind] || toKinds["*"] {
			result = append(result, comp)
		}
	}
	return result
}
