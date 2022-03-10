package utils

import (
	"fmt"
	"strconv"

	"github.com/layer5io/meshkit/errors"
)

var (
	ErrFailRequestCode     = "1044"
	ErrFailReqStatusCode   = "1045"
	ErrAttachAuthTokenCode = "1046"
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
		[]string{"Authentication token not found. please supply a valid user token with the --token (or -t) flag. or login with `mesheryctl system login`"}, []string{}, []string{})
}

func ErrFailRequest(err error) error {
	return errors.New(ErrFailRequestCode, errors.Alert, []string{},
		[]string{"Failed to make a request"}, []string{}, []string{})
}

func ErrUnauthenticated() error {
	return errors.New(ErrFailRequestCode, errors.Alert, []string{},
		[]string{"Authentication token is invalid. Please supply a valid user token. Login with `mesheryctl system login`"}, []string{}, []string{})
}

func InvalidToken() error {
	return errors.New(ErrFailRequestCode, errors.Alert, []string{},
		[]string{"Authentication token is expired/invalid. Please login with `mesheryctl system login` to generate a new token"}, []string{}, []string{})
}

func ErrFailReqStatus(statusCode int) error {
	return errors.New(ErrFailReqStatusCode, errors.Alert, []string{},
		[]string{"Response Status Code " + strconv.Itoa(statusCode) + ", possible Server Error"}, []string{}, []string{})
}
