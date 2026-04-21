package models

import (
	"fmt"
	"testing"

	meshkitErrors "github.com/meshery/meshkit/errors"
)

func TestErrSaveUserCredential_IsMeshKitError(t *testing.T) {
	err := ErrSaveUserCredential(fmt.Errorf("db write failed"))
	mkErr, ok := err.(*meshkitErrors.Error)
	if !ok {
		t.Fatal("ErrSaveUserCredential did not return a *meshkitErrors.Error")
	}
	if mkErr.Code != ErrSaveUserCredentialCode {
		t.Errorf("expected code %s, got %s", ErrSaveUserCredentialCode, mkErr.Code)
	}
}

func TestErrGetUserCredential_IsMeshKitError(t *testing.T) {
	err := ErrGetUserCredential(fmt.Errorf("db read failed"))
	mkErr, ok := err.(*meshkitErrors.Error)
	if !ok {
		t.Fatal("ErrGetUserCredential did not return a *meshkitErrors.Error")
	}
	if mkErr.Code != ErrGetUserCredentialCode {
		t.Errorf("expected code %s, got %s", ErrGetUserCredentialCode, mkErr.Code)
	}
}

func TestErrDeleteUserCredential_IsMeshKitError(t *testing.T) {
	err := ErrDeleteUserCredential(fmt.Errorf("db delete failed"))
	mkErr, ok := err.(*meshkitErrors.Error)
	if !ok {
		t.Fatal("ErrDeleteUserCredential did not return a *meshkitErrors.Error")
	}
	if mkErr.Code != ErrDeleteUserCredentialCode {
		t.Errorf("expected code %s, got %s", ErrDeleteUserCredentialCode, mkErr.Code)
	}
}

func TestErrUpdateUserCredential_IsMeshKitError(t *testing.T) {
	err := ErrUpdateUserCredential(fmt.Errorf("db update failed"))
	mkErr, ok := err.(*meshkitErrors.Error)
	if !ok {
		t.Fatal("ErrUpdateUserCredential did not return a *meshkitErrors.Error")
	}
	if mkErr.Code != ErrUpdateUserCredentialCode {
		t.Errorf("expected code %s, got %s", ErrUpdateUserCredentialCode, mkErr.Code)
	}
}
