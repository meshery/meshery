package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/meshery/meshery/server/models"
)

var testSchema2 = `{
  "$id": "https://example.com/person.schema.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "Person",
  "type": "object",
  "properties": {
    "firstName": {
      "type": "string",
      "description": "The person's first name."
    },
    "lastName": {
      "type": "string",
      "description": "The person's last name."
    },
    "age": {
      "description": "Age in years which must be equal to or greater than zero.",
      "type": "integer",
      "minimum": 0
    }
  }
}`
var testSchema3 = `{
  "$id": "https://example.com/person.schema.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "Person",
  "type": "object",
  "properties": {
    "firstName": {
      "type": "string",
      "description": "The person's first name."
    },
    "lastName": {
      "type": "string",
      "description": "The person's last name."
    },
    "age": {
      "description": "Age in years which must be equal to or greater than zero.",
      "type": "integer",
      "minimum": 0
    }
  }, 
  "required": ["age", "lastName"]
}`
var testSchema4 = `{
  "$id": "https://example.com/person.schema.json",
  "title": "Person",
  "type": "object",
  "properties": {
    "firstName": {
      "type": "string",
      "description": "The person's first name."
    },
    "lastName": {
      "type": "string",
      "description": "The person's last name."
    },
    "age": {
      "description": "Age in years which must be equal to or greater than zero.",
      "type": "integer",
      "minimum": 0
    }
  }
}`
var testSchema5 = `
name: string
task: string
properties?: {
  item: "array"
}
`

var value1 = `
firstName: "Meshery"
lastName: "Project"
`
var value2 = `
firstName: "Meshery"
lastName: "Project"
age: "string"
`
var value3 = `
name: "meshery"
task: "nothing"
`
var value4 = `
firstName: "Meshery"
lastName: "Project"
age: 34
`
var value5 = `
firstName: 34
lastName: 54
`

