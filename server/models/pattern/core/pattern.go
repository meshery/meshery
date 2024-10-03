package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/models/pattern/utils"
	"github.com/layer5io/meshkit/encoding"
	"github.com/layer5io/meshkit/logger"
	registry "github.com/layer5io/meshkit/models/meshmodel/registry"
	regv1beta1 "github.com/layer5io/meshkit/models/meshmodel/registry/v1beta1"
	mutils "github.com/layer5io/meshkit/utils"
	"github.com/layer5io/meshkit/utils/manifests"
	"github.com/meshery/schemas/models/v1beta1"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/pattern"
	cytoscapejs "gonum.org/v1/gonum/graph/formats/cytoscapejs"
	"gopkg.in/yaml.v2"
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

type DryRunResponseWrapper struct {
	//When success is true, error will be nil and Component will contain the structure of the component as it will look after deployment
	//When success is false, error will contain the errors. And Component will be set to Nil
	Success   bool                           `json:"success"`
	Error     *DryRunResponse                `json:"error"`
	Component *component.ComponentDefinition `json:"component"` //component.ComponentDefinition is synonymous with Component. Later component.ComponentDefinition is to be changed to "Component"
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

// NewPatternFile takes in raw yaml and encodes it into a construct
func NewPatternFile(yml []byte) (patternFile pattern.PatternFile, err error) {
	err = encoding.Unmarshal(yml, &patternFile)
	if err != nil {
		return patternFile, err
	}
	for _, component := range patternFile.Components {
		// If an explicit name is not given to the service then use
		// the service identifier as its name
		if component.DisplayName == "" {
			component.DisplayName = component.Id.String()
		}

		component.Configuration = utils.RecursiveCastMapStringInterfaceToMapStringInterface(component.Configuration)

		if component.Configuration == nil {
			component.Configuration = map[string]interface{}{}
		}
	}

	return
}

// AssignAdditionalLabels adds labels to identify resources deployed by meshery.
func AssignAdditionalLabels(comp *component.ComponentDefinition) error {

	if comp.Configuration == nil {
		comp.Configuration = make(map[string]interface{})
	}

	existingLabels := map[string]interface{}{}
	var err error

	if comp.Configuration["labels"] == nil {
		comp.Configuration["labels"] = make(map[string]interface{})
	} else {
		existingLabels, err = mutils.Cast[map[string]interface{}](comp.Configuration["labels"])
		if err != nil {
			return err
		}
	}

	existingLabels["resource.pattern.meshery.io/id"] = comp.Id.String() //set the patternID to track back the object
	comp.Configuration["labels"] = existingLabels
	return nil
}

// ToCytoscapeJS converts pattern file into cytoscape object
func ToCytoscapeJS(patternFile *pattern.PatternFile, log logger.Handler) (cytoscapejs.GraphElem, error) {
	var cy cytoscapejs.GraphElem

	// Not specifying any cytoscapejs layout
	// should fallback to "default" layout

	// Not specifying styles, may get applied on the
	// client side

	// Set up the nodes
	for _, cmp := range patternFile.Components {
		elemData := cytoscapejs.ElemData{
			ID: getCytoscapeElementID(cmp.Id.String(), cmp, log),
		}

		elemPosition, err := getCytoscapeJSPosition(cmp, log)
		if err != nil {
			return cy, err
		}

		elem := cytoscapejs.Element{
			Data:       elemData,
			Position:   &elemPosition,
			Selectable: true,
			Grabbable:  true,
			Scratch: map[string]component.ComponentDefinition{
				"_data": *cmp,
			},
		}

		cy.Elements = append(cy.Elements, elem)
	}

	return cy, nil
}

// NewPatternFileFromCytoscapeJSJSON takes in CytoscapeJS JSON
// and creates a PatternFile from it.
// This function always returns meshkit error
func NewPatternFileFromCytoscapeJSJSON(name string, byt []byte) (pattern.PatternFile, error) {
	// Unmarshal data into cytoscape struct
	var cy cytoscapejs.GraphElem
	if err := json.Unmarshal(byt, &cy); err != nil {
		return pattern.PatternFile{}, ErrPatternFromCytoscape(err)
	}

	if name == "" {
		name = "MesheryGeneratedPattern"
	}

	id, _ := uuid.NewV4()
	// Convert cytoscape struct to patternfile
	pf := pattern.PatternFile{
		Id:         id,
		Name:       name,
		Components: []*component.ComponentDefinition{},
	}

	// dependsOnMap := make(map[string][]string, 0) //used to figure out dependencies from componentmetadata.additionalProperties["dependsOn"]
	// eleToSvc := make(map[string]string)          //used to map cyto element ID uniquely to the name of the service created.
	countDuplicates := make(map[string]int)
	//store the names of services and their count
	err := processCytoElementsWithPattern(cy.Elements, func(comp component.ComponentDefinition, ele cytoscapejs.Element) error {
		name := comp.DisplayName
		countDuplicates[name]++
		return nil
	})

	if err != nil {
		return pf, ErrPatternFromCytoscape(err)
	}

	//Populate the dependsOn field with appropriate unique service names
	// err = processCytoElementsWithPattern(cy.Elements, func(declaration component.ComponentDefinition, ele cytoscapejs.Element) error {
	//Extract dependsOn, if present

	// dependencies, err := mutils.Cast[[]string](declaration.Metadata.AdditionalProperties["dependsOn"])

	// if err == nil {
	// 	declId := declaration.Id.String()
	// 	dependsOnMap[declId] = append(dependsOnMap[declId], dependencies...)
	// }

	// As client and server both depends on Id for determining uniqueness.
	// It isn't a problem if declaration have the same name.

	// Only make the name unique when duplicates are encountered. This allows clients to preserve and propagate the unique name they want to give to their workload
	// uniqueName := declaration.DisplayName
	// if countDuplicates[uniqueName] > 1 {
	// 	//set appropriate unique service name
	// 	uniqueName = strings.ToLower(uniqueName)
	// 	uniqueName += "-" + utils.GetRandomAlphabetsOfDigit(5)
	// }
	// eleToSvc[ele.Data.ID] = uniqueName //will be used while adding depends-on
	// pf.Services[uniqueName] = &svc
	// 	return nil
	// })

	// if err != nil {
	// 	return pf, ErrPatternFromCytoscape(err)
	// }

	// // add depends-on field
	// for child, parents := range dependsOnMap {
	// 	childSvc := eleToSvc[child]
	// 	if childSvc != "" {
	// 		for _, parent := range parents {
	// 			if eleToSvc[parent] != "" {
	// 				pf.Services[childSvc].DependsOn = append(pf.Services[childSvc].DependsOn, eleToSvc[parent])
	// 			}
	// 		}
	// 	}
	// }
	return pf, nil
}

// processCytoElementsWithPattern iterates over all the cyto elements, convert each into a patternfile service and exposes a callback to handle that service
func processCytoElementsWithPattern(eles []cytoscapejs.Element, callback func(svc component.ComponentDefinition, ele cytoscapejs.Element) error) error {
	for _, elem := range eles {
		// Try to create component.ComponentDefinition object from the elem.scratch's _data field
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
		declarationByt, err := json.Marshal(&data)
		if err != nil {
			return fmt.Errorf("failed to serialize component declaration from the metadata in the scratch")
		}

		// Unmarshal the JSON into a component declaration
		declaration := component.ComponentDefinition{
			Configuration: map[string]interface{}{},
		}

		// Add position
		declaration.Styles.Position.X = elem.Position.X
		declaration.Styles.Position.Y = elem.Position.Y

		if err := json.Unmarshal(declarationByt, &declaration); err != nil {
			return fmt.Errorf("failed to create component declaration from the metadata in the scratch")
		}
		if declaration.DisplayName == "" {
			return fmt.Errorf("cannot save design with empty name")
		}
		err = callback(declaration, elem)
		if err != nil {
			return err
		}
	}
	return nil
}

// Note: If modified, make sure this function always returns a meshkit error
func NewPatternFileFromK8sManifest(data string, fileName string, ignoreErrors bool, reg *registry.RegistryManager) (pattern.PatternFile, error) {
	if fileName == "" {
		fileName = "Autogenerated"
	}

	pattern := pattern.PatternFile{
		SchemaVersion: v1beta1.DesignSchemaVersion,
		Name:          fileName,
		Components:    []*component.ComponentDefinition{},
	}

	buffer := bytes.NewBufferString(data)
	decoder := yaml.NewDecoder(buffer)
	for {
		manifest := map[string]interface{}{}

		err := decoder.Decode(manifest)
		if err != nil {
			if err == io.EOF {
				if len(pattern.Components) == 0 {
					return pattern, ErrParseK8sManifest(fmt.Errorf("kubernetes manifest is empty"))
				}
				return pattern, nil
			} else {
				return pattern, ErrParseK8sManifest(fmt.Errorf("kubernetes manifest is empty"))
			}
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

		declaration, err := createPatternDeclarationFromK8s(manifest, reg)
		if err != nil {
			if ignoreErrors {
				continue
			}
			return pattern, ErrCreatePatternService(fmt.Errorf("failed to create design service from kubernetes component: %s", err))
		}

		pattern.Components = append(pattern.Components, &declaration)
	}

}

func createPatternDeclarationFromK8s(manifest map[string]interface{}, regManager *registry.RegistryManager) (component.ComponentDefinition, error) {
	fmt.Printf("%+#v\n", manifest)

	apiVersion, err := mutils.Cast[string](manifest["apiVersion"])
	if err != nil {
		return component.ComponentDefinition{}, ErrCreatePatternService(fmt.Errorf("invalid or empty apiVersion in manifest"))
	}

	kind, err := mutils.Cast[string](manifest["kind"])
	if err != nil {
		return component.ComponentDefinition{}, ErrCreatePatternService(fmt.Errorf("invalid or empty kind in manifest"))
	}

	metadata, _ := manifest["metadata"].(map[string]interface{})
	name, _ := metadata["name"].(string)

	rest := map[string]interface{}{}
	// rest will store a map of everything other than the above mentioned fields

	// Labels and annotations are stored inside the "component.Coonfiguration" hence not skipped while assiging manifest properties.
	for k, v := range manifest {
		// Ignore a few fields
		if k == "apiVersion" || k == "kind" || k == "status" {
			continue
		}

		rest[k] = v
	}

	// Get MeshModel entity with the selectors
	componentList, _, _, _ := regManager.GetEntities(&regv1beta1.ComponentFilter{
		Name:       kind,
		APIVersion: apiVersion,
	})

	if len(componentList) == 0 {
		return component.ComponentDefinition{}, ErrCreatePatternService(fmt.Errorf("no resources found for APIVersion: %s Kind: %s", apiVersion, kind))
	}

	// just needs the first entry to grab meshmodel-metadata and other model requirements
	comp, ok := componentList[0].(*component.ComponentDefinition)
	if !ok {
		return component.ComponentDefinition{}, ErrCreatePatternService(fmt.Errorf("cannot cast to the component-definition for APIVersion: %s Kind: %s", apiVersion, kind))
	}

	rest = Format.Prettify(rest, false)
	uuidV4, _ := uuid.NewV4()
	declaration := component.ComponentDefinition{
		Id:            uuidV4,
		SchemaVersion: comp.SchemaVersion,
		Version:       comp.Version,
		DisplayName:   name,
		Component: component.Component{
			Version: comp.Component.Version,
			Kind:    comp.Component.Kind,
		},
		Model:         comp.Model,
		Configuration: rest,
		Metadata:      comp.Metadata,
		Capabilities:  comp.Capabilities,
		Styles:        comp.Styles,
		Status:        comp.Status,
	}

	assignNamespaceForNamespacedScopedComp(&declaration, metadata, comp)
	return declaration, nil
}

func assignNamespaceForNamespacedScopedComp(declaration *component.ComponentDefinition, metadata map[string]interface{}, compDef *component.ComponentDefinition) *component.ComponentDefinition {
	if isNamespacedComponent(compDef) {
		namespace, _ := mutils.Cast[string](metadata["namespace"])
		if namespace == "" {
			metadata["namespace"] = "default"
		}
	}

	declaration.Configuration["metadata"] = metadata
	return declaration
}

// Checks whether the component is namespaced scope or not.
// While determining if an error occurs, the conversion process skips assigning a namespace value. If comp is originally namespaced scope, then k8s automatically assign a "default" namespace.
func isNamespacedComponent(comp *component.ComponentDefinition) bool {
	isNamespaced, _ := mutils.Cast[bool](comp.Metadata.AdditionalProperties["isNamespaced"])
	return isNamespaced
}

// getCytoscapeElementID returns the element id for a given service
func getCytoscapeElementID(name string, component *component.ComponentDefinition, log logger.Handler) string {
	return component.Id.String()
}

func getCytoscapeJSPosition(component *component.ComponentDefinition, log logger.Handler) (cytoscapejs.Position, error) {

	x, y := component.Styles.Position.X, component.Styles.Position.Y

	X := float64(x)
	Y := float64(y)

	pos := cytoscapejs.Position{
		X: X,
		Y: Y,
	}

	return pos, nil
}
