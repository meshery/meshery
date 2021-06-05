package resolver

import (
	"github.com/layer5io/meshkit/errors"
)

const (
	ErrInvalidRequestCode           = "1000"
	ErrNilClientCode                = "1001"
	ErrCreateDataCode               = "1002"
	ErrQueryCode                    = "1003"
	ErrMeshsyncSubscriptionCode     = "1004"
	ErrOperatorSubscriptionCode     = "1005"
	ErrAddonSubscriptionCode        = "1006"
	ErrControlPlaneSubscriptionCode = "1007"
	ErrMesheryClientCode            = "1008"
	ErrSubscribeChannelCode         = "1009"
	ErrNoMeshSyncCode               = "1010"
)

var (
	ErrNilClient      = errors.New(ErrNilClientCode, errors.Alert, []string{"Kubernetes client not initialized"}, []string{}, []string{}, []string{})
	ErrInvalidRequest = errors.New(ErrInvalidRequestCode, errors.Alert, []string{"Invalid query, please check syntax"}, []string{}, []string{}, []string{})
	ErrNoMeshSync     = errors.New(ErrNoMeshSyncCode, errors.Alert, []string{"MeshSync disabled"}, []string{}, []string{}, []string{})
	ErrEmptyHandler   = errors.New(ErrNoMeshSyncCode, errors.Alert, []string{"Database handler not initialized"}, []string{}, []string{}, []string{})
)

func ErrCreateData(err error) error {
	return errors.New(ErrCreateDataCode, errors.Alert, []string{"Error while writing meshsync data", err.Error()}, []string{}, []string{}, []string{})
}

func ErrUpdateData(err error) error {
	return errors.New(ErrCreateDataCode, errors.Alert, []string{"Error while updating meshsync data", err.Error()}, []string{}, []string{}, []string{})
}

func ErrDeleteData(err error) error {
	return errors.New(ErrCreateDataCode, errors.Alert, []string{"Error while deleting meshsync data", err.Error()}, []string{}, []string{}, []string{})
}

func ErrQuery(err error) error {
	return errors.New(ErrQueryCode, errors.Alert, []string{"Error while querying data", err.Error()}, []string{}, []string{}, []string{})
}

func ErrMeshsyncSubscription(err error) error {
	return errors.New(ErrMeshsyncSubscriptionCode, errors.Alert, []string{"MeshSync Subscription failed", err.Error()}, []string{}, []string{}, []string{})
}

func ErrOperatorSubscription(err error) error {
	return errors.New(ErrOperatorSubscriptionCode, errors.Alert, []string{"Operator Subscription failed", err.Error()}, []string{}, []string{}, []string{})
}

func ErrAddonSubscription(err error) error {
	return errors.New(ErrAddonSubscriptionCode, errors.Alert, []string{"Addons Subscription failed", err.Error()}, []string{}, []string{}, []string{})
}

func ErrControlPlaneSubscription(err error) error {
	return errors.New(ErrControlPlaneSubscriptionCode, errors.Alert, []string{"Control Plane Subscription failed", err.Error()}, []string{}, []string{}, []string{})
}

func ErrSubscribeChannel(err error) error {
	return errors.New(ErrSubscribeChannelCode, errors.Alert, []string{"Unable to subscribe to channel", err.Error()}, []string{}, []string{}, []string{})
}

func ErrMesheryClient(err error) error {
	if err != nil {
		return errors.New(ErrMesheryClientCode, errors.Alert, []string{"Meshery kubernetes client not initialized", err.Error()}, []string{}, []string{}, []string{})
	}
	return errors.New(ErrMesheryClientCode, errors.Alert, []string{"Meshery kubernetes client not initialized"}, []string{}, []string{}, []string{})
}
