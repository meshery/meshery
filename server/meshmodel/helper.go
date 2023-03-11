package meshmodel

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/layer5io/meshery/server/helpers/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/models/meshmodel"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
	"github.com/pkg/errors"
)

const ArtifactHubComponentsHandler = "kubernetes" //The components generated in output directory will be handled by kubernetes

type ComponentHelper struct {
	handlerConfig    *models.HandlerConfig
	regManager       *meshmodel.RegistryManager
	componentChan    chan v1alpha1.ComponentDefinition
	relationshipChan chan v1alpha1.RelationshipDefinition
	doneSignal       chan bool
	errorChan        chan error
}

func NewComponentHelper(hc *models.HandlerConfig, rm *meshmodel.RegistryManager) *ComponentHelper {
	return &ComponentHelper{
		handlerConfig:    hc,
		regManager:       rm,
		componentChan:    make(chan v1alpha1.ComponentDefinition, 1),
		relationshipChan: make(chan v1alpha1.RelationshipDefinition, 1),
		doneSignal:       make(chan bool),
		errorChan:        make(chan error),
	}
}

// seed the local meshmodel components
func (ch *ComponentHelper) SeedComponents() {
	// Watch channels and register components and relationships with the registry manager
	go ch.watchComponents()

	// Read component and relationship definitions from files and send them to respective channels
	ch.generateComponents("/components")
	ch.generateRelationships("/relationships")
	ch.doneSignal <- true
}

func (ch *ComponentHelper) generateComponents(pathToComponents string) {
	path, err := filepath.Abs(pathToComponents)
	if err != nil {
		ch.errorChan <- errors.Wrapf(err, "error while getting absolute path for generating components")
		return
	}

	err = filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if info == nil {
			return nil
		}

		if !info.IsDir() {
			var comp v1alpha1.ComponentDefinition
			byt, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			err = json.Unmarshal(byt, &comp)
			if err != nil {
				return err
			}
			if comp.Metadata != nil && comp.Metadata["published"] == true {
				utils.WriteSVGsOnFileSystem(&comp)
				ch.componentChan <- comp
			}
		}
		return nil
	})
	ch.errorChan <- errors.Wrapf(err, "error while generating components")
	return
}

func (ch *ComponentHelper) generateRelationships(pathToComponents string) {
	path, err := filepath.Abs(pathToComponents)
	if err != nil {
		ch.errorChan <- errors.Wrapf(err, "error while getting absolute path for generating relationships")
		return
	}

	err = filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		if !info.IsDir() {
			var rel v1alpha1.RelationshipDefinition
			byt, err := os.ReadFile(path)
			if err != nil {
				return nil
			}
			err = json.Unmarshal(byt, &rel)
			if err != nil {
				return nil
			}
			ch.relationshipChan <- rel
		}
		return nil
	})

	ch.errorChan <- errors.Wrapf(err, "error while generating relationships")
	return
}

func (ch *ComponentHelper) watchComponents() {
	var err error
	for {
		select {
		case comp := <-ch.componentChan:
			err = ch.regManager.RegisterEntity(meshmodel.Host{
				Hostname: ArtifactHubComponentsHandler,
			}, comp)
		case rel := <-ch.relationshipChan:
			err = ch.regManager.RegisterEntity(meshmodel.Host{
				Hostname: ArtifactHubComponentsHandler,
			}, rel)
		case mhErr := <-ch.errorChan:
			fmt.Println("err: ", mhErr.Error())

		case <-ch.doneSignal:
			go ch.handlerConfig.MeshModelSummaryChannel.Publish()
			return
		}

		if err != nil {
			ch.errorChan <- err
		}
	}
}
