package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gofrs/uuid" // Use gofrs UUID package for the target struct
	// Use google UUID package for conversion

	"github.com/layer5io/meshkit/models/meshmodel/entity"
	"github.com/meshery/schemas/models/v1alpha1/capability"
	"github.com/meshery/schemas/models/v1beta1/category"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/meshery/schemas/models/v1beta1/model"
)

// Job struct for file processing
type Job struct {
	filePath string
}
type VersionMeta struct {
	SchemaVersion string `json:"schemaVersion,omitempty" yaml:"schemaVersion" gorm:"column:schema_version"`
	Version       string `json:"version,omitempty" yaml:"version" gorm:"column:version"`
}

type TypeMeta struct {
	Kind    string `json:"kind,omitempty" yaml:"kind"`
	Version string `json:"version,omitempty" yaml:"version"`
}

type ComponentFormat string

const (
	JSON ComponentFormat = "JSON"
	YAML ComponentFormat = "YAML"
	CUE  ComponentFormat = "CUE"

	ComponentSchemaVersion = "components.meshery.io/v1beta1"
)

// Contains information as extracted from the core underlying component eg: Pod's apiVersion, kind and schema
type ComponentEntity struct {
	TypeMeta
	Schema string `json:"schema,omitempty" yaml:"schema"`
}

// swagger:response ComponentDefinition
type ComponentDefinition struct {
	ID          uuid.UUID `json:"id"`
	VersionMeta `json:",inline" yaml:",inline"`
	DisplayName string                 `json:"displayName" gorm:"column:display_name"`
	Description string                 `json:"description" gorm:"column:description"`
	Format      ComponentFormat        `json:"format" yaml:"format"`
	ModelID     uuid.UUID              `json:"-" gorm:"index:idx_component_definition_dbs_model_id,column:model_id"`
	Model       Model                  `json:"model" gorm:"foreignKey:ModelID;references:ID"`
	Metadata    map[string]interface{} `json:"metadata" yaml:"metadata" gorm:"type:bytes;serializer:json"`
	Component   ComponentEntity        `json:"component,omitempty" yaml:"component" gorm:"type:bytes;serializer:json"`
}
type ModelEntity struct {
	Version string `json:"version,omitempty" yaml:"version"`
}

