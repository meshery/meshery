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
	"github.com/layer5io/meshery/server/models"

	mesheryctlUtils "github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshkit/encoding"

	"github.com/layer5io/meshkit/models/events"
	"github.com/layer5io/meshkit/models/meshmodel/entity"
	"github.com/layer5io/meshkit/models/registration"
	meshkitutils "github.com/layer5io/meshkit/utils"
	"github.com/meshery/schemas/models/v1alpha3/relationship"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/model"
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
		var c component.ComponentDefinition
		if err := encoding.Unmarshal((content), &c); err == nil {
			entry := map[string]interface{}{
				"Model":       c.Model.Name,
				"Category":    c.Model.Category,
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
		var r relationship.RelationshipDefinition
		if err := encoding.Unmarshal((content), &r); err == nil {
			entry := map[string]interface{}{
				"Model":            r.Model.Name,
				"Kind":             r.Kind,
				"Subtype":          r.SubType,
				"Selectors":        r.Selectors,
				"RelationshipType": r.RelationshipType,
			}
			response.EntityTypeSummary.SuccessfulRelationships = append(response.EntityTypeSummary.SuccessfulRelationships, entry)
			if !isDisplayNamePresent(r.Model.Name) {
				response.ModelName = append(response.ModelName, r.Model.Name)
			}
		}
	case entity.Model:
		var m model.ModelDefinition
		if err := encoding.Unmarshal((content), &m); err == nil {
			entry := map[string]interface{}{
				"Model":       m.Name,
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
func CreateTemp(filename string, data []byte) (*os.File, error) {
	tempFile, err := os.CreateTemp("", filename)
	if err != nil {
		err = meshkitutils.ErrCreateFile(err, "Error creating temp file")
		return nil, err
	}

	_, err = tempFile.Write(data)
	if err != nil {
		err = meshkitutils.ErrWriteFile(err, filename)
		return nil, err
	}
	err = tempFile.Close()
	if err != nil {
		err = meshkitutils.ErrCloseFile(err)
		return nil, err
	}
	return tempFile, nil
}
func handleRegistrationAndError(registrationHelper registration.RegistrationHelper, mu *sync.Mutex, response *models.RegistryAPIResponse, regErrorStore *models.RegistrationFailureLog) {
	for _, pkg := range registrationHelper.PkgUnits {
		registeredModel := pkg.Model
		for _, comp := range pkg.Components {
			incrementCountersOnSuccess(mu, entity.ComponentDefinition, &response.EntityCount.CompCount, &response.EntityCount.RelCount, &response.EntityCount.ModelCount)
			componentBytes, _ := json.Marshal(comp)
			addSuccessfulEntry(componentBytes, entity.ComponentDefinition, response)
		}
		for _, rel := range pkg.Relationships {
			incrementCountersOnSuccess(mu, entity.RelationshipDefinition, &response.EntityCount.CompCount, &response.EntityCount.RelCount, &response.EntityCount.ModelCount)
			relationshipBytes, _ := json.Marshal(rel)
			addSuccessfulEntry(relationshipBytes, entity.RelationshipDefinition, response)
		}
		modelBytes, _ := json.Marshal(registeredModel)
		incrementCountersOnSuccess(mu, entity.Model, &response.EntityCount.CompCount, &response.EntityCount.RelCount, &response.EntityCount.ModelCount)
		addSuccessfulEntry(modelBytes, entity.Model, response)

	}

	if regErrorStore != nil {
		for _, errEntry := range regErrorStore.GetEntityRegErrors() {
			if errEntry.EntityType == "Unknown" {
				errEntry.EntityType = ""
			}
			incrementCountersOnErr(mu, errEntry.EntityType, response)
			path := errEntry.EntityName
			err := errEntry.Err
			entityTypeStr := string(errEntry.EntityType)
			addUnsuccessfulEntry(path, response, err, entityTypeStr)
		}
	}
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
			if component["Model"] == modelName {
				modelData["Components"] = append(modelData["Components"].([]interface{}), component)
			}
		}

		for _, relationship := range response.EntityTypeSummary.SuccessfulRelationships {
			if relationship["Model"] == modelName {
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
	} else if response.EntityCount.ModelCount > 0 && response.EntityCount.RelCount == 0 && response.EntityCount.CompCount == 0 {
		eventType = events.Warning
	} else if response.EntityCount.TotalErrCount == 0 {
		eventType = events.Success
	}
	description := getFirst42Chars(response.ErrMsg)
	if len(description) == 42 {
		description = description + "..."
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
func (h *Handler) sendEventForImport(userID uuid.UUID, provider models.Provider, compsCount int, modelName string) {
	componentWord := determinePluralWord(compsCount, "component")
	description := fmt.Sprintf("Generated %d %s for model %s", compsCount, componentWord, modelName)

	metadata := map[string]interface{}{
		"Description": fmt.Sprintf("Extracted %v %s for model %s\nModel can be accessed from `.meshery/models`", compsCount, componentWord, modelName),
	}
	event := events.NewEvent().
		ActedUpon(userID).
		FromUser(userID).
		FromSystem(*h.SystemID).
		WithAction("generate").
		WithDescription(description).
		WithSeverity(events.Informational).
		WithMetadata(metadata).
		Build()
	_ = provider.PersistEvent(event)
	go h.config.EventBroadcaster.Publish(userID, event)
}
func writeMessageString(response *models.RegistryAPIResponse) string {
	var message strings.Builder
	if response.EntityCount.ModelCount > 0 {
		if response.EntityCount.CompCount == 0 && response.EntityCount.RelCount == 0 && response.EntityCount.TotalErrCount == 0 {
			message.WriteString("Model won't be registered because there was no Component or Relationship.")
			return message.String()
		}
		modelName := ModelNames(response)
		modelWord := determinePluralWord(response.EntityCount.ModelCount, "model")
		message.WriteString(fmt.Sprintf("Imported %s %s ", modelWord, modelName))
	}
	if response.EntityCount.CompCount > 0 || response.EntityCount.RelCount > 0 {
		message.WriteString("(")
	}
	if response.EntityCount.CompCount > 0 {
		componentWord := determinePluralWord(response.EntityCount.CompCount, "component")
		message.WriteString(fmt.Sprintf("%d %s", response.EntityCount.CompCount, componentWord))
	}
	if response.EntityCount.RelCount > 0 && response.EntityCount.CompCount > 0 {
		if message.Len() > 0 {
			message.WriteString(" and ")
		}
		relationshipWord := determinePluralWord(response.EntityCount.RelCount, "relationship")
		message.WriteString(fmt.Sprintf("%d %s", response.EntityCount.RelCount, relationshipWord))
	}
	if response.EntityCount.CompCount == 0 && response.EntityCount.RelCount > 0 {
		relationshipWord := determinePluralWord(response.EntityCount.RelCount, "relationship")
		message.WriteString(fmt.Sprintf("%d %s", response.EntityCount.RelCount, relationshipWord))
	}
	if response.EntityCount.CompCount > 0 || response.EntityCount.RelCount > 0 {
		message.WriteString(")")
	}
	return message.String()
}

func ErrMsgContruct(response *models.RegistryAPIResponse) string {
	entityWord := determinePluralWord(response.EntityCount.TotalErrCount, "entity")

	msg := fmt.Sprintf("encountered error for %d %s (", response.EntityCount.TotalErrCount, entityWord)
	componentsPresent := response.EntityCount.ErrCompCount > 0
	relationshipsPresent := response.EntityCount.ErrRelCount > 0
	modelsPresent := response.EntityCount.ErrModelCount > 0
	knownErrors := response.EntityCount.ErrCompCount + response.EntityCount.ErrRelCount + response.EntityCount.ErrModelCount
	unknownErrors := response.EntityCount.TotalErrCount - knownErrors

	// Collect errors in a slice for dynamic message construction
	errors := []string{}
	if modelsPresent {
		modelWord := determinePluralWord(response.EntityCount.ErrModelCount, "model")
		errors = append(errors, fmt.Sprintf("%d %s", response.EntityCount.ErrModelCount, modelWord))
	}
	if componentsPresent {
		componentWord := determinePluralWord(response.EntityCount.ErrCompCount, "component")
		errors = append(errors, fmt.Sprintf("%d %s", response.EntityCount.ErrCompCount, componentWord))
	}
	if relationshipsPresent {
		relationshipWord := determinePluralWord(response.EntityCount.ErrRelCount, "relationship")
		errors = append(errors, fmt.Sprintf("%d %s", response.EntityCount.ErrRelCount, relationshipWord))
	}
	if unknownErrors > 0 {
		unknownEntityWord := determinePluralWord(unknownErrors, "unknown entity")
		errors = append(errors, fmt.Sprintf("%d %s", unknownErrors, unknownEntityWord))
	}
	if len(errors) > 1 {
		msg += fmt.Sprintf("%s and %s", strings.Join(errors[:len(errors)-1], ", "), errors[len(errors)-1])
	} else if len(errors) == 1 {
		msg += errors[0]
	}
	msg += ")"
	return msg
}
func determinePluralWord(count int, word string) string {
	if count > 1 {
		if strings.HasSuffix(word, "y") {
			return word[:len(word)-1] + "ies"
		} else if strings.HasSuffix(word, "s") {
			return word + "es"
		}
		return word + "s"
	}
	return word
}
func setDefaultValues(model *mesheryctlUtils.ModelCSV) {
	svgColor := `<svg xmlns="http://www.w3.org/2000/svg" id="Layer_1" data-name="Layer 1" viewBox="0 0 134.95 135.02"><defs><style>.cls-1{fill:#00d3a9}.cls-2{fill:#00b39f}</style></defs><title>meshery-logo-light</title><polygon points="69.49 31.82 69.49 64.07 97.44 47.89 69.49 31.82" class="cls-1"/><polygon points="69.49 70.81 69.49 103.22 97.7 87.09 69.49 70.81" class="cls-1"/><polygon points="65.47 63.85 65.47 32.09 37.87 47.92 65.47 63.85" class="cls-2"/><path d="M10.1,103.1a67.79,67.79,0,0,0,21.41,21.55V90.71Z" class="cls-2"/><polygon points="65.47 103.06 65.47 71.05 37.8 87.07 65.47 103.06" class="cls-2"/><polygon points="35.54 122.63 63.56 106.61 35.54 90.41 35.54 122.63" class="cls-1"/><polygon points="99.61 122.8 99.61 90.63 71.63 106.63 99.61 122.8" class="cls-2"/><path d="M127,99.37a67.22,67.22,0,0,0,7.91-28.94L105.78,87.11Z" class="cls-2"/><polygon points="103.64 83.69 131.76 67.61 103.64 51.45 103.64 83.69" class="cls-1"/><polygon points="99.61 44.5 99.61 12.52 71.76 28.49 99.61 44.5" class="cls-2"/><polygon points="99.61 83.55 99.61 51.28 71.7 67.44 99.61 83.55" class="cls-2"/><polygon points="67.48 135.02 67.49 135.02 67.48 135.02 67.48 135.02" class="cls-2"/><polygon points="35.54 51.22 35.54 83.73 63.66 67.45 35.54 51.22" class="cls-1"/><path d="M65.47,0A67.2,67.2,0,0,0,35.83,7.83l29.64,17Z" class="cls-2"/><polygon points="35.54 12.3 35.54 44.62 63.68 28.48 35.54 12.3" class="cls-1"/><path d="M31.51,10.34A67.89,67.89,0,0,0,10.1,31.89L31.51,44.25Z" class="cls-2"/><path d="M99.43,8A67.23,67.23,0,0,0,69.49,0V25.15Z" class="cls-1"/><path d="M0,69.87A67.27,67.27,0,0,0,8.07,99.63L29.76,87.07Z" class="cls-1"/><path d="M8.07,35.37A67.16,67.16,0,0,0,0,65L29.79,47.91Z" class="cls-1"/><path d="M35.78,127.13A67.13,67.13,0,0,0,65.47,135V110.15Z" class="cls-2"/><path d="M124.92,32a67.9,67.9,0,0,0-21.28-21.52V44.3Z" class="cls-1"/><path d="M103.64,124.54A68,68,0,0,0,125,102.86L103.64,90.52Z" class="cls-1"/><path d="M135,64.81a67.06,67.06,0,0,0-8-29.35L105.49,47.88Z" class="cls-2"/><path d="M69.49,135a67.12,67.12,0,0,0,29.63-7.83L69.49,110Z" class="cls-1"/><polygon points="31.51 83.44 31.51 51.56 3.83 67.43 31.51 83.44" class="cls-2"/></svg>`
	svgWhite := `<svg width="32" height="32" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M16.405 8.732v6.57l5.694-3.297-5.694-3.273Zm0 7.942v6.602l5.747-3.285-5.747-3.317Z" fill="#fff"/><path d="M15.586 15.256v-6.47l-5.622 3.225 5.622 3.245ZM4.307 23.252a13.809 13.809 0 0 0 4.362 4.39v-6.914l-4.362 2.524Zm11.279-.008v-6.52L9.95 19.985l5.636 3.258Z" fill="#fff" fill-opacity=".8"/><path d="m9.49 27.23 5.707-3.263-5.707-3.3v6.563Z" fill="#fff"/><path d="M22.54 27.265v-6.553l-5.699 3.259 5.7 3.294Zm5.58-4.773a13.697 13.697 0 0 0 1.612-5.895l-5.934 3.397 4.323 2.498Z" fill="#fff" fill-opacity=".8"/><path d="m23.362 19.298 5.728-3.276-5.728-3.291v6.567Z" fill="#fff"/><path d="M22.541 11.315V4.8l-5.673 3.253 5.673 3.262Zm0 7.955v-6.574l-5.685 3.292 5.685 3.281Z" fill="#fff" fill-opacity=".8"/><path d="M9.49 12.684v6.622l5.728-3.316-5.728-3.306Z" fill="#fff"/><path d="M15.586 2.25a13.69 13.69 0 0 0-6.037 1.595l6.037 3.463V2.25Z" fill="#fff" fill-opacity=".8"/><path d="M9.49 4.756v6.583l5.732-3.288L9.49 4.756Z" fill="#fff"/><path d="M8.669 4.356a13.83 13.83 0 0 0-4.362 4.39l4.362 2.518V4.356Z" fill="#fff" fill-opacity=".8"/><path d="M22.504 3.88a13.695 13.695 0 0 0-6.099-1.63v5.123l6.1-3.493ZM2.25 16.483c.071 2.12.634 4.196 1.644 6.062l4.418-2.559-6.062-3.503Zm1.644-7.028a13.68 13.68 0 0 0-1.644 6.036l6.068-3.482-4.424-2.554Z" fill="#fff"/><path d="M9.539 28.147a13.673 13.673 0 0 0 6.047 1.603v-5.062L9.54 28.147Z" fill="#fff" fill-opacity=".8"/><path d="M27.697 8.768a13.83 13.83 0 0 0-4.335-4.383v6.889l4.335-2.506ZM23.362 27.62a13.851 13.851 0 0 0 4.351-4.417l-4.351-2.514v6.93Z" fill="#fff"/><path d="M29.75 15.452a13.659 13.659 0 0 0-1.63-5.979l-4.381 2.53 6.011 3.45Z" fill="#fff" fill-opacity=".8"/><path d="M16.405 29.75a13.673 13.673 0 0 0 6.036-1.595l-6.036-3.498v5.093Z" fill="#fff"/><path d="M8.669 19.247v-6.494L3.03 15.986l5.639 3.261Z" fill="#fff" fill-opacity=".8"/></svg>`
	setIfEmpty(&model.PrimaryColor, "#00B39F")
	setIfEmpty(&model.SecondaryColor, "#00D3A9")
	setIfEmpty(&model.SVGColor, svgColor)
	setIfEmpty(&model.SVGWhite, svgWhite)
	setIfEmpty(&model.Category, "Uncategorized")
	setIfEmpty(&model.Shape, "Circle")
	setIfEmpty(&model.SubCategory, "Uncategorized")
}

func setIfEmpty(field *string, defaultValue string) {
	if *field == "" {
		*field = defaultValue
	}
}
