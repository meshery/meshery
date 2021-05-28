package resolver

import (
	"sync"

	"github.com/layer5io/meshery/internal/graphql/model"
	"github.com/layer5io/meshkit/broker"
	"github.com/layer5io/meshkit/utils"
	meshsyncmodel "github.com/layer5io/meshsync/pkg/model"
)

var (
	controlPlaneNamespace = map[model.MeshType]string{
		model.MeshTypeIstio:              "istio-system",
		model.MeshTypeLinkerd:            "linkerd-system",
		model.MeshTypeConsul:             "consul-system",
		model.MeshTypeOctarine:           "octarine-system",
		model.MeshTypeTraefikMesh:        "traefik-system",
		model.MeshTypeOpenServiceMesh:    "osm-system",
		model.MeshTypeKuma:               "kuma-system",
		model.MeshTypeNginxServiceMesh:   "nginx-system",
		model.MeshTypeNetworkServiceMesh: "nsm-system",
		model.MeshTypeCitrixServiceMesh:  "ctrix-system",
	}

	addonPortSelector = map[string]string{
		"grafana":          "service",
		"prometheus":       "http",
		"jaeger-collector": "jaeger-collector-http",
		"kiali":            "http",
		"zipkin":           "http-query",
	}
)

// listernToEvents - scale this function with the number of channels
func (r *Resolver) listernToEvents() {
	var wg sync.WaitGroup
	wg.Wait()
	for {
		select {
		case msg := <-r.brokerChannel:
			wg.Add(1)
			switch msg.ObjectType {
			case broker.MeshSync:
				go r.persistData(*msg, &wg)
			case broker.LogStreamObject:
				go r.processLogs(msg.Object)
			}
		}
	}
}

// persistData - scale this function with the number of events to persist
func (r *Resolver) persistData(msg broker.Message,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	objectJSON, _ := utils.Marshal(msg.Object)
	switch msg.ObjectType {
	case broker.MeshSync:
		object := meshsyncmodel.Object{}
		err := utils.Unmarshal(string(objectJSON), &object)
		if err != nil {
			r.Log.Error(err)
			return
		}

		// persist the object
		r.Log.Info("Incoming object: ", object.ObjectMeta.Name, ", kind: ", object.Kind)
		err = recordMeshSyncData(msg.EventType, r.DBHandler, &object)
		if err != nil {
			r.Log.Error(err)
			return
		}
		r.MeshSyncChannel <- struct{}{}
	case broker.SMI:
		r.Log.Info("Received SMI Result")
	}
}
