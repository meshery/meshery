---
layout: default
title: SMI Conformance Capability
permalink: functionality/smi-conformance
type: functionality
---

## Meshery + SMI

Meshery encaptures the SMI spec and practices conformance through a service mesh's workloads without interacting with the specific mesh's APIs.

We understand that different workloads (applications) and the varied types and sizes of infrastructure resources pose a substantial need for a cross-mesh comparison between service meshes so that users can predict and understand distinct behavioural differences, and run performance tests to validate their mesh's abilities.

This allows Meshery to provide an independent, unbiased, and credible analysis in the form of a compatibility matrix, identifying the SMI features supported by each service mesh and the tabulated results of a suite of repeatable conformance tests. This is achieved by using a sample application to produce a consistent workload by deploying a pre-configured sample application [`Learn Layer5`](https://github.com/layer5io/learn-layer5).

#### Testing the capability of your mesh

Keeping par with SMI's ability to differentiate between full implementation of a specification and compliance with the portions that it implements, Meshery designates three possible capabilities to a service mesh:

- **Full** - The service mesh has this capability.

- **Partial** - The service mesh has a portion of this capability implemented (may or may not have this full capability in the future).

- **None** - The service mesh does not have this capability (may or may not have this capability in the future).

#### Approach to Conformance Testing

Meshery has a unique approach to testing, involving a suite of tests which will be carried out for each:

- Kubernetes version 
- Service Mesh version 
- SMI category

Each test will be carried out in an automated and concurrent fashion by invoking Meshery to execute the conformance tests.

##### Here are the steps

- [Install Meshery]({{ site.baseurl }}/installation/quick-start)
- [Setup a specific version of a service mesh]({{ site.baseurl }}/service-meshes)
- (**Optional**) [Setup Prometheus accordingly]({{ site.baseurl }}/guides/meshery-metrics)
- [Deploy a chosen sample app]({{ site.baseurl }}/guides/sample-apps)
- Deploy the needed SMI operator like smi-metrics, the test to run, which can involve calling an API endpoint
- Validation of the response

Not all tests can be validated by just using the response, in those cases we can also connect to a prometheus instance, which is preconfigured to collect all the metrics for the test, and use these metrics to define expectations. In this guide, we will be working on specific examples.
Finally, the results for the test runs from Meshery are persisted (same GitHub repository) and published on the conformance web page. 

#### Conformance Test Definitions

Conformance tests are categorized by SMI specification type. A set of tests are defined for each SMI specification. Within each test set, two types of assertion tests are defined a presence assertion and a capability assertion.

#### Validating Conformance

Conformance to SMI specifications will be done through automated provisioning of individual service meshes and deployment of a common workload. A simple, instrumented, sample application is used as the workload to test.

#### Defining Conformance

Conformance with SMI specifications is defined as a series of test assertions. A test assertion is a condition that must be tested to confirm conformance to a requirement. A test assertion is a condition that from the perspective of validation testing, determining conformance will require any number of conditions to be tested. The collection of test assertions categorized by SMI specification collectively define the suite of SMI conformance tests. Meshery is the test harness used to fit SMI conformance tests to different service meshes and different workloads.

### Steps To Executing Performance Tests

#### Preconditions

- A given service mesh’s ability to adhere to the SMI specification is validated by running a workload on top of the service mesh.
- Workload deployments are configured specific to the onboarding requirements of the given service mesh.
- Tests are defined to validate conformance for each type of SMI specification (e.g. metrics, access, traffic… ).

#### Invocation

- Test assertions are defined in a workload-specific way and deployed with the workloads being tested (test assertions are packaged).
- A test result is collected with the evaluation of each assertion.
- Future Test results will be individually streamed to Meshery after each assertion is evaluated.
- Once all assertions are evaluated, test results are returned for visual presentation in Meshery

### Reporting Conformance

#### Provenance of Test Results

Each participating service mesh project will be asked to incorporate the conformance tool, Meshery, into their CI pipelines, or alternatively, to manually run the conformance test suite when a release of the service mesh project is made. The conformance tool will run the test suite and automatically update the conformance dashboard when the test is complete.
In order to ensure provenance of test results that represent a given service mesh, each project will be asked to identify one or more github accounts that will be used for publishing the tests. Ideally, this github account is a servicerobot account used within the project’s CI pipeline.

This method of providing verification of results is similarly used for those same service mesh projects that also use Meshery to provide their performance test results.

Each service mesh project needs to identify their servicerobot account, by updating this list httpsmeshery.iosmi-conformancesm-service-accounts. Identify a given Meshery Provider user and designate their “CI service account”.

### Publishing Test Results

A public-facing report will display the current and historical status of individual service mesh capability with each of the SMI specifications. The report will be visual in nature, but also be available as yaml. Dashboard to be published publicly here httpsmeshery.iosmi-conformance (currently, listed on httpslayer5.iolandscape#smi).

### Historical Conformance

Reports will track the history of service mesh versions and SMI spec versions and their compatibility.
