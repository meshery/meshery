package models

import (
	"encoding/json"
	"testing"

	mkerrors "github.com/meshery/meshkit/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMesheryResultsPersister_NilReceiverSafety(t *testing.T) {
	var mrp *MesheryResultsPersister

	assert.NotPanics(t, func() {
		res, err := mrp.GetResults(0, 10, "profile-id", nil)
		assert.Nil(t, res)
		assert.Equal(t, ErrDBConnection, err)
	})

	assert.NotPanics(t, func() {
		res, err := mrp.GetAllResults(0, 10, nil)
		assert.Nil(t, res)
		assert.Equal(t, ErrDBConnection, err)
	})

	assert.NotPanics(t, func() {
		res, err := mrp.GetResult(mustUUID(t), nil)
		assert.Nil(t, res)
		assert.Equal(t, ErrDBConnection, err)
	})

	assert.NotPanics(t, func() {
		err := mrp.WriteResult(mustUUID(t), []byte("{}"))
		assert.Equal(t, ErrDBConnection, err)
	})

	mrpWithNilDB := &MesheryResultsPersister{DB: nil}
	assert.NotPanics(t, func() {
		res, err := mrpWithNilDB.GetResult(mustUUID(t), nil)
		assert.Nil(t, res)
		assert.Equal(t, ErrDBConnection, err)
	})
}

func TestMesheryResultsPersister_GetResult_NotFound(t *testing.T) {
	db := newMigratedDB(t)
	mrp := &MesheryResultsPersister{DB: db}

	// 1. Populate DB with an existing result A to ensure lookup doesn't return
	// the wrong row (guarding the original Find(...).Where(...) bug in #22110).
	existingID := mustUUID(t)
	existingResult := &MesheryResult{
		ID:     existingID,
		Name:   "existing-result",
		Result: map[string]interface{}{"status": "ok"},
	}
	existingBytes, err := json.Marshal(existingResult)
	require.NoError(t, err)
	require.NoError(t, mrp.WriteResult(existingID, existingBytes))

	// Verify existingID lookup succeeds
	found, err := mrp.GetResult(existingID, nil)
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, existingID, found.ID)

	// 2. Querying a non-existent ID B in a populated table must return ErrResultNotFound
	nonExistentID := mustUUID(t)
	result, err := mrp.GetResult(nonExistentID, nil)

	assert.Nil(t, result)
	require.Error(t, err)
	assert.Equal(t, ErrResultNotFoundCode, mkerrors.GetCode(err), "error code should be ErrResultNotFoundCode")
}

func TestMesheryResultsPersister_PreservesServerBoardConfigAndMetrics(t *testing.T) {
	db := newMigratedDB(t)
	mrp := &MesheryResultsPersister{DB: db}

	id := mustUUID(t)
	profileID := mustUUID(t)

	serverMetrics := map[string]interface{}{
		"metric_source": "prometheus",
		"queries":       []interface{}{"node_cpu_seconds_total"},
	}
	serverBoardConfig := map[string]interface{}{
		"board_id": "grafana-dashboard-42",
		"panels":   []interface{}{"panel-1", "panel-2"},
	}

	result := &MesheryResult{
		ID:                 id,
		Name:               "load-test-benchmark",
		Mesh:               "istio",
		PerformanceProfile: &profileID,
		Result: map[string]interface{}{
			"status": "success",
			"qps":    float64(500),
		},
		ServerMetrics:     serverMetrics,
		ServerBoardConfig: serverBoardConfig,
	}

	data, err := json.Marshal(result)
	require.NoError(t, err)

	err = mrp.WriteResult(id, data)
	require.NoError(t, err)

	retrieved, err := mrp.GetResult(id, nil)
	require.NoError(t, err)
	require.NotNil(t, retrieved)

	assert.Equal(t, id, retrieved.ID)
	assert.Equal(t, "load-test-benchmark", retrieved.Name)
	assert.Equal(t, "istio", retrieved.Mesh)

	// Verify that ServerBoardConfig is NOT overwritten by ServerMetrics
	metricsBytes, err := json.Marshal(retrieved.ServerMetrics)
	require.NoError(t, err)
	boardBytes, err := json.Marshal(retrieved.ServerBoardConfig)
	require.NoError(t, err)

	assert.NotEqual(t, string(metricsBytes), string(boardBytes), "ServerMetrics and ServerBoardConfig must remain distinct")
	assert.Contains(t, string(metricsBytes), "prometheus")
	assert.Contains(t, string(boardBytes), "grafana-dashboard-42")
}

func TestMesheryResultsPersister_GetAllResults_And_GetResults(t *testing.T) {
	db := newMigratedDB(t)
	mrp := &MesheryResultsPersister{DB: db}

	profileID1 := mustUUID(t)
	profileID2 := mustUUID(t)

	id1 := mustUUID(t)
	res1 := &MesheryResult{
		ID:                 id1,
		Name:               "run-profile-1",
		PerformanceProfile: &profileID1,
		Result:             map[string]interface{}{"run": 1},
	}
	b1, err := json.Marshal(res1)
	require.NoError(t, err)
	require.NoError(t, mrp.WriteResult(id1, b1))

	id2 := mustUUID(t)
	res2 := &MesheryResult{
		ID:                 id2,
		Name:               "run-profile-2",
		PerformanceProfile: &profileID2,
		Result:             map[string]interface{}{"run": 2},
	}
	b2, err := json.Marshal(res2)
	require.NoError(t, err)
	require.NoError(t, mrp.WriteResult(id2, b2))

	// Test GetAllResults
	allBytes, err := mrp.GetAllResults(0, 10, nil)
	require.NoError(t, err)
	var allPage MesheryResultPage
	err = json.Unmarshal(allBytes, &allPage)
	require.NoError(t, err)
	assert.Equal(t, 2, allPage.TotalCount)
	assert.Len(t, allPage.Results, 2)

	// Test GetResults for profileID1
	profile1Bytes, err := mrp.GetResults(0, 10, profileID1.String(), nil)
	require.NoError(t, err)
	var p1Page MesheryResultPage
	err = json.Unmarshal(profile1Bytes, &p1Page)
	require.NoError(t, err)
	assert.Equal(t, 1, p1Page.TotalCount)
	require.Len(t, p1Page.Results, 1)
	assert.Equal(t, id1, p1Page.Results[0].ID)
}

func TestMesheryResultsPersister_WriteResult_InvalidJSON(t *testing.T) {
	db := newMigratedDB(t)
	mrp := &MesheryResultsPersister{DB: db}

	err := mrp.WriteResult(mustUUID(t), []byte("invalid-json"))
	assert.Error(t, err)
}

func TestMesheryResultsPersister_StringFieldPreservesType(t *testing.T) {
	db := newMigratedDB(t)
	mrp := &MesheryResultsPersister{DB: db}

	id := mustUUID(t)
	result := &MesheryResult{
		ID:                id,
		Name:              "string-type-test",
		Result:            map[string]interface{}{"ok": true},
		ServerBoardConfig: "null", // string "null", must not become nil or map
		ServerMetrics:     `{"json":"string"}`,
	}

	data, err := json.Marshal(result)
	require.NoError(t, err)

	err = mrp.WriteResult(id, data)
	require.NoError(t, err)

	retrieved, err := mrp.GetResult(id, nil)
	require.NoError(t, err)
	require.NotNil(t, retrieved)

	assert.Equal(t, "null", retrieved.ServerBoardConfig)
	assert.Equal(t, `{"json":"string"}`, retrieved.ServerMetrics)
}
