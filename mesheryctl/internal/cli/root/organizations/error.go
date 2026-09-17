package organizations

import (
	"github.com/meshery/meshkit/errors"
)

const ErrCountWithOutputFormatCode = "mesheryctl-1256"

func ErrCountWithOutputFormat() error {
	return errors.New(
		ErrCountWithOutputFormatCode,
		errors.Alert,
		[]string{"--count cannot be combined with --output-format"},
		[]string{"Both --count and --output-format were provided"},
		[]string{"--count only has a defined meaning for the default table output"},
		[]string{"Use --count on its own for a total, or --output-format on its own for structured (json|yaml) output of the full page"},
	)
}
