package namespace_integration_test

import rego.v1

import data.relationship_evaluation_policy

# ================================================================================
# SAMPLE DESIGN FIXTURE: Multi-Namespace Kubernetes Design
# ================================================================================
# This test file demonstrates the Namespace-to-Namespace deny relationship policy.
# The Kubernetes hierarchical parent relationship definition explicitly denies
# Namespace components from being parents or children of other Namespace components.
#
# Design structure for testing:
# - 4 Kubernetes Namespace components (ns-1, ns-2, ns-3, ns-4)
# - 1 Kubernetes Deployment component (deploy-1)
# - Expected relationships: Only deploy-1 -> ns-1 (Deployment to Namespace)
# - Denied relationships: Any namespace to namespace relationships
# ================================================================================

# Helper function to create a standard relationship definition for namespace tests
namespace_hierarchical_relationship := {
	"id": "test-rel-id",
	"kind": "hierarchical",
	"type": "parent",
	"subType": "inventory",
	"metadata": {"isAnnotation": false},
	"selectors": [{
		"allow": {
			"from": [{
				"kind": "*",
				"model": {"name": "*", "registrant": {"kind": ""}},
				"patch": {"patchStrategy": "replace", "mutatedRef": [["configuration", "metadata", "namespace"]]},
			}],
			"to": [{
				"kind": "Namespace",
				"model": {"name": "kubernetes", "registrant": {"kind": "github"}},
				"patch": {"patchStrategy": "replace", "mutatorRef": [["displayName"]]},
			}],
		},
		# Deny selector: Prevents Kubernetes Namespace components from having
		# hierarchical parent-child relationships with other Namespace components.
		# This is semantically correct because Namespaces in Kubernetes are flat
		# and do not support nesting - one Namespace cannot be inside another.
		"deny": {
			"from": [{"kind": "Namespace", "model": {"name": "kubernetes", "registrant": {"kind": "github"}}}],
			"to": [{"kind": "Namespace", "model": {"name": "kubernetes", "registrant": {"kind": "github"}}}],
		},
	}],
}

# Test 1: Deployment can have Namespace as parent (should be allowed)
test_deployment_can_have_namespace_as_parent if {
	design_file := {"components": [
		{
			"id": "namespace-1",
			"displayName": "test-namespace",
			"component": {"kind": "Namespace"},
			"model": {"name": "kubernetes", "registrant": {"kind": "github"}},
			"configuration": {"metadata": {"name": "test-namespace"}},
		},
		{
			"id": "deployment-1",
			"displayName": "test-deployment",
			"component": {"kind": "Deployment"},
			"model": {"name": "kubernetes", "registrant": {"kind": "github"}},
			"configuration": {"metadata": {"name": "test-deployment", "namespace": "test-namespace"}},
		},
	]}

	results := relationship_evaluation_policy.identify_relationship(design_file, namespace_hierarchical_relationship)
	count(results) == 1
}

# Test 2: Namespace cannot have Namespace as parent (should be denied)
test_namespace_cannot_have_namespace_as_parent if {
	design_file := {"components": [
		{
			"id": "namespace-1",
			"displayName": "parent-namespace",
			"component": {"kind": "Namespace"},
			"model": {"name": "kubernetes", "registrant": {"kind": "github"}},
			"configuration": {"metadata": {"name": "parent-namespace"}},
		},
		{
			"id": "namespace-2",
			"displayName": "child-namespace",
			"component": {"kind": "Namespace"},
			"model": {"name": "kubernetes", "registrant": {"kind": "github"}},
			"configuration": {"metadata": {"name": "child-namespace", "namespace": "parent-namespace"}},
		},
	]}

	results := relationship_evaluation_policy.identify_relationship(design_file, namespace_hierarchical_relationship)
	count(results) == 0
}

