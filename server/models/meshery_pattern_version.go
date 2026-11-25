package models

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/schemas/models/core"
	"gorm.io/gorm"
)

// MesheryPatternVersion is one saved revision of a design's file. A revision is
// written in the same transaction as the design it belongs to, so the two can
// never disagree.
//
// Every path that deletes a design deletes its revisions in that same
// transaction: the two persister deletes, deletePublishedPatterns (seed and
// reseed), and Cleanup, which drops this table on shutdown. A new
// design-deletion path must do the same, or it leaves orphaned full copies of
// the design file behind - see TestReplaceSeededPatterns_DropsHistoryOfReplacedDesigns.
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
