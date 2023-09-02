// Package handlers :collection of handlers (aka "HTTP middleware")
package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"encoding/json"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/meshes"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/errors"
	_events "github.com/layer5io/meshkit/models/events"
	"github.com/layer5io/meshkit/utils/events"

	"github.com/sirupsen/logrus"
)

var (
	flusherMap map[string]http.Flusher
)

// swagger:route GET /api/events EventsAPI idGetEventStreamer
// Handle GET request for events.
// Listens for events across all of Meshery's components like adapters and server, streaming them to the UI via Server Side Events
// This API call never terminates and establishes a persistent keep-alive connection over which `EventsResponse`s are pushed.
// responses:
// 	200:

// EventStreamHandler endpoint is used for streaming events to the frontend
func (h *Handler) EventStreamHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, p models.Provider) {
	// if req.Method != http.MethodGet {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	log := logrus.WithField("file", "events_streamer")
	client := "ui"
	if req.URL.Query().Get("client") != "" {
		client = req.URL.Query().Get("client")
	}

	if flusherMap == nil {
		flusherMap = make(map[string]http.Flusher, 0)
	}

	flusher, ok := w.(http.Flusher)
	flusherMap[client] = flusher

	if !ok {
		log.Error("Event streaming not supported.")
		http.Error(w, "Event streaming is not supported at the moment.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	notify := req.Context()

	var err error

	localMeshAdapters := map[string]*meshes.MeshClient{}
	localMeshAdaptersLock := &sync.Mutex{}

	respChan := make(chan []byte, 100)

	newAdaptersChan := make(chan *meshes.MeshClient)

	go func() {
		for mClient := range newAdaptersChan {
			log.Debug("received a new mesh client, listening for events")
			go func(mClient *meshes.MeshClient) {
				listenForAdapterEvents(req.Context(), mClient, respChan, log, p,  h.config.EventsChannel, *h.SystemID, user.ID)
				_ = mClient.Close()
			}(mClient)
		}

		log.Debug("new adapters channel closed")
	}()
	go listenForCoreEvents(req.Context(), h.EventsBuffer, respChan, log, p)
	go func(flusher http.Flusher) {
		for data := range respChan {
			log.Debug("received new data on response channel")
			_, _ = fmt.Fprintf(w, "data: %s\n\n", data)
			if flusher != nil {
				flusher.Flush()
				log.Debugf("Flushed the messages on the wire...")
			}
		}
		log.Debug("response channel closed")
	}(flusherMap[client])

STOP:
	for {
		select {
		case <-notify.Done():
			log.Debugf("received signal to close connection and channels")
			close(newAdaptersChan)
			break STOP
		default:
			meshAdapters := prefObj.MeshAdapters
			if meshAdapters == nil {
				meshAdapters = []*models.Adapter{}
			}

			adaptersLen := len(meshAdapters)
			if adaptersLen == 0 {
				// Clear the adapter cache
				localMeshAdapters = closeAdapterConnections(localMeshAdaptersLock, localMeshAdapters)
			} else {
				localMeshAdaptersLock.Lock()
				for _, ma := range meshAdapters {
					mClient, ok := localMeshAdapters[ma.Location]
					if !ok {
						mClient, err = meshes.CreateClient(req.Context(), ma.Location)
						if err == nil {
							localMeshAdapters[ma.Location] = mClient
						}
					}
					if mClient != nil {
						_, err = mClient.MClient.MeshName(req.Context(), &meshes.MeshNameRequest{})
						if err != nil {
							_ = mClient.Close()
							delete(localMeshAdapters, ma.Location)
						} else {
							if !ok { // reusing the map check, only when ok is false a new entry will be added
								newAdaptersChan <- mClient
							}
						}
					}
				}
				localMeshAdaptersLock.Unlock()
			}
		}
		time.Sleep(5 * time.Second)
	}
	close(respChan)
	defer log.Debug("events handler closed")
}
func listenForCoreEvents(ctx context.Context, eb *events.EventStreamer, resp chan []byte, log *logrus.Entry, _ models.Provider) {
	datach := make(chan interface{}, 10)
	go eb.Subscribe(datach)
	for {
		select {
		case datap := <-datach:
			event, ok := datap.(*meshes.EventsResponse)
			if !ok {
				continue
			}
			data, err := json.Marshal(event)
			if err != nil {
				log.Error(models.ErrMarshal(err, "event"))
				continue
			}
			resp <- data

		case <-ctx.Done():
			return
		}
	}
}
func listenForAdapterEvents(ctx context.Context, mClient *meshes.MeshClient, respChan chan []byte, log *logrus.Entry, p models.Provider, ec *models.Signal, systemID uuid.UUID, userID string) {
	log.Debugf("Received a stream client...")
	userUUID := uuid.FromStringOrNil(userID)
	streamClient, err := mClient.MClient.StreamEvents(ctx, &meshes.EventsRequest{})
	if err != nil {
		log.Error(ErrStreamEvents(err))
		// errChan <- err
		// http.Error(w, "There was an error connecting to the backend to get events", http.StatusInternalServerError)
		return
	}

	for {
		log.Debugf("Waiting to receive events.")
		event, err := streamClient.Recv()
		if err != nil {
			if err == io.EOF {
				log.Error(ErrStreamClient(err))
				return
			}
			log.Error(ErrStreamClient(err))
			return
		}
		// log.Debugf("received an event: %+#v", event)
		log.Debugf("Received an event.")
		eventType := event.EventType.String()
		eventBuilder := _events.NewEvent().FromSystem(uuid.FromStringOrNil(event.Component)).
		WithSeverity(_events.EventSeverity(eventType)).WithDescription(event.Summary).WithCategory(event.ComponentName).WithAction("deploy").FromUser(userUUID)
		if strings.Contains(event.Summary, "removed") {
			eventBuilder.WithAction("undeploy")
		}

		if strings.Contains(event.Summary, "Smi conformance test") {
			result := &models.SmiResult{}
			err := json.Unmarshal([]byte(event.Details), result)
			if err != nil {
				log.Error(models.ErrUnmarshal(err, "event"))
				return
			}

			id, err := p.PublishSmiResults(result)
			if err != nil {
				log.Error(ErrPublishSmiResults(err))
				return
			}
			event.Details = fmt.Sprintf("Result-Id: %s", id)
		}
		if eventType == "ERROR" {
			err := errors.New(event.ErrorCode, errors.Alert, []string{event.Summary}, []string{event.Details}, []string{event.ProbableCause}, []string{event.SuggestedRemediation})
			eventBuilder.WithMetadata(map[string]interface{}{
				"error": err,
			})
		}

		_event := eventBuilder.Build()
		_ = p.PersistEvent(_event)
		ec.Publish(userUUID, _event)

		data, err := json.Marshal(event)
		if err != nil {
			log.Error(models.ErrMarshal(err, "event"))
			return
		}
		respChan <- data
	}
}

func closeAdapterConnections(localMeshAdaptersLock *sync.Mutex, localMeshAdapters map[string]*meshes.MeshClient) map[string]*meshes.MeshClient {
	localMeshAdaptersLock.Lock()
	for _, mcl := range localMeshAdapters {
		_ = mcl.Close()
	}
	localMeshAdaptersLock.Unlock()

	return map[string]*meshes.MeshClient{}
}
