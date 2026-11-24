package handlers

import (
	"encoding/json"
	"errors"
	"strconv"
	"testing"

	"github.com/meshery/meshery/server/models"
)

type deletionContextPages struct {
	pages               [][]*models.K8sContext
	failPage            int
	calls               int
	changedTotal        bool
	totalOverride       *int
	seenWithCredentials *bool
}

func (p *deletionContextPages) GetK8sContexts(_, page, _, _, _, status string, withCredentials bool) ([]byte, error) {
	if status != "" {
		return nil, errors.New("context snapshot was filtered by status")
	}
	p.seenWithCredentials = &withCredentials
	n, err := strconv.Atoi(page)
	if err != nil {
		return nil, err
	}
	p.calls++
	if n == p.failPage {
		return nil, errors.New("read failure")
	}
	total := 0
	if p.totalOverride != nil {
		total = *p.totalOverride
	} else {
		for _, page := range p.pages {
			total += len(page)
		}
	}
	if p.changedTotal && n > 0 {
		total--
	}
	return json.Marshal(models.MesheryK8sContextPage{Contexts: p.pages[n], TotalCount: total})
}

func TestLoadContextsForDeletion(t *testing.T) {
	p := &deletionContextPages{pages: [][]*models.K8sContext{{{ID: "one"}}, {{ID: "two"}}}, failPage: -1}
	got, err := loadContextsForDeletion(p, "token")
	if err != nil || len(got) != 2 || p.calls != 2 {
		t.Fatalf("contexts=%v calls=%d err=%v", got, p.calls, err)
	}
	p.failPage = 1
	got, err = loadContextsForDeletion(p, "token")
	if err == nil || got != nil {
		t.Fatalf("partial snapshot accepted: contexts=%v err=%v", got, err)
	}
}

func TestLoadContextsForDeletionRejectsChangingOrDuplicatePages(t *testing.T) {
	for _, tc := range []struct {
		name string
		p    *deletionContextPages
	}{
		{"changing count", &deletionContextPages{pages: [][]*models.K8sContext{{{ID: "one"}}, {{ID: "two"}}}, failPage: -1, changedTotal: true}},
		{"duplicate context", &deletionContextPages{pages: [][]*models.K8sContext{{{ID: "one"}}, {{ID: "one"}}}, failPage: -1}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := loadContextsForDeletion(tc.p, "token")
			if err == nil || got != nil {
				t.Fatalf("accepted inconsistent snapshot: %v, %v", got, err)
			}
		})
	}
}

// The snapshot must be metadata-only so a rejected credential cannot drop a
// context from it (issue #24279).
func TestLoadContextsForDeletionRequestsMetadataOnly(t *testing.T) {
	p := &deletionContextPages{pages: [][]*models.K8sContext{{{ID: "one"}}}, failPage: -1}
	if _, err := loadContextsForDeletion(p, "token"); err != nil {
		t.Fatal(err)
	}
	if p.seenWithCredentials == nil || *p.seenWithCredentials {
		t.Fatalf("withCredentials = %v, want false", p.seenWithCredentials)
	}
}

// A short list must never be accepted: a missing context could be the one that
// still uses the cluster, and a flush would then wipe its data.
func TestLoadContextsForDeletionRejectsShortList(t *testing.T) {
	total := 3
	p := &deletionContextPages{
		pages:         [][]*models.K8sContext{{{ID: "one"}}, {{ID: "two"}}, {}},
		failPage:      -1,
		totalOverride: &total,
	}
	if got, err := loadContextsForDeletion(p, "token"); err == nil || got != nil {
		t.Fatalf("short list accepted: contexts=%v err=%v", got, err)
	}
}
