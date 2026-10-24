package handlers

import (
	"fmt"
	"sort"

	"github.com/meshery/meshery/server/models/pattern/patterns"
	"github.com/meshery/meshery/server/models/pattern/stages"
	events "github.com/meshery/meshkit/models/events"
)

// failedComponent records the failure details of a component during deployment or undeployment.
type failedComponent struct {
	Name     string `json:"name"`
	Kind     string `json:"kind"`
	Model    string `json:"model"`
	Location string `json:"location"`
	Error    string `json:"error"`
}

// deploymentOutcomeTally summarizes the components processed during deployment or undeployment.
type deploymentOutcomeTally struct {
	Total            int               // components with a provision outcome
	Succeeded        int               // components for which stages.ProvisionSucceeded is true
	Failed           int               // Total - Succeeded
	FailedComponents []failedComponent // one entry per failed non-prerequisite outcome, sorted
}

// tallyDeploymentOutcome tallies deployment results from the provision stage.
// Whether a component failed is decided by stages.ProvisionSucceeded, the same
// rule the provision stage uses to withhold dependents, so the event cannot
// disagree with the deployment itself. A component is counted once however
// many clusters it failed on; failedComponents lists every failed entry.
// Values that are not provision outcomes, and the "" key, are ignored.
func tallyDeploymentOutcome(resp map[string]interface{}) deploymentOutcomeTally {
	tally := deploymentOutcomeTally{
		FailedComponents: make([]failedComponent, 0),
	}
	if len(resp) == 0 {
		return tally
	}

	keys := make([]string, 0, len(resp))
	for k := range resp {
		if k == "" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		val := resp[k]
		contexts, ok := val.([]patterns.DeploymentMessagePerContext)
		if !ok {
			continue
		}

		tally.Total++
		if stages.ProvisionSucceeded(contexts) {
			tally.Succeeded++
			continue
		}

		tally.Failed++
		tally.FailedComponents = append(tally.FailedComponents, failedEntries(k, contexts)...)
	}

	sort.Slice(tally.FailedComponents, func(i, j int) bool {
		if tally.FailedComponents[i].Name != tally.FailedComponents[j].Name {
			return tally.FailedComponents[i].Name < tally.FailedComponents[j].Name
		}
		if tally.FailedComponents[i].Location != tally.FailedComponents[j].Location {
			return tally.FailedComponents[i].Location < tally.FailedComponents[j].Location
		}
		return tally.FailedComponents[i].Error < tally.FailedComponents[j].Error
	})

	return tally
}

// addOutcomeMetadata adds the tally to an outcome event's metadata under the
// camelCase wire keys total, succeeded, failed and failedComponents. Keys
// already in metadata are left as they are.
func addOutcomeMetadata(metadata map[string]interface{}, tally deploymentOutcomeTally) {
	metadata["total"] = tally.Total
	metadata["succeeded"] = tally.Succeeded
	metadata["failed"] = tally.Failed
	metadata["failedComponents"] = tally.FailedComponents
}

// failedEntries lists the failed non-prerequisite outcomes of one component,
// one per cluster it failed on. key names the component when an outcome
// carries no CompName. The error is stored as a string because a raw error
// value in the event metadata would JSON-encode as {}.
func failedEntries(key string, contexts []patterns.DeploymentMessagePerContext) []failedComponent {
	var entries []failedComponent
	for _, ctx := range contexts {
		for _, summary := range ctx.Summary {
			if summary.IsPrerequisite || summary.Success {
				continue
			}

			name := summary.CompName
			if name == "" {
				name = key
			}

			errStr := summary.Message
			if summary.Error != nil {
				errStr = summary.Error.Error()
			}

			entries = append(entries, failedComponent{
				Name:     name,
				Kind:     summary.Kind,
				Model:    summary.Model,
				Location: ctx.Location,
				Error:    errStr,
			})
		}
	}
	return entries
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