# Test 3: Multiple Namespaces - no parent relationships between them
test_multiple_namespaces_no_relationships if {
	design_file := {"components": [
		{
			"id": "ns-1",
			"displayName": "namespace-1",
			"component": {"kind": "Namespace"},
			"model": {"name": "kubernetes", "registrant": {"kind": "github"}},
			"configuration": {"metadata": {"name": "namespace-1"}},
		},
		{
			"id": "ns-2",
			"displayName": "namespace-2",
			"component": {"kind": "Namespace"},
			"model": {"name": "kubernetes", "registrant": {"kind": "github"}},
			"configuration": {"metadata": {"name": "namespace-2"}},
		},
		{
			"id": "ns-3",
			"displayName": "namespace-3",
			"component": {"kind": "Namespace"},
			"model": {"name": "kubernetes", "registrant": {"kind": "github"}},
			"configuration": {"metadata": {"name": "namespace-3"}},
		},
	]}

	results := relationship_evaluation_policy.identify_relationship(design_file, namespace_hierarchical_relationship)
	count(results) == 0
}

# Test 4: Mixed scenario - only Deployment-to-Namespace relationship should be created
test_mixed_components_only_valid_relationships if {
	design_file := {"components": [
		{
			"id": "ns-1",
			"displayName": "namespace-1",
			"component": {"kind": "Namespace"},
			"model": {"name": "kubernetes", "registrant": {"kind": "github"}},
			"configuration": {"metadata": {"name": "namespace-1"}},
		},
		{
			"id": "ns-2",
			"displayName": "namespace-2",
			"component": {"kind": "Namespace"},
			"model": {"name": "kubernetes", "registrant": {"kind": "github"}},
			"configuration": {"metadata": {"name": "namespace-2"}},
		},
		{
			"id": "deploy-1",
			"displayName": "deployment-1",
			"component": {"kind": "Deployment"},
			"model": {"name": "kubernetes", "registrant": {"kind": "github"}},
			"configuration": {"metadata": {"name": "deployment-1", "namespace": "namespace-1"}},
		},
	]}

	results := relationship_evaluation_policy.identify_relationship(design_file, namespace_hierarchical_relationship)
	count(results) == 1
}

# ================================================================================
# ENHANCED TESTS: Extended registrant objects (as used in real component definitions)
# ================================================================================
# These tests verify the deny logic works with extended registrant objects
# that include additional fields like id, name, status, type - matching the
# actual structure of component declarations from the Kubernetes model.
# ================================================================================

# Test 5: Namespace deny with extended registrant objects (real-world scenario)
test_namespace_deny_with_extended_registrant_objects if {
	design_file := {"components": [
		{
			"id": "ns-1",
			"displayName": "namespace-1",
			"component": {"kind": "Namespace"},
			"model": {
				"name": "kubernetes",
				"registrant": {
					"id": "00000000-0000-0000-0000-000000000001",
					"kind": "github",
					"name": "github",
					"status": "discovered",
					"type": "registry",
				},
			},
			"configuration": {"metadata": {"name": "namespace-1"}},
		},
		{
			"id": "ns-2",
			"displayName": "namespace-2",
			"component": {"kind": "Namespace"},
			"model": {
				"name": "kubernetes",
				"registrant": {
					"id": "00000000-0000-0000-0000-000000000002",
					"kind": "github",
					"name": "github",
					"status": "discovered",
					"type": "registry",
				},
			},
			"configuration": {"metadata": {"name": "namespace-2", "namespace": "namespace-1"}},
		},
	]}

	results := relationship_evaluation_policy.identify_relationship(design_file, namespace_hierarchical_relationship)

	# No namespace-to-namespace relationships should be created
	count(results) == 0
}

