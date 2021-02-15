package oam

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/layer5io/meshery/models/oam/core/v1alpha1"
	"github.com/sirupsen/logrus"
	cytoscapejs "gonum.org/v1/gonum/graph/formats/cytoscapejs"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Pattern is the golang representation of the Pattern
// config file model
type Pattern struct {
	Name     string              `yaml:"name,omitempty"`
	Services map[string]*Service `yaml:"services,omitempty"`
}

// Service represents the services defined within the appfile
type Service struct {
	Type      string   `yaml:"type,omitempty"`
	Namespace string   `yaml:"namespace,omitempty"`
	DependsOn []string `yaml:"dependsOn,omitempty"`

	Settings map[string]interface{} `yaml:"settings,omitempty"`
	Traits   map[string]interface{} `yaml:"traits,omitempty"`
}

// NewPatternFile takes in raw yaml and encodes it into a construct
func NewPatternFile(yml []byte) (af Pattern, err error) {
	err = yaml.Unmarshal(yml, &af)

	for _, svc := range af.Services {
		svc.Settings = RecursiveCastMapStringInterfaceToMapStringInterface(svc.Settings)
		svc.Traits = RecursiveCastMapStringInterfaceToMapStringInterface(svc.Traits)

		if svc.Settings == nil {
			svc.Settings = map[string]interface{}{}
		}
		if svc.Traits == nil {
			svc.Traits = map[string]interface{}{}
		}
	}

	return
}

// GetApplicationComponent generates OAM Application Components from the
// the given Pattern file
func (p *Pattern) GetApplicationComponent(name string) (v1alpha1.Component, error) {
	svc, ok := p.Services[name]
	if !ok {
		return v1alpha1.Component{}, fmt.Errorf("invalid service name")
	}

	comp := v1alpha1.Component{
		TypeMeta:   v1.TypeMeta{Kind: "Component", APIVersion: "core.oam.dev/v1alpha2"},
		ObjectMeta: v1.ObjectMeta{Name: name, Namespace: svc.Namespace},
		Spec: v1alpha1.ComponentSpec{
			Type:     svc.Type,
			Settings: svc.Settings,
		},
	}

	return comp, nil
}

// GenerateApplicationConfiguration generates OAM Application Configuration from the
// the given Pattern file for a particular deploymnet
func (p *Pattern) GenerateApplicationConfiguration() (v1alpha1.Configuration, error) {
	config := v1alpha1.Configuration{
		TypeMeta:   v1.TypeMeta{Kind: "ApplicationConfiguration", APIVersion: "core.oam.dev/v1alpha2"},
		ObjectMeta: v1.ObjectMeta{Name: p.Name},
	}

	// Create configs for each component
	for k, v := range p.Services {
		// Indicates that map for properties is not empty
		if len(v.Traits) > 0 {
			specComp := v1alpha1.ConfigurationSpecComponent{
				ComponentName: k,
			}

			for k2, v2 := range v.Traits {
				castToMap, ok := v2.(map[string]interface{})

				trait := v1alpha1.ConfigurationSpecComponentTrait{
					Name: k2,
				}

				if !ok {
					castToMap = map[string]interface{}{}
				}

				trait.Properties = castToMap

				specComp.Traits = append(specComp.Traits, trait)
			}

			config.Spec.Components = append(config.Spec.Components, specComp)
		}
	}

	return config, nil
}

// GetServiceType returns the type of the service
func (p *Pattern) GetServiceType(name string) string {
	return p.Services[name].Type
}

// ToCytoscapeJS converts pattern file into cytoscape object
func (p *Pattern) ToCytoscapeJS() (cytoscapejs.GraphElem, error) {
	var cy cytoscapejs.GraphElem

	// Not specifying any cytoscapejs layout
	// should fallback to "default" layout

	// Not specifying styles, may get applied on the
	// client side

	// Set up the nodes
	for name, svc := range p.Services {
		// Skip if type is either prometheus or grafana
		if !notIn(svc.Type, []string{"prometheus", "grafana"}) {
			continue
		}

		elemData := cytoscapejs.ElemData{
			ID: name, // Assuming that the service names are unique
		}

		elemPosition := getCytoscapeJSPosition(svc)

		elem := cytoscapejs.Element{
			Data:       elemData,
			Position:   &elemPosition,
			Selectable: true,
			Grabbable:  true,
			Scratch: map[string]Service{
				"_data": *svc,
			},
		}

		cy.Elements = append(cy.Elements, elem)
	}

	return cy, nil
}

