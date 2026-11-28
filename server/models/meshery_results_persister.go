package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/schemas/models/core"
	"gorm.io/gorm"
)

type MesheryResultsPersister struct {
	DB *database.Handler
}

// ErrInvalidPerformanceResultDate indicates a malformed date filter.
var ErrInvalidPerformanceResultDate = errors.New("invalid performance result date")

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

	ServerMetrics     interface{} `json:"serverMetrics,omitempty" gorm:"type:JSONB"`
	ServerBoardConfig interface{} `json:"serverBoardConfig,omitempty" gorm:"type:JSONB"`

	TestStartTime          *time.Time         `json:"testStartTime,omitempty"`
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

// GetAllResults returns a page of results. from and to, when set, are the first
// and last days (YYYY-MM-DD) of the test start times to return.
func (mrp *MesheryResultsPersister) GetAllResults(page, pageSize uint64, from, to string, log logger.Handler) ([]byte, error) {
	if from != "" {
		if _, err := time.Parse(time.DateOnly, from); err != nil {
			return nil, fmt.Errorf("%w: from must use YYYY-MM-DD", ErrInvalidPerformanceResultDate)
		}
	}
	if to != "" {
		if _, err := time.Parse(time.DateOnly, to); err != nil {
			return nil, fmt.Errorf("%w: to must use YYYY-MM-DD", ErrInvalidPerformanceResultDate)
		}
	}

	var res []*localMesheryResultDBRepresentation
	var count int64
	// Take the date from the timestamp text to preserve its offset, as the calendar does.
	// DATETIME validates the timestamp for the query without supplying the comparison date.
	const runnerStart = "json_extract(CASE WHEN json_valid(result) THEN result END, '$.StartTime')"
	const startDay = "SUBSTR(CASE WHEN DATETIME(" + runnerStart + ") IS NOT NULL" +
		" AND SUBSTR(" + runnerStart + ", 11, 1) = 'T'" +
		" AND (SUBSTR(" + runnerStart + ", -1) = 'Z' OR SUBSTR(" + runnerStart + ", -6, 1) IN ('+', '-'))" +
		" AND " + runnerStart + " <> '0001-01-01T00:00:00Z' THEN " + runnerStart + " ELSE test_start_time END, 1, 10)"
	err := mrp.DB.Transaction(func(tx *gorm.DB) error {
		query := tx.Table("meshery_results")
		if from != "" {
			query = query.Where(startDay+" >= ?", from)
		}
		if to != "" {
			query = query.Where(startDay+" <= ?", to)
		}
		if err := query.Count(&count).Error; err != nil {
			return err
		}
		return Paginate(uint(page), uint(pageSize))(query).Find(&res).Error
	})
	if err != nil {
		return nil, err
	}

	results := convertLocalRepresentationSliceToMesheryResultSlice(res, log)
	if results == nil {
		results = []*MesheryResult{}
	}
	for _, result := range results {
		raw, _ := result.Result["StartTime"].(string)
		if start, err := time.Parse(time.RFC3339Nano, raw); err == nil && !start.IsZero() {
			result.TestStartTime = &start
		}
	}
	resultPage := &MesheryResultPage{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: int(count),
		Results:    results,
	}

	return marshalMesheryResultsPage(resultPage), nil
}

func (mrp *MesheryResultsPersister) GetResult(key core.Uuid, log logger.Handler) (*MesheryResult, error) {
	var lres localMesheryResultDBRepresentation

	err := mrp.DB.Table("meshery_results").Where("id = ?", key).First(&lres).Error
	if err != nil {
		return nil, err
	}
	res, err := convertLocalRepresentationToMesheryResult(&lres, log)
	return res, err
}

func (mrp *MesheryResultsPersister) WriteResult(key core.Uuid, result []byte) error {
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
		r, _ := convertLocalRepresentationToMesheryResult(val, log)
		if r != nil {
			res = append(res, r)
		}
	}

	return
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
		ServerMetrics:      local.ServerMetrics,
		ServerBoardConfig:  local.ServerMetrics,
		TestStartTime:      local.TestStartTime,
	}

	return res, nil
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
