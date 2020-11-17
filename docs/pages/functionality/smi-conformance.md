---
layout: default
title: SMI Conformance Capability
permalink: functionality/smi-conformance
type: functionality
---

* toc
{:toc}

### **Meshery + SMI**

Meshery encaptures the SMI spec by validating your service mesh's conformance to the SMI specifications. We achieve this by running through a service mesh's workloads without interacting with the specific mesh's APIs by using a sample application to provide a consistent and constant workload. This testing practice is followed because we, at Meshery, understand your meshing needs and acknowledge that different workloads (applications) and the varied types and sizes of infrastructure resources pose a substantial need for a cross-mesh comparison between service meshes so that users can predict and understand distinct behavioural differences, and run performance tests to validate their mesh's abilities.

This allows Meshery to provide an independent, unbiased and credible analysis in the form of a compatibility matrix, identifying the SMI features supported by each service mesh and the tabulated results of a suite of repeatable conformance tests. This is achieved by using a sample application to produce a consistent workload by deploying a pre-configured sample application [**Learn Layer5**](https://github.com/layer5io/learn-layer5).

#### **Features**

**User** - For users, Meshery allows tests to be scheduled and invoked programmatically. Meshery will also store these test results and allow you to retrieve them later.

**Provider** - For providers, Meshery guarantees provenance of these tests and facilitates the public publicing of this suite of tests results into a versioned, public matrix of conformance status (consisting of both supported capabilities and test compliance).

### **Testing the capability of your mesh**

Keeping par with SMI's ability to differentiate between full implementation of a specification and compliance with the portions that it implements, Meshery designates three possible capabilities to a service mesh:

- **Full** - The service mesh has this capability.

- **Partial** - The service mesh has a portion of this capability implemented (may or may not have this full capability in the future).

- **None** - The service mesh does not have this capability (may or may not have this capability in the future).

### **Approach to Conformance Testing**

Meshery has a unique approach to testing, involving a suite of tests which will be carried out for each:

- Kubernetes version 
- Service Mesh version 
- SMI category

Conformance tests are categorized by SMI specification type. Each test will be carried out in an automated and concurrent fashion by invoking Meshery to execute the conformance tests. Within a set of defined tests for each SMI specification, two types of assertion tests have been defined:

- Presence assertion  
- Capability assertion

##### **Steps**

- [Install Meshery]({{ site.baseurl }}/installation/quick-start)
- [Setup a specific version of a service mesh]({{ site.baseurl }}/service-meshes)
- (**Optional**) [Setup Prometheus accordingly]({{ site.baseurl }}/guides/meshery-metrics)
- [Deploy a chosen sample app]({{ site.baseurl }}/guides/sample-apps)
- Deploy the needed SMI operator like smi-metrics, the test to run, which can involve calling an API endpoint
- Validation of the response

**Note**: All tests cannot be validated by using responses. In those cases, we can also connect to a Prometheus instance, which is preconfigured to collect all the metrics for individual tests, and use these metrics to define expectations. In this guide, we will be working on specific examples.
Finally, the results for the test runs from Meshery are persisted (same GitHub repository) and published on the conformance web page. 

### **Executing Performance Tests**

##### Preconditions

- A given service mesh’s ability to adhere to the SMI specification is validated by running a workload on top of the service mesh.
- Workload deployments are configured specific to the onboarding requirements of the given service mesh.
- Tests are defined to validate conformance for each type of SMI specification (e.g. metrics, access, traffic).

##### Invocation

- Test assertions are defined in a workload-specific way and deployed with the workloads being tested (test assertions are packaged).
- A test result is collected with the evaluation of each assertion.
- Future Test results will be individually streamed to Meshery after each assertion is evaluated.
- Once all assertions are evaluated, test results are returned and can be viewed on the Meshery dashboard.

### **Reporting Conformance**

##### Provenance of Test Results

Each participating service mesh project will be asked to incorporate the conformance tool, Meshery, into their CI pipelines, or alternatively, to manually run the conformance test suite when a release of the service mesh project is made. The conformance tool will run the test suite and automatically update the conformance dashboard when the test is complete.
In order to ensure provenance of test results that represent a given service mesh, each project will be asked to identify one or more github accounts that will be used for publishing the tests. Ideally, this github account is a *servicerobot* account used within the project’s CI pipeline.

This method of providing verification of results is similarly used for the service mesh projects that also use Meshery to provide their performance test results.

### **Publishing Test Results**

1. A public-facing report will display the current and historical status of individual service mesh capabilities with each of the SMI specifications.
1.  The report will be visual in nature, but also be available as a yaml file. The dashboard will be published publicly on our [SMI Conformance](https/meshery.io/smi-conformance) page. It is currently listed on the [Landscape](https/layer5.io/landscape#smi) page
1. Reports tracking the history of service mesh versions, SMI spec versions and their compatibility will be stored.