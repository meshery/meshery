package main

import (
	"strings"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/logger"
	"github.com/spf13/viper"
)

// configureInstanceID decides the ID this Meshery Server instance identifies
// itself by and stores it in v under INSTANCE_ID.
//
// INSTANCE_ID may be set in the environment or in server-config.env, as the
// environment variable docs describe. A valid UUID there is used as is.
// Otherwise a new ID is generated, with a warning if a value was given.
//
// The ID is stored with Set, not SetDefault. With AutomaticEnv on, an
// INSTANCE_ID environment variable outranks a default, so viper.Get would
// return the raw string instead of the *core.Uuid that handlers type-assert,
// and the server would panic while building its handler. Set outranks the
// environment, so every reader gets the same *core.Uuid, and GetString
// returns its string form.
func configureInstanceID(v *viper.Viper, log logger.Handler) (*uuid.UUID, error) {
	if configured := strings.TrimSpace(v.GetString("INSTANCE_ID")); configured != "" {
		id, err := uuid.FromString(configured)
		if err == nil && id != uuid.Nil {
			v.Set("INSTANCE_ID", &id)
			return &id, nil
		}
		log.Warnf("INSTANCE_ID %q is not a valid UUID; generating a new instance ID instead", configured)
	}

	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	v.Set("INSTANCE_ID", &id)
	return &id, nil
}
