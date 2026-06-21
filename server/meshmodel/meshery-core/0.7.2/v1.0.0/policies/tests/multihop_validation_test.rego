package multihop_validation_test

import rego.v1

import data.eval

# Test 1: Complete TGW VPC attachment hop (both TGW and VPC connections are present).
# The relationship should be valid (relationship_is_invalid should be false).
test_multihop_tgw_vpc_attachment_complete if {
	design_file := {
		"components": [
			{
				"id": "tgw-1",
				"component": {"kind": "TransitGateway"},
				"model": {"name": "aws-ec2-controller", "registrant": {"kind": "github"}},
				"configuration": {"metadata": {"name": "my-tgw"}},
			},
			{
				"id": "vpc-1",
				"component": {"kind": "VPC"},
				"model": {"name": "aws-ec2-controller", "registrant": {"kind": "github"}},
				"configuration": {"metadata": {"name": "my-vpc"}},
			},
			{
				"id": "tgw-attachment-1",
				"component": {"kind": "TransitGatewayVPCAttachment"},
				"model": {"name": "aws-ec2-controller", "registrant": {"kind": "github"}},
				"configuration": {
					"metadata": {"name": "my-attachment"},
					"spec": {
						"transitGatewayRef": {"from": {"name": "my-tgw"}},
						"vpcRef": {"from": {"name": "my-vpc"}},
					},
				},
			},
		],
		"relationships": [
			{
				"id": "rel-tgw-to-att",
				"kind": "edge",
				"type": "binding",
				"subType": "network",
				"status": "approved",
				"selectors": [{
					"allow": {
						"from": [{"id": "tgw-1", "kind": "TransitGateway"}],
						"to": [{"id": "tgw-attachment-1", "kind": "TransitGatewayVPCAttachment"}],
					},
				}],
			},
			{
				"id": "rel-vpc-to-att",
				"kind": "edge",
				"type": "binding",
				"subType": "network",
				"status": "approved",
				"selectors": [{
					"allow": {
						"from": [{"id": "vpc-1", "kind": "VPC"}],
						"to": [{"id": "tgw-attachment-1", "kind": "TransitGatewayVPCAttachment"}],
					},
				}],
			},
		],
	}

	# Ensure it implicates the relationship
	eval.relationship_is_implicated_by_policy(design_file.relationships[0], eval.multihop_policy_identifier)

	# The relationship should NOT be marked invalid since both hops exist
	not eval.relationship_is_invalid(design_file.relationships[0], design_file, eval.multihop_policy_identifier)
}

# Test 2: Incomplete TGW VPC attachment (missing the VPC connection).
# The relationship should be invalid.
test_multihop_tgw_vpc_attachment_missing_vpc if {
	design_file := {
		"components": [
			{
				"id": "tgw-1",
				"component": {"kind": "TransitGateway"},
				"model": {"name": "aws-ec2-controller", "registrant": {"kind": "github"}},
				"configuration": {"metadata": {"name": "my-tgw"}},
			},
			{
				"id": "tgw-attachment-1",
				"component": {"kind": "TransitGatewayVPCAttachment"},
				"model": {"name": "aws-ec2-controller", "registrant": {"kind": "github"}},
				"configuration": {
					"metadata": {"name": "my-attachment"},
					"spec": {
						"transitGatewayRef": {"from": {"name": "my-tgw"}},
					},
				},
			},
		],
		"relationships": [
			{
				"id": "rel-tgw-to-att",
				"kind": "edge",
				"type": "binding",
				"subType": "network",
				"status": "approved",
				"selectors": [{
					"allow": {
						"from": [{"id": "tgw-1", "kind": "TransitGateway"}],
						"to": [{"id": "tgw-attachment-1", "kind": "TransitGatewayVPCAttachment"}],
					},
				}],
			},
		],
	}

	# The relationship should be invalid because VPC relationship is missing
	eval.relationship_is_invalid(design_file.relationships[0], design_file, eval.multihop_policy_identifier)
}

# Test 3: Incomplete TGW VPC attachment (missing the TGW connection).
# The relationship should be invalid.
test_multihop_tgw_vpc_attachment_missing_tgw if {
	design_file := {
		"components": [
			{
				"id": "vpc-1",
				"component": {"kind": "VPC"},
				"model": {"name": "aws-ec2-controller", "registrant": {"kind": "github"}},
				"configuration": {"metadata": {"name": "my-vpc"}},
			},
			{
				"id": "tgw-attachment-1",
				"component": {"kind": "TransitGatewayVPCAttachment"},
				"model": {"name": "aws-ec2-controller", "registrant": {"kind": "github"}},
				"configuration": {
					"metadata": {"name": "my-attachment"},
					"spec": {
						"vpcRef": {"from": {"name": "my-vpc"}},
					},
				},
			},
		],
		"relationships": [
			{
				"id": "rel-vpc-to-att",
				"kind": "edge",
				"type": "binding",
				"subType": "network",
				"status": "approved",
				"selectors": [{
					"allow": {
						"from": [{"id": "vpc-1", "kind": "VPC"}],
						"to": [{"id": "tgw-attachment-1", "kind": "TransitGatewayVPCAttachment"}],
					},
				}],
			},
		],
	}

	# The relationship should be invalid because TransitGateway relationship is missing
	eval.relationship_is_invalid(design_file.relationships[0], design_file, eval.multihop_policy_identifier)
}
