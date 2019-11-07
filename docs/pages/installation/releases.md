---
layout: page
title: Releases
permalink: /releases
---

A list of the [releases of Meshery](https://github.com/layer5io/meshery/releases). See the [Build and Release Strategy](https://docs.google.com/document/d/11nAxYtz2SUusCYZ0JeNRrOLIxkgmmbUVWz63MBZV2oE/edit?usp=sharing) document for details.

<table class="responsive-table hover striped">
  <thead>
    <th class="centered">Version</th>
    <th>Description</th>
    <th>Release Date</th>
  </thead>
  <tbody>
    <tr>
      <td class="centered">0.2.4</td>
      <td>

*Meshery*
  - Meshery adapter for Octarine released as stable.

*mesheryctl*
  - now available through homebrew.

*Documentation*
 - revised quick start for Mac, Linux _and_ Windows.
 - WSL2 support published.
 - GKE kubeconfig generation script switched to `--decode`.
       </td>
       <td>Nov 6, 2019</td>
    </tr>
    <tr>
      <td class="centered">0.2.3</td>
      <td>

_Mesheryctl_
  - improved `status` output on Windows

_Meshery_
  - Ability to deploy Meshery on Istio. 
  - Adapter Chips: Move adapter port number into tooltip

_Docs_
   - /search no longer redirecting to github.io.
      </td>
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
    <tr>
      <td class="centered">0.0.9</td>
      <td>Documentation site segragated and content laided out.</td>
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
