<img align="right" src="https://layer5.io/assets/images/cube-sh-small.png" />

# [Meshery](https://layer5.io/meshery)

A service mesh playground to faciliate learning about functionality and performance of different service meshes. [Meshery](https://layer5.io/meshery) incorporates the collection and display of metrics from applications running in the playground.

- [Functionality](#functionality)
- [Running Meshery](#running)
- [Architecture](https://docs.google.com/presentation/d/1UbuYMpn-e-mWVYwEASy4dzyZlrSgZX6MUfNtokraT9o/edit?usp=sharing)
  - [Design document](https://docs.google.com/document/d/1nV8TunLmVC8j5cBELT42YfEXYmhG3ZqFtHxeG3-w9t0/edit?usp=sharing)
- [Contributing](CONTRIBUTING.md/#contributing)
  - [Write an adapter](CONTRIBUTING.md/#adapter)
  - [Build the project](CONTRIBUTING.md/#building)
  
In an effort to produce service mesh agnostic tooling, Meshery uses a [common performance benchmark specification](https://github.com/layer5io/service-mesh-benchmark-spec) to capture and share environment information and test configuration. 

## Presentations
- Service Mesh Day - [Slides](https://docs.google.com/presentation/d/1T0w5sXiUYtjHhmwJYF7VI-q5lgYAN46-Yn8ey0EZV-A/edit?usp=sharing)
- Innotech San Antonio - [Slides](https://calcotestudios.com/talks/decks/slides-innotech-san-antonio-2019-the-enterprise-path-to-service-mesh.html)

## <a name="functionality">Functionality</a>
<img align="right" src="./ui/static/img/meshery.png?raw=true" alt="Service Mesh Playground" width="50%" />

### Multi-mesh Performance Benchmark

Meshery is intended to be a vendor and project-neutral utility for uniformly benchmarking the performance of service meshes. Between service mesh and proxy projects (and surprisingly, [within a single project](https://layer5.io/landscape#tools)), a number of different tools *and results* exist. 

### Multi-mesh Functionalty Playground

A service mesh playground to faciliate learning about functionality of different service meshes. Meshery incorporates a visual interface for manipulating traffic routing rules. Sample applications will be included in Meshery. 

## <a name="running">Running Meshery</a>

Please checkout [https://layer5.io/meshery](https://layer5.io/meshery) for quick start instructions and [https://meshery.layer5.io/docs](https://meshery.layer5.io/docs) for a complete set of Meshery documentation.

## License

This repository and site are available as open source under the terms of the [Apache 2.0 License](https://opensource.org/licenses/Apache-2.0).

## About Layer5
[Layer5.io](https://layer5.io) is a service mesh community, serving as a repository for information pertaining to the surrounding technology ecosystem (service meshes, api gateways, edge proxies, ingress and egress controllers) of microservice management in cloud native environments.
