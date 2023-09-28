package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/models/events"
	"github.com/layer5io/meshkit/models/meshmodel/registry"
	helm "github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/layer5io/meshkit/utils/patterns"
	"gopkg.in/yaml.v2"
)

var helmIndexes = []string{"index.yml", "index.yaml"}

type Helm struct {
	log      logger.Handler
	userID   *uuid.UUID
	systemID *uuid.UUID
	eb       *models.Broadcast
	provider models.Provider
}

func (h *Helm) Register(res http.ResponseWriter, req *http.Request) {
	var description string

	requestBody := make(map[string]string, 3)

	defer func() {
		_ = req.Body.Close()
	}()

	res.Header().Set("content/type", "application/json")
	eb := events.NewEvent().ActedUpon(*h.userID).FromSystem(*h.systemID).WithCategory("connection").WithAction("create").FromUser(*h.userID)

	err := json.NewDecoder(req.Body).Decode(&requestBody)
	if err != nil {
		unmarshalErr := models.ErrUnmarshal(err, "connection request body")
		h.log.Error(unmarshalErr)
		
		description = "Cannot read connection payload from the request body."
		eb.WithDescription(description).WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": unmarshalErr,
		})
		event := eb.Build()

		err = h.provider.PersistEvent(event)
		if err != nil {
			h.log.Warn(fmt.Errorf("unable to create event: %s", err.Error()))
		}

		res.WriteHeader(http.StatusInternalServerError)
		data, _ := json.Marshal(event)
		_, err = res.Write(data)
		if err != nil {
			h.log.Error(ErrWriteResponse)
		}
		go h.eb.Publish(*h.userID, event)
		return
	}

	repoURL, _err := extractURL(requestBody["repoURL"], eb)
	if _err != nil {
		err = h.provider.PersistEvent(_err)
		if err != nil {
			h.log.Warn(fmt.Errorf("unable to create event: %s", err.Error()))
		}
		go h.eb.Publish(*h.userID, _err)
		res.WriteHeader(http.StatusBadRequest)
		data, _ := json.Marshal(_err)
		_, err = res.Write(data)
		if err != nil {
			h.log.Error(ErrWriteResponse)
		}
	}

	connectionName, ok := requestBody["name"]
	if !ok || connectionName == "" {
		invalidReqBodyErr := ErrInvalidRequestObject("Connection name is empty.")
		h.log.Error(invalidReqBodyErr)
		description = fmt.Sprintf("Invalid connection name provided %v.", connectionName)
		eb.WithDescription(description).WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": invalidReqBodyErr,
		})
		event := eb.Build()

		err = h.provider.PersistEvent(event)
		if err != nil {
			h.log.Warn(fmt.Errorf("unable to create event: %s", err.Error()))
		}
		go h.eb.Publish(*h.userID, event)
		res.WriteHeader(http.StatusBadRequest)
		data, _ := json.Marshal(event)
		_, err = res.Write(data)
		if err != nil {
			h.log.Error(ErrWriteResponse)
		}
		return
	}
	description = fmt.Sprintf("Helm Connection %s with repository %s created.", connectionName, repoURL)

	connectionID := uuid.NewV5(uuid.NamespaceOID, sanitizeURL(repoURL))
	h.log.Info("Generated id for repo: ", connectionID, repoURL)

	connDescription := requestBody["description"]
	connection := &models.Connection{
		ID:      connectionID,
		Name:    connectionName,
		Kind:    string(models.Helm),
		Type:    string(models.PLATFORM),
		SubType: models.ORCHESTRATOR,
		Status:  models.REGISTERED,
		UserID:  h.userID,
		Metadata: map[string]interface{}{
			"description": connDescription,
			"repoURL": repoURL,
		},
	}

	registeredConnection, err := h.provider.RegisterConnection(req, connection.Kind, connection)

	if err != nil {
		errRegisterConnection := ErrFailToSave(err, "connection")
		h.log.Error(errRegisterConnection)
		description = fmt.Sprintf("Unable to register connection %s.", connectionName)
		eb.WithDescription(description).WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": errRegisterConnection,
		})
		event := eb.Build()

		err = h.provider.PersistEvent(event)
		if err != nil {
			h.log.Warn(fmt.Errorf("unable to create event: %s", err.Error()))
		}
		go h.eb.Publish(*h.userID, event)
		res.WriteHeader(http.StatusInternalServerError)
		data, _ := json.Marshal(event)
		_, err = res.Write(data)
		if err != nil {
			h.log.Error(ErrWriteResponse)
		}
		return
	}

	event := eb.WithSeverity(events.Success).WithDescription(description).ActedUpon(connection.ID).Build()
	err = h.provider.PersistEvent(event)
	go h.eb.Publish(*h.userID, event)
	if err != nil {
		h.log.Warn(fmt.Errorf("unable to create event: %s", err.Error()))
	}
	_ = json.NewEncoder(res).Encode(registeredConnection)
}

