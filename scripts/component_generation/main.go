package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"time"

	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
	"github.com/layer5io/meshkit/utils/artifacthub"
	"gopkg.in/yaml.v3"
)

var (
	AhSearchEndpoint = artifacthub.AhHelmExporterEndpoint

	OutputDirectoryPath     = "./output"
	ComponentsFileName      = "components.yaml"
	ComponentModelsFileName = "component_models.yaml"
)

type ComponentModel struct {
	SystemName string `yaml:"systemName"`
	AhRepo     string `yaml:"ahRepo"`
	Version    string `yaml:"version"`
	RepoUrl    string `yaml:"repoUrl"`
}

func convertPackagesToCompModels(pkgs []artifacthub.AhPackage) []ComponentModel {
	models := make([]ComponentModel, 0)
	for _, ap := range pkgs {
		models = append(models, ComponentModel{
			SystemName: ap.Name,
			Version:    ap.Version,
			AhRepo:     ap.Repository,
			RepoUrl:    ap.RepoUrl,
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
			RepoUrl:    cm.RepoUrl,
		})
	}
	return pkgs
}

func main() {
	compsFd, err := os.OpenFile(filepath.Join(OutputDirectoryPath, ComponentsFileName), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
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
	content, err = io.ReadAll(modelsFd)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = yaml.Unmarshal(content, &models)
	if err != nil {
		fmt.Println(err)
	}
	pkgs = convertCompModelsToPackages(models)
	if pkgs == nil || len(pkgs) == 0 {
		pkgs, err = artifacthub.GetAllAhHelmPackages()
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
	for i := 0; i <= 10; i++ {
		StartPipeline(inputChan, &compsWriter)
	}
	inputChan <- pkgs[:10]
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
				fmt.Println("[DEBUG] Updating package data for: ", ap.Name)
				err := ap.UpdatePackageData()
				if err != nil {
					fmt.Println(err)
					continue
				}
				ahPkgs = append(ahPkgs, ap)
			}
			pkgsChan <- ahPkgs
		}
	}()
	// generation of components
	go func() {
		for pkgs := range pkgsChan {
			compStructs := make([]ComponentStruct, 0)
			for _, ap := range pkgs {
				fmt.Println("[DEBUG] Generating components for: ", ap.Name)
				comps, err := ap.GenerateComponents()
				if err != nil {
					fmt.Println(err)
					continue
				}
				compStructs = append(compStructs, wrapComponentsInWritableStruct(comps, ap.Name))
			}
			compsChan <- compStructs
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
	PackageName string               `yaml:"name"`
	Components  []v1alpha1.Component `yaml:"components"`
}

func wrapComponentsInWritableStruct(comps []v1alpha1.Component, pkgName string) ComponentStruct {
	return ComponentStruct{
		PackageName: pkgName,
		Components:  comps,
	}

}

type Writer struct {
	file *os.File
	m    sync.Mutex
}

func writeComponentModels(models []ComponentModel, writer *Writer) error {
	writer.m.Lock()
	defer writer.m.Unlock()
	val, err := yaml.Marshal(models)
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
	writer.m.Lock()
	defer writer.m.Unlock()
	compsToWrite := make([]ComponentStruct, 0)
	content, err := io.ReadAll(writer.file)
	if err != nil {
		return err
	}
	// 1. the file is empty
	if string(content) == "" {
		for _, cs := range cmps {
			if len(cs.Components) != 0 {
				compsToWrite = append(compsToWrite, cs)
			}
		}
	}
	// 2. the file already has some components in it
	if string(content) != "" {
		err = yaml.Unmarshal(content, &compsToWrite)
		if err != nil {
			return err
		}
		for _, cs := range cmps {
			if len(cs.Components) != 0 {
				fmt.Println("[DEBUG] Writing components for package: ", cs.PackageName)
				compsToWrite = append(compsToWrite, cs)
			}
		}
	}
	if len(compsToWrite) == 0 {
		return nil
	}
	// write components
	val, err := yaml.Marshal(compsToWrite)
	if err != nil {
		return err
	}
	_, err = writer.file.Write(val)
	if err != nil {
		return err
	}
	return nil
}
