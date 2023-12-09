<!-- ---
layout: default
title: Meshery and Service Mesh Interface (SMI)
permalink: tasks/service-mesh-interface
type: tasks
language: en
---
{% include alert.html type="info" title="Validating SMI Conformance with Meshery" content="<p>Meshery is SMI's 
<a href='https://smi-spec.io/blog/validating-smi-conformance-with-meshery'>official tool for validating conformance</a>.</p> <p>The <a href='https://meshery.io/service-mesh-interface'>SMI Conformance dashboard</a> displays the current and historical test results of each of the service mesh projects.</p>" %}

As service mesh providers and the surrounding ecosystem start adopting Service Mesh Interface, there is a growing need for validating the SMI implementation.

As a service mesh agnostic tool that provides lifecycle and performance management of a large number of service meshes, sample applications, service mesh patterns and WebAssembly filters, Meshery is the ideal tool for the job when it comes to SMI Conformance.

## Defining Conformance

Conformance of SMI specifications is defined as a series of test assertions. A test assertion is a condition that must be positively verified in order for an implementation to be considered conformant. A test assertion may involve any number of conditions. Sets of test assertions are categorized by SMI specification. Collectively, these test sets comprise the complete suite of SMI conformance tests. Just like SMI itself, the suite of SMI Conformance tests is versioned, and with each new version of SMI, as interfaces are added and specifications changed, the Conformance requirements will change as appropriate.

## Validating Conformance

Conformance validation is performed through automated provisioning of individual service meshes, deployment of a common workload, and generation of service request load as necessary. To facilitate a common set of tests, a simple, instrumented, sample application has been developed for purposes of providing a consistent workload to apply SMI specs against. Deployment of the sample application has been fitted to each service mesh.

## Capability

Given that some service mesh implementations, never intend to fully implement the SMI specifications, for each individual test, three possible capability designations exist

- Full - service mesh has this capability.

- Partial - service mesh has a portion of this capability implemented (may or may not have this full capability in the future).

- None - currently service mesh does not have this capability (may or may not have this capability in the future).

## Steps in Running Conformance Tests

The following list highlights the sequence of steps taken to perform conformance testing of one type of service mesh.

* Preconditions
   * A given service mesh’s ability to adhere to the SMI specification is validated by running a workload on top of the service mesh.
   * Workload deployments are configured specific to the onboarding requirements of the given service mesh.
   * Tests are defined to validate conformance for each type of SMI specification (e.g. metrics, access, traffic… ).
* Invocation
   * Test assertions are defined in a workload-specific way and deployed with the workloads being tested (test assertions are packaged).
   * A test result is collected with the evaluation of each assertion. Future: Test results will be individually streamed to Meshery after each assertion is evaluated.
   * Once all assertions are evaluated, test results are returned for visual presentation in Meshery.
   * Tests that are run by the service mesh project team are identified by a pre-approved service account and are publicly published the SMI Conformance dashboard.


## Reporting Conformance

Each of the participating service mesh projects can leverage the [Meshery SMI Conformance GitHub action](https://github.com/layer5io/meshery-smi-conformance-action) to incorporate SMI conformance into their CI/CD pipelines.

The conformance tool will run the test suite and automatically update the conformance dashboard when the test is complete. In order to ensure provenance of test results that represent a given service mesh, each project will be asked to identify one or more GitHub accounts that will be used for publishing the tests. 

Ideally, this GitHub account is a service-robot account used within the project’s CI pipeline.

The [SMI Conformance dashboard](https://meshery.io/service-mesh-interface) displays the current and historical test results of each of the service mesh projects. These results are automatically updated through the projects official accounts as mentioned above. -->