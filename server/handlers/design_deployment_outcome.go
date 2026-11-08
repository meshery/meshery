package handlers

import (
	"fmt"

	events "github.com/meshery/meshkit/models/events"
)

// addOutcomeCounts sets the component counts of a deploy or undeploy summary
// on the event metadata under the camelCase wire keys total, succeeded and
// failed, overwriting any existing values. They count components, not
// outcomes: see patternProcessSummary for how failed differs from errorCount.
func addOutcomeCounts(metadata map[string]interface{}, summary patternProcessSummary) {
	metadata["total"] = summary.Total
	metadata["succeeded"] = summary.Succeeded
	metadata["failed"] = summary.Failed
}

// deploymentOutcomeEvent returns the severity and description of a deploy or
// undeploy outcome event, chosen by the summary's Status so the event and its
// status metadata cannot disagree. A successful outcome, including one with
// nothing to provision, keeps the existing severity (Informational for deploy,
// Success for undeploy) and okDescription; a partial outcome is a Warning and a
// failed one an Error.
func deploymentOutcomeEvent(action string, summary patternProcessSummary, designName, okDescription string) (events.EventSeverity, string) {
	actionTitle := "Deployment"
	if action == "undeploy" {
		actionTitle = "Undeployment"
	}

	switch summary.Status {
	case patternProcessStatusFailed:
		return events.Error, fmt.Sprintf("%s failed for design '%s': %d of %d components failed.", actionTitle, designName, summary.Failed, summary.Total)
	case patternProcessStatusPartial:
		return events.Warning, fmt.Sprintf("%s partially completed for design '%s': %d of %d components failed.", actionTitle, designName, summary.Failed, summary.Total)
	}

	if action == "undeploy" {
		return events.Success, okDescription
	}
	return events.Informational, okDescription
}
