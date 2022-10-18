package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	// "time"

	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
	"github.com/layer5io/meshkit/utils/artifacthub"
)

const ArtifactHubApiEndpoint = "https://artifacthub.io/api/v1/"

var (
	AhSearchEndpoint = ArtifactHubApiEndpoint + "helm-exporter"

	OutputDirectoryPath     = "./output"
	ComponentsFileName      = "components.json"
	ComponentModelsFileName = "component_models.json"
)

type ComponentModel struct {
	SystemName string `json:"systemName"`
	AhRepo     string `json:"ahRepo"`
	Version    string `json:"version"`
}

func convertPackagesToCompModels(pkgs []artifacthub.AhPackage) []ComponentModel {
	models := make([]ComponentModel, 0)
	for _, ap := range pkgs {
		models = append(models, ComponentModel{
			SystemName: ap.Name,
			Version:    ap.Version,
			AhRepo:     ap.Repository,
		})
	}
	return models
}

func convertCompModelsToPackages(models []ComponentModel) []artifacthub.AhPackage {
	pkgs := make([]artifacthub.AhPackage, 0)
	for _, cm := range models {
		pkgs = append(pkgs, artifacthub.AhPackage{
			Name:       cm.SystemName,
			Version:    cm.Version,
			Repository: cm.AhRepo,
		})
	}
	return pkgs
}

func main() {
	compsFd, err := os.OpenFile(filepath.Join(OutputDirectoryPath, ComponentsFileName), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	compsWriter := Writer{
		file: compsFd,
	}
	modelsFd, err := os.OpenFile(filepath.Join(OutputDirectoryPath, ComponentModelsFileName), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	modelsWriter := Writer{
		file: modelsFd,
	}
	defer modelsFd.Close()
	defer compsFd.Close()
	// move to a new function: getHelmPackages
	content := make([]byte, 0)
	models := make([]ComponentModel, 0)
	pkgs := make([]artifacthub.AhPackage, 0)
	content, err = ioutil.ReadAll(modelsFd)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(content, &models)
	if err != nil {
		fmt.Println(err)
	}
	pkgs = convertCompModelsToPackages(models)
	if pkgs == nil || len(pkgs) == 0 {
		pkgs, err = GetAllAhHelmPackages()
		if err != nil {
			fmt.Println(err)
			return
		}
		err = writeComponentModels(convertPackagesToCompModels(pkgs), &modelsWriter)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	inputChan := make(chan []artifacthub.AhPackage)
	// for i := 0; i <= 10; i++ {
	StartPipeline(inputChan, &compsWriter)
	// }
	for len(pkgs) != 0 {
		if len(pkgs) < 50 {
			inputChan <- pkgs
			time.Sleep(10 * time.Second)
			return
		}
		inputChan <- pkgs[:50]
		pkgs = pkgs[50:]
	}
	time.Sleep(20 * time.Second)
}

func StartPipeline(in chan []artifacthub.AhPackage, writer *Writer) error {
	pkgsChan := make(chan []artifacthub.AhPackage)
	compsChan := make(chan []ComponentStruct)
	// updating pacakge data
	go func() {
		for pkgs := range in {
			ahPkgs := make([]artifacthub.AhPackage, 0)
			for _, ap := range pkgs {
				err := ap.UpdatePackageData()
				if err != nil {
					fmt.Println(err)
				}
				ahPkgs = append(ahPkgs, ap)
			}
			pkgsChan <- ahPkgs
		}
	}()
	// generation of components
	go func() {
		for pkgs := range pkgsChan {
			compsChan <- GenerateComponents(pkgs)
		}
	}()
	// writer
	go func() {
		for comps := range compsChan {
			err := writeComponents(comps, writer)
			if err != nil {
				fmt.Println(err)
			}
		}
	}()
	return nil
}

type ComponentStruct struct {
	PackageName string               `json:"name"`
	Components  []v1alpha1.Component `json:"components"`
}

func GetAllAhHelmPackages() ([]artifacthub.AhPackage, error) {
	pkgs := make([]artifacthub.AhPackage, 0)
	fmt.Println("[DEBUG] Getting packages from ArtifactHub")
	resp, err := http.Get(AhSearchEndpoint)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("status code %d for %s", resp.StatusCode, AhSearchEndpoint)
		return nil, err
	}
	defer resp.Body.Close()
	var res []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	for _, p := range res {
		pkgs = append(pkgs, artifacthub.AhPackage{
			Name:       p["name"].(string),
			Version:    p["version"].(string),
			Repository: p["repository"].(map[string]interface{})["name"].(string),
		})
	}
	return pkgs, nil
}

func GenerateComponents(pkgs []artifacthub.AhPackage) []ComponentStruct {
	result := make([]ComponentStruct, 0)
	for _, ap := range pkgs {
		comps, err := ap.GenerateComponents()
		if err != nil {
			fmt.Println(err)
			continue
		}
		result = append(result, ComponentStruct{PackageName: ap.Name, Components: comps})
	}
	return result
}

type Writer struct {
	file *os.File
	m    sync.Mutex
}

func writeComponentModels(models []ComponentModel, writer *Writer) error {
	writer.m.Lock()
	defer writer.m.Unlock()
	val, err := json.MarshalIndent(models, "", " ")
	if err != nil {
		return err
	}
	_, err = writer.file.Write(val)
	if err != nil {
		return err
	}
	return nil
}

func writeComponents(cmps []ComponentStruct, writer *Writer) error {
	out := make([]ComponentStruct, 0)
	var content []byte
	writer.m.Lock()
	defer writer.m.Unlock()
	l, err := writer.file.Read(content)
	if err == nil && l != 0 {
		err = json.Unmarshal(content, &out)
		if err != nil {
			return err
		}
	}
	for _, c := range cmps {
		if len(c.Components) != 0 {
			cmps = append(cmps, c)
		}
	}
	val, err := json.MarshalIndent(out, "", " ")
	if err != nil {
		return err
	}
	_, err = writer.file.Write(val)
	if err != nil {
		return err
	}
	return nil
}
