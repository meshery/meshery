package core

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/models/pattern/utils"
	"github.com/layer5io/meshkit/models/meshmodel"
	meshmodelv1alpha1 "github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
	"github.com/layer5io/meshkit/models/oam/core/v1alpha1"
	"github.com/layer5io/meshkit/utils/manifests"
	"github.com/sirupsen/logrus"
	cytoscapejs "gonum.org/v1/gonum/graph/formats/cytoscapejs"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type prettifier bool

/*
The logic and principle for prettification/deprettification.
1. Specific considerations are made when schema is passed to be prettified like handling kubernetes specific fields.
2. A general rule of thumb is to never prettify or deprettify the end-user fields, the ones which are entered by USER.
For non schema files, it would be all end string fields (user input) and for schema files it would be ENUMS as they are the only system defined fields that are used at end user input.
*/
func (p prettifier) Prettify(m map[string]interface{}, isSchema bool) map[string]interface{} {
	res := ConvertMapInterfaceMapString(m, true, isSchema)
	out, ok := res.(map[string]interface{})
	if !ok {
		fmt.Println("failed to cast")
	}
	return out
}
func (p prettifier) DePrettify(m map[string]interface{}, isSchema bool) map[string]interface{} {
	res := ConvertMapInterfaceMapString(m, false, isSchema)
	out, ok := res.(map[string]interface{})
	if !ok {
		fmt.Println("failed to cast")
	}

	return out
}

// ConvertMapInterfaceMapString converts map[interface{}]interface{} => map[string]interface{}
//
// It will also convert []interface{} => []string
func ConvertMapInterfaceMapString(v interface{}, prettify bool, isSchema bool) interface{} {
	switch x := v.(type) {
	case map[interface{}]interface{}:
		m := map[string]interface{}{}
		for k, v2 := range x {
			switch k2 := k.(type) {
			case string:
				if isSchema && k2 == "enum" { //While schema prettification, ENUMS are end system defined end user input and therefore should not be prettified/deprettified
					m[k2] = v2
					continue
				}
				newmap := ConvertMapInterfaceMapString(v2, prettify, isSchema)
				if isSchema && isSpecialKey(k2) { //Few special keys in schema should not be prettified
					m[k2] = newmap
				} else if prettify {
					m[manifests.FormatToReadableString(k2)] = newmap
				} else {
					m[manifests.DeFormatReadableString(k2)] = newmap
				}
			default:
				m[fmt.Sprint(k)] = ConvertMapInterfaceMapString(v2, prettify, isSchema)
			}
		}
		return m

	case []interface{}:
		x2 := make([]interface{}, len(x))
		for i, v2 := range x {
			x2[i] = ConvertMapInterfaceMapString(v2, prettify, isSchema)
		}
		return x2
	case map[string]interface{}:
		m := map[string]interface{}{}
		// foundFormatIntOrString := false
		for k, v2 := range x {
			if isSchema && k == "enum" { //While schema prettification, ENUMS are end system defined end user input and therefore should not be prettified/deprettified
				m[k] = v2
				continue
			}
			newmap := ConvertMapInterfaceMapString(v2, prettify, isSchema)
			if isSchema && isSpecialKey(k) {
				m[k] = newmap
			} else if prettify {
				m[manifests.FormatToReadableString(k)] = newmap
			} else {
				m[manifests.DeFormatReadableString(k)] = newmap
			}
		}
		return m
	case string:
		if isSchema {
			if prettify {
				return manifests.FormatToReadableString(x) //Whitespace formatting should be done at the time of prettification only
			}
			return manifests.DeFormatReadableString(x)
		}
	}
	return v
}

// These keys should not be prettified to "any Of", "all Of" and "one Of"
var keysToNotPrettifyOnSchema = []string{"anyOf", "allOf", "oneOf"}

func isSpecialKey(k string) bool {
	for _, k0 := range keysToNotPrettifyOnSchema {
		if k0 == k {
			return true
		}
	}
	return false
}

// In case of any breaking change or bug caused by this, set this to false and the whitespace addition in schema generated/consumed would be removed(will go back to default behavior)
const Format prettifier = true

type DryRunResponse2 struct {
	//When success is true, error will be nil and Component will contain the structure of the component as it will look after deployment
	//When success is false, error will contain the errors. And Component will be set to Nil
	Success   bool            `json:"success"`
	Error     *DryRunResponse `json:"error"`
	Component *Service        `json:"component"` //Service is synonymous with Component. Later Service is to be changed to "Component"
}
type DryRunResponse struct {
	Status string
	Causes []DryRunFailureCause
}

type DryRunFailureCause struct {
	Type      string //Type of error
	Message   string //Error message
	FieldPath string //Dot separated field path inside service. (For eg: <name>.settings.spec.containers (for pod) or <name>.annotations ) where <name> is the name of service/component
}

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
	Name         string            `yaml:"name,omitempty" json:"name,omitempty"`
	Type         string            `yaml:"type,omitempty" json:"type,omitempty"`
	APIVersion   string            `yaml:"apiVersion,omitempty" json:"apiVersion,omitempty"`
	Namespace    string            `yaml:"namespace,omitempty" json:"namespace,omitempty"`
	Version      string            `yaml:"version,omitempty" json:"version,omitempty"`
	Model        string            `yaml:"model,omitempty" json:"model,omitempty"`
	IsAnnotation bool              `yaml:"isAnnotation,omitempty" json:"isAnnotation,omitempty"`
	Labels       map[string]string `yaml:"labels,omitempty" json:"labels,omitempty"`
	Annotations  map[string]string `yaml:"annotations,omitempty" json:"annotations,omitempty"`
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

