package models

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
type PatternsAPIResponse struct {
=======
type PatternsApiResponse struct {
>>>>>>> 55a519b9... fix merge conflict
	Page       uint             `json:"page"`
	PageSize   uint             `json:"page_size"`
	TotalCount uint             `json:"total_size"`
	Patterns   []MesheryPattern `json:"patterns"`
<<<<<<< HEAD
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
=======
type PatternsAPIResponse struct {
=======
type PatternsApiResponse struct {
>>>>>>> e56351c0e3ec180813a2d3638991292a528bd34d
>>>>>>> 9af444385e81347cabaf254dcad319009aef2028
	Page       uint             `json:"page"`
	PageSize   uint             `json:"page_size"`
	TotalCount uint             `json:"total_size"`
	Patterns   []MesheryPattern `json:"patterns"`
<<<<<<< HEAD
>>>>>>> 4c2ebc10... adds struct representation for patterns api response
=======
>>>>>>> 0753cdeb... adds api response model
=======
>>>>>>> 55a519b9... fix merge conflict
=======
>>>>>>> 9af444385e81347cabaf254dcad319009aef2028
}
