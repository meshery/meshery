package models

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/schemas/models/core"
	"gorm.io/gorm"
)

// MesheryPatternVersion is one saved revision of a design's file.
type MesheryPatternVersion struct {
	ID          core.Uuid `json:"id" gorm:"primaryKey"`
	DesignID    core.Uuid `json:"designId" gorm:"index;not null"`
	Version     string    `json:"version"`
	PatternFile string    `json:"patternFile"`
	CreatedAt   time.Time `json:"createdAt"`
}

func recordPatternVersion(tx *gorm.DB, p *MesheryPattern, version string) error {
	id, err := uuid.NewV4()
	if err != nil {
		return ErrGenerateUUID(err)
	}
	revision := MesheryPatternVersion{
		ID:          id,
		DesignID:    *p.ID,
		Version:     version,
		PatternFile: p.PatternFile,
	}
	if err := tx.Create(&revision).Error; err != nil {
		return ErrDBCreate(err)
	}
	return nil
}
