package academy

import (
	"fmt"

	"github.com/meshery/meshkit/errors"
)

var (
	ErrInvalidNestingCode = "mesheryctl-1254"
	ErrTaxonomyTypeCode   = "mesheryctl-1255"
	ErrScaffoldExistsCode = "mesheryctl-1256"
	ErrMissingOrgIDCode   = "mesheryctl-1257"
	ErrMissingIntoCode    = "mesheryctl-1258"
	ErrInvalidLevelCode   = "mesheryctl-1259"
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
		[]string{"A scaffold file already exists at the target path."},
		[]string{"You are attempting to overwrite an existing node."},
		[]string{"Use the --force flag to overwrite the existing file."},
	)
}

func errMissingOrgID() error {
	return errors.New(ErrMissingOrgIDCode, errors.Alert,
		[]string{"Missing organization ID"},
		[]string{"An organization ID is required when scaffolding a fresh top-level content type."},
		[]string{"No --org flag was provided."},
		[]string{"Provide the --org flag."},
	)
}

func errMissingInto() error {
	return errors.New(ErrMissingIntoCode, errors.Alert,
		[]string{"Missing target directory"},
		[]string{"The --into flag is required for non-top-level content types."},
		[]string{"You are attempting to scaffold a node that is not a learning-path, certification, or challenge without specifying where it belongs."},
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
