package policies

import (
	"strings"

	"github.com/meshery/schemas/models/v1beta1/pattern"
	"github.com/meshery/schemas/models/v1beta2/relationship"
)

// MultihopValidationPolicy handles validation of multi-hop network routing connections.
type MultihopValidationPolicy struct{}

// Identifier returns the unique identifier for this policy.
func (p *MultihopValidationPolicy) Identifier() string {
	return "multihop-routing-validation"
}

// IsImplicatedBy checks if the relationship is implicated by the multi-hop policy.
func (p *MultihopValidationPolicy) IsImplicatedBy(rel *relationship.RelationshipDefinition) bool {
	if !strings.EqualFold(string(rel.Kind), "edge") || !strings.EqualFold(rel.RelationshipType, "binding") {
		return false
	}
	if rel.Selectors == nil {
		return false
	}
	for _, ss := range *rel.Selectors {
		for _, toSel := range ss.Allow.To {
			if strings.EqualFold(toSel.Kind, "TransitGatewayVPCAttachment") {
				return true
			}
		}
	}
	return false
}

// IsInvalid checks if the relationship forms an invalid multi-hop routing chain.
func (p *MultihopValidationPolicy) IsInvalid(rel *relationship.RelationshipDefinition, design *pattern.PatternFile) bool {
	if rel.Selectors == nil {
		return false
	}
	for _, ss := range *rel.Selectors {
		for _, toSel := range ss.Allow.To {
			if strings.EqualFold(toSel.Kind, "TransitGatewayVPCAttachment") {
				if toSel.ID == nil {
					continue
				}
				attachmentID := toSel.ID.String()
				if p.isIncompleteTgwVpcAttachment(attachmentID, design) {
					return true
				}
			}
		}
	}
	return false
}

// AlreadyExists checks if the relationship already exists in the design.
func (p *MultihopValidationPolicy) AlreadyExists(rel *relationship.RelationshipDefinition, design *pattern.PatternFile) bool {
	return relationshipAlreadyExists(design, rel)
}

// IdentifyRelationship is handled by the general EdgeBindingPolicy.
func (p *MultihopValidationPolicy) IdentifyRelationship(relDef *relationship.RelationshipDefinition, design *pattern.PatternFile) []*relationship.RelationshipDefinition {
	return nil
}

// SideEffects returns update actions to mark incomplete/dangling attachments as deleted.
func (p *MultihopValidationPolicy) SideEffects(rel *relationship.RelationshipDefinition, design *pattern.PatternFile) []PolicyAction {
	if getRelStatus(rel) == StatusApproved {
		if rel.Selectors == nil {
			return nil
		}
		for _, ss := range *rel.Selectors {
			for _, toSel := range ss.Allow.To {
				if strings.EqualFold(toSel.Kind, "TransitGatewayVPCAttachment") {
					if toSel.ID == nil {
						continue
					}
					attachmentID := toSel.ID.String()
					if p.isIncompleteTgwVpcAttachment(attachmentID, design) {
						return []PolicyAction{
							newUpdateRelationshipAction(rel.ID.String(), "/status", StatusDeleted),
						}
					}
				}
			}
		}
	}
	return nil
}

// isIncompleteTgwVpcAttachment checks if the attachment is missing either TransitGateway or VPC connections.
func (p *MultihopValidationPolicy) isIncompleteTgwVpcAttachment(attachmentID string, design *pattern.PatternFile) bool {
	return !p.hasTransitGatewayConnection(attachmentID, design) || !p.hasVpcConnection(attachmentID, design)
}

// hasTransitGatewayConnection checks for an active TransitGateway connection to the attachment.
func (p *MultihopValidationPolicy) hasTransitGatewayConnection(attachmentID string, design *pattern.PatternFile) bool {
	for _, rel := range design.Relationships {
		if getRelStatus(rel) == StatusDeleted {
			continue
		}
		if rel.Selectors == nil {
			continue
		}
		for _, ss := range *rel.Selectors {
			for _, fromSel := range ss.Allow.From {
				if strings.EqualFold(fromSel.Kind, "TransitGateway") {
					for _, toSel := range ss.Allow.To {
						if toSel.ID != nil && toSel.ID.String() == attachmentID {
							return true
						}
					}
				}
			}
		}
	}
	return false
}

// hasVpcConnection checks for an active VPC connection to the attachment.
func (p *MultihopValidationPolicy) hasVpcConnection(attachmentID string, design *pattern.PatternFile) bool {
	for _, rel := range design.Relationships {
		if getRelStatus(rel) == StatusDeleted {
			continue
		}
		if rel.Selectors == nil {
			continue
		}
		for _, ss := range *rel.Selectors {
			for _, fromSel := range ss.Allow.From {
				if strings.EqualFold(fromSel.Kind, "VPC") {
					for _, toSel := range ss.Allow.To {
						if toSel.ID != nil && toSel.ID.String() == attachmentID {
							return true
						}
					}
				}
			}
		}
	}
	return false
}
