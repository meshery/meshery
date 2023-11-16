package actions

import (
	"github.com/layer5io/meshery/server/machines"
)

type DisconnectAction struct {}

func(da *DisconnectAction) ExecuteOnEntry(ctx interface{}) (machines.EventType, error) {
	var mt machines.EventType
	return mt, nil
}
func(da *DisconnectAction) Execute(ctx interface{}) (machines.State, machines.EventType, error) {
	var mt machines.EventType
	var state machines.State
	return state, mt, nil
}

func(da *DisconnectAction) ExecuteOnExit(ctx interface{}) (machines.EventType, error) {
	var mt machines.EventType
	return mt, nil
}
