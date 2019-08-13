package handlers

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/layer5io/meshery/meshes"
	"github.com/layer5io/meshery/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"strconv"
)

func init() {
	gob.Register([]*models.Adapter{})
}

func (h *Handler) GetAllAdaptersHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, err := h.config.SessionStore.Get(req, h.config.SessionName)
	if err != nil {
		logrus.Errorf("error getting session: %v", err)
		http.Error(w, "unable to get session", http.StatusUnauthorized)
		return
	}

	err = json.NewEncoder(w).Encode(h.config.AdapterTracker.GetAdapters(req.Context()))
	if err != nil {
		logrus.Errorf("error marshalling data: %v", err)
		http.Error(w, "unable to retrieve the requested data", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) MeshAdapterConfigHandler(w http.ResponseWriter, req *http.Request) {
	session, err := h.config.SessionStore.Get(req, h.config.SessionName)
	if err != nil {
		logrus.Errorf("error getting session: %v", err)
		http.Error(w, "unable to get session", http.StatusUnauthorized)
		return
	}

	var user *models.User
	user, _ = session.Values["user"].(*models.User)

	// h.config.SessionPersister.Lock(user.UserId)
	// defer h.config.SessionPersister.Unlock(user.UserId)

	sessObj, err := h.config.SessionPersister.Read(user.UserId)
	if err != nil {
		logrus.Warn("unable to read session from the session persister, starting with a new one")
	}

	if sessObj == nil {
		sessObj = &models.Session{}
	}

	meshAdapters := sessObj.MeshAdapters
	if meshAdapters == nil {
		meshAdapters = []*models.Adapter{}
	}

	switch req.Method {
	case http.MethodPost:
		meshLocationURL := req.FormValue("meshLocationURL")

		logrus.Debugf("meshLocationURL: %s", meshLocationURL)
		if strings.TrimSpace(meshLocationURL) == "" {
			err := errors.New("meshLocationURL cannot be empty to add an adapter")
			logrus.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if sessObj.K8SConfig == nil || !sessObj.K8SConfig.InClusterConfig && (sessObj.K8SConfig.Config == nil || len(sessObj.K8SConfig.Config) == 0) {
			err := errors.New("no valid kubernetes config found")
			logrus.Error(err)
			http.Error(w, "No valid kubernetes config found.", http.StatusBadRequest)
			return
		}

		meshAdapters, err = h.addAdapter(req.Context(), meshAdapters, sessObj, meshLocationURL)
		if err != nil {
			http.Error(w, "unable to retrieve the requested data", http.StatusInternalServerError)
			return // error is handled appropriately in the relevant method
		}
	case http.MethodDelete:
		meshAdapters, err = h.deleteAdapter(meshAdapters, w, req)
		if err != nil {
			return // error is handled appropriately in the relevant method
		}
	default:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	sessObj.MeshAdapters = meshAdapters
	err = h.config.SessionPersister.Write(user.UserId, sessObj)
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

func (h *Handler) addAdapter(ctx context.Context, meshAdapters []*models.Adapter, sessObj *models.Session, meshLocationURL string) ([]*models.Adapter, error) {
	alreadyConfigured := false
	for _, adapter := range meshAdapters {
		if adapter.Location == meshLocationURL {
			// err := errors.New("Adapter with the given meshLocationURL already exists.")
			// logrus.Error(err)
			// http.Error(w, err.Error(), http.StatusForbidden)
			// return nil, err
			alreadyConfigured = true
			break
		}
	}

	if alreadyConfigured {
		logrus.Debugf("adapter already configured. . . ")
		return meshAdapters, nil
	}

	mClient, err := meshes.CreateClient(ctx, sessObj.K8SConfig.Config, sessObj.K8SConfig.ContextName, meshLocationURL)
	if err != nil {
		err = errors.Wrapf(err, "error creating a mesh client")
		logrus.Error(err)
		// http.Error(w, "Unable to connect to the Mesh adapter using the given config, please try again", http.StatusInternalServerError)
		return nil, err
	}
	defer mClient.Close()
	respOps, err := mClient.MClient.SupportedOperations(ctx, &meshes.SupportedOperationsRequest{})
	if err != nil {
		logrus.Errorf("error getting operations for the mesh: %v", err)
		// http.Error(w, "unable to retrieve the requested data", http.StatusInternalServerError)
		return nil, err
	}

	meshNameOps, err := mClient.MClient.MeshName(ctx, &meshes.MeshNameRequest{})
	if err != nil {
		err = errors.Wrapf(err, "error getting mesh name")
		logrus.Error(err)
		// http.Error(w, "unable to retrieve the requested data", http.StatusInternalServerError)
		return nil, err
	}

	result := &models.Adapter{
		Location: meshLocationURL,
		Name:     meshNameOps.GetName(),
		Ops:      respOps.Ops,
	}

	h.config.AdapterTracker.AddAdapter(ctx, meshLocationURL)

	return append(meshAdapters, result), nil
}

func (h *Handler) deleteAdapter(meshAdapters []*models.Adapter, w http.ResponseWriter, req *http.Request) ([]*models.Adapter, error) {

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
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	session, err := h.config.SessionStore.Get(req, h.config.SessionName)
	if err != nil {
		logrus.Error("unable to get session data")
		http.Error(w, "unable to get user data", http.StatusUnauthorized)
		return
	}

	var user *models.User
	user, _ = session.Values["user"].(*models.User)

	// h.config.SessionPersister.Lock(user.UserId)
	// defer h.config.SessionPersister.Unlock(user.UserId)

	sessObj, err := h.config.SessionPersister.Read(user.UserId)
	if err != nil {
		logrus.Warn("unable to read session from the session persister, starting with a new one")
	}

	if sessObj == nil {
		sessObj = &models.Session{}
	}

	meshAdapters := sessObj.MeshAdapters
	if meshAdapters == nil {
		meshAdapters = []*models.Adapter{}
	}

	adapterID := req.PostFormValue("index")
	logrus.Debugf("adapterID to execute ops on: %s", adapterID)

	adaptersLen := len(meshAdapters)

	aId, err := strconv.Atoi(adapterID)
	if err != nil {
		err = errors.Wrapf(err, "unable to parse the given adapter id")
		logrus.Error(err)
		http.Error(w, "given adapter id is not valid", http.StatusBadRequest)
		return
	}
	if aId < 0 || aId >= adaptersLen {
		err := errors.New("given adapter id is outside the valid range")
		logrus.Error(err)
		http.Error(w, "given adapter id is not valid", http.StatusBadRequest)
		return
	}
	opName := req.PostFormValue("query")
	customBody := req.PostFormValue("customBody")
	namespace := req.PostFormValue("namespace")
	delete := req.PostFormValue("deleteOp")
	if namespace == "" {
		namespace = "default"
	}

	if sessObj.K8SConfig == nil || !sessObj.K8SConfig.InClusterConfig && (sessObj.K8SConfig.Config == nil || len(sessObj.K8SConfig.Config) == 0) {
		logrus.Error("no valid kubernetes config found")
		http.Error(w, `No valid kubernetes config found.`, http.StatusBadRequest)
		return
	}

	mClient, err := meshes.CreateClient(req.Context(), sessObj.K8SConfig.Config, sessObj.K8SConfig.ContextName, meshAdapters[aId].Location)
	if err != nil {
		logrus.Errorf("error creating a mesh client: %v", err)
		http.Error(w, "Unable to create a mesh client", http.StatusBadRequest)
		return
	}
	defer mClient.Close()

	_, err = mClient.MClient.ApplyOperation(req.Context(), &meshes.ApplyRuleRequest{
		OpName:     opName,
		Username:   user.UserId,
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
