package output

import (
	"fmt"

	"github.com/meshery/meshkit/errors"
)

var ErrEncodingDataCode = "replace_me"

func ErrEncodingData(err error, encoder string) error {
	return errors.New(ErrEncodingDataCode, errors.Alert, []string{fmt.Sprintf("error occured while trying to encode data in %s", encoder)}, []string{fmt.Sprintf("Encoding the data provided failed in %s format", encoder)}, []string{"Non supported characters in the data"}, []string{"Ensure the content of the data provide does not contains invalid supported characters"})
}
