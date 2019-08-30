package handlers

import (
	"fmt"
	"io"
	"net/http"

	"encoding/json"

	"github.com/layer5io/meshery/meshes"
	"github.com/layer5io/meshery/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// EventStreamHandler endpoint is used for streaming events to the frontend
func (h *Handler) EventStreamHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	session, err := h.config.SessionStore.Get(req, h.config.SessionName)
	if err != nil {
		logrus.Errorf("Error getting session: %v.", err)
		http.Error(w, "Unable to get session.", http.StatusUnauthorized)
		return
	}

	var user *models.User
	user, _ = session.Values["user"].(*models.User)

	// h.config.SessionPersister.Lock(user.UserID)
	// defer h.config.SessionPersister.Unlock(user.UserID)

	sessObj, err := h.config.SessionPersister.Read(user.UserID)
	if err != nil {
		logrus.Warn("Unable to read session from the session persister. Starting with a new session.")
	}

	if sessObj == nil {
		sessObj = &models.Session{}
	}

	meshAdapters := sessObj.MeshAdapters
	if meshAdapters == nil {
		meshAdapters = []*models.Adapter{}
	}

	if sessObj.K8SConfig == nil || !sessObj.K8SConfig.InClusterConfig && (sessObj.K8SConfig.Config == nil || len(sessObj.K8SConfig.Config) == 0) {
		logrus.Error("No valid Kubernetes config found.")
		http.Error(w, `No valid Kubernetes config found.`, http.StatusBadRequest)
		return
	}

	adaptersLen := len(meshAdapters)
	if adaptersLen == 0 {
		logrus.Error("No valid mesh adapter(s) found.")
		http.Error(w, `No valid mesh adapter(s) found.`, http.StatusBadRequest)
		return
	}

	flusher, ok := w.(http.Flusher)

	if !ok {
		logrus.Errorf("Event streaming not supported: %v.", err)
		http.Error(w, "Event streaming is not supported at the moment.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	notify := w.(http.CloseNotifier).CloseNotify()

	// go func() {
	// 	<-notify
	// 	// an attempt to re-establish connection
	// 	// mClient, _ = meshes.CreateClient(req.Context(), k8sConfigBytes, contextName, meshLocationURL)
	// }()

	respChan := make(chan []byte, 100)
	errChan := make(chan error)
	defer close(respChan)

	for _, ma := range meshAdapters {
		go func(ma *models.Adapter) {
			mClient, err := meshes.CreateClient(req.Context(), sessObj.K8SConfig.Config, sessObj.K8SConfig.ContextName, ma.Location)
			if err != nil {
				err = errors.Wrapf(err, "Error creating a mesh client.")
				logrus.Error(err)
				errChan <- err
				// http.Error(w, "Unable to create a mesh client", http.StatusBadRequest)
				return
			}
			defer mClient.Close()

			logrus.Debugf("Received a stream client...")

			streamClient, err := mClient.MClient.StreamEvents(req.Context(), &meshes.EventsRequest{})
			if err != nil {
				err = errors.Wrapf(err, "There was an error connecting to the backend to get events.")
				logrus.Error(err)
				errChan <- err
				// http.Error(w, "There was an error connecting to the backend to get events", http.StatusInternalServerError)
				return
			}

			for {
				logrus.Debugf("Waiting to receive events.")
				event, err := streamClient.Recv()
				if err != nil {
					if err == io.EOF {
						err = errors.Wrapf(err, "Event streaming ended.")
						logrus.Error(err)
						errChan <- nil
						// break
						return
					}
					err = errors.Wrapf(err, "Event streaming ended with an unknown error.")
					logrus.Error(err)
					// http.Error(w, "streaming events was interrupted", http.StatusInternalServerError)
					// return
					errChan <- err
					// break
					return
				}
				// logrus.Debugf("received an event: %+#v", event)
				logrus.Debugf("Received an event.")
				data, err := json.Marshal(event)
				if err != nil {
					err = errors.Wrapf(err, "Error marshalling event to json.")
					logrus.Error(err)
					errChan <- err
					// logrus.Errorf(
					// http.Error(w, "error while sending event to client", http.StatusInternalServerError)
					return
				}
				respChan <- data
			}
		}(ma)
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorf("Recovered from panic: %v.", r)
			}
		}()
		for {
			select {
			case data := <-respChan:
				fmt.Fprintf(w, "data: %s\n\n", data)
				if flusher != nil {
					flusher.Flush()
					logrus.Debugf("Flushed the messages on the wire...")
				}
			}
		}
	}()

	select {
	case <-notify:
	case err := <-errChan:
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// close(respChan)
	// ... close all events and their channels
}
