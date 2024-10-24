package utils

import (
	"testing"

	"github.com/layer5io/meshkit/models/meshmodel/core/v1beta1"
	"github.com/layer5io/meshkit/models/meshmodel/entity"
	"github.com/stretchr/testify/assert"
)

func TestUpdateModelDefinition(t *testing.T) {
	modelCSV := &ModelCSV{
		IsAnnotation: "true",
		PrimaryColor: "#FFFFFF",
	}

	modelDef := &v1beta1.Model{
		Metadata: map[string]interface{}{
			"existingKey": "existingValue",
		},
	}

	err := modelCSV.UpdateModelDefinition(modelDef)
	assert.NoError(t, err)
	assert.Equal(t, true, modelDef.Metadata["isAnnotation"])
	assert.Equal(t, "#FFFFFF", modelDef.Metadata["primaryColor"])
	assert.Equal(t, "existingValue", modelDef.Metadata["existingKey"])
}

func TestCreateModelDefinition(t *testing.T) {
	modelCSV := &ModelCSV{
		Registrant:        "Test Registrant",
		ModelDisplayName:  "Test Model",
		Model:             "test-model",
		Category:          "Test Category",
		SubCategory:       "Test SubCategory",
		PublishToRegistry: "true",
	}

	model := modelCSV.CreateModelDefinition("v1", "v1beta1")
	assert.Equal(t, "test-model", model.Name)
	assert.Equal(t, "Test Model", model.DisplayName)
	assert.Equal(t, entity.Enabled, model.Status)
	assert.Equal(t, "testregistrant", model.Registrant.Hostname)
	assert.Equal(t, "Test Category", model.Category.Name)
	assert.Equal(t, "Test SubCategory", model.SubCategory)
}

func TestParseModelsSheet(t *testing.T) {
	// Initialize ModelCSVHelper with non-existent CSV path
	mch := &ModelCSVHelper{
		CSVPath: "/non/existent/path/models.csv",
		Models:  []ModelCSV{},
	}

	err := mch.ParseModelsSheet(false, "xy")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no such file or directory")
}

func TestCreateMarkDownForMDXStyle(t *testing.T) {
	m := ModelCSV{
		Model:             "Test Model",
		ModelDisplayName:  "Test Model Display Name",
		PageSubtTitle:     "Test Subtitle",
		DocsURL:           "http://example.com",
		Description:       "Test Description",
		Category:          "Test Category",
		SubCategory:       "Test SubCategory",
		Registrant:        "Test Registrant",
		Feature1:          "Feature 1",
		Feature2:          "Feature 2",
		Feature3:          "Feature 3",
		HowItWorks:        "How it works",
		HowItWorksDetails: "How it works details",
		PublishToSites:    "true",
		AboutProject:      "About the project",
		StandardBlurb:     "Standard blurb",
	}

	expectedMarkdown := `---
title: Test Model Display Name
subtitle: Test Subtitle
integrationIcon: icons/color/test-model-color.svg
darkModeIntegrationIcon: icons/white/test-model-white.svg
docURL: http://example.com
description: Test Description
category: Test Category
subcategory: Test SubCategory
registrant: Test Registrant
components: test components
featureList: [
  "Feature 1",
  "Feature 2",
  "Feature 3"
]
workingSlides: [
  ../_images/meshmap-visualizer.png,
  ../_images/meshmap-designer.png
]
howItWorks: "How it works"
howItWorksDetails: "How it works details"
published: true
---
<p>
About the project
</p>
Standard blurb
`

	markdown := m.CreateMarkDownForMDXStyle("test components")
	assert.Equal(t, expectedMarkdown, markdown)
}

func TestCreateJSONItem(t *testing.T) {
	m := ModelCSV{
		Model:    "Test Model",
		SVGColor: "color.svg",
		SVGWhite: "white.svg",
		DocsURL:  "http://example.com",
	}

	expectedJSON := `{"color":"icons/icons/color/test-model-color.svg","name":"Test Model","permalink":"http://example.com","white":"icons/icons/white/test-model-white.svg"}`

	jsonItem := m.CreateJSONItem("icons")
	assert.JSONEq(t, expectedJSON, jsonItem)
}

func TestCreateMarkDownForMDStyle(t *testing.T) {
	m := ModelCSV{
		Model:             "Test Model",
		ModelDisplayName:  "Test Model Display Name",
		PageSubtTitle:     "Test Subtitle",
		DocsURL:           "http://example.com",
		Description:       "Test Description",
		Category:          "Test Category",
		SubCategory:       "Test SubCategory",
		Registrant:        "Test Registrant",
		Feature1:          "Feature 1",
		Feature2:          "Feature 2",
		Feature3:          "Feature 3",
		HowItWorks:        "How it works",
		HowItWorksDetails: "How it works details",
	}

	expectedMarkdown := `---
layout: integration
title: Test Model Display Name
subtitle: Test Subtitle
image: /assets/img/integrations/test-model/icons/color/test-model-color.svg
permalink: extensibility/integrations/test-model
docURL: http://example.com
description: Test Description
integrations-category: Test Category
integrations-subcategory: Test SubCategory
registrant: Test Registrant
components: test components
featureList: [
  "Feature 1",
  "Feature 2",
  "Feature 3"
]
howItWorks: "How it works"
howItWorksDetails: "How it works details"
language: en
list: include
type: extensibility
category: integrations
---
`

	markdown := m.CreateMarkDownForMDStyle("test components")
	assert.Equal(t, expectedMarkdown, markdown)
}
