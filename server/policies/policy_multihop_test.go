package policies

import (
	"testing"

	"github.com/google/uuid"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/pattern"
	"github.com/meshery/schemas/models/v1beta2/relationship"
)

func TestMultihopValidationPolicy_IsImplicatedBy(t *testing.T) {
	policy := &MultihopValidationPolicy{}

	tgwKind := "TransitGateway"
	attachmentKind := "TransitGatewayVPCAttachment"

	implicatedRel := &relationship.RelationshipDefinition{
		Kind:             "edge",
		RelationshipType: "binding",
		Selectors: &relationship.SelectorSet{
			{
				Allow: relationship.Selector{
					From: []relationship.SelectorItem{{Kind: &tgwKind}},
					To:   []relationship.SelectorItem{{Kind: &attachmentKind}},
				},
			},
		},
	}

	nonImplicatedRel := &relationship.RelationshipDefinition{
		Kind:             "edge",
		RelationshipType: "non-binding",
		Selectors: &relationship.SelectorSet{
			{
				Allow: relationship.Selector{
					From: []relationship.SelectorItem{{Kind: &tgwKind}},
					To:   []relationship.SelectorItem{{Kind: &attachmentKind}},
				},
			},
		},
	}

	if !policy.IsImplicatedBy(implicatedRel) {
		t.Error("Expected relationship to be implicated by MultihopValidationPolicy")
	}

	if policy.IsImplicatedBy(nonImplicatedRel) {
		t.Error("Expected relationship to not be implicated by MultihopValidationPolicy")
	}
}

func TestMultihopValidationPolicy_IsInvalid(t *testing.T) {
	policy := &MultihopValidationPolicy{}

	tgwID := uuid.New()
	vpcID := uuid.New()
	attID := uuid.New()

	tgwKind := "TransitGateway"
	vpcKind := "VPC"
	attKind := "TransitGatewayVPCAttachment"

	designComplete := &pattern.PatternFile{
		Components: []*component.ComponentDefinition{
			{
				ID: tgwID,
				Component: component.Component{
					Kind: "TransitGateway",
				},
			},
			{
				ID: vpcID,
				Component: component.Component{
					Kind: "VPC",
				},
			},
			{
				ID: attID,
				Component: component.Component{
					Kind: "TransitGatewayVPCAttachment",
				},
			},
		},
		Relationships: []*relationship.RelationshipDefinition{
			{
				Kind:             "edge",
				RelationshipType: "binding",
				Selectors: &relationship.SelectorSet{
					{
						Allow: relationship.Selector{
							From: []relationship.SelectorItem{{Kind: &tgwKind, ID: &tgwID}},
							To:   []relationship.SelectorItem{{Kind: &attKind, ID: &attID}},
						},
					},
				},
			},
			{
				Kind:             "edge",
				RelationshipType: "binding",
				Selectors: &relationship.SelectorSet{
					{
						Allow: relationship.Selector{
							From: []relationship.SelectorItem{{Kind: &vpcKind, ID: &vpcID}},
							To:   []relationship.SelectorItem{{Kind: &attKind, ID: &attID}},
						},
					},
				},
			},
		},
	}

	relToTest := designComplete.Relationships[0]

	// Complete connections: should NOT be invalid
	if policy.IsInvalid(relToTest, designComplete) {
		t.Error("Expected relationship to be valid when both TGW and VPC connections exist")
	}

	// Missing VPC connection: should be invalid
	designIncomplete := &pattern.PatternFile{
		Components: []*component.ComponentDefinition{
			{
				ID: tgwID,
				Component: component.Component{
					Kind: "TransitGateway",
				},
			},
			{
				ID: attID,
				Component: component.Component{
					Kind: "TransitGatewayVPCAttachment",
				},
			},
		},
		Relationships: []*relationship.RelationshipDefinition{
			{
				Kind:             "edge",
				RelationshipType: "binding",
				Selectors: &relationship.SelectorSet{
					{
						Allow: relationship.Selector{
							From: []relationship.SelectorItem{{Kind: &tgwKind, ID: &tgwID}},
							To:   []relationship.SelectorItem{{Kind: &attKind, ID: &attID}},
						},
					},
				},
			},
		},
	}

	if !policy.IsInvalid(relToTest, designIncomplete) {
		t.Error("Expected relationship to be invalid when VPC connection is missing")
	}
}