func (h *Helm) Status(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content/type", "application/json")
	eb := events.NewEvent().ActedUpon(*h.userID).FromSystem(*h.systemID).WithCategory("connection").WithAction("verify").FromUser(*h.userID)

	var description string

	// The repo url serves as a id for helm connection
	id := req.URL.Query().Get("id")
	repoURL, err := url.QueryUnescape(id)
	if err != nil {
		invalidQueryParamErr := ErrQueryGet("id")
		description = "Query parameter \"id\" is invalid."
		eb.WithDescription(description).WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": invalidQueryParamErr,
		})
		event := eb.Build()

		err = h.provider.PersistEvent(event)
		if err != nil {
			h.log.Warn(fmt.Errorf("unable to create event: %s", err.Error()))
		}
		
		res.WriteHeader(http.StatusBadRequest)
		data, _ := json.Marshal(invalidQueryParamErr)
		go h.eb.Publish(*h.userID, invalidQueryParamErr)
		_, err = res.Write(data)
		if err != nil {
			h.log.Error(ErrWriteResponse)
		}
	}
	
	repoURL, _err := extractURL(repoURL, eb)
	if _err != nil {		
		err = h.provider.PersistEvent(_err)
		if err != nil {
			h.log.Warn(fmt.Errorf("unable to create event: %s", err.Error()))
		}

		res.WriteHeader(http.StatusBadRequest)
		data, _ := json.Marshal(_err)
		go h.eb.Publish(*h.userID, _err)
		_, err = res.Write(data)
		if err != nil {
			h.log.Error(ErrWriteResponse)
		}
		return
	}

	repoURL = sanitizeURL(repoURL)

	connectionID := uuid.NewV5(uuid.NamespaceOID, sanitizeURL(repoURL))
	h.log.Info("Generated id for repo: ", connectionID, repoURL)

	existingConnection, err := h.provider.GetConnectionById(req, &connectionID, models.Helm)
	if err == nil {
		_ = json.NewEncoder(res).Encode(existingConnection)
		return
	}
}

func (h *Helm) Verify(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("content/type", "application/json")
	eb := events.NewEvent().ActedUpon(*h.userID).FromSystem(*h.systemID).WithCategory("connection").WithAction("verify").FromUser(*h.userID)

	var description string

	requestBody := make(map[string]string, 1)
	defer func() {
		_ = req.Body.Close()
	}()

	err := json.NewDecoder(req.Body).Decode(&requestBody)
	if err != nil {
		unmarshalErr := models.ErrUnmarshal(err, "connection request body")
		h.log.Error(unmarshalErr)
		description = "Cannot read repository url from the request body."
		eb.WithDescription(description).WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": unmarshalErr,
		})
		event := eb.Build()

		err = h.provider.PersistEvent(event)
		if err != nil {
			h.log.Warn(fmt.Errorf("unable to create event: %s", err.Error()))
		}

		res.WriteHeader(http.StatusInternalServerError)
		data, _ := json.Marshal(event)
		go h.eb.Publish(*h.userID, event)
		_, err = res.Write(data)
		if err != nil {
			h.log.Error(ErrWriteResponse)
		}
		return
	}

	repoURL, _err := extractURL(requestBody["repoURL"], eb)
	if _err != nil {
		err = h.provider.PersistEvent(_err)
		if err != nil {
			h.log.Warn(fmt.Errorf("unable to create event: %s", err.Error()))
		}

		res.WriteHeader(http.StatusBadRequest)
		data, _ := json.Marshal(_err)
		go h.eb.Publish(*h.userID, _err)
		_, err = res.Write(data)
		if err != nil {
			h.log.Error(ErrWriteResponse)
		}
	}
	
	repoURL = sanitizeURL(repoURL)

	helmIndex := getHelmIndex(repoURL, eb)
	if helmIndex == nil {
		event := eb.Build()		
		err = h.provider.PersistEvent(event)
		if err != nil {
			h.log.Warn(fmt.Errorf("unable to create event: %s", err.Error()))
		}

		res.WriteHeader(http.StatusInternalServerError)
		data, _ := json.Marshal(event)
		_, err = res.Write(data)
		if err != nil {
			h.log.Error(ErrWriteResponse)
		}
		go h.eb.Publish(*h.userID, event)
		return
	}

	if len(helmIndex.Entries) == 0 {
		event := eb.WithDescription(fmt.Sprintf("Repsository %s is empty.", repoURL)).WithSeverity(events.Informational).Build()
		err = h.provider.PersistEvent(event)
		if err != nil {
			h.log.Warn(fmt.Errorf("unable to create event: %s", err.Error()))
		}

		data, _ := json.Marshal(event)
		go h.eb.Publish(*h.userID, event)
		_, err = res.Write(data)
		if err != nil {
			h.log.Error(ErrWriteResponse)
		}
		return
	}

	h.log.Info("HELM ENTRIES: ", helmIndex)
}

