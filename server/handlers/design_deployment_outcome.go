package handlers

import (
	"fmt"

	"github.com/meshery/meshery/server/models/pattern/patterns"
	"github.com/meshery/meshery/server/models/pattern/stages"
	events "github.com/meshery/meshkit/models/events"
)

// deploymentOutcomeTally counts the components processed during deployment or undeployment.
type deploymentOutcomeTally struct {
	Total     int // components with a provision outcome
	Succeeded int // components for which stages.ProvisionSucceeded is true
	Failed    int // Total - Succeeded
}

// tallyDeploymentOutcome tallies deployment results from the provision stage.
// Whether a component failed is decided by stages.ProvisionSucceeded, the same
// rule the provision stage uses to withhold dependents, so the event cannot
// disagree with the deployment itself. A component is counted once however
// many clusters it failed on. Values that are not provision outcomes, and the
// "" key, are ignored.
//
// Only counts are tallied: the per-component errors are reported by
// summarizePatternProcessResponse, capped so the event metadata stays within
// models.MaxEventMetadataBytes.
func tallyDeploymentOutcome(resp map[string]interface{}) deploymentOutcomeTally {
	var tally deploymentOutcomeTally
	for k, val := range resp {
		if k == "" {
			continue
		}
		contexts, ok := val.([]patterns.DeploymentMessagePerContext)
		if !ok {
			continue
		}

		tally.Total++
		if stages.ProvisionSucceeded(contexts) {
			tally.Succeeded++
		} else {
			tally.Failed++
		}
	}
	return tally
}

// addOutcomeMetadata sets the tally on an outcome event's metadata under the
// camelCase wire keys total, succeeded and failed, overwriting any existing
// values.
func addOutcomeMetadata(metadata map[string]interface{}, tally deploymentOutcomeTally) {
	metadata["total"] = tally.Total
	metadata["succeeded"] = tally.Succeeded
	metadata["failed"] = tally.Failed
}

// deploymentOutcomeEvent returns the severity and description for the outcome event.
// When nothing failed, or there was nothing to provision, it keeps the
// existing severity (Informational for deploy, Success for undeploy) and
// okDescription; partial failure is a Warning and total failure an Error.
func deploymentOutcomeEvent(action string, tally deploymentOutcomeTally, designName, okDescription string) (events.EventSeverity, string) {
	if tally.Total == 0 || tally.Failed == 0 {
		if action == "undeploy" {
			return events.Success, okDescription
		}
		return events.Informational, okDescription
	}

	actionTitle := "Deployment"
	if action == "undeploy" {
		actionTitle = "Undeployment"
	}

	if tally.Failed == tally.Total {
		return events.Error, fmt.Sprintf("%s failed for design '%s': %d of %d components failed.", actionTitle, designName, tally.Failed, tally.Total)
	}

	return events.Warning, fmt.Sprintf("%s partially completed for design '%s': %d of %d components failed.", actionTitle, designName, tally.Failed, tally.Total)
}
