package helpers

import (
	"context"

	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/machines/grafana"
	"github.com/meshery/meshery/server/machines/kubernetes"
	"github.com/meshery/meshery/server/machines/prometheus"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/utils"
	"github.com/meshery/schemas/models/core"
)

func StatusToEvent(status connections.ConnectionStatus) machines.EventType {
	switch status {
	case connections.DISCOVERED:
		return machines.Discovery
	case connections.REGISTERED:
		return machines.Register
	case connections.CONNECTED:
		return machines.Connect
	case connections.DISCONNECTED:
		return machines.Disconnect
	case connections.IGNORED:
		return machines.Ignore
	case connections.DELETED:
		return machines.Delete
	case connections.NOTFOUND:
		return machines.NotFound
	}
	return machines.EventType(machines.DefaultState)
}

// getMachine builds the lifecycle state machine for a connection kind.
//
// Kubernetes has a bespoke machine (extra states, cluster-specific actions).
// Every other kind runs the default machine, whose transitions mirror the
// `transitionMap` authored on the kind's connection definition
// (models/.../connections/*.json): discovered -> registered/ignored,
// registered -> connected/ignored, connected -> disconnected/deleted,
// disconnected -> connected/deleted.
//
// Kinds differ only in whether the REGISTERED state carries a verification
// action. Grafana and Prometheus verify their endpoint is reachable before
// advancing; kinds without a reachability probe (e.g. artifacthub, github)
// register with no extra action and persist on connect.
//
// This deliberately does NOT allowlist kinds. The set of registerable kinds is
// the set of registered connection definitions - that is what drives the Create
// Connection wizard (buildConnectionWizardKindConfigs) - so an allowlist here
// would be a second, silently divergent source of truth: shipping a definition
// would surface the kind in the wizard and then fail registration with
// meshery-server-1218.
func getMachine(initialState machines.StateType, mtype, id string, userID core.Uuid, log logger.Handler, dbHandler *database.Handler) (*machines.StateMachine, error) {
	if mtype == "kubernetes" {
		return kubernetes.New(id, userID, log)
	}

	mch, err := machines.New(initialState, id, userID, log, mtype)
	if err != nil {
		return mch, err
	}

	if action := registerActionForKind(mtype); action != nil {
		register := mch.States[machines.REGISTERED]
		mch.States[machines.REGISTERED] = *register.RegisterAction(action)
	}

	connect := mch.States[machines.CONNECTED]
	mch.States[machines.CONNECTED] = *connect.RegisterAction(&machines.DefaultConnectAction{})

	return mch, nil
}

// registerActionForKind returns the kind-specific action run on entry to the
// REGISTERED state, or nil when the kind has no verification step.
func registerActionForKind(mtype string) machines.Action {
	switch mtype {
	case "grafana":
		return &grafana.RegisterAction{}
	case "prometheus":
		return &prometheus.RegisterAction{}
	}
	return nil
}

// HasMachineContext reports whether a state machine instance exists and carries
// an assigned Context. Callers that must resolve that Context - anything doing
// utils.Cast[*kubernetes.MachineCtx](inst.Context), or driving an action that
// does - should gate on this.
//
// It rejects two distinct shapes of a failed initialization:
//
//   - a nil instance, which InitializeMachineWithContext returns when the
//     machine could not be built at all (e.g. the cluster's API server was
//     unreachable, so the client set could not be generated and
//     AssignInitialCtx returned an error);
//   - a non-nil instance whose Context is nil. InitializeMachineWithContext
//     caches the instance via smInstanceTracker.Add *before* checking the Start
//     error, so every later call for that connection takes the cache-hit path
//     and gets back the same half-built instance, this time paired with a nil
//     error (meshery#20820). The tracker is only ever cleared by an explicit
//     user action (deleting the connection, cancelling registration), so this
//     state persists for the life of the process.
//
// Driving either shape nil-dereferences, or type-asserts a nil interface and
// logs meshkit-11180 ("nil interface cannot be type casted") - which, on the
// ~5s controller-status poll, is a log line every five seconds per broken
// connection.
//
// NOT a general "is this machine usable" test: a machine started with a nil
// InitFunc - as the connection-registration path does for non-kubernetes kinds
// - never gets a Context assigned, yet Start succeeds and the machine is
// perfectly drivable. Gating that path on this predicate would skip every such
// registration. This is for callers that genuinely require a Context.
//
// The Context test uses utils.IsInterfaceNil rather than a bare `== nil` so a
// boxed typed-nil pointer is caught too, matching the check utils.Cast performs
// internally; a bare nil-interface comparison would silently stop working if an
// InitFunc ever returned a partially-built context instead of a literal nil.
func HasMachineContext(inst *machines.StateMachine) bool {
	return inst != nil && !utils.IsInterfaceNil(inst.Context)
}

func InitializeMachineWithContext(
	machineCtx interface{},
	ctx context.Context,
	ID core.Uuid,
	userID core.Uuid,
	smInstanceTracker *machines.ConnectionToStateMachineInstanceTracker,
	log logger.Handler,
	provider models.Provider,
	initialState machines.StateType,
	mtype string,
	initFunc connections.InitFunc,
) (*machines.StateMachine, error) {
	inst, ok := smInstanceTracker.Get(ID)
	if ok {
		return inst, nil
	}

	inst, err := getMachine(initialState, mtype, ID.String(), userID, log, provider.GetGenericPersister())
	if err != nil {
		log.Error(err)
		return nil, err
	}
	inst.Provider = provider
	_, err = inst.Start(ctx, machineCtx, log, initFunc)
	smInstanceTracker.Add(ID, inst)
	if err != nil {
		return nil, err
	}

	return inst, nil
}
