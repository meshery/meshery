//Package handlers :collection of handlers (aka "HTTP middleware")
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

	"github.com/layer5io/meshery/meshes"
	"github.com/layer5io/meshery/models"
	"github.com/sirupsen/logrus"
)

// EventStreamHandler endpoint is used for streaming events to the frontend
func (h *Handler) EventStreamHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, p models.Provider) {
	// if req.Method != http.MethodGet {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	log := logrus.WithField("file", "events_streamer")

	flusher, ok := w.(http.Flusher)

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
				listenForAdapterEvents(req.Context(), mClient, respChan, log, p)
				_ = mClient.Close()
			}(mClient)
		}

		log.Debug("new adapters channel closed")
	}()

	go func() {
		for data := range respChan {
			log.Debug("received new data on response channel")
			_, _ = fmt.Fprintf(w, "data: %s\n\n", data)
			if flusher != nil {
				flusher.Flush()
				log.Debugf("Flushed the messages on the wire...")
			}
		}
		log.Debug("response channel closed")
	}()

STOP:
	for {
		select {
		case <-notify.Done():
			log.Debugf("received signal to close connection and channels")
			close(newAdaptersChan)
			close(respChan)
			break STOP
		default:
			meshAdapters := prefObj.MeshAdapters
			if meshAdapters == nil {
				meshAdapters = []*models.Adapter{}
			}

			if prefObj.K8SConfig == nil || !prefObj.K8SConfig.InClusterConfig && (prefObj.K8SConfig.Config == nil || len(prefObj.K8SConfig.Config) == 0) {
				log.Debug("No valid Kubernetes config found.") // switching from Error to Debug to prevent it from filling up the logs
				// http.Error(w, `No valid Kubernetes config found.`, http.StatusBadRequest)
				// return
				localMeshAdaptersLock.Lock()
				for _, mcl := range localMeshAdapters {
					_ = mcl.Close()
				}
				localMeshAdapters = map[string]*meshes.MeshClient{}
				localMeshAdaptersLock.Unlock()
			} else {
				adaptersLen := len(meshAdapters)
				if adaptersLen == 0 {
					log.Debug("No valid mesh adapter(s) found.") // switching from Error to Debug to prevent it from filling up the logs
					// http.Error(w, `No valid mesh adapter(s) found.`, http.StatusBadRequest)
					// return
					localMeshAdaptersLock.Lock()
					for _, mcl := range localMeshAdapters {
						_ = mcl.Close()
					}
					localMeshAdapters = map[string]*meshes.MeshClient{}
					localMeshAdaptersLock.Unlock()
				} else {
					localMeshAdaptersLock.Lock()
					for _, ma := range meshAdapters {
						mClient, ok := localMeshAdapters[ma.Location]
						if !ok {
							mClient, err = meshes.CreateClient(req.Context(), prefObj.K8SConfig.Config, prefObj.K8SConfig.ContextName, ma.Location)
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
		}
		time.Sleep(5 * time.Second)
	}
	defer log.Debug("events handler closed")
}

func listenForAdapterEvents(ctx context.Context, mClient *meshes.MeshClient, respChan chan []byte, log *logrus.Entry, p models.Provider) {
	log.Debugf("Received a stream client...")

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
		if strings.Contains(event.Summary, "Smi conformance test") {
			result := &models.SmiResult{}
			err := json.Unmarshal([]byte(event.Details), result)
			if err != nil {
				log.Error(ErrUnmarshal(err, "event"))
				return
			}

			id, err := p.PublishSmiResults(result)
			if err != nil {
				log.Error(ErrPublishSmiResults(err))
				return
			}
			event.Details = fmt.Sprintf("Result-Id: %s", id)
		}

		data, err := json.Marshal(event)
		if err != nil {
			log.Error(ErrMarshal(err, "event"))
			return
		}
		respChan <- data
	}
}
