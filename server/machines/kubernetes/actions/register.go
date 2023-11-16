package actions

import (
	"github.com/layer5io/meshery/server/machines"
)

type RegisterAction struct {}

func(ra *RegisterAction) ExecuteOnEntry(ctx interface{}) (machines.EventType, error) {
	var mt machines.EventType
	return mt, nil
}
func(ra *RegisterAction) Execute(ctx interface{}) (machines.State, machines.EventType, error) {
	var mt machines.EventType
	var state machines.State
	return state, mt, nil
}

func(ra *RegisterAction) ExecuteOnExit(ctx interface{}) (machines.EventType, error) {
	var mt machines.EventType
	return mt, nil
}
