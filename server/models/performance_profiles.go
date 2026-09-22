package models

import (
	perfprofile "github.com/meshery/schemas/models/v1beta3/performance_profile"
)

type PerformanceProfile = perfprofile.PerformanceProfile
type PerformanceProfilePage = perfprofile.PerformanceProfilePage

// PerformanceTestConfigFile is the uploadable envelope that wraps a runtime
// PerformanceTestConfig together with the technology under test. It replaces
// the legacy SMP {test, mesh} file format; "mesh" is now a free-form Meshery
// Registry model name rather than the constraining SMP ServiceMesh enum.
type PerformanceTestConfigFile struct {
	Config      *perfprofile.PerformanceTestConfig `json:"test,omitempty"`
	ServiceMesh string                             `json:"mesh,omitempty"`
}
