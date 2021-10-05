package models

import (
	"encoding/json"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshkit/database"
)

// BitCaskSmiResultsPersister assists with persisting session in a Bitcask store
type BitCaskSmiResultsPersister struct {
	DB *database.Handler
}

type SmiResultWithId struct {
	ID        uuid.UUID
	SmiResult `gorm:"embedded"`
}

// SmiResultPage - represents a page of meshery results
type SmiResultPage struct {
	Page       uint64       `json:"page"`
	PageSize   uint64       `json:"page_size"`
	TotalCount int          `json:"total_count"`
	Results    []*SmiResult `json:"results"`
}

// GetSmiResults - gets result for the page and pageSize
func (s *BitCaskSmiResultsPersister) GetResults(page, pageSize uint64) ([]byte, error) {
	if s.DB == nil {
		return nil, ErrDBConnection
	}

	total := int64(0)
	s.DB.Model(&SmiResultWithId{}).Count(&total)
	results := []*SmiResult{}
	order := "updated_at desc"
	query := s.DB.Order(order)
	Paginate(uint(page), uint(pageSize))(query).Find(&results)
	bd, err := json.Marshal(&SmiResultPage{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: int(total),
		Results:    results,
	})
	if err != nil {
		obj := "result data"
		//err = errors.Wrapf(err, "Unable to marshal result data.")
		return nil, ErrMarshal(err, obj)
	}

	return bd, nil
}

// WriteSmiResult persists the result
func (s *BitCaskSmiResultsPersister) WriteResult(key uuid.UUID, result []byte) error {
	if s.DB == nil {
		return ErrDBConnection
	}

	if result == nil {
		return ErrResultData()
	}
	var r SmiResultWithId
	if err := s.DB.Model(&PerformanceTestConfig{}).Where("id = ?", key).First(&r).Error; err == nil {
		err = s.DeleteResult(key)
		if err != nil {
			return err
		}
	}
	err := json.Unmarshal(result, &r.SmiResult)
	if err != nil {
		return err
	}
	r.ID = key
	return s.DB.Model(&PerformanceTestConfig{}).Create(&r).Error
}

func (s *BitCaskSmiResultsPersister) DeleteResult(key uuid.UUID) error {
	if s.DB == nil {
		return ErrDBConnection
	}
	return s.DB.Model(&PerformanceTestConfig{}).Where("id = ?", key).Delete(&PerformanceTestConfig{}).Error
}

// CloseSmiResultPersister closes the badger store
// func (s *BitCaskSmiResultsPersister) CloseResultPersister() {
// 	if s.db == nil {
// 		return
// 	}
// 	_ = s.db.Close()
// }
