package eval

import rego.v1

import data.actions
import data.eval_rules

# Policy Identifier for Multi-hop Routing validation
multihop_policy_identifier := "multihop-routing-validation"

# Implicate relationships where the target is a TransitGatewayVPCAttachment.
# This ensures we validate the multi-hop hub-and-spoke connection integrity.
relationship_is_implicated_by_policy(relationship, policy_identifier) if {
	policy_identifier == multihop_policy_identifier
	lower(relationship.kind) == "edge"
	lower(relationship.type) == "binding"
	
	some selector in relationship.selectors
	some to_clause in selector.allow.to
	to_clause.kind == "TransitGatewayVPCAttachment"
}

# A relationship is invalid if its referenced from or to components no longer exist in the design.
relationship_is_invalid(relationship, design_file, policy_identifier) if {
	policy_identifier == multihop_policy_identifier
	eval_rules.from_or_to_components_dont_exist(relationship, design_file)
}

# A relationship targetting a TransitGatewayVPCAttachment is invalid if it forms an incomplete connection hop
# (i.e. it lacks either the corresponding TransitGateway or VPC connection).
relationship_is_invalid(relationship, design_file, policy_identifier) if {
	policy_identifier == multihop_policy_identifier
	
	some selector in relationship.selectors
	some to_clause in selector.allow.to
	to_clause.kind == "TransitGatewayVPCAttachment"
	attachment_id := to_clause.id
	attachment_id != null
	
	is_incomplete_tgw_vpc_attachment(attachment_id, design_file)
}

# A TransitGatewayVPCAttachment connection is incomplete if it does not have BOTH a TransitGateway
# and a VPC connected to it in the current design.
is_incomplete_tgw_vpc_attachment(attachment_id, design_file) if {
	not has_transit_gateway_connection(attachment_id, design_file)
}

is_incomplete_tgw_vpc_attachment(attachment_id, design_file) if {
	not has_vpc_connection(attachment_id, design_file)
}

# Verify if there is an active relationship connecting a TransitGateway to the attachment ID.
has_transit_gateway_connection(attachment_id, design_file) if {
	some rel in design_file.relationships
	rel.status != "deleted"
	some selector in rel.selectors
	some from_clause in selector.allow.from
	from_clause.kind == "TransitGateway"
	some to_clause in selector.allow.to
	to_clause.id == attachment_id
}

# Verify if there is an active relationship connecting a VPC to the attachment ID.
has_vpc_connection(attachment_id, design_file) if {
	some rel in design_file.relationships
	rel.status != "deleted"
	some selector in rel.selectors
	some from_clause in selector.allow.from
	from_clause.kind == "VPC"
	some to_clause in selector.allow.to
	to_clause.id == attachment_id
}

# Cascading deletion: If a relationship is active/approved but is now dangling/incomplete,
# trigger a side effect action to mark its status as "deleted".
relationship_side_effects(relationship, design_file, policy_identifier) := side_effects if {
	policy_identifier == multihop_policy_identifier
	relationship.status == "approved"
	
	some selector in relationship.selectors
	some to_clause in selector.allow.to
	to_clause.kind == "TransitGatewayVPCAttachment"
	attachment_id := to_clause.id
	attachment_id != null
	
	is_incomplete_tgw_vpc_attachment(attachment_id, design_file)
	
	side_effects := {
		{
			"op": actions.update_relationship_op,
			"value": {
				"id": relationship.id,
				"path": "/status",
				"value": "deleted",
			},
		}
	}
} else := {}
