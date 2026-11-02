package handlers

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	patterncore "github.com/meshery/meshery/server/models/pattern/core"
	"github.com/meshery/meshery/server/models/pattern/patterns"
	"github.com/meshery/meshery/server/models/pattern/stages"
	"github.com/meshery/schemas/models/v1beta2/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// bannedDesignEventMetadataKeys may not appear at any nesting level of a
// design-related event's metadata (#22764).
var bannedDesignEventMetadataKeys = []string{
	"design", "designFile", "patternFile", "pattern_file", "components",
	"relationships", "trace", "evaluationResponse", "summary",
}

// assertNoBannedDesignEventMetadataKeys walks the JSON form of metadata - the
// form that is persisted - and fails on any banned key at any depth.
func assertNoBannedDesignEventMetadataKeys(t *testing.T, metadata map[string]interface{}) {
	t.Helper()
	data, err := json.Marshal(metadata)
	require.NoError(t, err)
	var decoded interface{}
	require.NoError(t, json.Unmarshal(data, &decoded))

	var walk func(path string, value interface{})
	walk = func(path string, value interface{}) {
		switch v := value.(type) {
		case map[string]interface{}:
			for key, child := range v {
				for _, banned := range bannedDesignEventMetadataKeys {
					if key == banned {
						t.Errorf("banned key %q present at %s", key, path)
					}
				}
				walk(path+"."+key, child)
			}
		case []interface{}:
			for i, child := range v {
				walk(fmt.Sprintf("%s[%d]", path, i), child)
			}
		}
	}
	walk("metadata", decoded)
}

// assertMetadataExcludes fails when needle appears anywhere in the JSON form
// of metadata.
func assertMetadataExcludes(t *testing.T, metadata map[string]interface{}, needle string) {
	t.Helper()
	data, err := json.Marshal(metadata)
	require.NoError(t, err)
	assert.NotContains(t, string(data), needle)
}

func provisionOutcome(success bool, err error, prerequisite bool) []patterns.DeploymentMessagePerContext {
	return []patterns.DeploymentMessagePerContext{{
		Location: "ctx",
		Summary: []patterns.DeploymentMessagePerComp{{
			Kind:           "Pod",
			CompName:       "must-not-leak",
			Success:        success,
			Message:        "must-not-leak",
			Metadata:       map[string]interface{}{"components": "must-not-leak"},
			Error:          err,
			IsPrerequisite: prerequisite,
		}},
	}}
}

func TestBuildPatternProcessEventMetadata_DeployOutcomes(t *testing.T) {
	designID := uuid.Must(uuid.NewV4())
	okID, failedID, plainErrID, prereqID := uuid.Must(uuid.NewV4()).String(), uuid.Must(uuid.NewV4()).String(), uuid.Must(uuid.NewV4()).String(), uuid.Must(uuid.NewV4()).String()
	meshkitErr := stages.ErrDependencyNotSatisfied("a", "b", false)

	response := map[string]interface{}{
		okID:       provisionOutcome(true, nil, false),
		failedID:   provisionOutcome(false, meshkitErr, false),
		plainErrID: provisionOutcome(false, fmt.Errorf("plain"), false),
		// a failed prerequisite install is fail-forward and not a component failure
		prereqID: append(provisionOutcome(true, nil, false), provisionOutcome(false, meshkitErr, true)...),
	}

	metadata := buildPatternProcessEventMetadata(response, "http://meshery.test", designID, "my-design", "0.0.2")

	assertNoBannedDesignEventMetadataKeys(t, metadata)
	assertMetadataExcludes(t, metadata, "must-not-leak")
	assert.Equal(t, designID, metadata["designId"])
	assert.Equal(t, "my-design", metadata["designName"])
	assert.Equal(t, "0.0.2", metadata["designVersion"])
	link := "http://meshery.test/extension/meshmap?mode=operator&type=view&design_id=" + designID.String()
	assert.Equal(t, link, metadata["designLink"])
	assert.Equal(t, link, metadata["viewLink"])

	assert.Equal(t, patternProcessStatusPartial, metadata["status"])
	assert.Equal(t, 2, metadata["errorCount"])
	assert.Equal(t, false, metadata["truncated"])
	errs, ok := metadata["errors"].([]patternEventError)
	require.True(t, ok)
	assert.ElementsMatch(t, []patternEventError{
		{Code: stages.ErrDependencyNotSatisfiedCode, ComponentID: failedID},
		{Code: "", ComponentID: plainErrID},
	}, errs)
}

