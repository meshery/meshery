package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/meshery/schemas/models/core"
	schemasWorkspace "github.com/meshery/schemas/models/v1beta3/workspace"

	"github.com/gofrs/uuid"
	"gopkg.in/yaml.v2"

	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/models/patterns"
)

// MesheryPatternPersister is the persister for persisting
// performance profiles on the database
type MesheryPatternPersister struct {
	DB *database.Handler
}

// stampLocalProviderOwner sets the built-in provider's single user as the
// design owner. MesheryPattern.UserID/User are gorm:"-" (never a DB column),
// and the local provider is single-user, so every design it persists is owned
// by the local "meshery" user. Emitting userId plus the embedded user profile
// matches the schemas v1beta3 design.MesheryPattern wire contract, so the UI
// resolves the owner without a second lookup - this is the root-cause fix for
// the design Info modal's "Owner: undefined undefined".
//
// Published designs are left untouched: on the built-in provider those are the
// seeded community catalog, which the local user neither authored nor owns.
// The guard lives here rather than at the call sites so no read path can opt
// out of it - GetMesheryPatterns accepts a visibility filter and can therefore
// return published designs too.
func stampLocalProviderOwner(p *MesheryPattern) {
	if p == nil || !LocalProviderOwnsContent(p.Visibility) {
		return
	}
	id := LocalProviderUserID
	p.UserID = &id
	p.User = LocalProviderContentUser()
}

// MesheryPatternPage represents a page of performance profiles
type MesheryPatternPage struct {
	Page       uint64            `json:"page"`
	PageSize   uint64            `json:"pageSize"`
	TotalCount int               `json:"totalCount"`
	Patterns   []*MesheryPattern `json:"patterns"`
}

// GetMesheryPatterns returns all of the 'private' patterns. Though private has no meaning here since there is only
// one local user. We make this distinction to be consistent with the remote provider
func (mpp *MesheryPatternPersister) GetMesheryPatterns(search, order string, page, pageSize uint64, updatedAfter string, visibility []string) ([]byte, error) {
	order = SanitizeOrderInput(order, []string{"created_at", "updated_at", "name"})

	if order == "" {
		order = defaultOrderUpdatedAtDesc
	}

	count := int64(0)
	patterns := []*MesheryPattern{}

	query := mpp.DB.Table("meshery_patterns")

	if len(visibility) > 0 {
		query = query.Where("visibility in (?)", visibility)
	}

	query = query.Where("updated_at > ?", updatedAfter).Order(order)

	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		query = query.Where("(lower(meshery_patterns.name) like ?)", like)
	}

	query.Count(&count)
	Paginate(uint(page), uint(pageSize))(query).Find(&patterns)

	for _, p := range patterns {
		stampLocalProviderOwner(p)
	}

	mesheryPatternPage := &MesheryPatternPage{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: int(count),
		Patterns:   patterns,
	}

	return marshalMesheryPatternPage(mesheryPatternPage), nil
}

// GetMesheryCatalogPatterns returns all of the published patterns
func (mpp *MesheryPatternPersister) GetMesheryCatalogPatterns(page, pageSize, search, order string) ([]byte, error) {
	var err error
	order = SanitizeOrderInput(order, []string{"created_at", "updated_at", "name"})

	if order == "" {
		order = defaultOrderUpdatedAtDesc
	}

	var pg int
	if page != "" {
		pg, err = strconv.Atoi(page)

		if err != nil || pg < 0 {
			pg = 0
		}
	} else {
		pg = 0
	}

	// 0 page size is for all records
	var pgSize int
	if pageSize != "" {
		pgSize, err = strconv.Atoi(pageSize)

		if err != nil || pgSize < 0 {
			pgSize = 0
		}
	} else {
		pgSize = 0
	}

	patterns := []MesheryPattern{}

	query := mpp.DB.Where("visibility = ?", Published).Order(order)

	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		query = query.Where("(lower(meshery_patterns.name) like ?)", like)
	}

	var count int64
	err = query.Model(&MesheryPattern{}).Count(&count).Error

	if err != nil {
		return nil, err
	}

	if pgSize != 0 {
		Paginate(uint(pg), uint(pgSize))(query).Find(&patterns)
	} else {
		query.Find(&patterns)
	}

	response := PatternsAPIResponse{
		Page:       uint(pg),
		PageSize:   uint(pgSize),
		TotalCount: uint(count),
		Patterns:   patterns,
	}

	marshalledResponse, _ := json.Marshal(response)
	return marshalledResponse, nil
}

