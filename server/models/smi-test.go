package models

import (
	"github.com/meshery/schemas/models/core"
)

type SmiResult struct {
	ID                core.Uuid `json:"id,omitempty"`
	Date              string    `json:"date,omitempty"`
	MeshName          string    `json:"meshName,omitempty"`
	MeshVersion       string    `json:"meshVersion,omitempty"`
	CasesPassed       string    `json:"casesPassed,omitempty"`
	PassingPercentage string    `json:"passingPercentage,omitempty"`
	Status            string    `json:"status,omitempty"`
	MoreDetails       []*Detail `json:"moreDetails,omitempty" gorm:"type:detail[]"`
}

type Detail struct {
	SmiSpecification string `json:"smiSpecification,omitempty"`
	SmiVersion       string `json:"smiVersion,omitempty"`
	Time             string `json:"time,omitempty"`
	Assertions       string `json:"assertions,omitempty"`
	Result           string `json:"result,omitempty"`
	Reason           string `json:"reason,omitempty"`
	Capability       string `json:"capability,omitempty"`
	Status           string `json:"status,omitempty"`
}
