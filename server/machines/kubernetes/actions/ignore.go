package actions

import (
	"github.com/layer5io/meshery/server/machines"
)

type IgnoreAction struct {}

func(ia *IgnoreAction) ExecuteOnEntry(ctx interface{}) (machines.EventType, error) {
	var mt machines.EventType
	return mt, nil
}
func(ia *IgnoreAction) Execute(ctx interface{}) (machines.State, machines.EventType, error) {
	var mt machines.EventType
	var state machines.State
	return state, mt, nil
}

func(ia *IgnoreAction) ExecuteOnExit(ctx interface{}) (machines.EventType, error) {
	var mt machines.EventType
	return mt, nil
}