func (h *Helm) Details(res http.ResponseWriter, req *http.Request) {
	
	res.Header().Set("content/type", "application/json")
	eb := events.NewEvent().ActedUpon(*h.userID).FromSystem(*h.systemID).FromUser(*h.userID).WithCategory("connection").WithAction("fetch") // change

	var description string

	requestBody := make(map[string]string, 1)
	defer func() {
		_ = req.Body.Close()
	}()

	err := json.NewDecoder(req.Body).Decode(&requestBody)
	if err != nil {
		unmarshalErr := models.ErrUnmarshal(err, "connection request body")
		h.log.Error(unmarshalErr)
		description = "Cannot read repository url from the request body, unable to connect."
		eb.WithDescription(description).WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": err,
		})
		event := eb.Build()

		err = h.provider.PersistEvent(event)
		if err != nil {
			h.log.Warn(fmt.Errorf("unable to create event: %s", err.Error()))
		}
		data, _ := json.Marshal(event)
		res.Write(data)
		go h.eb.Publish(*h.userID, event)
		return
	}

	repoURL, _err := extractURL(requestBody["repoURL"], eb)
	if _err != nil {		
		err = h.provider.PersistEvent(_err)
		if err != nil {
			h.log.Warn(fmt.Errorf("unable to create event: %s", err.Error()))
		}

		res.WriteHeader(http.StatusBadRequest)
		data, _ := json.Marshal(_err)
		go h.eb.Publish(*h.userID, _err)
		_, err = res.Write(data)
		if err != nil {
			h.log.Error(ErrWriteResponse)
		}
		return
	}
	
	repoURL = sanitizeURL(repoURL)
	h.getHelmEntries(repoURL, eb, res)
}

