package models

type PatternsApiResponse struct {
<<<<<<< HEAD
	Page       uint             `json:"page"`
	PageSize   uint             `json:"page_size"`
	TotalCount uint             `json:"total_size"`
	Patterns   []MesheryPattern `json:"patterns"`
=======
	Page       uint              `json:"page"`
	PageSize   uint              `json:"page_size"`
	TotalCount uint              `json:"total_size"`
	Patterns   *[]MesheryPattern `json:"patterns"`
>>>>>>> 0753cdebf0bdc94854c8d3b13dea8017790aea4f
}
