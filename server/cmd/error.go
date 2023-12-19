package main

import (
	"fmt"

	"github.com/layer5io/meshkit/errors"
)

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
	ErrMarshalingRegisteryAttemptsCode            = "1014"
	ErrWritingRegisteryAttemptsCode               = "1015"
	ErrRegisteringEntityCode                      = "1016"
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
func ErrMarshalingRegisteryAttempts(err error) error {
	return errors.New(ErrMarshalingRegisteryAttemptsCode, errors.Alert, []string{"Error marshaling RegisterAttempts to JSON"}, []string{"Error marshaling RegisterAttempts to JSON: ", err.Error()}, []string{}, []string{})
}
func ErrWritingRegisteryAttempts(err error) error {
	return errors.New(ErrWritingRegisteryAttemptsCode, errors.Alert, []string{"Error writing RegisteryAttempts JSON data to file"}, []string{"Error writing RegisteryAttempts JSON data to file:", err.Error()}, []string{}, []string{})
}
func ErrRegisteringEntity(failedMsg string, hostName string) error {
	return errors.New(ErrRegisteringEntityCode, errors.Alert, []string{fmt.Sprintf("The import process for a registrant %s encountered difficulties,due to which %s. Specific issues during the import process resulted in certain entities not being successfully registered in the table.", hostName, failedMsg)}, []string{fmt.Sprintf("For registrant %s %s", hostName, failedMsg)}, []string{"Could be because of empty schema or some issue with the json or yaml file"}, []string{"Check /server/cmd/registery_attempts.json for futher details"})
}
