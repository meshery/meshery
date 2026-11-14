package models

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/schemas/models/core"
	"gorm.io/gorm"
)

type MesheryResultsPersister struct {
	DB *database.Handler
}

// MesheryResultPage - represents a page of meshery results
type MesheryResultPage struct {
	Page       uint64           `json:"page"`
	PageSize   uint64           `json:"pageSize"`
	TotalCount int              `json:"totalCount"`
	Results    []*MesheryResult `json:"results"`
}

type localMesheryResultDBRepresentation struct {
	ID                 core.Uuid  `json:"mesheryId,omitempty"`
	Name               string     `json:"name,omitempty"`
	Mesh               string     `json:"mesh,omitempty"`
	PerformanceProfile *core.Uuid `json:"performanceProfile,omitempty"`
	Result             []byte     `json:"runnerResults,omitempty" gorm:"type:JSONB"`

	ServerMetrics     []byte `json:"serverMetrics,omitempty" gorm:"type:JSONB"`
	ServerBoardConfig []byte `json:"serverBoardConfig,omitempty" gorm:"type:JSONB"`

	TestStartTime          *time.Time         `json:"testStartTime,omitempty"`
	PerformanceProfileInfo PerformanceProfile `json:"-" gorm:"constraint:OnDelete:SET NULL;foreignKey:PerformanceProfile"`
}

func (mrp *MesheryResultsPersister) GetResults(page, pageSize uint64, profileID string, log logger.Handler) ([]byte, error) {
	if mrp == nil || mrp.DB == nil {
		return nil, ErrDBConnection
	}

	var res []*localMesheryResultDBRepresentation
	var count int64
	query := mrp.DB.Where("performance_profile = ?", profileID)

	err := query.Table("meshery_results").Count(&count).Error
	if err != nil {
		return nil, err
	}
	err = Paginate(uint(page), uint(pageSize))(query).Find(&res).Error
	if err != nil {
		return nil, err
	}

	resultPage := &MesheryResultPage{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: int(count),
		Results:    convertLocalRepresentationSliceToMesheryResultSlice(res, log),
	}

	return marshalMesheryResultsPage(resultPage), nil
}

func (mrp *MesheryResultsPersister) GetAllResults(page, pageSize uint64, log logger.Handler) ([]byte, error) {
	if mrp == nil || mrp.DB == nil {
		return nil, ErrDBConnection
	}

	var res []*localMesheryResultDBRepresentation
	var count int64
	query := mrp.DB.Table("meshery_results")

	err := query.Count(&count).Error
	if err != nil {
		return nil, err
	}
	err = Paginate(uint(page), uint(pageSize))(query).Find(&res).Error
	if err != nil {
		return nil, err
	}

	resultPage := &MesheryResultPage{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: int(count),
		Results:    convertLocalRepresentationSliceToMesheryResultSlice(res, log),
	}

	return marshalMesheryResultsPage(resultPage), nil
}

func (mrp *MesheryResultsPersister) GetResult(key core.Uuid, log logger.Handler) (*MesheryResult, error) {
	if mrp == nil || mrp.DB == nil {
		return nil, ErrDBConnection
	}

	var lres localMesheryResultDBRepresentation

	err := mrp.DB.Table("meshery_results").Where("id = ?", key).First(&lres).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrResultNotFound(err)
		}
		return nil, err
	}
	res, err := convertLocalRepresentationToMesheryResult(&lres, log)
	return res, err
}

func (mrp *MesheryResultsPersister) WriteResult(key core.Uuid, result []byte) error {
	if mrp == nil || mrp.DB == nil {
		return ErrDBConnection
	}

	var data MesheryResult
	if err := json.Unmarshal(result, &data); err != nil {
		return err
	}

	data.ID = key

	t := time.Now()
	data.TestStartTime = &t

	localRep, err := convertMesheryResultToLocalRepresentation(&data)
	if err != nil {
		return err
	}

	return mrp.DB.Table("meshery_results").Save(localRep).Error
}

func marshalMesheryResultsPage(mrp *MesheryResultPage) []byte {
	res, _ := json.Marshal(mrp)

	return res
}

func convertLocalRepresentationSliceToMesheryResultSlice(local []*localMesheryResultDBRepresentation, log logger.Handler) (res []*MesheryResult) {
	for _, val := range local {
		r, _ := convertLocalRepresentationToMesheryResult(val, log)
		if r != nil {
			res = append(res, r)
		}
	}

	return
}

func marshalJSONField(v interface{}) ([]byte, error) {
	if v == nil {
		return nil, nil
	}
	if b, ok := v.([]byte); ok {
		return b, nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func unmarshalJSONField(b []byte) interface{} {
	if len(b) == 0 {
		return nil
	}
	var out interface{}
	if err := json.Unmarshal(b, &out); err == nil {
		return out
	}
	return string(b)
}

func convertLocalRepresentationToMesheryResult(local *localMesheryResultDBRepresentation, log logger.Handler) (*MesheryResult, error) {
	var jsonmap map[string]interface{}
	if err := json.Unmarshal(local.Result, &jsonmap); err != nil {
		err = ErrUnmarshal(err, "MesheryResult")
		if log != nil {
			log.Error(err)
		}
		return nil, err
	}

	res := &MesheryResult{
		ID:                 local.ID,
		Name:               local.Name,
		Mesh:               local.Mesh,
		PerformanceProfile: local.PerformanceProfile,
		Result:             jsonmap,
		ServerMetrics:      unmarshalJSONField(local.ServerMetrics),
		ServerBoardConfig:  unmarshalJSONField(local.ServerBoardConfig),
		TestStartTime:      local.TestStartTime,
	}

	return res, nil
}

func convertMesheryResultToLocalRepresentation(mr *MesheryResult) (*localMesheryResultDBRepresentation, error) {
	byt, err := json.Marshal(mr.Result)
	if err != nil {
		return nil, err
	}

	serverMetricsBytes, err := marshalJSONField(mr.ServerMetrics)
	if err != nil {
		return nil, err
	}

	serverBoardConfigBytes, err := marshalJSONField(mr.ServerBoardConfig)
	if err != nil {
		return nil, err
	}

	res := &localMesheryResultDBRepresentation{
		ID:                 mr.ID,
		Name:               mr.Name,
		Mesh:               mr.Mesh,
		PerformanceProfile: mr.PerformanceProfile,
		Result:             byt,
		ServerMetrics:      serverMetricsBytes,
		ServerBoardConfig:  serverBoardConfigBytes,
		TestStartTime:      mr.TestStartTime,
	}

	return res, nil
}
