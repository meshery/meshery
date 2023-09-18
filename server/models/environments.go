package models

type EnvironmentPage struct {
	Connections []Connection `json:"connections"`
	TotalCount  int          `json:"total_count"`
	Page        int          `json:"page"`
	PageSize    int          `json:"page_size"`
}