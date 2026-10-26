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
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/helpers/utils"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/models/events"
)

type mockCSVProvider struct {
	*models.DefaultLocalProvider
}

func (m *mockCSVProvider) GetProviderToken(_ *http.Request) (string, error) {
	return "", nil
}

func (m *mockCSVProvider) PersistEvent(_ events.Event, _ string) error {
	return nil
}

func TestRegisterMeshmodels_TemporaryCSVCleanup(t *testing.T) {
	systemID := uuid.Must(uuid.NewV4())
	h := &Handler{
		config: &models.HandlerConfig{
			EventBroadcaster: models.NewBroadcaster("Events"),
		},
		log:      newTestLogger(t),
		SystemID: &systemID,
	}

	provider := &mockCSVProvider{DefaultLocalProvider: &models.DefaultLocalProvider{}}
	user := &models.User{ID: uuid.Must(uuid.NewV4())}

	tempPattern := filepath.Join(os.TempDir(), "model-*.csv")
	initialMatches, err := filepath.Glob(tempPattern)
	if err != nil {
		t.Fatalf("failed to list temp files: %v", err)
	}
	initialMap := make(map[string]bool, len(initialMatches))
	for _, m := range initialMatches {
		initialMap[m] = true
	}

	payload := map[string]interface{}{
		"uploadType": "csv",
		"importBody": map[string]interface{}{
			"modelCsv":     "data:text/csv;base64," + base64.StdEncoding.EncodeToString([]byte("name,model\ntest,test-model")),
			"componentCsv": "invalid-csv-data",
		},
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/meshmodels/register", bytes.NewReader(payloadBytes))
	rec := httptest.NewRecorder()

	h.RegisterMeshmodels(rec, req, nil, user, provider)

	afterMatches, err := filepath.Glob(tempPattern)
	if err != nil {
		t.Fatalf("failed to list temp files after handler run: %v", err)
	}

	for _, match := range afterMatches {
		if !initialMap[match] {
			t.Errorf("temporary file %s was leaked by RegisterMeshmodels", match)
		}
	}
}

func TestRegisterMeshmodels_CrossPlatformRegistryLocation(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skipf("skipping: user home directory not available: %v", err)
	}

	modelLocation := filepath.Join(home, utils.RegistryLocation)
	if modelLocation == "" {
		t.Fatalf("expected non-empty model location")
	}

	invalidRootPath := filepath.Join("", utils.RegistryLocation)
	if modelLocation == invalidRootPath {
		t.Errorf("model location %q incorrectly resolved to system root without user home", modelLocation)
	}

	if !strings.HasPrefix(modelLocation, home) {
		t.Errorf("expected model location %q to start with home directory %q", modelLocation, home)
	}
}
