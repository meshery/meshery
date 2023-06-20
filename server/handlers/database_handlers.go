package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/layer5io/meshery/server/models"
	"gorm.io/gorm/clause"
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
	TotalTables  int64           `json:"totalTables"`
}

const defaultPageSize = 10

func (h *Handler) GetSystemDatabase(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	var tables []*sqliteSchema
	var totalRecords int
	var totalTables int64

	limitstr := r.URL.Query().Get("pagesize")
	var limit int
	if limitstr != "all" {
		limit, _ = strconv.Atoi(limitstr)
		if limit == 0 {
			limit = defaultPageSize
		}
	}
	pagestr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pagestr)
	offset := (page) * limit
	order := r.URL.Query().Get("order")
	sort := r.URL.Query().Get("sort")
	search := r.URL.Query().Get("search")

	tableFinder := h.dbHandler.DB.Table("sqlite_schema").
		Where("type = ?", "table")

	if search != "" {
		tableFinder = tableFinder.Where("name LIKE ?", "%"+search+"%")
	}

	tableFinder.Count(&totalTables)

	if limit != 0 {
		tableFinder = tableFinder.Limit(limit)
	}
	if offset != 0 {
		tableFinder = tableFinder.Offset(offset)
	}
	if order != "" {
		if sort == "desc" {
			tableFinder = tableFinder.Order(clause.OrderByColumn{Column: clause.Column{Name: order}, Desc: true})
		} else {
			tableFinder = tableFinder.Order(order)
		}
	}

	tableFinder.Find(&tables)

	for _, table := range tables {
		h.dbHandler.DB.Table(table.Name).Count(&table.Count)
		totalRecords += int(table.Count)
	}

	databaseSummary := &databaseSummary{
		Tables:       tables,
		TotalRecords: totalRecords,
		TotalSize:    0,
		TotalTables:  totalTables,
	}

	w.Header().Set("Content-Type", "application/json")

	val, err := json.Marshal(databaseSummary)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprint(w, string(val))
}
