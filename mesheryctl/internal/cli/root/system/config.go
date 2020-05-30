// Copyright 2020 The Meshery Authors
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

package system

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

const URL1 = "http://localhost:9081/api/k8sconfig/contexts"
const URL2 = "http://localhost:9081/api/k8sconfig"

const tokenName = "token"
const providerName = "meshery-provider"
const paramName = "k8sfile"
const contextName = "contextName"

func UploadFileWithParams(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, fi.Name())
	if err != nil {
		return nil, err
	}
	part.Write(fileContents)

	if params != nil {
		for key, val := range params {
			_ = writer.WriteField(key, val)
		}
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", writer.FormDataContentType())
	return request, nil
}

func getContext(configFile, tokenPath string) (string, error) {
	client := &http.Client{}

	req, err := UploadFileWithParams(URL1, nil, paramName, configFile)
	if err != nil {
		return "", err
	}

	err = utils.AddAuthDetails(req, tokenPath)
	if err != nil {
		return "", err
	}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var val []map[string]string
	json.Unmarshal(body, &val)
	if len(val) < 1 {
		return "", nil
	}

	return val[0]["contextName"], nil
}

// TODO: show used the available context and let user choose which context to set
func SetDefaultContext(configFile, cname, tokenPath string) error {
	client := &http.Client{}
	extraParams1 := map[string]string{
		"contextName": cname,
	}
	req, err := UploadFileWithParams(URL2, extraParams1, paramName, configFile)
	if err != nil {
		return err
	}
	err = utils.AddAuthDetails(req, tokenPath)
	if err != nil {
		return err
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", string(body))
	return nil
}

// resetCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure Meshery",
	Long:  `Configure the Kubernetes cluster used by Meshery.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		path := "/home/kanishkarj/Downloads/config_minikube.yaml"
		tokenPath := "/home/kanishkarj/Downloads/auth.json"

		cName, err := getContext(path, tokenPath)
		if err != nil {
			log.Printf("Error getting contexts : %s", err.Error())
		}
		err = SetDefaultContext(path, cName, tokenPath)
		if err != nil {
			log.Printf("Error setting context : %s", err.Error())
		}
	},
}
