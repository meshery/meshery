package handlers

import (
	"fmt"
	"sort"

	"github.com/meshery/meshery/server/models/pattern/patterns"
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
	Total            int
	Succeeded        int
	Failed           int
	FailedComponents []failedComponent
}

// tallyDeploymentOutcome tallies deployment results from the provision stage.
// It counts components once as failed if any non-prerequisite entry in any context failed.
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
		hasFailed := false
		var compFailures []failedComponent

		for _, ctx := range contexts {
			for _, summary := range ctx.Summary {
				if summary.IsPrerequisite {
					continue
				}

				if !summary.Success {
					hasFailed = true
					name := summary.CompName
					if name == "" {
						name = k
					}

					errStr := ""
					if summary.Error != nil {
						errStr = summary.Error.Error()
					} else if summary.Message != "" {
						errStr = summary.Message
					}

					compFailures = append(compFailures, failedComponent{
						Name:     name,
						Kind:     summary.Kind,
						Model:    summary.Model,
						Location: ctx.Location,
						Error:    errStr,
					})
				}
			}
		}

		if hasFailed {
			tally.Failed++
			tally.FailedComponents = append(tally.FailedComponents, compFailures...)
		} else {
			tally.Succeeded++
		}
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

// deploymentOutcomeEvent returns the severity and description for the outcome event.
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
