package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/pattern/patterns"
	events "github.com/meshery/meshkit/models/events"
)

func dmsg(comp string, ok bool, err error, prereq bool, msg string) patterns.DeploymentMessagePerComp {
	return patterns.DeploymentMessagePerComp{CompName: comp, Success: ok, Error: err, IsPrerequisite: prereq, Message: msg}
}

func TestDeploymentOutcome_EdgeCases(t *testing.T) {
	tests := []struct {
		name, action, design, okDesc string
		resp                         map[string]interface{}
		total, succ, fail            int
		severity                     events.EventSeverity
		desc                         string
	}{
		{
			name: "all ok deploy", action: "deploy", design: "d1", okDesc: "Deployment completed for design 'd1'.",
			resp: map[string]interface{}{
				"c1": []patterns.DeploymentMessagePerContext{{Summary: []patterns.DeploymentMessagePerComp{dmsg("c1", true, nil, false, "")}}},
				"c2": []patterns.DeploymentMessagePerContext{{Summary: []patterns.DeploymentMessagePerComp{dmsg("c2", true, nil, false, "")}}},
			},
			total: 2, succ: 2, fail: 0, severity: events.Informational, desc: "Deployment completed for design 'd1'.",
		},
		{
			name: "all ok undeploy", action: "undeploy", design: "d1", okDesc: "undeployed design 'd1'.",
			resp: map[string]interface{}{
				"c1": []patterns.DeploymentMessagePerContext{{Summary: []patterns.DeploymentMessagePerComp{dmsg("c1", true, nil, false, "")}}},
			},
			total: 1, succ: 1, fail: 0, severity: events.Success, desc: "undeployed design 'd1'.",
		},
		{
			name: "1 of N fails deploy", action: "deploy", design: "d1", okDesc: "Deployment completed for design 'd1'.",
			resp: map[string]interface{}{
				"c1": []patterns.DeploymentMessagePerContext{{Summary: []patterns.DeploymentMessagePerComp{dmsg("c1", true, nil, false, "")}}},
				"c2": []patterns.DeploymentMessagePerContext{{Summary: []patterns.DeploymentMessagePerComp{dmsg("c2", false, errors.New("err"), false, "")}}},
			},
			total: 2, succ: 1, fail: 1, severity: events.Warning, desc: "Deployment partially completed for design 'd1': 1 of 2 components failed.",
		},
		{
			name: "1 of N fails undeploy", action: "undeploy", design: "d1", okDesc: "undeployed design 'd1'.",
			resp: map[string]interface{}{
				"c1": []patterns.DeploymentMessagePerContext{{Summary: []patterns.DeploymentMessagePerComp{dmsg("c1", true, nil, false, "")}}},
				"c2": []patterns.DeploymentMessagePerContext{{Summary: []patterns.DeploymentMessagePerComp{dmsg("c2", false, errors.New("err"), false, "")}}},
			},
			total: 2, succ: 1, fail: 1, severity: events.Warning, desc: "Undeployment partially completed for design 'd1': 1 of 2 components failed.",
		},
		{
			name: "all fail deploy", action: "deploy", design: "d1", okDesc: "Deployment completed for design 'd1'.",
			resp: map[string]interface{}{
				"c1": []patterns.DeploymentMessagePerContext{{Summary: []patterns.DeploymentMessagePerComp{dmsg("c1", false, errors.New("e1"), false, "")}}},
			},
			total: 1, succ: 0, fail: 1, severity: events.Error, desc: "Deployment failed for design 'd1': 1 of 1 components failed.",
		},
		{
			name: "all fail undeploy", action: "undeploy", design: "d1", okDesc: "undeployed design 'd1'.",
			resp: map[string]interface{}{
				"c1": []patterns.DeploymentMessagePerContext{{Summary: []patterns.DeploymentMessagePerComp{dmsg("c1", false, errors.New("e1"), false, "")}}},
			},
			total: 1, succ: 0, fail: 1, severity: events.Error, desc: "Undeployment failed for design 'd1': 1 of 1 components failed.",
		},
		{
			name: "empty design (no 0 of 0)", action: "deploy", design: "d1", okDesc: "Deployment completed for design 'd1'.",
			resp:  map[string]interface{}{},
			total: 0, succ: 0, fail: 0, severity: events.Informational, desc: "Deployment completed for design 'd1'.",
		},
		{
			name: "empty design undeploy", action: "undeploy", design: "d1", okDesc: "undeployed design 'd1'.",
			resp:  map[string]interface{}{},
			total: 0, succ: 0, fail: 0, severity: events.Success, desc: "undeployed design 'd1'.",
		},
		{
			name: "failed prerequisite but component ok = success", action: "deploy", design: "d1", okDesc: "Deployment completed for design 'd1'.",
			resp: map[string]interface{}{
				"c1": []patterns.DeploymentMessagePerContext{{Summary: []patterns.DeploymentMessagePerComp{
					dmsg("c1-crd", false, errors.New("crd err"), true, ""),
					dmsg("c1", true, nil, false, ""),
				}}},
			},
			total: 1, succ: 1, fail: 0, severity: events.Informational, desc: "Deployment completed for design 'd1'.",
		},
		{
			name: "withheld = failed", action: "deploy", design: "d1", okDesc: "Deployment completed for design 'd1'.",
			resp: map[string]interface{}{
				"c1": []patterns.DeploymentMessagePerContext{{Summary: []patterns.DeploymentMessagePerComp{
					dmsg("c1", false, errors.New("dep unsatisfied"), false, "Withheld c1"),
				}}},
			},
			total: 1, succ: 0, fail: 1, severity: events.Error, desc: "Deployment failed for design 'd1': 1 of 1 components failed.",
		},
		{
			name: "dispatch-failure and no-outcome entries", action: "deploy", design: "d1", okDesc: "Deployment completed for design 'd1'.",
			resp: map[string]interface{}{
				"c1": []patterns.DeploymentMessagePerContext{{Summary: []patterns.DeploymentMessagePerComp{dmsg("c1", false, errors.New("dispatch err"), false, "Could not deploy")}}},
				"c2": []patterns.DeploymentMessagePerContext{{Summary: []patterns.DeploymentMessagePerComp{dmsg("c2", false, nil, false, "Could not confirm")}}},
			},
			total: 2, succ: 0, fail: 2, severity: events.Error, desc: "Deployment failed for design 'd1': 2 of 2 components failed.",
		},
		{
			name: "multi-cluster with A ok and B fail = 1 failed", action: "deploy", design: "d1", okDesc: "Deployment completed for design 'd1'.",
			resp: map[string]interface{}{
				"c1": []patterns.DeploymentMessagePerContext{
					{Location: "cluster-a", Summary: []patterns.DeploymentMessagePerComp{dmsg("c1", true, nil, false, "")}},
					{Location: "cluster-b", Summary: []patterns.DeploymentMessagePerComp{dmsg("c1", false, errors.New("errB"), false, "")}},
				},
			},
			total: 1, succ: 0, fail: 1, severity: events.Error, desc: "Deployment failed for design 'd1': 1 of 1 components failed.",
		},
		{
			name: "non-slice values and the \"\" key ignored", action: "deploy", design: "d1", okDesc: "Deployment completed for design 'd1'.",
			resp: map[string]interface{}{
				"":          []patterns.DeploymentMessagePerContext{{Summary: []patterns.DeploymentMessagePerComp{dmsg("empty", false, errors.New("err"), false, "")}}},
				"nonSlice":  "notASlice",
				"validComp": []patterns.DeploymentMessagePerContext{{Summary: []patterns.DeploymentMessagePerComp{dmsg("validComp", true, nil, false, "")}}},
			},
			total: 1, succ: 1, fail: 0, severity: events.Informational, desc: "Deployment completed for design 'd1'.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			summary := summarizePatternProcessResponse(tt.resp)
			if summary.Total != tt.total || summary.Succeeded != tt.succ || summary.Failed != tt.fail {
				t.Fatalf("got total=%d, succ=%d, fail=%d, want total=%d, succ=%d, fail=%d", summary.Total, summary.Succeeded, summary.Failed, tt.total, tt.succ, tt.fail)
			}
			sev, desc := deploymentOutcomeEvent(tt.action, summary, tt.design, tt.okDesc)
			if sev != tt.severity || desc != tt.desc {
				t.Errorf("got (%v, %q), want (%v, %q)", sev, desc, tt.severity, tt.desc)
			}
			if tt.total == 0 && strings.Contains(desc, "0 of 0") {
				t.Errorf("empty design should not produce '0 of 0'")
			}
		})
	}
}

