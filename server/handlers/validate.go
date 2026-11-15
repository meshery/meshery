package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	cueerrors "cuelang.org/go/cue/errors"
	"github.com/meshery/meshkit/utils"
)

type validationInputType string

const (
	jsontype       validationInputType = "JSON"
	yamltype       validationInputType = "YAML"
	jsonschematype validationInputType = "JSONSCHEMA"
	cuetype        validationInputType = "CUE"
)

type validationItem struct {
	Schema    string              `json:"schema"`
	Value     string              `json:"value"`
	ValueType validationInputType `json:"valueType,omitempty"`
}

type payload struct {
	ValidationItems map[string]validationItem `json:"validationItems"`
}

type validationResponse struct {
	IsValid bool     `json:"isValid"`
	Errors  []string `json:"errors"`
}

type jsonSchemaValidationType struct {
	Schema string `json:"$schema,omitempty"`
}

// maxValidationPayloadSize limits the maximum allowed request body size (4 MB)
// for validation requests, preventing unbounded memory allocation and denial-of-service
// vulnerabilities via oversized payloads.
const maxValidationPayloadSize = 4 << 20 // 4MB

// ValidationHandler expects the request body to be JSON.
// request body should be of format - {validationItems: {[id]:{schema: string, value: string, valueType: "JSON"|"YAML"|"CUE"}}}
// response format - {[id]: {isValid: bool, error: string}}
func (h *Handler) ValidationHandler(rw http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	body, err := io.ReadAll(http.MaxBytesReader(rw, r.Body, maxValidationPayloadSize))
	if err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			h.log.Error(ErrRequestBody(err))
			writeMeshkitError(rw, ErrRequestBody(err), http.StatusRequestEntityTooLarge)
			return
		}
		h.log.Error(ErrRequestBody(err))
		writeMeshkitError(rw, ErrRequestBody(err), http.StatusBadRequest)
		return
	}
	body = bytes.Trim(body, " \t\r\n")
	if len(body) == 0 || bytes.Equal(body, []byte("null")) {
		err := fmt.Errorf("validation request body is empty or null")
		h.log.Error(ErrRequestBody(err))
		writeMeshkitError(rw, ErrRequestBody(err), http.StatusBadRequest)
		return
	}
	// Unmarshal request body
	pld := payload{}
	err = json.Unmarshal(body, &pld)
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		writeMeshkitError(rw, ErrRequestBody(err), http.StatusBadRequest)
		return
	}
	if len(pld.ValidationItems) == 0 {
		err := fmt.Errorf("no validation items provided")
		h.log.Error(ErrRequestBody(err))
		writeMeshkitError(rw, ErrRequestBody(err), http.StatusBadRequest)
		return
	}
	// Validate
	validationResults := validate(pld.ValidationItems)
	// Send response
	rw.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(rw).Encode(struct {
		Result map[string]validationResponse `json:"result"`
	}{
		Result: validationResults,
	})
	if err != nil {
		// Encode streams directly to rw; a failure here means the body has
		// already started — log only, don't double-write.
		if isClientDisconnect(err) {
			h.log.Debug(ErrValidate(err))
		} else {
			h.log.Error(ErrValidate(err))
		}
	}
}

func validate(validationItems map[string]validationItem) map[string]validationResponse {
	validationResults := make(map[string]validationResponse, 0)
	for id, vi := range validationItems {
		// Parse the schema as CUE value
		schemaType := findSchemaType(vi.Schema)
		cueSchema, err := parseSchema(vi.Schema, schemaType)
		if err != nil {
			// if there is an error, push it into the map and continue
			validationResults[id] = validationResponse{IsValid: false, Errors: []string{err.Error()}}
			continue
		}
		// Parse the value as CUE value
		cueValue, err := parseValue(vi.Value, vi.ValueType)
		if err != nil {
			// if there is an error, push it into the map and continue
			validationResults[id] = validationResponse{IsValid: false, Errors: []string{err.Error()}}
			continue
		}
		// Validate the value against the schema
		isValid, errs := utils.Validate(cueSchema, cueValue)
		if len(errs) != 0 {
			errs := convertCueErrorsToStrings(errs)
			validationResults[id] = validationResponse{IsValid: false, Errors: errs}
		}
		if isValid {
			// empty string means that the item is valid
			validationResults[id] = validationResponse{IsValid: true, Errors: make([]string, 0)}
		}
	}
	return validationResults
}

func convertCueErrorsToStrings(errs []cueerrors.Error) []string {
	var res []string
	// sanitize
	for _, err := range errs {
		_ = cueerrors.Sanitize(err)
	}
	for _, err2 := range errs {
		res = append(res, err2.Error())
	}
	return res
}

// if schema is not a JSONSCHEMA, we assume that it is CUE
func findSchemaType(schema string) validationInputType {
	jsValType := jsonSchemaValidationType{}
	err := json.Unmarshal([]byte(schema), &jsValType)
	if err != nil {
		// schema is not a valid JSON
		return cuetype
	}
	if jsValType.Schema == "" {
		return cuetype
	}
	return jsonschematype
}

// NOTE: does not return meshkit error - make sure to wrap it in meshkit errors before using
func parseSchema(schema string, schemaType validationInputType) (cue.Value, error) {
	if schemaType == jsonschematype {
		cueVal, err := utils.JsonSchemaToCue(schema)
		if err != nil {
			return cue.Value{}, err
		}
		return cueVal, nil
	}
	if schemaType == cuetype {
		// if not jsonschema, then it should be CUE
		cuectx := cuecontext.New()
		cueVal := cuectx.CompileString(schema)
		if cueVal.Err() != nil {
			return cue.Value{}, cueVal.Err()
		}
		return cueVal, nil
	}
	return cue.Value{}, fmt.Errorf("given schema is not in JSONSCHEMA or CUE format")
}

// NOTE: does not return meshkit error - make sure to wrap it in meshkit errors before using
func parseValue(value string, valueType validationInputType) (cue.Value, error) {
	if valueType == jsontype {
		cueVal, err := utils.JsonToCue([]byte(value))
		if err != nil {
			return cue.Value{}, err
		}
		return cueVal, nil
	}
	if valueType == yamltype {
		cueVal, err := utils.YamlToCue(value)
		if err != nil {
			return cue.Value{}, err
		}
		return cueVal, nil
	}
	if valueType == cuetype {
		cuectx := cuecontext.New()
		cueVal := cuectx.CompileString(value)
		if cueVal.Err() != nil {
			return cue.Value{}, cueVal.Err()
		}
		return cueVal, nil
	}
	return cue.Value{}, fmt.Errorf("given value is not in JSON,YAML or CUE format")
}