func TestFindSchemaType(t *testing.T) {
	var tests = []struct {
		schema string
		want   validationInputType
	}{
		{testSchema2, jsonschematype},
		{testSchema4, cuetype},
		{testSchema5, cuetype},
	}

	for _, tt := range tests {
		t.Run("schema type", func(t *testing.T) {
			ans := findSchemaType(tt.schema)
			if ans != tt.want {
				t.Errorf("got %v, want %v", ans, tt.want)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	var tests = []struct {
		validationItems map[string]validationItem
		want            map[string]bool
	}{
		{
			map[string]validationItem{
				"1": {Schema: testSchema2, Value: value1, ValueType: cuetype},
				"2": {Schema: testSchema2, Value: value2, ValueType: cuetype},
				"3": {Schema: testSchema3, Value: value1, ValueType: cuetype},
				"4": {Schema: testSchema3, Value: value4, ValueType: cuetype},
				"5": {Schema: testSchema5, Value: value3, ValueType: cuetype},
				"6": {Schema: testSchema3, Value: value5, ValueType: cuetype},
			},
			map[string]bool{
				"1": true,
				"2": false,
				"3": false,
				"4": true,
				"5": true,
				"6": false,
			},
		},
	}
	for _, tt := range tests {
		t.Run("meshmodel validate", func(t *testing.T) {
			ans := validate(tt.validationItems)
			for id, vr := range ans {
				if vr.IsValid != tt.want[id] {
					t.Errorf("got %v, want %v for id: %v", ans, tt.want, id)
				}
			}
		})
	}
}

func TestValidationHandler_ValidPayload(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")

	reqBody := payload{
		ValidationItems: map[string]validationItem{
			"item1": {
				Schema:    testSchema2,
				Value:     value1,
				ValueType: cuetype,
			},
		},
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/registry/validate", bytes.NewReader(bodyBytes))
	rec := httptest.NewRecorder()

	h.ValidationHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d (body: %s)", rec.Code, rec.Body.String())
	}

	var resp struct {
		Result map[string]validationResponse `json:"result"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	item1, exists := resp.Result["item1"]
	if !exists {
		t.Fatalf("expected item1 in result, got: %v", resp.Result)
	}
	if !item1.IsValid {
		t.Errorf("expected item1 to be valid, got errors: %v", item1.Errors)
	}
}

func TestValidationHandler_EmptyBody(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")

	req := httptest.NewRequest(http.MethodPost, "/api/registry/validate", bytes.NewReader([]byte("   ")))
	rec := httptest.NewRecorder()

	h.ValidationHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 Bad Request, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), ErrRequestBodyCode) {
		t.Errorf("expected response to contain %s, got: %s", ErrRequestBodyCode, rec.Body.String())
	}
}

func TestValidationHandler_NullBody(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")

	req := httptest.NewRequest(http.MethodPost, "/api/registry/validate", bytes.NewReader([]byte("null")))
	rec := httptest.NewRecorder()

	h.ValidationHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 Bad Request, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), ErrRequestBodyCode) {
		t.Errorf("expected response to contain %s, got: %s", ErrRequestBodyCode, rec.Body.String())
	}
}

func TestValidationHandler_MalformedJSON(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")

	req := httptest.NewRequest(http.MethodPost, "/api/registry/validate", bytes.NewReader([]byte("{invalid-json")))
	rec := httptest.NewRecorder()

	h.ValidationHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 Bad Request, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), ErrRequestBodyCode) {
		t.Errorf("expected response to contain %s, got: %s", ErrRequestBodyCode, rec.Body.String())
	}
}

func TestValidationHandler_PayloadTooLarge(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")

	oversizedBody := bytes.Repeat([]byte("a"), maxValidationPayloadSize+10)
	req := httptest.NewRequest(http.MethodPost, "/api/registry/validate", bytes.NewReader(oversizedBody))
	rec := httptest.NewRecorder()

	h.ValidationHandler(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status 413 Request Entity Too Large, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), ErrRequestBodyCode) {
		t.Errorf("expected response to contain %s, got: %s", ErrRequestBodyCode, rec.Body.String())
	}
}

func TestValidationHandler_ExactBoundaryPayload(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")

	// Construct a valid payload and pad with trailing whitespace to exactly maxValidationPayloadSize bytes.
	raw := `{"validationItems":{"item1":{"schema":"name: string","value":"name: \"meshery\"","valueType":"CUE"}}}`
	if len(raw) > maxValidationPayloadSize {
		t.Fatalf("raw payload exceeds limit")
	}
	padding := bytes.Repeat([]byte(" "), maxValidationPayloadSize-len(raw))
	exactBoundaryBody := append([]byte(raw), padding...)

	if len(exactBoundaryBody) != maxValidationPayloadSize {
		t.Fatalf("expected exact boundary size %d, got %d", maxValidationPayloadSize, len(exactBoundaryBody))
	}

	req := httptest.NewRequest(http.MethodPost, "/api/registry/validate", bytes.NewReader(exactBoundaryBody))
	rec := httptest.NewRecorder()

	h.ValidationHandler(rec, req)

	// MaxBytesReader allows exactly maxValidationPayloadSize bytes, so it must NOT return 413
	if rec.Code == http.StatusRequestEntityTooLarge {
		t.Fatalf("expected payload at exactly maxValidationPayloadSize to be accepted, but got 413")
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK for exact boundary payload, got %d (body: %s)", rec.Code, rec.Body.String())
	}
}

func TestValidationHandler_ExceedBoundaryByOneByte(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")

	// Body of exactly maxValidationPayloadSize + 1 byte must be rejected with 413
	oversizedBody := bytes.Repeat([]byte("a"), maxValidationPayloadSize+1)
	req := httptest.NewRequest(http.MethodPost, "/api/registry/validate", bytes.NewReader(oversizedBody))
	rec := httptest.NewRecorder()

	h.ValidationHandler(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status 413 Request Entity Too Large at limit+1, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), ErrRequestBodyCode) {
		t.Errorf("expected response to contain %s, got: %s", ErrRequestBodyCode, rec.Body.String())
	}
}

func TestValidationHandler_EmptyValidationItems(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")

	// Empty object {}
	req := httptest.NewRequest(http.MethodPost, "/api/registry/validate", bytes.NewReader([]byte("{}")))
	rec := httptest.NewRecorder()

	h.ValidationHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 Bad Request for empty object, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), ErrRequestBodyCode) {
		t.Errorf("expected response to contain %s, got: %s", ErrRequestBodyCode, rec.Body.String())
	}

	// Explicit empty validationItems: {}
	reqEmptyItems := httptest.NewRequest(http.MethodPost, "/api/registry/validate", bytes.NewReader([]byte(`{"validationItems":{}}`)))
	recEmptyItems := httptest.NewRecorder()

	h.ValidationHandler(recEmptyItems, reqEmptyItems)

	if recEmptyItems.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 Bad Request for empty validationItems, got %d (body: %s)", recEmptyItems.Code, recEmptyItems.Body.String())
	}
	if !strings.Contains(recEmptyItems.Body.String(), ErrRequestBodyCode) {
		t.Errorf("expected response to contain %s, got: %s", ErrRequestBodyCode, recEmptyItems.Body.String())
	}
}

func TestValidationHandler_NonJSONWhitespaceRejected(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")

	// 0x0B (vertical tab) is not RFC 8259 JSON whitespace, so it must not be trimmed into valid JSON
	body := []byte("\x0b" + `{"validationItems":{"item1":{"schema":"name: string","value":"name: \"meshery\"","valueType":"CUE"}}}`)
	req := httptest.NewRequest(http.MethodPost, "/api/registry/validate", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.ValidationHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 Bad Request for non-JSON whitespace prefix (0x0B), got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), ErrRequestBodyCode) {
		t.Errorf("expected response to contain %s, got: %s", ErrRequestBodyCode, rec.Body.String())
	}
}
