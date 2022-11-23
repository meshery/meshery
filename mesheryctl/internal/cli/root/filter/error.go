package filter

import (
	"strconv"

	"github.com/layer5io/meshkit/errors"
)

const (
	ErrInvalidAuthTokenCode = "1000"
	ErrInvalidAPICallCode   = "1001"
	ErrReadAPIResponseCode  = "1002"
	ErrUnmarshalCode        = "1003"
)

func ErrInvalidAuthToken() error {
	return errors.New(ErrInvalidAuthTokenCode, errors.Alert, []string{"authentication token not found. please supply a valid user token with the --token (or -t) flag"}, []string{}, []string{}, []string{})
}

func ErrInvalidAPICall(statusCode int) error {
	return errors.New(ErrInvalidAPICallCode, errors.Alert, []string{"Response Status Code ", strconv.Itoa(statusCode), " possible Server Error"}, []string{"Response Status Code " + strconv.Itoa(statusCode) + " possible Server Error"}, []string{}, []string{})
}

func ErrReadAPIResponse(err error) error {
	return errors.New(ErrReadAPIResponseCode, errors.Alert, []string{"failed to read response body"}, []string{"Failed to read response body", " " + err.Error()}, []string{}, []string{})
}

func ErrUnmarshal(err error) error {
	return errors.New(ErrUnmarshalCode, errors.Alert, []string{"Error unmarshalling response"}, []string{"Error processing JSON response from server.\n" + err.Error()}, []string{}, []string{})
}
