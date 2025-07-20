package environments

import "flag"

var update = flag.Bool("update", false, "update golden files")

var testConstants = map[string]string{
	"orgID":           "3f8319e0-33a9-4736-b248-12bc3def4a56",
	"environmentName": "test-environment",
}
