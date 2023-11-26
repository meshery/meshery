package machines

// Represents an state in the system/machine
type StateType string

type State struct {
	Events Events
	Action Action
}

// Represents mapping between state name and the state
type States map[StateType]State

func (s *State) RegisterEvent(eventType EventType, stateType StateType) *State {
	if s.Events == nil {
		s.Events = make(Events)
	}
	s.Events[eventType] = stateType
	return s
}

func (s *State) RegisterAction(action Action) *State {
	s.Action = action
	return s
}
