package machines

import (
	"fmt"

	"github.com/layer5io/meshery/server/machines/grafana"
	"github.com/layer5io/meshery/server/machines/kubernetes"
	"github.com/layer5io/meshery/server/machines/prometheus"
	"github.com/layer5io/meshery/server/models/connections"
	"github.com/layer5io/meshery/server/models/machines"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/utils"
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

func GetMachine(initialState machines.StateType, mtype, id string, log logger.Handler) (*machines.StateMachine, error) {
	switch mtype {
		case "kubernetes":
			return kubernetes.New(id, log)
		case "grafana":
			return grafana.New(initialState, id, log)
		case "prometheus":
			t, err := prometheus.New(initialState, id, log)
			fmt.Println("inside GetMachine", t)
			return t, err
	}
	return nil, machines.ErrInvalidType(fmt.Errorf("invlaid type requested"))
}

func MarshalAndUnmarshal[k any, v any](val k) (unmarshalledvalue v, err error){
	data, err := utils.Marshal(val)
	if err != nil {
		return
	}

	err = utils.Unmarshal(data, &unmarshalledvalue)
	if err != nil {
		return
	}
	return
}
