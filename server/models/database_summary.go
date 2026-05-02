package models

type SqliteSchema struct {
	Name  string `json:"name,omitempty"`
	Type  string `json:"type,omitempty"`
	Count int64  `json:"count"`
}

type DatabaseSummary struct {
	Page        int             `json:"page"`
	PageSize    int             `json:"pageSize"`
	TotalTables int             `json:"totalTables"`
	RecordCount int             `json:"recordCount"`
	Tables      []*SqliteSchema `json:"tables"`
}
