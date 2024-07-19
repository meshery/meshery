package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/helpers"
	"github.com/layer5io/meshery/server/helpers/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/models/events"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha2"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1beta1"
	"github.com/layer5io/meshkit/models/meshmodel/entity"
	meshkitutils "github.com/layer5io/meshkit/utils"
)

func (h *Handler) sendSuccessResponse(rw http.ResponseWriter, userID uuid.UUID, provider models.Provider, message string, errMsg string, response *models.RegisterMeshmodelAPIResponse) {
	if errMsg != "" {
		if message != "" {
			h.log.Info(message + " and " + errMsg)
			response.ErrMsg = message + "\n" + errMsg
		} else {
			h.log.Info(errMsg)
			response.ErrMsg = errMsg
		}
	} else {
		response.ErrMsg = message
		h.log.Info(response.ErrMsg)
	}

	if len(response.EntityTypeSummary.UnsuccessfulEntityNameWithError) > 0 {
		h.log.Info("Unsuccessful Entities: ", response.EntityTypeSummary.UnsuccessfulEntityNameWithError)
	}
	h.sendFileEvent(userID, provider, response)
	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(response); err != nil {
		h.log.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func createTempFile(dirPath string) (*os.File, error) {
	tempFile, err := os.CreateTemp("", "upload-*.tar.gz")
	if err != nil {
		return nil, ErrCreateFile(err, "/tmp/upload-*.tar.gz")
	}

	if _, err = tempFile.Write([]byte(dirPath)); err != nil {
		return nil, meshkitutils.ErrWriteFile(err, tempFile.Name())
	}
	return tempFile, nil
}

func processUploadedFile(filePath string, tempDir string, h *Handler, response *models.RegisterMeshmodelAPIResponse, provider models.Provider) error {

	if err := utils.ExtractFile(filePath, tempDir); err != nil {
		h.sendErrorEvent(uuid.Nil, provider, "Error creating temp directory", err)
		return err
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	if err := filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return meshkitutils.ErrFileWalkDir(err, path)
		}
		if !info.IsDir() {
			wg.Add(1)
			go func(path string) {
				defer wg.Done()
				processFile(path, h, &mu, response, provider)
			}(path)
		}
		return nil
	}); err != nil {
		return err
	}
	wg.Wait()
	return nil
}

func processFile(path string, h *Handler, mu *sync.Mutex, response *models.RegisterMeshmodelAPIResponse, provider models.Provider) {
	if meshkitutils.IsZip(path) || meshkitutils.IsTarGz(path) {
		newTempDir, err := os.MkdirTemp("", "nested-extracted-")
		if err != nil {
			incrementCounter(mu, &response.EntityCount.TotalErrCount)
			h.log.Error(meshkitutils.ErrCreateDir(err, "Error creating nested temp directory"))
			mu.Lock()
			addUnsuccessfulEntry(path, response, meshkitutils.ErrCreateDir(err, "Error creating nested temp directory"), "")
			mu.Unlock()
			return
		}
		defer os.RemoveAll(newTempDir)
		if err := processUploadedFile(path, newTempDir, h, response, provider); err != nil {
			incrementCounter(mu, &response.EntityCount.TotalErrCount)
			h.log.Error(err)
			mu.Lock()
			addUnsuccessfulEntry(path, response, err, "")
			mu.Unlock()
			return
		}
		return
	}

	content, err := os.ReadFile(path)
	if err != nil {
		incrementCounter(mu, &response.EntityCount.TotalErrCount)
		h.log.Error(meshkitutils.ErrReadFile(err, path))
		mu.Lock()
		addUnsuccessfulEntry(path, response, meshkitutils.ErrReadFile(err, path), "")
		mu.Unlock()
		return
	}

	entityType, err := meshkitutils.FindEntityType(content)
	if err != nil {
		incrementCounter(mu, &response.EntityCount.TotalErrCount)
		mu.Lock()
		addUnsuccessfulEntry(path, response, err, "")
		mu.Unlock()
	}
	if entityType != "" {
		path, err := RegisterEntity(content, entityType, h, response, mu)
		if err != nil {
			incrementCountersOnErr(mu, entityType, response)
			h.log.Error(err)
			mu.Lock()
			addUnsuccessfulEntry(path, response, err, string(entityType))
			mu.Unlock()
		} else {
			if path != "" {
				incrementCountersOnSuccess(mu, entityType, &response.EntityCount.CompCount, &response.EntityCount.RelCount, &response.EntityCount.ModelCount)
				addSuccessfulEntry(content, entityType, response)
			}

		}
	}
}

