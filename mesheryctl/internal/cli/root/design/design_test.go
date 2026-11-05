package design

import "flag"

var (
	update                 = flag.Bool("update", false, "update golden files")
	invalidFilePath        = "/invalid/path/design.yaml"
	validDesignSourceTypes = []string{"Helm Chart", "Kubernetes Manifest", "Docker Compose", "Meshery Design"}
)

func resetVariables() {
	skipSave = false
	designFile = ""
	file = ""
}
