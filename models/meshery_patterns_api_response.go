package models

<<<<<<< HEAD
<<<<<<< HEAD
type PatternsAPIResponse struct {
	Page       uint             `json:"page"`
	PageSize   uint             `json:"page_size"`
	TotalCount uint             `json:"total_size"`
	Patterns   []MesheryPattern `json:"patterns"`
=======
type PatternsApiResponse struct {
<<<<<<< HEAD
=======
type PatternsApiResponse struct {
>>>>>>> 0753cdeb... adds api response model
	Page       uint              `json:"page"`
	PageSize   uint              `json:"page_size"`
	TotalCount uint              `json:"total_size"`
	Patterns   *[]MesheryPattern `json:"patterns"`
<<<<<<< HEAD
>>>>>>> 87573286... adds api response model
=======
	Page       uint             `json:"page"`
	PageSize   uint             `json:"page_size"`
	TotalCount uint             `json:"total_size"`
	Patterns   []MesheryPattern `json:"patterns"`
>>>>>>> 4c2ebc10... adds struct representation for patterns api response
=======
>>>>>>> 0753cdeb... adds api response model
}