// isValidPattern checks if the pattern file is valid or not
func IsValidPattern(stringifiedFile string) (err error) {
	pattern := Pattern{}

	if err = yaml.Unmarshal([]byte(stringifiedFile), &pattern); err != nil {
		return err
	}

	if pattern.Services == nil {
		return errors.New("invalid design-file format: missing services field")
	}

	// for serviceName, service := range pattern.Services {
	// 	if service.Traits == nil {
	// 		return errors.New("missing traits field for:" + serviceName)
	// 	}
	// }
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
			Type:       svc.Type,
			Version:    svc.Version,
			Model:      svc.Model,
			APIVersion: svc.APIVersion,
			Settings:   svc.Settings,
		},
	}
	if comp.ObjectMeta.Labels == nil {
		comp.ObjectMeta.Labels = make(map[string]string)
	}
	comp.ObjectMeta.Labels["resource.pattern.meshery.io/id"] = svc.ID.String() //set the patternID to track back the object
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
// and creates a PatternFile from it.
// This function always returns meshkit error
func NewPatternFileFromCytoscapeJSJSON(name string, byt []byte) (Pattern, error) {
	// Unmarshal data into cytoscape struct
	var cy cytoscapejs.GraphElem
	if err := json.Unmarshal(byt, &cy); err != nil {
		return Pattern{}, ErrPatternFromCytoscape(err)
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
	err := processCytoElementsWithPattern(cy.Elements, func(svc Service, ele cytoscapejs.Element) error {
		name := svc.Name
		countDuplicates[name]++
		return nil
	})
	if err != nil {
		return pf, ErrPatternFromCytoscape(err)
	}

	//Populate the dependsOn field with appropriate unique service names
	err = processCytoElementsWithPattern(cy.Elements, func(svc Service, ele cytoscapejs.Element) error {
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

		//Only make the name unique when duplicates are encountered. This allows clients to preserve and propagate the unique name they want to give to their workload
		uniqueName := svc.Name
		if countDuplicates[uniqueName] > 1 {
			//set appropriate unique service name
			uniqueName = strings.ToLower(svc.Name)
			uniqueName += "-" + utils.GetRandomAlphabetsOfDigit(5)
		}
		eleToSvc[ele.Data.ID] = uniqueName //will be used while adding depends-on
		pf.Services[uniqueName] = &svc
		return nil
	})
	if err != nil {
		return pf, ErrPatternFromCytoscape(err)
	}
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
	return pf, nil
}

// processCytoElementsWithPattern iterates over all the cyto elements, convert each into a patternfile service and exposes a callback to handle that service
func processCytoElementsWithPattern(eles []cytoscapejs.Element, callback func(svc Service, ele cytoscapejs.Element) error) error {
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
		if svc.Name == "" {
			return fmt.Errorf("cannot save service with empty name")
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

// Note: If modified, make sure this function always returns a meshkit error
func NewPatternFileFromK8sManifest(data string, ignoreErrors bool, reg *meshmodel.RegistryManager) (Pattern, error) {
	pattern := Pattern{
		Name:     "Autogenerated",
		Services: map[string]*Service{},
	}

	manifests := strings.Split(data, "\n---\n")
	//For `---` separated manifests, even if only one manifest is there followed/preceded by multiple `\n---\n`- the manifest be will be valid
	//If there is no data present (except \n---\n) , then the yaml will be marked as empty and error will be thrown
	if manifestIsEmpty(manifests) {
		return pattern, ErrParseK8sManifest(fmt.Errorf("kubernetes manifest is empty"))
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

		name, svc, err := createPatternServiceFromK8s(manifest, reg)
		if err != nil {
			if ignoreErrors {
				continue
			}
			return pattern, ErrCreatePatternService(fmt.Errorf("failed to create pattern service from kubernetes component: %s", err))
		}

		pattern.Services[name] = &svc
	}

	return pattern, nil
}

func createPatternServiceFromK8s(manifest map[string]interface{}, regManager *meshmodel.RegistryManager) (string, Service, error) {
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

	// Get MeshModel entity with the selectors
	componentList, _, _ := regManager.GetEntities(&meshmodelv1alpha1.ComponentFilter{
		Name:       kind,
		APIVersion: apiVersion,
	})
	if len(componentList) == 0 {
		return "", Service{}, ErrCreatePatternService(fmt.Errorf("no resources found for APIVersion: %s Kind: %s", apiVersion, kind))
	}
	// just needs the first entry to grab meshmodel-metadata and other model requirements
	comp, ok := componentList[0].(meshmodelv1alpha1.ComponentDefinition)
	if !ok {
		return "", Service{}, ErrCreatePatternService(fmt.Errorf("cannot cast to the component-definition for APIVersion: %s Kind: %s", apiVersion, kind))
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
	rest = Format.Prettify(rest, false)
	uuidV4, _ := uuid.NewV4()
	svc := Service{
		Name:        name,
		Type:        comp.Kind,
		APIVersion:  comp.APIVersion,
		Namespace:   namespace,
		Model:       comp.Model.Name,
		Labels:      castedLabel,
		Annotations: castedAnnotation,
		Settings:    rest,
		Traits: map[string]interface{}{
			"meshmap": map[string]interface{}{
				"id":                 uuidV4,
				"meshmodel-metadata": comp.Metadata,
			},
		},
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
