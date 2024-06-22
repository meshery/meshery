package models

import (
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshkit/database"
	"github.com/layer5io/meshkit/logger"
)

type MesheryResultsPersister struct {
	DB *database.Handler
}

// MesheryResultPage - represents a page of meshery results
type MesheryResultPage struct {
	Page       uint64           `json:"page"`
	PageSize   uint64           `json:"page_size"`
	TotalCount int              `json:"total_count"`
	Results    []*MesheryResult `json:"results"`
}

type localMesheryResultDBRepresentation struct {
	ID                 uuid.UUID  `json:"meshery_id,omitempty"`
	Name               string     `json:"name,omitempty"`
	Mesh               string     `json:"mesh,omitempty"`
	PerformanceProfile *uuid.UUID `json:"performance_profile,omitempty"`
	Result             []byte     `json:"runner_results,omitempty" gorm:"type:JSONB"`

	ServerMetrics     interface{} `json:"server_metrics,omitempty" gorm:"type:JSONB"`
	ServerBoardConfig interface{} `json:"server_board_config,omitempty" gorm:"type:JSONB"`

	TestStartTime          *time.Time         `json:"test_start_time,omitempty"`
	PerformanceProfileInfo PerformanceProfile `json:"-" gorm:"constraint:OnDelete:SET NULL;foreignKey:PerformanceProfile"`
}

func (mrp *MesheryResultsPersister) GetResults(page, pageSize uint64, profileID string, log logger.Handler) ([]byte, error) {
	var res []*localMesheryResultDBRepresentation
	var count int64
	query := mrp.DB.Where("performance_profile = ?", profileID)

	err := query.Table("meshery_results").Count(&count).Error
	if err != nil {
		return nil, err
	}
	err = Paginate(uint(page), uint(pageSize))(query).Find(&res).Error

	resultPage := &MesheryResultPage{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: int(count),
		Results:    convertLocalRepresentationSliceToMesheryResultSlice(res, log),
	}

	return marshalMesheryResultsPage(resultPage), err
}

func (mrp *MesheryResultsPersister) GetAllResults(page, pageSize uint64, log logger.Handler) ([]byte, error) {
	var res []*localMesheryResultDBRepresentation
	var count int64
	query := mrp.DB.Table("meshery_results")

	err := query.Count(&count).Error
	if err != nil {
		return nil, err
	}
	err = Paginate(uint(page), uint(pageSize))(query).Find(&res).Error

	resultPage := &MesheryResultPage{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: int(count),
		Results:    convertLocalRepresentationSliceToMesheryResultSlice(res, log),
	}

	return marshalMesheryResultsPage(resultPage), err
}

func (mrp *MesheryResultsPersister) GetResult(key uuid.UUID, log logger.Handler) (*MesheryResult, error) {
	var lres localMesheryResultDBRepresentation

	err := mrp.DB.Table("meshery_results").Find(&lres).Where("id = ?", key).Error
	res := convertLocalRepresentationToMesheryResult(&lres, log)
	return res, err
}

func (mrp *MesheryResultsPersister) WriteResult(key uuid.UUID, result []byte) error {
	var data MesheryResult
	if err := json.Unmarshal(result, &data); err != nil {
		return err
	}

	data.ID = key

	t := time.Now()
	data.TestStartTime = &t
	return mrp.DB.Table("meshery_results").Save(convertMesheryResultToLocalRepresentation(&data)).Error
}

func marshalMesheryResultsPage(mrp *MesheryResultPage) []byte {
	res, _ := json.Marshal(mrp)

	return res
}

func convertLocalRepresentationSliceToMesheryResultSlice(local []*localMesheryResultDBRepresentation, log logger.Handler) (res []*MesheryResult) {
	for _, val := range local {
		res = append(res, convertLocalRepresentationToMesheryResult(val, log))
	}

	return
}

func convertLocalRepresentationToMesheryResult(local *localMesheryResultDBRepresentation, log logger.Handler) *MesheryResult {
	var jsonmap map[string]interface{}
	if err := json.Unmarshal(local.Result, &jsonmap); err != nil {
		err = ErrUnmarshal(err, "MesheryResult")
		log.Error(err)
		return nil
	}

	res := &MesheryResult{
		ID:                 local.ID,
		Name:               local.Name,
		Mesh:               local.Mesh,
		PerformanceProfile: local.PerformanceProfile,
		Result:             jsonmap,
		ServerMetrics:      local.ServerMetrics,
		ServerBoardConfig:  local.ServerMetrics,
		TestStartTime:      local.TestStartTime,
	}

	return res
}

func convertMesheryResultToLocalRepresentation(mr *MesheryResult) *localMesheryResultDBRepresentation {
	byt, _ := json.Marshal(mr.Result)

	res := &localMesheryResultDBRepresentation{
		ID:                 mr.ID,
		Name:               mr.Name,
		Mesh:               mr.Mesh,
		PerformanceProfile: mr.PerformanceProfile,
		Result:             byt,
		ServerMetrics:      mr.ServerMetrics,
		ServerBoardConfig:  mr.ServerMetrics,
		TestStartTime:      mr.TestStartTime,
	}

	return res
}
