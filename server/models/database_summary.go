package models

type TableSchema struct {
	Name  string `json:"name,omitempty"`
	Count int64  `json:"count"`
}

type DatabaseSummary struct {
	Page        int            `json:"page"`
	PageSize    int            `json:"page_size"`
	TotalTables int            `json:"total_tables"`
	RecordCount int            `json:"record_count"`
	Tables      []*TableSchema `json:"tables"`
}
