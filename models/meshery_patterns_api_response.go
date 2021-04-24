package models

<<<<<<< HEAD
type PatternsAPIResponse struct {
=======
type PatternsApiResponse struct {
>>>>>>> e56351c0e3ec180813a2d3638991292a528bd34d
	Page       uint             `json:"page"`
	PageSize   uint             `json:"page_size"`
	TotalCount uint             `json:"total_size"`
	Patterns   []MesheryPattern `json:"patterns"`
}
