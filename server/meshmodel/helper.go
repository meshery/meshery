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

type meshmodelHelper struct {
	handlerConfig    *models.HandlerConfig
	regManager       *meshmodel.RegistryManager
	componentChan    chan v1alpha1.ComponentDefinition
	relationshipChan chan v1alpha1.RelationshipDefinition
	doneSignal       chan bool
	errorChan        chan error
}

// seed the local meshmodel components
func (mh meshmodelHelper) SeedComponents() {
	// Read component and relationship definitions from files and send them to channels
	go func() {
		mh.generateComponents("/components")
		mh.generateRelationships("/relationships")
	}()

	// Register components and relationships with the registry manager
	go mh.watchComponents()
}

func (mh meshmodelHelper) generateComponents(pathToComponents string) {
	path, err := filepath.Abs(pathToComponents)
	if err != nil {
		mh.errorChan <- errors.Wrapf(err, "error while getting absolute path for generating components")
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
				mh.componentChan <- comp
			}
		}
		return nil
	})
	mh.errorChan <- errors.Wrapf(err, "error while generating components")
	return
}

func (mh meshmodelHelper) generateRelationships(pathToComponents string) {
	path, err := filepath.Abs(pathToComponents)
	if err != nil {
		mh.errorChan <- errors.Wrapf(err, "error while getting absolute path for generating relationships")
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
			mh.relationshipChan <- rel
		}
		return nil
	})

	mh.errorChan <- errors.Wrapf(err, "error while generating relationships")
	mh.doneSignal <- true
	return
}

func (mh meshmodelHelper) watchComponents() {
	var err error
	for {
		select {
		case comp := <-mh.componentChan:
			err = mh.regManager.RegisterEntity(meshmodel.Host{
				Hostname: ArtifactHubComponentsHandler,
			}, comp)
		case rel := <-mh.relationshipChan:
			err = mh.regManager.RegisterEntity(meshmodel.Host{
				Hostname: ArtifactHubComponentsHandler,
			}, rel)
		case mhErr := <-mh.errorChan:
			fmt.Println("err: ", mhErr.Error())

		case <-mh.doneSignal:
			go mh.handlerConfig.MeshModelSummaryChannel.Publish()
			return
		}

		if err != nil {
			mh.errorChan <- err
		}
	}
}
