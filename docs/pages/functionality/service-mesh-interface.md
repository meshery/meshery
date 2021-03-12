---
layout: default
title: Meshery and Service Mesh Interface (SMI)
permalink: functionality/service-mesh-interface
type: functionality
---

## Conformance

Defining “Conformance” - It’s important to acknowledge that conformance consists of both capabilities and compliance status. We define conformance as a combination of these two concepts.

1. SMI Conformance acknowledges that
   ...some participating service meshes may conscientiously never fully implement functions (SMI specs)...

2. SMI Conformance identifies
   ...a difference between full implementation of a specification and compliance with the portions that it implements...

## Capability

Given that some service mesh implementations, never intend to fully implement the SMI specifications, for each individual test, three possible capability designations exist

- Full - service mesh has this capability.

- Partial - service mesh has a portion of this capability implemented (may or may not have this full capability in the future).

- None - currently service mesh does not have this capability (may or may not have this capability in the future).

## Approach to Conformance Testing

Each Kubernetes version, service mesh version and SMI category will undergo several tests. Each test will be carried out in an automated and concurrent fashion, mostly, by invoking Meshery to execute the conformance tests.

#### Here are the steps

- Setup a specific version of a service mesh
- [Optional] Setup Prometheus accordingly
- Deploy a chosen sample app
- Deploy the needed SMI operator like smi-metrics
  the test to run, which can involve calling an API endpoint
- Validation of the response

Not all tests can be validated by just using the response, in those cases we can also connect to a prometheus instance, which is preconfigured to collect all the metrics for the test, and use these metrics to define expectations. We will be working on specific examples.

Finally, the results for the test runs from Meshery are persisted (same GitHub repository) and published on the conformance web page. Eventually, we can build a system which will allow us to run granular tests on demand for a chosen Kubernetes, service mesh and SMI Operator versions.

## Conformance Test Definitions

Conformance tests are categorized by SMI specification type. A set of tests are defined for each SMI specification. Within each test set, two types of assertion tests are defined a presence assertion and a capability assertion.

## Validating Conformance

Conformance to SMI specifications will be done through automated provisioning of individual service meshes and deployment of a common workload. A simple, instrumented, sample application is used as the workload to test.

## Defining Conformance

Conformance with SMI specifications is defined as a series of test assertions. A test assertion is a condition that must be tested to confirm conformance to a requirement. A test assertion is a condition that from the perspective of validation testing, determining conformance will require any number of conditions to be tested. The collection of test assertions categorized by SMI specification collectively define the suite of SMI conformance tests. Meshery is the test harness used to fit SMI conformance tests to different service meshes and different workloads.

## Steps To Executing Performance Tests

### Preconditions

- A given service mesh’s ability to adhere to the SMI specification is validated by running a workload on top of the service mesh.
- Workload deployments are configured specific to the onboarding requirements of the given service mesh.
- Tests are defined to validate conformance for each type of SMI specification (e.g. metrics, access, traffic… ).

### Invocation

- Test assertions are defined in a workload-specific way and deployed with the workloads being tested (test assertions are packaged).
- A test result is collected with the evaluation of each assertion.
- Future Test results will be individually streamed to Meshery after each assertion is evaluated.
- Once all assertions are evaluated, test results are returned for visual presentation in Meshery

## Reporting Conformance

### Provenance of Test Results

Each participating service mesh project will be asked to incorporate the conformance tool, Meshery, into their CI pipelines, or alternatively, to manually run the conformance test suite when a release of the service mesh project is made. The conformance tool will run the test suite and automatically update the conformance dashboard when the test is complete.
In order to ensure provenance of test results that represent a given service mesh, each project will be asked to identify one or more github accounts that will be used for publishing the tests. Ideally, this github account is a servicerobot account used within the project’s CI pipeline.

This method of providing verification of results is similarly used for those same service mesh projects that also use Meshery to provide their performance test results.

Each service mesh project needs to identify their servicerobot account, by updating this list httpsmeshery.iosmi-conformancesm-service-accounts. Identify a given Meshery Provider user and designate their “CI service account”.

## Publishing Test Results

A public-facing report will display the current and historical status of individual service mesh capability with each of the SMI specifications. The report will be visual in nature, but also be available as yaml. Dashboard to be published publicly here httpsmeshery.iosmi-conformance (currently, listed on httpslayer5.iolandscape#smi).

### Historical Conformance

Reports will track the history of service mesh versions and SMI spec versions and their compatibility.
