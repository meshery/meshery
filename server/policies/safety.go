package policies

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// TieInfo records a conflict between two or more actions that targeted the
// same field on the same component with different values at equal priority.
type TieInfo struct {
	ComponentID string
	Path        []string
}

// resolveTies drops component-update actions that conflict on the same target.
// Since all policies currently share the same priority, any divergence on the
// same (componentID, path) is treated as a tie: both actions are dropped.
// Returns the filtered list and the ties that were detected.
func resolveTies(actions []PolicyAction) ([]PolicyAction, []TieInfo) {
	type fp struct {
		op    string
		value string
	}
	groups := make(map[string][]int)
	prints := make(map[string][]fp)

	for i, a := range actions {
		if !isComponentPatch(a.Op) || a.ID == "" {
			continue
		}
		key := a.ID + "\x00" + strings.Join(a.UpdatePath, ".")
		groups[key] = append(groups[key], i)
		prints[key] = append(prints[key], fp{op: a.Op, value: canonicalJSON(a.UpdateValue)})
	}

	// Iterate groups in sorted key order so dropped-action order and the
	// resulting ties slice are deterministic across runs.
	sortedKeys := make([]string, 0, len(groups))
	for k := range groups {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	drop := make(map[int]bool)
	var ties []TieInfo
	for _, key := range sortedKeys {
		idxs := groups[key]
		if len(idxs) < 2 {
			continue
		}
		first := prints[key][0]
		conflict := false
		for _, p := range prints[key][1:] {
			if p != first {
				conflict = true
				break
			}
		}
		if !conflict {
			// All agree: keep the first, drop the rest (dedupe).
			for _, i := range idxs[1:] {
				drop[i] = true
			}
			continue
		}
		// Equal-priority tie: drop every conflicting action.
		for _, i := range idxs {
			drop[i] = true
		}
		ties = append(ties, TieInfo{ComponentID: actions[idxs[0]].ID, Path: actions[idxs[0]].UpdatePath})
	}

	if len(drop) == 0 {
		return actions, ties
	}
	filtered := make([]PolicyAction, 0, len(actions)-len(drop))
	for i, a := range actions {
		if !drop[i] {
			filtered = append(filtered, a)
		}
	}
	return filtered, ties
}

func isComponentPatch(op string) bool {
	return op == UpdateComponentConfigurationOp || op == UpdateComponentOp
}

// canonicalJSON produces a stable string representation for comparison.
// On marshal failure it returns a type-tagged sentinel so unmarshalable values
// of different types do not collapse into the same key (and thus falsely tie).
func canonicalJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("<unmarshalable:%T:%v>", v, err)
	}
	return string(b)
}

// actionsFingerprint returns a deterministic hash of the actions produced by
// one evaluation. It captures adds/removes of relationships and all component
// patches, sorted so order does not affect the fingerprint.
func actionsFingerprint(actions []PolicyAction) string {
	lines := make([]string, 0, len(actions))
	for _, a := range actions {
		switch a.Op {
		case AddRelationshipOp, DeleteRelationshipOp, UpdateRelationshipOp:
			id := a.ID
			if id == "" && a.Relationship != nil {
				id = a.Relationship.ID.String()
			}
			lines = append(lines, a.Op+"|"+id)
		case UpdateComponentOp, UpdateComponentConfigurationOp:
			lines = append(lines, a.Op+"|"+a.ID+"|"+strings.Join(a.UpdatePath, ".")+"|"+canonicalJSON(a.UpdateValue))
		case AddComponentOp:
			if a.Component != nil {
				lines = append(lines, a.Op+"|"+a.Component.ID.String())
			}
		case DeleteComponentOp:
			lines = append(lines, a.Op+"|"+a.ID)
		}
	}
	sort.Strings(lines)
	sum := sha256.Sum256([]byte(strings.Join(lines, "\n")))
	return hex.EncodeToString(sum[:])
}

// flapHistoryLimit is the number of most-recent fingerprints kept per design.
const flapHistoryLimit = 4

// recordAndCheckFlap appends current to the history for designID and returns
// true if the last four evaluations form an A-B-A-B alternation (i.e. three
// consecutive inversions).
func recordAndCheckFlap(history map[string][]string, designID, current string) bool {
	h := append(history[designID], current)
	if len(h) > flapHistoryLimit {
		h = h[len(h)-flapHistoryLimit:]
	}
	history[designID] = h
	if len(h) < flapHistoryLimit {
		return false
	}
	return h[0] != h[1] && h[0] == h[2] && h[1] == h[3]
}
