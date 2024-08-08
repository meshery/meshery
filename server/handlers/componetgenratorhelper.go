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
	"gopkg.in/yaml.v2"
)

func (h *Handler) handleError(rw http.ResponseWriter, err error, logMsg string) {
	h.log.Error(err)
	http.Error(rw, logMsg, http.StatusInternalServerError)
}

func (h *Handler) sendSuccessResponse(rw http.ResponseWriter, userID uuid.UUID, provider models.Provider, message string, errMsg string, response *models.RegistryAPIResponse) {
	if errMsg != "" {
		if message != "" {
			response.ErrMsg = message + ", " + errMsg
			h.log.Info(response.ErrMsg)
		} else {
			h.log.Info(errMsg)
			response.ErrMsg = "Import " + errMsg + "."
		}
	} else {
		response.ErrMsg = message
		h.log.Info(response.ErrMsg)
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

func processUploadedFile(filePath string, tempDir string, h *Handler, response *models.RegistryAPIResponse, provider models.Provider, checkOCI bool) error {

	if !checkOCI {
		if err := utils.ExtractFile(filePath, tempDir); err != nil {
			h.sendErrorEvent(uuid.Nil, provider, "Error creating temp directory", err)
			return err
		}
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
		if err := processUploadedFile(path, newTempDir, h, response, provider, false); err != nil {
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
	var jsonData interface{}
	if err := yaml.Unmarshal(content, &jsonData); err == nil {
		jsonData = utils.ConvertToJSONCompatible(jsonData)
		convertedContent, err := json.Marshal(jsonData)
		if err == nil {
			content = convertedContent
		}
	}
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
	// Helper function to check if a display name is already present
	isDisplayNamePresent := func(displayName string) bool {
		for _, name := range response.ModelName {
			if name == displayName {
				return true
			}
		}
		return false
	}

	switch entityType {
	case entity.ComponentDefinition:
		var c v1beta1.ComponentDefinition
		if err := meshkitutils.Unmarshal(string(content), &c); err == nil {
			entry := map[string]interface{}{
				"Model":       c.Model,
				"Metadata":    c.Component.Kind,
				"DisplayName": c.DisplayName,
				"Version":     c.Model.Version,
			}
			response.EntityTypeSummary.SuccessfulComponents = append(response.EntityTypeSummary.SuccessfulComponents, entry)
			if !isDisplayNamePresent(c.Model.Name) {
				response.ModelName = append(response.ModelName, c.Model.Name)
			}
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
			if !isDisplayNamePresent(r.Model.Name) {
				response.ModelName = append(response.ModelName, r.Model.Name)
			}
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
			if !isDisplayNamePresent(m.Name) {
				response.ModelName = append(response.ModelName, m.Name)
			}
		}
	}
}

func addUnsuccessfulEntry(path string, response *models.RegistryAPIResponse, err error, entityType string) {
	isDisplayNamePresent := func(displayName string) bool {
		for _, name := range response.ModelName {
			if name == displayName {
				return true
			}
		}
		return false
	}
	filename := filepath.Base(path)
	entryFound := false
	if entityType == "" {
		entityType = "Unknown"
	}
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
	if !isDisplayNamePresent(filename) {
		response.ModelName = append(response.ModelName, filename)
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
	// Initialize metadata map
	metadata := map[string]interface{}{
		"ModelImportMessage": response.ErrMsg,
		"ModelDetails":       map[string]interface{}{},
	}

	// Add successful components and relationships to their respective model entries
	for _, modelName := range response.ModelName {
		modelData := map[string]interface{}{
			"Components":    []interface{}{},
			"Relationships": []interface{}{},
			"Errors":        []interface{}{},
		}

		for _, component := range response.EntityTypeSummary.SuccessfulComponents {
			if component["Model"].(v1beta1.Model).Name == modelName {
				modelData["Components"] = append(modelData["Components"].([]interface{}), component)
			}
		}

		for _, relationship := range response.EntityTypeSummary.SuccessfulRelationships {
			if relationship["Model"].(v1beta1.Model).Name == modelName {
				modelData["Relationships"] = append(modelData["Relationships"].([]interface{}), relationship)
			}
		}
		for _, errorEntry := range response.EntityTypeSummary.UnsuccessfulEntityNameWithError {
			if errorMap, ok := errorEntry.(map[string]interface{}); ok {
				if names, ok := errorMap["name"].([]string); ok {
					for _, name := range names {
						if name == modelName {
							modelData["Errors"] = append(modelData["Errors"].([]interface{}), errorMap)
							break
						}
					}
				}
			}
		}
		metadata["ModelDetails"].(map[string]interface{})[modelName] = modelData
	}

	eventType := events.Informational
	if response.EntityCount.CompCount == 0 && response.EntityCount.RelCount == 0 && response.EntityCount.ModelCount == 0 {
		eventType = events.Error
	} else if response.EntityCount.TotalErrCount > 0 && (response.EntityCount.CompCount > 0 || response.EntityCount.RelCount > 0 || response.EntityCount.ModelCount > 0) {
		eventType = events.Warning
	}
	description := getFirst42Chars(response.ErrMsg)
	if len(description) == 42 {
		description = description + "..."
	}
	if response.EntityCount.TotalErrCount == 0 {
		eventType = events.Success
	}
	event := events.NewEvent().
		ActedUpon(userID).
		FromUser(userID).
		FromSystem(*h.SystemID).
		WithAction("register").
		WithDescription(description).
		WithSeverity(eventType).
		WithMetadata(metadata).
		Build()

	_ = provider.PersistEvent(event)
	go h.config.EventBroadcaster.Publish(userID, event)
}
func getFirst42Chars(s string) string {
	if len(s) > 42 {
		return s[:42]
	}
	return s
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
			rels = make([]v1alpha2.RelationshipDefinition, len(slice))
			for i, v := range slice {
				bytes, err := json.Marshal(v)
				if err != nil {
					continue
				}
				var rel v1alpha2.RelationshipDefinition
				err = json.Unmarshal(bytes, &rel)
				if err != nil {
					continue
				}
				rels[i] = rel
			}
		}

		if len(components) > 0 || len(rels) > 0 {
			checkBool = true
		}
		for _, comp := range components {
			comp.Model = m
			compBytes, _ := json.Marshal(comp)
			_, err := meshkitutils.FindEntityType([]byte(compBytes))
			if err != nil {
				incrementCountersOnErr(mu, entity.ComponentDefinition, response)
				addUnsuccessfulEntry(m.Name, response, err, string(entity.ComponentDefinition))
				continue
			}
			utils.WriteSVGsOnFileSystem(&comp)
			isRegistrantError, isModelError, err := h.registryManager.RegisterEntity(v1beta1.Host{Hostname: comp.Model.Registrant.Hostname}, &comp)
			helpers.HandleError(v1beta1.Host{Hostname: comp.Model.Registrant.Hostname}, &comp, err, isModelError, isRegistrantError)
			if err != nil {
				incrementCountersOnErr(mu, entity.ComponentDefinition, response)
				addUnsuccessfulEntry(m.Name, response, err, string(entity.ComponentDefinition))
				continue
			}
			incrementCountersOnSuccess(mu, entity.ComponentDefinition, &response.EntityCount.CompCount, &response.EntityCount.RelCount, &response.EntityCount.ModelCount)

			addSuccessfulEntry(compBytes, entity.ComponentDefinition, response)

		}
		for _, r := range rels {
			r.Model = m
			relBytes, _ := json.Marshal(r)
			_, err := meshkitutils.FindEntityType([]byte(relBytes))
			if err != nil {
				incrementCountersOnErr(mu, entity.RelationshipDefinition, response)
				addUnsuccessfulEntry(m.Name, response, err, string(entity.RelationshipDefinition))
				continue
			}
			isRegistrantError, isModelError, err := h.registryManager.RegisterEntity(v1beta1.Host{Hostname: r.Model.Registrant.Hostname}, &r)
			helpers.HandleError(v1beta1.Host{Hostname: r.Model.Registrant.Hostname}, &r, err, isModelError, isRegistrantError)
			if err != nil {
				incrementCountersOnErr(mu, entity.RelationshipDefinition, response)
				addUnsuccessfulEntry(m.Name, response, err, string(entity.RelationshipDefinition))
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
		if response.EntityCount.CompCount == 0 && response.EntityCount.RelCount == 0 && response.EntityCount.TotalErrCount == 0 {
			message.WriteString("Model won't be registered because there was no Component(s) or Relationship(s).")
			return message
		}
		modelName := ModelNames(response)
		model := "model"
		if response.EntityCount.ModelCount > 1 {
			model = "models"
		}
		message.WriteString(fmt.Sprintf("Imported %s %s ", model, modelName))
	}
	if response.EntityCount.CompCount > 0 || response.EntityCount.RelCount > 0 {
		message.WriteString("(")
	}
	if response.EntityCount.CompCount > 0 {
		component := "component"
		if response.EntityCount.CompCount > 1 {
			component = "components"
		}
		message.WriteString(fmt.Sprintf("%d %s", response.EntityCount.CompCount, component))
	}
	if response.EntityCount.RelCount > 0 && response.EntityCount.CompCount > 0 {
		if message.Len() > 0 {
			message.WriteString(" and ")
		}
		relationship := "relationship"
		if response.EntityCount.RelCount > 1 {
			relationship = "relationships"
		}
		message.WriteString(fmt.Sprintf("%d %s", response.EntityCount.RelCount, relationship))
	}
	if response.EntityCount.CompCount == 0 && response.EntityCount.RelCount > 0 {
		relationship := "relationship"
		if response.EntityCount.RelCount > 1 {
			relationship = "relationships"
		}
		message.WriteString(fmt.Sprintf("%d %s", response.EntityCount.RelCount, relationship))
	}
	if response.EntityCount.CompCount > 0 || response.EntityCount.RelCount > 0 {
		message.WriteString(")")
	}
	return message
}

func ErrMsgContruct(response *models.RegistryAPIResponse) string {
	component := "component"
	if response.EntityCount.ErrCompCount > 1 {
		component = "components"
	}
	relationship := "relationship"
	if response.EntityCount.ErrRelCount > 1 {
		relationship = "relationships"
	}
	model := "model"
	if response.EntityCount.ErrModelCount > 1 {
		model = "models"
	}
	entity := "entity"
	if response.EntityCount.TotalErrCount > 1 {
		entity = "entities"
	}

	msg := fmt.Sprintf("encountered error for %d %s (", response.EntityCount.TotalErrCount, entity)
	componentsPresent := response.EntityCount.ErrCompCount > 0
	relationshipsPresent := response.EntityCount.ErrRelCount > 0
	modelsPresent := response.EntityCount.ErrModelCount > 0
	knownErrors := response.EntityCount.ErrCompCount + response.EntityCount.ErrRelCount + response.EntityCount.ErrModelCount
	unknownErrors := response.EntityCount.TotalErrCount - knownErrors

	// Collect errors in a slice for dynamic message construction
	errors := []string{}
	if modelsPresent {
		errors = append(errors, fmt.Sprintf("%d %s", response.EntityCount.ErrModelCount, model))
	}
	if componentsPresent {
		errors = append(errors, fmt.Sprintf("%d %s", response.EntityCount.ErrCompCount, component))
	}
	if relationshipsPresent {
		errors = append(errors, fmt.Sprintf("%d %s", response.EntityCount.ErrRelCount, relationship))
	}
	if unknownErrors > 0 {
		unknownEntity := "entity"
		if unknownErrors > 1 {
			unknownEntity = "entities"
		}
		errors = append(errors, fmt.Sprintf("%d unknown %s", unknownErrors, unknownEntity))
	}
	if len(errors) > 1 {
		msg += fmt.Sprintf("%s and %s", strings.Join(errors[:len(errors)-1], ", "), errors[len(errors)-1])
	} else if len(errors) == 1 {
		msg += errors[0]
	}
	msg += ")"
	return msg
}
