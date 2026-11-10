package models

// MesheryPatternDeleteRequestBody refers to the type of request body
// that DeleteMultiMesheryPatternsHandler would receive
type MesheryPatternDeleteRequestBody struct {
	Patterns []deletePatternModel `json:"patterns,omitempty"`
}

// DeletePatternModel is the model for individual patterns to be deleted
type deletePatternModel struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
