package academy

import (
	"fmt"

	"github.com/meshery/meshkit/errors"
)

var (
	ErrInvalidNestingCode = "mesheryctl-exp-1001"
	ErrTaxonomyTypeCode   = "mesheryctl-exp-1002"
	ErrScaffoldExistsCode = "mesheryctl-exp-1003"
	ErrMissingOrgIDCode   = "mesheryctl-exp-1004"
	ErrMissingIntoCode    = "mesheryctl-exp-1005"
	ErrInvalidLevelCode   = "mesheryctl-exp-1006"
)

func errInvalidNesting(parent, child string) error {
	return errors.New(ErrInvalidNestingCode, errors.Alert,
		[]string{fmt.Sprintf("Invalid nesting: cannot nest %s under %s", child, parent)},
		[]string{"The academy taxonomy rules forbid this hierarchy."},
		[]string{"Ensure you are scaffolding content in the correct parent directory."},
		[]string{fmt.Sprintf("Check allowed children for %s.", parent)},
	)
}

func errTaxonomyType(invalidType string) error {
	return errors.New(ErrTaxonomyTypeCode, errors.Alert,
		[]string{"Invalid taxonomy type"},
		[]string{"The provided content type '" + invalidType + "' is not supported."},
		[]string{"You are attempting to scaffold an unsupported taxonomy type. Valid root types are: learning-path, certification, challenge. Structural types (course, module, etc) must be created via their subcommands."},
		[]string{"Provide a valid root type (e.g. learning-path, certification, challenge) or use the appropriate subcommand."},
	)
}

func errScaffoldExists(path string) error {
	return errors.New(ErrScaffoldExistsCode, errors.Alert,
		[]string{fmt.Sprintf("Scaffold already exists at %s", path)},
		[]string{"An _index.md file already exists in the target directory."},
		[]string{"You are attempting to overwrite an existing node."},
		[]string{"Use the --force flag to overwrite the existing file."},
	)
}

func errMissingOrgID() error {
	return errors.New(ErrMissingOrgIDCode, errors.Alert,
		[]string{"Missing organization ID"},
		[]string{"An organization ID is required when scaffolding a fresh top-level content type."},
		[]string{"No --org flag was provided and no default organization context could be resolved."},
		[]string{"Pass the --org flag or set a default organization in your meshconfig."},
	)
}

func errMissingInto() error {
	return errors.New(ErrMissingIntoCode, errors.Alert,
		[]string{"Missing target directory"},
		[]string{"The --into flag is required for non-top-level content types."},
		[]string{"You are attempting to scaffold a node that is not a learning-path or certification without specifying where it belongs."},
		[]string{"Provide the --into flag with the path to the parent directory."},
	)
}

func errInvalidLevel(invalidLevel string) error {
	return errors.New(ErrInvalidLevelCode, errors.Alert,
		[]string{"Invalid level"},
		[]string{"The provided level '" + invalidLevel + "' is not supported."},
		[]string{"You are attempting to set an unsupported level. Valid levels are: beginner, intermediate, advanced."},
		[]string{"Provide a valid --level argument (beginner, intermediate, or advanced)."},
	)
}
