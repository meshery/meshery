package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/models"
	pCore "github.com/layer5io/meshery/server/models/pattern/core"
	"github.com/layer5io/meshkit/models/events"
)

func (h *Handler) HandleCytoPattern(
	res http.ResponseWriter,
	userID uuid.UUID,
	provider models.Provider,
	eb *events.EventBuilder,
	body *MesheryPatternPOSTRequestBody,
	format,
	token,
	action string,
) {
	sourcetype := models.Design

	pf, err := pCore.NewPatternFileFromCytoscapeJSJSON(body.Name, []byte(body.CytoscapeJSON))
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(res, "%s", err)
		event := eb.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": ErrSavePattern(err),
		}).WithDescription("Pattern save failed, cytoJSON could be malformed.").Build()

		_ = provider.PersistEvent(event)
		go h.config.EventBroadcaster.Publish(userID, event)
		return
	}

	pfByt, err := pf.ToYAML()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(res, "%s", err)
		event := eb.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": ErrSavePattern(err),
		}).WithDescription(ErrSavePattern(err).Error()).Build()

		_ = provider.PersistEvent(event)
		go h.config.EventBroadcaster.Publish(userID, event)
		return
	}

	patternName := pf.Name
	fmt.Println("TEST >>>>>>>>>>")

	mesheryPattern := &models.MesheryPattern{
		Name:        patternName,
		PatternFile: string(pfByt),
		Location: map[string]interface{}{
			"host": "",
			"path": "",
			"type": "local",
		},
		CatalogData: body.PatternData.CatalogData,
		Type: sql.NullString{
			String: string(models.Design),
			Valid:  true,
		},
	}
	if body.PatternData != nil {
		mesheryPattern.ID = body.PatternData.ID
	}

	var patternFileInBytes []byte
	if body.Save {
		patternFileInBytes, err = h.SavePattern(
			res,
			provider,
			eb,
			mesheryPattern,
			userID,
			token,
		)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		patternFileInBytes, err = json.Marshal([]models.MesheryPattern{*mesheryPattern})
		if err != nil {
			http.Error(res, ErrEncodePattern(err).Error(), http.StatusInternalServerError)
			return
		}
	}

	mesheryPatterns, patternNames, err := h.convertPatternToCyto(provider, userID, patternFileInBytes, eb)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(&mesheryPatterns)
	if err != nil {
		obj := "pattern file"
		http.Error(res, models.ErrMarshal(err, obj).Error(), http.StatusInternalServerError)
		return
	}

	var actionDesc string
	if action == models.Update {
		actionDesc = "updated"
	} else {
		actionDesc = "created"
	}
	response := fmt.Sprintf("%s \"%s\" %s.", sourcetype, strings.Join(patternNames, ","), actionDesc)
	event := eb.WithDescription(response).Build()
	_ = provider.PersistEvent(event)

	res.Header().Set("Content-Type", "application/json")
	res.Write(data) // or 	fmt.Fprint(rw, string(data))
}

func (h *Handler) convertPatternToCyto(provider models.Provider, userID uuid.UUID, resp []byte, eb *events.EventBuilder) ([]models.MesheryPattern, []string, error) {
	contentMesheryPatternSlice := make([]models.MesheryPattern, 0)
	if err := json.Unmarshal(resp, &contentMesheryPatternSlice); err != nil {
		return nil, nil, ErrDecodePattern(err)
	}
	result := []models.MesheryPattern{}
	names := []string{}
	for _, content := range contentMesheryPatternSlice {
		if content.ID != nil {
			eb.ActedUpon(*content.ID)
		}

		patternFile, err := pCore.NewPatternFile([]byte(content.PatternFile))
		if err != nil {
			err = ErrParsePattern(err)
			event := eb.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
				"error": err,
			}).WithDescription("Unable to parse pattern file, pattern could be malformed.").Build()
			h.log.Error(err)
			_ = provider.PersistEvent(event)
			go h.config.EventBroadcaster.Publish(userID, event)

			return nil, nil, err
		}
		newpatternfile := evalImportAndReferenceStage(&patternFile)
		cyjs, _ := newpatternfile.ToCytoscapeJS()
		bytes, err := json.Marshal(&cyjs)
		if err != nil {
			err = ErrConvertPattern(err)
			event := eb.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
				"error": err,
			}).WithDescription("Unable to convert pattern file to cytoscape format, pattern file could be malformed.").Build()
			h.log.Error(err)
			_ = provider.PersistEvent(event)
			go h.config.EventBroadcaster.Publish(userID, event)
			return nil, nil, err
		}
		// Replace the patternfile with cytoscape type data
		content.PatternFile = string(bytes)
		result = append(result, content)
		names = append(names, content.Name)
	}
	return result, names, nil
}
