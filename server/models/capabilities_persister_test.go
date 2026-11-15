package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserCapabilitiesPersister_NilReceiverSafety(t *testing.T) {
	var u *UserCapabilitiesPersister

	assert.NotPanics(t, func() {
		props, err := u.ReadCapabilitiesForUser("test-user")
		assert.Nil(t, props)
		assert.Equal(t, ErrDBConnection, err)
	})

	assert.NotPanics(t, func() {
		err := u.WriteCapabilitiesForUser("test-user", &ProviderProperties{})
		assert.Equal(t, ErrDBConnection, err)
	})

	assert.NotPanics(t, func() {
		err := u.DeleteCapabilitiesForUser("test-user")
		assert.Equal(t, ErrDBConnection, err)
	})
}
