package models

import (
	"reflect"
	"testing"
)

func TestMesheryResultConversionsPreserveServerFields(t *testing.T) {
	metrics := map[string]interface{}{"source": "metrics"}
	boardConfig := map[string]interface{}{"source": "board"}
	result := &MesheryResult{
		Result:            map[string]interface{}{"status": "complete"},
		ServerMetrics:     metrics,
		ServerBoardConfig: boardConfig,
	}

	local := convertMesheryResultToLocalRepresentation(result)
	if !reflect.DeepEqual(local.ServerMetrics, metrics) {
		t.Fatalf("write conversion server metrics = %#v, want %#v", local.ServerMetrics, metrics)
	}
	if !reflect.DeepEqual(local.ServerBoardConfig, boardConfig) {
		t.Fatalf("write conversion server board config = %#v, want %#v", local.ServerBoardConfig, boardConfig)
	}

	roundTripped := convertLocalRepresentationToMesheryResult(local, nil)
	if !reflect.DeepEqual(roundTripped.ServerMetrics, metrics) {
		t.Fatalf("read conversion server metrics = %#v, want %#v", roundTripped.ServerMetrics, metrics)
	}
	if !reflect.DeepEqual(roundTripped.ServerBoardConfig, boardConfig) {
		t.Fatalf("read conversion server board config = %#v, want %#v", roundTripped.ServerBoardConfig, boardConfig)
	}
}
