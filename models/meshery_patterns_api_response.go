package models

// PatternsAPIResponse response retruned by patternfile endpoint on meshery server
type PatternsAPIResponse struct {
	Page       uint             `json:"page"`
	PageSize   uint             `json:"page_size"`
	TotalCount uint             `json:"total_size"`
	Patterns   []MesheryPattern `json:"patterns"`
}
