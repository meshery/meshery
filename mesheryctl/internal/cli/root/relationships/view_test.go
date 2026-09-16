package relationships

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

// runRelationshipViewTest handles all shared scaffolding for relationship view tests.
func runRelationshipViewTest(t *testing.T, args []string) error {
	t.Helper()
	mesheryctlflags.InitValidators(RelationshipCmd)
	testContext := utils.InitTestEnvironment(t)
	t.Cleanup(func() { utils.StopMockery(t) })
	t.Cleanup(func() { utils.ResetCommandFlags(RelationshipCmd, t) })

	utils.TokenFlag = utils.GetToken(t)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot determine current working directory")
	}
	fixturesDir := filepath.Join(filepath.Dir(filename), "fixtures")

	apiResponse := utils.NewGoldenFile(t, "view.relationship.save.api.response.golden", fixturesDir).Load()
	httpmock.RegisterResponder("GET",
		testContext.BaseURL+"/api/registry/models/kubernetes/relationships",
		httpmock.NewStringResponder(200, apiResponse))

	buf := &bytes.Buffer{}
	RelationshipCmd.SetOut(buf)
	RelationshipCmd.SetErr(buf)
	_ = utils.SetupMeshkitLoggerTesting(t, false)
	RelationshipCmd.SetArgs(args)
	return RelationshipCmd.Execute()
}

func TestRelationshipViewSaveCreatesFileWithExtension(t *testing.T) {
	tmpDir := t.TempDir()
	origMesheryFolder := utils.MesheryFolder
	utils.MesheryFolder = tmpDir
	t.Cleanup(func() { utils.MesheryFolder = origMesheryFolder })

	if err := runRelationshipViewTest(t, []string{"view", "kubernetes", "--output-format", "json", "--save"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedFile := filepath.Join(tmpDir, "relationship_kubernetes_aaaabbbb.json")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		entries, _ := os.ReadDir(tmpDir)
		names := make([]string, 0, len(entries))
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Errorf("--save: expected file %q to exist, got: %v", expectedFile, names)
	}
}