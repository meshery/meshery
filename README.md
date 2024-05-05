<p style="text-align:center;" align="center">
  <a href="https://meshery.io">
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/meshery/meshery/master/.github/assets/images/readme/meshery-logo-light-text-side.svg">
      <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/meshery/meshery/master/.github/assets/images/readme/meshery-logo-dark-text-side.svg">
      <img src="https://raw.githubusercontent.com/meshery/meshery/master/.github/assets/images/readme/meshery-logo-dark-text-side.svg" alt="Meshery Logo" width="70%" />
    </picture>
  </a>
  <br/>
</p>

<h5 align="center" ><i>Meshery is a Cloud Native Computing Foundation Project</i></h5>
<hr>


<p align="center">
  <a href="https://hub.docker.com/r/meshery/meshery" alt="Docker pulls">
    <img src="https://img.shields.io/docker/pulls/layer5/meshery.svg" />
  </a>
  <a href="https://github.com/issues?q=is%3Aopen+is%3Aissue+archived%3Afalse+org%3Alayer5io+org%3Ameshery+org%3Aservice-mesh-performance+org%3Aservice-mesh-patterns+org%3A+label%3A%22help+wanted%22+" alt="GitHub issues by-label">
    <img src="https://img.shields.io/github/issues/layer5io/meshery/help%20wanted.svg?color=informational" />
  </a>
  <a href="https://github.com/meshery/meshery/blob/master/LICENSE" alt="LICENSE">
    <img src="https://img.shields.io/github/license/meshery/meshery?color=brightgreen" />
  </a>
  <a href="https://artifacthub.io/packages/helm/meshery/meshery" alt="Artifact Hub Meshery">
    <img src="https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/meshery" />
  </a>
  <a href="https://goreportcard.com/report/github.com/meshery/meshery" alt="Go Report Card">
    <img src="https://goreportcard.com/badge/github.com/meshery/meshery" />
  </a>
  <a href="https://github.com/meshery/meshery/actions" alt="Build Status">
    <img src="https://img.shields.io/github/actions/workflow/status/meshery/meshery/release-drafter.yml" >
  </a>
  <a href="https://bestpractices.coreinfrastructure.org/projects/3564" alt="CLI Best Practices">
    <img src="https://bestpractices.coreinfrastructure.org/projects/3564/badge" />
  </a>
  <a href="https://github.com/meshery/meshery/releases" alt="Meshery Downloads">
    <img src="https://img.shields.io/github/downloads/meshery/meshery/total" />
  </a>
</p>

<h4><p align="center"><i>Connect with Us</i></p></h4>
<h5 align="center">
    <a href="https://slack.meshery.io" alt="Join Slack">
     <img src="https://img.shields.io/badge/slack-meshery-darkred.svg?logo=slack"/>
  </a>
  <a href="https://twitter.com/mesheryio" alt="Follow meshery on X">
    <img src="https://img.shields.io/twitter/follow/meshery?style=social" />
  </a>
  <a href="https://www.linkedin.com/company/meshery" alt="Follow Meshery on LinkedIn">
    <img src="https://img.shields.io/badge/LinkedIn-meshery-blue.svg?logo=linkedin"/>
  </a>
  <a href="http://discuss.meshery.io" alt="Discuss Users">
    <img src="https://img.shields.io/discourse/users?label=discuss&logo=discourse&server=https%3A%2F%2Fdiscuss.layer5.io" />
  </a>
</h5>

<h4><p align="center"><i>Give us a star <a href="https://github.com/meshery/meshery/stargazers">⭐️</a> to show your support!</i></p></h4>
<hr>

## What is Meshery?

*<a href="https://meshery.io/" style="text-decoration: none;"><strong>Meshery</strong></a>, the open-source, cloud-native platform, simplifies Kubernetes and service mesh management. Developed by the <a href="https://github.com/layer5io" style="text-decoration: none;"><strong>Layer5</strong></a> community and now a part of **CNCF**, Meshery enables the design and management of Kubernetes-based infrastructure and applications across multi-cloud environments. With visual and collaborative GitOps, it frees users from YAML constraints, offering smooth, extensible management of multi-cluster deployments.*


