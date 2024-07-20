package meshmodel

import (
	"github.com/layer5io/meshkit/errors"
)

const (
	ErrSeedingComponentsCode = "replace-me"
)

func ErrSeedingComponents(err error) error {
	return errors.New(ErrSeedingComponentsCode, errors.Alert, []string{"Failed to seed meshery components."}, []string{err.Error()}, []string{}, []string{"Make sure the models being seeded are valid."})
}

