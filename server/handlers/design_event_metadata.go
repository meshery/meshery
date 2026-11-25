package handlers

import (
	"sort"

	"github.com/meshery/meshery/server/models"
	patterncore "github.com/meshery/meshery/server/models/pattern/core"
	"github.com/meshery/meshery/server/models/pattern/patterns"
	"github.com/meshery/meshery/server/models/pattern/stages"
	meshkiterrors "github.com/meshery/meshkit/errors"
	"github.com/meshery/schemas/models/core"
)

// maxPatternEventErrors caps the errors[] list of a pattern deploy, undeploy
// or dry-run event; errorCount always carries the full count.
const maxPatternEventErrors = 20

// Status values of a pattern deploy, undeploy or dry-run event.
const (
	patternProcessStatusSuccess = "success"
	patternProcessStatusPartial = "partial"
	patternProcessStatusFailed  = "failed"
)

// dryRunFailureCode is the code of a dry-run failure that carries no
// per-cause type: it is the Kubernetes Status value of a failed request.
const dryRunFailureCode = "Failure"

// patternEventError is one entry of the errors[] list of a pattern event. It
// names the failing component and the error code, never the component itself.
type patternEventError struct {
	Code        string `json:"code"`
	ComponentID string `json:"componentId"`
}

// patternProcessSummary is what a pattern event reports about the response of
// _processPattern.
type patternProcessSummary struct {
	Status     string
	ErrorCount int
	Errors     []patternEventError
	Truncated  bool
}

// summarizePatternProcessResponse derives the status and per-component errors
// from the response of _processPattern, which holds either
//
//   - one entry per component id, carrying the []patterns.DeploymentMessagePerContext
//     the Provision stage recorded for it (deploy and undeploy), or
//   - stages.DryRunResponseKey, carrying component name -> context id ->
//     patterncore.DryRunResponseWrapper (dry-run).
//
// A provisioning entry fails when a non-prerequisite outcome reports
// Success: false, mirroring the Provision stage's own decision; its code is the
// MeshKit code of the outcome's error, empty when that error is not a MeshKit
// error. A dry-run entry fails when Success is false; it contributes one error
// per failure cause, coded by the cause type (dryRunFailureCode when there are
// none).
func summarizePatternProcessResponse(response map[string]interface{}) patternProcessSummary {
	var errs []patternEventError
	failed := map[string]bool{}

	keys := make([]string, 0, len(response))
	for key := range response {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		switch value := response[key].(type) {
		case []patterns.DeploymentMessagePerContext:
			if key == "" {
				continue
			}
			if _, seen := failed[key]; !seen {
				failed[key] = false
			}
			for _, perContext := range value {
				for _, outcome := range perContext.Summary {
					if outcome.IsPrerequisite || outcome.Success {
						continue
					}
					failed[key] = true
					errs = append(errs, patternEventError{Code: meshkitCode(outcome.Error), ComponentID: key})
				}
			}
		case map[string]map[string]patterncore.DryRunResponseWrapper:
			if key != stages.DryRunResponseKey {
				continue
			}
			errs = append(errs, dryRunErrors(value, failed)...)
		}
	}

	summary := patternProcessSummary{
		ErrorCount: len(errs),
		Errors:     errs,
	}
	failedCount := 0
	for _, isFailed := range failed {
		if isFailed {
			failedCount++
		}
	}
	switch {
	case failedCount == 0:
		summary.Status = patternProcessStatusSuccess
	case failedCount < len(failed):
		summary.Status = patternProcessStatusPartial
	default:
		summary.Status = patternProcessStatusFailed
	}
	if len(summary.Errors) > maxPatternEventErrors {
		summary.Errors = summary.Errors[:maxPatternEventErrors]
		summary.Truncated = true
	}
	if summary.Errors == nil {
		summary.Errors = []patternEventError{}
	}
	return summary
}

func dryRunErrors(resp map[string]map[string]patterncore.DryRunResponseWrapper, failed map[string]bool) []patternEventError {
	var errs []patternEventError

	names := make([]string, 0, len(resp))
	for name := range resp {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		perContext := resp[name]
		contexts := make([]string, 0, len(perContext))
		for ctxID := range perContext {
			contexts = append(contexts, ctxID)
		}
		sort.Strings(contexts)

		for _, ctxID := range contexts {
			wrapper := perContext[ctxID]
			componentID := name
			if wrapper.Component != nil {
				componentID = wrapper.Component.ID.String()
			}
			failed[componentID] = failed[componentID] || !wrapper.Success
			if wrapper.Success {
				continue
			}
			if wrapper.Error == nil || len(wrapper.Error.Causes) == 0 {
				errs = append(errs, patternEventError{Code: dryRunFailureCode, ComponentID: componentID})
				continue
			}
			for _, cause := range wrapper.Error.Causes {
				code := cause.Type
				if code == "" {
					code = dryRunFailureCode
				}
				errs = append(errs, patternEventError{Code: code, ComponentID: componentID})
			}
		}
	}
	return errs
}

// meshkitCode returns the MeshKit error code of err, or "" when err is nil or
// not a MeshKit error.
func meshkitCode(err error) string {
	if err == nil {
		return ""
	}
	code := meshkiterrors.GetCode(err)
	if code == "None" {
		return ""
	}
	return code
}

// buildPatternProcessEventMetadata builds the metadata of a successful
// pattern deploy, undeploy or dry-run event: the design's identity, the
// operator view viewLink kept for existing readers, and the status and per-component errors
// of the response - never the raw response, which embeds the components
// (#22764).
func buildPatternProcessEventMetadata(response map[string]interface{}, serverURL string, designID core.Uuid, name, version string) map[string]interface{} {
	metadata := models.DesignEventIdentity(serverURL, designID, name, version)
	metadata["viewLink"] = models.DesignOperatorViewLink(serverURL, designID)

	summary := summarizePatternProcessResponse(response)
	metadata["status"] = summary.Status
	metadata["errorCount"] = summary.ErrorCount
	metadata["errors"] = summary.Errors
	metadata["truncated"] = summary.Truncated
	return metadata
}
