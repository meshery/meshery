// Copyright Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
)

func newComponentTestHandler(t *testing.T) *Handler {
	t.Helper()
	return &Handler{
		config: &models.HandlerConfig{},
		log:    newTestLogger(t),
	}
}

func TestDeleteModelsByRegistrant_MissingConnectionID(t *testing.T) {
	h := newComponentTestHandler(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/registry/registrants//models", nil)
	req = mux.SetURLVars(req, map[string]string{"connectionId": ""})
	rec := httptest.NewRecorder()

	h.DeleteModelsByRegistrant(rec, req, nil, nil, nil)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400 for missing connection ID; body: %s", rec.Code, rec.Body.String())
	}
}

func TestDeleteModelsByRegistrant_InvalidConnectionUUID(t *testing.T) {
	h := newComponentTestHandler(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/registry/registrants/invalid-uuid/models", nil)
	req = mux.SetURLVars(req, map[string]string{"connectionId": "invalid-uuid"})
	rec := httptest.NewRecorder()

	h.DeleteModelsByRegistrant(rec, req, nil, nil, nil)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400 for invalid connection UUID; body: %s", rec.Code, rec.Body.String())
	}
}
