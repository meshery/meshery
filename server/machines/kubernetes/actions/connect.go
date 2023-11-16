package actions

import (
	"github.com/layer5io/meshery/server/machines"
)

type ConnectAction struct {}

func(ca *ConnectAction) ExecuteOnEntry(ctx interface{}) (machines.EventType, error) {
	var mt machines.EventType
	return mt, nil
}
func(ca *ConnectAction) Execute(ctx interface{}) (machines.State, machines.EventType, error) {
	var mt machines.EventType
	var state machines.State
	return state, mt, nil
}

func(ca *ConnectAction) ExecuteOnExit(ctx interface{}) (machines.EventType, error) {
	var mt machines.EventType
	return mt, nil
}
