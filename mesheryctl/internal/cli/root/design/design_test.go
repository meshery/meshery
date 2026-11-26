package design

import (
	"flag"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

var (
	update                 = flag.Bool("update", false, "update golden files")
	invalidFilePath        = "/invalid/path/design.yaml"
	validDesignSourceTypes = []string{"Helm Chart", "Kubernetes Manifest", "Docker Compose", "Meshery Design"}
)

func resetVariables() {
	skipSave = false
	utils.SilentFlag = false
	designFile = ""
	file = ""
}
