package core

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/models/pattern/utils"
	"github.com/layer5io/meshkit/models/oam/core/v1alpha1"
	"github.com/sirupsen/logrus"
	cytoscapejs "gonum.org/v1/gonum/graph/formats/cytoscapejs"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Pattern is the golang representation of the Pattern
// config file model
type Pattern struct {
	// Name is the human-readable, display-friendly descriptor of the pattern
	Name string `yaml:"name,omitempty" json:"name,omitempty"`
	//Vars will be used to configure the pattern when it is imported from other patterns.
	Vars map[string]interface{} `yaml:"vars,omitempty" json:"vars,omitempty"`
	// PatternID is the moniker use to uniquely identify any given pattern
	// Convention: SMP-###-v#.#.#
	PatternID string              `yaml:"patternID,omitempty" json:"patternID,omitempty"`
	Services  map[string]*Service `yaml:"services,omitempty" json:"services,omitempty"`
}

// Service represents the services defined within the appfile
type Service struct {
	// ID is the id of the service and is completely internal to
	// Meshery Server and meshery providers
	ID *uuid.UUID `yaml:"id,omitempty" json:"id,omitempty"`
	// Name is the name of the service and is an optional parameter
	// If given then this supercedes the name of the service inherited
	// from the parent
	Name        string            `yaml:"name,omitempty" json:"name,omitempty"`
	Type        string            `yaml:"type,omitempty" json:"type,omitempty"`
	Namespace   string            `yaml:"namespace" json:"namespace"`
	Version     string            `yaml:"version,omitempty" json:"version,omitempty"`
	Labels      map[string]string `yaml:"labels,omitempty" json:"labels,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty" json:"annotations,omitempty"`
	// DependsOn correlates one or more objects as a required dependency of this service
	// DependsOn is used to determine sequence of operations
	DependsOn []string `yaml:"dependsOn,omitempty" json:"dependsOn,omitempty"`

	Settings map[string]interface{} `yaml:"settings,omitempty" json:"settings,omitempty"`
	Traits   map[string]interface{} `yaml:"traits,omitempty" json:"traits,omitempty"`
}

