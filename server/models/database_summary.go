package models

type SqliteSchema struct {
	Name  string `json:"name,omitempty"`
	Type  string `json:"type,omitempty"`
	Count int64  `json:"count"`
}

type DatabaseSummary struct {
	Page        int             `json:"page"`
	PageSize    int             `json:"page_size"`
	TotalTables int             `json:"total_tables"`
	RecordCount int             `json:"record_count"`
	Tables      []*SqliteSchema `json:"tables"`
}
