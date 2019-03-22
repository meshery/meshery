package handlers

import (
	"fmt"
	"io"
	"net/http"

	"encoding/json"

	"github.com/layer5io/meshery/meshes"
	"github.com/sirupsen/logrus"
)

func (h *Handler) EventStreamHandler(w http.ResponseWriter, req *http.Request) {
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

	meshLocationURL := ""
	meshLocationURLI, ok := session.Values["meshLocationURL"]
	if ok && meshLocationURLI != nil {
		meshLocationURL, _ = meshLocationURLI.(string)
	} else {
		logrus.Error("no valid url for mesh adapter found")
		http.Error(w, `No valid url for mesh adapter found.`, http.StatusBadRequest)
		return
	}

	mClient, err := meshes.CreateClient(req.Context(), k8sConfigBytes, contextName, meshLocationURL)
	if err != nil {
		logrus.Errorf("error creating a mesh client: %v", err)
		http.Error(w, "Unable to create a mesh client", http.StatusBadRequest)
		return
	}
	defer mClient.Close()

	streamClient, err := mClient.MClient.StreamEvents(req.Context(), &meshes.EventsRequest{})
	if err != nil {
		logrus.Error(err)
		http.Error(w, "There was an error connecting to the backend to get events", http.StatusInternalServerError)
		return
	}
	logrus.Debugf("received a stream client. . .")
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

	// notify := w.(http.CloseNotifier).CloseNotify()

	// go func() {
	// 	<-notify
	// 	// an attempt to re-establish connection
	// 	// mClient, _ = meshes.CreateClient(req.Context(), k8sConfigBytes, contextName, meshLocationURL)
	// }()
	for {
		logrus.Debugf("waiting to receive events")
		event, err := streamClient.Recv()
		if err != nil {
			if err == io.EOF {
				logrus.Errorf("streaming ended: %v", err)
				break
			} else {
				logrus.Errorf("streaming ended with an unknown error: %v", err)
				// http.Error(w, "streaming events was interrupted", http.StatusInternalServerError)
				return
			}
		}
		// logrus.Debugf("received an event: %+#v", event)
		logrus.Debugf("received an event")
		data, err := json.Marshal(event)
		if err != nil {
			logrus.Errorf("error marshalling event to json: %v", err)
			http.Error(w, "error while sending event to client", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
		logrus.Debugf("flushed the messages on the wire. . .")
	}
}
