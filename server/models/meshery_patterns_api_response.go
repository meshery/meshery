package models

// PatternsAPIResponse response retruned by patternfile endpoint on meshery server
type PatternsAPIResponse struct {
	Page       uint             `json:"page"`
	PageSize   uint             `json:"pageSize"`
	TotalCount uint             `json:"totalCount"`
	Patterns   []MesheryPattern `json:"patterns"`
}

type PatternSourceTypesAPIResponse struct {
	DesignType          string   `json:"designType"`
	SupportedExtensions []string `json:"supportedExtensions"`
}
