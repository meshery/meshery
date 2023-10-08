package handlers

import (
	"github.com/gofrs/uuid"
	"github.com/layer5io/meshkit/database"
	"gorm.io/gorm"
)

type MesheryHostsPersistor struct {
	DB *database.Handler
}

type MesheryHosts struct {
	ID       uuid.UUID
	Hostname string
	Port     int
}
type MesheryHostsDisplay struct {
	ID       uuid.UUID           `json:"id"`
	Hostname string              `json:"hostname"`
	Port     int                 `json:"port"`
	Summary  hostIndividualCount `json:"summary"`
}
type MesheryHostsContextPage struct {
	Page        int                   `json:"page"`
	PageSize    int                   `json:"page_size"`
	Count       int64                 `json:"total_count"`
	Registrants []MesheryHostsDisplay `json:"registrants"`
}
type hostIndividualCount struct {
	Models        int64 `json:"models"`
	Components    int64 `json:"components"`
	Relationships int64 `json:"relationships"`
}

func (h Handler) GetMesheryHosts(limit int, page int, search string) ([]MesheryHostsDisplay, int64, error) {
	db := h.dbHandler.DB
	var hosts []*MesheryHosts
	var hostsCount int64
	if limit == -1 {
		if err := db.Table("hosts").Select("*").Count(&hostsCount).Error; err != nil {
			return nil, 0, err
		}
		limit = int(hostsCount)
	}
	if err := db.Table("hosts").
		Select("hosts.id ,hosts.hostname,hosts.port").
		Where("hostname LIKE ?", "%"+search+"%").
		Offset(limit * (page - 1)).
		Limit(limit).
		Scan(&hosts).
		Count(&hostsCount).
		Error; err != nil {
		return nil, 0, err
	}
	var result []MesheryHostsDisplay
	for _, host := range hosts {
		var hostCounts hostIndividualCount
		if count, err := getCountForHostAndType(db, host.ID, "model"); err != nil {
			return nil, 0, err
		} else {
			hostCounts.Models = count
		}

		if count, err := getCountForHostAndType(db, host.ID, "component"); err != nil {
			return nil, 0, err
		} else {
			hostCounts.Components = count
		}

		if count, err := getCountForHostAndType(db, host.ID, "relationship"); err != nil {
			return nil, 0, err
		} else {
			hostCounts.Relationships = count
		}
		res := MesheryHostsDisplay{
			ID:       host.ID,
			Hostname: host.Hostname,
			Port:     host.Port,
			Summary: hostIndividualCount{
				Models:        hostCounts.Models,
				Components:    hostCounts.Components,
				Relationships: hostCounts.Relationships,
			},
		}
		result = append(result, res)
	}
	return result, hostsCount, nil
}

func getCountForHostAndType(db *gorm.DB, hostID uuid.UUID, typeName string) (int64, error) {
	var count int64
	err := db.Table("hosts").
		Select("hosts.id").
		Joins("LEFT JOIN registries ON hosts.id = registries.registrant_id").
		Where("hosts.id = ? AND registries.type LIKE ?", hostID, "%"+typeName+"%").
		Count(&count).
		Error

	return count, err
}
