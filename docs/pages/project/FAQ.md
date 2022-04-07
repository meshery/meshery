---
layout: page
title: FAQ
permalink: project/faq
description: General commonly asked questions and answers about Meshery.
language: en
type: project
---

#### Question: What is service mesh?
**Answer:** _A service mesh is a way to control how different parts of an application share data with one another. Unlike other systems for managing this communication, a service mesh is a dedicated infrastructure layer built right into an app. This visible infrastructure layer can document how well (or not) different parts of an app interact, so it becomes easier to optimize communication and avoid downtime as an app grows._ 
#### Question: What is Meshery?
**Answer:** _Meshery is the open source, service mesh management plane that enables the adoption, operation, and management of any service mesh and their workloads._

#### Question: Why was Meshery created?

 **Answer:** _As an open source, vendor neutral project, Meshery was created out of the necessity to enable service mesh adopters to overcome the challenge of complex virtual networking; to come to expect more from their infrastructure; to enable the world to understand and operate any service mesh with confidence._

#### Question: What does Meshery do?
**Answer:** 
 - _Operates service meshes with confidence,_
   - _with operational best practices._ 
 - _Compare apples-to-apples performance across service meshes._
 - _Understand behavioral differences between service meshes._
 - _Track your application performance from version to version._

#### Question: What is mesheryctl?

**Answer:** _A command line interface to manage Meshery._

#### Question: How to install Meshery?

**Answer:** _Meshery runs on a number of platforms._
_Popular Installers:_
- _Bash user → Meshery Installation [Quick Start](https://docs.meshery.io/installation/quick-start)._
- _[Brew user](https://github.com/meshery/homebrew-tap)._
- _Kubernetes user._
- _[Scoop user](https://github.com/meshery/scoop-bucket)._ 


#### Question: What is the difference between `make server` and `mesheryctl system start`? Do they both run Meshery on my local machine?

  **Answer:** _Yes, both of them do run Meshery on your local machine. `make server` builds Meshery from source and runs it on your local OS, while `mesheryctl system start` runs Meshery as a set of containers in Docker or in Kubernetes on your local machine._ 

#### Question: What systems is Meshery compatible with?

  **Answer:** _Many. See Meshery's [Compatibility Matrix]({{site.baseurl}}/project/compatibility-matrix)._

{% include discuss.html %}

<!--Add other questions-->