// NewPatternFile takes in raw yaml and encodes it into a construct
func NewPatternFile(yml []byte) (af Pattern, err error) {
	err = yaml.Unmarshal(yml, &af)
	if err != nil {
		return af, err
	}
	for svcName, svc := range af.Services {
		// If an explicit name is not given to the service then use
		// the service identifier as its name
		if svc.Name == "" {
			svc.Name = svcName
		}

		svc.Settings = utils.RecursiveCastMapStringInterfaceToMapStringInterface(svc.Settings)
		svc.Traits = utils.RecursiveCastMapStringInterfaceToMapStringInterface(svc.Traits)

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
		TypeMeta: v1.TypeMeta{Kind: "Component", APIVersion: "core.oam.dev/v1alpha2"},
		ObjectMeta: v1.ObjectMeta{
			Name:        svc.Name,
			Namespace:   svc.Namespace,
			Labels:      svc.Labels,
			Annotations: svc.Annotations,
		},
		Spec: v1alpha1.ComponentSpec{
			Type:     svc.Type,
			Version:  svc.Version,
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
		elemData := cytoscapejs.ElemData{
			ID: getCytoscapeElementID(name, svc),
		}

		elemPosition, err := getCytoscapeJSPosition(svc)
		if err != nil {
			return cy, err
		}

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
func NewPatternFileFromCytoscapeJSJSON(name string, byt []byte) (Pattern, error) {
	// Unmarshal data into cytoscape struct
	var cy cytoscapejs.GraphElem
	if err := json.Unmarshal(byt, &cy); err != nil {
		return Pattern{}, err
	}
	if name == "" {
		name = "MesheryGeneratedPattern"
	}
	// Convert cytoscape struct to patternfile
	pf := Pattern{
		Name:     name,
		Services: make(map[string]*Service),
	}
	dependsOnMap := make(map[string][]string, 0) //used to figure out dependencies from traits.meshmap.parent
	eleToSvc := make(map[string]string)          //used to map cyto element ID uniquely to the name of the service created.
	countDuplicates := make(map[string]int)
	//store the names of services and their count
	err := processCytoElementsWithPattern(cy.Elements, &pf, func(svc Service, ele cytoscapejs.Element) error {
		name, ok := svc.Settings["name"].(string)
		if !ok {
			return fmt.Errorf("missing name in service settings")
		}
		countDuplicates[name]++
		return nil
	})
	if err != nil {
		return pf, err
	}

	//Populate the dependsOn field with appropriate unique service names
	err = processCytoElementsWithPattern(cy.Elements, &pf, func(svc Service, ele cytoscapejs.Element) error {
		//Extract parents, if present
		m, ok := svc.Traits["meshmap"].(map[string]interface{})
		if ok {
			parentID, ok := m["parent"].(string)
			if ok { //If it does not have a parent then we can skip and we dont make it depend on anything
				elementID, ok := m["id"].(string)
				if !ok {
					return fmt.Errorf("required meshmap trait field: \"id\" missing")
				}
				dependsOnMap[elementID] = append(dependsOnMap[elementID], parentID)
			}
		}
		svc.Name, ok = svc.Settings["name"].(string)
		if !ok {
			return fmt.Errorf("required service setting: \"name\" missing")
		}
		//Only make the name unique when duplicates are encountered. This allows clients to preserve and propagate the unique name they want to give to their workload
		if countDuplicates[svc.Name] > 1 {
			//set appropriate unique service name
			svc.Name = strings.ToLower(svc.Name)
			svc.Name += "-" + getRandomAlphabetsOfDigit(5)
		}
		eleToSvc[ele.Data.ID] = svc.Name //will be used while adding depends-on
		pf.Services[svc.Name] = &svc
		return nil
	})
	if err == nil {
		//add depends-on field
		for child, parents := range dependsOnMap {
			childSvc := eleToSvc[child]
			if childSvc != "" {
				for _, parent := range parents {
					if eleToSvc[parent] != "" {
						pf.Services[childSvc].DependsOn = append(pf.Services[childSvc].DependsOn, eleToSvc[parent])
					}
				}
			}
		}
	}
	return pf, err
}

func getRandomAlphabetsOfDigit(length int) (s string) {
	charSet := "abcdedfghijklmnopqrstuvwxyz"
	for i := 0; i < length; i++ {
		random := mathrand.Intn(len(charSet))
		randomChar := charSet[random]
		s += string(randomChar)
	}
	return
}

// processCytoElementsWithPattern iterates over all the cyto elements, convert each into a patternfile service and exposes a callback to handle that service
func processCytoElementsWithPattern(eles []cytoscapejs.Element, pf *Pattern, callback func(svc Service, ele cytoscapejs.Element) error) error {
	for _, elem := range eles {
		// Try to create Service object from the elem.scratch's _data field
		// if this fails then immediately fail the process and return an error
		castedScratch, ok := elem.Scratch.(map[string]interface{})
		if !ok {
			return fmt.Errorf("empty scratch field is not allowed, must contain \"_data\" field holding metadata")
		}

		data, ok := castedScratch["_data"]
		if !ok {
			return fmt.Errorf("\"_data\" cannot be empty")
		}

		// Convert data to JSON for easy serialization
		svcByt, err := json.Marshal(&data)
		if err != nil {
			return fmt.Errorf("failed to serialize service from the metadata in the scratch")
		}

		// Unmarshal the JSON into a service
		svc := Service{
			Settings: map[string]interface{}{},
			Traits:   map[string]interface{}{},
		}

		// Add meshmap position
		svc.Traits["meshmap"] = map[string]interface{}{
			"position": map[string]float64{
				"posX": elem.Position.X,
				"posY": elem.Position.Y,
			},
		}

		if err := json.Unmarshal(svcByt, &svc); err != nil {
			return fmt.Errorf("failed to create service from the metadata in the scratch")
		}
		err = callback(svc, elem)
		if err != nil {
			return err
		}
	}
	return nil
}

func manifestIsEmpty(manifests []string) bool {
	for _, m := range manifests {
		x := strings.TrimSpace(strings.Trim(m, "\n"))
		if x != "---" && x != "" {
			return false
		}
	}
	return true
}

func NewPatternFileFromK8sManifest(data string, ignoreErrors bool) (Pattern, error) {
	pattern := Pattern{
		Name:     "Autogenerated",
		Services: map[string]*Service{},
	}

	manifests := strings.Split(data, "\n---\n")
	//For `---` separated manifests, even if only one manifest is there followed/preceded by multiple `\n---\n`- the manifest be will be valid
	//If there is no data present (except \n---\n) , then the yaml will be marked as empty and error will be thrown
	if manifestIsEmpty(manifests) {
		return pattern, ErrParseK8sManifest(fmt.Errorf("manifest is empty"))
	}
	for _, manifestYAML := range manifests {
		manifest := map[string]interface{}{}

		if err := yaml.Unmarshal([]byte(manifestYAML), &manifest); err != nil {
			if ignoreErrors {
				continue
			}
			return pattern, ErrParseK8sManifest(err)
		}
		if len(manifest) == 0 {
			continue
		}
		// Recursive casting
		manifest = utils.RecursiveCastMapStringInterfaceToMapStringInterface(manifest)
		if manifest == nil {
			if ignoreErrors {
				continue
			}
			return pattern, ErrParseK8sManifest(fmt.Errorf("failed to parse manifest into an internal representation"))
		}

		name, svc, err := createPatternServiceFromK8s(manifest)
		if err != nil {
			if ignoreErrors {
				continue
			}

			return pattern, ErrCreatePatternService(fmt.Errorf("failed to create pattern service from extended kubernetes component: %s", err))
		}

		pattern.Services[name] = &svc
	}

	return pattern, nil
}

func createPatternServiceFromK8s(manifest map[string]interface{}) (string, Service, error) {
	apiVersion, _ := manifest["apiVersion"].(string)
	kind, _ := manifest["kind"].(string)
	metadata, _ := manifest["metadata"].(map[string]interface{})
	name, _ := metadata["name"].(string)
	namespace, _ := metadata["namespace"].(string)
	labels, _ := metadata["labels"].(map[string]interface{})
	annotations, _ := metadata["annotations"].(map[string]interface{})
	if namespace == "" {
		namespace = "default"
	}
	fmt.Printf("%+#v\n", manifest)
	// rest will store a map of everything other than the above mentioned fields
	rest := map[string]interface{}{}
	for k, v := range manifest {
		// Ignore a few fields
		if k == "apiVersion" || k == "kind" || k == "metadata" || k == "status" {
			continue
		}

		rest[k] = v
	}

	id := name
	uid, err := uuid.NewV4()
	if err == nil {
		id = uid.String()
	}
	if apiVersion == "" || kind == "" {
		return "", Service{}, ErrCreatePatternService(fmt.Errorf("empty apiVersion or kind in manifest"))
	}
	w := GetWorkloadsByK8sAPIVersionKind(apiVersion, kind)

	if len(w) == 0 {
		return "", Service{}, ErrCreatePatternService(fmt.Errorf("no resources found for APIVersion: %s Kind: %s", apiVersion, kind))
	}

	// Setup labels
	castedLabel := map[string]string{}
	for k, v := range labels {
		cv, ok := v.(string)
		if ok {
			castedLabel[k] = cv
		}
	}

	// Setup annotations
	castedAnnotation := map[string]string{}
	for k, v := range annotations {
		cv, ok := v.(string)
		if ok {
			castedAnnotation[k] = cv
		}
	}

	svc := Service{
		Name:        name,
		Type:        w[0].OAMDefinition.Name,
		Namespace:   namespace,
		Labels:      castedLabel,
		Annotations: castedAnnotation,
		Settings:    rest,
	}

	return id, svc, nil
}

// getCytoscapeElementID returns the element id for a given service
func getCytoscapeElementID(name string, svc *Service) string {
	mpi, ok := svc.Traits["meshmap"] // check if service has meshmap as trait
	if !ok {
		return name // Assuming that the service names are unique
	}

	mpStrInterface, ok := mpi.(map[string]interface{})
	if !ok {
		logrus.Debugf("failed to cast meshmap trait (MPI): %+#v", mpi)
		return name // Assuming that the service names are unique
	}

	mpID, ok := mpStrInterface["id"].(string)
	if !ok {
		logrus.Debugf("Meshmap id not present in Meshmap interface")
		return name // Assuming that the service names are unique
	}

	return mpID
}

func getCytoscapeJSPosition(svc *Service) (cytoscapejs.Position, error) {
	pos := cytoscapejs.Position{}

	// Check if the service has "meshmap" as a trait
	mpi, ok := svc.Traits["meshmap"]

	if !ok {
		randX, err := rand.Int(rand.Reader, big.NewInt(100))
		if err != nil {
			return pos, err
		}
		randY, err := rand.Int(rand.Reader, big.NewInt(100))
		if err != nil {
			return pos, err
		}

		pos := cytoscapejs.Position{}
		pos.X, _ = big.NewFloat(0).SetInt(randX).Float64()
		pos.Y, _ = big.NewFloat(0).SetInt(randY).Float64()

		return pos, nil
	}

	mpStrInterface, ok := mpi.(map[string]interface{})
	if !ok {
		logrus.Debugf("failed to cast meshmap trait (MPI): %+#v", mpi)
		return pos, nil
	}

	posInterface, ok := mpStrInterface["position"]
	if !ok {
		logrus.Debugf("failed to cast meshmap trait (posInterface): %+#v", mpStrInterface)
		return pos, nil
	}

	posMap, ok := posInterface.(map[string]interface{})
	if !ok {
		logrus.Debugf("failed to cast meshmap trait (posMap): %+#v", posInterface)
		return pos, nil
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

	return pos, nil
}

func init() {
	mathrand.Seed(time.Now().Unix())
}
