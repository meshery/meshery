package context

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var b *bytes.Buffer

type CmdTestInput struct {
	Name             string
	Args             []string
	ExpectedResponse string
}

func SetupContextEnv(t *testing.T, testfile string) {
	path, err := os.Getwd()
	if err != nil {
		t.Error("unable to locate meshery directory")
	}
	viper.Reset()
	viper.SetConfigFile(path + testfile)
	err = viper.ReadInConfig()
	if err != nil {
		t.Errorf("unable to read configuration from %v, %v", viper.ConfigFileUsed(), err.Error())
	}

	configuration, err = config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		t.Error("error processing config", err)
	}
}
func SetupFunc(t *testing.T) {
	b = bytes.NewBufferString("")
	logrus.SetOutput(b)
	utils.SetupLogrusFormatter()
	ContextCmd.SetOut(b)
}
func TestViewContextCmd(t *testing.T) {
	SetupContextEnv(t, "/fixtures/.meshery/TestContext.yaml")

	tests := []CmdTestInput{
		{
			Name:             "view for default context",
			Args:             []string{"view"},
			ExpectedResponse: viewExpected,
		},
		{
			Name:             "view with specified context through context flag",
			Args:             []string{"view", "--context", "local2"},
			ExpectedResponse: viewWithContextExpected,
		},
		{
			Name: "Error for viewing a non-existing context",
			Args: []string{"view", "local3"},
			ExpectedResponse: `context "local3" doesn't exists, run the following to create:

mesheryctl system context create local3
`,
		},
		{
			Name:             "view with specified context as argument",
			Args:             []string{"view", "local2"},
			ExpectedResponse: viewWithContextExpected,
		},
		{
			Name:             "view with all flag set",
			Args:             []string{"view", "--all"},
			ExpectedResponse: viewAllExpected,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			SetupFunc(t)
			ContextCmd.SetArgs(tt.Args)
			err := ContextCmd.Execute()
			if err != nil {
				t.Error(err)
			}

			actualResponse := b.String()
			expectedResponse := tt.ExpectedResponse
			if expectedResponse != actualResponse {
				t.Errorf("expected response [%v] and actual response [%v] don't match", expectedResponse, actualResponse)
			}
		})
	}
}
func TestListContextCmd(t *testing.T) {
	SetupContextEnv(t, "/fixtures/.meshery/TestContext.yaml")

	tests := []CmdTestInput{
		{
			Name:             "list all contexts",
			Args:             []string{"list"},
			ExpectedResponse: listExpected,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			SetupFunc(t)
			ContextCmd.SetArgs(tt.Args)
			err := ContextCmd.Execute()
			if err != nil {
				t.Error(err)
			}

			actualResponse := b.String()
			expectedResponse := tt.ExpectedResponse

			if expectedResponse != actualResponse {
				t.Errorf("expected response %v and actual response %v don't match", expectedResponse, actualResponse)
			}
		})
	}
}

func TestDeleteContextCmd(t *testing.T) {
	SetupContextEnv(t, "/testdata/context/ExpectedDelete.yaml")

	tests := []CmdTestInput{
		{
			Name:             "delete given context",
			Args:             []string{"delete", "local2"},
			ExpectedResponse: "deleted context local2\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			SetupFunc(t)
			ContextCmd.SetArgs(tt.Args)
			err := ContextCmd.Execute()
			if err != nil {
				t.Error(err)
			}

			actualResponse := b.String()
			expectedResponse := tt.ExpectedResponse

			if expectedResponse != actualResponse {
				t.Errorf("expected response [%v] and actual response [%v] don't match", expectedResponse, actualResponse)
			}
			path, err := os.Getwd()
			if err != nil {
				t.Error("unable to locate meshery directory")
			}
			filepath := path + "/testdata/context/ExpectedDelete.yaml"

			content, err := ioutil.ReadFile(filepath)
			if err != nil {
				t.Error(err)
			}
			actualResponse = string(content)
			if actualResponse != deleteExpected {
				t.Errorf("expected response [%v] and actual response [%v] don't match", deleteExpected, actualResponse)
			}

			//Repopulating Expected yaml
			if err := copy(path+"/fixtures/.meshery/TestContext.yaml", filepath); err != nil {
				t.Error(err, "Could not complete test. Unable to configure delete test file")
			}
		})
	}
}
func TestAddContextCmd(t *testing.T) {
	SetupContextEnv(t, "/testdata/context/ExpectedAdd.yaml")

	tests := []CmdTestInput{
		{
			Name:             "add given context",
			Args:             []string{"create", "local3"},
			ExpectedResponse: "Added `local3` context\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			SetupFunc(t)
			ContextCmd.SetArgs(tt.Args)
			err := ContextCmd.Execute()
			if err != nil {
				t.Error(err)
			}

			actualResponse := b.String()
			expectedResponse := tt.ExpectedResponse

			if expectedResponse != actualResponse {
				t.Errorf("expected response [%v] and actual response [%v] don't match", expectedResponse, actualResponse)
			}
			path, err := os.Getwd()
			if err != nil {
				t.Error("unable to locate meshery directory")
			}
			filepath := path + "/testdata/context/ExpectedAdd.yaml"

			content, err := ioutil.ReadFile(filepath)
			if err != nil {
				t.Error(err)
			}
			actualResponse = string(content)
			if actualResponse != addExpected {
				t.Errorf("expected response [%v] and actual response [%v] don't match", addExpected, actualResponse)
			}

			//Repopulating Expected yaml
			if err := copy(path+"/fixtures/.meshery/TestContext.yaml", filepath); err != nil {
				t.Error(err, "Could not complete test. Unable to configure delete test file")
			}
		})
	}
}

