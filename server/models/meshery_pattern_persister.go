package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"gorm.io/gorm"

	"github.com/meshery/schemas/models/core"

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

// CloneMesheryPattern clones meshery pattern to private
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

	mesheryPattern.Visibility = Private
	mesheryPattern.ID = &id
	mesheryPattern.Name = clonePatternRequest.Name

	return mpp.SaveMesheryPattern(&mesheryPattern)
}

// DeleteMesheryPattern takes in a profile id and delete it if it already exists
func (mpp *MesheryPatternPersister) DeleteMesheryPattern(id core.Uuid) ([]byte, error) {
	pattern := MesheryPattern{ID: &id}
	mpp.DB.Delete(&pattern)

	return marshalMesheryPattern(&pattern), nil
}

// DeleteMesheryPatterns takes in a meshery-patterns and delete those if exist
func (mpp *MesheryPatternPersister) DeleteMesheryPatterns(patterns MesheryPatternDeleteRequestBody) ([]byte, error) {
	var deletedMaptterns []MesheryPattern
	for _, pObj := range patterns.Patterns {
		id := uuid.FromStringOrNil(pObj.ID)
		pattern := MesheryPattern{ID: &id}
		mpp.DB.Delete(&pattern)
		deletedMaptterns = append(deletedMaptterns, pattern)
	}

	return marshalMesheryPatterns(deletedMaptterns), nil
}

func (mpp *MesheryPatternPersister) SaveMesheryPattern(pattern *MesheryPattern) ([]byte, error) {
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

	return marshalMesheryPatterns([]MesheryPattern{*pattern}), mpp.DB.Save(pattern).Error
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