// TestDeploymentOutcome_FailedCountsComponentsNotOutcomes pins that failed and
// errorCount use different units: a component that fails on two clusters is
// one failed component but two failing outcomes.
func TestDeploymentOutcome_FailedCountsComponentsNotOutcomes(t *testing.T) {
	resp := map[string]interface{}{
		"c1": []patterns.DeploymentMessagePerContext{
			{Location: "cluster-a", Summary: []patterns.DeploymentMessagePerComp{dmsg("c1", false, errors.New("errA"), false, "")}},
			{Location: "cluster-b", Summary: []patterns.DeploymentMessagePerComp{dmsg("c1", false, errors.New("errB"), false, "")}},
		},
		"c2": []patterns.DeploymentMessagePerContext{
			{Location: "cluster-a", Summary: []patterns.DeploymentMessagePerComp{dmsg("c2", true, nil, false, "")}},
			{Location: "cluster-b", Summary: []patterns.DeploymentMessagePerComp{dmsg("c2", true, nil, false, "")}},
		},
	}

	summary := summarizePatternProcessResponse(resp)
	if summary.Total != 2 || summary.Succeeded != 1 || summary.Failed != 1 {
		t.Errorf("got total=%d succeeded=%d failed=%d, want 2/1/1", summary.Total, summary.Succeeded, summary.Failed)
	}
	if summary.ErrorCount != 2 {
		t.Errorf("got errorCount=%d, want 2 (one per failing cluster)", summary.ErrorCount)
	}
	if summary.Status != patternProcessStatusPartial {
		t.Errorf("got status %q, want %q", summary.Status, patternProcessStatusPartial)
	}
}

