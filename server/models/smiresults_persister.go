package models

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/meshery/meshkit/database"
	"github.com/meshery/schemas/models/core"
	"gorm.io/gorm"
)

// SMIResultsPersister assists with persisting SMI results in store
type SMIResultsPersister struct {
	DB *database.Handler
}

type SmiResultWithID struct {
	ID        core.Uuid
	SmiResult `gorm:"embedded"`
}

// SmiResultPage - represents a page of meshery results
type SmiResultPage struct {
	Page       uint64             `json:"page"`
	PageSize   uint64             `json:"pageSize"`
	TotalCount int                `json:"totalCount"`
	Results    []*SmiResultWithID `json:"results"`
}

// GetResults - gets result for the page and pageSize
func (s *SMIResultsPersister) GetResults(page, pageSize uint64) ([]byte, error) {
	if s.DB == nil {
		return nil, ErrDBConnection
	}

	total := int64(0)
	if err := s.DB.Model(&SmiResultWithID{}).Count(&total).Error; err != nil {
		return nil, ErrDBRead(err)
	}

	results := []*SmiResultWithID{}
	order := "date desc"
	query := s.DB.Order(order)
	if err := Paginate(uint(page), uint(pageSize))(query).Find(&results).Error; err != nil {
		return nil, ErrDBRead(err)
	}

	bd, err := json.Marshal(&SmiResultPage{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: int(total),
		Results:    results,
	})
	if err != nil {
		return nil, ErrMarshal(err, "result data")
	}

	return bd, nil
}

// WriteResult persists the result
func (s *SMIResultsPersister) WriteResult(key core.Uuid, result []byte) error {
	if s.DB == nil {
		return ErrDBConnection
	}

	if result == nil {
		return ErrResultData()
	}
	var r SmiResultWithID
	if err := s.DB.Model(&SmiResultWithID{}).Where("id = ?", key).First(&r).Error; err == nil {
		if err := json.Unmarshal(result, &r.SmiResult); err != nil {
			return ErrUnmarshal(err, "smi result")
		}
		return s.UpdateResult(key, r)
	}
	if err := json.Unmarshal(result, &r.SmiResult); err != nil {
		return ErrUnmarshal(err, "smi result")
	}
	r.ID = key
	if err := s.DB.Model(&SmiResultWithID{}).Create(&r).Error; err != nil {
		return ErrDBCreate(err)
	}
	return nil
}

func (s *SMIResultsPersister) DeleteResult(key core.Uuid) error {
	if s.DB == nil {
		return ErrDBConnection
	}
	res := s.DB.Model(&SmiResultWithID{}).Where("id = ?", key).Delete(&SmiResultWithID{})
	if res.Error != nil {
		return ErrDBDelete(res.Error, key.String())
	}
	if res.RowsAffected == 0 {
		return ErrResultNotFound(fmt.Errorf("smi result %s not found", key.String()))
	}
	return nil
}

func (s *SMIResultsPersister) GetResult(key core.Uuid) ([]byte, error) {
	if s.DB == nil {
		return nil, ErrDBConnection
	}
	var result SmiResultWithID
	err := s.DB.Model(&SmiResultWithID{}).Where("id = ?", key).First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrResultNotFound(fmt.Errorf("smi result %s not found", key.String()))
		}
		return nil, ErrDBRead(err)
	}
	result.SmiResult.ID = result.ID
	bd, err := json.Marshal(result.SmiResult)
	if err != nil {
		return nil, ErrMarshal(err, "smi result data")
	}

	return bd, nil
}

func (s *SMIResultsPersister) UpdateResult(key core.Uuid, res SmiResultWithID) error {
	if s.DB == nil {
		return ErrDBConnection
	}
	if err := s.DB.Model(&SmiResultWithID{}).Where("id = ?", key).UpdateColumns(res).Error; err != nil {
		return ErrDBPut(err)
	}
	return nil
}
