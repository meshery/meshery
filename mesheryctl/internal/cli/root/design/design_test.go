package design

import "flag"

var (
	update                 = flag.Bool("update", false, "update golden files")
	invalidFilePath        = "/invalid/path/design.yaml"
	validDesignSourceTypes = []string{"Helm Chart", "Kubernetes Manifest", "Docker Compose", "Meshery Design"}
)

// validationErrorDetail is the longDescription Meshery Server sends in the 500
// body of deploy.validationError.response.golden.
const validationErrorDetail = "Deploying design failed: invalid component configuration for Deployment: spec.replicas: Invalid type. Expected: integer, given: string"

func resetVariables() {
	skipSave = false
	skipValidation = false
	designFile = ""
	file = ""
	designDeployFlags = cmdDesignDeployFlags{}
}
