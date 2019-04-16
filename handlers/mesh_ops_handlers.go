package handlers

import (
	"encoding/gob"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/layer5io/meshery/meshes"
	"github.com/layer5io/meshery/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"strconv"
)

func init() {
	gob.Register([]*models.Adapter{})
}

func (h *Handler) MeshAdapterConfigHandler(w http.ResponseWriter, req *http.Request) {
	session, err := h.config.SessionStore.Get(req, h.config.SessionName)
	if err != nil {
		logrus.Errorf("error getting session: %v", err)
		http.Error(w, "unable to get session", http.StatusUnauthorized)
		return
	}

	var meshAdapters []*models.Adapter

	meshAdaptersI, ok := session.Values["meshAdapters"]
	if ok && meshAdaptersI != nil {
		meshAdapters, _ = meshAdaptersI.([]*models.Adapter)
	}

	if meshAdapters == nil {
		meshAdapters = []*models.Adapter{}
	}

	switch req.Method {
	case http.MethodPost:
		meshAdapters, err = h.addAdapter(meshAdapters, session, w, req)
		if err != nil {
			return // error is handled appropriately in the relevant method
		}
	case http.MethodDelete:
		meshAdapters, err = h.deleteAdapter(meshAdapters, session, w, req)
		if err != nil {
			return // error is handled appropriately in the relevant method
		}
	default:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	session.Values["meshAdapters"] = meshAdapters

	err = session.Save(req, w)
	if err != nil {
		logrus.Errorf("unable to save session: %v", err)
		http.Error(w, "unable to save session", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(meshAdapters)
	if err != nil {
		logrus.Errorf("error marshalling data: %v", err)
		http.Error(w, "unable to retrieve the requested data", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) addAdapter(meshAdapters []*models.Adapter, session *sessions.Session, w http.ResponseWriter, req *http.Request) ([]*models.Adapter, error) {

	meshLocationURL := req.FormValue("meshLocationURL")
	logrus.Debugf("meshLocationURL: %s", meshLocationURL)
	if strings.TrimSpace(meshLocationURL) == "" {
		err := errors.New("meshLocationURL cannot be empty to add an adapter")
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}

	for _, adapter := range meshAdapters {
		if adapter.Location == meshLocationURL {
			err := errors.New("Adapter with the given meshLocationURL already exists.")
			logrus.Error(err)
			http.Error(w, err.Error(), http.StatusForbidden)
			return nil, err
		}
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
		err := errors.New("no valid kubernetes config found")
		logrus.Error(err)
		http.Error(w, "No valid kubernetes config found.", http.StatusBadRequest)
		return nil, err
	}
	k8sConfigBytes, _ := k8sConfigBytesI.([]byte)

	mClient, err := meshes.CreateClient(req.Context(), k8sConfigBytes, contextName, meshLocationURL)
	if err != nil {
		err = errors.Wrapf(err, "error creating a mesh client")
		logrus.Error(err)
		http.Error(w, "Unable to connect to the Mesh adapter using the given config, please try again", http.StatusInternalServerError)
		return nil, err
	}
	defer mClient.Close()
	respOps, err := mClient.MClient.SupportedOperations(req.Context(), &meshes.SupportedOperationsRequest{})
	if err != nil {
		logrus.Errorf("error getting operations for the mesh: %v", err)
		http.Error(w, "unable to retrieve the requested data", http.StatusInternalServerError)
		return nil, err
	}

	meshNameOps, err := mClient.MClient.MeshName(req.Context(), &meshes.MeshNameRequest{})
	if err != nil {
		err = errors.Wrapf(err, "error getting mesh name")
		logrus.Error(err)
		http.Error(w, "unable to retrieve the requested data", http.StatusInternalServerError)
		return nil, err
	}

	result := &models.Adapter{
		Location: meshLocationURL,
		Name:     meshNameOps.GetName(),
		Ops:      respOps.Ops,
	}
	return append(meshAdapters, result), nil
}

func (h *Handler) deleteAdapter(meshAdapters []*models.Adapter, session *sessions.Session, w http.ResponseWriter, req *http.Request) ([]*models.Adapter, error) {

	adapterID := req.URL.Query().Get("adapterID")
	logrus.Debugf("adapterID of adapter to be removed: %s", adapterID)

	adaptersLen := len(meshAdapters)

	aId, err := strconv.Atoi(adapterID)
	if err != nil {
		err = errors.Wrapf(err, "unable to parse the given id")
		logrus.Error(err)
		http.Error(w, "given id is not valid", http.StatusBadRequest)
		return nil, err
	}
	if aId < 0 || aId >= adaptersLen {
		err := errors.New("given id is outside the valid range")
		logrus.Error(err)
		http.Error(w, "given id is not valid", http.StatusBadRequest)
		return nil, err
	}
	newMeshAdapters := []*models.Adapter{}
	if aId == 0 {
		newMeshAdapters = append(newMeshAdapters, meshAdapters[1:]...)
	} else if aId == adaptersLen-1 {
		newMeshAdapters = append(newMeshAdapters, meshAdapters[:adaptersLen-1]...)
	} else {
		newMeshAdapters = append(newMeshAdapters, meshAdapters[0:aId]...)
		newMeshAdapters = append(newMeshAdapters, meshAdapters[aId+1:]...)
	}
	if logrus.GetLevel() == logrus.DebugLevel {
		b, _ := json.Marshal(meshAdapters)
		logrus.Debugf("Old adapters: %s", b)
		b, _ = json.Marshal(newMeshAdapters)
		logrus.Debugf("New adapters: %s", b)
	}
	return newMeshAdapters, nil
}

func (h *Handler) MeshOpsHandler(w http.ResponseWriter, req *http.Request) {
	opName := req.PostFormValue("query")
	customBody := req.PostFormValue("customBody")
	namespace := req.PostFormValue("namespace")
	delete := req.PostFormValue("deleteOp")
	if namespace == "" {
		namespace = "default"
	}

	session, err := h.config.SessionStore.Get(req, h.config.SessionName)
	if err != nil {
		logrus.Error("unable to get session data")
		http.Error(w, "unable to get user data", http.StatusUnauthorized)
		return
	}
	userNameI, ok := session.Values["user"]
	if !ok && userNameI == nil {
		logrus.Error("unable to get session data")
		http.Error(w, "unable to get user data", http.StatusUnauthorized)
		return
	}
	userName, _ := userNameI.(string)

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

	_, err = mClient.MClient.ApplyOperation(req.Context(), &meshes.ApplyRuleRequest{
		OpName:     opName,
		Username:   userName,
		Namespace:  namespace,
		CustomBody: customBody,
		DeleteOp:   (delete != ""),
	})
	if err != nil {
		logrus.Error(err)
		http.Error(w, "There was an error applying the change", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("{}"))
}