func (h *Helm) AddMetadata(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content/type", "application/json")
	eb := events.NewEvent().ActedUpon(*h.userID).FromSystem(*h.systemID).WithCategory("connection").WithAction("update").FromUser(*h.userID)

	var description string

	requestBody := make(map[string]interface{}, 2)
	defer func() {
		_ = req.Body.Close()
	}()

	err := json.NewDecoder(req.Body).Decode(&requestBody)
	if err != nil {
		unmarshalErr := models.ErrUnmarshal(err, "connection request body")
		h.log.Error(unmarshalErr)

		description := "Cannot read request body."
		eb.WithDescription(description).WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": err,
		})
		event := eb.Build()

		err = h.provider.PersistEvent(event)
		if err != nil {
			h.log.Warn(fmt.Errorf("unable to create event: %s", err.Error()))
		}

		res.WriteHeader(http.StatusInternalServerError)
		data, _ := json.Marshal(event)
		go h.eb.Publish(*h.userID, event)
		_, err = res.Write(data)
		if err != nil {
			h.log.Error(ErrWriteResponse)
		}
		return
	}

	repoURL := requestBody["repoURL"]
	updatedConnection, err := h.provider.UpdateConnectionMetadata(req, string(models.Helm), requestBody)
	if updatedConnection == nil && err == nil {
		description = fmt.Sprintf("%s connection with repository %s doesn't exist or has been deleted.", models.Helm, repoURL)
		event := eb.WithDescription(description).WithSeverity(events.Informational).WithMetadata(map[string]interface{}{
			"summary": fmt.Sprintf("Helm connection with %s doesn't exist or has been deleted, please create one from the connections page.", repoURL),
		}).Build()
		
		err = h.provider.PersistEvent(event)
		go h.eb.Publish(*h.userID, event)
		if err != nil {
			h.log.Warn(fmt.Errorf("unable to create event: %s", err.Error()))
		}
	}

	if err != nil {
		h.log.Error(err)
		description := "Cannot update connection."
		eb.WithDescription(description).WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": err,
		})
		event := eb.Build()

		err = h.provider.PersistEvent(event)
		if err != nil {
			h.log.Warn(fmt.Errorf("unable to create event: %s", err.Error()))
		}

		go h.eb.Publish(*h.userID, event)
		res.WriteHeader(http.StatusInternalServerError)
		data, _ := json.Marshal(event)
		_, err = res.Write(data)
		if err != nil {
			h.log.Error(ErrWriteResponse)
		}
		return
	}
	
	_ = json.NewEncoder(res).Encode(updatedConnection)
	description = fmt.Sprintf("Connection %s updated.", updatedConnection.Name)
	event := eb.WithDescription(description).WithSeverity(events.Success).Build()
	err = h.provider.PersistEvent(event)
	if err != nil {
		h.log.Warn(fmt.Errorf("unable to create event: %s", err.Error()))
	}
	go h.eb.Publish(*h.userID, event)
}

func (h *Helm) Configure(res http.ResponseWriter, req *http.Request) {
	
	res.Header().Set("content/type", "application/json")
	eb := events.NewEvent().ActedUpon(*h.userID).FromSystem(*h.systemID).WithCategory("connection").WithAction("update").FromUser(*h.userID)

	requestBody := make(map[string]interface{}, 2)

	defer func(){
		_ = req.Body.Close()
	}()

	err := json.NewDecoder(req.Body).Decode(&requestBody)
	if err != nil {
		unmarshalErr := models.ErrUnmarshal(err, "connection request body")
		h.log.Error(unmarshalErr)

		description := "Cannot read request body."
		eb.WithDescription(description).WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": err,
		})
		event := eb.Build()

		err = h.provider.PersistEvent(event)
		if err != nil {
			h.log.Warn(fmt.Errorf("unable to create event: %s", err.Error()))
		}

		res.WriteHeader(http.StatusInternalServerError)
		data, _ := json.Marshal(event)
		_, err = res.Write(data)
		if err != nil {
			h.log.Error(ErrWriteResponse)
		}
		return
	}
	
	repoURL, _err := extractURL(requestBody["repoURL"], eb)
	if _err != nil {		
		err = h.provider.PersistEvent(_err)
		if err != nil {
			h.log.Warn(fmt.Errorf("unable to create event: %s", err.Error()))
		}

		res.WriteHeader(http.StatusBadRequest)
		data, _ := json.Marshal(_err)
		go h.eb.Publish(*h.userID, _err)
		_, err = res.Write(data)
		if err != nil {
			h.log.Error(ErrWriteResponse)
		}
		return
	}

	repoURL = sanitizeURL(repoURL)

	connectionID := uuid.NewV5(uuid.NamespaceOID, sanitizeURL(repoURL))
	h.log.Info("Generated id for repo: ", connectionID, repoURL)

	existingConnection, err := h.provider.GetConnectionById(req, &connectionID, models.Helm)
	if err != nil {
		description := fmt.Sprintf("Cannot retrieve %s connection for repository.", repoURL)
		eb.WithDescription(description).WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": err,
		})
		event := eb.Build()

		err = h.provider.PersistEvent(event)
		if err != nil {
			h.log.Warn(fmt.Errorf("unable to create event: %s", err.Error()))
		}

		res.WriteHeader(http.StatusInternalServerError)
		data, _ := json.Marshal(event)
		_, err = res.Write(data)
		if err != nil {
			h.log.Error(ErrWriteResponse)
		}
		return
	}
	charts, ok := existingConnection.Metadata["charts"]
	
	marshalledCharts, err := json.Marshal(charts)
	if !ok || err != nil {
		marshallErr := models.ErrMarshal(err, "helm charts")
		h.log.Error(marshallErr)
		description := fmt.Sprintf("Unable to import helm chart %s as meshery application", existingConnection.Name)
		event := eb.WithDescription(description).WithSeverity(events.Error).
		WithMetadata(map[string]interface{}{
			"error": marshallErr,
		}).
		Build()
		_ = h.provider.PersistEvent(event)
		
		go h.eb.Publish(*h.userID, event)
		return
	}

	var chartsToInstall []*helm.HelmEntryMetadata 
	err = json.Unmarshal(marshalledCharts, &chartsToInstall)
	if err != nil {
		unmarshallErr := models.ErrUnmarshal(err, "helm charts")
		h.log.Error(unmarshallErr)
		description := fmt.Sprintf("Unable to import helm chart %s as meshery application", existingConnection.Name)
		event := eb.WithDescription(description).WithSeverity(events.Error).
		WithMetadata(map[string]interface{}{
			"error": unmarshallErr,
		}).
		Build()
		_ = h.provider.PersistEvent(event)
		
		go h.eb.Publish(*h.userID, event)
		return
	}
	
	results := h.importChartsAsApps(req.Context(), chartsToInstall, repoURL)
	_ = json.NewEncoder(res).Encode(results)
}

