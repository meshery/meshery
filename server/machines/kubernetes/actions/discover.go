package actions

import "github.com/layer5io/meshery/server/machines"


type DiscoverAction struct {}

func(da *DiscoverAction) ExecuteOnEntry(ctx interface{}) (machines.EventType, error) {
	var mt machines.EventType
	return mt, nil
}
func(da *DiscoverAction) Execute(ctx interface{}) (machines.State, machines.EventType, error) {
	var mt machines.EventType
	var state machines.State
	return state, mt, nil
}

func(da *DiscoverAction) ExecuteOnExit(ctx interface{}) (machines.EventType, error) {
	var mt machines.EventType
	return mt, nil
}
