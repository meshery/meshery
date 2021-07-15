// Package models - Error codes for Meshery models
package models

import (
	"github.com/layer5io/meshkit/errors"
)

const (
	ErrGrafanaClientCode          = "3000"
	ErrPageSizeCode               = "replace_me"
	ErrPageNumberCode             = "replace_me"
	ErrResultIDCode               = "replace_me"
	ErrPerfIDCode                 = "replace_me"
	ErrMarshalCode                = "replace_me"
	ErrUnmarshalCode              = "replace_me"
	ErrGenerateUUIDCode           = "replace_me"
	ErrLocalProviderSupportCode   = "replace_me"
	ErrGrafanaOrgCode             = "replace_me"
	ErrGrafanaBoardsCode          = "replace_me"
	ErrGrafanaDashboardCode       = "replace_me"
	ErrGrafanaDataSourceCode      = "replace_me"
	ErrNilQueryCode               = "replace_me"
	ErrGrafanaDataCode            = "replace_me"
	ErrApplicationFileNameCode    = "replace_me"
	ErrFilterFileNameCode         = "replace_me"
	ErrPatternFileNameCode        = "replace_me"
	ErrMakeDirCode                = "replace_me"
	ErrFolderStatCode             = "replace_me"
	ErrUserIDCode                 = "replace_me"
	ErrDBConnectionCode           = "replace_me"
	ErrNilConfigDataCode          = "replace_me"
	ErrDBOpenCode                 = "replace_me"
	ErrDBRLockCode                = "replace_me"
	ErrDBLockCode                 = "replace_me"
	ErrDBReadCode                 = "replace_me"
	ErrDBDeleteCode               = "replace_me"
	ErrCopyCode                   = "replace_me"
	ErrDBPutCode                  = "replace_me"
	ErrResultDataCode             = "replace_me"
	ErrUnableToPersistsResultCode = "replace_me"
)

var (
	ErrResultID             = errors.New(ErrResultIDCode, errors.Alert, []string{"Given resultID is not valid"}, []string{"Given resultID is nil"}, []string{}, []string{})
	ErrLocalProviderSupport = errors.New(ErrLocalProviderSupportCode, errors.Alert, []string{"Method not supported by local provider"}, []string{}, []string{}, []string{})
	ErrNilQuery             = errors.New(ErrNilQueryCode, errors.Alert, []string{"Query data passed is nil"}, []string{}, []string{}, []string{})
	ErrApplicationFileName  = errors.New(ErrApplicationFileNameCode, errors.Alert, []string{"Invalid Applicationfile"}, []string{"Name field is either not present or is not valid"}, []string{}, []string{})
	ErrFilterFileName       = errors.New(ErrFilterFileNameCode, errors.Alert, []string{"Invalid Filterfile"}, []string{"Name field is either not present or is not valid"}, []string{}, []string{})
	ErrPatternFileName      = errors.New(ErrPatternFileNameCode, errors.Alert, []string{"Invalid Patternfile"}, []string{"Name field is either not present or is not valid"}, []string{}, []string{})
	ErrUserID               = errors.New(ErrUserIDCode, errors.Alert, []string{"User ID is empty"}, []string{}, []string{}, []string{})
	ErrDBConnection         = errors.New(ErrDBConnectionCode, errors.Alert, []string{"Connection to DataBase does not exist"}, []string{}, []string{}, []string{})
	ErrNilConfigData        = errors.New(ErrNilConfigDataCode, errors.Alert, []string{"Given config data is nil"}, []string{}, []string{}, []string{})
)

func ErrGrafanaClient(err error) error {
	return errors.New(ErrGrafanaClientCode, errors.Alert, []string{"Unable to initialize Grafana Client"}, []string{"Unable to initializes client for interacting with an instance of Grafana server", err.Error()}, []string{"Invalid Grafana Endpoint or API-Key"}, []string{"Update your Grafana URL and API-Key from the settings page in the UI"})
}

func ErrPageSize(err error) error {
	return errors.New(ErrPageSizeCode, errors.Alert, []string{"Unable to prase the Page Size"}, []string{err.Error()}, []string{}, []string{})
}

func ErrPageNumber(err error) error {
	return errors.New(ErrPageNumberCode, errors.Alert, []string{"Unable to prase the Page Numer"}, []string{err.Error()}, []string{}, []string{})
}

