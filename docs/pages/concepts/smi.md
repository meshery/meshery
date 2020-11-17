---
layout: default
title: SMI specification
permalink: concepts/smi-spec
type: concepts
---

#### **Service Mesh Interface (SMI)**

The [SMI specification](https://smi-spec.io/) is a standard interface for service meshes on Kubernetes. It aims to define a common base standard for enabling users and developers to compare and contrast service mesh technologies by using a set of pre-defined and consistent APIs, thereby facilitating a comparison of behavioral differences between service meshes. SMI currently supports a set of four API's:

- **Traffic Specs** - Used to define traffic, currently only supports TCP, HTTP traffic.
- **Traffic Access Control** - Used to specify whether a particular form of traffic is allowed or not
- **Traffic Split** - Used to redirect/divide a request for a resource between 2 or more resources. Useful in canary testing
- **Traffic Metrics** - Used to expose common traffic metrics like p99 in a specific format that can be utilized by single dashboard for all the service meshes.

##### **How does SMI work?**

SMI delivers its API's with the intention of allowing providers to define its implementations as needed.This allows:
1. The direct usage of SMI API's
2. The usage of operators which can implement SMI to native API's present in their environment

SMI also acknowledges that Conformance to any specification is highly subjective and dependent on the existing infrastructure and specific workloads. To establish a level playing field, SMI recognizes that some participating service meshes may conscientiously never fully implement functions (based on SMI specs) and identifies the difference between full implementation of a specification and compliance with the portions that it implements.

Multiple mesh providers have implemented SMI compatibility as an industry standard by:
- Using an adapter that defines SMI CRDs as a wrapper over the native APIs ( <img src="{{ site.baseurl }}/assets/img/service-meshes/consul.svg" style="width:20px" />, <img src="{{ site.baseurl }}/assets/img/service-meshes/istio.svg" style="width:15px" />)

- Directly implementing and adhering to the SMI spec APIs internally ( <img src="{{ site.baseurl }}/assets/img/service-meshes/linkerd.svg" style="width:20px" />, <img src="{{ site.baseurl }}/assets/img/service-meshes/osm.svg" style="width:20px" />, <img src="{{ site.baseurl }}/assets/img/service-meshes/maesh.svg" style="width:15px" />)

**Meshery utilizes SMI in its capacity as an industry standard to provide users with an independent, unbiased, and credible analysis so that they can make an informed choice.**

#### **Defining Conformance with Meshery**

Conformance with SMI specifications is defined as a series of test assertions. A test assertion is a condition that must be tested to confirm conformance to a requirement. A test assertion is a condition that from the perspective of validation testing, determines whether the conformance will require any number of conditions to be tested. The collection of test assertions categorized by SMI specification collectively defines the suite of SMI conformance tests.

**Meshery is the test harness used to fit SMI conformance tests to different service meshes and different workloads.**

Meshery validates your service mesh's conformance to the SMI specifications through automated provisioning of individual service meshes and deployment of a consistent workload produced by deploying [**Learn Layer5**](https://github.com/layer5io/learn-layer5) as a sample application. 

##### **Suggested Reading**

Read about how Meshery works with the SMI spec in our guide on [SMI Conformance]({{ site.baseurl }}/functionality/smi-conformance)