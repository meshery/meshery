# Meshery Adapter for App Mesh

Source: /pr-preview/pr-21670/extensions/adapters/app-mesh/

### Features

1. Lifecycle management of App Mesh
1. Lifecycle management of sample applications
1. Performance testing

### Sample Applications

- Emojivoto

  - A microservice application that allows users to vote for their favorite emoji, and tracks votes received on a leaderboard.

- Bookinfo
  - The sample BookInfo application displays information about a book, similar to a single catalog entry of an online book store.
- Httpbin

  - Httpbin is a simple HTTP request and response service.

### Performance Management

Identify overhead involved in running the App Mesh, various App Mesh configurations while running different workloads and on different infrastructure. The adapter facilitates data plane and control plane performance testing.

1. Prometheus integration
1. Grafana integration

The [Meshery Adapter for App Mesh](https://github.com/meshery-extensions/meshery-app-mesh) will connect to App Mesh’s Prometheus and Grafana instances running in the control plane.
