package models

// ApplicationsAPIResponse response retruned by patternfile endpoint on meshery server
type ApplicationsAPIResponse struct {
	Page         uint                 `json:"page"`
	PageSize     uint                 `json:"pageSize"`
	TotalCount   uint                 `json:"totalCount"`
	Applications []MesheryApplication `json:"applications"`
}

type ApplicationSourceTypesAPIResponse struct {
	ApplicationType     string   `json:"applicationType"`
	SupportedExtensions []string `json:"supportedExtensions"`
}