// swagger:response Model
type Model struct {
	ID          uuid.UUID `json:"id"`
	VersionMeta `yaml:",inline"`

	Name          string                 `json:"name" gorm:"modelName"`
	DisplayName   string                 `json:"displayName"`
	Description   string                 `json:"description" gorm:"description"`
	Status        entity.EntityStatus    `json:"status" gorm:"status"`
	RegistrantID  uuid.UUID              `json:"hostID" gorm:"column:host_id"` // make as a foreign refer to host's table
	Registrant    Host                   `json:"registrant" gorm:"foreignKey:RegistrantID;references:ID"`
	CategoryID    uuid.UUID              `json:"-" gorm:"categoryID"`
	Category      Category               `json:"category" gorm:"foreignKey:CategoryID;references:ID"`
	SubCategory   string                 `json:"subCategory" gorm:"subCategory"`
	Metadata      map[string]interface{} `json:"metadata" gorm:"type:bytes;serializer:json"`
	Model         ModelEntity            `json:"model,omitempty" gorm:"model;type:bytes;serializer:json"`
	Components    []ComponentDefinition  `json:"components" gorm:"-"`
	Relationships interface{}            `json:"relationships" gorm:"-"`
}
type Category struct {
	ID       uuid.UUID              `json:"-"`
	Name     string                 `json:"name" gorm:"name"`
	Metadata map[string]interface{} `json:"metadata"  yaml:"metadata" gorm:"type:bytes;serializer:json"`
}
type Host struct {
	ID        uuid.UUID `json:"-"`
	Hostname  string    `json:"hostname"`
	Port      int       `json:"port,omitempty"`
	Metadata  string    `json:"metadata,omitempty"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func main() {
	dirPath := "server/meshmodel" // Change this to your directory path

	// Create a channel for jobs and a wait group for synchronization
	jobs := make(chan Job, 100)
	var wg sync.WaitGroup

	// Start a fixed number of worker goroutines
	for w := 1; w <= 5; w++ { // You can adjust the number of workers
		wg.Add(1)
		go worker(jobs, &wg)
	}

	// Walk through the directory and add jobs to the channel
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			fileData, err := os.ReadFile(path)
			if err != nil {
				log.Printf("Failed to read file: %v", err)
				return nil
			}

			// Check if the file contains the schemaVersion
			if containsSchemaVersion(fileData, "components.meshery.io/v1beta1") {
				jobs <- Job{filePath: path}
			} else if containsSchemaVersion(fileData, "models.meshery.io/v1beta1") {
				createModel(path, fileData) // Directly create the model
			}
		}
		return nil
	})

	close(jobs)
	wg.Wait()

	if err != nil {
		log.Fatalf("Failed to walk directory: %v", err)
	}
}

func worker(jobs <-chan Job, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		processFile(job.filePath)
	}
}
func createModel(filePath string, fileData []byte) {
	// Unmarshal the old JSON structure into the old Model struct
	var oldModel Model
	if err := json.Unmarshal(fileData, &oldModel); err != nil {
		log.Printf("Failed to unmarshal JSON in file %s: %v", filePath, err)
		return
	}

	// Convert the google/uuid.UUID to gofrs/uuid.UUID
	oldID, err := uuid.FromString(oldModel.ID.String())
	if err != nil {
		log.Printf("Failed to convert UUID in file %s: %v", filePath, err)
		return
	}

	oldRegistrantID, err := uuid.FromString(oldModel.RegistrantID.String())
	if err != nil {
		log.Printf("Failed to convert RegistrantID UUID in file %s: %v", filePath, err)
		return
	}

	// Map metadata efficiently
	newModelMetadata := mapMetadata(oldModel.Metadata)

	// Map the old Model to the new ModelDefinition
	newRegistrant := createNewRegistrant(oldRegistrantID, oldModel.Registrant.Hostname)

	oldCategoryID, err := uuid.FromString(oldModel.CategoryID.String())
	if err != nil {
		log.Printf("Failed to convert UUID in file %s: %v", filePath, err)
		return
	}
	cat := category.CategoryDefinition{
		Name: oldModel.Category.Name,
	}
	// Map the old Model to the new ModelDefinition
	newModel := model.ModelDefinition{
		Id:            oldID,
		Name:          oldModel.Name,
		DisplayName:   oldModel.DisplayName,
		Model:         model.Model(oldModel.Model),
		Description:   oldModel.Description,
		Status:        model.ModelDefinitionStatus(oldModel.Status),
		CategoryId:    oldCategoryID,
		Category:      cat,
		SchemaVersion: oldModel.SchemaVersion,
		SubCategory:   oldModel.SubCategory,
		Metadata:      &newModelMetadata,
		Registrant:    newRegistrant,
		Version:       oldModel.Version,
	}
	modifiedData, err := json.MarshalIndent(newModel, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal JSON in file %s: %v", filePath, err)
		return
	}

	// Write the updated JSON to a new file
	if err := os.WriteFile(filePath, modifiedData, 0644); err != nil {
		log.Printf("Failed to write file %s: %v", filePath, err)
		return
	}

	fmt.Println("File processed and saved:", filePath)
}
func processFile(filePath string) {
	// Read the JSON file
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Failed to read file %s: %v", filePath, err)
		return
	}

	// Unmarshal the old JSON structure into the old ComponentDefinition struct
	var oldCompDef ComponentDefinition
	if err := json.Unmarshal(fileData, &oldCompDef); err != nil {
		log.Printf("Failed to unmarshal JSON in file %s: %v", filePath, err)
		return
	}

	// Convert the google/uuid.UUID to gofrs/uuid.UUID
	oldID, err := uuid.FromString(oldCompDef.Model.ID.String())
	if err != nil {
		log.Printf("Failed to convert UUID in file %s: %v", filePath, err)
		return
	}

	oldRegistrantID, err := uuid.FromString(oldCompDef.Model.RegistrantID.String())
	if err != nil {
		log.Printf("Failed to convert RegistrantID UUID in file %s: %v", filePath, err)
		return
	}

	// Prepare new capabilities data
	newCapabilities := createNewCapabilities()

	// Map metadata efficiently
	newModelMetadata := mapMetadata(oldCompDef.Model.Metadata)

	// Map the old Model to the new ModelDefinition
	newRegistrant := createNewRegistrant(oldRegistrantID, oldCompDef.Model.Registrant.Hostname)

	oldCategoryID, err := uuid.FromString(oldCompDef.Model.CategoryID.String())
	if err != nil {
		log.Printf("Failed to convert UUID in file %s: %v", filePath, err)
		return
	}
	cat := category.CategoryDefinition{
		Name: oldCompDef.Model.Category.Name,
	}
	// Map the old Model to the new ModelDefinition
	newModel := model.ModelDefinition{
		Id:            oldID,
		Name:          oldCompDef.Model.Name,
		DisplayName:   oldCompDef.Model.DisplayName,
		Model:         model.Model(oldCompDef.Model.Model),
		Description:   oldCompDef.Model.Description,
		Status:        model.ModelDefinitionStatus(oldCompDef.Model.Status),
		CategoryId:    oldCategoryID,
		Category:      cat,
		SchemaVersion: oldCompDef.Model.SchemaVersion,
		SubCategory:   oldCompDef.Model.SubCategory,
		Metadata:      &newModelMetadata,
		Registrant:    newRegistrant,
		Version:       oldCompDef.Model.Version,
	}

	oldCompId, err := uuid.FromString(oldCompDef.ID.String())
	if err != nil {
		log.Printf("Failed to convert UUID in file %s: %v", filePath, err)
		return
	}

	// Map the old ComponentDefinition to the new structure
	newCompDef := createNewComponentDef(oldCompDef, newModel, oldCompId, newCapabilities)

	// Marshal the new structure back to JSON
	modifiedData, err := json.MarshalIndent(newCompDef, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal JSON in file %s: %v", filePath, err)
		return
	}

	// Write the updated JSON to a new file
	if err := os.WriteFile(filePath, modifiedData, 0644); err != nil {
		log.Printf("Failed to write file %s: %v", filePath, err)
		return
	}

	fmt.Println("File processed and saved:", filePath)
}

func containsSchemaVersion(fileData []byte, schemaVersion string) bool {
	var temp map[string]interface{}
	if err := json.Unmarshal(fileData, &temp); err != nil {
		return false
	}

	if val, ok := temp["schemaVersion"].(string); ok {
		return val == schemaVersion
	}

	return false
}

func createNewCapabilities() []capability.Capability {
	return []capability.Capability{
		{
			SchemaVersion: "capability.meshery.io/v1alpha1",
			Version:       "0.7.0",
			DisplayName:   "Performance Test",
			Description:   "Initiate a performance test. Meshery will execute the load generation, collect metrics, and present the results.",
			Kind:          "action",
			Type:          "operator",
			SubType:       "perf-test",
			EntityState:   &[]capability.CapabiliyEntityState{"declaration", "instance"},
			Status:        "enabled",
		},
		{
			SchemaVersion: "capability.meshery.io/v1alpha1",
			Version:       "0.7.0",
			DisplayName:   "Workload Configuration",
			Description:   "Configure the workload-specific setting of a component",
			Kind:          "mutate",
			Type:          "configuration",
			SubType:       "config",
			EntityState:   &[]capability.CapabiliyEntityState{"declaration"},
			Status:        "enabled",
		},
		{
			SchemaVersion: "capability.meshery.io/v1alpha1",
			Version:       "0.7.0",
			DisplayName:   "Labels and Annotations Configuration",
			Description:   "Configure Labels And Annotations for the component",
			Kind:          "mutate",
			Type:          "configuration",
			SubType:       "labels-and-annotations",
			EntityState:   &[]capability.CapabiliyEntityState{"declaration"},
			Status:        "enabled",
		},
		{
			SchemaVersion: "capability.meshery.io/v1alpha1",
			Version:       "0.7.0",
			DisplayName:   "Relationships",
			Description:   "View relationships for the component",
			Kind:          "view",
			Type:          "configuration",
			SubType:       "relationship",
			EntityState:   &[]capability.CapabiliyEntityState{"declaration", "instance"},
			Status:        "enabled",
		},
	}
}

func mapMetadata(oldMetadata map[string]interface{}) model.ModelDefinition_Metadata {
	newModelMetadata := model.ModelDefinition_Metadata{
		AdditionalProperties: make(map[string]interface{}),
	}

	for key, value := range oldMetadata {
		switch key {
		case "capabilities":
			if capabilities, ok := value.([]capability.Capability); ok {
				newModelMetadata.Capabilities = &capabilities
			}
		case "isAnnotation":
			if isAnnotation, ok := value.(bool); ok {
				newModelMetadata.IsAnnotation = &isAnnotation
			}
		case "primaryColor":
			if primaryColor, ok := value.(string); ok {
				newModelMetadata.PrimaryColor = &primaryColor
			}
		case "secondaryColor":
			if secondaryColor, ok := value.(string); ok {
				newModelMetadata.SecondaryColor = &secondaryColor
			}
		case "svgColor":
			if svgColor, ok := value.(string); ok {
				newModelMetadata.SvgColor = svgColor
			}
		case "svgComplete":
			if svgComplete, ok := value.(string); ok {
				newModelMetadata.SvgComplete = &svgComplete
			}
		case "svgWhite":
			if svgWhite, ok := value.(string); ok {
				newModelMetadata.SvgWhite = svgWhite
			}
		default:
			newModelMetadata.AdditionalProperties[key] = value
		}
	}
	return newModelMetadata
}

func createNewRegistrant(oldRegistrantID uuid.UUID, hostname string) connection.Connection {
	newRegistrant := connection.Connection{
		Id:     oldRegistrantID,
		Kind:   hostname,
		Status: "connection.Registered",
		Type:   "registry",
	}

	switch newRegistrant.Kind {
	case "artifacthub":
		newRegistrant.Name = "Artifact Hub"
	case "github":
		newRegistrant.Name = "GitHub"
	case "kubernetes":
		newRegistrant.Name = "Kubernetes"
	default:
		newRegistrant.Kind = "Unknown"
	}

	return newRegistrant
}

func createNewComponentDef(oldCompDef ComponentDefinition, newModel model.ModelDefinition, oldCompId uuid.UUID, newCapabilities []capability.Capability) component.ComponentDefinition {
	status := getString(oldCompDef.Metadata, "status")

	return component.ComponentDefinition{
		Id:            oldCompId,
		SchemaVersion: oldCompDef.SchemaVersion,
		Version:       oldCompDef.Version,
		DisplayName:   oldCompDef.DisplayName,
		Description:   oldCompDef.Description,
		Model:         newModel,
		Format:        component.ComponentDefinitionFormat(oldCompDef.Format),
		Status:        (*component.ComponentDefinitionStatus)(&status),
		Metadata: component.ComponentDefinition_Metadata{
			Genealogy: getString(oldCompDef.Metadata, "genealogy"),
			IsAnnotation: func() bool {
				if value, ok := oldCompDef.Metadata["isAnnotation"].(bool); ok {
					return value
				}
				return false
			}(),
			Published: func() bool {
				if value, ok := oldCompDef.Metadata["published"].(bool); ok {
					return value
				}
				return false
			}(),
			AdditionalProperties: func() map[string]interface{} {
				additionalProps := make(map[string]interface{})
				for key, value := range oldCompDef.Metadata {
					if key != "genealogy" && key != "isAnnotation" && key != "published" && key != "capabilities" {
						switch key {
						case "active-bg-color", "active-bg-opacity", "active-bg-size", "background-blacken", "background-color",
							"background-opacity", "body-text", "body-text-background-color", "body-text-color", "body-text-decoration",
							"body-text-font-size", "body-text-horizontal-align", "body-text-max-width", "body-text-opacity",
							"body-text-vertical-align", "body-text-weight", "body-text-wrap", "border-color", "border-opacity",
							"border-style", "border-width", "color", "font-family", "font-size", "font-style", "font-weight", "ghost",
							"height", "label", "opacity", "outside-texture-bg-color", "outside-texture-bg-opacity", "padding", "position",
							"primaryColor", "secondaryColor", "selection-box-border-width", "selection-box-color", "selection-box-opacity",
							"shape", "svgColor", "svgComplete", "svgWhite", "text-halign", "text-opacity", "text-transform", "text-valign",
							"width", "z-index", "styles", "subCategory", "styleOverrides", "status":
							// Do not add these keys to additionalProps
						default:
							additionalProps[key] = value
						}
					}
				}
				return additionalProps
			}(),
		},
		Component: component.Component{
			Kind:    oldCompDef.Component.Kind,
			Schema:  oldCompDef.Component.Schema,
			Version: oldCompDef.Component.Version,
		},
		Capabilities: &newCapabilities,
		Styles:       mapStyles(oldCompDef.Metadata),
	}
}

func mapStyles(metadata map[string]interface{}) *component.Styles {
	return &component.Styles{
		ActiveBgColor:           getStringPointer(metadata, "active-bg-color"),
		ActiveBgOpacity:         getStringPointer(metadata, "active-bg-opacity"),
		ActiveBgSize:            getStringPointer(metadata, "active-bg-size"),
		BackgroundBlacken:       getFloat32Pointer(metadata, "background-blacken"),
		BackgroundColor:         getStringPointer(metadata, "background-color"),
		BackgroundOpacity:       getFloat32Pointer(metadata, "background-opacity"),
		BodyText:                getStringPointer(metadata, "body-text"),
		BodyTextBackgroundColor: getStringPointer(metadata, "body-text-background-color"),
		BodyTextColor:           getStringPointer(metadata, "body-text-color"),
		BodyTextDecoration:      getStringPointer(metadata, "body-text-decoration"),
		BodyTextFontSize:        getStringPointer(metadata, "body-text-font-size"),
		BodyTextHorizontalAlign: getStringPointer(metadata, "body-text-horizontal-align"),
		BodyTextMaxWidth:        getStringPointer(metadata, "body-text-max-width"),
		BodyTextOpacity:         getFloat32Pointer(metadata, "body-text-opacity"),
		BodyTextVerticalAlign:   getStringPointer(metadata, "body-text-vertical-align"),
		BodyTextWeight:          getStringPointer(metadata, "body-text-weight"),
		BodyTextWrap:            getStringPointer(metadata, "body-text-wrap"),
		BorderColor:             getStringPointer(metadata, "border-color"),
		BorderOpacity:           getFloat32Pointer(metadata, "border-opacity"),
		BorderStyle:             getBorderStylePointer(metadata, "border-style"),
		BorderWidth:             getFloat32Pointer(metadata, "border-width"),
		Color:                   getStringPointer(metadata, "color"),
		FontFamily:              getStringPointer(metadata, "font-family"),
		FontSize:                getStringPointer(metadata, "font-size"),
		FontStyle:               getStringPointer(metadata, "font-style"),
		FontWeight:              getStringPointer(metadata, "font-weight"),
		Ghost:                   getGhostPointer(metadata, "ghost"),
		Height:                  getFloat32Pointer(metadata, "height"),
		Label:                   getStringPointer(metadata, "label"),
		Opacity:                 getFloat32Pointer(metadata, "opacity"),
		OutsideTextureBgColor:   getStringPointer(metadata, "outside-texture-bg-color"),
		OutsideTextureBgOpacity: getFloat32Pointer(metadata, "outside-texture-bg-opacity"),
		Padding:                 getFloat32Pointer(metadata, "padding"),
		PrimaryColor:            getString(metadata, "primaryColor"),
		SecondaryColor:          getStringPointer(metadata, "secondaryColor"),
		SelectionBoxBorderWidth: getFloat32Pointer(metadata, "selection-box-border-width"),
		SelectionBoxColor:       getStringPointer(metadata, "selection-box-color"),
		SelectionBoxOpacity:     getFloat32Pointer(metadata, "selection-box-opacity"),
		Shape:                   getShapePointer(metadata, "shape"),
		SvgColor:                getString(metadata, "svgColor"),
		SvgComplete: func() string {
			if val := getStringPointer(metadata, "svgComplete"); val != nil {
				return *val
			}
			return "" // or some default value
		}(), SvgWhite: getString(metadata, "svgWhite"),
		TextHalign:    getTextHalignPointer(metadata, "text-halign"),
		TextOpacity:   getFloat32Pointer(metadata, "text-opacity"),
		TextTransform: getTextTransformPointer(metadata, "text-transform"),
		TextValign:    getTextValignPointer(metadata, "text-valign"),
		Width:         getFloat32Pointer(metadata, "width"),
		ZIndex:        getIntPointer(metadata, "z-index"),
	}
}

func getString(metadata map[string]interface{}, key string) string {
	if value, ok := metadata[key].(string); ok {
		return value
	}
	return "" // default value
}

func getStringPointer(metadata map[string]interface{}, key string) *string {
	if value, ok := metadata[key].(string); ok {
		return &value
	}
	return nil
}

func getFloat32Pointer(metadata map[string]interface{}, key string) *float32 {
	if value, ok := metadata[key].(float32); ok {
		return &value
	}
	return nil
}

func getIntPointer(metadata map[string]interface{}, key string) *int {
	if value, ok := metadata[key].(int); ok {
		return &value
	}
	return nil
}

func getBorderStylePointer(metadata map[string]interface{}, key string) *component.ComponentDefinitionStylesBorderStyle {
	if value, ok := metadata[key].(component.ComponentDefinitionStylesBorderStyle); ok {
		return &value
	}
	return nil
}

func getGhostPointer(metadata map[string]interface{}, key string) *component.ComponentDefinitionStylesGhost {
	if value, ok := metadata[key].(component.ComponentDefinitionStylesGhost); ok {
		return &value
	}
	return nil
}

func getShapePointer(metadata map[string]interface{}, key string) *component.ComponentDefinitionStylesShape {
	if value, ok := metadata[key].(component.ComponentDefinitionStylesShape); ok {
		return &value
	}
	return nil
}

func getTextHalignPointer(metadata map[string]interface{}, key string) *component.ComponentDefinitionStylesTextHalign {
	if value, ok := metadata[key].(component.ComponentDefinitionStylesTextHalign); ok {
		return &value
	}
	return nil
}

func getTextTransformPointer(metadata map[string]interface{}, key string) *component.ComponentDefinitionStylesTextTransform {
	if value, ok := metadata[key].(component.ComponentDefinitionStylesTextTransform); ok {
		return &value
	}
	return nil
}

func getTextValignPointer(metadata map[string]interface{}, key string) *component.ComponentDefinitionStylesTextValign {
	if value, ok := metadata[key].(component.ComponentDefinitionStylesTextValign); ok {
		return &value
	}
	return nil
}