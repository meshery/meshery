package main

import "github.com/layer5io/meshkit/errors"

// Please reference the following before contributing an error code:
// https://docs.meshery.io/project/contributing/contributing-error
// https://github.com/meshery/meshkit/blob/master/errors/errors.go
const (
	ErrCreatingUUIDInstanceCode                   = "1001"
	ErrRegisteringMesheryOAMTraitsCode            = "1002"
	ErrRegisteringMesheryOAMWorkloadsCode         = "1003"
	ErrRetrievingUserHomeDirectoryCode            = "1004"
	ErrCreatingUserDataDirectoryCode              = "1005"
	ErrCreatingMapPreferencePersisterInstanceCode = "1006"
	ErrDatabaseAutoMigrationCode                  = "1008"
	ErrInvalidURLSkippingProviderCode             = "1009"
	ErrListenAndServeCode                         = "1010"
	ErrCleaningUpLocalProviderCode                = "1011"
	ErrClosingDatabaseInstanceCode                = "1012"
	ErrInitializingRegistryManagerCode            = "1013"
)

func ErrInitializingRegistryManager(err error) error {
	return errors.New(ErrInitializingRegistryManagerCode, errors.Fatal, []string{"could not initialize registry manager"}, []string{err.Error()}, []string{"could not migrate tables into the database"}, []string{"make sure the database instance passed is not nil"})
}

func ErrCreatingUUIDInstance(err error) error {
	return errors.New(ErrCreatingUUIDInstanceCode, errors.Fatal, []string{"Unable to create UUID instance"}, []string{"Unable to create UUID instance: ", err.Error()}, []string{}, []string{})
}

func ErrRegisteringMesheryOAMTraits(err error) error {
	return errors.New(ErrRegisteringMesheryOAMTraitsCode, errors.Alert, []string{"Error registering local OAM traits"}, []string{"Error registering local OAM traits: ", err.Error()}, []string{}, []string{})
}

func ErrRegisteringMesheryOAMWorkloads(err error) error {
	return errors.New(ErrRegisteringMesheryOAMWorkloadsCode, errors.Alert, []string{"Error registering local OAM workloads"}, []string{"Error registering local OAM workloads: ", err.Error()}, []string{}, []string{})
}

func ErrRetrievingUserHomeDirectory(err error) error {
	return errors.New(ErrRetrievingUserHomeDirectoryCode, errors.Fatal, []string{"Unable to retrieve the user's home directory"}, []string{"Unable to retrieve the user's home directory: ", err.Error()}, []string{}, []string{})
}

func ErrCreatingUserDataDirectory(dir string) error {
	return errors.New(ErrCreatingUserDataDirectoryCode, errors.Fatal, []string{"Unable to create the directory for storing user data at: ", dir}, []string{"Unable to create the directory for storing user data at: ", dir}, []string{}, []string{})
}

func ErrCreatingMapPreferencePersisterInstance(err error) error {
	return errors.New(ErrCreatingMapPreferencePersisterInstanceCode, errors.Fatal, []string{"Unable to create a new MapPreferencePersister instance"}, []string{"Unable to create a new MapPreferencePersister instance: ", err.Error()}, []string{}, []string{})
}

func ErrDatabaseAutoMigration(err error) error {
	return errors.New(ErrDatabaseAutoMigrationCode, errors.Fatal, []string{"Unable to auto migrate to database"}, []string{"Unable to auto migrate to database: ", err.Error()}, []string{}, []string{})
}

func ErrInvalidURLSkippingProvider(url string) error {
	return errors.New(ErrInvalidURLSkippingProviderCode, errors.Alert, []string{url, " is invalid url skipping provider"}, []string{url, " is invalid url skipping provider"}, []string{}, []string{})
}

func ErrListenAndServe(err error) error {
	return errors.New(ErrListenAndServeCode, errors.Fatal, []string{"ListenAndServe Error"}, []string{"ListenAndServe Error: ", err.Error()}, []string{}, []string{})
}

func ErrCleaningUpLocalProvider(err error) error {
	return errors.New(ErrCleaningUpLocalProviderCode, errors.Alert, []string{"Error cleaning up local provider"}, []string{"Error cleaning up local provider: ", err.Error()}, []string{}, []string{})
}

func ErrClosingDatabaseInstance(err error) error {
	return errors.New(ErrClosingDatabaseInstanceCode, errors.Alert, []string{"Error closing database instance"}, []string{"Error closing database instance: ", err.Error()}, []string{}, []string{})
}
