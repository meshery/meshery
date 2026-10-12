package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/schemas/models/core"
)

type resultErrorProvider struct {
	models.Provider
	err error
}

func (p *resultErrorProvider) GetResult(_ string, _ core.Uuid) (*models.MesheryResult, error) {
	return nil, p.err
}

func TestGetResultHandlerErrorStatus(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{
			name:       "missing result",
			err:        models.ErrResultNotFound(errors.New("record not found")),
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "other persistence error",
			err:        errors.New("database unavailable"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := newTestHandler(t, nil, "")
			provider := &resultErrorProvider{err: test.err}
			req := httptest.NewRequest(http.MethodGet, "/api/perf/profile/result/5f92b476-bbb5-44f5-a99e-88bca9a40613", nil)
			req = mux.SetURLVars(req, map[string]string{"id": "5f92b476-bbb5-44f5-a99e-88bca9a40613"})
			req = req.WithContext(context.WithValue(req.Context(), models.TokenCtxKey, "test-token"))
			rec := httptest.NewRecorder()

			h.GetResultHandler(rec, req, nil, nil, provider)

			if rec.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d (body=%q)", rec.Code, test.wantStatus, rec.Body.String())
			}
		})
	}
}
