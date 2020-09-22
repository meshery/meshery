---
layout: page
title: Releases
permalink: project/releases
---

<!-- A list of the [releases of Meshery](https://github.com/layer5io/meshery/releases). See the [Build and Release Strategy](https://docs.google.com/document/d/11nAxYtz2SUusCYZ0JeNRrOLIxkgmmbUVWz63MBZV2oE/edit?usp=sharing) document for details. -->

<!-- <div class="card">
  <div class="card-header">
    Featured
  </div>
  <div class="card-body">
    <h5 class="card-title">Special title treatment</h5>
    <p class="card-text">With supporting text below as a natural lead-in to additional content.</p>
    <a href="#" class="btn btn-primary">Go somewhere</a>
  </div>
</div>
<br>
<div class="card">
  <div class="card-header">
    Featured
  </div>
  <div class="card-body">
    <h5 class="card-title">Special title treatment</h5>
    <p class="card-text">With supporting text below as a natural lead-in to additional content.</p>
    <a href="#" class="btn btn-primary">Go somewhere</a>
  </div>
</div>
<br>
<div class="card">
  <div class="card-header">
    Featured
  </div>
  <div class="card-body">
    <h5 class="card-title">Special title treatment</h5>
    <p class="card-text">With supporting text below as a natural lead-in to additional content.</p>
    <a href="#" class="btn btn-primary">Go somewhere</a>
  </div>
</div>
<div class="card">
  <div class="card-header">
    Featured
  </div>
  <div class="card-body">
    <h5 class="card-title">Special title treatment</h5>
    <p class="card-text">With supporting text below as a natural lead-in to additional content.</p>
    <a href="#" class="btn btn-primary">Go somewhere</a>
  </div>
</div>
<div class="card">
  <div class="card-header">
    Featured
  </div>
  <div class="card-body">
    <h5 class="card-title">Special title treatment</h5>
    <p class="card-text">With supporting text below as a natural lead-in to additional content.</p>
    <a href="#" class="btn btn-primary">Go somewhere</a>
  </div>
</div> -->

# Meshery Releases


### [v0.4.1](https://github.com/layer5io/meshery/releases/tag/v0.4.1)
Released on 8 August 

**What's new ?**
- Documentation
   + Typo and error fixes



### [v0.4.0-beta.4](https://github.com/layer5io/meshery/releases/tag/v0.4.0-beta.4)
Released on 9 July 2020

**What's new ?**
- Bug Fixes
    + hotfix: Meshery settings page not building

### [v0.4.0-beta.3](https://github.com/layer5io/meshery/releases/tag/v0.4.0-beta.3)
Released on 8 July 2020

**What's new ?**
- Bug Fixes
    + Hotfix of bug introduced in Release v0.4.0-beta.2

### [v0.4.0-beta.2](https://github.com/layer5io/meshery/releases/tag/v0.4.0-beta.2)
Released on 8 July 2020

**What's new ?**
- Meshery
    + Added initial kubeconfig.
- Mesheryctl
    + Removed fatal error due to missing `meshery.yaml`.
- Documentation
    + Updated Kuma adapter logo in Meshery Docs.
    + Added network ports to the adapters table.
    + Added providers screenshot.
    + Added all supported platforms to list.
    + Added nighthawk as a load generator.
    + Fixed broken link in EKS page.

### [v0.3.19](https://github.com/layer5io/meshery/releases/tag/v0.3.19)
Released on 7 July 2020

**What's new ?**
- Mesheryctl
    + Important bug fix for failure to initialize app folder (`~/meshery`) and app config file (`meshery.yaml`).

### [v0.4.0-beta.1](https://github.com/layer5io/meshery/releases/tag/v0.4.0-beta.1)
Released on 5 July 2020

**What's new ?**
- Meshery
    + "stable-latest" tag for new release channel.

### [v0.3.18](https://github.com/layer5io/meshery/releases/tag/v0.3.18)
Released on 4 July 2020

**What's new ?**
- Mesheryctl
    + `mesheryctl version` now retrieves server version information, including release tag and git sha.
    + `meshery perf` now supports SPMS as a package through the `--file` flag.
- Meshery
    + Extensibility enhancement: an abstract interface, Load Generator Interface, is now provided for ease of extending Meshery to support other load generators, like Nighthawk.
- Docs
    + uplifted EKS and GKE instructions

### [v0.3.17](https://github.com/layer5io/meshery/releases/tag/v0.3.17)
Released on 1 July 2020

**What's new ?**
- Meshery
    + Feature: New "stable" release channel established.
    + Feature: Enhanced UX for advanced performance test HTTP headers.

### [v0.3.16](https://github.com/layer5io/meshery/releases/tag/v0.3.16)
Released on 22 June 2020

**What's new ?**
- Mesheryctl
    + Feature: Addition of `mesheryctl system` as the new command for lifecycle management of Meshery.
        + All Meshery lifecycle management commands like `start`, `stop`, `reset`, `logs` and so on are now found under `system`.
    + Feature: Addition of `mesheryctl system config` as a new subcommand to assist with configuring Meshery to communicate with Kubernetes deployments in GKE, AKS, and EKS.

### [v0.3.15](https://github.com/layer5io/meshery/releases/tag/v0.3.15)
Released on 1 May 2020

**What's new ?**
- Mesheryctl
    + Feature: Addition of `mesheryctl perf --file <smp.yaml>' support for [SMP formatted](https://github.com/layer5io/service-mesh-performance) performance test profiles.
    + Feature: Behavior change `mesheryctl perf` for using tokens and supporting JWT.

### [v0.3.14](https://github.com/layer5io/meshery/releases/tag/v0.3.14) 

**What's new ?**
- Mesheryctl
    + Behavior change in `mesheryctl start` to check for new Meshery server images by default.

### [v0.3.13](https://github.com/layer5io/meshery/releases/tag/v0.3.13)

**What's new ?**
- Mesheryctl
    + Point release for mesheryctl's initial ARM support.

### [v0.3.12](https://github.com/layer5io/meshery/releases/tag/v0.3.12)
**What's new ?**
- Mesheryctl
    + Initial Scoop Bucket and Scoop package for `mesheryctl`.

### [v0.3.11](https://github.com/layer5io/meshery/releases/tag/v0.3.11)

**What's new ?**
- Mesheryctl
    + Fixes minor user experience issues on `mesheryctl perf` command. See [Meshery CLI Commands & Documentation](https://docs.google.com/document/d/1xRlFpElRmybJ3WacgPKXgCSiQ2poJl3iCCV1dAalf0k/edit#) for reference.

### [v0.3.10](https://github.com/layer5io/meshery/releases/tag/v0.3.10)

**What's new ?**
- Mesheryctl
    + Introduces `mesheryctl perf` command. See [Meshery CLI Commands & Documentation](https://docs.google.com/document/d/1xRlFpElRmybJ3WacgPKXgCSiQ2poJl3iCCV1dAalf0k/edit#) for reference.


### [v0.3.9](https://github.com/layer5io/meshery/releases/tag/v0.3.12)

**What's new ?**
- Mesheryctl
    + Addition of `mesheryctl version` to provide server-side version number.
    + Improvement of `mesheryctl logs | stop | start` to provide appropriate grammar in situations when Meshery is stopped or Docker is not present.
- Meshery Server
    + Strip off extraneous information beyond IP address and port in Grafana and Prometheus endpoints ([#612](https://github.com/layer5io/meshery/issues/612))

### [v0.3.8](https://github.com/layer5io/meshery/releases/tag/v0.3.8)

**What's new ?**
- Mesheryctl
    + Extraneous command line output removed. Clarity of CLI interaction with `mesheryctl` is improved in this release.

### [v0.3.7](https://github.com/layer5io/meshery/releases/tag/v0.3.7)

**What's new ?**
- Meshery Server
    + Ad-hoc connectivity tests for Prometheus is now supported. Users can click the Prometheus chip and have Meshery verify its ability to connect to the configured Prometheus instance.

### [v0.3.6](https://github.com/layer5io/meshery/releases/tag/v0.3.6)

**What's new ?**
- Mesheryctl
    + Ad-hoc connectivity tests for Grafana is now supported. Users can click the Grafana chip and have Meshery verify its ability to connect to the configured Grafana instance.

### [v0.3.5](https://github.com/layer5io/meshery/releases/tag/v0.3.5)

**What's new ?**
- Mesheryctl
    + Removal of `init` as a command exposed to users. This command's functionality is used internal to `mesheryctl start`. A new `start --check` command will provide preflight check functionality in `init`'s place.

### [v0.3.4](https://github.com/layer5io/meshery/releases/tag/v0.3.4)

**What's new ?**
- Mesheryctl
    + `mesheryctl version` is now enhanced with the addition of displaying the git commit (sha) of the mesheryctl release.

### [v0.3.3](https://github.com/layer5io/meshery/releases/tag/v0.3.3)

**What's new ?**
- Meshery Server
    + Providers - a new project construct that allows users to select authentication, long-term storage, and so on provider.

### [v0.3.2](https://github.com/layer5io/meshery/releases/tag/v0.3.2)

**What's new ?**
- Mesheryctl
    + Addition of `mesheryctl version` as a new subcommand.

### [v0.3.1](https://github.com/layer5io/meshery/releases/tag/v0.3.1)

**What's new ?**
- Meshery Server
    + Support for wrk2 as an alternative load generator.

### [v0.2.4](https://github.com/layer5io/meshery/releases/tag/v0.2.4)

**What's new ?**
- Mesheryctl
    + Homebrew support available for mesheryctl.
- Meshery Server
    + Meshery adapter for Octarine released as stable.
- Documentation
    + Revised quick start guide for Mac, Linux and Windows.
    + WSL2 support published.
    + GKE kubeconfig generation script switched to `--decode`.

### [v0.2.3](https://github.com/layer5io/meshery/releases/tag/v0.2.3)

**What's new ?**
- Mesheryctl
    + Improved `status` command output on Windows.
- Meshery Server
    + Ability to deploy Meshery on Istio.
    + Adapter Chips: Move adapter port number into tooltip.

### [v0.2.2](https://github.com/layer5io/meshery/releases/tag/v0.2.2)

**What's new ?**
- Mesheryctl
    + Improved `update` command, no overwriting in local configuration while `start` command.
- Meshery UI
    + A new modal view to organize and display performance results in a tabular format.

### [v0.2.1](https://github.com/layer5io/meshery/releases/tag/v0.2.1)

**What's new ?**
- Mesheryctl
    + Overhaul of in-cluster vs out-of-cluster Kubernetes setup.

## [v0.2.0](https://github.com/layer5io/meshery/releases/tag/v0.2.1)

**What's new?**
- Meshery Adapters
    + Adapter introduced for Network Service Mesh.
- Meshery Server
    + Ability to run performance tests asynchronously.
    + Collect & persist node metrics.

## [v0.1.6](https://github.com/layer5io/meshery/releases/tag/v0.1.6)

**What's new?**
- Changelog
 + New UI for managing Meshery's connection to Kubernetes cluster.
 + New `mesheryctl` compatibility for Windows for opening default browser upon `start`.
 + New sample app added to the istio adapter - Hipster app (Google's Microservices demo application).

## [v0.1.5](https://github.com/layer5io/meshery/releases/tag/v0.1.5)

**What's new?**
- Changelog
- UX Improvements
 + `mesheryctl start` now waits for meshery application containers to be up before launching the user's browser.
 + `mesheryctl stop` now shows command progress akin to the experience when using `meshery` bash script.

## [v0.1.4](https://github.com/layer5io/meshery/releases/tag/v0.1.4)

**What's new?**
- Changelog

## [v0.1.3](https://github.com/layer5io/meshery/releases/tag/v0.1.3)

**What's new?**
- Migrate from Configure Meshery to Settings page.

## [v0.1.2](https://github.com/layer5io/meshery/releases/tag/v0.1.2)

**What's new?**
- Synchronization of browser local storage
 + Meshery in-memory session storage.

## [v0.1.1](https://github.com/layer5io/meshery/releases/tag/v0.1.1)

**What's new?**
- Patch for Alpine bug.

## [v0.1.0](https://github.com/layer5io/meshery/releases/tag/v0.1.0)

**What's new**
- Initial fully functional Istio.

## [v0.0.9](https://github.com/layer5io/meshery/releases/tag/v0.0.9)

**What's new?**
- Documentation site segragated
 + content laid out.

## [v0.0.8](https://github.com/layer5io/meshery/releases/tag/v0.0.8)

**What's new?**
- Ability to import Grafana board json
 + Integrating with Prometheus directly for metrics.

## [v0.0.7](https://github.com/layer5io/meshery/releases/tag/v0.0.7)

**What's new?**
- Migrated away from embedded iframe grafana charts
 + Now using Chartjs for charting.
 + Now using C3 for charting.

## [v0.0.6](https://github.com/layer5io/meshery/releases/tag/v0.0.6)

**What's new?**
- Pre-alpha Consul adapter.
 + Support for deploying SMI Istio adapter.
 + Ensuring compatibility to apply SMI manifests through Meshery.

## [v0.0.5](https://github.com/layer5io/meshery/releases/tag/v0.0.5)

**What's new?**
- Pre-alpha Linkerd adapter.
 + Ability to filter results.

## [v0.0.4](https://github.com/layer5io/meshery/releases/tag/v0.0.4)

**What's new?**
-Ability to view persisted results.
 + Integration and support for Grafana charts.
 + Embedding panels in iframe.

## [v0.0.3](https://github.com/layer5io/meshery/releases/tag/v0.0.3)

**What's new?**
-  Initial version with support of Meshery adapters.
 + Release of a pre-alpha version of Istio adapter.

## [v0.0.2](https://github.com/layer5io/meshery/releases/tag/v0.0.2)

**What's new?**
- Ability to support running custom yaml on Kubernetes with Istio.

## [v0.0.1](https://github.com/layer5io/meshery/releases/tag/v0.0.1)

**New release**
- Initial version of Meshery
  + Connect to Kubernetes.
  + Run preconfigured commands on Kubernetes with Istio.

<!-- <table class="responsive-table hover striped">
  <thead>
    <th class="centered">Version</th>
    <th>Description</th>
    <th style="white-space: nowrap;">Release Date</th>
  </thead>
  <tbody>
      <tr>
        <td class="centered">0.3.15</td>
        <td>
          <em>Mesheryctl</em> -Feature: Addition of mesheryctl perf --file <smp.yaml> support for SMP formatted performance test profiles. Behavior changemesheryctl perf for using tokens and supporting JWT.
        </td>
        <td>May 1, 2020</td>
      </tr>
      <tr>
      <td class="centered">0.3.14</td>
      <td>
        <em>Mesheryctl</em> - Behavior change in mesheryctl start to check for new Meshery server images by default.
      </td>
      <td>Apr 24, 2020</td>
    </tr>
      <tr>
      <td class="centered">0.3.13</td>
      <td>
        <em>Mesheryctl</em> - Point release for mesheryctl's initial ARM support.
      </td>
      <td>Apr 16, 2020</td>
    </tr>
      <tr>
      <td class="centered">0.3.12</td>
      <td>
        <em>Mesheryctl</em> - Initial Scoop Bucket and Scoop package for mesheryctl.
      </td>
      <td>Apr 10, 2020</td>
    </tr>
      <tr>
      <td class="centered">0.3.11</td>
      <td>
        <em>Mesheryctl</em> - Fixes minor user experience issues on mesheryctl perf command. See Meshery CLI Commands & Documentation for reference.
      </td>
      <td>Mar 25, 2020</td>
    </tr>
      <tr>
      <td class="centered">0.3.10</td>
      <td>
        <em>Mesheryctl</em> - Introduces mesheryctl perf command. See Meshery CLI Commands & Documentation for reference.
      </td>
      <td>Mar 9, 2020</td>
    </tr>
      <tr>
      <td class="centered">0.3.9</td>
      <td>
        <em>Mesheryctl</em> - Addition of mesheryctl version to provide server-side version number, Improvement of mesheryctl logs | stop | start to provide appropriate grammar in situations when Meshery is stopped or Docker is not present.
        <em>Meshery</em> - Striped off extraneous information beyond IP address and port in Grafana and Prometheus endpoints. 
      </td>
      <td>Feb 3, 2020</td>
    </tr>
    <tr>
      <td class="centered">0.3.8</td>
      <td>
        <em>Mesheryctl</em> - Extraneous command line output removed. Clarity of CLI interaction with mesheryctl is improved in this release.
      </td>
      <td>Jan 17, 2020</td>
    </tr>
    <tr>
      <td class="centered">0.3.7</td>
      <td>
        <em>Meshery</em> - Ad-hoc connectivity tests for Prometheus is now supported. Users can click the Prometheus chip and have Meshery verify its ability to connect to the configured Prometheus instance.
      </td>
      <td>Jan 15, 2020</td>
    </tr>
     <tr>
      <td class="centered">0.3.6</td>
      <td>
        <em>Meshery</em> - Ad-hoc connectivity tests for Grafana is now supported. Users can click the Grafana chip and have Meshery verify its ability to connect to the configured Grafana instance.
      </td>
      <td>Jan 12, 2020</td>
    </tr>
     <tr>
      <td class="centered">0.3.5</td>
      <td>
        <em>Mesheryctl</em> - Removal of init as a command exposed to users. This command's functionality is used internal to mesheryctl start. A new start --check command will provide preflight check functionality in init's place.
      </td>
      <td>Jan 12, 2020</td>
    </tr>
     <tr>
      <td class="centered">0.3.4</td>
      <td>
        <em>Mesheryctl</em> - mesheryctl version is now enhanced with the addition of displaying the git commit (sha) of the mesheryctl release. 
      </td>
      <td>Dec 30, 2019</td>
    </tr>
     <tr>
      <td class="centered">0.3.3</td>
      <td>
        <em>Meshery</em> - Providers (a new project construct that allows users to select authentication, long-term storage, and so on provider). 
      </td>
      <td>Dec 20, 2019</td>
    </tr>
    <tr>
      <td class="centered">0.3.2</td>
      <td>
        <em>Mesheryctl</em> - adds mesheryctl version as a new subcommand. 
      </td>
      <td>Nov 29, 2019</td>
    </tr>
    <tr>
      <td class="centered">0.3.1</td>
      <td>
        <em>Meshery</em> - Support for wrk2 as an alternative load generator.
      </td>
      <td>Nov 12, 2019</td>
    </tr>
    <tr><td colspan="3"><strong>v0.3.0</strong></td></tr>
    <tr>
      <td class="centered">0.2.4</td>
      <td>
        <em>Meshery</em>
          - Meshery adapter for Octarine released as stable.
        <br />
        <em>mesheryctl</em>
          - now available through homebrew.
        <br />
        <em>Documentation</em>
        - revised quick start for Mac, Linux _and_ Windows.
        - WSL2 support published.
        - GKE kubeconfig generation script switched to `--decode`.
      </td>
      <td>Nov 5, 2019</td>
    </tr>
    <tr>
      <td class="centered">0.2.3</td>
      <td>
        <em>mesheryctl</em>
          - improved `status` output on Windows
        <br />
        <em>Meshery</em>
          - Ability to deploy Meshery on Istio. 
          - Adapter Chips: Move adapter port number into tooltip
        <br />
        <em>Docs</em>
          - /search no longer redirecting to github.io.
      </td>
      <td>Nov 3, 2019</td>
    </tr>
    <tr>
      <td class="centered">0.2.2</td>
      <td>Mesheryctl: improved verbosity of update command; no longer overwriting local .meshery.yml file when running start or logs. Docs: overhaul of docs site with a new jekyll theme (thanks @venilnoronha). Performance Testing: A new modal view to organize and display performance results in a tabular format.</td>
      <td>Oct 26, 2019</td>
    </tr>
    <tr>
      <td class="centered">0.2.1</td>
      <td>Meshery Installation: overhaul of in-cluster vs out-of-cluster Kubernetes setup.</td>
      <td>Oct 23, 2019</td>
    </tr>
    <tr>
      <td class="centered">0.2.0</td>
      <td>Meshery adapter for Network Service Mesh: adapter is now in beta; NSM provisioning. Performance Testing enhancements: performance tests run asynchronously, notifying the user of when test results are available; collect and persist node metrics. Service Mesh Sync: support for discovering service mesh type. Performance enhancements through memory tweaks and code profiling.</td>
      <td>Oct 22, 2019</td>
    </tr>
        <tr><td colspan="3"><strong>v0.2.0</strong></td></tr>
    <tr>
      <td class="centered">0.1.6</td>
      <td>New UI for managing Meshery's connection to Kubernetes cluster. New mesheryctl compatibility for Windows for opening default browser upon start.</td>
      <td>Oct 13, 2019</td>
    </tr>
    <tr>
      <td class="centered">0.1.5</td>
      <td>UX Improvements: mesheryctl start now waits for meshery application containers to be up before launching the user's browser. This new behavior ensures that users do not experience a 404 message; mesheryctl stop now shows command progress akin to the experience when using meshery bash script.</td>
      <td>Sep 20, 2019</td>
    </tr>
    <tr>
      <td class="centered">0.1.4</td>
      <td>Update README.md for the release.</td>
      <td>Sep 12, 2019</td>
    </tr>
    <tr>
      <td class="centered">0.1.3</td>
      <td>Migrate from Configure Meshery to Settings page.</td>
      <td>Jun 27, 2019</td>
    </tr>
    <tr>
      <td class="centered">0.1.2</td>
      <td>Synchronization of browser local storage with Meshery in-memory session storage.</td>
      <td>Jun 14, 2019</td>
    </tr>
    <tr>
      <td class="centered">0.1.1</td>
      <td>Patch for Alpine bug.</td>
      <td>May 31, 2019</td>
    </tr>
    <tr><td colspan="3"><strong>v0.1.0</strong></td></tr>
    <tr>
      <td class="centered">0.0.9</td>
      <td>Documentation site segregated and content laided out.</td>
      <td>May 2, 2019</td>
    </tr>
    <tr>
      <td class="centered">0.0.8</td>
      <td>Ability to import Grafana board json but integrating with Prometheus directly for metrics.</td>
      <td>Apr 15, 2019</td>
    </tr>
    <tr>
      <td class="centered">0.0.7</td>
      <td>Migrated away from embedded iframe grafana charts to using Chartjs and C3 for charting.</td>
      <td>Mar 20, 2019</td>
    </tr>
    <tr>
      <td class="centered">0.0.5</td>
      <td>Pre-alpha Linkerd adapter. Ability to filter results.</td>
      <td>Feb 28, 2019</td>
    </tr>
    <tr>
      <td class="centered">0.0.4</td>
      <td>Ability to view persisted results. Integration and support for Grafana charts and embedding panels in iframe.</td>
      <td>Feb 28, 2019</td>
    </tr>
    <tr>
      <td class="centered">0.0.3</td>
      <td>Initial version with support of Meshery adapters and release of a pre-alpha version of Istio adapter.</td>
      <td>Jan 21, 2019</td>
    </tr>
    <tr>
      <td class="centered">0.0.2</td>
      <td>Ability to support running custom yaml on Kubernetes with Istio.</td>
      <td>Nov 30, 2018</td>
    </tr>
    <tr>
      <td class="centered">0.0.1</td>
      <td>Initial version of Meshery. Connect to Kubernetes and run preconfigured commands on Kubernetes with Istio.</td>
      <td>Nov 16, 2018</td>
    </tr>
  </tbody>
</table>
 -->
