package validation

// https://docs.google.com/document/d/1y-LHhOcQGA9fiUbN5CR42CXvjaPxKAJHftGrLgxJPTo
type ValidationErrorType string

const (
	JSONSchemaValidationError ValidationErrorType = "JsonSchemaValidationError"
	PolicyValidationError     ValidationErrorType = "PolicyValidationError"

	K8sAdminssionControllerValidationError ValidationErrorType = "AdmissionControllerValidationError"
	K8sInternalValidationError             ValidationErrorType = "K8sInternalValidationError"
)

type JSONSchemaError struct {
	instancePath string
	schemaPath   string
	propertyName string
	message      string
}
type PolicyError struct{}
type K8sAdmissionControllerError struct{}
type K8sInternalError struct{}