func addSuccessfulEntry(content []byte, entityType entity.EntityType, response *models.RegisterMeshmodelAPIResponse) {
	switch entityType {
	case entity.ComponentDefinition:
		var c v1beta1.ComponentDefinition
		if err := meshkitutils.Unmarshal(string(content), &c); err == nil {
			entry := map[string]interface{}{
				"Model":       c.Model,
				"Metadata":    c.Metadata,
				"DisplayName": c.DisplayName,
				"Version":     c.Model.Version,
			}
			response.EntityTypeSummary.SuccessfulComponents = append(response.EntityTypeSummary.SuccessfulComponents, entry)
			response.ModelName = append(response.ModelName, c.Model.DisplayName)
		}
	case entity.RelationshipDefinition:
		var r v1alpha2.RelationshipDefinition
		if err := meshkitutils.Unmarshal(string(content), &r); err == nil {
			entry := map[string]interface{}{
				"Model":   r.Model,
				"Kind":    r.Kind,
				"Subtype": r.SubType,
				// "RelationshipType": r.RelationshipType, //future when we support type
			}
			response.EntityTypeSummary.SuccessfulRelationships = append(response.EntityTypeSummary.SuccessfulRelationships, entry)
			response.ModelName = append(response.ModelName, r.Model.DisplayName)
		}
	case entity.Model:
		var m v1beta1.Model
		if err := meshkitutils.Unmarshal(string(content), &m); err == nil {
			entry := map[string]interface{}{
				"Model":       m.Model,
				"Metadata":    m.Metadata,
				"DisplayName": m.DisplayName,
				"Version":     m.Model.Version,
			}
			response.EntityTypeSummary.SuccessfulModels = append(response.EntityTypeSummary.SuccessfulModels, entry)
			response.ModelName = append(response.ModelName, m.DisplayName)
		}
	}
}

func addUnsuccessfulEntry(path string, response *models.RegisterMeshmodelAPIResponse, err error, entityType string) {
	filename := path

	filename = filepath.Base(path)

	entry := map[string]interface{}{
		"name":       filename,
		"entityType": entityType,
		"error":      err,
	}
	response.EntityTypeSummary.UnsuccessfulEntityNameWithError = append(response.EntityTypeSummary.UnsuccessfulEntityNameWithError, entry)
}

func incrementCounter(mu *sync.Mutex, counter *int) {
	mu.Lock()
	defer mu.Unlock()
	*counter++
}

func incrementCountersOnErr(mu *sync.Mutex, entityType entity.EntityType, response *models.RegisterMeshmodelAPIResponse) {
	mu.Lock()
	defer mu.Unlock()
	response.EntityCount.TotalErrCount++
	if entityType == entity.ComponentDefinition {
		response.EntityCount.ErrCompCount++
	} else if entityType == entity.RelationshipDefinition {
		response.EntityCount.ErrRelCount++
	} else if entityType == entity.Model {
		response.EntityCount.ErrModelCount++
	}
}

func incrementCountersOnSuccess(mu *sync.Mutex, entityType entity.EntityType, compCount *int, relCount *int, modelCount *int) {
	mu.Lock()
	defer mu.Unlock()
	if entityType == entity.ComponentDefinition {
		*compCount++
	} else if entityType == entity.RelationshipDefinition {
		*relCount++
	} else if entityType == entity.Model {
		*modelCount++
	}
}

func (h *Handler) sendErrorEvent(userID uuid.UUID, provider models.Provider, description string, err error) {
	event := events.NewEvent().ActedUpon(userID).FromUser(userID).FromSystem(*h.SystemID).WithAction("register").WithSeverity(events.Error).WithDescription(description).WithMetadata(map[string]interface{}{
		"error": err,
	}).Build()
	_ = provider.PersistEvent(event)
	go h.config.EventBroadcaster.Publish(userID, event)
}

func ModelNames(response *models.RegisterMeshmodelAPIResponse) string {
	msg := ""
	seen := make(map[string]bool) // map to track seen model names

	for _, model := range response.ModelName {
		if !seen[model] {
			msg += model + " "
			seen[model] = true
		}
	}

	return msg
}
func (h *Handler) sendFileEvent(userID uuid.UUID, provider models.Provider, response *models.RegisterMeshmodelAPIResponse) {
	s := ModelNames(response)
	description := fmt.Sprintf("Imported model(s) %s", s)
	metadata := map[string]interface{}{
		"ImportedModelName":               s,
		"ImportedComponent":               response.EntityTypeSummary.SuccessfulComponents,
		"ImportedRelationship":            response.EntityTypeSummary.SuccessfulRelationships,
		"ImportedModel":                   response.EntityTypeSummary.SuccessfulModels,
		"UnsuccessfulEntityNameWithError": response.EntityTypeSummary.UnsuccessfulEntityNameWithError,
	}
	eventType := events.Informational
	if response.EntityCount.CompCount == 0 && response.EntityCount.RelCount == 0 && response.EntityCount.ModelCount == 0 {
		eventType = events.Error
	}
	event := events.NewEvent().
		ActedUpon(userID).
		FromUser(userID).
		FromSystem(*h.SystemID).
		WithAction("register").
		WithSeverity(eventType).
		WithDescription(description).
		WithMetadata(metadata).
		Build()
	_ = provider.PersistEvent(event)
	go h.config.EventBroadcaster.Publish(userID, event)
}

