package models

import (
	"context"

	"github.com/layer5io/meshery/meshes"
)

type Adapter struct {
	Location string                       `json:"adapter_location"`
	Name     string                       `json:"name"`
	Ops      []*meshes.SupportedOperation `json:"ops"`
}

type AdaptersTrackerInterface interface {
	AddAdapter(context.Context, string)
	RemoveAdapter(context.Context, string)
	GetAdapters(context.Context) []string
}
