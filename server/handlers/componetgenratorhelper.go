package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/models"

	"github.com/layer5io/meshkit/encoding"
	"github.com/layer5io/meshkit/models/events"
	"github.com/layer5io/meshkit/models/meshmodel/entity"
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
				"Model":            r.Model.DisplayName,
				"Kind":             r.Kind,
				"Subtype":          r.SubType,
				"Selectors":        r.Selectors,
				"RelationshipType": r.RelationshipType, //future when we support type
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
	} else if response.EntityCount.ModelCount > 0 {
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
