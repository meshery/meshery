package core

import (
	"context"

	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/meshmodel/registry"
	pattern "github.com/meshery/schemas/models/v1beta3/design"
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
