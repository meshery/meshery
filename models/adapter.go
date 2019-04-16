package models

type Adapter struct {
	Location string            `json:"adapter_location"`
	Name     string            `json:"name"`
	Ops      map[string]string `json:"ops"`
}
