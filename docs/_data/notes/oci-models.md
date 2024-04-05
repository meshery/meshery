# OCI & Meshery Entities

packaging and distribution

Design Prologue 3
Design Goals and Objectives 3
FAQ 3
Architectural Considerations 4
Available Libraries and Reference Implementations 4
Sourcing Image Data from Registries 6
OCI Registry APIs: Standard Capabilities of the OCI Distribution API 6
Docker Hub 6
Functional Architecture 6
Architecture Diagram 6
Sequence Diagram 7
System Flows 7
Flow: Meshery Server Connects to an OCI-compliant Registry unauthenticated 7
Flow: Meshery Server Connects to an OCI-compliant Registry authenticated 8
Action Items: 8
Flow: <short title> 10
Flow: <short title> 10
Flow: <short title> 11

Design Prologue
Meshery Models represent a schema-based description of cloud native infrastructure and operations of which Meshery is capable of managing. Models and Designs need to be portable between Meshery deployments as well as easily versionable in external repositories.

Related Issues
<https://github.com/meshery/meshery/issues/6447>
<https://github.com/meshery/meshery/issues/8855>
[Feat]: Support models export as OCI format #10165
Support OCI export of Designs, Patterns, Filters #10195

Design Goals and Objectives
The designs in this specification should result in upholding the following goals and should achieve these specific functions:

Goal 1: Ensure Meshery Models and Designs enjoy a high degree of portability
OCI is the broadly accepted standard image definition and distribution format.
Objective 1: Import and export of Meshery Entities
Import and export of Meshery Designs in OCI format.
Import and export of Meshery Models in OCI format.
Objective 2: Support for OCI-compatible registries, in priority order:
Docker Hub, Artifact Hub, AWS ECR, GCP GCR, Azure ACR, JFrog Artifactory, Harbor,
Objective 3: Native support for Meshery Designs in Artifact Hub

Goal 2: Overarching idea
Supporting statement
Objective 1: Measurable and specific item
Objective 2: Measurable and specific item
FAQ
How does Meshery version Models and Designs?
And its answer
Common question
And its answer
Architectural Considerations
The following section outlines factors considered when identifying and determining a best-fit Golang library to use as the primary library for interaction with OCI Registries.
Available Libraries and Reference Implementations
ORAS (library)
GitHub, Website
Library used: (name, description, and link to library needed)
Mainly uses libraries
github.com/opencontainers/image-spec/specs-go/v1
github.com/opencontainers/go-digest
There is go-client for ORAS
github.com/oras-project/oras-go
FluxCD (utility)
GitHub, Website
Library used: (name, description, and link to library needed) (golang container registry client)
Mainly uses libraries
github.com/distribution/distribution
github.com/google/go-containerregistry
github.com/fluxcd/pkg/
KubeApp (utility)
GitHub, Website (link needed)Library used: (name, description, and link to library needed)
Mainly uses libraries
github.com/fluxcd/pkg
github.com/fluxcd/source-controller
oras.land/oras-go/
github.com/opencontainers/image-spec
Docker Hub CLI tool(utility)
GitHub, Website
Library used: (name, description, and link to library needed)
Mainly uses library
github.com/opencontainers/image-spec
github.com/docker/distribution

Does docker CLI support pushing and pulling OCI artifacts? Yes, Docker CLI does support pulling and pushing OCI artifacts from/to OCI registries.
If so, only from Docker Hub or other compliant registries? Docker CLI’s support for OCI artifacts is not limited to Docker Hub. You can use Docker CLI to manage OCI artifacts with any OCI-compliant registry. This includes Docker Hub, Azure Container Registry, and others. Docker Hub supports OCI artifacts by leveraging the config property on the image manifest.

Does AWS CLI support pushing and pulling OCI artifacts? Yes, the AWS CLI does support pushing and pulling OCI artifacts in Amazon ECR. ECR stores Docker images, OCI images, and OCI-compatible artifacts
If so, only from ECR or other compliant registries? No, only ECR.

Does KubeApp Have a cli? No, Kubeapp doesn't have its own dedicated CLI tool. However, It does integrate with the kubectl command-line tool and uses Kubernetes APIs for its operations.

Reference
<https://docs.docker.com/docker-hub/oci-artifacts/>
<https://www.docker.com/blog/announcing-docker-hub-oci-artifacts-support/>
<https://aws.amazon.com/blogs/containers/oci-artifact-support-in-amazon-ecr/>
<https://docs.aws.amazon.com/AmazonECR/latest/userguide/what-is-ecr.html>

What library does ORAS use under the hood?
ORAS mainly uses two libraries:
github.com/opencontainers/image-spec/specs-go/v1
github.com/opencontainers/go-digest

What library does KubeApps use under the hood? Why is it used?
KubeApps uses:
github.com/fluxcd/pkg/
The Flux helm controller allows you to manage Helm chart releases declaratively with Kubernetes manifests. Similarly, the Flux HelmRepository can be used to manage Helm repositories declaratively.
github.com/oras-project/oras-go
KubeApps uses ORAS CLI oras is used to consume Tanzu Application Catalog metadata from an OCI registry.

