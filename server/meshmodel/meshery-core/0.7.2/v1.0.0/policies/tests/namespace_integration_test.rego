package namespace_integration_test

import rego.v1

import data.relationship_evaluation_policy

# Test 1: Deployment can have Namespace as parent (should be allowed)
test_deployment_can_have_namespace_as_parent if {
design_file := {
"components": [
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
],
}

relationship := {
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
"deny": {
"from": [{"kind": "Namespace", "model": {"name": "kubernetes", "registrant": {"kind": "github"}}}],
"to": [{"kind": "Namespace", "model": {"name": "kubernetes", "registrant": {"kind": "github"}}}],
},
}],
}

results := relationship_evaluation_policy.identify_relationship(design_file, relationship)
count(results) == 1
}

# Test 2: Namespace cannot have Namespace as parent (should be denied)
test_namespace_cannot_have_namespace_as_parent if {
design_file := {
"components": [
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
],
}

relationship := {
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
"deny": {
"from": [{"kind": "Namespace", "model": {"name": "kubernetes", "registrant": {"kind": "github"}}}],
"to": [{"kind": "Namespace", "model": {"name": "kubernetes", "registrant": {"kind": "github"}}}],
},
}],
}

results := relationship_evaluation_policy.identify_relationship(design_file, relationship)
count(results) == 0
}

# Test 3: Multiple Namespaces - no parent relationships between them
test_multiple_namespaces_no_relationships if {
design_file := {
"components": [
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
],
}

relationship := {
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
"deny": {
"from": [{"kind": "Namespace", "model": {"name": "kubernetes", "registrant": {"kind": "github"}}}],
"to": [{"kind": "Namespace", "model": {"name": "kubernetes", "registrant": {"kind": "github"}}}],
},
}],
}

results := relationship_evaluation_policy.identify_relationship(design_file, relationship)
count(results) == 0
}

# Test 4: Mixed scenario - only Deployment-to-Namespace relationship should be created
test_mixed_components_only_valid_relationships if {
design_file := {
"components": [
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
],
}

relationship := {
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
"deny": {
"from": [{"kind": "Namespace", "model": {"name": "kubernetes", "registrant": {"kind": "github"}}}],
"to": [{"kind": "Namespace", "model": {"name": "kubernetes", "registrant": {"kind": "github"}}}],
},
}],
}

results := relationship_evaluation_policy.identify_relationship(design_file, relationship)
count(results) == 1
}
