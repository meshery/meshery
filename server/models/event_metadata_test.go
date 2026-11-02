package models

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/models/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDesignEventIdentity(t *testing.T) {
	id := uuid.Must(uuid.NewV4())

	withVersion := DesignEventIdentity("http://meshery.test", id, "d", "1.0.0")
	assert.Equal(t, map[string]interface{}{
		"designId":      id,
		"designName":    "d",
		"designVersion": "1.0.0",
		"designLink":    "http://meshery.test/extension/meshmap?mode=operator&type=view&design_id=" + id.String(),
	}, withVersion)

	withoutVersion := DesignEventIdentity("", id, "d", "")
	assert.NotContains(t, withoutVersion, "designVersion")
}

func TestGuardEventMetadataSize_OversizedKeepsIdentityAndHistoryTitle(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	metadata := DesignEventIdentity("http://meshery.test", id, "big", "2.0.0")
	metadata["historyTitle"] = "3 changes made at version 2.0.0"
	metadata["components"] = strings.Repeat("x", MaxEventMetadataBytes)
	metadata["viewLink"] = "kept only when small"

	event := events.Event{Category: "pattern", Action: "deploy", Metadata: metadata}
	guarded := GuardEventMetadataSize(event, nil)

	assert.Equal(t, map[string]interface{}{
		"designId":      id,
		"designName":    "big",
		"designVersion": "2.0.0",
		"designLink":    metadata["designLink"],
		"historyTitle":  "3 changes made at version 2.0.0",
	}, guarded.Metadata)
	// The caller's map is shared with the published event and stays intact.
	assert.Contains(t, event.Metadata, "components")
	assert.Contains(t, metadata, "components")
}

func TestGuardEventMetadataSize_SmallUntouched(t *testing.T) {
	metadata := map[string]interface{}{
		"designId": uuid.Must(uuid.NewV4()),
		"status":   "success",
		"errors":   []string{},
	}
	event := events.Event{Category: "pattern", Action: "deploy", Metadata: metadata}

	guarded := GuardEventMetadataSize(event, nil)
	assert.Equal(t, metadata, guarded.Metadata)

	assert.Nil(t, GuardEventMetadataSize(events.Event{}, nil).Metadata)
}

func TestGuardEventMetadataSize_Boundary(t *testing.T) {
	// {"k":"<n x>"} is n+8 bytes.
	atLimit := map[string]interface{}{"k": strings.Repeat("x", MaxEventMetadataBytes-8)}
	data, err := json.Marshal(atLimit)
	require.NoError(t, err)
	require.Len(t, data, MaxEventMetadataBytes)
	assert.Equal(t, atLimit, GuardEventMetadataSize(events.Event{Metadata: atLimit}, nil).Metadata)

	overLimit := map[string]interface{}{"k": strings.Repeat("x", MaxEventMetadataBytes-7)}
	assert.Empty(t, GuardEventMetadataSize(events.Event{Metadata: overLimit}, nil).Metadata)
}