<div align="center" width="100%">
 <p style="clear:both;">&nbsp;</p>
<h4>Try Meshery in the <a href="https://play.meshery.io">Cloud Native Playground</a></h4>

[videoplayback.webm](https://github.com/praduman8435/JAVA-DSA/assets/118579284/80f786fb-a3d2-41f5-8373-db58f7c0f48e)


</div>



<p style="clear:both;">&nbsp;</p>
<hr>

## Features of Meshery
### Configuration Management

*As an cloud native management platform, Meshery has built-in support infastructure and application configuration management. Use Meshery to configure your multi-cluster Kubernetes environments on-premises or across clouds. Meshery's configuration management is based on designs, which are documents that describe the desired state of an environment. Users can create designs using the design configurator in Meshery's UI or extensions like MeshMap, or they can hand code their designs using the mesh model spec. Meshery also offers a catalog of design templates that include configuration best practices.*

***Meshery's configuration management features include:***
<ul>

<li>
<details>
  <summary><a href="https://docs.meshery.io/concepts/logical/policies"><strong>Context-Aware Policies For Applications</strong></a></summary>
 <hr>

 *Leverage built-in relationships to enforce configuration best practices consistently from code to Kubernetes. Enhance development process by building custom rules in Open Policy Agent's Rego query language. Policies offer an evaluation algorithm to ensure desired behavior enforcement. Policies can be applied to components and relationships, defining rules and actions based on predefined conditions.*
     <p style="clear:both;">&nbsp;</p>
    <img alt="Meshery cloud native management" src="https://docs.meshery.io/assets/img/concepts/meshery-models-policy-evaluation.svg" style="margin-left:10px; margin-bottom:10px;" width="100%" align="center" />
     <p style="clear:both;">&nbsp;</p>
     To know more about it read [docs](https://docs.meshery.io/concepts/logical/policies)
     <hr>
</details>

</li>
    <li>
<details><summary><a href="https://docs.meshery.io/guides/configuration-management/"><strong>Multi-cluster management for Kubernetes and Clouds</strong></a></summary>
<hr>

*Meshery empowers seamless multi-cluster Kubernetes management, simplifying operations across diverse environments. With effortless task execution spanning clusters, intuitive contextual navigation, and precision configuration options, Meshery ensures optimized performance tailored to your needs. Its automated deployment mechanism, powered by Meshery's operator, guarantees synchronized operations and resource management efficiency across managed clusters. Ready for the future, Meshery's extensible architecture accommodates multi-cluster-specific components, promising streamlined operations and scalability. Experience enhanced productivity and control over your cloud-native infrastructure with Meshery's comprehensive multi-cluster management features.*
<br/><br/>
        <img alt="Wasm" src="https://meshery.io/assets/images/screens/multi-cluster-management.gif" style="margin-left:10px; margin-bottom:10px;" width="100%" align="center" />
       <hr>
</details>
  </li>
<li>
<details><summary><a href="https://docs.meshery.io/guides/configuration-management/creating-a-meshery-design"><strong>Cloud Native Design Patterns</strong></a></summary>
  <hr>

  *Patterns are essentially atomic designs with one or more components made in composed into an atomic, reusable design. Patterns are a way to apply the DRY principle when managing the configuration of cloud native infrastructure. Design and manage all of your cloud native infrastructure using the design configurator in Meshery or start from a template using the patterns from the [catalog](https://meshery.io/catalog).* <br/>

  [videoplayback.webm](https://github.com/praduman8435/JAVA-DSA/assets/118579284/80f786fb-a3d2-41f5-8373-db58f7c0f48e)
  <hr/>
  </details>
</li>

</ul>



### Performance Management

*Meshery helps users weigh the value of their cloud native deployments against the overhead incurred in running different deployment scenarios and different configruations. Meshery provides statistical analysis of the request latency and throughput seen across various permutations of your workload, infrastructure and infrastructure configuration. In addition to request latency and throughput, Meshery also tracks memory and CPU overhead in of the nodes in your cluster. Establish a performance benchmark and track performance against this baseline as your environment changes over time.*

***Meshery's performance management features include:***
<ul>
<li>
<details><summary><a href="https://smp-spec.io/"><strong>Service Mesh Performance</strong></a></summary><hr>

<picture align="left">
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/layer5io/layer5/master/src/assets/images/service-mesh-performance/stacked/smp-light-text.svg"  width="18%" align="left" style="margin-left:10px;" />
  <img alt="Shows an illustrated light mode meshery logo in light color mode and a dark mode meshery logo dark color mode." src="https://raw.githubusercontent.com/layer5io/layer5/master/src/assets/images/service-mesh-performance/stacked/smp-light-text.svg" width="18%" align="left" style="margin-left:10px;" />
</picture>

*In an effort to produce service mesh agnostic tooling, Meshery uses the ***<a href="https://smp-spec.io">Service Mesh Performance</a>*** specification as a common format to capture and measure your infrastructure's performance against a universal cloud native performance index. As an implementation of Service Mesh Interface (SMI), Meshery participates in advancing cloud native infrastructure adoption through the standardization of API*
<br /><br />
<p style="clear:both;">

### Standardized Cloud Native Performance Management

*Meshery is the cloud native utility for uniformly managing the performance of microservices and the infrastructure that run them. As an implementation of the Service Mesh Performance ***[(SMP)](https://github.com/service-mesh-performance/service-mesh-performance/issues)***, Meshery enables you to measure the value provided by Docker, Kubernetes, or a service mesh in the context of the overhead incurred*

<ul>
<li>See your performance regressions before you merge
</li>
  <li>Connect Performance Management to your GitHub repo and see changes pull request-to-pull request
</li>
    <li>Red light performance regressions</li>
      <li>Baseline and analyze the performance of your services is key to efficient operation of any application</li>
    <li>Meshery is the canonical implementation of the Cloud Native Performance specification</li>
    <li>Define your performance profiles upfront. See statistcal analysis with microservice latency and throughput quartiles</li>
    <li>Meshery packages all these features into an easy-to-use GitHub Action</li>
    <li>Meshery includes your choice of load generator, so that you can meausure your way
</li>

</ul>
<hr>
  </details>
</li>

  <li>
<details>
  <summary><a href="https://docs.meshery.io/guides/performance-management/managing-performance"><strong>Manage infra performance with Performance Profiles</strong></a></summary>
 <hr>

 ***[Performance Profile](https://docs.meshery.io/guides/performance-management/managing-performance)** is a new and improved way to save performance test configurations, schedule performance tests, etc. Meshery Performance Profiles offers the following features:*
 <ul>
  <li>Saving Test Configurations</li>
  <li>Create schedules for running tests automatically</li>
  <li>Meshery offers UI for describing the schedules which are as powerful as CRON but without the hassle of remembering its syntax!</li>
  <li>All the tests performed are under Profiles which has multiple advantages like:<ul>
   <li>Access test results ran under a certain profiles at one place</li>
   <li>If running ad-hoc tests, user can anytime come back and rename the performance profiles for re-use</li>
  </ul></li>

 </ul>
    <p style="clear:both;">&nbsp;</p>
    <img alt="Meshery cloud native management" src="https://meshery.io/assets/images/features/performance-metrics.gif" style="margin-left:10px; margin-bottom:10px;" width="100%" align="center" />
    <p style="clear:both;">&nbsp;</p>
    <hr>
</details>
  </li>


  </ul>


### Lifecycle Management

*Meshery simplifies the lifecycle management of your Kubernetes clusters and cloud-native infrastructure, offering robust provisioning, configuration, and operational capabilities. With support for a wide array of cloud-native integrations, Meshery ensures seamless compatibility across diverse environments.*

***Meshery's Lifecycle management features include:***
<ul>

  <li>
<details><summary><a href="https://docs.meshery.io/guides/configuration-management/filter-management"><strong>WASM Envoy Filter Management</strong></a></summary>
<hr>

*Meshery offers WebAssembly-based network traffic management filters in Envoy, bringing a new level of flexibility and control to your Istio deployments. Allows users to manage WASM filters for data plane proxies to improve agility, maintainability, and more.*
<br/><br/>
        <img alt="Wasm" src="https://meshery.io/assets/images/meshery-wasm.png" style="margin-left:10px; margin-bottom:10px;" width="100%" align="center" />


 <h5>Retrieves and presents the contents of a designated filter through identification by name or ID</h5><br/>
<ul><li>Displays the contents of a specific filter based on name or id</li>

 ```bash
  mesheryctl filter view [filter-name | ID]
 ```
 <li>View all filter files
 </li>

 ```bash
  mesheryctl filter view --all
 ```
</ul>

*You can read ***[docs](https://docs.meshery.io/guides/configuration-management/filter-management)*** to know more about it*
 <hr>
</details>


  </li>

 <li>
<details>
  <summary><a href="https://docs.meshery.io/concepts/logical/environments"><strong>Operation and Monitoring</strong></a></summary>
 <hr>


  *Elevate your infrastructure's performance with Meshery. Ensure smooth operations by checking your configurations with the validator, keep an eye on performance metrics, and optimize resource usage effortlessly. Meshery makes it easy to manage and enhance your system's performance.*
    <p style="clear:both;">&nbsp;</p>
    <img alt="Meshery cloud native management" src="https://meshery.io/assets/images/screens/configuration-best-practices.gif" style="margin-left:10px; margin-bottom:10px;" width="100%" align="center" />
      <p style="clear:both;">&nbsp;</p>
      <hr>
</details>

</li>

<li>
<details><summary><a href="https://docs.meshery.io/guides/configuration-management/creating-a-meshery-design"><strong>Infrastructure Managment with GitOps Snapshots</strong></a></summary>
<hr>

  *Streamlining Platform Engineering with GitOps and Meshery. Preview your deployment, view changes pull request-to-pull request and get infrastructure snapshots within your PRs by connecting MeshMap to your GitHub repositories.*
  <ul>
    <li>See your deployment before you merge</li>
    <li>Connect MeshMap to your GitHub repo and see changes pull request-to-pull request</li>
    <li>Get snapshots of your infrastructure directly in your PRs</li>

  </ul>

<p style="clear:both;">&nbsp;</p>
    <img alt="Meshery cloud native management" src="https://docs.meshery.io/assets/img/meshmap/meshmap-snapshot.png" style="margin-left:10px; margin-bottom:10px;" width="100%" align="center" />
     <p style="clear:both;">&nbsp;</p>

  See ***[Extension: MeshMap Snapshot](https://docs.meshery.io/extensions/snapshot)*** for more details.
  <hr>
  </details>
</li>


<li>
<details>
  <summary><a href="https://docs.meshery.io/concepts/logical/environments"><strong>Manage your connections with Environments</strong></a></summary>
 <hr>

 ***[Meshery Environments](https://docs.meshery.io/concepts/logical/environments)** allow you to logically group related Connections and their associated Credentials. Environments make it easier for you to manage, share, and work with a collection of resources as a group, instead of dealing with all your Connections and Credentials on an individual basis.*
     <p style="clear:both;">&nbsp;</p>
    <img alt="Meshery cloud native management" src="https://meshery.io/assets/images/features/environments.gif" style="margin-left:10px; margin-bottom:10px;" width="100%" align="center" />
    <p style="clear:both;">&nbsp;</p>
</details>

</li>
</ul>
<hr>

## Collaboration

 ***[Meshery Workspaces](https://docs.meshery.io/concepts/logical/workspaces)** serve as a virtual space for your team-based work. Create a Workspace to organize your work and to serve as the central point of collaboration for you and your teams and a central point of access control to Environments and their resources*

*You may create Workspaces to organize project-based work or to create domains of responsibility for your teams or segregate Designs and Environments and track team activity*

<img alt="Meshery cloud native management" src="https://raw.githubusercontent.com/layer5io/meshery.io/master/assets/images/features/workspace.gif" style="margin-left:10px; margin-bottom:10px;" width="100%" align="center" />
     <p style="clear:both;">&nbsp;</p>

<hr>

## Extension Points

 *Extend Meshery's functionality to suit your unique requirements with its vast set of extensibility features. Leverage gRPC adapters, hot-loadable ReactJS packages, Golang plugins, and more to customize and enhance Meshery's capabilities. With consumable and extendable API interfaces via REST and GraphQL, Meshery serves as your ultimate self-service engineering platform, adaptable to your evolving needs.The great number of **[extension points](https://docs.meshery.io/extensibility)** in Meshery make it ideal as the foundation of your internal developer platform*
    <img alt="Meshery cloud native management" src="https://docs.meshery.io/assets/img/architecture/meshery_extension_points.svg" style="margin-left:10px; margin-bottom:10px;" width="100%" align="center" />
        <p style="clear:both;">&nbsp;</p>
   <hr>


## Meshery Architecture

*Meshery is deployed as a set of containers that can run on Docker or Kubernetes. Meshery architecture is categorized into three components, the Meshery server, Meshery UI and Meshery Providers. These components can communicate with one another through the help of the gRPC request*

*Meshery Server store the location of the other components and connects to them as needed. Using Meshery UI, CTL or scripts, an administrator deploys the service mesh components to a target Kubernetes cluster. When an administrator requests to retrieve information from the Adapter or invoke an Adapter's operation, this initiates a connection between the Meshery Server to Meshery Adapters. Each Meshery Adapter manages its service mesh*

<p align="center"><a href="https://docs.meshery.io/concepts/architecture"><img src="https://docs.meshery.io/assets/img/architecture/Meshery-client-architecture.svg" width="90%" align="center" /></a></p>

***Learn more about*** **<a href="https://docs.meshery.io/architecture">Meshery's architecture</a>**

<p style="clear:both;">&nbsp;</p>
<hr>

## Get Started with Meshery
Getting Meshery up and running locally on a Docker-enabled system or in Kubernetes is easy. Meshery deploys as a set of Docker containers, which can be deployed to either a Docker host or Kubernetes cluster
<br/>
### All Supported Platforms
*Download, install, and run Meshery in a single command. See all supported* **<a href="https://docs.meshery.io/installation">platforms</a>**

### Installation Guide
*To get Meshery running up in a docker or Kubernetes system. you need to Install Meshery command line client ***mesheryctl*** with the command:*
- **For Kubernetes**
```bash
curl -L https://meshery.io/install | PLATFORM=kubernetes bash -
```

- **For Docker**
```bash
curl -L https://meshery.io/install | PLATFORM=docker bash -
```

### Access Meshery
*Your default browser will be opened and directed to Meshery’s web-based user interface typically found at http://localhost:9081*




### Select a Provider

  *Select from the list of Providers in order to login to Meshery. Authenticate with your chosen Provider*

### Configure Connections to your Kubernetes Clusters

  *If you have deployed Meshery out-of-cluster, Meshery will automatically connect to any available Kubernetes clusters found in your kubeconfig (under $HOME/.kube/config). If you have deployed Meshery out-of-cluster, Meshery will automatically connect to the Kubernetes API Server availabe in the control plane. Ensure that Meshery is connected to one or more of your Kubernetes clusters*
<img alt="Meshery cloud native management" src="https://docs.meshery.io/assets/img/platforms/meshery-settings.png" style="margin-left:10px; margin-bottom:10px;" width="100%" align="center" />

*If your config has not been autodetected, you can manually upload your kubeconfig file (or any number of kubeconfig files). By default, Meshery will attempt to connect to and deploy Meshery Operator to each reachable context contained in the imported kubeconfig files. See Managing Kubernetes Clusters for more information*

### Verify Deployment

  *Run connectivity tests and verify the health of your Meshery system. Verify Meshery’s connection to your Kubernetes clusters by clicking on the connection chip. A quick connectivity test will run and inform you of Meshery’s ability to reach and authenticate to your Kubernetes control plane(s). You will be notified of your connection status. You can also verify any other connection between Meshery and either its components (like **<a href="https://docs.meshery.io/concepts/architecture/adapters">Meshery Adapters</a>**) or other managed infrastructure by clicking on any of the connection chips. When clicked, a chip will perform an ad hoc connectivity test*


### Design and operate Kubernetes clusters and their workloads

  *You may now proceed to managed any cloud native infrastructure supported by Meshery. See all integrations for a complete list of supported infrastructure*
<img alt="Meshery cloud native management" src="https://docs.meshery.io/assets/img/platforms/meshery-designs.png" style="margin-left:10px; margin-bottom:10px;" width="100%" align="center" />

### Additional Guides
- [Troubleshooting Meshery Installations](https://docs.meshery.io/guides/troubleshooting/installation)
- [Meshery Error Code Reference](https://docs.meshery.io/reference/error-codes)
- [Mesheryctl system check](https://docs.meshery.io/reference/mesheryctl/system/check)


<p style="clear:both;">&nbsp;</p>
<hr>

## Join the Meshery community!

Our projects are all about teamwork and we're eager to collaborate! 🚀 Explore our ***<a href="https://layer5.io/community/newcomers">Contributor's Journey Map</a>*** and ***<a href="https://layer5.io/community/handbook">Community Handbook</a>*** for a guided tour of resources available to you. Want a sneak peek at each repository? Check out the ***<a href="https://layer5.io/community/handbook/repository-overview">Repository Overview</a>***, categorized by technology and programming language. Ready to dive in? Join our ***<a href="https://slack.meshery.io">community Slack</a>*** Or ***<a href="http://discuss.meshery.io"> discussion forum</a>*** and start participating! 💬

<p style="clear:both;">&nbsp;</p>



<a href ="https://layer5.io/community/meshmates"><img alt="MeshMates" src=".github/assets/images/readme/layer5-community-sign.png" style="margin-right:36px; margin-bottom:7px;" width="140px" align="left" /></a>
***<h3>Find your MeshMate</h3>***

<p><a href="https://layer5.io/community/meshmates"><strong>MeshMates</strong></a> are experienced Layer5 community members, who will help you learn your way around, discover live projects, and expand your community network. Connect with a Meshmate today!</p>
Layer5 MeshMates are committed to helping community members be successful contributors. MeshMates aid in identifying areas of projects to engage within, working groups to join, and in helping community members grow in their open source and cloud native knowledge. By connecting one-on-one, MeshMates will share tips on how to have the best community experience possible
</p>

<p style="clear:both;">&nbsp;</p>

<a href="https://meshery.io/community"><img alt="Layer5 Cloud Native Community" src="https://docs.meshery.io/assets/img/readme/community.png" width="140px" style="margin-right:36px; margin-bottom:7px;" width="140px" align="left"/></a>
***<h3>Participate</h3>***


<p>
✔️ <strong>Join</strong> any or all of the weekly meetings on <a href="https://meshery.io/calendar">Community calendar</a>. <br />
✔️ <strong>Watch</strong> Community <a href="https://www.youtube.com/playlist?list=PL3A-A6hPO2IMPPqVjuzgqNU5xwnFFn3n0">meeting recordings</a>. <br />
✔️ <strong>Fill-in</strong> a <a href="https://layer5.io/newcomers">Community member form</a> to gain access to community resources. <br />
✔️ <strong>Discuss</strong> in the <a href="http://discuss.meshery.io">Community Forum</a>.<br />
✔️ <strong>Explore more</strong> in the <a href="https://layer5.io/community/handbook">Community Handbook</a>. <br />
</p>

<p style="clear:both;">&nbsp;</p>

<a href="https://slack.meshery.io">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/meshery/meshery/master/.github/assets/images/readme/slack.svg"  style="margin-right:36px; margin-bottom:7px;" width="140px" align="left"/>
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/meshery/meshery/master/.github/assets/images/readme/slack.svg" style="margin-right:36px; margin-bottom:7px;" width="140px" align="left" />
  <img alt="Shows an illustrated light mode meshery logo in light color mode and a dark mode meshery logo dark color mode." src="https://raw.githubusercontent.com/meshery/meshery/master/.github/assets/images/readme/slack.svg" width="140px" style="margin-right:36px; margin-bottom:7px;" width="140px" align="left" />
</picture>
</a>
<p style="clear:both;">&nbsp;</p>

<p align="left"><i>Not sure where to start?</i> Grab an open issue with the <a href="https://github.com/meshery/meshery/labels/good%20first%20issue"><strong>good-first-issue</strong></a>.
</p>
<p>Join <a href="https://slack.meshery.io"><strong>Slack workspace</strong></a> to get engaged with the Community</p>

<p style="clear:both;">&nbsp;</p>
<p style="clear:both;">&nbsp;</p>
<hr>

## :handshake: Contributing

Please do! We're a warm and welcoming community of open source contributors. Whether you're a seasoned developer or just starting out, your contributions are invaluable to us. Be sure to read the **[Contribution Guides](https://docs.meshery.io/project/contributing#general-contribution-flow)** for a tour of resources available to you and how to get started.

<h4>
<details>
  <summary>Follow these steps to start your contributing journey</summary>
 <hr>
 <ol>
  <li>See the <a href="https://layer5.io/community/newcomers"><strong>Newcomers Guide</strong></a> for how, where, and why to contribute</li><br/>
  <li>Sign up for a <a href="https://layer5.io/community/meshmates"><strong> MeshMate</strong></a> to find the perfect Mentor to help you explore the Layer5 projects and find your place in the community :
  <ul><br/>
   <li><strong>Familiarize </strong>yourself with the broader set of community projects (take a look at the<a href="https://layer5.io/community/handbook/repository-overview"><strong>Repository Overview:</strong></a> Spend time understanding each of the initiatives through high-level overviews available in the community drive and through discussions with your MeshMate</li><br/>
   <li><strong>Identify your area of interest:</strong> Use the time with your MeshMate to familiarize yourself with the architecture and technologies used in the projects. Inform your MeshMate of your current skills and what skills you aim to develop</li><br/>
   <li><strong>Run Meshery: </strong>Put on your user hat and walk-through all of Meshery’s features and functions as a user</li><br/>
   <li><strong>Build Meshery: </strong>Confirm that you have a usable development environment</li><br/>
   <li><strong>Discuss </strong>with the community by engaging in the <a href="http://discuss.meshery.io/"><strong>discussion forum</strong></a></li><br/>
   <li><strong>Contribute </strong>by grabbing any open issue with the <a href="https://github.com/meshery/meshery/issues/"><strong>help-wanted label </strong></a>and jump in. If needed, create a <a href="https://github.com/meshery/meshery/issues/new/choose"><strong>new issue </strong></a>All <a href="https://github.com/meshery/meshery/pulls"><strong>pull requests </strong></a>should reference an open issue. Include keywords in your pull request descriptions, as well as commit messages, to <a href="https://help.github.com/en/github/managing-your-work-on-github/closing-issues-using-keywords"><strong></strong>automatically close issues in GitHub</strong></a></li><br/>
   <li><strong>Fill-in </strong>a <a href="https://layer5.io/newcomers"><strong>community member form </strong></a>to gain access to community resources</li><br/>

  </ul>
  </li>
 </ol>
<hr>
 <p style="clear:both;">&nbsp;</p>

</details>
</h4>

*Didn't find an answer to your question? **[Just Ask](http://discuss.meshery.io/)** 💬*
### Contributors:

<p Style="align:center"><i>🙌 We are deeply grateful to all our amazing contributors! 🌟</i></p>

<br/>

<a href="https://github.com/meshery/meshery/graphs/contributors">
 <img src="https://contrib.rocks/image?repo=meshery/meshery" />
</a>
<p style="clear:both;">&nbsp;</p>
<hr>


## :bookmark: License

<i>This repository and site are available as open-source under the terms of the</i> **[Apache 2.0 License](https://opensource.org/licenses/Apache-2.0)**
