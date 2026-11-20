package models

import (
	"testing"

	mkerrors "github.com/meshery/meshkit/errors"
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

func TestUserCapabilitiesPersister_WriteNilData(t *testing.T) {
	u := &UserCapabilitiesPersister{DB: newMigratedDB(t)}

	assert.NotPanics(t, func() {
		err := u.WriteCapabilitiesForUser("test-user", nil)
		assert.Equal(t, ErrNilConfigData, err)
	})
}

func TestUserCapabilitiesPersister_NotFound(t *testing.T) {
	u := &UserCapabilitiesPersister{DB: newMigratedDB(t)}

	props, err := u.ReadCapabilitiesForUser("missing-user")
	assert.Nil(t, props)
	assert.Equal(t, ErrResultNotFoundCode, mkerrors.GetCode(err))

	err = u.DeleteCapabilitiesForUser("missing-user")
	assert.Equal(t, ErrResultNotFoundCode, mkerrors.GetCode(err))

	assert.NoError(t, u.WriteCapabilitiesForUser("test-user", &ProviderProperties{ProviderName: "Meshery"}))
	props, err = u.ReadCapabilitiesForUser("test-user")
	assert.NoError(t, err)
	if assert.NotNil(t, props) {
		assert.Equal(t, "Meshery", props.ProviderName)
	}

	assert.NoError(t, u.DeleteCapabilitiesForUser("test-user"))
	_, err = u.ReadCapabilitiesForUser("test-user")
	assert.Equal(t, ErrResultNotFoundCode, mkerrors.GetCode(err))
}
