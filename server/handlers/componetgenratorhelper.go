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
	"github.com/mitchellh/mapstructure"
)

func (h *Handler) handleError(rw http.ResponseWriter, err error, logMsg string) {
	h.log.Error(err)
	http.Error(rw, logMsg, http.StatusInternalServerError)
}

func (h *Handler) sendSuccessResponse(rw http.ResponseWriter, userID uuid.UUID, provider models.Provider, message string, errMsg string, response *models.RegistryAPIResponse) {
	if errMsg != "" {
		if message != "" {
			response.ErrMsg = message + ". " + errMsg
			h.log.Info(response.ErrMsg)
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

func processUploadedFile(filePath string, tempDir string, h *Handler, response *models.RegistryAPIResponse, provider models.Provider) error {

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

func processFile(path string, h *Handler, mu *sync.Mutex, response *models.RegistryAPIResponse, provider models.Provider) {
	if meshkitutils.IsZip(path) || meshkitutils.IsTarGz(path) {
		newTempDir, err := os.MkdirTemp("", "nested-extracted-")
		if err != nil {
			incrementCounter(mu, &response.EntityCount.TotalErrCount)
			h.log.Error(meshkitutils.ErrCreateDir(err, "Error creating nested temp directory"))
			addUnsuccessfulEntry(path, response, meshkitutils.ErrCreateDir(err, "Error creating nested temp directory"), "")
			return
		}
		defer os.RemoveAll(newTempDir)
		if err := processUploadedFile(path, newTempDir, h, response, provider); err != nil {
			incrementCounter(mu, &response.EntityCount.TotalErrCount)
			h.log.Error(err)
			addUnsuccessfulEntry(path, response, err, "")
			return
		}
		return
	}

	content, err := os.ReadFile(path)
	if err != nil {
		incrementCounter(mu, &response.EntityCount.TotalErrCount)
		h.log.Error(meshkitutils.ErrReadFile(err, path))
		addUnsuccessfulEntry(path, response, meshkitutils.ErrReadFile(err, path), "")
		return
	}
	processFileToRegistry(response, content, mu, path, h)
}
func processFileToRegistry(response *models.RegistryAPIResponse, content []byte, mu *sync.Mutex, path string, h *Handler) {
	entityType, err := meshkitutils.FindEntityType(content)
	if err != nil {
		incrementCounter(mu, &response.EntityCount.TotalErrCount)
		addUnsuccessfulEntry(path, response, err, "")
	}
	if entityType != "" {
		path, err := RegisterEntity(content, entityType, h, response, mu)
		if err != nil {
			incrementCountersOnErr(mu, entityType, response)
			h.log.Error(err)
			addUnsuccessfulEntry(path, response, err, string(entityType))
		} else {
			if path != "" {
				incrementCountersOnSuccess(mu, entityType, &response.EntityCount.CompCount, &response.EntityCount.RelCount, &response.EntityCount.ModelCount)
				addSuccessfulEntry(content, entityType, response)
			}

		}
	}
}
func addSuccessfulEntry(content []byte, entityType entity.EntityType, response *models.RegistryAPIResponse) {
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
				"Model":     r.Model,
				"Kind":      r.Kind,
				"Subtype":   r.SubType,
				"Selectors": r.Selectors,
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

func addUnsuccessfulEntry(path string, response *models.RegistryAPIResponse, err error, entityType string) {
	filename := filepath.Base(path)
	entryFound := false

	// Loop through existing entries to check if the error already exists
	for i, entry := range response.EntityTypeSummary.UnsuccessfulEntityNameWithError {
		if entryMap, ok := entry.(map[string]interface{}); ok {
			if existingErr, ok := entryMap["error"]; ok && existingErr == err {
				// Append new filename and entityType to the existing entry
				if names, ok := entryMap["name"].([]string); ok {
					entryMap["name"] = append(names, filename)
				} else {
					entryMap["name"] = []string{entryMap["name"].(string), filename}
				}

				if entityTypes, ok := entryMap["entityType"].([]string); ok {
					entryMap["entityType"] = append(entityTypes, entityType)
				} else {
					entryMap["entityType"] = []string{entryMap["entityType"].(string), entityType}
				}
				response.EntityTypeSummary.UnsuccessfulEntityNameWithError[i] = entryMap
				entryFound = true
				break
			}
		}
	}

	// If error not found, create a new entry
	if !entryFound {
		entry := map[string]interface{}{
			"name":       []string{filename},
			"entityType": []string{entityType},
			"error":      err,
		}
		response.EntityTypeSummary.UnsuccessfulEntityNameWithError = append(response.EntityTypeSummary.UnsuccessfulEntityNameWithError, entry)
	}
}

func incrementCounter(mu *sync.Mutex, counter *int) {
	mu.Lock()
	defer mu.Unlock()
	*counter++
}

func incrementCountersOnErr(mu *sync.Mutex, entityType entity.EntityType, response *models.RegistryAPIResponse) {
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

func ModelNames(response *models.RegistryAPIResponse) string {
	var builder strings.Builder
	seen := make(map[string]bool) // map to track seen model names

	for _, model := range response.ModelName {
		if model != "" {
			if !seen[model] {
				if builder.Len() > 0 {
					builder.WriteString(", ")
				}
				builder.WriteString(model)
				seen[model] = true
			}
		}
	}
	return builder.String()
}

func (h *Handler) sendFileEvent(userID uuid.UUID, provider models.Provider, response *models.RegistryAPIResponse) {
	s := ModelNames(response)
	description := fmt.Sprintf("Imported model(s) %s", s)
	metadata := map[string]interface{}{
		"ImportedModelName":               s,
		"ImportedComponent":               response.EntityTypeSummary.SuccessfulComponents,
		"ImportedRelationship":            response.EntityTypeSummary.SuccessfulRelationships,
		"UnsuccessfulEntityNameWithError": response.EntityTypeSummary.UnsuccessfulEntityNameWithError,
		"ModelImportMessage":              response.ErrMsg,
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

func RegisterEntity(content []byte, entityType entity.EntityType, h *Handler, response *models.RegistryAPIResponse, mu *sync.Mutex) (string, error) {
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
		var rels []v1alpha2.RelationshipDefinition
		components := m.Components
		if m.Relationships != nil {
			slice, ok := m.Relationships.([]interface{})
			if !ok {
				return "", meshkitutils.ErrUnmarshal(meshkitutils.ErrInvalidSchemaVersion)
			}
			rels = make([]v1alpha2.RelationshipDefinition, 0, len(slice))
			for _, v := range slice {
				mapVal, ok := v.(map[string]interface{})
				if !ok {
					return "", meshkitutils.ErrUnmarshal(meshkitutils.ErrInvalidSchemaVersion)
				}
				var rel v1alpha2.RelationshipDefinition
				if err := mapstructure.Decode(mapVal, &rel); err != nil {
					return "", fmt.Errorf("failed to decode relationship: %v", err)
				}
				if schemaVersion, ok := mapVal["schemaVersion"].(string); ok {
					rel.SchemaVersion = schemaVersion
				}
				rels = append(rels, rel)
			}
		}

		if len(components) > 0 || len(rels) > 0 {
			checkBool = true
		}
		for _, comp := range components {
			compBytes, _ := json.Marshal(comp)
			_, err := meshkitutils.FindEntityType([]byte(compBytes))
			if err != nil {
				incrementCountersOnErr(mu, entity.ComponentDefinition, response)
				addUnsuccessfulEntry(m.DisplayName, response, err, string(entity.ComponentDefinition))
				continue
			}
			utils.WriteSVGsOnFileSystem(&comp)
			isRegistrantError, isModelError, err := h.registryManager.RegisterEntity(v1beta1.Host{Hostname: comp.Model.Registrant.Hostname}, &comp)
			helpers.HandleError(v1beta1.Host{Hostname: comp.Model.Registrant.Hostname}, &comp, err, isModelError, isRegistrantError)
			if err != nil {
				componetName := m.DisplayName + " of component name " + comp.DisplayName
				incrementCountersOnErr(mu, entity.ComponentDefinition, response)
				addUnsuccessfulEntry(componetName, response, err, string(entity.ComponentDefinition))
				continue
			}
			incrementCountersOnSuccess(mu, entity.ComponentDefinition, &response.EntityCount.CompCount, &response.EntityCount.RelCount, &response.EntityCount.ModelCount)

			addSuccessfulEntry(compBytes, entity.ComponentDefinition, response)

		}
		for _, r := range rels {
			relBytes, _ := json.Marshal(r)
			_, err := meshkitutils.FindEntityType([]byte(relBytes))
			if err != nil {
				incrementCountersOnErr(mu, entity.RelationshipDefinition, response)
				addUnsuccessfulEntry(m.DisplayName, response, err, string(entity.RelationshipDefinition))
				continue
			}
			isRegistrantError, isModelError, err := h.registryManager.RegisterEntity(v1beta1.Host{Hostname: r.Model.Registrant.Hostname}, &r)
			helpers.HandleError(v1beta1.Host{Hostname: r.Model.Registrant.Hostname}, &r, err, isModelError, isRegistrantError)
			if err != nil {
				relName := m.DisplayName + " of Kind :" + r.Kind
				incrementCountersOnErr(mu, entity.RelationshipDefinition, response)
				addUnsuccessfulEntry(relName, response, err, string(entity.RelationshipDefinition))
				continue
			}
			incrementCountersOnSuccess(mu, entity.RelationshipDefinition, &response.EntityCount.CompCount, &response.EntityCount.RelCount, &response.EntityCount.ModelCount)
			addSuccessfulEntry(relBytes, entity.RelationshipDefinition, response)

		}
		if checkBool {
			if response.EntityCount.CompCount == 0 && response.EntityCount.RelCount == 0 {
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

func writeMessageString(response *models.RegistryAPIResponse) strings.Builder {
	var message strings.Builder
	if response.EntityCount.ModelCount > 0 {
		if response.EntityCount.CompCount == 0 && response.EntityCount.RelCount == 0 {
			message.WriteString("Model(s) won't be registered because there was no Component(s) or Relationship(s)")
			return message
		}
		modelName := ModelNames(response)
		message.WriteString(fmt.Sprintf("Model(s) imported: %s. ", modelName))
	}
	if response.EntityCount.CompCount > 0 || response.EntityCount.RelCount > 0 {
		message.WriteString("Total ")
	}
	if response.EntityCount.CompCount > 0 {
		message.WriteString(fmt.Sprintf("Component(s) imported: %d", response.EntityCount.CompCount))
	}
	if response.EntityCount.RelCount > 0 && response.EntityCount.CompCount > 0 {
		if message.Len() > 0 {
			message.WriteString(" and ")
		}
		message.WriteString(fmt.Sprintf("Relationship(s) imported: %d", response.EntityCount.RelCount))
	}
	return message
}

func ErrMsgContruct(totalErrCount int, errCompCount int, errRelCount int) string {
	msg := fmt.Sprintf("Failed to import %d entity(s) that is", totalErrCount)
	if errCompCount > 0 && errRelCount > 0 {
		msg = fmt.Sprintf("%s %d Component(s) and %d Relationship(s)", msg, errCompCount, errRelCount)
	} else if errCompCount > 0 {
		msg = fmt.Sprintf("%s %d Component(s)", msg, errCompCount)
	} else if errRelCount > 0 {
		msg = fmt.Sprintf("%s %d Relationship(s)", msg, errRelCount)
	}
	return msg
}
func findTarFile(directory string) (string, error) {
	var tarFilePath string
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return meshkitutils.ErrFileWalkDir(err, path)
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".tar.gz") {
			tarFilePath = path
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	if tarFilePath == "" {
		return "", ErrNoTarInsideOCi(err)
	}
	return tarFilePath, nil
}
