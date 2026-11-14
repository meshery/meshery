package handlers

import (
	"encoding/json"
	"errors"
	"strconv"
	"testing"

	"github.com/meshery/meshery/server/models"
)

type deletionContextPages struct {
	pages    [][]*models.K8sContext
	failPage int
	calls    int
}

func (p *deletionContextPages) GetK8sContexts(_, page, _, _, _, status string, _ bool) ([]byte, error) {
	if status != "" {
		return nil, errors.New("context snapshot was filtered by status")
	}
	n, err := strconv.Atoi(page)
	if err != nil {
		return nil, err
	}
	p.calls++
	if n == p.failPage {
		return nil, errors.New("read failure")
	}
	total := 0
	for _, page := range p.pages {
		total += len(page)
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
