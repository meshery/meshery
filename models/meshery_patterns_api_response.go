package models

<<<<<<< HEAD
type PatternsAPIResponse struct {
	Page       uint             `json:"page"`
	PageSize   uint             `json:"page_size"`
	TotalCount uint             `json:"total_size"`
	Patterns   []MesheryPattern `json:"patterns"`
=======
type PatternsApiResponse struct {
	Page       uint              `json:"page"`
	PageSize   uint              `json:"page_size"`
	TotalCount uint              `json:"total_size"`
	Patterns   *[]MesheryPattern `json:"patterns"`
>>>>>>> 87573286... adds api response model
}