// TestDeploymentOutcome_CountsSurviveMetadataSizeGuard pins that a large
// partial failure keeps its counts: the outcome adds only three integers to
// the capped metadata, so the event stays under models.MaxEventMetadataBytes
// and is never reduced to the bare design identity.
func TestDeploymentOutcome_CountsSurviveMetadataSizeGuard(t *testing.T) {
	longErr := errors.New(strings.Repeat("admission webhook denied the request: ", 20))
	resp := map[string]interface{}{}
	for i := 0; i < 500; i++ {
		ok := i%2 == 0
		var err error
		if !ok {
			err = longErr
		}
		resp[fmt.Sprintf("component-%03d", i)] = []patterns.DeploymentMessagePerContext{{
			Location: "https://127.0.0.1:6443",
			Summary:  []patterns.DeploymentMessagePerComp{dmsg(fmt.Sprintf("component-%03d", i), ok, err, false, "")},
		}}
	}

	summary := summarizePatternProcessResponse(resp)
	metadata := buildPatternProcessEventMetadata(summary, "http://meshery.test", uuid.Must(uuid.NewV4()), "big", "")
	addOutcomeCounts(metadata, summary)

	b, err := json.Marshal(metadata)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	if len(b) > models.MaxEventMetadataBytes {
		t.Fatalf("metadata is %d bytes, over the %d-byte guard", len(b), models.MaxEventMetadataBytes)
	}
	if metadata["total"] != 500 || metadata["succeeded"] != 250 || metadata["failed"] != 250 {
		t.Errorf("got total=%v succeeded=%v failed=%v, want 500/250/250", metadata["total"], metadata["succeeded"], metadata["failed"])
	}
}

// TestDeploymentOutcome_ExampleEvent pins the event a real partial deployment
// produces: three Kubernetes components on one cluster, one rejected by the
// API server. Run with -v to see the event JSON.
func TestDeploymentOutcome_ExampleEvent(t *testing.T) {
	const host = "https://127.0.0.1:6443"
	ok := func(name, kind string) []patterns.DeploymentMessagePerContext {
		return []patterns.DeploymentMessagePerContext{{SystemName: "kind-meshery", Location: host, Summary: []patterns.DeploymentMessagePerComp{
			{Kind: kind, Model: "kubernetes", CompName: name, Success: true, DesignName: "bookinfo", Message: "deployed bookinfo/" + name},
		}}}
	}
	resp := map[string]interface{}{
		"productpage": ok("productpage", "Deployment"),
		"ratings":     ok("ratings", "Service"),
		"reviews": []patterns.DeploymentMessagePerContext{{SystemName: "kind-meshery", Location: host, Summary: []patterns.DeploymentMessagePerComp{
			{Kind: "Deployment", Model: "kubernetes", CompName: "reviews", DesignName: "bookinfo", Message: "Error deploying bookinfo/reviews",
				Error: errors.New(`Deployment.apps "reviews" is invalid: spec.template.spec.containers[0].image: Required value`)},
		}}},
	}

	summary := summarizePatternProcessResponse(resp)
	severity, description := deploymentOutcomeEvent("deploy", summary, "bookinfo", "Deployment completed for design 'bookinfo'.")
	designID := uuid.FromStringOrNil("8c5a1f4e-2b7d-4e0a-9f3c-6d1e2a4b7c90")
	metadata := buildPatternProcessEventMetadata(summary, "http://localhost:9081", designID, "bookinfo", "")
	addOutcomeCounts(metadata, summary)

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(map[string]interface{}{
		"severity":    severity,
		"description": description,
		"metadata":    metadata,
	}); err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	got := strings.TrimSuffix(buf.String(), "\n")
	t.Logf("example event:\n%s", got)

	want := `{
  "description": "Deployment partially completed for design 'bookinfo': 1 of 3 components failed.",
  "metadata": {
    "designId": "8c5a1f4e-2b7d-4e0a-9f3c-6d1e2a4b7c90",
    "designLink": "http://localhost:9081/extension/meshmap?mode=design&design=8c5a1f4e-2b7d-4e0a-9f3c-6d1e2a4b7c90",
    "designName": "bookinfo",
    "errorCount": 1,
    "errors": [
      {
        "code": "",
        "componentId": "reviews"
      }
    ],
    "failed": 1,
    "status": "partial",
    "succeeded": 2,
    "total": 3,
    "truncated": false,
    "viewLink": "http://localhost:9081/extension/meshmap?mode=operator&type=view&design_id=8c5a1f4e-2b7d-4e0a-9f3c-6d1e2a4b7c90"
  },
  "severity": "warning"
}`
	if got != want {
		t.Errorf("example event changed:\n got: %s\nwant: %s", got, want)
	}
}
