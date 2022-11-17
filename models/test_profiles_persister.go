package models

import (
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshkit/database"
	SMP "github.com/layer5io/service-mesh-performance/spec"
)

// TestProfilesPersister assists with persisting session in store
type TestProfilesPersister struct {
	DB *database.Handler
}

type PerformanceTestConfig struct {
	ID                         uuid.UUID
	PerformanceTestConfigBytes []byte
	UpdatedAt                  time.Time
}

// UserTestProfiles - represents a page of user test configs
type UserTestProfiles struct {
	Page        uint64                       `json:"page"`
	PageSize    uint64                       `json:"page_size"`
	TotalCount  int                          `json:"total_count"`
	TestConfigs []*SMP.PerformanceTestConfig `json:"test_configs"`
}

// GetTestConfigs - gets result for the page and pageSize
func (s *TestProfilesPersister) GetTestConfigs(page, pageSize uint64) ([]byte, error) {
	if s.DB == nil {
		return nil, ErrDBConnection
	}
	order := "updated_at desc"
	query := s.DB.Order(order)
	total := int64(0)
	s.DB.Model(&PerformanceTestConfig{}).Count(&total)
	testConfigs := []*SMP.PerformanceTestConfig{}
	var p []*PerformanceTestConfig
	Paginate(uint(page), uint(pageSize))(query).Find(&p)
	for _, config := range p {
		var testConfig SMP.PerformanceTestConfig
		err := json.Unmarshal(config.PerformanceTestConfigBytes, &testConfig)
		if err == nil {
			testConfigs = append(testConfigs, &testConfig)
		}
	}
	bd, err := json.Marshal(&UserTestProfiles{
		Page:        page,
		PageSize:    pageSize,
		TotalCount:  int(total),
		TestConfigs: testConfigs,
	})
	if err != nil {
		return nil, ErrMarshal(err, "result data")
	}

	return bd, nil
}

// GetTestConfig - gets result for a specific key
func (s *TestProfilesPersister) GetTestConfig(key uuid.UUID) (*SMP.PerformanceTestConfig, error) {
	if s.DB == nil {
		return nil, ErrDBConnection
	}
	testConfig := &SMP.PerformanceTestConfig{}
	var u PerformanceTestConfig
	err := s.DB.Model(&PerformanceTestConfig{}).Where("id = ?", key).First(&u).Error
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(u.PerformanceTestConfigBytes, &testConfig)
	if err != nil {
		return nil, err
	}
	return testConfig, nil
}

// DeleteTestConfig - delete result for a specific key
func (s *TestProfilesPersister) DeleteTestConfig(key uuid.UUID) error {
	if s.DB == nil {
		return ErrDBConnection
	}

	return s.DB.Model(&PerformanceTestConfig{}).Where("id = ?", key).Delete(&PerformanceTestConfig{}).Error
}

// WriteTestConfig persists the result
func (s *TestProfilesPersister) WriteTestConfig(key uuid.UUID, result []byte) error {
	if s.DB == nil {
		return ErrDBConnection
	}

	if result == nil {
		return ErrResultData()
	}
	var p PerformanceTestConfig
	if err := s.DB.Model(&PerformanceTestConfig{}).Where("id = ?", key).First(&p).Error; err == nil {
		return s.UpdateTestConfig(key, p)
	}
	p.ID = key
	p.PerformanceTestConfigBytes = result
	return s.DB.Model(&PerformanceTestConfig{}).Create(&p).Error
}

func (s *TestProfilesPersister) UpdateTestConfig(key uuid.UUID, p PerformanceTestConfig) error {
	return s.DB.Model(&PerformanceTestConfig{}).Where("id = ?", key).UpdateColumns(p).Error
}