func ErrPerfID(err error) error {
	return errors.New(ErrPerfIDCode, errors.Alert, []string{"Invalid peformance profile ID"}, []string{err.Error()}, []string{}, []string{})
}

func ErrMarshal(err error, obj string) error {
	return errors.New(ErrMarshalCode, errors.Alert, []string{"Unable to marshal the object", obj}, []string{err.Error()}, []string{}, []string{})
}

func ErrUnmarshal(err error, obj string) error {
	return errors.New(ErrUnmarshalCode, errors.Alert, []string{"Unable to unmarshal the object", obj}, []string{err.Error()}, []string{}, []string{})
}

func ErrResultData() error {
	return errors.New(ErrResultDataCode, errors.Alert, []string{"given result data is nil"}, []string{}, []string{}, []string{})
}

func ErrUnableToPersistsResult(err error) error {
	return errors.New(ErrUnableToPersistsResultCode, errors.Alert, []string{"unable to persists the result data"}, []string{err.Error()}, []string{}, []string{})
}

func ErrGenerateUUID(err error) error {
	return errors.New(ErrGenerateUUIDCode, errors.Alert, []string{"Unable to generate a new UUID"}, []string{err.Error()}, []string{}, []string{})
}

func ErrGrafanaOrg(err error) error {
	return errors.New(ErrGrafanaOrgCode, errors.Alert, []string{"Failed to get Org data from Grafana"}, []string{err.Error()}, []string{"Invalid URL", "Invalid API-Key"}, []string{})
}

func ErrGrafanaBoards(err error) error {
	return errors.New(ErrGrafanaBoardsCode, errors.Alert, []string{"Unable to get Grafana Boards"}, []string{err.Error()}, []string{}, []string{})
}

func ErrGrafanaDashboard(err error, UID string) error {
	return errors.New(ErrGrafanaDashboardCode, errors.Alert, []string{"Error getting grafana dashboard from UID", UID}, []string{err.Error()}, []string{}, []string{})
}

func ErrGrafanaDataSource(err error, ds string) error {
	return errors.New(ErrGrafanaDataSourceCode, errors.Alert, []string{"Error getting Grafana Board's Datasource", ds}, []string{err.Error()}, []string{}, []string{})
}

func ErrGrafanaData(err error, apiEndpoint string) error {
	return errors.New(ErrGrafanaDataCode, errors.Alert, []string{"Error getting data from Grafana API", apiEndpoint}, []string{err.Error()}, []string{}, []string{})
}

func ErrMakeDir(err error, dir string) error {
	return errors.New(ErrMakeDirCode, errors.Alert, []string{"Unable to create directory/folder", dir}, []string{err.Error()}, []string{}, []string{})
}

func ErrFolderStat(err error, dir string) error {
	return errors.New(ErrFolderStatCode, errors.Alert, []string{"Unable to find (os.stat) the folder", dir}, []string{err.Error()}, []string{}, []string{})
}

func ErrDBOpen(err error) error {
	return errors.New(ErrDBOpenCode, errors.Alert, []string{"Unable to open the database"}, []string{err.Error()}, []string{}, []string{})
}

func ErrDBRLock(err error) error {
	return errors.New(ErrDBRLockCode, errors.Alert, []string{"Unable to obtain read lock from bitcask store"}, []string{err.Error()}, []string{}, []string{})
}

func ErrDBLock(err error) error {
	return errors.New(ErrDBLockCode, errors.Alert, []string{"Unable to obtain write lock from bitcask store"}, []string{err.Error()}, []string{}, []string{})
}

func ErrDBRead(err error) error {
	return errors.New(ErrDBReadCode, errors.Alert, []string{"Unable to read data from bitcast store"}, []string{err.Error()}, []string{}, []string{})
}

func ErrDBPut(err error) error {
	return errors.New(ErrDBPutCode, errors.Alert, []string{"Unable to Persist config data."}, []string{err.Error()}, []string{}, []string{})
}

func ErrDBDelete(err error, user string) error {
	return errors.New(ErrDBDeleteCode, errors.Alert, []string{"Unable to delete config data for the user", user}, []string{err.Error()}, []string{}, []string{})
}

func ErrCopy(err error, obj string) error {
	return errors.New(ErrCopyCode, errors.Alert, []string{"Error occurred while copying", obj}, []string{err.Error()}, []string{}, []string{})
}
