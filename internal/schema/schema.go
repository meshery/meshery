package schema

import (
	"github.com/layer5io/meshery/internal/schema/meshsync"
	// "github.com/layer5io/meshery/internal/schema/smi"
)

type Schema struct {
	// SmiConformance smi.Schema
	Meshsync meshsync.Schema
}