// ToYAML converts a patternfile to yaml
func (p *Pattern) ToYAML() ([]byte, error) {
	return yaml.Marshal(p)
}

// NewPatternFileFromCytoscapeJSJSON takes in CytoscapeJS JSON
// and creates a PatternFile from it
func NewPatternFileFromCytoscapeJSJSON(byt []byte) (Pattern, error) {
	// Unmarshal data into cytoscape struct
	var cy cytoscapejs.GraphElem
	if err := json.Unmarshal(byt, &cy); err != nil {
		return Pattern{}, err
	}

	// Convert cytoscape struct to patternfile
	pf := Pattern{
		Name:     "MesheryGeneratedPatternFile",
		Services: make(map[string]*Service),
	}
	for _, elem := range cy.Elements {
		// Try to create Service object from the elem.scratch's _data field
		// if this fails then immediately fail the process and return an error
		castedScratch, ok := elem.Scratch.(map[string]interface{})
		if !ok {
			return pf, fmt.Errorf("empty scratch field is not allowed, must containe \"_data\" field holding metadata")
		}

		data, ok := castedScratch["_data"]
		if !ok {
			return pf, fmt.Errorf("\"_data\" cannot be empty")
		}

		// Convert data to JSON for easy serialization
		svcByt, err := json.Marshal(&data)
		if err != nil {
			return pf, fmt.Errorf("failed to serialize service from the metadata in the scratch")
		}

		// Unmarshal the JSON into a service
		var svc Service
		if err := json.Unmarshal(svcByt, &svc); err != nil {
			return pf, fmt.Errorf("failed to create service from the metadata in the scratch")
		}

		// Add other meshmap specific data into service
		svc.Traits["meshmap"] = map[string]map[string]float64{
			"position": {
				"posX": elem.Position.X,
				"posY": elem.Position.Y,
			},
		}

		pf.Services[elem.Data.ID] = &svc
	}

	return pf, nil
}

func getCytoscapeJSPosition(svc *Service) (pos cytoscapejs.Position) {
	// Check if the service has "meshmap" as a trait
	mpi, ok := svc.Traits["meshmap"]
	if !ok {
		rand.Seed(time.Now().UnixNano())
		pos.X = float64(rand.Intn(100))
		pos.Y = float64(rand.Intn(100))

		return
	}

	mpStrInterface, ok := mpi.(map[string]interface{})
	if !ok {
		logrus.Debugf("failed to cast meshmap trait (MPI): %+#v", mpi)
		return
	}

	posInterface, ok := mpStrInterface["position"]
	if !ok {
		logrus.Debugf("failed to cast meshmap trait (posInterface): %+#v", mpStrInterface)
		return
	}

	posMap, ok := posInterface.(map[string]interface{})
	if !ok {
		logrus.Debugf("failed to cast meshmap trait (posMap): %+#v", posInterface)
		return
	}

	pos.X, ok = posMap["posX"].(float64)
	if !ok {
		logrus.Debugf("failed to cast meshmap trait (posMap): %T\n", posMap["posX"])

		// Attempt to cast as int
		intX, ok := posMap["posX"].(int)
		if !ok {
			logrus.Debugf("failed to cast meshmap trait (posMap): %T\n", posMap["posX"])
		}

		pos.X = float64(intX)
	}
	pos.Y, ok = posMap["posY"].(float64)
	if !ok {
		logrus.Debugf("failed to cast meshmap trait (posMap): %T\n", posMap["posY"])

		// Attempt to cast as int
		intY, ok := posMap["posY"].(int)
		if !ok {
			logrus.Debugf("failed to cast meshmap trait (posMap): %T\n", posMap["posY"])
		}

		pos.Y = float64(intY)
	}

	return
}

func notIn(name string, prohibited []string) bool {
	for _, p := range prohibited {
		if strings.HasPrefix(strings.ToLower(name), p) {
			return false
		}
	}

	return true
}
