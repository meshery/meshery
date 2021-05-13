package models

// PerformanceAPIResponse response retruned by performance endpoint on meshery server
type PerformanceAPIResponse struct {
	Page       uint                `json:"page"`
	PageSize   uint                `json:"page_size"`
	TotalCount uint                `json:"total_size"`
	Results    []PerformanceResult `json:"results"`
}
