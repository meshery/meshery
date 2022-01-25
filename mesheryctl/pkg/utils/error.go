package utils

import (
	"fmt"

	"github.com/layer5io/meshkit/errors"
)

var (
	ErrAttachAuthTokenCode = "1043"
)

// RootError returns a formatted error message with a link to 'root' command usage page at
// in addition to the error message
func RootError(msg string) string {
	return formatError(msg, cmdRoot)
}

// PerfError returns a formatted error message with a link to 'perf' command usage page at
// in addition to the error message
func PerfError(msg string) string {
	return formatError(msg, cmdPerf)
}

// SystemError returns a formatted error message with a link to 'system' command usage page
// in addition to the error message
func SystemError(msg string) string {
	return formatError(msg, cmdSystem)
}

// MeshError returns a formatted error message with a link to 'mesh' command usage page in addition to the error message
//func MeshError(msg string) string {
//	return formatError(msg, cmdMesh)
//}

// formatError returns a formatted error message with a link to the meshery command URL
func formatError(msg string, cmd cmdType) string {
	switch cmd {
	case cmdRoot:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, rootUsageURL)
	case cmdPerf:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, perfUsageURL)
	case cmdMesh:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, meshUsageURL)
	case cmdSystem:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, systemUsageURL)
	}
	return fmt.Sprintf("%s\n", msg)
}

func ErrAttachAuthToken(err error) error {
	return errors.New(ErrAttachAuthTokenCode, errors.Alert, []string{err.Error()},
		[]string{"authentication token not found. please supply a valid user token with the --token (or -t) flag. or login with `mesheryctl system login`"}, []string{}, []string{})
}