func TestSwitchContextCmd(t *testing.T) {
	SetupContextEnv(t, "/testdata/context/ExpectedSwitch.yaml")

	tests := []CmdTestInput{
		{
			Name:             "switch to a different context",
			Args:             []string{"switch", "local2"},
			ExpectedResponse: "switched to context 'local2'\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			SetupFunc(t)
			ContextCmd.SetArgs(tt.Args)
			err := ContextCmd.Execute()
			if err != nil {
				t.Error(err)
			}

			actualResponse := b.String()
			expectedResponse := tt.ExpectedResponse

			if expectedResponse != actualResponse {
				t.Errorf("expected response %v and actual response %v don't match", expectedResponse, actualResponse)
			}
			path, err := os.Getwd()
			if err != nil {
				t.Error("unable to locate meshery directory")
			}
			filepath := path + "/testdata/context/ExpectedSwitch.yaml"
			content, err := ioutil.ReadFile(filepath)
			if err != nil {
				t.Error(err)
			}
			actualResponse = string(content)
			if actualResponse != switchExpected {
				t.Errorf("expected response %v and actual response %v don't match", switchExpected, actualResponse)
			}
			if err := copy(path+"/fixtures/.meshery/TestContext.yaml", filepath); err != nil {
				t.Error(err, "Could not complete test. Unable to configure delete test file")
			}
		})
	}
}

//utility function to repopulate config
func copy(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}

// Expected response for different cases hardcoded below for View command.
var viewExpected = `
Current Context: local

endpoint: http://localhost:9081
token: Default
token-location: auth.json
platform: kubernetes
adapters:
- meshery-istio
channel: stable
version: latest

`
var viewWithContextExpected = `
Current Context: local2

endpoint: http://localhost:32242
token: Default2
token-location: auth.json
platform: docker
adapters:
- meshery-istio
channel: stable
version: latest

`
var viewAllExpected = `local:
  endpoint: http://localhost:9081
  token: Default
  token-location: auth.json
  platform: kubernetes
  adapters:
  - meshery-istio
  channel: stable
  version: latest
local2:
  endpoint: http://localhost:32242
  token: Default2
  token-location: auth.json
  platform: docker
  adapters:
  - meshery-istio
  channel: stable
  version: latest

`

// Expected response for different cases hardcoded below for List command.
var listExpected = `Current context: local

Available contexts:

- local
- local2
`

// Expected response for different cases hardcoded below for Delete command.
var deleteExpected = `contexts:
  local:
    endpoint: http://localhost:9081
    token: Default
    platform: kubernetes
    adapters:
    - meshery-istio
    channel: stable
    version: latest
current-context: local
tokens:
- location: auth.json
  name: Default
- location: auth.json
  name: Default2
`

var addExpected = `contexts:
  local:
    endpoint: http://localhost:9081
    token: Default
    platform: kubernetes
    adapters:
    - meshery-istio
    channel: stable
    version: latest
  local2:
    endpoint: http://localhost:32242
    token: Default2
    platform: docker
    adapters:
    - meshery-istio
    channel: stable
    version: latest
  local3:
    endpoint: http://localhost:9081
    token: Default
    platform: docker
    adapters:
    - meshery-istio
    - meshery-linkerd
    - meshery-consul
    - meshery-nsm
    - meshery-kuma
    - meshery-cpx
    - meshery-osm
    - meshery-traefik-mesh
    channel: stable
    version: latest
current-context: local
tokens:
- name: Default
  location: auth.json
- name: Default2
  location: auth.json
`

var switchExpected = `contexts:
  local:
    adapters:
    - meshery-istio
    channel: stable
    endpoint: http://localhost:9081
    platform: kubernetes
    token: Default
    version: latest
  local2:
    adapters:
    - meshery-istio
    channel: stable
    endpoint: http://localhost:32242
    platform: docker
    token: Default2
    version: latest
current-context: local2
tokens:
- location: auth.json
  name: Default
- location: auth.json
  name: Default2
`
