package handlers

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

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
		fcCount                      int
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
			total: 2, succ: 1, fail: 1, fcCount: 1, severity: events.Warning, desc: "Deployment partially completed for design 'd1': 1 of 2 components failed.",
		},
		{
			name: "1 of N fails undeploy", action: "undeploy", design: "d1", okDesc: "undeployed design 'd1'.",
			resp: map[string]interface{}{
				"c1": []patterns.DeploymentMessagePerContext{{Summary: []patterns.DeploymentMessagePerComp{dmsg("c1", true, nil, false, "")}}},
				"c2": []patterns.DeploymentMessagePerContext{{Summary: []patterns.DeploymentMessagePerComp{dmsg("c2", false, errors.New("err"), false, "")}}},
			},
			total: 2, succ: 1, fail: 1, fcCount: 1, severity: events.Warning, desc: "Undeployment partially completed for design 'd1': 1 of 2 components failed.",
		},
		{
			name: "all fail deploy", action: "deploy", design: "d1", okDesc: "Deployment completed for design 'd1'.",
			resp: map[string]interface{}{
				"c1": []patterns.DeploymentMessagePerContext{{Summary: []patterns.DeploymentMessagePerComp{dmsg("c1", false, errors.New("e1"), false, "")}}},
			},
			total: 1, succ: 0, fail: 1, fcCount: 1, severity: events.Error, desc: "Deployment failed for design 'd1': 1 of 1 components failed.",
		},
		{
			name: "all fail undeploy", action: "undeploy", design: "d1", okDesc: "undeployed design 'd1'.",
			resp: map[string]interface{}{
				"c1": []patterns.DeploymentMessagePerContext{{Summary: []patterns.DeploymentMessagePerComp{dmsg("c1", false, errors.New("e1"), false, "")}}},
			},
			total: 1, succ: 0, fail: 1, fcCount: 1, severity: events.Error, desc: "Undeployment failed for design 'd1': 1 of 1 components failed.",
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
			total: 1, succ: 0, fail: 1, fcCount: 1, severity: events.Error, desc: "Deployment failed for design 'd1': 1 of 1 components failed.",
		},
		{
			name: "dispatch-failure and no-outcome entries", action: "deploy", design: "d1", okDesc: "Deployment completed for design 'd1'.",
			resp: map[string]interface{}{
				"c1": []patterns.DeploymentMessagePerContext{{Summary: []patterns.DeploymentMessagePerComp{dmsg("c1", false, errors.New("dispatch err"), false, "Could not deploy")}}},
				"c2": []patterns.DeploymentMessagePerContext{{Summary: []patterns.DeploymentMessagePerComp{dmsg("c2", false, nil, false, "Could not confirm")}}},
			},
			total: 2, succ: 0, fail: 2, fcCount: 2, severity: events.Error, desc: "Deployment failed for design 'd1': 2 of 2 components failed.",
		},
		{
			name: "multi-cluster with A ok and B fail = 1 failed naming B's location", action: "deploy", design: "d1", okDesc: "Deployment completed for design 'd1'.",
			resp: map[string]interface{}{
				"c1": []patterns.DeploymentMessagePerContext{
					{Location: "cluster-a", Summary: []patterns.DeploymentMessagePerComp{dmsg("c1", true, nil, false, "")}},
					{Location: "cluster-b", Summary: []patterns.DeploymentMessagePerComp{dmsg("c1", false, errors.New("errB"), false, "")}},
				},
			},
			total: 1, succ: 0, fail: 1, fcCount: 1, severity: events.Error, desc: "Deployment failed for design 'd1': 1 of 1 components failed.",
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
			tally := tallyDeploymentOutcome(tt.resp)
			if tally.Total != tt.total || tally.Succeeded != tt.succ || tally.Failed != tt.fail || len(tally.FailedComponents) != tt.fcCount {
				t.Fatalf("got tally %+v, want total=%d, succ=%d, fail=%d, fc=%d", tally, tt.total, tt.succ, tt.fail, tt.fcCount)
			}
			sev, desc := deploymentOutcomeEvent(tt.action, tally, tt.design, tt.okDesc)
			if sev != tt.severity || desc != tt.desc {
				t.Errorf("got (%v, %q), want (%v, %q)", sev, desc, tt.severity, tt.desc)
			}
			if tt.total == 0 && strings.Contains(desc, "0 of 0") {
				t.Errorf("empty design should not produce '0 of 0'")
			}
		})
	}
}

func TestDeploymentOutcome_SuccessFalseErrorNilFallback(t *testing.T) {
	resp := map[string]interface{}{
		"c1": []patterns.DeploymentMessagePerContext{{Summary: []patterns.DeploymentMessagePerComp{dmsg("c1", false, nil, false, "msg fallback")}}},
	}
	tally := tallyDeploymentOutcome(resp)
	if len(tally.FailedComponents) != 1 || tally.FailedComponents[0].Error != "msg fallback" {
		t.Fatalf("expected fallback to message, got %+v", tally.FailedComponents)
	}
}

func TestDeploymentOutcome_ErrorJSONSerialization(t *testing.T) {
	resp := map[string]interface{}{
		"c1": []patterns.DeploymentMessagePerContext{{Summary: []patterns.DeploymentMessagePerComp{dmsg("c1", false, errors.New("raw err"), false, "")}}},
	}
	tally := tallyDeploymentOutcome(resp)
	b, err := json.Marshal(map[string]interface{}{"failedComponents": tally.FailedComponents})
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	fc := raw["failedComponents"].([]interface{})[0].(map[string]interface{})
	if fc["error"] != "raw err" {
		t.Errorf("expected string 'raw err', got %v (%T)", fc["error"], fc["error"])
	}
}

func TestDeploymentOutcome_FailedComponentsSortedByNameThenLocation(t *testing.T) {
	resp := map[string]interface{}{
		"z": []patterns.DeploymentMessagePerContext{
			{Location: "loc2", Summary: []patterns.DeploymentMessagePerComp{dmsg("z", false, errors.New("e"), false, "")}},
			{Location: "loc1", Summary: []patterns.DeploymentMessagePerComp{dmsg("z", false, errors.New("e"), false, "")}},
		},
		"a": []patterns.DeploymentMessagePerContext{
			{Location: "loc2", Summary: []patterns.DeploymentMessagePerComp{dmsg("a", false, errors.New("e"), false, "")}},
			{Location: "loc1", Summary: []patterns.DeploymentMessagePerComp{dmsg("a", false, errors.New("e"), false, "")}},
		},
	}
	fc := tallyDeploymentOutcome(resp).FailedComponents
	if len(fc) != 4 || fc[0].Name != "a" || fc[0].Location != "loc1" || fc[1].Location != "loc2" || fc[2].Name != "z" || fc[2].Location != "loc1" || fc[3].Location != "loc2" {
		t.Fatalf("sorting mismatch: %+v", fc)
	}
}