func(h *Helm) importChartsAsApps(ctx context.Context, chartsToInstall []*helm.HelmEntryMetadata, repoURL string) []*helm.HelmEntryMetadata {
	mx := sync.Mutex{}
	
	regManager, _ := ctx.Value(models.RegistryManagerKey).(*registry.RegistryManager)
	
	token := ctx.Value(models.TokenCtxKey).(string)
	msgs := []*helm.HelmEntryMetadata{}
	wg := sync.WaitGroup{}
	repoURL = sanitizeURL(repoURL)
	for _, chart := range chartsToInstall {
		wg.Add(1)
		go func(chart helm.HelmEntryMetadata) {
			defer wg.Done()
			eb := events.NewEvent().ActedUpon(*h.userID).FromSystem(*h.systemID).WithCategory("application").WithAction("import").FromUser(*h.userID)
			cfg := helm.ApplyHelmChartConfig{
				ChartLocation: helm.HelmChartLocation{
					Repository: repoURL,
					Chart: chart.Name,
					Version: chart.Version,
					AppVersion: chart.AppVersion,
				},
			}

			helmManifest, err := helm.ConvertHelmChartToK8sManifest(cfg)
			if err != nil {
				h.log.Error(err)
				description := fmt.Sprintf("Unable to import helm chart %s as meshery application", chart.Name)
				event := eb.WithDescription(description).WithSeverity(events.Error).
				WithMetadata(map[string]interface{}{
					"error": err,
				}).
				Build()
				_ = h.provider.PersistEvent(event)
				mx.Lock()
				msgs = append(msgs, &chart)
				mx.Unlock()
				go h.eb.Publish(*h.userID, event)
				return
			}
			pattern, err := patterns.NewPatternFileFromK8sManifest(string(helmManifest), false, regManager)
			if err != nil {
				h.log.Error(err)
				description := fmt.Sprintf("Unable to import helm chart %s as meshery application", chart.Name)
				event := eb.WithDescription(description).WithSeverity(events.Error).
				WithMetadata(map[string]interface{}{
					"error": err,
				}).
				Build()
				_ = h.provider.PersistEvent(event)
				mx.Lock()
				msgs = append(msgs, &chart)
				mx.Unlock()
				go h.eb.Publish(*h.userID, event)
				return
			}
			application, err := convertAppToPattern(&pattern, repoURL, chart.Name, chart.Version)
			if err != nil {
				h.log.Error(err)
				description := fmt.Sprintf("Unable to import helm chart %s as meshery application", chart.Name)
				event := eb.WithDescription(description).WithSeverity(events.Error).
				WithMetadata(map[string]interface{}{
					"error": err,
				}).
				Build()
				_ = h.provider.PersistEvent(event)
				mx.Lock()
				msgs = append(msgs, &chart)
				mx.Unlock()
				go h.eb.Publish(*h.userID, event)
				return
			}
			resp, err := h.provider.SaveMesheryApplication(token, application)
			if err != nil {
				h.log.Error(err)
				errAppSave := ErrSaveApplication(err)
				h.log.Error(errAppSave)
				
				event := eb.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
					"error": errAppSave,
				}).WithDescription(fmt.Sprintf("Error saving application %s", application.Name)).Build()

				_ = h.provider.PersistEvent(event)
				go h.eb.Publish(*h.userID, event)
				mx.Lock()
				msgs = append(msgs, &chart)
				mx.Unlock()
				return
			}

			savedMesheryApplication := []models.MesheryApplication{}
			err = json.Unmarshal(resp, &savedMesheryApplication)
			// Ignoring error since if we are unable to save actual source content it is ok, since the actual source location is captured already.
			if err == nil {
				err = h.provider.SaveApplicationSourceContent(token, savedMesheryApplication[0].ID.String(), helmManifest)
				if err != nil {
					errSaveAppSourceContent := ErrApplicationSourceContent(err, "persist")
					h.log.Error(errSaveAppSourceContent)
					
					event := eb.WithSeverity(events.Warning).WithMetadata(map[string]interface{}{
						"error": errSaveAppSourceContent,
					}).WithDescription(fmt.Sprintf("Error saving application source content %s", application.Name)).Build()
	
					_ = h.provider.PersistEvent(event)
					go h.eb.Publish(*h.userID, event)
					return
				}
			}
			event := eb.WithSeverity(events.Informational).WithDescription(fmt.Sprintf("Helm chart %s/%s imported as meshery application %s", chart.Name, chart.Version, application.Name)).Build()
			_ = h.provider.PersistEvent(event)
			go h.eb.Publish(*h.userID, event)
		}(*chart)
	}
	wg.Wait()
	return msgs
}