# Test 6: Four Namespaces with extended registrant - sample design fixture test
# This test represents the exact scenario from the sample design fixture
test_four_namespaces_sample_design_no_relationships if {
	design_file := {"components": [
		{
			"id": "11111111-1111-1111-1111-111111111111",
			"displayName": "parent-namespace",
			"component": {"kind": "Namespace"},
			"model": {
				"name": "kubernetes",
				"registrant": {
					"id": "00000000-0000-0000-0000-000000000001",
					"kind": "github",
					"name": "github",
					"status": "discovered",
				},
			},
			"configuration": {"metadata": {"name": "parent-namespace"}},
		},
		{
			"id": "22222222-2222-2222-2222-222222222222",
			"displayName": "child-namespace-1",
			"component": {"kind": "Namespace"},
			"model": {
				"name": "kubernetes",
				"registrant": {
					"id": "00000000-0000-0000-0000-000000000002",
					"kind": "github",
					"name": "github",
					"status": "discovered",
				},
			},
			"configuration": {"metadata": {"name": "child-namespace-1", "namespace": "parent-namespace"}},
		},
		{
			"id": "33333333-3333-3333-3333-333333333333",
			"displayName": "child-namespace-2",
			"component": {"kind": "Namespace"},
			"model": {
				"name": "kubernetes",
				"registrant": {
					"id": "00000000-0000-0000-0000-000000000003",
					"kind": "github",
					"name": "github",
					"status": "discovered",
				},
			},
			"configuration": {"metadata": {"name": "child-namespace-2", "namespace": "parent-namespace"}},
		},
		{
			"id": "44444444-4444-4444-4444-444444444444",
			"displayName": "sibling-namespace",
			"component": {"kind": "Namespace"},
			"model": {
				"name": "kubernetes",
				"registrant": {
					"id": "00000000-0000-0000-0000-000000000004",
					"kind": "github",
					"name": "github",
					"status": "discovered",
				},
			},
			"configuration": {"metadata": {"name": "sibling-namespace"}},
		},
	]}

	results := relationship_evaluation_policy.identify_relationship(design_file, namespace_hierarchical_relationship)

	# No namespace-to-namespace relationships should be created
	count(results) == 0
}

# Test 7: Sample design fixture - four namespaces + one deployment
# Only the Deployment-to-Namespace relationship should be created
test_sample_design_fixture_deployment_to_namespace_only if {
	design_file := {"components": [
		{
			"id": "11111111-1111-1111-1111-111111111111",
			"displayName": "parent-namespace",
			"component": {"kind": "Namespace"},
			"model": {
				"name": "kubernetes",
				"registrant": {
					"id": "00000000-0000-0000-0000-000000000001",
					"kind": "github",
					"status": "discovered",
				},
			},
			"configuration": {"metadata": {"name": "parent-namespace"}},
		},
		{
			"id": "22222222-2222-2222-2222-222222222222",
			"displayName": "child-namespace-1",
			"component": {"kind": "Namespace"},
			"model": {
				"name": "kubernetes",
				"registrant": {
					"id": "00000000-0000-0000-0000-000000000002",
					"kind": "github",
					"status": "discovered",
				},
			},
			"configuration": {"metadata": {"name": "child-namespace-1", "namespace": "parent-namespace"}},
		},
		{
			"id": "33333333-3333-3333-3333-333333333333",
			"displayName": "child-namespace-2",
			"component": {"kind": "Namespace"},
			"model": {
				"name": "kubernetes",
				"registrant": {
					"id": "00000000-0000-0000-0000-000000000003",
					"kind": "github",
					"status": "discovered",
				},
			},
			"configuration": {"metadata": {"name": "child-namespace-2", "namespace": "parent-namespace"}},
		},
		{
			"id": "44444444-4444-4444-4444-444444444444",
			"displayName": "sibling-namespace",
			"component": {"kind": "Namespace"},
			"model": {
				"name": "kubernetes",
				"registrant": {
					"id": "00000000-0000-0000-0000-000000000004",
					"kind": "github",
					"status": "discovered",
				},
			},
			"configuration": {"metadata": {"name": "sibling-namespace"}},
		},
		{
			"id": "55555555-5555-5555-5555-555555555555",
			"displayName": "test-deployment",
			"component": {"kind": "Deployment"},
			"model": {
				"name": "kubernetes",
				"registrant": {
					"id": "00000000-0000-0000-0000-000000000005",
					"kind": "github",
					"status": "discovered",
				},
			},
			"configuration": {"metadata": {"name": "test-deployment", "namespace": "parent-namespace"}},
		},
	]}

	results := relationship_evaluation_policy.identify_relationship(design_file, namespace_hierarchical_relationship)

	# Only ONE relationship should be created: Deployment -> Namespace
	count(results) == 1

	# Verify the relationship connects a non-Namespace component to a Namespace
	# The 'from' selector will keep the original wildcard kind "*" from the relationship definition
	# while 'to' will have kind "Namespace". We verify by checking the IDs reference
	# components from our design (deployment and namespace)
	some result in results
	from_id := result.selectors[0].allow.from[0].id
	to_id := result.selectors[0].allow.to[0].id

	# Verify from component is the deployment (not a namespace)
	some from_comp in design_file.components
	from_comp.id == from_id
	from_comp.component.kind == "Deployment"

	# Verify to component is a namespace
	some to_comp in design_file.components
	to_comp.id == to_id
	to_comp.component.kind == "Namespace"
}

