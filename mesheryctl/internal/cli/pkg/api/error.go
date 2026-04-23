package api

import "github.com/meshery/meshkit/errors"

var ErrGenerateDataForInvalidResponseCode = "mesheryctl-1231"

func ErrGenerateDataForInvalidResponse() error {
	return errors.New(
		ErrGenerateDataForInvalidResponseCode,
		errors.Fatal,
		[]string{"Invalid response received from Meshery server. ", "The response is nil or does not contain a body."},
		[]string{"The response from the Meshery server is invalid and cannot be processed. "},
		[]string{"An error occurred while processing the response from the Meshery server. "},
		[]string{"Ensure that the Meshery server is running and returning valid responses. "})
}
