package meshmodel

// Different types of issues encountered while processing a Design are defined here

type IssueType string

const (
	JSONSchemaValidationError IssueType = "JsonSchemaValidationError"
	PolicyError               IssueType = "PolicyError"

	K8sAdminssionControllerError IssueType = "AdmissionControllerError"
	K8sInternalError             IssueType = "K8sInternalError"
)

type JSONSchemaErrorDetail struct{}
type PolicyErrorDetail struct{}
type K8sAdmissionControllerErrorDetail struct{}
type K8sInternalErrorDetail struct{}
