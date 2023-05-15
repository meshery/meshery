package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/layer5io/meshery/server/models"
)

type sqliteSchema struct {
	Name  string `json:"name,omitempty"`
	Type  string `json:"type,omitempty"`
	Count int64  `json:"count"`
}

type databaseSummary struct {
	Tables       []*sqliteSchema `json:"tables"`
	TotalRecords int             `json:"totalRecords"`
	TotalSize    int             `json:"totalSize"`
	TotalTables  int             `json:"totalTables"`
}

func (h *Handler) GetSystemDatabase(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	var tables []*sqliteSchema
	var totalRecords int

	h.dbHandler.DB.Table("sqlite_schema").
		Where("type = ?", "table").
		Find(&tables)

	for _, table := range tables {
		h.dbHandler.DB.Table(table.Name).Count(&table.Count)
		totalRecords += int(table.Count)
	}

	databaseSummary := &databaseSummary{
		Tables:       tables,
		TotalRecords: totalRecords,
		TotalSize:    0,
		TotalTables:  len(tables),
	}

	w.Header().Set("Content-Type", "application/json")

	val, err := json.Marshal(databaseSummary)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprint(w, string(val))
}
