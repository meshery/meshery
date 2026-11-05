package model

import (
	"errors"
	"testing"

	"github.com/jarcoal/httpmock"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models"
	"github.com/stretchr/testify/assert"
)

func TestHasCSVs(t *testing.T) {
	tests := []struct {
		name           string
		dirPath        string
		expectedResult bool
	}{
		{
			name:           "directory with CSVs",
			dirPath:        "./fixtures/with_csvs",
			expectedResult: true,
		},
		{
			name:           "directory without CSVs",
			dirPath:        "./fixtures/without_csvs",
			expectedResult: false,
		},
		{
			name:           "non-existent directory",
			dirPath:        "./fixtures/invalid_path",
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := hasCSVs(tc.dirPath)

			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestImportModelReturnsErrorForURLImportFailure(t *testing.T) {
	defer utils.ResetCommandFlags(ModelCmd, t)
	testContext := utils.InitTestEnvironment(t)
	defer utils.StopMockery(t)
	errMessage := "internal server error"
	utils.TokenFlag = utils.GetToken(t)
	mesheryctlflags.InitValidators(ModelCmd)

	httpmock.RegisterResponder("POST", testContext.BaseURL+"/api/registry/register",
		httpmock.NewStringResponder(500, errMessage))

	buf := utils.SetupMeshkitLoggerTesting(t, false)
	ModelCmd.SetArgs([]string{"import", "https://example.com/model"})
	ModelCmd.SetOut(buf)

	err := ModelCmd.Execute()
	if err == nil {
		t.Fatal("expected an error but command succeeded")
	}

	utils.AssertMeshkitErrorsEqual(t, err, utils.ErrMesheryServerInternalError(errors.New(errMessage)))
}

func TestDisplaySuccessfulRelationships_PanicGuards(t *testing.T) {
	_ = utils.SetupMeshkitLoggerTesting(t, false)

	tests := []struct {
		name     string
		response *models.RegistryAPIResponse
		model    string
	}{
		{
			name: "valid relationship with from and to",
			response: &models.RegistryAPIResponse{
				EntityTypeSummary: models.EntityTypeSummary{
					SuccessfulRelationships: []map[string]interface{}{
						{
							"Kind":             "hierarchical",
							"Subtype":          "inventory",
							"RelationshipType": "parent",
							"Model":            "kubernetes",
							"Selectors": []interface{}{
								map[string]interface{}{
									"allow": map[string]interface{}{
										"from": []interface{}{
											map[string]interface{}{"kind": "Namespace"},
										},
										"to": []interface{}{
											map[string]interface{}{"kind": "Pod"},
										},
									},
								},
							},
						},
					},
				},
			},
			model: "kubernetes",
		},
		{
			name: "empty from slice does not panic",
			response: &models.RegistryAPIResponse{
				EntityTypeSummary: models.EntityTypeSummary{
					SuccessfulRelationships: []map[string]interface{}{
						{
							"Kind":             "hierarchical",
							"Subtype":          "inventory",
							"RelationshipType": "parent",
							"Model":            "kubernetes",
							"Selectors": []interface{}{
								map[string]interface{}{
									"allow": map[string]interface{}{
										"from": []interface{}{},
										"to": []interface{}{
											map[string]interface{}{"kind": "Pod"},
										},
									},
								},
							},
						},
					},
				},
			},
			model: "kubernetes",
		},
		{
			name: "empty to slice does not panic",
			response: &models.RegistryAPIResponse{
				EntityTypeSummary: models.EntityTypeSummary{
					SuccessfulRelationships: []map[string]interface{}{
						{
							"Kind":             "hierarchical",
							"Subtype":          "inventory",
							"RelationshipType": "parent",
							"Model":            "kubernetes",
							"Selectors": []interface{}{
								map[string]interface{}{
									"allow": map[string]interface{}{
										"from": []interface{}{
											map[string]interface{}{"kind": "Namespace"},
										},
										"to": []interface{}{},
									},
								},
							},
						},
					},
				},
			},
			model: "kubernetes",
		},
		{
			name: "missing or empty kind in from and to maps does not panic",
			response: &models.RegistryAPIResponse{
				EntityTypeSummary: models.EntityTypeSummary{
					SuccessfulRelationships: []map[string]interface{}{
						{
							"Kind":             "hierarchical",
							"Subtype":          "inventory",
							"RelationshipType": "parent",
							"Model":            "kubernetes",
							"Selectors": []interface{}{
								map[string]interface{}{
									"allow": map[string]interface{}{
										"from": []interface{}{
											map[string]interface{}{},
										},
										"to": []interface{}{
											map[string]interface{}{"kind": ""},
										},
									},
								},
								map[string]interface{}{
									"allow": map[string]interface{}{
										"from": []interface{}{
											map[string]interface{}{"kind": nil},
										},
										"to": []interface{}{
											map[string]interface{}{"kind": 12345},
										},
									},
								},
							},
						},
					},
				},
			},
			model: "kubernetes",
		},
		{
			name: "malformed and non-slice selectors do not panic",
			response: &models.RegistryAPIResponse{
				EntityTypeSummary: models.EntityTypeSummary{
					SuccessfulRelationships: []map[string]interface{}{
						{
							"Kind":      123,
							"Selectors": "not-a-slice",
							"Model":     "kubernetes",
						},
						{
							"Kind":  "edge",
							"Model": "kubernetes",
							"Selectors": []interface{}{
								"not-a-map",
								map[string]interface{}{
									"allow": "not-a-map",
								},
								map[string]interface{}{
									"allow": map[string]interface{}{
										"from": "not-a-slice",
										"to":   123,
									},
								},
								map[string]interface{}{
									"allow": map[string]interface{}{
										"from": []interface{}{"not-a-map"},
										"to":   []interface{}{123},
									},
								},
							},
						},
					},
				},
			},
			model: "kubernetes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("displaySuccessfulRelationships panicked on %s: %v", tt.name, r)
				}
			}()
			displaySuccessfulRelationships(tt.response, tt.model)
		})
	}
}

func TestBuildEntityTypeLine_PanicGuards(t *testing.T) {
	_ = utils.SetupMeshkitLoggerTesting(t, false)

	tests := []struct {
		name        string
		names       []interface{}
		entityTypes []interface{}
		modelName   string
		expected    string
	}{
		{
			name:        "valid components and relationships",
			names:       []interface{}{"test-model", "test-model"},
			entityTypes: []interface{}{"component", "relationship"},
			modelName:   "test-model",
			expected:    " 1 entity of type component and 1 entity of type relationship",
		},
		{
			name:        "non-string name and entity type do not panic",
			names:       []interface{}{123, nil, "test-model"},
			entityTypes: []interface{}{true, nil, "component"},
			modelName:   "test-model",
			expected:    " 1 entity of type component",
		},
		{
			name:        "empty inputs",
			names:       []interface{}{},
			entityTypes: []interface{}{},
			modelName:   "test-model",
			expected:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("buildEntityTypeLine panicked on %s: %v", tt.name, r)
				}
			}()
			line := buildEntityTypeLine(tt.names, tt.entityTypes, "desc", "cause", "remediation", tt.modelName)
			if tt.expected != "" {
				assert.Equal(t, tt.expected, line)
			}
		})
	}
}