// CloneMesheryPattern saves a private copy of a design under a new ID. The copy
// is a new design: it keeps the source's content and import source, but gets
// its own creation time, zero usage counters and no catalog data, and its
// design file carries the copy's id and name at the first version.
func (mpp *MesheryPatternPersister) CloneMesheryPattern(patternID string, clonePatternRequest *MesheryClonePatternRequestBody) ([]byte, error) {
	var mesheryPattern MesheryPattern
	patternUUID, err := uuid.FromString(patternID)
	if err != nil {
		return nil, ErrInvalidUUID(err)
	}
	err = mpp.DB.First(&mesheryPattern, patternUUID).Error
	if err != nil || *mesheryPattern.ID == uuid.Nil {
		return nil, fmt.Errorf("unable to get design: %w", err)
	}

	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	pf, err := patterns.GetPatternFormat(mesheryPattern.PatternFile)
	if err != nil {
		return nil, err
	}
	pf.ID = id
	pf.Name = clonePatternRequest.Name
	patterns.AssignVersion(pf)
	patternFile, err := yaml.Marshal(pf)
	if err != nil {
		return nil, err
	}

	// Build the row from scratch so nothing that belongs to the source row,
	// such as created_at or the usage counters, carries over. Timestamps are
	// left nil for gorm to set on insert. It is inserted directly because
	// SaveMesheryPattern would bump the version of a file already at 0.0.1.
	clone := MesheryPattern{
		ID:            &id,
		Name:          clonePatternRequest.Name,
		PatternFile:   string(patternFile),
		Location:      mesheryPattern.Location,
		Visibility:    Private,
		Type:          mesheryPattern.Type,
		SourceContent: mesheryPattern.SourceContent,
	}
	if err := mpp.DB.Create(&clone).Error; err != nil {
		return nil, ErrDBCreate(err)
	}

	stampLocalProviderOwner(&clone)
	return marshalMesheryPatterns([]MesheryPattern{clone}), nil
}

// DeleteMesheryPattern deletes a pattern by ID and removes any associated
// workspace design mappings in a database transaction.
func (mpp *MesheryPatternPersister) DeleteMesheryPattern(id core.Uuid) ([]byte, error) {
	pattern := MesheryPattern{ID: &id}
	err := mpp.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&pattern).Error; err != nil {
			return err
		}
		if tx.Migrator().HasTable(&schemasWorkspace.WorkspacesDesignsMapping{}) {
			if err := tx.Where("design_id = ?", id).Delete(&schemasWorkspace.WorkspacesDesignsMapping{}).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, ErrDBDelete(err, "")
	}

	return marshalMesheryPattern(&pattern), nil
}

// DeleteMesheryPatterns batch deletes the specified patterns and cleans up their
// associated workspace design mappings within a database transaction.
func (mpp *MesheryPatternPersister) DeleteMesheryPatterns(patterns MesheryPatternDeleteRequestBody) ([]byte, error) {
	var deletedMaptterns []MesheryPattern
	err := mpp.DB.Transaction(func(tx *gorm.DB) error {
		for _, pObj := range patterns.Patterns {
			id := uuid.FromStringOrNil(pObj.ID)
			pattern := MesheryPattern{ID: &id}
			if err := tx.Delete(&pattern).Error; err != nil {
				return err
			}
			if tx.Migrator().HasTable(&schemasWorkspace.WorkspacesDesignsMapping{}) {
				if err := tx.Where("design_id = ?", id).Delete(&schemasWorkspace.WorkspacesDesignsMapping{}).Error; err != nil {
					return err
				}
			}
			deletedMaptterns = append(deletedMaptterns, pattern)
		}
		return nil
	})
	if err != nil {
		return nil, ErrDBDelete(err, "")
	}

	return marshalMesheryPatterns(deletedMaptterns), nil
}

// SaveMesheryPattern inserts a new design, or updates an existing one with only
// the non-zero fields of pattern, and bumps the design file's version.
func (mpp *MesheryPatternPersister) SaveMesheryPattern(pattern *MesheryPattern) ([]byte, error) {
	pf, err := patterns.GetPatternFormat(pattern.PatternFile)
	if err != nil {
		return nil, err
	}

	isNew := pattern.ID == nil
	if isNew {
		id, err := uuid.NewV4()
		if err != nil {
			return nil, ErrGenerateUUID(err)
		}

		patterns.AssignVersion(pf)

		pattern.ID = &id
	} else {
		nextVersion, err := patterns.GetNextVersion(pf)
		if err != nil {
			return nil, err
		}
		pf.Version = nextVersion
		byt, err := yaml.Marshal(pf)
		if err != nil {
			return nil, err
		}
		pattern.PatternFile = string(byt)
	}

	// Stamp before marshalling so the save/clone response carries the same
	// owner contract as the GET paths. UserID/User are gorm:"-", so this is
	// response shaping only - nothing extra is persisted.
	stampLocalProviderOwner(pattern)

	// An ID does not mean the row exists: clone and seeding save new rows
	// under IDs they minted themselves.
	exists := false
	if !isNew {
		var count int64
		if err := mpp.DB.Model(&MesheryPattern{}).Where("id = ?", *pattern.ID).Count(&count).Error; err != nil {
			return nil, ErrDBRead(err)
		}
		exists = count > 0
	}

	updated := exists
	if exists {
		err = mpp.updateSavedColumns(pattern)
	} else {
		suppliedVisibility := pattern.Visibility
		if pattern.Visibility == "" {
			pattern.Visibility = Private
		}
		res := mpp.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(pattern)
		if res.Error != nil {
			return marshalMesheryPatterns([]MesheryPattern{*pattern}), ErrDBCreate(res.Error)
		}
		if res.RowsAffected == 0 {
			// A concurrent save inserted this ID first: update it like any
			// existing row, without the insert-only visibility default.
			pattern.Visibility = suppliedVisibility
			updated = true
			err = mpp.updateSavedColumns(pattern)
		}
	}

	// The saved struct holds only what the request carried; respond with the
	// stored row so the response shows the columns the update preserved.
	if updated && err == nil {
		var stored MesheryPattern
		if err = mpp.DB.First(&stored, "id = ?", *pattern.ID).Error; err != nil {
			err = ErrDBRead(err)
		} else {
			stampLocalProviderOwner(&stored)
			pattern = &stored
		}
	}

	return marshalMesheryPatterns([]MesheryPattern{*pattern}), err
}

