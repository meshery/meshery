package machines

import (
	"fmt"

	"github.com/layer5io/meshery/server/machines/grafana"
	"github.com/layer5io/meshery/server/machines/kubernetes"
	"github.com/layer5io/meshery/server/machines/prometheus"
	"github.com/layer5io/meshery/server/models/connections"
	"github.com/layer5io/meshery/server/models/machines"
	"github.com/layer5io/meshkit/logger"
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
		mch, err := New(initialState, id, log, mtype)
		if err != nil {
			return mch, err
		}
		register := mch.States[machines.REGISTERED]
		register.RegisterAction(&grafana.RegisterAction{})

		connect := mch.States[machines.CONNECTED]
		connect.RegisterAction(&machines.DefaultConnectAction{})
		return mch, nil
	case "prometheus":
		mch, err := New(initialState, id, log, mtype)
		if err != nil {
			return mch, err
		}
		register := mch.States[machines.REGISTERED]
		register.RegisterAction(&prometheus.RegisterAction{})

		connect := mch.States[machines.CONNECTED]
		connect.RegisterAction(&machines.DefaultConnectAction{})

		return mch, nil
	}
	return nil, machines.ErrInvalidType(fmt.Errorf("invlaid type requested"))
}
