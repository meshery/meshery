package resolver

import (
	"context"
	"time"
)

func (r *Resolver) subscribeBrokerConnection(ctx context.Context) (<-chan bool, error) {
	ch := make(chan bool)
	previousIsConnected := r.getBrokerConnection()
	go func(outCh chan bool) {
		outCh <- r.getBrokerConnection()
		for {
			connected := r.getBrokerConnection()

			// only send message to the client if current state is different from last state
			if connected != previousIsConnected {
				// send message
				outCh <- connected

				// change reference of last to the current one for next iterations
				previousIsConnected = connected
			}
			select {
			// flush subscription
			case <-ctx.Done():
				close(outCh)

				r.Log.Info("Broker subscription flushed")
				return
			default:
			}
			time.Sleep(10 * time.Second)
		}
	}(ch)

	return ch, nil
}

func (r *Resolver) getBrokerConnection() bool {
	connected := !r.BrokerConn.IsEmpty()

	if connected {
		err := r.BrokerConn.Publish("connection.test", nil)
		if err != nil {
			connected = false
		}
	}

	return connected
}
