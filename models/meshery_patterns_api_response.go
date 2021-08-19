package models

// PatternsAPIResponse response retruned by patternfile endpoint on meshery server
type PatternsAPIResponse struct {
	// example: 0
	Page       uint             `json:"page"`
	// example: 0
	PageSize   uint             `json:"page_size"`
	// example: 0
	TotalCount uint             `json:"total_count"`
	Patterns   []MesheryPattern `json:"patterns"`
}
