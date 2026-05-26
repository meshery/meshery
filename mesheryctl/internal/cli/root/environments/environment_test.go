package environments

import (
	"flag"
)

var update = flag.Bool("update", false, "update golden files")

var testConstants = map[string]string{
	"orgID":           "2d2c0b60-076a-4f0a-8a63-de538570a553",
	"environmentName": "test-environment",
}