# Test 8: Verify is_relationship_denied with extended registrant structure
test_is_relationship_denied_extended_registrant if {
	from_declaration := {
		"id": "from-id",
		"component": {"kind": "Namespace"},
		"model": {
			"name": "kubernetes",
			"registrant": {
				"id": "00000000-0000-0000-0000-000000000001",
				"kind": "github",
				"name": "github",
				"status": "discovered",
				"type": "registry",
			},
		},
	}

	to_declaration := {
		"id": "to-id",
		"component": {"kind": "Namespace"},
		"model": {
			"name": "kubernetes",
			"registrant": {
				"id": "00000000-0000-0000-0000-000000000002",
				"kind": "github",
				"name": "github",
				"status": "discovered",
				"type": "registry",
			},
		},
	}

	deny_selectors := {
		"from": [{
			"kind": "Namespace",
			"model": {
				"name": "kubernetes",
				"registrant": {"kind": "github"},
			},
		}],
		"to": [{
			"kind": "Namespace",
			"model": {
				"name": "kubernetes",
				"registrant": {"kind": "github"},
			},
		}],
	}

	# This should return true (relationship IS denied)
	relationship_evaluation_policy.is_relationship_denied(from_declaration, to_declaration, deny_selectors)
}

# Test 9: Verify is_relationship_denied allows non-matching kinds
test_is_relationship_not_denied_different_kinds if {
	from_declaration := {
		"id": "from-id",
		"component": {"kind": "Deployment"},
		"model": {
			"name": "kubernetes",
			"registrant": {
				"id": "00000000-0000-0000-0000-000000000001",
				"kind": "github",
				"status": "discovered",
			},
		},
	}

	to_declaration := {
		"id": "to-id",
		"component": {"kind": "Namespace"},
		"model": {
			"name": "kubernetes",
			"registrant": {
				"id": "00000000-0000-0000-0000-000000000002",
				"kind": "github",
				"status": "discovered",
			},
		},
	}

	deny_selectors := {
		"from": [{
			"kind": "Namespace",
			"model": {
				"name": "kubernetes",
				"registrant": {"kind": "github"},
			},
		}],
		"to": [{
			"kind": "Namespace",
			"model": {
				"name": "kubernetes",
				"registrant": {"kind": "github"},
			},
		}],
	}

	# This should return false (relationship is NOT denied - Deployment to Namespace is allowed)
	not relationship_evaluation_policy.is_relationship_denied(from_declaration, to_declaration, deny_selectors)
}

# Test 10: Multiple deployments in different namespaces - all should have relationships
test_multiple_deployments_in_namespaces if {
	design_file := {"components": [
		{
			"id": "ns-1",
			"displayName": "namespace-1",
			"component": {"kind": "Namespace"},
			"model": {"name": "kubernetes", "registrant": {"kind": "github"}},
			"configuration": {"metadata": {"name": "namespace-1"}},
		},
		{
			"id": "ns-2",
			"displayName": "namespace-2",
			"component": {"kind": "Namespace"},
			"model": {"name": "kubernetes", "registrant": {"kind": "github"}},
			"configuration": {"metadata": {"name": "namespace-2"}},
		},
		{
			"id": "deploy-1",
			"displayName": "deployment-1",
			"component": {"kind": "Deployment"},
			"model": {"name": "kubernetes", "registrant": {"kind": "github"}},
			"configuration": {"metadata": {"name": "deployment-1", "namespace": "namespace-1"}},
		},
		{
			"id": "deploy-2",
			"displayName": "deployment-2",
			"component": {"kind": "Deployment"},
			"model": {"name": "kubernetes", "registrant": {"kind": "github"}},
			"configuration": {"metadata": {"name": "deployment-2", "namespace": "namespace-2"}},
		},
	]}

	results := relationship_evaluation_policy.identify_relationship(design_file, namespace_hierarchical_relationship)

	# Two relationships: deploy-1 -> ns-1, deploy-2 -> ns-2
	count(results) == 2
}
