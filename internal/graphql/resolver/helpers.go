package resolver

import (
	"sync"

	"github.com/layer5io/meshery/internal/graphql/model"
	"github.com/layer5io/meshkit/broker"
	"github.com/layer5io/meshkit/database"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/utils"
	meshsyncmodel "github.com/layer5io/meshsync/pkg/model"
)

var (
	controlPlaneNamespace = map[model.MeshType]string{
		model.MeshTypeIstio:       "istio-system",
		model.MeshTypeLinkerd:     "linkerd-system",
		model.MeshTypeConsul:      "consul-system",
		model.MeshTypeOctarine:    "octarine-system",
		model.MeshTypeTraefikmesh: "traefik-system",
		model.MeshTypeOsm:         "osm-system",
		model.MeshTypeKuma:        "kuma-system",
		model.MeshTypeNginxsm:     "nginx-system",
		model.MeshTypeNsm:         "nsm-system",
		model.MeshTypeCitrix:      "ctrix-system",
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
func listernToEvents(log logger.Handler,
	handler *database.Handler,
	datach chan *broker.Message,
	meshsyncCh chan struct{},
) {
	var wg sync.WaitGroup
	wg.Wait()
	for {
		select {
		case msg := <-datach:
			wg.Add(1)
			go persistData(*msg, log, handler, meshsyncCh, &wg)
		}
	}
}

// persistData - scale this function with the number of events to persist
func persistData(msg broker.Message,
	log logger.Handler,
	handler *database.Handler,
	meshsyncCh chan struct{},
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	objectJSON, _ := utils.Marshal(msg.Object)
	switch msg.ObjectType {
	case broker.MeshSync:
		object := meshsyncmodel.Object{}
		err := utils.Unmarshal(string(objectJSON), &object)
		if err != nil {
			log.Error(err)
			return
		}

		// persist the object
		log.Info("Incoming object: ", object.ObjectMeta.Name, ", kind: ", object.Kind)
		err = recordMeshSyncData(msg.EventType, handler, &object)
		if err != nil {
			log.Error(err)
			return
		}
		meshsyncCh <- struct{}{}
	case broker.SMI:
		log.Info("Received SMI Result")
	}
}
