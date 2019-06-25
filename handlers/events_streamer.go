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

func (h *Handler) EventStreamHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	session, err := h.config.SessionStore.Get(req, h.config.SessionName)
	if err != nil {
		logrus.Errorf("error getting session: %v", err)
		http.Error(w, "unable to get session", http.StatusUnauthorized)
		return
	}

	contextName := ""
	contextNameI, ok := session.Values["k8sContext"]
	if ok && contextNameI != nil {
		contextName, _ = contextNameI.(string)
	}

	inClusterConfig := ""
	inClusterConfigI, ok := session.Values["k8sInCluster"]
	if ok && contextNameI != nil {
		inClusterConfig, _ = inClusterConfigI.(string)
	}

	// logrus.Debugf("session values: %v", session.Values)
	k8sConfigBytesI, ok := session.Values["k8sConfig"]
	if inClusterConfig == "" && (!ok || k8sConfigBytesI == nil) {
		logrus.Error("no valid kubernetes config found")
		http.Error(w, `No valid kubernetes config found.`, http.StatusBadRequest)
		return
	}
	k8sConfigBytes, _ := k8sConfigBytesI.([]byte)

	var meshAdapters []*models.Adapter

	meshAdaptersI, ok := session.Values["meshAdapters"]
	if ok && meshAdaptersI != nil {
		meshAdapters, _ = meshAdaptersI.([]*models.Adapter)
	}

	if meshAdapters == nil {
		meshAdapters = []*models.Adapter{}
	}

	adaptersLen := len(meshAdapters)
	if adaptersLen == 0 {
		logrus.Error("no valid mesh adapter(s) found")
		http.Error(w, `No valid mesh adapter(s) found.`, http.StatusBadRequest)
		return
	}

	flusher, ok := w.(http.Flusher)

	if !ok {
		logrus.Errorf("streaming not supported: %v", err)
		http.Error(w, "streaming events is not supported at the moment", http.StatusInternalServerError)
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

	for _, ma := range meshAdapters {
		go func(ma *models.Adapter) {
			mClient, err := meshes.CreateClient(req.Context(), k8sConfigBytes, contextName, ma.Location)
			if err != nil {
				err = errors.Wrapf(err, "error creating a mesh client")
				logrus.Error(err)
				errChan <- err
				// http.Error(w, "Unable to create a mesh client", http.StatusBadRequest)
				return
			}
			defer mClient.Close()

			logrus.Debugf("received a stream client. . .")

			streamClient, err := mClient.MClient.StreamEvents(req.Context(), &meshes.EventsRequest{})
			if err != nil {
				err = errors.Wrapf(err, "There was an error connecting to the backend to get events")
				logrus.Error(err)
				errChan <- err
				// http.Error(w, "There was an error connecting to the backend to get events", http.StatusInternalServerError)
				return
			}

			for {
				logrus.Debugf("waiting to receive events")
				event, err := streamClient.Recv()
				if err != nil {
					if err == io.EOF {
						err = errors.Wrapf(err, "streaming ended")
						logrus.Error(err)
						errChan <- nil
						// break
						return
					} else {
						err = errors.Wrapf(err, "streaming ended with an unknown error")
						logrus.Error(err)
						// http.Error(w, "streaming events was interrupted", http.StatusInternalServerError)
						// return
						errChan <- err
						// break
						return
					}
				}
				// logrus.Debugf("received an event: %+#v", event)
				logrus.Debugf("received an event")
				data, err := json.Marshal(event)
				if err != nil {
					err = errors.Wrapf(err, "error marshalling event to json")
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
				logrus.Errorf("Recovered from panic: %v", r)
			}
		}()
		for {
			select {
			case data := <-respChan:
				fmt.Fprintf(w, "data: %s\n\n", data)
				if flusher != nil {
					flusher.Flush()
					logrus.Debugf("flushed the messages on the wire. . .")
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
	close(respChan)
	// ... close all events and their channels
}
