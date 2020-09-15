package models

import (
	"github.com/gofrs/uuid"
)

type SmiResult struct {
	ID                    uuid.UUID `json:"id,string,omitempty"`
	MeshName              string    `json:"mesh_name,omitempty"`
	MeshVersion           string    `json:"mesh_version,omitempty"`
	CasesPassed           string    `json:"cases_passed,omitempty"`
	ConformanceCapability string    `json:"conformance_capability,omitempty"`
	Status                string    `json:"status,omitempty"`
	MoreDetails           []*Detail `json:"more_details,omitempty"`
}

type Detail struct {
	SmiSpecification string `json:"smi_specification,omitempty"`
	SmiVersion       string `json:"smi_version,omitempty"`
	Time             string `json:"time,omitempty"`
	Assertions       string `json:"assertions,omitemtpy"`
	Result           string `json:"result,omitempty"`
	Reason           string `json:"reason,omitempty"`
}
