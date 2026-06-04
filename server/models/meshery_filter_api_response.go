package models

// FiltersAPIResponse response retruned by filterfile endpoint on meshery server
type FiltersAPIResponse struct {
	Page       uint            `json:"page"`
	PageSize   uint            `json:"pageSize"`
	TotalCount uint            `json:"totalCount"`
	Filters    []MesheryFilter `json:"filters"`
}