// updateSavedColumns updates an existing design with the columns the save
// carries. A save carries only some of them - POST /api/pattern builds the
// record from DesignPostPayload, which has no source content, type or creation
// time - so Updates skips zero-valued fields and those columns keep their
// stored values. created_at and the usage counters are never written by a save,
// and pattern_file is not written when the request carried no design file.
func (mpp *MesheryPatternPersister) updateSavedColumns(pattern *MesheryPattern) error {
	omit := []string{"created_at", "view_count", "share_count", "download_count", "clone_count", "deployment_count"}
	if pattern.PatternFileOmitted {
		omit = append(omit, "pattern_file")
	}
	res := mpp.DB.Model(pattern).Omit(omit...).Updates(pattern)
	// updated_at always changes, so no affected row means the design was
	// deleted after the existence check.
	if res.Error == nil && res.RowsAffected == 0 {
		res.Error = fmt.Errorf("design %s was deleted while saving", pattern.ID)
	}
	if res.Error != nil {
		return ErrDBPut(res.Error)
	}
	return nil
}

// SaveMesheryPatterns batch inserts the given patterns
func (mpp *MesheryPatternPersister) SaveMesheryPatterns(mesheryPatterns []MesheryPattern) ([]byte, error) {
	finalPatterns := []MesheryPattern{}
	for _, pattern := range mesheryPatterns {

		pf, err := patterns.GetPatternFormat(pattern.PatternFile)
		if err != nil {
			return nil, err
		}

		if pattern.Visibility == "" {
			pattern.Visibility = Private
		}
		if pattern.ID == nil {
			id, err := uuid.NewV4()
			if err != nil {
				return nil, ErrGenerateUUID(err)
			}
			patterns.AssignVersion(pf)
			pattern.ID = &id
		} else {
			nextVersion, err := patterns.GetNextVersion(pf)
			if err != nil {
				return nil, err
			}
			pf.Version = nextVersion
		}

		stampLocalProviderOwner(&pattern)

		finalPatterns = append(finalPatterns, pattern)
	}

	return marshalMesheryPatterns(finalPatterns), mpp.DB.Create(finalPatterns).Error
}

// ReplaceSeededPatterns makes the published designs in the database mirror the
// supplied set, which SeedContent derives from the on-disk catalog directory.
//
// Replacement rather than upsert, for two reasons: seeded designs are minted
// with a fresh uuid.NewV4() on every pass, so re-saving them inserts duplicates
// instead of overwriting; and deleting is what drops designs whose catalog file
// has since been removed upstream, which an upsert would leave orphaned.
//
// Delete and inserts share one transaction so a failure part-way through rolls
// back, rather than committing the delete and leaving the catalog empty.
//
// Safe to scope by visibility because a local user cannot create a published
// design - PublishCatalogPattern returns ErrLocalProviderSupport - so published
// rows here are always seeded, and private user-authored designs never match.
func (mpp *MesheryPatternPersister) ReplaceSeededPatterns(seeded []*MesheryPattern) error {
	return mpp.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("visibility = ?", Published).
			Delete(&MesheryPattern{}).Error; err != nil {
			return err
		}

		txPersister := &MesheryPatternPersister{
			DB: &database.Handler{DB: tx, Mutex: &sync.Mutex{}},
		}
		for _, pattern := range seeded {
			if _, err := txPersister.SaveMesheryPattern(pattern); err != nil {
				return err
			}
		}
		return nil
	})
}

func (mpp *MesheryPatternPersister) GetMesheryPattern(id core.Uuid) ([]byte, error) {
	var mesheryPattern MesheryPattern

	err := mpp.DB.First(&mesheryPattern, id).Error
	stampLocalProviderOwner(&mesheryPattern)
	return marshalMesheryPattern(&mesheryPattern), err
}

func (mpp *MesheryPatternPersister) GetMesheryPatternSource(id core.Uuid) ([]byte, error) {
	var mesheryPattern MesheryPattern
	err := mpp.DB.First(&mesheryPattern, id).Error
	return mesheryPattern.SourceContent, err
}

func marshalMesheryPatternPage(mpp *MesheryPatternPage) []byte {
	res, _ := json.Marshal(mpp)

	return res
}

func marshalMesheryPattern(mp *MesheryPattern) []byte {
	res, _ := json.Marshal(mp)

	return res
}

func marshalMesheryPatterns(mps []MesheryPattern) []byte {
	res, _ := json.Marshal(mps)

	return res
}