Why using multiple library?
Two libraries were uses for different reasons. Oras is used for retrieve all Tanzu Application Catalog metadata from a given registry for containers and charts.but flux cd is primarily focused on managing helm charts and repositories.

What library does Dockerhub cli use under the hood?
DockerHub mainly uses library:
github.com/opencontainers/image-spec
github.com/docker/distribution
Sourcing Image Data from Registries
Container Registry Metadata: Container registries often allow storing additional metadata alongside images, which can be used to document build information, dependencies, or usage instructions. This information serves a similar purpose to Dockerfile comments and instructions.
OCI Registry APIs: Standard Capabilities of the OCI Distribution API
<https://github.com/opencontainers/distribution-spec/blob/main/spec.md>
Docker Hub
Docker Hub API
Docker Hub Tool
Docker Hub pull rate limiting
Functional Architecture
Architecture Diagram
<here>

Sequence Diagram
<here>
System Flows
Flow: Meshery Server Connects to an OCI-compliant Registry unauthenticated
Goal: As a system,
Meshery Server needs to connect to OCI-compliant registries,
so that it can offer an index of container images and their details to Meshery users.
Actors:
Meshery Server, OCI Registry
Preconditions:
User or System has created a connection for any given OCI Registry.
Assumptions:
User does not provide credentials.
Registry is indexed without authentication.
Rate limiting will ensue.
Behavior:
Retrieve the following data:
Repo
Image
Tags (how many? - - latest, then paginate)
Dockerfile
All exposed ports
All environment variables
???
Username?

Implementation:
Server stores container image details and Registry Connection index inside of Meshery Database.

Error Handling:

Postconditions:
Upon disconnection of the OCI Registry Connection, Meshery Server does not retain the registry’s index of container images.

Acceptance Criteria:

Flow: Meshery Server Connects to an OCI-compliant Registry authenticated
Goal: As a [developer/integrator/mesheryctl/operator] user,
I would like to …….. ,
so that ………..
Actors:

Preconditions:

Assumptions:

Behavior:
Validate user credentials.
Implementation:

Error Handling:

Acceptance Criteria:

Action Items:
Support OCI export of Contents (Meshery Designs, Patterns, Filters)
Meshery UI (Completed)
mesheryctl: Create a new issue for supporting this functionality.
Support pulling and pushing OCI artifacts from/to OCI registries.
ORAS VS FLUX
References: <https://fluxcd.io/flux/cmd/flux_push/>
<https://github.com/google/go-containerregistry>
Do `docker` or `aws` CLIs support push/pull? If so, only from Hub and ECR or all compliant registries? <Althaf>
[Lee] the `docker` image exporter:
registry: exports the build result into a container image, and pushes it to the specified registry. See Docker Docs.

What library does ORAS use under the hood? <Althaf>
What library does KubeApps use under the hood? <Althaf>
What library does Docker Hub CLI tool use under the hood? <Althaf>
Number of registries supported.
How frequently does ORAS/Flux discover the latest tags?
Support mesheryctl design import as a new command.
Consolidate mesheryctl app command into mesheryctl pattern
Improve descriptions of parent and subcommands.
Pulling container images from OCI registries for
Meshery Models
Meshery Designs (tbd - base)
Do all OCI registries have Dockerfiles available for their images?
What manifest or index does OCI offer to find information similar as to what is in the Dockerfile?
Not in the case of non-runnable images.
Not in the case of non-Docker images.
What % of images have Dockerfiles?
What attributes or qualities of an image determine (identify) its runnability?
Is there a common mechanism for images on Docker Hub? What about other registries (a standard OCI Distribution spec?)
Other than Dockerfiles, what other source of image metadata is available?
The manifest inside of the OCI image.
Investigate API structure for the OCI registries.
See <https://github.com/opencontainers/distribution-spec>
Real-time?
How frequently to index the registry?
When searching from the registry refer to the actual registry, should the local server index be refreshed?
Familiarize with the metadata provided with each of these registries.
How does Kubernetes support sourcing multi-registry deployments?

Flow: <short title>
Goal: As a [developer/integrator/mesheryctl/operator] user,
I would like to …….. ,
so that ………..
Actors:

Preconditions:

Assumptions:

Behavior:

Implementation:

Error Handling:

Acceptance Criteria:

Flow: <short title>
Goal: As a [developer/integrator/mesheryctl/operator] user,
I would like to …….. ,
so that ………..
Actors:

Preconditions:

Assumptions:

Behavior:

Implementation:

Error Handling:

Acceptance Criteria:

Flow: <short title>
Goal: As a [developer/integrator/mesheryctl/operator] user,
I would like to …….. ,
so that ………..
Actors:

Preconditions:

Assumptions:

Behavior:

Implementation:

Error Handling:

Acceptance Criteria:
