package format

import "github.com/layer5io/meshkit/errors"

var (
	ErrOutputToJsonCode = "mesheryctl-1147"
)

func ErrOutputToJson() error {
	return errors.New(ErrOutputToJsonCode, errors.Alert,
		[]string{"Failed to format output in JSON"},
		[]string{"Error occurred while formatting to JSON"},
		[]string{"The content provided for formatting failed."},
		[]string{"Check the JSON structure you are providing for formatting."})
}
