package models

// FiltersAPIResponse response retruned by filterfile endpoint on meshery server
type FiltersAPIResponse struct {
	Page       uint            `json:"page"`
	PageSize   uint            `json:"page_size"`
	TotalCount uint            `json:"total_size"`
	Filters    []MesheryFilter `json:"filters"`
}
