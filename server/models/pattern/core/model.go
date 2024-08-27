package core

import (
	"context"

	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/models/meshmodel/registry"
	"github.com/meshery/schemas/models/v1beta1/pattern"
)

type ProcessPatternOptions struct {
	Context                context.Context
	Provider               models.Provider
	Pattern                pattern.PatternFile
	PrefObj                *models.Preference
	UserID                 string
	IsDelete               bool
	Validate               bool
	DryRun                 bool
	SkipCRDAndOperator     bool
	UpgradeExistingRelease bool
	SkipPrintLogs          bool
	Registry               *registry.RegistryManager
	EventBroadcaster       *models.Broadcast
	Log                    logger.Handler
}
