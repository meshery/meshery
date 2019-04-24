package models

import "context"

type Adapter struct {
	Location string            `json:"adapter_location"`
	Name     string            `json:"name"`
	Ops      map[string]string `json:"ops"`
}

type AdaptersTrackerInterface interface {
	AddAdapter(context.Context, string)
	RemoveAdapter(context.Context, string)
	GetAdapters(context.Context) []string
}
