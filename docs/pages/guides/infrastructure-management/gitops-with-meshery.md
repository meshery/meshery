---
layout: default
title: GitOps with Meshery
abstract: Integrating your CI/CD pipelines with Meshery's GitHub Actions
permalink: guides/infrastructure-management/gitops-with-meshery
type: guides
category: infrastructure
language: en
list: include
published: true
---

## <img style="height: 4rem; width: 4rem;" src="{{site.baseurl}}/assets/img/kanvas-icon-color.svg" /> Kanvas Snapshot GitHub Action

- See your deployment before you merge
- Connect Kanvas to your GitHub repo and see changes pull request-to-pull request
- Get snapshots of your infrastructure directly in your PRs

<h3>Using Meshery's SnapShot GitHub Action</h3>
See your deployment before you merge. Connect Kanvas to your GitHub repo and see changes pull request-to-pull request. Get snapshots of your infrastructure directly in your PRs.

See [Extension: Kanvas Snapshot](/extensions/snapshot) for more details.

## <img src="{{ site.baseurl }}/assets/img/smp-dark-text-side.svg" className="image-left-no-shadow" alt="Infrastructure fperformance logo" style="width: 60%; max-width: 200px;vertical-align:middle" /> Meshery Performance Analysis GitHub Action

- See your performance regressions before you merge
- Connect Performance Management to your GitHub repo and see changes pull request-to-pull request
- Red light performance regressions
- Baseline and analyze the performance of your services is key to efficient operation of any application
  - Meshery is the canonical implementation of the Cloud Native Performance specification
- Define your performance profiles upfront. See statistical analysis with microservice latency and throughput quartiles
- Meshery includes your choice of load generator, so that you can meausure your way
- Meshery packages all these features into an easy-to-use GitHub Action

Measuring and managing the performance of your infrastructure is key to efficient operation. You can choose from multiple load generators and use a highly configurable set of load profiles with variable tunable facets to run a performance test. Meshery packages all these features into an easy-to-use GitHub Action.

<h3>Using Meshery's Performance Analysis GitHub Action</h3>

The <a href="https://github.com/marketplace/actions/performance-testing-with-meshery">Meshery Performance Analysis GitHub Action</a> is available in the GitHub Marketplace. You can create your own performance profiles to run repeatable tests with Meshery. You can configure this action to trigger with each of your releases, on every pull request or any GitHub workflow trigger event. A sample configuration of the action is shown below.

```yaml
name: Meshery SMP Action
on:
  push:
    branches:
      'master'

jobs:
  performance-test:
    name: Performance Test
    runs-on: ubuntu-latest
    steps:
      - name: checkout
        uses: actions/checkout@v2
        with:
          ref: 'perf'

      - name: Deploy k8s-minikube
        uses: manusa/actions-setup-minikube@v2.4.1
        with:
          minikube version: 'v1.21.0'
          kubernetes version: 'v1.20.7'
          driver: docker

      - name: Run Performance Test
        uses: layer5io/meshery-smp-action@master
        with:
          provider_token: ${{ secrets.PROVIDER_TOKEN }}
          platform: docker
          profile_name: soak-test
```

You can also define your test configuration in an SMP compatible configuration file as shown below.

```yaml
smp_version: v0.0.1
id:
name: Istio Performance Test
labels: {}
clients:
- internal: false
  load_generator: fortio
  protocol: 1
  connections: 2
  rps: 10
  headers: {}
  cookies: {}
  body: ""
  content_type: ""
  endpoint_urls:
  - http://localhost:2323/productpage
duration: "30m"
```

See this sample GitHub workflow (<a href="https://github.com/layer5io/meshery-smp-action/blob/master/action.yml">action.yml</a>) for more configuration details.

<a href="{{ site.baseurl }}/assets/img/service-mesh-performance-profile-test-results.png"><img src="{{ site.baseurl }}/assets/img/service-mesh-performance-profile-test-results.png" className="image-center" alt="performance management dashboard" /></a>

The results from the tests are updated on the Performance Management dashboard in Meshery. To learn more about interpreting the test results, check out <a href="https://docs.meshery.io/guides/performance-management/interpreting-performance-test-results">this guide</a>. You can always checkout the <a href="https://docs.meshery.io/guides">Meshery User Guides</a> to dive deep into these features.

<a href="https://smp-spec.io">Cloud Native Performance</a> standardizes infrastructure measurement, characterizing any deployment's performance by capturing the details of infrastructure capacity, configuration and workload metadata.