package models

import (
	"encoding/json"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshkit/database"
)

// SmiResultsPersister assists with persisting session in store
type SMIResultsPersister struct {
	DB *database.Handler
}

type SmiResultWithID struct {
	ID        uuid.UUID
	SmiResult `gorm:"embedded"`
}

// SmiResultPage - represents a page of meshery results
type SmiResultPage struct {
	Page       uint64             `json:"page"`
	PageSize   uint64             `json:"page_size"`
	TotalCount int                `json:"total_count"`
	Results    []*SmiResultWithID `json:"results"`
}

// GetSmiResults - gets result for the page and pageSize
func (s *SMIResultsPersister) GetResults(page, pageSize uint64) ([]byte, error) {
	if s.DB == nil {
		return nil, ErrDBConnection
	}

	total := int64(0)
	s.DB.Model(&SmiResultWithID{}).Count(&total)
	results := []*SmiResultWithID{}
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
func (s *SMIResultsPersister) WriteResult(key uuid.UUID, result []byte) error {
	if s.DB == nil {
		return ErrDBConnection
	}

	if result == nil {
		return ErrResultData()
	}
	var r SmiResultWithID
	if err := s.DB.Model(&SmiResultWithID{}).Where("id = ?", key).First(&r).Error; err == nil {
		return s.UpdateResult(key, r)
	}
	err := json.Unmarshal(result, &r.SmiResult)
	if err != nil {
		return err
	}
	r.ID = key
	return s.DB.Model(&SmiResultWithID{}).Create(&r).Error
}

func (s *SMIResultsPersister) DeleteResult(key uuid.UUID) error {
	if s.DB == nil {
		return ErrDBConnection
	}
	return s.DB.Model(&SmiResultWithID{}).Where("id = ?", key).Delete(&SmiResultWithID{}).Error
}

func (s *SMIResultsPersister) GetResult(page, pageSize uint64, key uuid.UUID) ([]byte, error) {
	if s.DB == nil {
		return nil, ErrDBConnection
	}
	var result SmiResult
	err := s.DB.First(&result, key).Error
	if err != nil {
		return nil, err
	}
	bd, _ := json.Marshal(result)

	return bd, nil
}

func (s *SMIResultsPersister) UpdateResult(key uuid.UUID, res SmiResultWithID) error {
	return s.DB.Model(&SmiResultWithID{}).Where("id = ?", key).UpdateColumns(res).Error
}
