package models

import "github.com/grafana-tools/sdk"

type GrafanaBoard struct {
	URI   string `json:"uri,omitempty"`
	Title string `json:"title,omitempty"`
	Slug  string `json:"slug,omitempty"`
	UID   string `json:"uid,omitempty"`
	OrgID uint   `json:"org_id,omitempty"`
	// Panels       []*GrafanaPanel        `json:"panels,omitempty"`
	Panels       []*sdk.Panel           `json:"panels,omitempty"`
	TemplateVars []*GrafanaTemplateVars `json:"template_vars,omitempty"`
}

type GrafanaTemplateVars struct {
	Name       string             `json:"name,omitempty"`
	Query      string             `json:"query,omitempty"`
	Datasource *GrafanaDataSource `json:"datasource,omitempty"`
	Hide       uint8              `json:"hide,omitempty"`
	Value      interface{}        `json:"value,omitempty"`
}

type GrafanaDataSource struct {
	ID   uint   `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type GrafanaPanel struct {
	PType string `json:"type,omitempty"`
	ID    uint   `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
}