func RegisterEntity(content []byte, entityType entity.EntityType, h *Handler, response *models.RegisterMeshmodelAPIResponse, mu *sync.Mutex) (string, error) {
	switch entityType {
	case entity.ComponentDefinition:
		var c v1beta1.ComponentDefinition
		if err := meshkitutils.Unmarshal(string(content), &c); err != nil {
			return "", err
		}
		utils.WriteSVGsOnFileSystem(&c)
		isRegistrantError, isModelError, err := h.registryManager.RegisterEntity(v1beta1.Host{Hostname: c.Model.Registrant.Hostname}, &c)
		helpers.HandleError(v1beta1.Host{Hostname: c.Model.Registrant.Hostname}, &c, err, isModelError, isRegistrantError)

		return c.DisplayName, err
	case entity.RelationshipDefinition:
		var r v1alpha2.RelationshipDefinition
		if err := meshkitutils.Unmarshal(string(content), &r); err != nil {
			return "", meshkitutils.ErrUnmarshal(err)
		}
		isRegistrantError, isModelError, err := h.registryManager.RegisterEntity(v1beta1.Host{Hostname: r.Model.Registrant.Hostname}, &r)
		helpers.HandleError(v1beta1.Host{Hostname: r.Model.Registrant.Hostname}, &r, err, isModelError, isRegistrantError)

		return r.Kind, err
	case entity.Model:
		var m v1beta1.Model
		checkBool := false
		if err := meshkitutils.Unmarshal(string(content), &m); err != nil {
			err = meshkitutils.ErrUnmarshal(err)
			return "", err
		}
		components := m.Components
		var rel []v1alpha2.RelationshipDefinition
		relationships, _ := meshkitutils.Cast[string](m.Relationships)
		if relationships != "" {
			if err := meshkitutils.Unmarshal((relationships), &rel); err != nil {
				return "", err
			}
		}
		if len(components) > 0 || len(rel) > 0 {
			checkBool = true
		}
		for _, comp := range components {
			utils.WriteSVGsOnFileSystem(&comp)
			isRegistrantError, isModelError, err := h.registryManager.RegisterEntity(v1beta1.Host{Hostname: comp.Model.Registrant.Hostname}, &comp)
			helpers.HandleError(v1beta1.Host{Hostname: comp.Model.Registrant.Hostname}, &comp, err, isModelError, isRegistrantError)
			if err != nil {
				componetName := m.DisplayName + "-" + comp.DisplayName
				mu.Lock()
				incrementCountersOnErr(mu, entity.ComponentDefinition, response)
				mu.Unlock()
				addUnsuccessfulEntry(componetName, response, err, string(entity.ComponentDefinition))
				continue
			}
			mu.Lock()
			incrementCountersOnSuccess(mu, entity.ComponentDefinition, &response.EntityCount.CompCount, &response.EntityCount.RelCount, &response.EntityCount.ModelCount)
			mu.Unlock()
			addSuccessfulEntry(content, entity.ComponentDefinition, response)
		}
		// for _, rel := range relationships {
		// 	//future when we fix from interface to array of relationship definition

		// }
		for _, r := range rel {
			isRegistrantError, isModelError, err := h.registryManager.RegisterEntity(v1beta1.Host{Hostname: r.Model.Registrant.Hostname}, &r)
			helpers.HandleError(v1beta1.Host{Hostname: r.Model.Registrant.Hostname}, &r, err, isModelError, isRegistrantError)
			if err != nil {
				relName := m.DisplayName + "-" + r.Kind
				mu.Lock()
				incrementCountersOnErr(mu, entity.RelationshipDefinition, response)
				mu.Unlock()
				addUnsuccessfulEntry(relName, response, err, string(entity.RelationshipDefinition))
				continue

			}
		}
		if checkBool {
			if response.EntityCount.CompCount == 0 || response.EntityCount.RelCount == 0 {
				return "", nil
			}
		}
		return m.DisplayName, nil
	case entity.PolicyDefinition:
		//future when we support policy
		return "", nil
	}
	return "", meshkitutils.ErrInvalidSchemaVersion
}

func writeMessageString(compCount, relCount int) strings.Builder {
	var message strings.Builder
	if compCount > 0 || relCount > 0 {
		message.WriteString("Total ")
	}
	if compCount > 0 {
		message.WriteString(fmt.Sprintf(" Components imported: %d", compCount))
	}
	if relCount > 0 {
		if message.Len() > 0 {
			message.WriteString(" and ")
		}
		message.WriteString(fmt.Sprintf("Relationships imported: %d", relCount))
	}
	return message
}

func ErrMsgContruct(totalErrCount int, errCompCount int, errRelCount int) string {
	msg := fmt.Sprintf("Failed to import %d entity(s)", totalErrCount)
	if errCompCount > 0 && errRelCount > 0 {
		msg = fmt.Sprintf("%s %d Component(s) and %d Relationship(s)", msg, errCompCount, errRelCount)
	} else if errCompCount > 0 {
		msg = fmt.Sprintf("%s %d Component(s)", msg, errCompCount)
	} else if errRelCount > 0 {
		msg = fmt.Sprintf("%s %d Relationship(s)", msg, errRelCount)
	}
	return msg
}
