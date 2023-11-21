package machines

import "github.com/layer5io/meshery/server/models/connections"

func StatusToEvent(status connections.ConnectionStatus) EventType {
	switch status {
	case connections.DISCOVERED:
		return Discovery
	case connections.REGISTERED:
		return Register
	case connections.CONNECTED:
		return Connect
	case connections.DISCONNECTED:
		return Disconnect
	case connections.IGNORED:
		return Ignore
	case connections.DELETED:
		return Delete
	case connections.NOTFOUND:
		return NotFound
	}
	return EventType(DefaultState)
}