func TestSummarizePatternProcessResponse_Status(t *testing.T) {
	id := uuid.Must(uuid.NewV4()).String()

	success := summarizePatternProcessResponse(map[string]interface{}{id: provisionOutcome(true, nil, false)})
	assert.Equal(t, patternProcessStatusSuccess, success.Status)
	assert.Equal(t, 0, success.ErrorCount)
	assert.NotNil(t, success.Errors, "errors serializes as [] rather than null")

	failed := summarizePatternProcessResponse(map[string]interface{}{id: provisionOutcome(false, nil, false)})
	assert.Equal(t, patternProcessStatusFailed, failed.Status)
	assert.Equal(t, 1, failed.ErrorCount)

	empty := summarizePatternProcessResponse(nil)
	assert.Equal(t, patternProcessStatusSuccess, empty.Status)
}

func TestSummarizePatternProcessResponse_CapsErrorsAt20(t *testing.T) {
	response := map[string]interface{}{}
	for i := 0; i < 25; i++ {
		response[uuid.Must(uuid.NewV4()).String()] = provisionOutcome(false, nil, false)
	}

	summary := summarizePatternProcessResponse(response)
	assert.Equal(t, 25, summary.ErrorCount)
	assert.Len(t, summary.Errors, maxPatternEventErrors)
	assert.True(t, summary.Truncated)

	exact := map[string]interface{}{}
	for i := 0; i < maxPatternEventErrors; i++ {
		exact[uuid.Must(uuid.NewV4()).String()] = provisionOutcome(false, nil, false)
	}
	atCap := summarizePatternProcessResponse(exact)
	assert.Len(t, atCap.Errors, maxPatternEventErrors)
	assert.False(t, atCap.Truncated)
}

func TestBuildPatternProcessEventMetadata_DryRun(t *testing.T) {
	designID := uuid.Must(uuid.NewV4())
	okComp := &component.ComponentDefinition{ID: uuid.Must(uuid.NewV4()), DisplayName: "ok", Configuration: map[string]interface{}{"spec": "must-not-leak"}}
	badComp := &component.ComponentDefinition{ID: uuid.Must(uuid.NewV4()), DisplayName: "bad"}
	noCauseComp := &component.ComponentDefinition{ID: uuid.Must(uuid.NewV4()), DisplayName: "nocause"}

	response := map[string]interface{}{
		stages.DryRunResponseKey: map[string]map[string]patterncore.DryRunResponseWrapper{
			"ok": {"ctx1": {Success: true, Component: okComp}},
			"bad": {"ctx1": {Success: false, Component: badComp, Error: &patterncore.DryRunResponse{
				Status: "Failure",
				Causes: []patterncore.DryRunFailureCause{
					{Type: "FieldValueInvalid", Message: "must-not-leak"},
					{Type: "FieldValueRequired", Message: "must-not-leak"},
				},
			}}},
			"nocause": {"ctx1": {Success: false, Component: noCauseComp, Error: &patterncore.DryRunResponse{Status: "must-not-leak"}}},
		},
	}

	metadata := buildPatternProcessEventMetadata(response, "", designID, "dry", "")

	assertNoBannedDesignEventMetadataKeys(t, metadata)
	assertMetadataExcludes(t, metadata, "must-not-leak")
	assert.NotContains(t, metadata, "designVersion")
	assert.NotContains(t, metadata, stages.DryRunResponseKey)
	assert.Equal(t, patternProcessStatusPartial, metadata["status"])
	assert.Equal(t, 3, metadata["errorCount"])
	assert.ElementsMatch(t, []patternEventError{
		{Code: "FieldValueInvalid", ComponentID: badComp.ID.String()},
		{Code: "FieldValueRequired", ComponentID: badComp.ID.String()},
		{Code: dryRunFailureCode, ComponentID: noCauseComp.ID.String()},
	}, metadata["errors"])
}

func TestPatternErrorEventMetadata_CarriesOnlyIdentityAndError(t *testing.T) {
	designID := uuid.Must(uuid.NewV4())
	err := ErrPatternDeploy(fmt.Errorf("boom"), "my-design")

	metadata := patternErrorEventMetadata(nil, designID, "my-design", "1.0.0", err)

	assertNoBannedDesignEventMetadataKeys(t, metadata)
	keys := make([]string, 0, len(metadata))
	for key := range metadata {
		keys = append(keys, key)
	}
	assert.ElementsMatch(t, []string{"designId", "designName", "designVersion", "designLink", "error"}, keys, strings.Join(keys, ","))
	assert.Equal(t, err, metadata["error"])
}
