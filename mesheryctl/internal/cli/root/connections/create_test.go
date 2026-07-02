package connections

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestConnectionCreateCmd(t *testing.T) {
	testContext := utils.NewTestHelper(t)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	tempDir := t.TempDir()
	fakeKubeConfigPath := filepath.Join(tempDir, "config")
	fakeMesheryConfigPath := filepath.Join(tempDir, "kubeconfig.yaml")

	fakeKubeConfig := `
apiVersion: v1
clusters:
- cluster:
    server: https://127.0.0.1:6443
  name: kind-kind
contexts:
- context:
    cluster: kind-kind
    user: kind-kind
  name: kind-kind
current-context: kind-kind
kind: Config
preferences: {}
users:
- name: kind-kind
  user:
    token: "fake-token"
`
	if err := os.WriteFile(fakeKubeConfigPath, []byte(fakeKubeConfig), 0644); err != nil {
		t.Fatal(err)
	}

	origKubeConfig := utils.KubeConfig
	origConfigPath := utils.ConfigPath
	defer func() {
		utils.KubeConfig = origKubeConfig
		utils.ConfigPath = origConfigPath
	}()

	utils.KubeConfig = fakeKubeConfigPath
	utils.ConfigPath = fakeMesheryConfigPath

	tests := []utils.MesheryMultiURLCommamdTest{
		{
			Name:          "given no type flag provided when running mesheryctl connection create then an error message is displayed",
			Args:          []string{"create"},
			ExpectError:   true,
			ExpectedError: utils.ErrInvalidArgument(fmt.Errorf("connection type is required. Use --type flag to specify the type (%s)", strings.Join(supportedConnectionTypes, "|"))),
		},
		{
			Name:          "given an invalid type provided when running mesheryctl connection create --type invalid then an error message is displayed",
			Args:          []string{"create", "--type", "invalid"},
			ExpectError:   true,
			ExpectedError: errInvalidConnectionType("invalid"),
		},
		{
			Name: "given a valid type (kind) provided when running mesheryctl connection create then it successfully connects",
			Args: []string{"create", "--type", "kind"},
			URLs: []utils.MockURL{
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/system/kubernetes/contexts",
					Response:     "create.contexts.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/system/kubernetes",
					Response:     "create.kubernetes.response.golden",
					ResponseCode: 200,
				},
			},
			ExpectError:      false,
			ExpectedResponse: "create.connection.kind.output.golden",
			IsOutputGolden:   true,
		},
		{
			Name: "given a valid type (minikube) provided when running mesheryctl connection create then it successfully connects",
			Args: []string{"create", "--type", "minikube"},
			URLs: []utils.MockURL{
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/system/kubernetes/contexts",
					Response:     "create.minikube.contexts.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/system/kubernetes",
					Response:     "create.kubernetes.response.golden",
					ResponseCode: 200,
				},
			},
			ExpectError:      false,
			ExpectedResponse: "create.connection.minikube.output.golden",
			IsOutputGolden:   true,
		},
	}

	resetCreateVariables := func() {
		connectionType = ""
		utils.KubeConfig = fakeKubeConfigPath
		utils.ConfigPath = fakeMesheryConfigPath
	}

	utils.RunMesheryctlMultiURLTests(t, update, ConnectionsCmd, tests, currDir, "connection", resetCreateVariables)
}