func convertAppToPattern(pattern *patterns.Pattern, repoURL, name, version string) (*models.MesheryApplication, error) {
	data, err := yaml.Marshal(pattern)
	if err != nil {
		return nil, models.ErrUnmarshal(err, "application")
	}

	application := &models.MesheryApplication{
		Name: name,
		ApplicationFile: string(data),
		Type: sql.NullString{
			String: string(models.HelmChart),
			Valid:  true,
		},
		Location: map[string]interface{}{
			"type":   "http",
			"host":   repoURL,
			"path":   fmt.Sprintf("%s/%s", name, version),
			"branch": "",
		},
	}
	return application, nil
}

func sanitizeURL(repoURL string) string {
	for _, val := range helmIndexes {
		repoURL, _ = strings.CutSuffix(repoURL, val)
	}
	return repoURL
}

func extractURL(url interface{}, eb *events.EventBuilder) (string, *events.Event) {
	repoURL, ok := url.(string)
	if !ok || repoURL == "" {
		err := fmt.Errorf("invalid request body")
		meshkitErr := ErrInvalidRepoURL(err)
		description := fmt.Sprintf("Invalid repository url provided %v.", repoURL)
		eb.WithDescription(description).WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": meshkitErr,
		})
		event := eb.Build()
		return "", event
	}
	return repoURL, nil
}

func getHelmIndex(repoURL string, eb *events.EventBuilder) (*helm.HelmIndex){
	helmEntries, err := helm.CreateHelmIndex(repoURL)
	if err != nil {
		description := fmt.Sprintf("Cannot read helm charts in the repository %s", repoURL)
		eb.WithDescription(description).WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": err,
		})
		return nil
	}
	return helmEntries
}

func(h *Helm) getHelmEntries(repoURL string, eb *events.EventBuilder, res http.ResponseWriter) {
	helmIndex := getHelmIndex(repoURL, eb)
	if helmIndex == nil {
		event := eb.Build()
		
		err := h.provider.PersistEvent(event)
		if err != nil {
			h.log.Warn(fmt.Errorf("unable to create event: %s", err.Error()))
		}

		res.WriteHeader(http.StatusInternalServerError)
		data, _ := json.Marshal(event)
		_, err = res.Write(data)
		if err != nil {
			h.log.Error(ErrWriteResponse)
		}
		return
	}

	helmEntries := []helm.HelmEntryMetadata{}

	for key := range helmIndex.Entries { // key here is the chart name, and val corresponds to the chart metadadata
		helmEntries = append(helmEntries, helmIndex.Entries[key]...)
	}
	_ = json.NewEncoder(res).Encode(helmEntries)
}