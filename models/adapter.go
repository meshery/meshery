package models

import (
	"context"

	"github.com/layer5io/meshery/meshes"
)

// Adapter represents an adapter in Meshery
type Adapter struct {
	Location     string                       `json:"adapter_location"`
	Name         string                       `json:"name"`
	Version      string                       `json:"version"`
	GitCommitSHA string                       `json:"git_commit_sha"`
	Ops          []*meshes.SupportedOperation `json:"ops"`
}

// AdaptersTrackerInterface defines the methods a type should implement to be an adapter tracker
type AdaptersTrackerInterface interface {
	AddAdapter(context.Context, Adapter)
	RemoveAdapter(context.Context, Adapter)
	GetAdapters(context.Context) []Adapter
}
