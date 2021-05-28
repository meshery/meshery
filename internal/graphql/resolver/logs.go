package resolver

import (
	"context"
	"fmt"

	"github.com/layer5io/meshery/internal/graphql/model"
	"github.com/layer5io/meshkit/broker"
	"github.com/layer5io/meshkit/utils"
	meshsyncmodel "github.com/layer5io/meshsync/pkg/model"
)

func (r *subscriptionResolver) subscribeLogStream(ctx context.Context, selector []*model.LogStreamRequest) (<-chan *model.LogStream, error) {
	if r.logsChannel == nil {
		r.logsChannel = make(chan *model.LogStream)
	}

	go func() {
		r.Log.Info("Logs subscription started")
		for _, req := range selector {
			sig := false
			if req.Signal == model.SignalStop {
				sig = true
			}

			err := r.brokerConn.Publish("meshery.meshsync.request", &broker.Message{
				ObjectType: broker.ObjectType(broker.Request),
				EventType:  broker.EventType(broker.Add),
				Request: &broker.RequestObject{
					Entity: broker.LogRequestEntity,
					Payload: meshsyncmodel.LogRequests{
						"operator": meshsyncmodel.LogRequest{
							ID:        req.ID,
							Name:      req.PodName,
							Namespace: req.Namespace,
							Container: req.ContainerName,
							Follow:    true,
							Previous:  false,
							TailLines: 1000,
							Stop:      sig,
						},
					},
				},
			})
			if err != nil {
				r.Log.Error(err)
			}
		}
	}()

	return r.logsChannel, nil
}

func (r *Resolver) processLogs(obj interface{}) {
	nobj := meshsyncmodel.LogObject{}
	d, err := utils.Marshal(obj)
	if err != nil {
		r.logsChannel <- &model.LogStream{
			Error: err.Error(),
		}
	}

	err = utils.Unmarshal(d, &nobj)
	if err != nil {
		r.logsChannel <- &model.LogStream{
			Error: err.Error(),
		}
	}

	if r.logsChannel == nil {
		fmt.Print("Channel closed")
	}
	r.logsChannel <- &model.LogStream{
		ID:   nobj.ID,
		Data: nobj.Data,
	}
}
